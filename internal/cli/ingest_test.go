package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	kconfig "github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/firecrawl"
	kgenerate "github.com/compozy/kb/internal/generate"
	kingest "github.com/compozy/kb/internal/ingest"
	"github.com/compozy/kb/internal/models"
	ktopic "github.com/compozy/kb/internal/topic"
	"github.com/compozy/kb/internal/youtube"
)

func TestIngestParentHelpListsSubcommands(t *testing.T) {
	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"ingest", "--help"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	for _, fragment := range []string{"url", "file", "youtube", "codebase", "bookmarks", "--vault"} {
		if !strings.Contains(stdout.String(), fragment) {
			t.Fatalf("expected help output to contain %q, got:\n%s", fragment, stdout.String())
		}
	}
}

func TestIngestCodebaseHelpIncludesSupportedLanguagesAndDryRun(t *testing.T) {
	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"ingest", "codebase", "--help"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	for _, fragment := range []string{supportedCodebaseLanguagesHelp(), "java", "--dry-run"} {
		if !strings.Contains(stdout.String(), fragment) {
			t.Fatalf("expected help output to contain %q, got:\n%s", fragment, stdout.String())
		}
	}
}

func TestIngestCommandsRequireTopicFlag(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
	}{
		{name: "url", args: []string{"ingest", "url", "https://example.com"}},
		{name: "file", args: []string{"ingest", "file", "/tmp/source.md"}},
		{name: "youtube", args: []string{"ingest", "youtube", "https://youtu.be/abcdefghijk"}},
		{name: "codebase", args: []string{"ingest", "codebase", "/tmp/repo"}},
		{name: "bookmarks", args: []string{"ingest", "bookmarks", "/tmp/bookmarks.md"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			command := newRootCommand()
			command.SetOut(new(bytes.Buffer))
			command.SetErr(new(bytes.Buffer))
			command.SetArgs(tt.args)

			err := command.ExecuteContext(context.Background())
			if err == nil {
				t.Fatal("expected missing topic flag error")
			}
			if !strings.Contains(err.Error(), `required flag(s) "topic" not set`) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestIngestCommandsRequirePositionalArg(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
	}{
		{name: "url", args: []string{"ingest", "url", "--topic", "systems-design"}},
		{name: "file", args: []string{"ingest", "file", "--topic", "systems-design"}},
		{name: "youtube", args: []string{"ingest", "youtube", "--topic", "systems-design"}},
		{name: "codebase", args: []string{"ingest", "codebase", "--topic", "systems-design"}},
		{name: "bookmarks", args: []string{"ingest", "bookmarks", "--topic", "systems-design"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			command := newRootCommand()
			command.SetOut(new(bytes.Buffer))
			command.SetErr(new(bytes.Buffer))
			command.SetArgs(tt.args)

			err := command.ExecuteContext(context.Background())
			if err == nil {
				t.Fatal("expected missing positional argument error")
			}
			if !strings.Contains(err.Error(), "accepts 1 arg(s)") {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestIngestFileCommandReturnsErrorForMissingFile(t *testing.T) {
	restoreIngestGlobals(t)

	runIngestTopicInfo = func(vaultPath, slug string) (models.TopicInfo, error) {
		return models.TopicInfo{Slug: slug, Title: "Systems Design", Domain: "systems"}, nil
	}

	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"ingest", "file", filepath.Join(t.TempDir(), "missing.md"), "--topic", "systems-design", "--vault", "/tmp/vault"})

	err := command.ExecuteContext(context.Background())
	if err == nil {
		t.Fatal("expected missing file error")
	}
	if !strings.Contains(err.Error(), "stat source path") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIngestURLCommandScrapesAndWritesJSON(t *testing.T) {
	restoreIngestGlobals(t)

	var gotFirecrawlConfig kconfig.FirecrawlConfig
	var gotScrapeURL string
	var gotIngest kingest.Options
	var gotTopicVault string
	var gotTopicSlug string

	loadIngestConfig = func() (kconfig.Config, error) {
		return kconfig.Config{
			Firecrawl: kconfig.FirecrawlConfig{
				APIKey: "firecrawl-key",
				APIURL: "https://firecrawl.test",
			},
		}, nil
	}
	newFirecrawlScraper = func(cfg kconfig.FirecrawlConfig) firecrawlScraper {
		gotFirecrawlConfig = cfg
		return fakeFirecrawlScraper{
			scrape: func(ctx context.Context, sourceURL string) (*firecrawl.ScrapeResult, error) {
				gotScrapeURL = sourceURL
				return &firecrawl.ScrapeResult{
					Markdown:  "# Latency Budget\n\nKeep the service fast.\n",
					Title:     "Latency Budget",
					SourceURL: "https://example.com/latency-budget",
				}, nil
			},
		}
	}
	runIngestTopicInfo = func(vaultPath, slug string) (models.TopicInfo, error) {
		gotTopicVault = vaultPath
		gotTopicSlug = slug
		return models.TopicInfo{
			Slug:     slug,
			Title:    "Systems Design",
			Domain:   "systems",
			RootPath: filepath.Join(vaultPath, slug),
		}, nil
	}
	runIngest = func(ctx context.Context, options kingest.Options) (models.IngestResult, error) {
		gotIngest = options
		return models.IngestResult{
			Topic:      options.Topic,
			SourceType: options.SourceKind,
			FilePath:   "systems-design/raw/articles/latency-budget.md",
			Title:      "Latency Budget",
		}, nil
	}

	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"ingest", "url", "https://example.com/latency-budget", "--topic", "systems-design", "--vault", "/tmp/vault"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	if gotFirecrawlConfig.APIKey != "firecrawl-key" || gotFirecrawlConfig.APIURL != "https://firecrawl.test" {
		t.Fatalf("firecrawl config = %#v", gotFirecrawlConfig)
	}
	if gotScrapeURL != "https://example.com/latency-budget" {
		t.Fatalf("scrape URL = %q, want source URL", gotScrapeURL)
	}
	if gotTopicVault != "/tmp/vault" || gotTopicSlug != "systems-design" {
		t.Fatalf("topic lookup = (%q, %q), want (/tmp/vault, systems-design)", gotTopicVault, gotTopicSlug)
	}
	if gotIngest.VaultPath != "/tmp/vault" {
		t.Fatalf("ingest vault path = %q, want /tmp/vault", gotIngest.VaultPath)
	}
	if gotIngest.Topic != "systems-design" {
		t.Fatalf("ingest topic = %q, want systems-design", gotIngest.Topic)
	}
	if gotIngest.SourceKind != models.SourceKindArticle {
		t.Fatalf("ingest source kind = %q, want %q", gotIngest.SourceKind, models.SourceKindArticle)
	}
	if gotIngest.SourceURL != "https://example.com/latency-budget" {
		t.Fatalf("ingest source URL = %q, want canonical source URL", gotIngest.SourceURL)
	}
	if gotIngest.Title != "Latency Budget" {
		t.Fatalf("ingest title = %q, want Latency Budget", gotIngest.Title)
	}
	if gotIngest.Markdown != "# Latency Budget\n\nKeep the service fast.\n" {
		t.Fatalf("ingest markdown = %q", gotIngest.Markdown)
	}

	var result models.IngestResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("stdout did not contain JSON: %v\n%s", err, stdout.String())
	}
	if result.FilePath != "systems-design/raw/articles/latency-budget.md" || result.Title != "Latency Budget" {
		t.Fatalf("unexpected result payload: %#v", result)
	}
}

func TestIngestFileCommandRoutesToOrchestrator(t *testing.T) {
	restoreIngestGlobals(t)

	sourcePath := filepath.Join(t.TempDir(), "whitepaper.md")
	if err := os.WriteFile(sourcePath, []byte("# Whitepaper\n"), 0o644); err != nil {
		t.Fatalf("write source file: %v", err)
	}

	var gotIngest kingest.Options
	runIngestTopicInfo = func(vaultPath, slug string) (models.TopicInfo, error) {
		return models.TopicInfo{Slug: slug, Title: "Systems Design", Domain: "systems"}, nil
	}
	runIngest = func(ctx context.Context, options kingest.Options) (models.IngestResult, error) {
		gotIngest = options
		return models.IngestResult{
			Topic:      options.Topic,
			SourceType: options.SourceKind,
			FilePath:   "systems-design/raw/articles/whitepaper.md",
			Title:      "Whitepaper",
		}, nil
	}

	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"ingest", "file", sourcePath, "--topic", "systems-design", "--vault", "/tmp/vault"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	if gotIngest.SourceKind != models.SourceKindDocument {
		t.Fatalf("source kind = %q, want %q", gotIngest.SourceKind, models.SourceKindDocument)
	}
	if gotIngest.SourcePath != sourcePath {
		t.Fatalf("source path = %q, want %q", gotIngest.SourcePath, sourcePath)
	}
	if gotIngest.Registry == nil {
		t.Fatal("expected converter registry to be provided")
	}

	var result models.IngestResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("stdout did not contain JSON: %v\n%s", err, stdout.String())
	}
	if result.SourceType != models.SourceKindDocument {
		t.Fatalf("unexpected result payload: %#v", result)
	}
}

func TestIngestYouTubeCommandAcceptsTranscribePolicy(t *testing.T) {
	restoreIngestGlobals(t)

	var gotExtractURL string
	var gotExtractOptions youtube.ExtractOptions
	var gotIngest kingest.Options
	viewCount := int64(3271)
	likeCount := int64(77)
	commentCount := int64(11)
	channelFollowers := int64(20000)
	wasLive := false
	wantYouTubeConfig := kconfig.YouTubeConfig{
		YTDLPPath:     "custom-yt-dlp",
		Proxy:         "http://proxy.internal:8080",
		CookiesFile:   "/tmp/youtube-cookies.txt",
		UserAgent:     "kb-test-agent",
		RetryAttempts: 4,
		RetryBackoff:  "250ms",
	}

	loadIngestConfig = func() (kconfig.Config, error) {
		return kconfig.Config{
			OpenRouter: kconfig.OpenRouterConfig{
				APIKey:   "openrouter-key",
				APIURL:   "https://openrouter.test",
				STTModel: "demo-stt",
			},
			YouTube: wantYouTubeConfig,
		}, nil
	}
	newYouTubeTranscriptExtractor = func(cfg kconfig.Config) youtubeTranscriptExtractor {
		if cfg.OpenRouter.APIKey != "openrouter-key" {
			t.Fatalf("openrouter config not passed to extractor: %#v", cfg.OpenRouter)
		}
		if !reflect.DeepEqual(cfg.YouTube, wantYouTubeConfig) {
			t.Fatalf("youtube config = %#v, want %#v", cfg.YouTube, wantYouTubeConfig)
		}
		return fakeYouTubeExtractor{
			extract: func(ctx context.Context, rawURL string, options youtube.ExtractOptions) (*youtube.Result, error) {
				gotExtractURL = rawURL
				gotExtractOptions = options
				return &youtube.Result{
					Metadata: youtube.Metadata{
						URL:                  "https://www.youtube.com/watch?v=abcdefghijk",
						Title:                "Queueing Theory Deep Dive",
						Channel:              "Changelog",
						ChannelID:            "UCZbTest",
						UploaderID:           "@Changelog",
						Duration:             6441 * time.Second,
						DurationString:       "1:47:21",
						PublishDate:          time.Date(2025, time.July, 2, 0, 0, 0, 0, time.UTC),
						ViewCount:            &viewCount,
						LikeCount:            &likeCount,
						CommentCount:         &commentCount,
						ChannelFollowerCount: &channelFollowers,
						Categories:           []string{"Science & Technology"},
						VideoTags:            []string{"go", "systems"},
						Language:             "en",
						LiveStatus:           "not_live",
						WasLive:              &wasLive,
						ChapterCount:         17,
					},
					Markdown: "# Queueing Theory Deep Dive\n\nTranscript.\n",
				}, nil
			},
		}
	}
	runIngestTopicInfo = func(vaultPath, slug string) (models.TopicInfo, error) {
		return models.TopicInfo{Slug: slug, Title: "Systems Design", Domain: "systems"}, nil
	}
	runIngest = func(ctx context.Context, options kingest.Options) (models.IngestResult, error) {
		gotIngest = options
		return models.IngestResult{
			Topic:      options.Topic,
			SourceType: options.SourceKind,
			FilePath:   "systems-design/raw/youtube/queueing-theory-deep-dive.md",
			Title:      "Queueing Theory Deep Dive",
		}, nil
	}

	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{
		"ingest", "youtube", "https://youtu.be/abcdefghijk",
		"--topic", "systems-design",
		"--vault", "/tmp/vault",
		"--transcribe", "auto",
	})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	if gotExtractURL != "https://youtu.be/abcdefghijk" {
		t.Fatalf("extract URL = %q, want source URL", gotExtractURL)
	}
	if gotExtractOptions.TranscriptionPolicy != youtube.TranscriptionPolicyAuto {
		t.Fatalf("transcription policy = %q, want auto", gotExtractOptions.TranscriptionPolicy)
	}
	if gotIngest.SourceKind != models.SourceKindYouTubeTranscript {
		t.Fatalf("ingest source kind = %q, want %q", gotIngest.SourceKind, models.SourceKindYouTubeTranscript)
	}
	if gotIngest.SourceURL != "https://www.youtube.com/watch?v=abcdefghijk" {
		t.Fatalf("ingest source URL = %q", gotIngest.SourceURL)
	}
	if gotIngest.Title != "Queueing Theory Deep Dive" {
		t.Fatalf("ingest title = %q", gotIngest.Title)
	}
	if got := gotIngest.ExtraFrontmatter["view_count"]; got != viewCount {
		t.Fatalf("view_count = %#v, want %d", got, viewCount)
	}
	if got := gotIngest.ExtraFrontmatter["upload_date"]; got != "2025-07-02" {
		t.Fatalf("upload_date = %#v, want 2025-07-02", got)
	}
	if got := gotIngest.ExtraFrontmatter["duration"]; got != int64(6441) {
		t.Fatalf("duration = %#v, want 6441", got)
	}
	if got := gotIngest.ExtraFrontmatter["duration_string"]; got != "1:47:21" {
		t.Fatalf("duration_string = %#v, want 1:47:21", got)
	}
	if got := gotIngest.ExtraFrontmatter["channel"]; got != "Changelog" {
		t.Fatalf("channel = %#v, want Changelog", got)
	}
	if got := gotIngest.ExtraFrontmatter["channel_id"]; got != "UCZbTest" {
		t.Fatalf("channel_id = %#v, want UCZbTest", got)
	}
	if got := gotIngest.ExtraFrontmatter["uploader_id"]; got != "@Changelog" {
		t.Fatalf("uploader_id = %#v, want @Changelog", got)
	}
	if got := gotIngest.ExtraFrontmatter["channel_follower_count"]; got != channelFollowers {
		t.Fatalf("channel_follower_count = %#v, want %d", got, channelFollowers)
	}
	if got := gotIngest.ExtraFrontmatter["youtube_tags"]; !reflect.DeepEqual(got, []string{"go", "systems"}) {
		t.Fatalf("youtube_tags = %#v", got)
	}
	if _, exists := gotIngest.ExtraFrontmatter["tags"]; exists {
		t.Fatalf("frontmatter should not override reserved tags: %#v", gotIngest.ExtraFrontmatter)
	}

	var result models.IngestResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("stdout did not contain JSON: %v\n%s", err, stdout.String())
	}
	if result.SourceType != models.SourceKindYouTubeTranscript {
		t.Fatalf("unexpected result payload: %#v", result)
	}
}

func TestYouTubeFrontmatterIncludesVideoMetricsAndTranscriptProvenance(t *testing.T) {
	viewCount := int64(3271)
	likeCount := int64(77)
	commentCount := int64(11)
	channelFollowers := int64(20000)
	wasLive := false
	values := youtubeFrontmatter(&youtube.Result{
		Metadata: youtube.Metadata{
			Channel:              "Changelog",
			ChannelID:            "UCZbTest",
			UploaderID:           "@Changelog",
			Duration:             6441 * time.Second,
			DurationString:       "1:47:21",
			PublishDate:          time.Date(2025, time.July, 2, 0, 0, 0, 0, time.UTC),
			ViewCount:            &viewCount,
			LikeCount:            &likeCount,
			CommentCount:         &commentCount,
			ChannelFollowerCount: &channelFollowers,
			Categories:           []string{"Science & Technology"},
			VideoTags:            []string{"go", "systems"},
			Language:             "en",
			LiveStatus:           "not_live",
			WasLive:              &wasLive,
			ChapterCount:         17,
		},
		Source:              youtube.TranscriptSourceSTT,
		TranscriptionPolicy: youtube.TranscriptionPolicyAuto,
		Language:            "en-US",
		CaptionKind:         youtube.CaptionKindAutomatic,
		STTProvider:         "openai",
		STTModel:            "gpt-4o-transcribe",
	})
	want := map[string]any{
		"view_count":             viewCount,
		"like_count":             likeCount,
		"comment_count":          commentCount,
		"upload_date":            "2025-07-02",
		"duration":               int64(6441),
		"duration_string":        "1:47:21",
		"channel":                "Changelog",
		"channel_id":             "UCZbTest",
		"uploader_id":            "@Changelog",
		"channel_follower_count": channelFollowers,
		"categories":             []string{"Science & Technology"},
		"youtube_tags":           []string{"go", "systems"},
		"language":               "en",
		"live_status":            "not_live",
		"was_live":               false,
		"chapter_count":          17,
		"transcript_source":      "stt",
		"transcription_policy":   "auto",
		"transcript_language":    "en-US",
		"caption_kind":           "automatic",
		"stt_provider":           "openai",
		"stt_model":              "gpt-4o-transcribe",
	}
	if !reflect.DeepEqual(values, want) {
		t.Fatalf("frontmatter = %#v, want %#v", values, want)
	}
	if _, exists := values["tags"]; exists {
		t.Fatalf("frontmatter should not include reserved tags: %#v", values)
	}
}

func TestIngestYouTubeCommandRejectsRemovedSTTFlag(t *testing.T) {
	restoreIngestGlobals(t)

	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	stderr := new(bytes.Buffer)
	command.SetErr(stderr)
	command.SetArgs([]string{
		"ingest", "youtube", "https://youtu.be/abcdefghijk",
		"--topic", "systems-design",
		"--vault", "/tmp/vault",
		"--stt",
	})

	err := command.ExecuteContext(context.Background())
	if err == nil || !strings.Contains(err.Error(), "unknown flag: --stt") {
		t.Fatalf("error = %v, want unknown --stt flag", err)
	}
}

func TestIngestCodebaseCommandPassesGenerateFlags(t *testing.T) {
	restoreIngestGlobals(t)

	var gotGenerate models.GenerateOptions
	runIngestTopicInfo = func(vaultPath, slug string) (models.TopicInfo, error) {
		return models.TopicInfo{
			Slug:     slug,
			Title:    "Systems Design",
			Domain:   "systems",
			RootPath: filepath.Join(vaultPath, slug),
		}, nil
	}
	runGenerate = func(ctx context.Context, opts models.GenerateOptions, observer kgenerate.Observer) (models.GenerationSummary, error) {
		gotGenerate = opts
		return models.GenerationSummary{
			Command:               "generate",
			TopicSlug:             opts.TopicSlug,
			VaultPath:             opts.VaultPath,
			TopicPath:             filepath.Join(opts.VaultPath, opts.TopicSlug),
			FilesScanned:          2,
			FilesParsed:           2,
			RawDocumentsWritten:   9,
			WikiDocumentsWritten:  10,
			IndexDocumentsWritten: 3,
		}, nil
	}

	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{
		"ingest", "codebase", "/tmp/repo",
		"--topic", "systems-design",
		"--vault", "/tmp/vault",
		"--include", "src/**/*.go",
		"--include", "web/**/*.ts",
		"--exclude", "vendor/**",
		"--dry-run",
		"--semantic",
		"--progress", "never",
	})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	expected := models.GenerateOptions{
		RootPath:        "/tmp/repo",
		VaultPath:       "/tmp/vault",
		TopicSlug:       "systems-design",
		Title:           "Systems Design",
		Domain:          "systems",
		IncludePatterns: []string{"src/**/*.go", "web/**/*.ts"},
		ExcludePatterns: []string{"vendor/**"},
		DryRun:          true,
		Semantic:        true,
	}
	if !reflect.DeepEqual(gotGenerate, expected) {
		t.Fatalf("generate options = %#v, want %#v", gotGenerate, expected)
	}

	var result codebaseIngestResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("stdout did not contain JSON: %v\n%s", err, stdout.String())
	}
	if result.Topic != "systems-design" || result.FilePath != "systems-design/raw/codebase" {
		t.Fatalf("unexpected result payload: %#v", result)
	}
	if result.Summary.TopicSlug != "systems-design" || result.Summary.RawDocumentsWritten != 9 {
		t.Fatalf("unexpected summary payload: %#v", result.Summary)
	}
}

