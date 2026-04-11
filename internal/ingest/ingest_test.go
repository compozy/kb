package ingest

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/user/go-devstack/internal/frontmatter"
	"github.com/user/go-devstack/internal/models"
	"github.com/user/go-devstack/internal/topic"
)

var fixedScrapeTime = time.Date(2026, 4, 11, 12, 0, 0, 0, time.UTC)

func TestIngestFileUsesRegistryAndWritesArticleDocument(t *testing.T) {
	t.Parallel()

	vaultPath, topicSlug := scaffoldTopic(t)
	sourcePath := writeSourceFile(t, filepath.Join(t.TempDir(), "latency.txt"), "original content\n")

	registry := &stubRegistry{
		result: &models.ConvertResult{
			Title:    "Latency Budget",
			Markdown: "# Latency Budget\n\nKeep service latency inside the SLO.\n",
		},
	}

	result, err := Ingest(context.Background(), Options{
		VaultPath:  vaultPath,
		Topic:      topicSlug,
		SourceKind: models.SourceKindArticle,
		SourcePath: sourcePath,
		Registry:   registry,
		ScrapedAt:  fixedScrapeTime,
	})
	if err != nil {
		t.Fatalf("Ingest returned error: %v", err)
	}

	if registry.calls != 1 {
		t.Fatalf("registry calls = %d, want 1", registry.calls)
	}
	if registry.lastInput.FilePath != sourcePath {
		t.Fatalf("registry file path = %q, want %q", registry.lastInput.FilePath, sourcePath)
	}

	wantFilePath := "systems-design/raw/articles/latency-budget.md"
	if result.Topic != topicSlug {
		t.Fatalf("topic = %q, want %q", result.Topic, topicSlug)
	}
	if result.SourceType != models.SourceKindArticle {
		t.Fatalf("source type = %q, want %q", result.SourceType, models.SourceKindArticle)
	}
	if result.FilePath != wantFilePath {
		t.Fatalf("file path = %q, want %q", result.FilePath, wantFilePath)
	}
	if result.Title != "Latency Budget" {
		t.Fatalf("title = %q, want Latency Budget", result.Title)
	}

	documentPath := filepath.Join(vaultPath, filepath.FromSlash(wantFilePath))
	values, body := parseMarkdownFile(t, documentPath)
	assertFrontmatter(t, values, map[string]any{
		"title":       "Latency Budget",
		"type":        "source",
		"stage":       "raw",
		"domain":      "systems",
		"source_kind": "article",
		"scraped":     "2026-04-11",
		"source_path": sourcePath,
		"tags":        []string{"systems", "raw", "article"},
	})
	if body != "# Latency Budget\n\nKeep service latency inside the SLO.\n" {
		t.Fatalf("body = %q", body)
	}

	logContent := readFile(t, filepath.Join(vaultPath, topicSlug, "log.md"))
	if !strings.Contains(logContent, "## [2026-04-11] ingest | latency-budget.md (article)") {
		t.Fatalf("log.md missing ingest entry:\n%s", logContent)
	}
}

func TestIngestReturnsErrorWhenTopicDoesNotExist(t *testing.T) {
	t.Parallel()

	vaultPath := t.TempDir()

	_, err := Ingest(context.Background(), Options{
		VaultPath:  vaultPath,
		Topic:      "missing-topic",
		SourceKind: models.SourceKindArticle,
		Title:      "Missing Topic",
		Markdown:   "# Missing Topic\n",
		ScrapedAt:  fixedScrapeTime,
	})
	if err == nil {
		t.Fatal("expected Ingest to fail")
	}
	if !strings.Contains(err.Error(), "validate topic") {
		t.Fatalf("error = %q, want topic validation context", err)
	}
}

func TestIngestGeneratesUniqueSlugWhenFileExists(t *testing.T) {
	t.Parallel()

	vaultPath, topicSlug := scaffoldTopic(t)
	topicPath := filepath.Join(vaultPath, topicSlug)

	writeSourceFile(t, filepath.Join(topicPath, "raw", "articles", "release-notes.md"), "existing\n")
	writeSourceFile(t, filepath.Join(topicPath, "raw", "articles", "release-notes-2.md"), "existing\n")

	result, err := Ingest(context.Background(), Options{
		VaultPath:  vaultPath,
		Topic:      topicSlug,
		SourceKind: models.SourceKindArticle,
		Title:      "Release Notes",
		Markdown:   "# Release Notes\n\nFresh content.\n",
		ScrapedAt:  fixedScrapeTime,
	})
	if err != nil {
		t.Fatalf("Ingest returned error: %v", err)
	}

	if result.FilePath != "systems-design/raw/articles/release-notes-3.md" {
		t.Fatalf("file path = %q, want release-notes-3 path", result.FilePath)
	}
	assertFileExists(t, filepath.Join(vaultPath, filepath.FromSlash(result.FilePath)))
}

