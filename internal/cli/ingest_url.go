package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	kconfig "github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/firecrawl"
	kingest "github.com/compozy/kb/internal/ingest"
	"github.com/compozy/kb/internal/models"
)

type firecrawlScraper interface {
	Scrape(ctx context.Context, sourceURL string) (*firecrawl.ScrapeResult, error)
}

var newFirecrawlScraper = func(cfg firecrawlConfig) firecrawlScraper {
	return firecrawl.NewClient(cfg)
}

type firecrawlConfig = kconfig.FirecrawlConfig

func newIngestURLCommand() *cobra.Command {
	var topic string

	command := &cobra.Command{
		Use:   "url <url>",
		Short: "Scrape a web URL and ingest it into a topic",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := resolveIngestTarget(cmd, "ingest url", topic)
			if err != nil {
				return err
			}

			cfg, err := loadIngestConfig()
			if err != nil {
				return fmt.Errorf("ingest url: %w", err)
			}

			scrapeResult, err := newFirecrawlScraper(firecrawlConfig(cfg.Firecrawl)).Scrape(commandContext(cmd), args[0])
			if err != nil {
				return fmt.Errorf("ingest url: %w", err)
			}

			sourceURL := strings.TrimSpace(scrapeResult.SourceURL)
			if sourceURL == "" {
				sourceURL = args[0]
			}

			result, err := runIngest(commandContext(cmd), kingest.Options{
				VaultPath:  target.VaultPath,
				Topic:      target.TopicInfo.Slug,
				SourceKind: models.SourceKindArticle,
				SourceURL:  sourceURL,
				Title:      scrapeResult.Title,
				Markdown:   scrapeResult.Markdown,
			})
			if err != nil {
				return fmt.Errorf("ingest url: %w", err)
			}

			return writeJSON(cmd, result)
		},
	}

	requireTopicFlag(command, &topic)

	return command
}