func TestIngestCodebaseCommandJSONContractRequiredKeysFullRun(t *testing.T) {
	restoreIngestGlobals(t)

	runIngestTopicInfo = func(vaultPath, slug string) (models.TopicInfo, error) {
		return models.TopicInfo{
			Slug:     slug,
			Title:    "Systems Design",
			Domain:   "systems",
			RootPath: filepath.Join(vaultPath, slug),
		}, nil
	}
	runGenerate = func(ctx context.Context, opts models.GenerateOptions, observer kgenerate.Observer) (models.GenerationSummary, error) {
		return models.GenerationSummary{
			Command:               "ingest codebase",
			RootPath:              opts.RootPath,
			VaultPath:             opts.VaultPath,
			TopicPath:             filepath.Join(opts.VaultPath, opts.TopicSlug),
			TopicSlug:             opts.TopicSlug,
			DryRun:                false,
			DetectedLanguages:     []string{"java"},
			SelectedAdapters:      []string{"adapter.JavaAdapter"},
			FilesScanned:          6,
			FilesParsed:           6,
			FilesSkipped:          0,
			SymbolsExtracted:      10,
			RelationsEmitted:      8,
			RawDocumentsWritten:   12,
			WikiDocumentsWritten:  10,
			IndexDocumentsWritten: 3,
			Timings: models.GenerationTimings{
				ScanMillis:           1,
				SelectAdaptersMillis: 1,
				ParseMillis:          1,
				NormalizeMillis:      1,
				MetricsMillis:        1,
				RenderMillis:         1,
				WriteMillis:          1,
				TotalMillis:          8,
			},
			Diagnostics: []models.StructuredDiagnostic{},
		}, nil
	}

	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{
		"ingest", "codebase", "/tmp/repo",
		"--topic", "systems-design",
		"--vault", "/tmp/vault",
		"--progress", "never",
	})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	payload := decodeJSONMap(t, stdout.Bytes())
	assertCodebaseIngestContractShape(t, payload)
	assertCodebaseIngestContractSemantics(t, payload, false)
}