func TestIngestWritesExpectedSubdirectories(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		kind     models.SourceKind
		title    string
		wantPath string
	}{
		{
			name:     "article",
			kind:     models.SourceKindArticle,
			title:    "Architectural Article",
			wantPath: "systems-design/raw/articles/architectural-article.md",
		},
		{
			name:     "document",
			kind:     models.SourceKindDocument,
			title:    "Whitepaper Upload",
			wantPath: "systems-design/raw/articles/whitepaper-upload.md",
		},
		{
			name:     "youtube transcript",
			kind:     models.SourceKindYouTubeTranscript,
			title:    "System Design Interview",
			wantPath: "systems-design/raw/youtube/system-design-interview.md",
		},
		{
			name:     "github readme",
			kind:     models.SourceKindGitHubREADME,
			title:    "Awesome Repo README",
			wantPath: "systems-design/raw/github/awesome-repo-readme.md",
		},
		{
			name:     "bookmark cluster",
			kind:     models.SourceKindBookmarkCluster,
			title:    "April Links",
			wantPath: "systems-design/raw/bookmarks/april-links.md",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			vaultPath, topicSlug := scaffoldTopic(t)

			result, err := Ingest(context.Background(), Options{
				VaultPath:  vaultPath,
				Topic:      topicSlug,
				SourceKind: tc.kind,
				Title:      tc.title,
				Markdown:   "# " + tc.title + "\n\ncontent\n",
				ScrapedAt:  fixedScrapeTime,
			})
			if err != nil {
				t.Fatalf("Ingest returned error: %v", err)
			}
			if result.FilePath != tc.wantPath {
				t.Fatalf("file path = %q, want %q", result.FilePath, tc.wantPath)
			}

			values, _ := parseMarkdownFile(t, filepath.Join(vaultPath, filepath.FromSlash(tc.wantPath)))
			if got := values["source_kind"]; got != string(tc.kind) {
				t.Fatalf("source_kind = %#v, want %q", got, tc.kind)
			}
		})
	}
}

func TestIngestWithPreConvertedContentSkipsRegistryAndWritesDirectly(t *testing.T) {
	t.Parallel()

	vaultPath, topicSlug := scaffoldTopic(t)
	registry := &stubRegistry{}

	result, err := Ingest(context.Background(), Options{
		VaultPath:  vaultPath,
		Topic:      topicSlug,
		SourceKind: models.SourceKindGitHubREADME,
		SourceURL:  "https://github.com/example/repo",
		Title:      "Example Repository",
		Markdown:   "# Example Repository\n\nREADME snapshot.\n",
		Registry:   registry,
		ScrapedAt:  fixedScrapeTime,
	})
	if err != nil {
		t.Fatalf("Ingest returned error: %v", err)
	}
	if registry.calls != 0 {
		t.Fatalf("registry calls = %d, want 0", registry.calls)
	}
	if result.FilePath != "systems-design/raw/github/example-repository.md" {
		t.Fatalf("file path = %q", result.FilePath)
	}

	values, body := parseMarkdownFile(t, filepath.Join(vaultPath, filepath.FromSlash(result.FilePath)))
	if got := values["source_url"]; got != "https://github.com/example/repo" {
		t.Fatalf("source_url = %#v, want GitHub URL", got)
	}
	if body != "# Example Repository\n\nREADME snapshot.\n" {
		t.Fatalf("body = %q", body)
	}
}

func TestIngestEndToEndWithScaffoldedTopic(t *testing.T) {
	t.Parallel()

	vaultPath, topicSlug := scaffoldTopic(t)
	sourcePath := writeSourceFile(t, filepath.Join(t.TempDir(), "draft.md"), "# Queueing Theory\n\nLittle's Law applies.\n")

	result, err := Ingest(context.Background(), Options{
		VaultPath:  vaultPath,
		Topic:      topicSlug,
		SourceKind: models.SourceKindDocument,
		SourcePath: sourcePath,
		ScrapedAt:  fixedScrapeTime,
	})
	if err != nil {
		t.Fatalf("Ingest returned error: %v", err)
	}

	if result.FilePath != "systems-design/raw/articles/queueing-theory.md" {
		t.Fatalf("file path = %q", result.FilePath)
	}

	values, body := parseMarkdownFile(t, filepath.Join(vaultPath, filepath.FromSlash(result.FilePath)))
	assertFrontmatter(t, values, map[string]any{
		"title":       "Queueing Theory",
		"type":        "source",
		"stage":       "raw",
		"domain":      "systems",
		"source_kind": "document",
		"scraped":     "2026-04-11",
		"source_path": sourcePath,
		"tags":        []string{"systems", "raw", "document"},
	})
	if body != "# Queueing Theory\n\nLittle's Law applies.\n" {
		t.Fatalf("body = %q", body)
	}

	logContent := readFile(t, filepath.Join(vaultPath, topicSlug, "log.md"))
	for _, fragment := range []string{
		"## [2026-04-11] ingest | queueing-theory.md (document)",
		"Ingested `Queueing Theory` into `systems-design/raw/articles/queueing-theory.md`.",
	} {
		if !strings.Contains(logContent, fragment) {
			t.Fatalf("log.md missing %q:\n%s", fragment, logContent)
		}
	}
}

