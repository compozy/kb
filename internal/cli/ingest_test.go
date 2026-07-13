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
	"github.com/compozy/kb/internal/urlfetch"
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
	command.SetArgs([]string{"ingest", "url", "https://example.com/latency-budget", "--provider", "firecrawl", "--topic", "systems-design", "--vault", "/tmp/vault"})

	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext returned error: %v", err)
	}

	if gotFirecrawlConfig.APIKey != "firecrawl-key" || gotFirecrawlConfig.APIURL != "https://firecrawl.test" {
		t.Fatalf("firecrawl config = %#v", gotFirecrawlConfig)
	}
	if gotScrapeURL != "https://example.com/latency-budget" {
		t.Fatalf("scrape URL = %q, want source URL", gotScrapeURL)
	}
	if gotTopicVault != absoluteTestPath(t, "/tmp/vault") || gotTopicSlug != "systems-design" {
		t.Fatalf("topic lookup = (%q, %q), want (/tmp/vault, systems-design)", gotTopicVault, gotTopicSlug)
	}
	if gotIngest.VaultPath != absoluteTestPath(t, "/tmp/vault") {
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

func TestIngestURLCommandUsesHTTPLocalByDefaultAndPreservesProvenance(t *testing.T) {
	restoreIngestGlobals(t)

	var gotIngest kingest.Options
	fetchedAt := time.Date(2026, 7, 12, 15, 0, 0, 0, time.UTC)
	newHTTPLocalFetcher = func() urlContentFetcher {
		return fakeURLContentFetcher{fetch: func(context.Context, string) (*urlfetch.Result, error) {
			return &urlfetch.Result{
				SourceURL:   "https://example.com/start",
				FinalURL:    "https://www.example.com/article",
				ContentType: "text/html",
				FileName:    "article.html",
				Body:        []byte("<html><title>Article</title><body>Body</body></html>"),
				FetchedAt:   fetchedAt,
				ContentHash: "sha256:abc123",
			}, nil
		}}
	}
	runIngestTopicInfo = func(vaultPath, slug string) (models.TopicInfo, error) {
		return models.TopicInfo{Slug: slug, Title: "Systems Design", Domain: "systems", RootPath: filepath.Join(vaultPath, slug)}, nil
	}
	runIngest = func(_ context.Context, options kingest.Options) (models.IngestResult, error) {
		gotIngest = options
		return models.IngestResult{Topic: options.Topic, SourceType: options.SourceKind, FilePath: "systems-design/raw/articles/article.md", Title: "Article"}, nil
	}

	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"ingest", "url", "https://example.com/start", "--topic", "systems-design", "--vault", "/tmp/vault"})
	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext: %v", err)
	}

	if gotIngest.SourceURL != "https://example.com/start" || gotIngest.ConvertFilePath != "article.html" {
		t.Fatalf("ingest URL input = %#v", gotIngest)
	}
	if string(gotIngest.SourceContent) == "" || gotIngest.ConvertOptions["content_type"] != "text/html" {
		t.Fatalf("ingest conversion input = %#v", gotIngest)
	}
	for key, want := range map[string]any{
		"final_url": "https://www.example.com/article", "fetched_at": fetchedAt.Format(time.RFC3339), "content_type": "text/html", "content_hash": "sha256:abc123",
	} {
		if got := gotIngest.ExtraFrontmatter[key]; got != want {
			t.Fatalf("frontmatter %s = %#v, want %#v", key, got, want)
		}
	}
}