func TestIngestCodebaseCommandJSONContractRequiredKeysDryRun(t *testing.T) {
	restoreIngestGlobals(t)

	runIngestTopicInfo = func(vaultPath, slug string) (models.TopicInfo, error) {
		return models.TopicInfo{
			Slug:     slug,
			Title:    "Systems Design",
			Domain:   "systems",
			RootPath: filepath.Join(vaultPath, slug),
		}, nil
	}
	runGenerate = func(ctx context.Context, opts models.GenerateOptions, observer kgenerate.Observer) (models.GenerationSummary, error) {
		return models.GenerationSummary{
			Command:               "ingest codebase",
			RootPath:              opts.RootPath,
			VaultPath:             opts.VaultPath,
			TopicPath:             filepath.Join(opts.VaultPath, opts.TopicSlug),
			TopicSlug:             opts.TopicSlug,
			DryRun:                true,
			DetectedLanguages:     []string{"java"},
			SelectedAdapters:      []string{"adapter.JavaAdapter"},
			FilesScanned:          6,
			FilesParsed:           6,
			FilesSkipped:          0,
			SymbolsExtracted:      10,
			RelationsEmitted:      8,
			RawDocumentsWritten:   0,
			WikiDocumentsWritten:  0,
			IndexDocumentsWritten: 0,
			Timings: models.GenerationTimings{
				ScanMillis:           1,
				SelectAdaptersMillis: 1,
				ParseMillis:          1,
				NormalizeMillis:      1,
				MetricsMillis:        1,
				RenderMillis:         1,
				WriteMillis:          0,
				TotalMillis:          7,
			},
			Diagnostics: []models.StructuredDiagnostic{},
		}, nil
	}

	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{
		"ingest", "codebase", "/tmp/repo",
		"--topic", "systems-design",
		"--vault", "/tmp/vault",
		"--progress", "never",
		"--dry-run",
	})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	payload := decodeJSONMap(t, stdout.Bytes())
	assertCodebaseIngestContractShape(t, payload)
	assertCodebaseIngestContractSemantics(t, payload, true)
}