func TestResolveMarkdownValidationAndRegistryFailures(t *testing.T) {
	t.Parallel()

	t.Run("requires source path or markdown", func(t *testing.T) {
		t.Parallel()

		_, _, err := resolveMarkdown(context.Background(), Options{SourceKind: models.SourceKindArticle})
		if err == nil || !strings.Contains(err.Error(), "source path or markdown content is required") {
			t.Fatalf("error = %v, want missing input validation", err)
		}
	})

	t.Run("rejects nil converter result", func(t *testing.T) {
		t.Parallel()

		sourcePath := writeSourceFile(t, filepath.Join(t.TempDir(), "source.txt"), "content")
		_, _, err := resolveMarkdown(context.Background(), Options{
			SourceKind: models.SourceKindArticle,
			SourcePath: sourcePath,
			Registry:   &stubRegistry{},
		})
		if err == nil || !strings.Contains(err.Error(), "converter returned nil result") {
			t.Fatalf("error = %v, want nil result validation", err)
		}
	})

	t.Run("rejects empty converter markdown", func(t *testing.T) {
		t.Parallel()

		sourcePath := writeSourceFile(t, filepath.Join(t.TempDir(), "source.txt"), "content")
		_, _, err := resolveMarkdown(context.Background(), Options{
			SourceKind: models.SourceKindArticle,
			SourcePath: sourcePath,
			Registry: &stubRegistry{
				result: &models.ConvertResult{Title: "Empty", Markdown: " \n\t"},
			},
		})
		if err == nil || !strings.Contains(err.Error(), "markdown output is empty") {
			t.Fatalf("error = %v, want empty markdown validation", err)
		}
	})

	t.Run("clones convert options before registry call", func(t *testing.T) {
		t.Parallel()

		sourcePath := writeSourceFile(t, filepath.Join(t.TempDir(), "source.txt"), "content")
		registry := &stubRegistry{
			result: &models.ConvertResult{Title: "Converted", Markdown: "# Converted\n"},
		}
		options := map[string]any{"mimeType": "text/plain"}

		_, _, err := resolveMarkdown(context.Background(), Options{
			SourceKind:     models.SourceKindArticle,
			SourcePath:     sourcePath,
			Registry:       registry,
			ConvertOptions: options,
		})
		if err != nil {
			t.Fatalf("resolveMarkdown returned error: %v", err)
		}
		if !reflect.DeepEqual(registry.lastInput.Options, options) {
			t.Fatalf("registry options = %#v, want %#v", registry.lastInput.Options, options)
		}
		options["mimeType"] = "application/json"
		if registry.lastInput.Options["mimeType"] != "text/plain" {
			t.Fatalf("registry options mutated with caller map: %#v", registry.lastInput.Options)
		}
	})
}

