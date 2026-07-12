package cli

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/compozy/kb/internal/browser"
	kconfig "github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/firecrawl"
	kingest "github.com/compozy/kb/internal/ingest"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/urlfetch"
)

const (
	urlProviderHTTPLocal = "http-local"
	urlProviderBrowser   = "browser"
	urlProviderFirecrawl = "firecrawl"
)

type firecrawlScraper interface {
	Scrape(ctx context.Context, sourceURL string) (*firecrawl.ScrapeResult, error)
}

type urlContentFetcher interface {
	Fetch(ctx context.Context, sourceURL string) (*urlfetch.Result, error)
}

var newFirecrawlScraper = func(cfg firecrawlConfig) firecrawlScraper {
	return firecrawl.NewClient(cfg)
}
var newHTTPLocalFetcher = func() urlContentFetcher { return urlfetch.NewClient() }
var newBrowserFetcher = func(cfg browserConfig) urlContentFetcher { return browser.NewClient(cfg.Command) }

type firecrawlConfig = kconfig.FirecrawlConfig
type browserConfig = kconfig.BrowserConfig

func newIngestURLCommand() *cobra.Command {
	var topic string
	var provider string
	var render bool

	command := &cobra.Command{
		Use:   "url <url>",
		Short: "Fetch a web URL and ingest it into a topic",
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

			selectedProvider := strings.ToLower(strings.TrimSpace(provider))
			if render {
				if cmd.Flags().Changed("provider") && selectedProvider != urlProviderHTTPLocal {
					return fmt.Errorf("ingest url: --render cannot be combined with --provider %s", selectedProvider)
				}
				selectedProvider = urlProviderBrowser
			}

			switch selectedProvider {
			case urlProviderHTTPLocal:
				return ingestURLContent(cmd, target, args[0], newHTTPLocalFetcher())
			case urlProviderBrowser:
				return ingestURLContent(cmd, target, args[0], newBrowserFetcher(browserConfig(cfg.Browser)))
			case urlProviderFirecrawl:
				return ingestURLFirecrawl(cmd, target, args[0], newFirecrawlScraper(firecrawlConfig(cfg.Firecrawl)))
			default:
				return fmt.Errorf("ingest url: unsupported provider %q (want http-local, browser, or firecrawl)", provider)
			}
		},
	}

	requireTopicFlag(command, &topic)
	command.Flags().StringVar(&provider, "provider", urlProviderHTTPLocal, "URL provider: http-local, browser, or firecrawl")
	command.Flags().BoolVar(&render, "render", false, "render JavaScript with the configured browser provider")

	return command
}

func ingestURLContent(cmd *cobra.Command, target ingestTarget, sourceURL string, fetcher urlContentFetcher) error {
	result, err := fetcher.Fetch(commandContext(cmd), sourceURL)
	if err != nil {
		return fmt.Errorf("ingest url: %w", err)
	}
	if result == nil {
		return fmt.Errorf("ingest url: provider returned no content")
	}

	ingestResult, err := runIngest(commandContext(cmd), kingest.Options{
		VaultPath:        target.VaultPath,
		Topic:            target.TopicInfo.Slug,
		SourceKind:       models.SourceKindArticle,
		SourceURL:        sourceURL,
		SourceContent:    result.Body,
		ConvertFilePath:  result.FileName,
		ConvertOptions:   map[string]any{"content_type": result.ContentType},
		ExtraFrontmatter: urlProvenance(result),
		ScrapedAt:        result.FetchedAt,
	})
	if err != nil {
		return fmt.Errorf("ingest url: %w", err)
	}
	return writeJSON(cmd, ingestResult)
}

func ingestURLFirecrawl(cmd *cobra.Command, target ingestTarget, sourceURL string, scraper firecrawlScraper) error {
	scrapeResult, err := scraper.Scrape(commandContext(cmd), sourceURL)
	if err != nil {
		return fmt.Errorf("ingest url: %w", err)
	}
	if scrapeResult == nil {
		return fmt.Errorf("ingest url: firecrawl returned no content")
	}
	contentHash := sha256.Sum256([]byte(scrapeResult.Markdown))
	now := time.Now().UTC()
	result, err := runIngest(commandContext(cmd), kingest.Options{
		VaultPath:  target.VaultPath,
		Topic:      target.TopicInfo.Slug,
		SourceKind: models.SourceKindArticle,
		SourceURL:  sourceURL,
		Title:      scrapeResult.Title,
		Markdown:   scrapeResult.Markdown,
		ExtraFrontmatter: map[string]any{
			"final_url":    firstNonEmptyURL(scrapeResult.SourceURL, sourceURL),
			"fetched_at":   now.Format(time.RFC3339),
			"content_type": "text/markdown",
			"content_hash": fmt.Sprintf("sha256:%x", contentHash),
		},
		ScrapedAt: now,
	})
	if err != nil {
		return fmt.Errorf("ingest url: %w", err)
	}
	return writeJSON(cmd, result)
}

func urlProvenance(result *urlfetch.Result) map[string]any {
	return map[string]any{
		"final_url":    firstNonEmptyURL(result.FinalURL, result.SourceURL),
		"fetched_at":   result.FetchedAt.UTC().Format(time.RFC3339),
		"content_type": result.ContentType,
		"content_hash": result.ContentHash,
	}
}

func firstNonEmptyURL(values ...string) string {
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}