func TestIngestCodebaseCommandBootstrapsMissingTopicWithDefaultVault(t *testing.T) {
	restoreIngestGlobals(t)

	var gotGenerate models.GenerateOptions
	runIngestTopicInfo = func(vaultPath, slug string) (models.TopicInfo, error) {
		return models.TopicInfo{}, fmt.Errorf(
			"topic info: topic %q is missing the expected KB skeleton: %w",
			slug,
			ktopic.ErrTopicNotFound,
		)
	}
	runGenerate = func(ctx context.Context, opts models.GenerateOptions, observer kgenerate.Observer) (models.GenerationSummary, error) {
		gotGenerate = opts
		return models.GenerationSummary{
			Command:   "generate",
			TopicSlug: opts.TopicSlug,
			VaultPath: opts.VaultPath,
			TopicPath: filepath.Join(opts.VaultPath, opts.TopicSlug),
		}, nil
	}

	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{
		"ingest", "codebase", "/tmp/demo-repo",
		"--topic", "systems-design",
		"--progress", "never",
	})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	expected := models.GenerateOptions{
		RootPath:  "/tmp/demo-repo",
		VaultPath: "/tmp/demo-repo/.kb/vault",
		TopicSlug: "systems-design",
		Title:     "Systems Design",
		Domain:    "systems-design",
	}
	if !reflect.DeepEqual(gotGenerate, expected) {
		t.Fatalf("generate options = %#v, want %#v", gotGenerate, expected)
	}

	var result codebaseIngestResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("stdout did not contain JSON: %v\n%s", err, stdout.String())
	}
	if result.Topic != "systems-design" || result.Title != "Systems Design" {
		t.Fatalf("unexpected result payload: %#v", result)
	}
}