func TestIngestURLCommandRenderUsesBrowserProvider(t *testing.T) {
	restoreIngestGlobals(t)

	called := false
	newBrowserFetcher = func(cfg browserConfig) urlContentFetcher {
		called = true
		return fakeURLContentFetcher{fetch: func(context.Context, string) (*urlfetch.Result, error) {
			return &urlfetch.Result{SourceURL: "https://example.com", FinalURL: "https://example.com", ContentType: "text/html", FileName: "page.html", Body: []byte("<html></html>"), FetchedAt: time.Now(), ContentHash: "sha256:rendered"}, nil
		}}
	}
	runIngestTopicInfo = func(vaultPath, slug string) (models.TopicInfo, error) {
		return models.TopicInfo{Slug: slug, RootPath: filepath.Join(vaultPath, slug)}, nil
	}
	runIngest = func(_ context.Context, options kingest.Options) (models.IngestResult, error) {
		return models.IngestResult{Topic: options.Topic, SourceType: options.SourceKind}, nil
	}

	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{"ingest", "url", "https://example.com", "--provider", "browser", "--render", "--topic", "systems-design", "--vault", "/tmp/vault"})
	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatalf("ExecuteContext: %v", err)
	}
	if !called {
		t.Fatal("expected --render to use the browser provider")
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
		YTDLPPath:               "custom-yt-dlp",
		Proxy:                   "http://proxy.internal:8080",
		CookiesFile:             "/tmp/youtube-cookies.txt",
		UserAgent:               "kb-test-agent",
		Transcription:           "captions",
		RetryAttempts:           4,
		RetryBackoff:            "250ms",
		CaptionLanguages:        []string{"orig"},
		AllowTranslatedCaptions: true,
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
		"--sub-langs", "pt, es",
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
	if !reflect.DeepEqual(gotExtractOptions.PreferredLanguages, []string{"pt", "es"}) {
		t.Fatalf("preferred languages = %#v, want [pt es]", gotExtractOptions.PreferredLanguages)
	}
	if !gotExtractOptions.AllowTranslatedCaptions {
		t.Fatal("expected allow_translated_captions to pass through")
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
			VideoID:              "abcdefghijk",
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
		"video_id":               "abcdefghijk",
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
		VaultPath:       absoluteTestPath(t, "/tmp/vault"),
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
		VaultPath: filepath.Join(absoluteTestPath(t, "/tmp/demo-repo"), ".kb", "vault"),
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
		VaultPath: filepath.Join(absoluteTestPath(t, "/tmp/chat"), ".kb", "vault"),
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

	if gotGenerate.VaultPath != absoluteTestPath(t, "/tmp/legacy-vault") {
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
	originalRunIngestTopicNew := runIngestTopicNew
	originalIngestGetwd := ingestGetwd
	originalLoadIngestConfig := loadIngestConfig
	originalFirecrawlScraper := newFirecrawlScraper
	originalHTTPLocalFetcher := newHTTPLocalFetcher
	originalBrowserFetcher := newBrowserFetcher
	originalYouTubeExtractor := newYouTubeTranscriptExtractor
	originalYouTubeChannelExtractor := newYouTubeChannelExtractor
	originalInstagramExtractor := newInstagramTranscriptExtractor
	originalExistingYouTubeVideoIDs := existingYouTubeVideoIDs
	originalRegistry := newIngestRegistry
	originalRunGenerate := runGenerate

	t.Cleanup(func() {
		runIngest = originalRunIngest
		runIngestTopicInfo = originalRunIngestTopicInfo
		runIngestTopicNew = originalRunIngestTopicNew
		ingestGetwd = originalIngestGetwd
		loadIngestConfig = originalLoadIngestConfig
		newFirecrawlScraper = originalFirecrawlScraper
		newHTTPLocalFetcher = originalHTTPLocalFetcher
		newBrowserFetcher = originalBrowserFetcher
		newYouTubeTranscriptExtractor = originalYouTubeExtractor
		newYouTubeChannelExtractor = originalYouTubeChannelExtractor
		newInstagramTranscriptExtractor = originalInstagramExtractor
		existingYouTubeVideoIDs = originalExistingYouTubeVideoIDs
		newIngestRegistry = originalRegistry
		runGenerate = originalRunGenerate
	})
}

type fakeFirecrawlScraper struct {
	scrape func(ctx context.Context, sourceURL string) (*firecrawl.ScrapeResult, error)
}

type fakeURLContentFetcher struct {
	fetch func(ctx context.Context, sourceURL string) (*urlfetch.Result, error)
}

func (fetcher fakeURLContentFetcher) Fetch(ctx context.Context, sourceURL string) (*urlfetch.Result, error) {
	return fetcher.fetch(ctx, sourceURL)
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

func TestIngestChannelCommand(t *testing.T) {
	t.Run("Should bulk ingest with resume", func(t *testing.T) {
		restoreIngestGlobals(t)

		video1 := youtube.ChannelVideo{VideoID: "vid00000001", Title: "Already Done", URL: "https://www.youtube.com/watch?v=vid00000001"}
		video2 := youtube.ChannelVideo{VideoID: "vid00000002", Title: "Fresh One", URL: "https://www.youtube.com/watch?v=vid00000002"}
		video3 := youtube.ChannelVideo{VideoID: "vid00000003", Title: "No Captions", URL: "https://www.youtube.com/watch?v=vid00000003"}

		var gotListURL string
		var gotLimit int
		var gotBulkVideos []youtube.ChannelVideo
		var gotBulkPolicy youtube.TranscriptionPolicy
		var gotBulkLanguages []string
		var gotBulkAllowTranslated bool
		var ingestCalls int

		loadIngestConfig = func() (kconfig.Config, error) {
			return kconfig.Config{YouTube: kconfig.Default().YouTube}, nil
		}
		runIngestTopicInfo = func(_, slug string) (models.TopicInfo, error) {
			return models.TopicInfo{Slug: slug, Title: "AI Engineer", Domain: "youtube-channel"}, nil
		}
		existingYouTubeVideoIDs = func(_, _ string) (map[string]struct{}, error) {
			return map[string]struct{}{"vid00000001": {}}, nil
		}
		newYouTubeChannelExtractor = func(kconfig.Config) youtubeChannelExtractor {
			return fakeYouTubeChannelExtractor{
				list: func(_ context.Context, normalizedURL string, limit int) (youtube.ChannelListing, error) {
					gotListURL = normalizedURL
					gotLimit = limit
					return youtube.ChannelListing{
						Channel: "AI Engineer",
						Videos:  []youtube.ChannelVideo{video1, video2, video3},
					}, nil
				},
				bulk: func(_ context.Context, videos []youtube.ChannelVideo, options youtube.BulkOptions, sink func(youtube.VideoOutcome)) error {
					gotBulkVideos = videos
					gotBulkPolicy = options.TranscriptionPolicy
					gotBulkLanguages = append([]string(nil), options.PreferredLanguages...)
					gotBulkAllowTranslated = options.AllowTranslatedCaptions
					for _, video := range videos {
						if video.VideoID == "vid00000003" {
							sink(youtube.VideoOutcome{Video: video, Err: errors.New("captions unavailable")})
							continue
						}
						sink(youtube.VideoOutcome{Video: video, Result: &youtube.Result{
							Metadata: youtube.Metadata{VideoID: video.VideoID, URL: video.URL, Title: video.Title},
							Markdown: "transcript",
						}})
					}
					return nil
				},
			}
		}
		runIngest = func(_ context.Context, options kingest.Options) (models.IngestResult, error) {
			ingestCalls++
			return models.IngestResult{
				Topic:      options.Topic,
				SourceType: options.SourceKind,
				FilePath:   "yt-channels/ai/raw/youtube/" + options.Title + ".md",
				Title:      options.Title,
			}, nil
		}

		vaultPath := t.TempDir()
		command := newRootCommand()
		var stdout bytes.Buffer
		command.SetOut(&stdout)
		command.SetErr(new(bytes.Buffer))
		command.SetArgs([]string{
			"ingest", "channel", "https://www.youtube.com/@aiDotEngineer",
			"--topic", "yt-channels/ai",
			"--vault", vaultPath,
			"--limit", "10",
			"--sub-langs", "pt",
		})

		if err := command.ExecuteContext(context.Background()); err != nil {
			t.Fatalf("ExecuteContext returned error: %v", err)
		}

		if gotListURL != "https://www.youtube.com/@aiDotEngineer/videos" {
			t.Fatalf("list URL = %q, want normalized /videos URL", gotListURL)
		}
		if gotLimit != 10 {
			t.Fatalf("limit = %d, want 10", gotLimit)
		}
		if len(gotBulkVideos) != 2 {
			t.Fatalf("bulk received %d videos, want 2 (resume skips the existing one)", len(gotBulkVideos))
		}
		if gotBulkPolicy != youtube.TranscriptionPolicyCaptions {
			t.Fatalf("bulk policy = %q, want captions", gotBulkPolicy)
		}
		if !reflect.DeepEqual(gotBulkLanguages, []string{"pt"}) {
			t.Fatalf("bulk preferred languages = %#v, want [pt]", gotBulkLanguages)
		}
		if gotBulkAllowTranslated {
			t.Fatal("expected channel bulk to default translated captions off")
		}
		if ingestCalls != 1 {
			t.Fatalf("runIngest calls = %d, want 1 (only the successful fresh video)", ingestCalls)
		}

		var summary channelIngestSummary
		if err := json.Unmarshal(stdout.Bytes(), &summary); err != nil {
			t.Fatalf("decode summary: %v\n%s", err, stdout.String())
		}
		if summary.Resolved != 3 {
			t.Fatalf("resolved = %d, want 3", summary.Resolved)
		}
		if len(summary.Skipped) != 1 || summary.Skipped[0].VideoID != "vid00000001" {
			t.Fatalf("skipped = %+v, want only vid00000001", summary.Skipped)
		}
		if len(summary.Ingested) != 1 || summary.Ingested[0].VideoID != "vid00000002" {
			t.Fatalf("ingested = %+v, want only vid00000002", summary.Ingested)
		}
		if summary.Ingested[0].FilePath == "" {
			t.Fatalf("ingested entry is missing a file path: %+v", summary.Ingested[0])
		}
		if len(summary.Failures) != 1 || summary.Failures[0].VideoID != "vid00000003" {
			t.Fatalf("failures = %+v, want only vid00000003", summary.Failures)
		}
	})

	t.Run("Should skip ingest during dry run", func(t *testing.T) {
		restoreIngestGlobals(t)

		loadIngestConfig = func() (kconfig.Config, error) {
			return kconfig.Config{YouTube: kconfig.Default().YouTube}, nil
		}
		runIngestTopicInfo = func(_, slug string) (models.TopicInfo, error) {
			t.Fatalf("runIngestTopicInfo must not run during a dry run; got slug %q", slug)
			return models.TopicInfo{}, nil
		}
		existingYouTubeVideoIDs = func(_, _ string) (map[string]struct{}, error) {
			t.Fatal("existingYouTubeVideoIDs must not run during a dry run")
			return nil, nil
		}
		newYouTubeChannelExtractor = func(kconfig.Config) youtubeChannelExtractor {
			return fakeYouTubeChannelExtractor{
				list: func(context.Context, string, int) (youtube.ChannelListing, error) {
					return youtube.ChannelListing{
						Channel: "Channel Name",
						Videos: []youtube.ChannelVideo{
							{VideoID: "vid00000001", Title: "One", URL: "https://www.youtube.com/watch?v=vid00000001"},
						},
					}, nil
				},
			}
		}
		runIngest = func(context.Context, kingest.Options) (models.IngestResult, error) {
			t.Fatal("runIngest must not run during a dry run")
			return models.IngestResult{}, nil
		}

		vaultPath := t.TempDir()
		command := newRootCommand()
		var stdout bytes.Buffer
		command.SetOut(&stdout)
		command.SetErr(new(bytes.Buffer))
		command.SetArgs([]string{
			"ingest", "channel", "https://www.youtube.com/@chan",
			"--topic", "yt-channels/chan",
			"--vault", vaultPath,
			"--dry-run",
		})

		if err := command.ExecuteContext(context.Background()); err != nil {
			t.Fatalf("ExecuteContext returned error: %v", err)
		}

		var summary channelIngestSummary
		if err := json.Unmarshal(stdout.Bytes(), &summary); err != nil {
			t.Fatalf("decode summary: %v\n%s", err, stdout.String())
		}
		if !summary.DryRun || summary.Resolved != 1 {
			t.Fatalf("dry-run summary = %+v, want DryRun=true resolved=1", summary)
		}
		if _, err := os.Stat(filepath.Join(vaultPath, "yt-channels", "chan")); !os.IsNotExist(err) {
			t.Fatalf("dry-run must not create topic directory, stat err = %v", err)
		}
	})

	t.Run("Should preserve existing topic metadata", func(t *testing.T) {
		restoreIngestGlobals(t)

		loadIngestConfig = func() (kconfig.Config, error) {
			return kconfig.Config{YouTube: kconfig.Default().YouTube}, nil
		}
		runIngestTopicNew = func(_, slug, _, _ string) (models.TopicInfo, error) {
			t.Fatalf("runIngestTopicNew must not run for existing topic %q", slug)
			return models.TopicInfo{}, nil
		}
		newYouTubeChannelExtractor = func(kconfig.Config) youtubeChannelExtractor {
			return fakeYouTubeChannelExtractor{
				list: func(context.Context, string, int) (youtube.ChannelListing, error) {
					return youtube.ChannelListing{
						Channel: "Updated Channel Name",
						Videos: []youtube.ChannelVideo{
							{VideoID: "vid00000006", Title: "Fresh Existing", URL: "https://www.youtube.com/watch?v=vid00000006"},
						},
					}, nil
				},
				bulk: func(_ context.Context, videos []youtube.ChannelVideo, _ youtube.BulkOptions, sink func(youtube.VideoOutcome)) error {
					for _, video := range videos {
						sink(youtube.VideoOutcome{Video: video, Result: &youtube.Result{
							Metadata: youtube.Metadata{VideoID: video.VideoID, URL: video.URL, Title: video.Title, Channel: "Updated Channel Name"},
							Markdown: "transcript",
						}})
					}
					return nil
				},
			}
		}

		vaultPath := t.TempDir()
		if _, err := ktopic.New(vaultPath, "yt-channels/existing", "Original Title", "original-domain"); err != nil {
			t.Fatalf("create existing topic: %v", err)
		}
		metadataPath := filepath.Join(vaultPath, "yt-channels", "existing", "topic.yaml")
		before, err := os.ReadFile(metadataPath)
		if err != nil {
			t.Fatalf("read topic metadata before ingest: %v", err)
		}

		command := newRootCommand()
		command.SetOut(new(bytes.Buffer))
		command.SetErr(new(bytes.Buffer))
		command.SetArgs([]string{
			"ingest", "channel", "https://www.youtube.com/@existing",
			"--topic", "yt-channels/existing",
			"--vault", vaultPath,
			"--limit", "1",
		})

		if err := command.ExecuteContext(context.Background()); err != nil {
			t.Fatalf("ExecuteContext returned error: %v", err)
		}

		after, err := os.ReadFile(metadataPath)
		if err != nil {
			t.Fatalf("read topic metadata after ingest: %v", err)
		}
		if !bytes.Equal(before, after) {
			t.Fatalf("existing topic metadata changed:\nbefore:\n%s\nafter:\n%s", string(before), string(after))
		}
	})

	t.Run("Should auto-create missing topic before real channel ingest", func(t *testing.T) {
		restoreIngestGlobals(t)

		loadIngestConfig = func() (kconfig.Config, error) {
			return kconfig.Config{YouTube: kconfig.Default().YouTube}, nil
		}
		newYouTubeChannelExtractor = func(kconfig.Config) youtubeChannelExtractor {
			return fakeYouTubeChannelExtractor{
				list: func(context.Context, string, int) (youtube.ChannelListing, error) {
					return youtube.ChannelListing{
						Channel: "Asimov Academy",
						Videos: []youtube.ChannelVideo{
							{VideoID: "vid00000004", Title: "Lesson One", URL: "https://www.youtube.com/watch?v=vid00000004"},
						},
					}, nil
				},
				bulk: func(_ context.Context, videos []youtube.ChannelVideo, _ youtube.BulkOptions, sink func(youtube.VideoOutcome)) error {
					for _, video := range videos {
						sink(youtube.VideoOutcome{Video: video, Result: &youtube.Result{
							Metadata: youtube.Metadata{VideoID: video.VideoID, URL: video.URL, Title: video.Title, Channel: "Asimov Academy"},
							Markdown: "transcript",
						}})
					}
					return nil
				},
			}
		}

		vaultPath := t.TempDir()
		command := newRootCommand()
		var stdout bytes.Buffer
		command.SetOut(&stdout)
		command.SetErr(new(bytes.Buffer))
		command.SetArgs([]string{
			"ingest", "channel", "https://www.youtube.com/@asimovacademy",
			"--topic", "yt-channels/asimov-academy",
			"--vault", vaultPath,
			"--limit", "1",
		})

		if err := command.ExecuteContext(context.Background()); err != nil {
			t.Fatalf("ExecuteContext returned error: %v", err)
		}

		topicPath := filepath.Join(vaultPath, "yt-channels", "asimov-academy")
		metadataContent, err := os.ReadFile(filepath.Join(topicPath, "topic.yaml"))
		if err != nil {
			t.Fatalf("read topic metadata: %v", err)
		}
		for _, fragment := range []string{
			"slug: asimov-academy",
			"title: Asimov Academy",
			"domain: youtube-channel",
			"category: yt-channels",
			"path: yt-channels/asimov-academy",
			"qmd_collection: asimov-academy",
		} {
			if !strings.Contains(string(metadataContent), fragment) {
				t.Fatalf("topic.yaml missing %q:\n%s", fragment, string(metadataContent))
			}
		}
		rawEntries, err := os.ReadDir(filepath.Join(topicPath, "raw", "youtube"))
		if err != nil {
			t.Fatalf("read raw/youtube: %v", err)
		}
		var rawMarkdown []byte
		for _, entry := range rawEntries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				rawMarkdown, err = os.ReadFile(filepath.Join(topicPath, "raw", "youtube", entry.Name()))
				if err != nil {
					t.Fatalf("read raw markdown: %v", err)
				}
				break
			}
		}
		if !strings.Contains(string(rawMarkdown), "video_id: vid00000004") {
			t.Fatalf("raw youtube file missing ingested video id:\n%s", string(rawMarkdown))
		}

		var summary channelIngestSummary
		if err := json.Unmarshal(stdout.Bytes(), &summary); err != nil {
			t.Fatalf("decode summary: %v\n%s", err, stdout.String())
		}
		if summary.Topic != "yt-channels/asimov-academy" || len(summary.Ingested) != 1 {
			t.Fatalf("summary = %+v, want categorized topic with one ingested video", summary)
		}
	})

	t.Run("Should keep missing topic error actionable when auto-create is disabled", func(t *testing.T) {
		restoreIngestGlobals(t)

		loadIngestConfig = func() (kconfig.Config, error) {
			return kconfig.Config{YouTube: kconfig.Default().YouTube}, nil
		}
		newYouTubeChannelExtractor = func(kconfig.Config) youtubeChannelExtractor {
			return fakeYouTubeChannelExtractor{
				list: func(context.Context, string, int) (youtube.ChannelListing, error) {
					return youtube.ChannelListing{
						Channel: "Asimov Academy",
						Videos: []youtube.ChannelVideo{
							{VideoID: "vid00000005", Title: "Lesson Two", URL: "https://www.youtube.com/watch?v=vid00000005"},
						},
					}, nil
				},
				bulk: func(context.Context, []youtube.ChannelVideo, youtube.BulkOptions, func(youtube.VideoOutcome)) error {
					t.Fatal("BulkExtract must not run when topic creation is disabled and topic is missing")
					return nil
				},
			}
		}

		vaultPath := t.TempDir()
		command := newRootCommand()
		command.SetOut(new(bytes.Buffer))
		command.SetErr(new(bytes.Buffer))
		command.SetArgs([]string{
			"ingest", "channel", "https://www.youtube.com/@asimovacademy",
			"--topic", "yt-channels/asimov-academy",
			"--vault", vaultPath,
			"--create-topic=false",
		})

		err := command.ExecuteContext(context.Background())
		if err == nil {
			t.Fatal("expected missing topic error")
		}
		want := `kb topic new yt-channels/asimov-academy "Asimov Academy" youtube-channel`
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error = %q, want actionable command %q", err.Error(), want)
		}
		if _, statErr := os.Stat(filepath.Join(vaultPath, "yt-channels", "asimov-academy")); !os.IsNotExist(statErr) {
			t.Fatalf("disabled auto-create must not create topic directory, stat err = %v", statErr)
		}
	})

	t.Run("Should reject single-video URL", func(t *testing.T) {
		restoreIngestGlobals(t)

		loadIngestConfig = func() (kconfig.Config, error) {
			return kconfig.Config{YouTube: kconfig.Default().YouTube}, nil
		}
		runIngestTopicInfo = func(_, slug string) (models.TopicInfo, error) {
			return models.TopicInfo{Slug: slug}, nil
		}

		vaultPath := t.TempDir()
		command := newRootCommand()
		command.SetOut(new(bytes.Buffer))
		command.SetErr(new(bytes.Buffer))
		command.SetArgs([]string{
			"ingest", "channel", "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			"--topic", "yt-channels/chan",
			"--vault", vaultPath,
		})

		if err := command.ExecuteContext(context.Background()); err == nil {
			t.Fatal("expected an error when a single-video URL is passed to ingest channel")
		} else if !strings.Contains(err.Error(), "expected a channel or playlist URL, not a video URL") {
			t.Fatalf("error = %q, want single-video rejection", err.Error())
		}
	})
}

