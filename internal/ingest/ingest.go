// Package ingest orchestrates multi-source content ingestion into knowledge base topics, handling frontmatter assembly, raw file writes, and operation log entries.
package ingest

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/compozy/kb/internal/convert"
	"github.com/compozy/kb/internal/frontmatter"
	"github.com/compozy/kb/internal/models"
	"github.com/compozy/kb/internal/topic"
	"github.com/compozy/kb/internal/vault"
)

var bookmarkURLPattern = regexp.MustCompile(`https?://[^\s<>()]+`)

// Registry converts file-backed inputs into markdown content.
type Registry interface {
	Convert(ctx context.Context, input models.ConvertInput) (*models.ConvertResult, error)
}

// Options configures a single-source ingest run.
type Options struct {
	VaultPath      string
	Topic          string
	SourceKind     models.SourceKind
	SourcePath     string
	SourceURL      string
	Title          string
	Markdown       string
	ConvertOptions map[string]any
	Registry       Registry
	ScrapedAt      time.Time
}

// Ingest validates the target topic, optionally converts the source, writes the
// raw markdown document, and appends a log entry.
func Ingest(ctx context.Context, options Options) (models.IngestResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	topicInfo, err := topic.Info(options.VaultPath, options.Topic)
	if err != nil {
		return models.IngestResult{}, fmt.Errorf("ingest: validate topic: %w", err)
	}
	if err := topic.EnsureCurrentSkeleton(topicInfo.RootPath); err != nil {
		return models.IngestResult{}, fmt.Errorf("ingest: ensure topic skeleton: %w", err)
	}

	sourceDirectory, err := rawDirectoryForSourceKind(options.SourceKind)
	if err != nil {
		return models.IngestResult{}, fmt.Errorf("ingest: %w", err)
	}

	title, markdown, err := resolveMarkdown(ctx, options)
	if err != nil {
		return models.IngestResult{}, fmt.Errorf("ingest: %w", err)
	}

	now := options.ScrapedAt.UTC()
	if options.ScrapedAt.IsZero() {
		now = time.Now().UTC()
	}

	targetDirectory := filepath.Join(topicInfo.RootPath, filepath.FromSlash(path.Join("raw", sourceDirectory)))
	if err := os.MkdirAll(targetDirectory, 0o755); err != nil {
		return models.IngestResult{}, fmt.Errorf("ingest: create target directory %q: %w", targetDirectory, err)
	}

	slug, err := uniqueSlug(targetDirectory, vault.SlugifySegment(title))
	if err != nil {
		return models.IngestResult{}, fmt.Errorf("ingest: allocate slug: %w", err)
	}

	relativeTopicPath := path.Join("raw", sourceDirectory, slug+".md")
	absoluteTargetPath := filepath.Join(topicInfo.RootPath, filepath.FromSlash(relativeTopicPath))

	values, err := buildFrontmatter(topicInfo, options, title, markdown, now)
	if err != nil {
		return models.IngestResult{}, fmt.Errorf("ingest: %w", err)
	}

	document, err := frontmatter.Generate(values, markdown)
	if err != nil {
		return models.IngestResult{}, fmt.Errorf("ingest: generate frontmatter: %w", err)
	}

	if err := os.WriteFile(absoluteTargetPath, []byte(document), 0o644); err != nil {
		return models.IngestResult{}, fmt.Errorf("ingest: write %q: %w", absoluteTargetPath, err)
	}

	result := models.IngestResult{
		Topic:      topicInfo.Slug,
		SourceType: options.SourceKind,
		FilePath:   path.Join(topicInfo.Slug, relativeTopicPath),
		Title:      title,
	}

	if err := appendLogEntry(filepath.Join(topicInfo.RootPath, "log.md"), now, slug, options.SourceKind, result.FilePath, title); err != nil {
		return models.IngestResult{}, fmt.Errorf("ingest: append log entry: %w", err)
	}

	return result, nil
}

func resolveMarkdown(ctx context.Context, options Options) (string, string, error) {
	if body := strings.TrimSpace(options.Markdown); body != "" {
		title := deriveTitle(strings.TrimSpace(options.Title), options.SourcePath, options.SourceURL, options.SourceKind)
		return title, options.Markdown, nil
	}

	sourcePath := strings.TrimSpace(options.SourcePath)
	if sourcePath == "" {
		return "", "", errors.New("source path or markdown content is required")
	}

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return "", "", fmt.Errorf("open source file %q: %w", sourcePath, err)
	}
	defer func() {
		_ = sourceFile.Close()
	}()

	registry := options.Registry
	if registry == nil {
		registry = convert.NewRegistry()
	}

	result, err := registry.Convert(ctx, models.ConvertInput{
		Reader:   sourceFile,
		FilePath: sourcePath,
		URL:      strings.TrimSpace(options.SourceURL),
		Options:  cloneMap(options.ConvertOptions),
	})
	if err != nil {
		return "", "", fmt.Errorf("convert source: %w", err)
	}
	if result == nil {
		return "", "", errors.New("convert source: converter returned nil result")
	}
	if strings.TrimSpace(result.Markdown) == "" {
		return "", "", errors.New("convert source: markdown output is empty")
	}

	title := deriveTitle(firstNonEmpty(strings.TrimSpace(options.Title), strings.TrimSpace(result.Title)), options.SourcePath, options.SourceURL, options.SourceKind)
	return title, result.Markdown, nil
}