func TestIngestCodebaseCommandBootstrapAcceptsTitleAndDomain(t *testing.T) {
	restoreIngestGlobals(t)

	var gotGenerate models.GenerateOptions
	runIngestTopicInfo = func(vaultPath, slug string) (models.TopicInfo, error) {
		return models.TopicInfo{}, fmt.Errorf(
			"topic info: topic %q is missing the expected KB skeleton: %w",
			slug,
			ktopic.ErrTopicNotFound,
		)
	}
	runGenerate = func(ctx context.Context, opts models.GenerateOptions, observer kgenerate.Observer) (models.GenerationSummary, error) {
		gotGenerate = opts
		return models.GenerationSummary{
			Command:   "generate",
			TopicSlug: opts.TopicSlug,
			VaultPath: opts.VaultPath,
			TopicPath: filepath.Join(opts.VaultPath, opts.TopicSlug),
		}, nil
	}

	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{
		"ingest", "codebase", "/tmp/chat",
		"--topic", "chat-sdk",
		"--title", "Chat SDK",
		"--domain", "messaging",
		"--progress", "never",
	})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	expected := models.GenerateOptions{
		RootPath:  "/tmp/chat",
		VaultPath: "/tmp/chat/.kb/vault",
		TopicSlug: "chat-sdk",
		Title:     "Chat SDK",
		Domain:    "messaging",
	}
	if !reflect.DeepEqual(gotGenerate, expected) {
		t.Fatalf("generate options = %#v, want %#v", gotGenerate, expected)
	}
}

