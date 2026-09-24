package cli

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/compozy/kb/internal/actions"
	kconfig "github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/firecrawl"
	"github.com/compozy/kb/internal/gate"
	kingest "github.com/compozy/kb/internal/ingest"
	"github.com/compozy/kb/internal/mediadl"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/review"
	"github.com/compozy/kb/internal/session"
	"github.com/compozy/kb/internal/youtube"
)

// ingestFlags are the decision flags shared by every gated ingest command.
type ingestFlags struct {
	Session session.Flags
	// Force skips gates 1–6 (classify and link still run).
	Force bool
	// Rescue ingests one skipped item by its skipped.jsonl id (bulk
	// commands only).
	Rescue string
}

// bindIngestFlags registers --budget, --decisions and --force, plus
// --rescue on bulk commands.
func bindIngestFlags(command *cobra.Command, flags *ingestFlags, bulk bool) {
	bindDecisionFlags(command, &flags.Session)
	command.Flags().BoolVar(&flags.Force, "force", false, "Skip the ingest gates (dedupe, relevance, quality, near-duplicate); classify and link still run")
	if bulk {
		command.Flags().StringVar(&flags.Rescue, "rescue", "", "Fetch and ingest one item skipped before fetch, by its skipped.jsonl id (bypasses only the pre-fetch gate)")
	}
}

// ingestRunner is the decision-backed run an ingest command drives
// (implemented by *ingest.Run; tests replace openIngestRunner).
type ingestRunner interface {
	Session() *session.Session
	Precheck(sourceURL, title string, kind models.SourceKind) (kingest.Result, bool, error)
	Prefetch(ctx context.Context, items []gate.Item) ([]gate.PrefetchDecision, error)
	Ingest(ctx context.Context, options kingest.Options) (kingest.Result, error)
	Skipped(id string) (gate.SkippedRow, error)
	MarkRescued(row gate.SkippedRow) error
	Finish(ctx context.Context) (kingest.Summary, error)
}

// openIngestRunner opens the session (the decision model is required) and
// starts one gated run.
var openIngestRunner = func(cmd *cobra.Command, command, topicSlug string, flags ingestFlags, batch, query string) (ingestRunner, error) {
	s, err := openSession(cmd, "kb ingest "+command, topicSlug, flags.Session)
	if err != nil {
		return nil, err
	}
	return kingest.NewRun(s, kingest.RunOptions{Command: command, Batch: batch, Query: query, Force: flags.Force}), nil
}

// finishIngestRun classifies and links what the run wrote, appends the run's
// log.md entry and prints the run summary (gate modes, counts, decisions) to
// stderr.
func finishIngestRun(cmd *cobra.Command, runner ingestRunner) error {
	summary, err := runner.Finish(commandContext(cmd))
	writeIngestSummary(cmd.ErrOrStderr(), runner.Session(), summary)
	return err
}

func writeIngestSummary(w io.Writer, s *session.Session, summary kingest.Summary) {
	for _, line := range summary.Lines() {
		_, _ = fmt.Fprintln(w, line)
	}
	if s != nil {
		s.WriteSummary(w)
	}
}

// scrapeRefetch returns the stage-4 refetch of a URL: one fresh scrape with
// the configured freshness options.
func scrapeRefetch(scraper firecrawlScraper, cfg kconfig.FirecrawlConfig, sourceURL string) func(context.Context) (*firecrawl.ScrapeResult, error) {
	return func(ctx context.Context) (*firecrawl.ScrapeResult, error) {
		return scraper.ScrapeWithOptions(ctx, sourceURL, firecrawl.RefetchOptions(cfg))
	}
}

// ingestScraped ingests one scraped web page through the run.
func ingestScraped(
	ctx context.Context,
	runner ingestRunner,
	target ingestTarget,
	scraper firecrawlScraper,
	cfg kconfig.Config,
	requested string,
	scraped *firecrawl.ScrapeResult,
	options kingest.Options,
) (kingest.Result, error) {
	sourceURL := strings.TrimSpace(scraped.SourceURL)
	if sourceURL == "" {
		sourceURL = requested
	}
	options.VaultPath = target.VaultPath
	options.Topic = target.TopicInfo.Slug
	options.SourceKind = models.SourceKindArticle
	options.SourceURL = sourceURL
	options.Title = scraped.Title
	options.Markdown = scraped.Markdown
	options.Gate.RequestedURL = requested
	options.Gate.FinalURL = scraped.FinalURL
	options.Gate.StatusCode = scraped.StatusCode
	options.Gate.SiteName = scraped.SiteName
	options.Gate.Refetch = scrapeRefetch(scraper, cfg.Firecrawl, requested)
	return runner.Ingest(ctx, options)
}

// rescueFetchers are the extractors a rescue may need.
type rescueFetchers struct {
	cfg     kconfig.Config
	scraper firecrawlScraper
}

