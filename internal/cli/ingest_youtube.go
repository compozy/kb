package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	kconfig "github.com/user/kb/internal/config"
	kingest "github.com/user/kb/internal/ingest"
	"github.com/user/kb/internal/models"
	"github.com/user/kb/internal/youtube"
)

type youtubeTranscriptExtractor interface {
	Extract(ctx context.Context, rawURL string, options youtube.ExtractOptions) (*youtube.Result, error)
}

var newYouTubeTranscriptExtractor = func(cfg kconfig.OpenRouterConfig) youtubeTranscriptExtractor {
	return youtube.NewExtractor(cfg)
}

func newIngestYouTubeCommand() *cobra.Command {
	var topic string
	var enableSTT bool

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

			extractResult, err := newYouTubeTranscriptExtractor(cfg.OpenRouter).Extract(
				commandContext(cmd),
				args[0],
				youtube.ExtractOptions{
					EnableSTTFallback: enableSTT,
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
				VaultPath:  target.VaultPath,
				Topic:      target.TopicInfo.Slug,
				SourceKind: models.SourceKindYouTubeTranscript,
				SourceURL:  sourceURL,
				Title:      extractResult.Metadata.Title,
				Markdown:   extractResult.Markdown,
			})
			if err != nil {
				return fmt.Errorf("ingest youtube: %w", err)
			}

			return writeJSON(cmd, result)
		},
	}

	requireTopicFlag(command, &topic)
	command.Flags().BoolVar(&enableSTT, "stt", false, "Enable OpenRouter speech-to-text fallback when captions are unavailable")

	return command
}