func TestRawDirectoryForSourceKind(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		sourceKind  models.SourceKind
		want        string
		wantErrText string
	}{
		{name: "article", sourceKind: models.SourceKindArticle, want: "articles"},
		{name: "document", sourceKind: models.SourceKindDocument, want: "articles"},
		{name: "github", sourceKind: models.SourceKindGitHubREADME, want: "github"},
		{name: "youtube", sourceKind: models.SourceKindYouTubeTranscript, want: "youtube"},
		{name: "bookmarks", sourceKind: models.SourceKindBookmarkCluster, want: "bookmarks"},
		{name: "codebase file", sourceKind: models.SourceKindCodebaseFile, want: "codebase/files"},
		{name: "codebase symbol", sourceKind: models.SourceKindCodebaseSymbol, want: "codebase/symbols"},
		{name: "missing kind", sourceKind: "", wantErrText: "source kind is required"},
		{name: "unsupported kind", sourceKind: models.SourceKind("podcast"), wantErrText: "unsupported source kind"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := rawDirectoryForSourceKind(tc.sourceKind)
			if tc.wantErrText != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErrText) {
					t.Fatalf("error = %v, want %q", err, tc.wantErrText)
				}
				return
			}
			if err != nil {
				t.Fatalf("rawDirectoryForSourceKind returned error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("directory = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestTitleHelpers(t *testing.T) {
	t.Parallel()

	if got := deriveTitle("Explicit Title", "", "", models.SourceKindArticle); got != "Explicit Title" {
		t.Fatalf("deriveTitle explicit = %q, want Explicit Title", got)
	}
	if got := deriveTitle("", "/tmp/distributed-systems-notes.md", "", models.SourceKindArticle); got != "Distributed Systems Notes" {
		t.Fatalf("deriveTitle path fallback = %q, want Distributed Systems Notes", got)
	}
	if got := deriveTitle("", "", "https://example.com/articles/queueing-theory", models.SourceKindArticle); got != "Queueing Theory" {
		t.Fatalf("deriveTitle URL fallback = %q, want Queueing Theory", got)
	}
	if got := deriveTitle("", "", "https://example.com/", models.SourceKindArticle); got != "Example Com" {
		t.Fatalf("deriveTitle host fallback = %q, want Example Com", got)
	}
	if got := deriveTitle("", "", "", models.SourceKindBookmarkCluster); got != "Bookmark Cluster" {
		t.Fatalf("deriveTitle kind fallback = %q, want Bookmark Cluster", got)
	}

	if got := humanizeSegment("  "); got != "" {
		t.Fatalf("humanizeSegment blank = %q, want empty string", got)
	}
	if got := humanizeSegment("latency_slo-v2"); got != "Latency Slo V2" {
		t.Fatalf("humanizeSegment = %q, want Latency Slo V2", got)
	}
	if got := firstNonEmpty("", "  ", "value", "other"); got != "value" {
		t.Fatalf("firstNonEmpty = %q, want value", got)
	}
}

func TestAppendLogEntryRequiresExistingLogPath(t *testing.T) {
	t.Parallel()

	logPath := filepath.Join(t.TempDir(), "missing", "log.md")
	err := appendLogEntry(logPath, fixedScrapeTime, "example", models.SourceKindArticle, "topic/raw/articles/example.md", "Example")
	if err == nil || !strings.Contains(err.Error(), "open") {
		t.Fatalf("error = %v, want open failure", err)
	}
}

type stubRegistry struct {
	result    *models.ConvertResult
	err       error
	calls     int
	lastInput models.ConvertInput
}

func (registry *stubRegistry) Convert(_ context.Context, input models.ConvertInput) (*models.ConvertResult, error) {
	registry.calls++
	registry.lastInput = input
	if registry.err != nil {
		return nil, registry.err
	}
	return registry.result, nil
}

func scaffoldTopic(t *testing.T) (string, string) {
	t.Helper()

	vaultPath := t.TempDir()
	info, err := topic.New(vaultPath, "systems-design", "Systems Design", "systems")
	if err != nil {
		t.Fatalf("topic.New returned error: %v", err)
	}

	return vaultPath, info.Slug
}

func writeSourceFile(t *testing.T, path string, content string) string {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create source dir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file %q: %v", path, err)
	}

	return path
}

func parseMarkdownFile(t *testing.T, filePath string) (map[string]any, string) {
	t.Helper()

	content := readFile(t, filePath)
	values, body, err := frontmatter.Parse(content)
	if err != nil {
		t.Fatalf("parse frontmatter %q: %v", filePath, err)
	}

	return values, body
}

func readFile(t *testing.T, filePath string) string {
	t.Helper()

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("read file %q: %v", filePath, err)
	}

	return string(content)
}

func assertFileExists(t *testing.T, filePath string) {
	t.Helper()

	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("stat %q: %v", filePath, err)
	}
	if info.IsDir() {
		t.Fatalf("%q is a directory, want file", filePath)
	}
}

func assertFrontmatter(t *testing.T, got map[string]any, want map[string]any) {
	t.Helper()

	for key, wantValue := range want {
		gotValue, exists := got[key]
		if !exists {
			t.Fatalf("frontmatter missing key %q", key)
		}

		switch typedWant := wantValue.(type) {
		case []string:
			gotSlice := frontmatter.GetStringSlice(got, key)
			if !reflect.DeepEqual(gotSlice, typedWant) {
				t.Fatalf("frontmatter[%q] = %#v, want %#v", key, gotSlice, typedWant)
			}
		default:
			if gotValue != typedWant {
				t.Fatalf("frontmatter[%q] = %#v, want %#v", key, gotValue, typedWant)
			}
		}
	}
}