func TestIngestCodebaseCommandRejectsBootstrapMetadataForExistingTopic(t *testing.T) {
	restoreIngestGlobals(t)

	runIngestTopicInfo = func(vaultPath, slug string) (models.TopicInfo, error) {
		return models.TopicInfo{
			Slug:     slug,
			Title:    "Systems Design",
			Domain:   "systems",
			RootPath: filepath.Join(vaultPath, slug),
		}, nil
	}

	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{
		"ingest", "codebase", "/tmp/repo",
		"--topic", "systems-design",
		"--vault", "/tmp/vault",
		"--title", "Renamed Topic",
		"--progress", "never",
	})

	err := command.ExecuteContext(context.Background())
	if err == nil {
		t.Fatal("expected bootstrap metadata rejection")
	}
	if !strings.Contains(err.Error(), "bootstrap-only") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIngestCodebaseCommandSupportsDeprecatedOutputAlias(t *testing.T) {
	restoreIngestGlobals(t)

	var gotGenerate models.GenerateOptions
	runIngestTopicInfo = func(vaultPath, slug string) (models.TopicInfo, error) {
		return models.TopicInfo{
			Slug:     slug,
			Title:    "Systems Design",
			Domain:   "systems",
			RootPath: filepath.Join(vaultPath, slug),
		}, nil
	}
	runGenerate = func(ctx context.Context, opts models.GenerateOptions, observer kgenerate.Observer) (models.GenerationSummary, error) {
		gotGenerate = opts
		return models.GenerationSummary{
			Command:   "generate",
			TopicSlug: opts.TopicSlug,
			VaultPath: opts.VaultPath,
			TopicPath: filepath.Join(opts.VaultPath, opts.TopicSlug),
		}, nil
	}

	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{
		"ingest", "codebase", "/tmp/repo",
		"--topic", "systems-design",
		"--output", "/tmp/legacy-vault",
		"--progress", "never",
	})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	if gotGenerate.VaultPath != "/tmp/legacy-vault" {
		t.Fatalf("VaultPath = %q, want /tmp/legacy-vault", gotGenerate.VaultPath)
	}
}

