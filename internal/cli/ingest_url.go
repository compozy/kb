package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	kconfig "github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/firecrawl"
	"github.com/compozy/kb/internal/gate"
	kingest "github.com/compozy/kb/internal/ingest"
	"github.com/compozy/kb/internal/models"
)

type firecrawlScraper interface {
	Scrape(ctx context.Context, sourceURL string) (*firecrawl.ScrapeResult, error)
	ScrapeWithOptions(ctx context.Context, sourceURL string, opts firecrawl.ScrapeOptions) (*firecrawl.ScrapeResult, error)
}

var newFirecrawlScraper = func(cfg firecrawlConfig) firecrawlScraper {
	return firecrawl.NewClient(cfg)
}

type firecrawlConfig = kconfig.FirecrawlConfig

// urlListEntry is one URL of a list, with its optional label.
type urlListEntry struct {
	URL   string
	Label string
}

// urlListSummary is the JSON result of a URL-list ingest.
type urlListSummary struct {
	Topic    string           `json:"topic"`
	Batch    string           `json:"batch"`
	Results  []kingest.Result `json:"results"`
	Skipped  []urlListSkip    `json:"skipped"`
	Failures []urlListFailure `json:"failures"`
}

type urlListSkip struct {
	ID        string  `json:"id,omitempty"`
	URL       string  `json:"url"`
	Title     string  `json:"title,omitempty"`
	Reason    string  `json:"reason"`
	POffTopic float64 `json:"p_off_topic,omitempty"`
	Detail    string  `json:"detail,omitempty"`
}

type urlListFailure struct {
	URL   string `json:"url"`
	Error string `json:"error"`
}

func newIngestURLCommand() *cobra.Command {
	var topic string
	var batch string
	var from string
	var flags ingestFlags

	command := &cobra.Command{
		Use:   "url <url>...",
		Short: "Scrape one or more web URLs and ingest them into a topic",
		Long: "Scrape web URLs through Firecrawl and ingest them through the ingest gates.\n" +
			"Several URLs, or --from <file> (one URL per line, optional label after whitespace),\n" +
			"make a bulk run: items are judged against the selection contract before fetch.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runIngestURLCommand(cmd, args, topic, batch, from, flags)
		},
	}

	requireTopicFlag(command, &topic)
	addBatchFlag(command, &batch)
	command.Flags().StringVar(&from, "from", "", "File with one URL per line (optional label after whitespace; # comments)")
	bindIngestFlags(command, &flags, true)

	return command
}

func runIngestURLCommand(cmd *cobra.Command, args []string, topicSlug, batch, from string, flags ingestFlags) error {
	const action = "ingest url"
	if len(args) == 0 && strings.TrimSpace(from) == "" && strings.TrimSpace(flags.Rescue) == "" {
		return fmt.Errorf("%s: requires at least one URL, --from <file> or --rescue <id>", action)
	}
	target, err := resolveIngestTarget(cmd, action, topicSlug)
	if err != nil {
		return err
	}
	cfg, err := loadIngestConfig()
	if err != nil {
		return fmt.Errorf("%s: %w", action, err)
	}
	entries := make([]urlListEntry, 0, len(args))
	for _, arg := range args {
		entries = append(entries, urlListEntry{URL: strings.TrimSpace(arg)})
	}
	query := ""
	if strings.TrimSpace(from) != "" {
		listed, err := readURLList(from)
		if err != nil {
			return fmt.Errorf("%s: %w", action, err)
		}
		entries = append(entries, listed...)
		query = filepath.Base(from)
	}

	runBatch := resolveIngestBatch("url", batch)
	runner, err := openIngestRunner(cmd, "url", target.TopicInfo.Slug, flags, runBatch, query)
	if err != nil {
		return err
	}
	scraper := newFirecrawlScraper(firecrawlConfig(cfg.Firecrawl))
	if rescue := strings.TrimSpace(flags.Rescue); rescue != "" {
		return runIngestRescue(cmd, runner, target, rescueFetchers{cfg: cfg, scraper: scraper}, rescue)
	}
	if len(entries) == 1 && strings.TrimSpace(from) == "" {
		return runSingleURL(cmd, runner, target, scraper, cfg, runBatch, entries[0].URL)
	}
	return runURLList(cmd, runner, target, scraper, cfg, runBatch, entries)
}

func runSingleURL(cmd *cobra.Command, runner ingestRunner, target ingestTarget, scraper firecrawlScraper, cfg kconfig.Config, batch, sourceURL string) error {
	ctx := commandContext(cmd)
	result, skipped, err := runner.Precheck(sourceURL, "", models.SourceKindArticle)
	if err != nil {
		return fmt.Errorf("ingest url: %w", err)
	}
	if !skipped {
		scraped, err := scraper.Scrape(ctx, sourceURL)
		if err != nil {
			return fmt.Errorf("ingest url: %w", err)
		}
		result, err = ingestScraped(ctx, runner, target, scraper, cfg, sourceURL, scraped, kingest.Options{Batch: batch})
		if err != nil {
			return fmt.Errorf("ingest url: %w", err)
		}
	}
	if err := writeJSON(cmd, result); err != nil {
		return err
	}
	return finishIngestRun(cmd, runner)
}

