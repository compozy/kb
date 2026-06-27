package cli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	kconfig "github.com/compozy/kb/internal/config"
	kingest "github.com/compozy/kb/internal/ingest"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/youtube"
)

type youtubeChannelExtractor interface {
	ListChannelVideos(ctx context.Context, normalizedURL string, limit int) ([]youtube.ChannelVideo, error)
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

func newIngestChannelCommand() *cobra.Command {
	var topicSlug string
	var transcribe string
	var limit int
	var all bool
	var concurrency int
	var throttle time.Duration
	var dryRun bool
	var subLangs string
	var lang string

	command := &cobra.Command{
		Use:   "channel <url>",
		Short: "Bulk-extract YouTube channel or playlist transcripts into a topic",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := resolveIngestTarget(cmd, "ingest channel", topicSlug)
			if err != nil {
				return err
			}
			cfg, err := loadIngestConfig()
			if err != nil {
				return fmt.Errorf("ingest channel: %w", err)
			}

			policyValue := strings.TrimSpace(cfg.YouTube.Transcription)
			if cmd.Flags().Changed("transcribe") {
				policyValue = transcribe
			}
			policy, err := youtube.ParseTranscriptionPolicy(policyValue)
			if err != nil {
				return fmt.Errorf("ingest channel: %w", err)
			}
			captionLanguages, err := resolveYouTubeCaptionLanguages(cmd, cfg.YouTube, subLangs, lang)
			if err != nil {
				return fmt.Errorf("ingest channel: %w", err)
			}

			normalizedURL, err := youtube.NormalizeChannelURL(args[0])
			if err != nil {
				return fmt.Errorf("ingest channel: %w", err)
			}

			resolvedLimit := limit
			if all {
				resolvedLimit = 0
			}

			ctx := commandContext(cmd)
			extractor := newYouTubeChannelExtractor(cfg)
			videos, err := extractor.ListChannelVideos(ctx, normalizedURL, resolvedLimit)
			if err != nil {
				return fmt.Errorf("ingest channel: %w", err)
			}

			summary := newChannelIngestSummary(target, args[0], normalizedURL, policy, captionLanguages, resolvedLimit, dryRun, videos)

			if dryRun {
				return writeJSON(cmd, summary)
			}

			existing, err := existingYouTubeVideoIDs(target.VaultPath, target.TopicInfo.Slug)
			if err != nil {
				return fmt.Errorf("ingest channel: %w", err)
			}

			toFetch := make([]youtube.ChannelVideo, 0, len(videos))
			for _, video := range videos {
				if _, done := existing[video.VideoID]; done {
					summary.Skipped = append(summary.Skipped, summarizeVideo(video))
					continue
				}
				toFetch = append(toFetch, video)
			}

			options, err := resolveChannelBulkOptions(cmd, cfg.YouTube, policy, captionLanguages, concurrency, throttle)
			if err != nil {
				return fmt.Errorf("ingest channel: %w", err)
			}

			sink := func(outcome youtube.VideoOutcome) {
				entry := summarizeVideo(outcome.Video)
				if outcome.Err != nil {
					entry.Error = outcome.Err.Error()
					summary.Failures = append(summary.Failures, entry)
					return
				}
				result, err := ingestChannelVideo(ctx, target, outcome)
				if err != nil {
					entry.Error = err.Error()
					summary.Failures = append(summary.Failures, entry)
					return
				}
				entry.FilePath = result.FilePath
				summary.Ingested = append(summary.Ingested, entry)
			}

			bulkErr := extractor.BulkExtract(ctx, toFetch, options, sink)

			if writeErr := writeJSON(cmd, summary); writeErr != nil {
				return writeErr
			}
			if bulkErr != nil {
				return fmt.Errorf("ingest channel: %w", bulkErr)
			}
			return nil
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

	return command
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

func ingestChannelVideo(ctx context.Context, target ingestTarget, outcome youtube.VideoOutcome) (models.IngestResult, error) {
	result := outcome.Result
	sourceURL := strings.TrimSpace(result.Metadata.URL)
	if sourceURL == "" {
		sourceURL = outcome.Video.URL
	}
	title := strings.TrimSpace(result.Metadata.Title)
	if title == "" {
		title = outcome.Video.Title
	}
	return runIngest(ctx, kingest.Options{
		VaultPath:        target.VaultPath,
		Topic:            target.TopicInfo.Slug,
		SourceKind:       models.SourceKindYouTubeTranscript,
		SourceURL:        sourceURL,
		Title:            title,
		Markdown:         result.Markdown,
		ExtraFrontmatter: youtubeFrontmatter(result),
	})
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