// rescueSkipped fetches and ingests one skipped item, bypassing only the
// pre-fetch gate: YouTube URLs go through the transcript extractor,
// Instagram URLs through the Instagram extractor, everything else through
// Firecrawl. The item keeps its original batch and query, and the
// skipped.jsonl row is marked rescued.
func rescueSkipped(ctx context.Context, runner ingestRunner, target ingestTarget, fetchers rescueFetchers, row gate.SkippedRow) (kingest.Result, error) {
	options := kingest.Options{VaultPath: target.VaultPath, Topic: target.TopicInfo.Slug, Batch: row.Batch, Query: row.Query}
	var result kingest.Result
	var err error
	switch platform := gate.PlatformID(row.URL); {
	case strings.HasPrefix(platform, "youtube:"):
		result, err = rescueYouTube(ctx, runner, fetchers.cfg, row, options)
	case strings.HasPrefix(platform, "instagram:"):
		result, err = rescueInstagram(ctx, runner, fetchers.cfg, row, options)
	default:
		var scraped *firecrawl.ScrapeResult
		scraped, err = fetchers.scraper.Scrape(ctx, row.URL)
		if err == nil {
			result, err = ingestScraped(ctx, runner, target, fetchers.scraper, fetchers.cfg, row.URL, scraped, options)
		}
	}
	if err != nil {
		return result, fmt.Errorf("rescue %s: %w", row.ID, err)
	}
	if err := runner.MarkRescued(row); err != nil {
		return result, err
	}
	return result, nil
}

func rescueYouTube(ctx context.Context, runner ingestRunner, cfg kconfig.Config, row gate.SkippedRow, options kingest.Options) (kingest.Result, error) {
	policy, err := youtube.ParseTranscriptionPolicy(cfg.YouTube.Transcription)
	if err != nil {
		return kingest.Result{}, err
	}
	languages := normalizeYouTubeCaptionLanguages(cfg.YouTube.CaptionLanguages)
	if len(languages) == 0 {
		languages = []string{"orig"}
	}
	extracted, err := newYouTubeTranscriptExtractor(cfg).Extract(ctx, row.URL, youtube.ExtractOptions{
		TranscriptionPolicy: policy, PreferredLanguages: languages, AllowTranslatedCaptions: cfg.YouTube.AllowTranslatedCaptions,
	})
	if err != nil {
		return kingest.Result{}, err
	}
	options.SourceKind = models.SourceKindYouTubeTranscript
	options.SourceURL = firstNonBlank(extracted.Metadata.URL, row.URL)
	options.Title = firstNonBlank(extracted.Metadata.Title, row.Title)
	options.Markdown = extracted.Markdown
	options.ExtraFrontmatter = youtubeFrontmatter(extracted)
	return runner.Ingest(ctx, options)
}

func rescueInstagram(ctx context.Context, runner ingestRunner, cfg kconfig.Config, row gate.SkippedRow, options kingest.Options) (kingest.Result, error) {
	policy, err := mediadl.ParseTranscriptionPolicy(cfg.Instagram.Transcription)
	if err != nil {
		return kingest.Result{}, err
	}
	extracted, err := newInstagramTranscriptExtractor(cfg).Extract(ctx, row.URL, mediadl.ExtractOptions{TranscriptionPolicy: policy})
	if err != nil {
		return kingest.Result{}, err
	}
	options.SourceKind = models.SourceKindInstagramVideo
	options.SourceURL = firstNonBlank(extracted.Metadata.URL, row.URL)
	options.Title = firstNonBlank(extracted.Metadata.Title, row.Title)
	options.Markdown = extracted.Markdown
	options.ExtraFrontmatter = instagramFrontmatter(extracted)
	return runner.Ingest(ctx, options)
}

// runIngestRescue is `--rescue <id>`: when the skip review item is pending
// it is accepted through the review actions (label + resolve); otherwise
// the item is rescued directly. It prints the ingest result as JSON.
func runIngestRescue(cmd *cobra.Command, runner ingestRunner, target ingestTarget, fetchers rescueFetchers, id string) error {
	ctx := commandContext(cmd)
	row, err := runner.Skipped(id)
	if err != nil {
		return err
	}
	var result kingest.Result
	rescue := func(ctx context.Context, row gate.SkippedRow) (string, error) {
		var err error
		result, err = rescueSkipped(ctx, runner, target, fetchers, row)
		return result.FilePath, err
	}
	item, pending, err := pendingSkipItem(runner.Session(), row.ID)
	if err != nil {
		return err
	}
	if pending {
		if _, err := actions.Apply(ctx, runner.Session(), actions.Deps{Rescue: rescue}, item, actions.Accept); err != nil {
			return err
		}
	} else if _, err := rescue(ctx, row); err != nil {
		return err
	}
	if err := writeJSON(cmd, result); err != nil {
		return err
	}
	return finishIngestRun(cmd, runner)
}

// pendingSkipItem finds the pending `skip` review item of a skipped id.
func pendingSkipItem(s *session.Session, id string) (review.Item, bool, error) {
	if s == nil {
		return review.Item{}, false, nil
	}
	items, err := review.Open(s.Root(), s.Now).Pending(review.QueueSkip)
	if err != nil {
		return review.Item{}, false, err
	}
	for _, item := range items {
		if item.Subject == id {
			return item, true, nil
		}
	}
	return review.Item{}, false, nil
}

func firstNonBlank(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