// runURLList is the bulk path: stage 1 on every URL, stage 2 over the
// rest (label as title, host as site), then fetch and ingest the survivors.
func runURLList(cmd *cobra.Command, runner ingestRunner, target ingestTarget, scraper firecrawlScraper, cfg kconfig.Config, batch string, entries []urlListEntry) error {
	ctx := commandContext(cmd)
	summary := urlListSummary{Topic: target.TopicInfo.Slug, Batch: batch, Results: []kingest.Result{}, Skipped: []urlListSkip{}, Failures: []urlListFailure{}}
	items := make([]gate.Item, 0, len(entries))
	labels := map[string]string{}
	for _, entry := range entries {
		result, skipped, err := runner.Precheck(entry.URL, entry.Label, models.SourceKindArticle)
		if err != nil {
			return fmt.Errorf("ingest url: %w", err)
		}
		if skipped {
			summary.Skipped = append(summary.Skipped, urlListSkip{URL: entry.URL, Title: entry.Label, Reason: result.TriageReason, Detail: "duplicate of " + result.DuplicateOf})
			continue
		}
		title := firstNonBlank(entry.Label, titleHintFromURL(entry.URL))
		labels[entry.URL] = entry.Label
		items = append(items, gate.Item{
			URL: entry.URL, Title: title, Channel: urlHost(entry.URL),
			Path: kingest.WouldBePath(models.SourceKindArticle, title), SourceKind: string(models.SourceKindArticle),
		})
	}
	decided, err := runner.Prefetch(ctx, items)
	if err != nil {
		return fmt.Errorf("ingest url: %w", err)
	}
	for index := range decided {
		decision := decided[index]
		item := decision.Item
		if decision.Action == gate.ActionSkip {
			summary.Skipped = append(summary.Skipped, urlListSkip{
				ID: item.ID, URL: item.URL, Title: labels[item.URL], Reason: gate.ReasonOffTopic, POffTopic: decision.POffTopic,
			})
			continue
		}
		scraped, err := scraper.Scrape(ctx, item.URL)
		if err != nil {
			summary.Failures = append(summary.Failures, urlListFailure{URL: item.URL, Error: err.Error()})
			continue
		}
		options := kingest.Options{Batch: batch, Gate: kingest.GateOptions{Prefetch: &decision}}
		if label := labels[item.URL]; label != "" {
			options.Query = label
		}
		result, err := ingestScraped(ctx, runner, target, scraper, cfg, item.URL, scraped, options)
		if err != nil {
			// Ingest errors are fatal (auth, cancellation, I/O): stop, but
			// still classify, link and log what the run already wrote.
			return errors.Join(fmt.Errorf("ingest url: %s: %w", item.URL, err), finishIngestRun(cmd, runner))
		}
		summary.Results = append(summary.Results, result)
	}
	if err := writeJSON(cmd, summary); err != nil {
		return err
	}
	return finishIngestRun(cmd, runner)
}

// readURLList parses a URL list file: one URL per line, an optional label
// after the first whitespace, blank lines and `#` comments skipped.
func readURLList(path string) ([]urlListEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open URL list: %w", err)
	}
	defer func() { _ = file.Close() }()
	entries := make([]urlListEntry, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		entry := urlListEntry{URL: fields[0]}
		if len(fields) > 1 {
			entry.Label = strings.TrimSpace(strings.TrimPrefix(line, fields[0]))
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read URL list: %w", err)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("URL list %s has no URLs", path)
	}
	return entries, nil
}

// titleHintFromURL is the last path segment of a URL, humanized, for the
// pre-fetch state of an unlabeled URL.
func titleHintFromURL(raw string) string {
	trimmed := strings.TrimRight(strings.SplitN(strings.SplitN(raw, "?", 2)[0], "#", 2)[0], "/")
	index := strings.LastIndex(trimmed, "/")
	if index < 0 || strings.HasSuffix(trimmed[:index], "/") {
		return ""
	}
	segment := strings.TrimSuffix(trimmed[index+1:], filepath.Ext(trimmed[index+1:]))
	return strings.TrimSpace(strings.NewReplacer("-", " ", "_", " ").Replace(segment))
}

func urlHost(raw string) string {
	key := gate.NormalizeURL(raw)
	host, _, _ := strings.Cut(key, "/")
	host, _, _ = strings.Cut(host, "?")
	return host
}