func TestIngestBookmarksCommandRoutesToOrchestrator(t *testing.T) {
	restoreIngestGlobals(t)

	sourcePath := filepath.Join(t.TempDir(), "bookmarks.md")
	if err := os.WriteFile(sourcePath, []byte("# April Links\n"), 0o644); err != nil {
		t.Fatalf("write source file: %v", err)
	}

	var gotIngest kingest.Options
	runIngestTopicInfo = func(vaultPath, slug string) (models.TopicInfo, error) {
		return models.TopicInfo{Slug: slug, Title: "Systems Design", Domain: "systems"}, nil
	}
	runIngest = func(ctx context.Context, options kingest.Options) (models.IngestResult, error) {
		gotIngest = options
		return models.IngestResult{
			Topic:      options.Topic,
			SourceType: options.SourceKind,
			FilePath:   "systems-design/raw/bookmarks/april-links.md",
			Title:      "April Links",
		}, nil
	}

	command := newRootCommand()
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"ingest", "bookmarks", sourcePath, "--topic", "systems-design", "--vault", "/tmp/vault"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	if gotIngest.SourceKind != models.SourceKindBookmarkCluster {
		t.Fatalf("source kind = %q, want %q", gotIngest.SourceKind, models.SourceKindBookmarkCluster)
	}
	if gotIngest.SourcePath != sourcePath {
		t.Fatalf("source path = %q, want %q", gotIngest.SourcePath, sourcePath)
	}
	if gotIngest.Registry == nil {
		t.Fatal("expected converter registry to be provided")
	}

	var result models.IngestResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("stdout did not contain JSON: %v\n%s", err, stdout.String())
	}
	if result.SourceType != models.SourceKindBookmarkCluster {
		t.Fatalf("unexpected result payload: %#v", result)
	}
}

