package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	kconfig "github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/gate"
	kingest "github.com/compozy/kb/internal/ingest"
	"github.com/compozy/kb/internal/models"
	ktopic "github.com/compozy/kb/internal/topic"
	"github.com/compozy/kb/internal/vault"
	"github.com/compozy/kb/internal/youtube"
)

type youtubeChannelExtractor interface {
	ListChannel(ctx context.Context, normalizedURL string, limit int) (youtube.ChannelListing, error)
	BulkExtract(
		ctx context.Context,
		videos []youtube.ChannelVideo,
		options youtube.BulkOptions,
		sink func(youtube.VideoOutcome),
	) error
}

var newYouTubeChannelExtractor = func(cfg kconfig.Config) youtubeChannelExtractor {
	return youtube.NewExtractorWithConfig(cfg.STT, cfg.OpenRouter, cfg.YouTube)
}

type channelVideoSummary struct {
	VideoID  string `json:"video_id"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	FilePath string `json:"file_path,omitempty"`
	Error    string `json:"error,omitempty"`
	// Triage and Reason report the gate outcome (kept, review,
	// quarantined; skipped / duplicate-skipped for items not fetched).
	Triage string `json:"triage,omitempty"`
	Reason string `json:"reason,omitempty"`
	// SkippedID is the skipped.jsonl id of an item skipped before fetch.
	SkippedID string `json:"skipped_id,omitempty"`
}

type channelIngestSummary struct {
	Topic                string                `json:"topic"`
	ChannelURL           string                `json:"channel_url"`
	NormalizedChannelURL string                `json:"normalized_channel_url"`
	Selection            string                `json:"selection"`
	Transcribe           string                `json:"transcribe"`
	CaptionLanguages     []string              `json:"caption_languages"`
	Resolved             int                   `json:"resolved"`
	DryRun               bool                  `json:"dry_run"`
	Videos               []channelVideoSummary `json:"videos"`
	Ingested             []channelVideoSummary `json:"ingested"`
	Skipped              []channelVideoSummary `json:"skipped"`
	Failures             []channelVideoSummary `json:"failures"`
}

type ingestChannelCommandOptions struct {
	topicSlug   string
	transcribe  string
	limit       int
	all         bool
	concurrency int
	throttle    time.Duration
	dryRun      bool
	createTopic bool
	subLangs    string
	lang        string
	batch       string
	flags       ingestFlags
}

// channelProvenance is the ingest_batch/ingest_query pair shared by every
// video of one channel run: one batch, and the channel/playlist URL argument
// as the query.
type channelProvenance struct {
	batch string
	query string
}

func newIngestChannelCommand() *cobra.Command {
	var topicSlug string
	var transcribe string
	var limit int
	var all bool
	var concurrency int
	var throttle time.Duration
	var dryRun bool
	var createTopic bool
	var subLangs string
	var lang string
	var batch string
	var flags ingestFlags

	command := &cobra.Command{
		Use:   "channel <url>",
		Short: "Bulk-extract YouTube channel or playlist transcripts into a topic",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options := ingestChannelCommandOptions{
				topicSlug:   topicSlug,
				transcribe:  transcribe,
				limit:       limit,
				all:         all,
				concurrency: concurrency,
				throttle:    throttle,
				dryRun:      dryRun,
				createTopic: createTopic,
				subLangs:    subLangs,
				lang:        lang,
				batch:       batch,
				flags:       flags,
			}
			if strings.TrimSpace(flags.Rescue) != "" {
				return runIngestChannelRescue(cmd, options)
			}
			if len(args) != 1 {
				return fmt.Errorf("ingest channel: requires a channel or playlist URL (or --rescue <id>)")
			}
			return runIngestChannelCommand(cmd, args[0], options)
		},
	}

	requireTopicFlag(command, &topicSlug)
	command.Flags().StringVar(&transcribe, "transcribe", "", "Transcription policy: captions, auto, or stt (default captions)")
	command.Flags().IntVar(&limit, "limit", 0, "Maximum newest uploads to ingest (0 = all)")
	command.Flags().BoolVar(&all, "all", false, "Ingest all uploads (overrides --limit)")
	command.Flags().IntVar(&concurrency, "concurrency", 0, "Concurrent transcript fetches (default from [youtube].bulk_concurrency)")
	command.Flags().DurationVar(&throttle, "throttle", 0, "Delay between transcript fetches (default from [youtube].bulk_throttle)")
	command.Flags().StringVar(&subLangs, "sub-langs", "", "Caption languages to request, comma-separated; use orig for the video's original language")
	command.Flags().StringVar(&lang, "lang", "", "Alias for --sub-langs")
	command.Flags().BoolVar(&dryRun, "dry-run", false, "Resolve and list videos without ingesting")
	command.Flags().BoolVar(&createTopic, "create-topic", true, "Create the target topic when it does not exist")
	addBatchFlag(command, &batch)
	bindIngestFlags(command, &flags, true)

	return command
}

func runIngestChannelCommand(cmd *cobra.Command, channelURL string, options ingestChannelCommandOptions) error {
	const action = "ingest channel"

	vaultPath, err := resolveCommandVaultPath(cmd, ingestGetwd, action)
	if err != nil {
		return err
	}
	cfg, err := loadIngestConfig()
	if err != nil {
		return fmt.Errorf("%s: %w", action, err)
	}

	policyValue := strings.TrimSpace(cfg.YouTube.Transcription)
	if cmd.Flags().Changed("transcribe") {
		policyValue = options.transcribe
	}
	policy, err := youtube.ParseTranscriptionPolicy(policyValue)
	if err != nil {
		return fmt.Errorf("%s: %w", action, err)
	}
	captionLanguages, err := resolveYouTubeCaptionLanguages(cmd, cfg.YouTube, options.subLangs, options.lang)
	if err != nil {
		return fmt.Errorf("%s: %w", action, err)
	}
	topicRef, err := ktopic.ParseTopicRef(options.topicSlug)
	if err != nil {
		return fmt.Errorf("%s: %w", action, err)
	}

	normalizedURL, err := youtube.NormalizeChannelURL(channelURL)
	if err != nil {
		return fmt.Errorf("%s: %w", action, err)
	}

	resolvedLimit := options.limit
	if options.all {
		resolvedLimit = 0
	}

	ctx := commandContext(cmd)
	extractor := newYouTubeChannelExtractor(cfg)
	listing, err := extractor.ListChannel(ctx, normalizedURL, resolvedLimit)
	if err != nil {
		return fmt.Errorf("%s: %w", action, err)
	}

	summaryTarget := ingestTarget{
		TopicInfo: models.TopicInfo{Slug: topicRef.Path},
		VaultPath: vaultPath,
	}
	summary := newChannelIngestSummary(
		summaryTarget,
		channelURL,
		normalizedURL,
		policy,
		captionLanguages,
		resolvedLimit,
		options.dryRun,
		listing.Videos,
	)
	if options.dryRun {
		return writeJSON(cmd, summary)
	}

	target, err := resolveChannelIngestTarget(
		action,
		vaultPath,
		topicRef,
		deriveChannelTopicTitle(listing, topicRef),
		options.createTopic,
	)
	if err != nil {
		return err
	}
	summary.Topic = target.TopicInfo.Slug

	provenance := channelProvenance{
		batch: resolveIngestBatch("channel", options.batch),
		query: strings.TrimSpace(channelURL),
	}
	runner, err := openIngestRunner(cmd, "channel", target.TopicInfo.Slug, options.flags, provenance.batch, provenance.query)
	if err != nil {
		return err
	}

	existing := map[string]struct{}{}
	if !options.flags.Force {
		existing, err = existingYouTubeVideoIDs(target.VaultPath, target.TopicInfo.Slug)
		if err != nil {
			return fmt.Errorf("%s: %w", action, err)
		}
	}
	candidates := newChannelVideosToFetch(listing.Videos, existing, &summary)
	decisions, err := prefetchChannelVideos(ctx, runner, listing, candidates, &summary)
	if err != nil {
		return fmt.Errorf("%s: %w", action, err)
	}
	toFetch := make([]youtube.ChannelVideo, 0, len(decisions))
	for _, video := range candidates {
		if decision, ok := decisions[video.VideoID]; ok && decision.Action != gate.ActionSkip {
			toFetch = append(toFetch, video)
		}
	}

	bulkOptions, err := resolveChannelBulkOptions(cmd, cfg.YouTube, policy, captionLanguages, options.concurrency, options.throttle)
	if err != nil {
		return fmt.Errorf("%s: %w", action, err)
	}

	var fatal error
	bulkErr := extractor.BulkExtract(ctx, toFetch, bulkOptions, func(outcome youtube.VideoOutcome) {
		if fatal != nil {
			return
		}
		decision := decisions[outcome.Video.VideoID]
		fatal = recordChannelOutcome(ctx, runner, target, provenance, outcome, &decision, &summary)
	})
	if writeErr := writeJSON(cmd, summary); writeErr != nil {
		return writeErr
	}
	finishErr := finishIngestRun(cmd, runner)
	if fatal != nil {
		return errors.Join(fmt.Errorf("%s: %w", action, fatal), finishErr)
	}
	if bulkErr != nil {
		return errors.Join(fmt.Errorf("%s: %w", action, bulkErr), finishErr)
	}
	return finishErr
}

// prefetchChannelVideos runs stage 1 (by video id, quarantined videos
// included) and stage 2 (title + channel against the contract) over the
// candidate videos. Flat channel listings carry no description or date, so
// the pre-fetch state holds the title, channel and provenance only.
func prefetchChannelVideos(
	ctx context.Context,
	runner ingestRunner,
	listing youtube.ChannelListing,
	candidates []youtube.ChannelVideo,
	summary *channelIngestSummary,
) (map[string]gate.PrefetchDecision, error) {
	channel := firstNonBlank(listing.Channel, listing.Uploader, listing.Title)
	items := make([]gate.Item, 0, len(candidates))
	videos := make([]youtube.ChannelVideo, 0, len(candidates))
	for _, video := range candidates {
		result, skipped, err := runner.Precheck(video.URL, video.Title, models.SourceKindYouTubeTranscript)
		if err != nil {
			return nil, err
		}
		if skipped {
			entry := summarizeVideo(video)
			entry.Triage, entry.Reason = result.Triage, result.TriageReason
			summary.Skipped = append(summary.Skipped, entry)
			continue
		}
		videos = append(videos, video)
		items = append(items, gate.Item{
			URL: video.URL, Title: video.Title, Channel: channel,
			Path:       kingest.WouldBePath(models.SourceKindYouTubeTranscript, video.Title),
			SourceKind: string(models.SourceKindYouTubeTranscript),
		})
	}
	decided, err := runner.Prefetch(ctx, items)
	if err != nil {
		return nil, err
	}
	byVideo := make(map[string]gate.PrefetchDecision, len(decided))
	for index, decision := range decided {
		video := videos[index]
		byVideo[video.VideoID] = decision
		if decision.Action == gate.ActionSkip {
			entry := summarizeVideo(video)
			entry.Triage, entry.Reason, entry.SkippedID = gate.TriageSkipped, gate.ReasonOffTopic, decision.Item.ID
			summary.Skipped = append(summary.Skipped, entry)
		}
	}
	return byVideo, nil
}

// runIngestChannelRescue is `kb ingest channel --rescue <id> --topic <t>`:
// it fetches and ingests one video skipped before fetch.
func runIngestChannelRescue(cmd *cobra.Command, options ingestChannelCommandOptions) error {
	const action = "ingest channel"
	target, err := resolveIngestTarget(cmd, action, options.topicSlug)
	if err != nil {
		return err
	}
	cfg, err := loadIngestConfig()
	if err != nil {
		return fmt.Errorf("%s: %w", action, err)
	}
	runner, err := openIngestRunner(cmd, "channel", target.TopicInfo.Slug, options.flags, resolveIngestBatch("channel", options.batch), "")
	if err != nil {
		return err
	}
	fetchers := rescueFetchers{cfg: cfg, scraper: newFirecrawlScraper(firecrawlConfig(cfg.Firecrawl))}
	return runIngestRescue(cmd, runner, target, fetchers, options.flags.Rescue)
}

func resolveChannelIngestTarget(
	action string,
	vaultPath string,
	topicRef ktopic.TopicRef,
	title string,
	createTopic bool,
) (ingestTarget, error) {
	topicInfo, err := runIngestTopicInfo(vaultPath, topicRef.Path)
	if err == nil {
		return ingestTarget{TopicInfo: topicInfo, VaultPath: vaultPath}, nil
	}
	if !errors.Is(err, ktopic.ErrTopicNotFound) {
		return ingestTarget{}, fmt.Errorf("%s: %w", action, err)
	}
	if !createTopic {
		return ingestTarget{}, fmt.Errorf(
			"%s: %w; create it with: %s",
			action,
			err,
			missingTopicCreateCommand(topicRef.Path, title, "youtube-channel"),
		)
	}

	topicInfo, err = runIngestTopicNew(vaultPath, topicRef.Path, title, "youtube-channel")
	if err != nil {
		return ingestTarget{}, fmt.Errorf("%s: create topic: %w", action, err)
	}
	return ingestTarget{TopicInfo: topicInfo, VaultPath: vaultPath}, nil
}

func deriveChannelTopicTitle(listing youtube.ChannelListing, topicRef ktopic.TopicRef) string {
	for _, candidate := range []string{listing.Channel, listing.Uploader, listing.Title} {
		if title := strings.TrimSpace(candidate); title != "" {
			return title
		}
	}
	return vault.DeriveTopicTitle(topicRef.Leaf)
}

func newChannelVideosToFetch(
	videos []youtube.ChannelVideo,
	existing map[string]struct{},
	summary *channelIngestSummary,
) []youtube.ChannelVideo {
	toFetch := make([]youtube.ChannelVideo, 0, len(videos))
	for _, video := range videos {
		if _, done := existing[video.VideoID]; done {
			summary.Skipped = append(summary.Skipped, summarizeVideo(video))
			continue
		}
		toFetch = append(toFetch, video)
	}
	return toFetch
}

// recordChannelOutcome ingests one extracted video; it returns only fatal
// errors (auth, cancellation, I/O), extraction failures are summarized.
func recordChannelOutcome(
	ctx context.Context,
	runner ingestRunner,
	target ingestTarget,
	provenance channelProvenance,
	outcome youtube.VideoOutcome,
	decision *gate.PrefetchDecision,
	summary *channelIngestSummary,
) error {
	entry := summarizeVideo(outcome.Video)
	if outcome.Err != nil {
		entry.Error = outcome.Err.Error()
		summary.Failures = append(summary.Failures, entry)
		return nil
	}
	result, err := ingestChannelVideo(ctx, runner, target, provenance, outcome, decision)
	if err != nil {
		entry.Error = err.Error()
		summary.Failures = append(summary.Failures, entry)
		return err
	}
	entry.FilePath = result.FilePath
	entry.Triage, entry.Reason = result.Triage, result.TriageReason
	if result.Triage == gate.TriageDuplicateSkipped {
		summary.Skipped = append(summary.Skipped, entry)
		return nil
	}
	summary.Ingested = append(summary.Ingested, entry)
	return nil
}

func newChannelIngestSummary(
	target ingestTarget,
	channelURL string,
	normalizedURL string,
	policy youtube.TranscriptionPolicy,
	captionLanguages []string,
	limit int,
	dryRun bool,
	videos []youtube.ChannelVideo,
) channelIngestSummary {
	summary := channelIngestSummary{
		Topic:                target.TopicInfo.Slug,
		ChannelURL:           channelURL,
		NormalizedChannelURL: normalizedURL,
		Selection:            channelSelectionLabel(limit),
		Transcribe:           string(policy),
		CaptionLanguages:     append([]string(nil), captionLanguages...),
		Resolved:             len(videos),
		DryRun:               dryRun,
		Videos:               make([]channelVideoSummary, 0, len(videos)),
		Ingested:             []channelVideoSummary{},
		Skipped:              []channelVideoSummary{},
		Failures:             []channelVideoSummary{},
	}
	for _, video := range videos {
		summary.Videos = append(summary.Videos, summarizeVideo(video))
	}
	return summary
}

func resolveChannelBulkOptions(
	cmd *cobra.Command,
	youtubeConfig kconfig.YouTubeConfig,
	policy youtube.TranscriptionPolicy,
	captionLanguages []string,
	concurrency int,
	throttle time.Duration,
) (youtube.BulkOptions, error) {
	throttleValue := throttle
	if !cmd.Flags().Changed("throttle") {
		parsed, err := youtubeConfig.BulkThrottleDuration()
		if err != nil {
			return youtube.BulkOptions{}, err
		}
		throttleValue = parsed
	}

	concurrencyValue := concurrency
	if !cmd.Flags().Changed("concurrency") {
		concurrencyValue = youtubeConfig.BulkConcurrency
	}

	backoffMax, err := youtubeConfig.BulkBackoffMaxDuration()
	if err != nil {
		return youtube.BulkOptions{}, err
	}

	return youtube.BulkOptions{
		TranscriptionPolicy:     policy,
		PreferredLanguages:      append([]string(nil), captionLanguages...),
		AllowTranslatedCaptions: youtubeConfig.AllowTranslatedCaptions,
		Concurrency:             concurrencyValue,
		Throttle:                throttleValue,
		BackoffMax:              backoffMax,
		MaxRetries:              youtubeConfig.BulkRetries,
	}, nil
}

func ingestChannelVideo(
	ctx context.Context,
	runner ingestRunner,
	target ingestTarget,
	provenance channelProvenance,
	outcome youtube.VideoOutcome,
	decision *gate.PrefetchDecision,
) (kingest.Result, error) {
	result := outcome.Result
	sourceURL := strings.TrimSpace(result.Metadata.URL)
	if sourceURL == "" {
		sourceURL = outcome.Video.URL
	}
	title := strings.TrimSpace(result.Metadata.Title)
	if title == "" {
		title = outcome.Video.Title
	}
	options := kingest.Options{
		VaultPath:        target.VaultPath,
		Topic:            target.TopicInfo.Slug,
		SourceKind:       models.SourceKindYouTubeTranscript,
		SourceURL:        sourceURL,
		Title:            title,
		Markdown:         result.Markdown,
		ExtraFrontmatter: youtubeFrontmatter(result),
		Batch:            provenance.batch,
		Query:            provenance.query,
		Gate:             kingest.GateOptions{PlatformID: "youtube:" + outcome.Video.VideoID},
	}
	if decision != nil && decision.Action != "" {
		options.Gate.Prefetch = decision
	}
	return runner.Ingest(ctx, options)
}

func summarizeVideo(video youtube.ChannelVideo) channelVideoSummary {
	return channelVideoSummary{VideoID: video.VideoID, Title: video.Title, URL: video.URL}
}

func channelSelectionLabel(limit int) string {
	if limit <= 0 {
		return "all uploads"
	}
	return fmt.Sprintf("latest %d uploads", limit)
}