func buildFrontmatter(
	topicInfo models.TopicInfo,
	options Options,
	title string,
	markdown string,
	scrapedAt time.Time,
) (map[string]any, error) {
	stamp := scrapedAt.Format(frontmatter.DateLayout)
	values := map[string]any{
		"title":       title,
		"type":        "source",
		"stage":       "raw",
		"domain":      topicInfo.Domain,
		"source_kind": string(options.SourceKind),
		"scraped":     stamp,
		"tags":        []string{topicInfo.Domain, "raw", string(options.SourceKind)},
	}

	if sourceURL := strings.TrimSpace(options.SourceURL); sourceURL != "" {
		values["source_url"] = sourceURL
	}
	if sourcePath := normalizedSourcePath(options.SourcePath); sourcePath != "" {
		values["source_path"] = sourcePath
	}

	if options.SourceKind == models.SourceKindBookmarkCluster {
		sourceURLs := extractSourceURLs(markdown)
		if len(sourceURLs) == 0 {
			return nil, errors.New("bookmark cluster must contain at least one URL")
		}

		values["status"] = "seeded"
		values["created"] = stamp
		values["updated"] = stamp
		values["source_urls"] = sourceURLs
	}

	return values, nil
}

func appendLogEntry(logPath string, when time.Time, slug string, sourceKind models.SourceKind, filePath string, title string) error {
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open %q: %w", logPath, err)
	}

	entry := strings.Join([]string{
		"",
		fmt.Sprintf("## [%s] ingest | %s.md (%s)", when.Format(frontmatter.DateLayout), slug, sourceKind),
		"",
		fmt.Sprintf("Ingested `%s` into `%s`.", title, filePath),
		"",
	}, "\n")

	if _, err := file.WriteString(entry); err != nil {
		if closeErr := file.Close(); closeErr != nil {
			return fmt.Errorf("write log entry: %w (close error: %v)", err, closeErr)
		}
		return fmt.Errorf("write log entry: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close %q: %w", logPath, err)
	}

	return nil
}

func rawDirectoryForSourceKind(sourceKind models.SourceKind) (string, error) {
	switch sourceKind {
	case models.SourceKindArticle, models.SourceKindDocument:
		return "articles", nil
	case models.SourceKindGitHubREADME:
		return "github", nil
	case models.SourceKindYouTubeTranscript:
		return "youtube", nil
	case models.SourceKindBookmarkCluster:
		return "bookmarks", nil
	case models.SourceKindCodebaseFile:
		return "codebase/files", nil
	case models.SourceKindCodebaseSymbol:
		return "codebase/symbols", nil
	case "":
		return "", errors.New("source kind is required")
	default:
		return "", fmt.Errorf("unsupported source kind %q", sourceKind)
	}
}

func uniqueSlug(directoryPath, baseSlug string) (string, error) {
	candidate := baseSlug
	for index := 2; ; index++ {
		targetPath := filepath.Join(directoryPath, candidate+".md")
		_, err := os.Stat(targetPath)
		if errors.Is(err, os.ErrNotExist) {
			return candidate, nil
		}
		if err != nil {
			return "", fmt.Errorf("stat %q: %w", targetPath, err)
		}

		candidate = fmt.Sprintf("%s-%d", baseSlug, index)
	}
}

func deriveTitle(initial, sourcePath, sourceURL string, sourceKind models.SourceKind) string {
	if title := strings.TrimSpace(initial); title != "" {
		return title
	}
	if sourcePath = strings.TrimSpace(sourcePath); sourcePath != "" {
		baseName := strings.TrimSuffix(filepath.Base(sourcePath), filepath.Ext(sourcePath))
		if title := humanizeSegment(baseName); title != "" {
			return title
		}
	}
	if title := titleFromURL(sourceURL); title != "" {
		return title
	}
	if title := humanizeSegment(string(sourceKind)); title != "" {
		return title
	}

	return "Untitled Source"
}

func titleFromURL(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ""
	}

	trimmed := strings.TrimSuffix(strings.TrimPrefix(rawURL, "https://"), "/")
	trimmed = strings.TrimSuffix(strings.TrimPrefix(trimmed, "http://"), "/")
	if trimmed == "" {
		return ""
	}

	lastSlash := strings.LastIndex(trimmed, "/")
	if lastSlash >= 0 && lastSlash+1 < len(trimmed) {
		return humanizeSegment(trimmed[lastSlash+1:])
	}

	host := trimmed
	if slash := strings.Index(host, "/"); slash >= 0 {
		host = host[:slash]
	}

	return humanizeSegment(host)
}

func humanizeSegment(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	slug := vault.SlugifySegment(value)
	if title := vault.HumanizeSlug(slug); title != "" {
		return title
	}

	return value
}

func normalizedSourcePath(sourcePath string) string {
	sourcePath = strings.TrimSpace(sourcePath)
	if sourcePath == "" {
		return ""
	}

	return vault.ToPosixPath(filepath.Clean(sourcePath))
}

func extractSourceURLs(markdown string) []string {
	matches := bookmarkURLPattern.FindAllString(markdown, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(matches))
	urls := make([]string, 0, len(matches))
	for _, match := range matches {
		cleaned := strings.TrimRight(strings.TrimSpace(match), ".,;:")
		if cleaned == "" {
			continue
		}
		if _, exists := seen[cleaned]; exists {
			continue
		}
		seen[cleaned] = struct{}{}
		urls = append(urls, cleaned)
	}

	return urls
}

func cloneMap(values map[string]any) map[string]any {
	if len(values) == 0 {
		return nil
	}

	cloned := make(map[string]any, len(values))
	for key, value := range values {
		cloned[key] = value
	}

	return cloned
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}

	return ""
}
