package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	kconfig "github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/gate"
	kingest "github.com/compozy/kb/internal/ingest"
	"github.com/compozy/kb/internal/instagram"
	"github.com/compozy/kb/internal/mediadl"
	"github.com/compozy/kb/internal/models"
)

type instagramTranscriptExtractor interface {
	Extract(ctx context.Context, rawURL string, options mediadl.ExtractOptions) (*mediadl.Result, error)
}

var newInstagramTranscriptExtractor = func(cfg kconfig.Config) instagramTranscriptExtractor {
	return instagram.NewExtractorWithConfig(cfg.STT, cfg.OpenRouter, cfg.Instagram)
}

func newIngestInstagramCommand() *cobra.Command {
	var topic string
	var transcribe string
	var batch string
	var flags ingestFlags

	command := &cobra.Command{
		Use:   "instagram <url>",
		Short: "Extract an Instagram reel/video caption and transcript and ingest it into a topic",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := resolveIngestTarget(cmd, "ingest instagram", topic)
			if err != nil {
				return err
			}

			cfg, err := loadIngestConfig()
			if err != nil {
				return fmt.Errorf("ingest instagram: %w", err)
			}
			policyValue := strings.TrimSpace(cfg.Instagram.Transcription)
			if cmd.Flags().Changed("transcribe") {
				policyValue = transcribe
			}
			policy, err := mediadl.ParseTranscriptionPolicy(policyValue)
			if err != nil {
				return fmt.Errorf("ingest instagram: %w", err)
			}

			runBatch := resolveIngestBatch("instagram", batch)
			runner, err := openIngestRunner(cmd, "instagram", target.TopicInfo.Slug, flags, runBatch, "")
			if err != nil {
				return err
			}
			ctx := commandContext(cmd)
			result, skipped, err := runner.Precheck(args[0], "", models.SourceKindInstagramVideo)
			if err != nil {
				return fmt.Errorf("ingest instagram: %w", err)
			}
			if !skipped {
				extractResult, err := newInstagramTranscriptExtractor(cfg).Extract(
					ctx,
					args[0],
					mediadl.ExtractOptions{
						TranscriptionPolicy: policy,
					},
				)
				if err != nil {
					return fmt.Errorf("ingest instagram: %w", err)
				}

				sourceURL := strings.TrimSpace(extractResult.Metadata.URL)
				if sourceURL == "" {
					sourceURL = args[0]
				}
				platformID := gate.PlatformID(args[0])
				if code := strings.TrimSpace(extractResult.Metadata.VideoID); code != "" {
					platformID = "instagram:" + code
				}

				result, err = runner.Ingest(ctx, kingest.Options{
					VaultPath:        target.VaultPath,
					Topic:            target.TopicInfo.Slug,
					SourceKind:       models.SourceKindInstagramVideo,
					SourceURL:        sourceURL,
					Title:            extractResult.Metadata.Title,
					Markdown:         extractResult.Markdown,
					ExtraFrontmatter: instagramFrontmatter(extractResult),
					Batch:            runBatch,
					Gate:             kingest.GateOptions{PlatformID: platformID},
				})
				if err != nil {
					return fmt.Errorf("ingest instagram: %w", err)
				}
			}

			if err := writeJSON(cmd, result); err != nil {
				return err
			}
			return finishIngestRun(cmd, runner)
		},
	}

	requireTopicFlag(command, &topic)
	command.Flags().StringVar(&transcribe, "transcribe", "", "Transcription policy: captions, auto, or stt")
	addBatchFlag(command, &batch)
	bindIngestFlags(command, &flags, false)

	return command
}

func instagramFrontmatter(result *mediadl.Result) map[string]any {
	if result == nil {
		return nil
	}
	metadata := result.Metadata
	values := map[string]any{
		"shortcode":       optionalString(metadata.VideoID),
		"uploader":        optionalString(metadata.Channel),
		"uploader_id":     optionalString(metadata.UploaderID),
		"like_count":      optionalInt64(metadata.LikeCount),
		"comment_count":   optionalInt64(metadata.CommentCount),
		"view_count":      optionalInt64(metadata.ViewCount),
		"upload_date":     optionalDate(metadata.PublishDate),
		"duration":        optionalDurationSeconds(metadata.Duration),
		"duration_string": optionalString(metadata.DurationString),
		"language":        optionalString(metadata.Language),
	}
	if result.Source != "" {
		values["transcript_source"] = string(result.Source)
	}
	if result.TranscriptionPolicy != "" {
		values["transcription_policy"] = string(result.TranscriptionPolicy)
	}
	if result.Language != "" {
		values["transcript_language"] = result.Language
	}
	if result.CaptionKind != "" {
		values["caption_kind"] = string(result.CaptionKind)
	}
	if result.STTProvider != "" {
		values["stt_provider"] = result.STTProvider
	}
	if result.STTModel != "" {
		values["stt_model"] = result.STTModel
	}
	return values
}
