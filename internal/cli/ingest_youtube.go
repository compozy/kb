package cli

import (
	"context"
	"fmt"
	"strings"

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
	values := map[string]any{}
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