func TestIngestYouTubeMissingTopicErrorIsActionable(t *testing.T) {
	restoreIngestGlobals(t)

	loadIngestConfig = func() (kconfig.Config, error) {
		return kconfig.Config{YouTube: kconfig.Default().YouTube}, nil
	}
	newYouTubeTranscriptExtractor = func(kconfig.Config) youtubeTranscriptExtractor {
		return fakeYouTubeExtractor{
			extract: func(context.Context, string, youtube.ExtractOptions) (*youtube.Result, error) {
				t.Fatal("YouTube extraction must not run when the target topic is missing")
				return nil, nil
			},
		}
	}

	command := newRootCommand()
	command.SetOut(new(bytes.Buffer))
	command.SetErr(new(bytes.Buffer))
	command.SetArgs([]string{
		"ingest", "youtube", "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		"--topic", "yt-channels/rick",
		"--vault", t.TempDir(),
	})

	err := command.ExecuteContext(context.Background())
	if err == nil {
		t.Fatal("expected missing topic error")
	}
	want := `kb topic new yt-channels/rick "<title>" youtube-channel`
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("error = %q, want actionable command %q", err.Error(), want)
	}
}

type fakeYouTubeChannelExtractor struct {
	list func(ctx context.Context, normalizedURL string, limit int) (youtube.ChannelListing, error)
	bulk func(ctx context.Context, videos []youtube.ChannelVideo, options youtube.BulkOptions, sink func(youtube.VideoOutcome)) error
}

func (extractor fakeYouTubeChannelExtractor) ListChannel(ctx context.Context, normalizedURL string, limit int) (youtube.ChannelListing, error) {
	if extractor.list == nil {
		return youtube.ChannelListing{}, errors.New("unexpected list call")
	}
	return extractor.list(ctx, normalizedURL, limit)
}

func (extractor fakeYouTubeChannelExtractor) BulkExtract(
	ctx context.Context,
	videos []youtube.ChannelVideo,
	options youtube.BulkOptions,
	sink func(youtube.VideoOutcome),
) error {
	if extractor.bulk == nil {
		return errors.New("unexpected bulk call")
	}
	return extractor.bulk(ctx, videos, options, sink)
}
