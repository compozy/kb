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

type youtubeTranscriptExtractor interface {
	Extract(ctx context.Context, rawURL string, options youtube.ExtractOptions) (*youtube.Result, error)
}

var newYouTubeTranscriptExtractor = func(cfg kconfig.Config) youtubeTranscriptExtractor {
	return youtube.NewExtractorWithConfig(cfg.STT, cfg.OpenRouter, cfg.YouTube)
}

func newIngestYouTubeCommand() *cobra.Command {
	var topic string
	var transcribe string

	command := &cobra.Command{
		Use:   "youtube <url>",
		Short: "Extract a YouTube transcript and ingest it into a topic",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := resolveIngestTarget(cmd, "ingest youtube", topic)
			if err != nil {
				return err
			}

			cfg, err := loadIngestConfig()
			if err != nil {
				return fmt.Errorf("ingest youtube: %w", err)
			}
			policyValue := strings.TrimSpace(cfg.YouTube.Transcription)
			if cmd.Flags().Changed("transcribe") {
				policyValue = transcribe
			}
			policy, err := youtube.ParseTranscriptionPolicy(policyValue)
			if err != nil {
				return fmt.Errorf("ingest youtube: %w", err)
			}

			extractResult, err := newYouTubeTranscriptExtractor(cfg).Extract(
				commandContext(cmd),
				args[0],
				youtube.ExtractOptions{
					TranscriptionPolicy: policy,
				},
			)
			if err != nil {
				return fmt.Errorf("ingest youtube: %w", err)
			}

			sourceURL := strings.TrimSpace(extractResult.Metadata.URL)
			if sourceURL == "" {
				sourceURL = args[0]
			}

			result, err := runIngest(commandContext(cmd), kingest.Options{
				VaultPath:        target.VaultPath,
				Topic:            target.TopicInfo.Slug,
				SourceKind:       models.SourceKindYouTubeTranscript,
				SourceURL:        sourceURL,
				Title:            extractResult.Metadata.Title,
				Markdown:         extractResult.Markdown,
				ExtraFrontmatter: youtubeFrontmatter(extractResult),
			})
			if err != nil {
				return fmt.Errorf("ingest youtube: %w", err)
			}

			return writeJSON(cmd, result)
		},
	}

	requireTopicFlag(command, &topic)
	command.Flags().StringVar(&transcribe, "transcribe", "", "Transcription policy: captions, auto, or stt")

	return command
}

func youtubeFrontmatter(result *youtube.Result) map[string]any {
	if result == nil {
		return nil
	}
	metadata := result.Metadata
	values := map[string]any{
		"view_count":             optionalInt64(metadata.ViewCount),
		"like_count":             optionalInt64(metadata.LikeCount),
		"comment_count":          optionalInt64(metadata.CommentCount),
		"upload_date":            optionalDate(metadata.PublishDate),
		"duration":               optionalDurationSeconds(metadata.Duration),
		"duration_string":        optionalString(metadata.DurationString),
		"channel":                optionalString(metadata.Channel),
		"channel_id":             optionalString(metadata.ChannelID),
		"uploader_id":            optionalString(metadata.UploaderID),
		"channel_follower_count": optionalInt64(metadata.ChannelFollowerCount),
		"categories":             cloneStringSlice(metadata.Categories),
		"youtube_tags":           cloneStringSlice(metadata.VideoTags),
		"language":               optionalString(metadata.Language),
		"live_status":            optionalString(metadata.LiveStatus),
		"was_live":               optionalBool(metadata.WasLive),
		"chapter_count":          metadata.ChapterCount,
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

func optionalString(value string) any {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return trimmed
}

func optionalDate(value time.Time) any {
	if value.IsZero() {
		return nil
	}
	return value.UTC().Format("2006-01-02")
}

func optionalDurationSeconds(value time.Duration) any {
	if value <= 0 {
		return nil
	}
	return int64(value / time.Second)
}

func optionalInt64(value *int64) any {
	if value == nil {
		return nil
	}
	return *value
}

func optionalBool(value *bool) any {
	if value == nil {
		return nil
	}
	return *value
}

func cloneStringSlice(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}
	return append([]string(nil), values...)
}
