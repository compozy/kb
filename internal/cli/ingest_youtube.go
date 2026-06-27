package cli

import (
	"context"
	"fmt"
	"reflect"
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
	var subLangs string
	var lang string

	command := &cobra.Command{
		Use:   "youtube <url>",
		Short: "Extract a YouTube transcript and ingest it into a topic",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := resolveIngestTargetWithMissingTopicHint(
				cmd,
				"ingest youtube",
				topic,
				"<title>",
				"youtube-channel",
			)
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
			captionLanguages, err := resolveYouTubeCaptionLanguages(cmd, cfg.YouTube, subLangs, lang)
			if err != nil {
				return fmt.Errorf("ingest youtube: %w", err)
			}

			extractResult, err := newYouTubeTranscriptExtractor(cfg).Extract(
				commandContext(cmd),
				args[0],
				youtube.ExtractOptions{
					TranscriptionPolicy:     policy,
					PreferredLanguages:      captionLanguages,
					AllowTranslatedCaptions: cfg.YouTube.AllowTranslatedCaptions,
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
	command.Flags().StringVar(&subLangs, "sub-langs", "", "Caption languages to request, comma-separated; use orig for the video's original language")
	command.Flags().StringVar(&lang, "lang", "", "Alias for --sub-langs")

	return command
}

func resolveYouTubeCaptionLanguages(
	cmd *cobra.Command,
	youtubeConfig kconfig.YouTubeConfig,
	subLangs string,
	lang string,
) ([]string, error) {
	configured := normalizeYouTubeCaptionLanguages(youtubeConfig.CaptionLanguages)
	if len(configured) == 0 {
		configured = []string{"orig"}
	}

	subLangsChanged := cmd.Flags().Changed("sub-langs")
	langChanged := cmd.Flags().Changed("lang")
	switch {
	case subLangsChanged && langChanged:
		left, err := parseYouTubeCaptionLanguages(subLangs)
		if err != nil {
			return nil, err
		}
		right, err := parseYouTubeCaptionLanguages(lang)
		if err != nil {
			return nil, err
		}
		if !reflect.DeepEqual(left, right) {
			return nil, fmt.Errorf("--sub-langs and --lang must not disagree")
		}
		return left, nil
	case subLangsChanged:
		return parseYouTubeCaptionLanguages(subLangs)
	case langChanged:
		return parseYouTubeCaptionLanguages(lang)
	default:
		return configured, nil
	}
}

func parseYouTubeCaptionLanguages(value string) ([]string, error) {
	languages := normalizeYouTubeCaptionLanguages(strings.Split(value, ","))
	if len(languages) == 0 {
		return nil, fmt.Errorf("caption language list must contain at least one language or orig")
	}
	return languages, nil
}

func normalizeYouTubeCaptionLanguages(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	normalized := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		value = strings.ToLower(strings.TrimSpace(value))
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		normalized = append(normalized, value)
	}
	return normalized
}

func youtubeFrontmatter(result *youtube.Result) map[string]any {
	if result == nil {
		return nil
	}
	metadata := result.Metadata
	values := map[string]any{
		"video_id":               optionalString(metadata.VideoID),
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