func restoreIngestGlobals(t *testing.T) {
	t.Helper()

	originalRunIngest := runIngest
	originalRunIngestTopicInfo := runIngestTopicInfo
	originalIngestGetwd := ingestGetwd
	originalLoadIngestConfig := loadIngestConfig
	originalFirecrawlScraper := newFirecrawlScraper
	originalYouTubeExtractor := newYouTubeTranscriptExtractor
	originalRegistry := newIngestRegistry
	originalRunGenerate := runGenerate

	t.Cleanup(func() {
		runIngest = originalRunIngest
		runIngestTopicInfo = originalRunIngestTopicInfo
		ingestGetwd = originalIngestGetwd
		loadIngestConfig = originalLoadIngestConfig
		newFirecrawlScraper = originalFirecrawlScraper
		newYouTubeTranscriptExtractor = originalYouTubeExtractor
		newIngestRegistry = originalRegistry
		runGenerate = originalRunGenerate
	})
}

type fakeFirecrawlScraper struct {
	scrape func(ctx context.Context, sourceURL string) (*firecrawl.ScrapeResult, error)
}

func (scraper fakeFirecrawlScraper) Scrape(ctx context.Context, sourceURL string) (*firecrawl.ScrapeResult, error) {
	if scraper.scrape == nil {
		return nil, errors.New("unexpected scrape call")
	}
	return scraper.scrape(ctx, sourceURL)
}

type fakeYouTubeExtractor struct {
	extract func(ctx context.Context, rawURL string, options youtube.ExtractOptions) (*youtube.Result, error)
}

func (extractor fakeYouTubeExtractor) Extract(ctx context.Context, rawURL string, options youtube.ExtractOptions) (*youtube.Result, error) {
	if extractor.extract == nil {
		return nil, errors.New("unexpected extract call")
	}
	return extractor.extract(ctx, rawURL, options)
}
