package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	kconfig "github.com/compozy/kb/internal/config"
	kingest "github.com/compozy/kb/internal/ingest"
	"github.com/compozy/kb/internal/mediadl"
	"github.com/compozy/kb/internal/models"
)

type fakeInstagramExtractor struct {
	extract func(ctx context.Context, rawURL string, options mediadl.ExtractOptions) (*mediadl.Result, error)
}

func (extractor fakeInstagramExtractor) Extract(ctx context.Context, rawURL string, options mediadl.ExtractOptions) (*mediadl.Result, error) {
	if extractor.extract == nil {
		return nil, errors.New("unexpected extract call")
	}
	return extractor.extract(ctx, rawURL, options)
}

func TestIngestInstagramCommandComposesCaptionAndIngests(t *testing.T) {
	t.Run("Should compose caption and ingest", func(t *testing.T) {
		restoreIngestGlobals(t)

		var gotExtractURL string
		var gotExtractOptions mediadl.ExtractOptions
		var gotIngest kingest.Options
		likeCount := int64(4200)
		commentCount := int64(88)
		viewCount := int64(150000)
		wantInstagramConfig := kconfig.InstagramConfig{
			YTDLPPath:     "custom-yt-dlp",
			CookiesFile:   "/tmp/instagram-cookies.txt",
			Transcription: "auto",
			RetryAttempts: 3,
			RetryBackoff:  "1s",
		}

		loadIngestConfig = func() (kconfig.Config, error) {
			return kconfig.Config{
				OpenRouter: kconfig.OpenRouterConfig{APIKey: "openrouter-key"},
				Instagram:  wantInstagramConfig,
			}, nil
		}
		newInstagramTranscriptExtractor = func(cfg kconfig.Config) instagramTranscriptExtractor {
			if !reflect.DeepEqual(cfg.Instagram, wantInstagramConfig) {
				t.Fatalf("instagram config = %#v, want %#v", cfg.Instagram, wantInstagramConfig)
			}
			return fakeInstagramExtractor{
				extract: func(_ context.Context, rawURL string, options mediadl.ExtractOptions) (*mediadl.Result, error) {
					gotExtractURL = rawURL
					gotExtractOptions = options
					return &mediadl.Result{
						Metadata: mediadl.Metadata{
							VideoID:      "C8Qh1Z6Iq3K",
							URL:          "https://www.instagram.com/reel/C8Qh1Z6Iq3K/",
							Title:        "A great reel",
							Description:  "A great reel caption",
							Channel:      "nasa",
							UploaderID:   "nasa",
							LikeCount:    &likeCount,
							CommentCount: &commentCount,
							ViewCount:    &viewCount,
							Duration:     32 * time.Second,
							Language:     "en",
						},
						Markdown:            "## Caption\nA great reel caption\n\n## Transcript\n## 00:00\nspoken words",
						Source:              mediadl.TranscriptSourceSTT,
						TranscriptionPolicy: mediadl.TranscriptionPolicyAuto,
						STTProvider:         "openai",
						STTModel:            "gpt-4o-transcribe",
					}, nil
				},
			}
		}
		runIngestTopicInfo = func(_, slug string) (models.TopicInfo, error) {
			return models.TopicInfo{Slug: slug, Title: "Space", Domain: "space"}, nil
		}
		runIngest = func(_ context.Context, options kingest.Options) (models.IngestResult, error) {
			gotIngest = options
			return models.IngestResult{
				Topic:      options.Topic,
				SourceType: options.SourceKind,
				FilePath:   "space/raw/instagram/a-great-reel.md",
				Title:      "A great reel",
			}, nil
		}

		vaultPath := t.TempDir()
		command := newRootCommand()
		var stdout bytes.Buffer
		command.SetOut(&stdout)
		command.SetErr(new(bytes.Buffer))
		command.SetArgs([]string{
			"ingest", "instagram", "https://www.instagram.com/reel/C8Qh1Z6Iq3K/",
			"--topic", "space",
			"--vault", vaultPath,
		})

		if err := command.ExecuteContext(context.Background()); err != nil {
			t.Fatalf("ExecuteContext returned error: %v", err)
		}

		if gotExtractURL != "https://www.instagram.com/reel/C8Qh1Z6Iq3K/" {
			t.Fatalf("extract URL = %q", gotExtractURL)
		}
		if gotExtractOptions.TranscriptionPolicy != mediadl.TranscriptionPolicyAuto {
			t.Fatalf("transcription policy = %q, want auto (config default)", gotExtractOptions.TranscriptionPolicy)
		}
		if gotIngest.SourceKind != models.SourceKindInstagramVideo {
			t.Fatalf("ingest source kind = %q, want %q", gotIngest.SourceKind, models.SourceKindInstagramVideo)
		}
		if gotIngest.SourceURL != "https://www.instagram.com/reel/C8Qh1Z6Iq3K/" {
			t.Fatalf("ingest source URL = %q", gotIngest.SourceURL)
		}
		if gotIngest.Markdown != "## Caption\nA great reel caption\n\n## Transcript\n## 00:00\nspoken words" {
			t.Fatalf("ingest markdown = %q", gotIngest.Markdown)
		}
		if got := gotIngest.ExtraFrontmatter["shortcode"]; got != "C8Qh1Z6Iq3K" {
			t.Fatalf("shortcode = %#v", got)
		}
		if got := gotIngest.ExtraFrontmatter["like_count"]; got != likeCount {
			t.Fatalf("like_count = %#v, want %d", got, likeCount)
		}
		if got := gotIngest.ExtraFrontmatter["uploader"]; got != "nasa" {
			t.Fatalf("uploader = %#v, want nasa", got)
		}
		if got := gotIngest.ExtraFrontmatter["transcript_source"]; got != "stt" {
			t.Fatalf("transcript_source = %#v, want stt", got)
		}

		var result models.IngestResult
		if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
			t.Fatalf("stdout did not contain JSON: %v\n%s", err, stdout.String())
		}
		if result.SourceType != models.SourceKindInstagramVideo {
			t.Fatalf("unexpected result payload: %#v", result)
		}
	})
}

func TestInstagramFrontmatterIncludesMetricsAndProvenance(t *testing.T) {
	t.Parallel()

	likeCount := int64(4200)
	commentCount := int64(88)
	viewCount := int64(150000)
	values := instagramFrontmatter(&mediadl.Result{
		Metadata: mediadl.Metadata{
			VideoID:        "C8Qh1Z6Iq3K",
			Channel:        "nasa",
			UploaderID:     "nasa",
			LikeCount:      &likeCount,
			CommentCount:   &commentCount,
			ViewCount:      &viewCount,
			Duration:       32 * time.Second,
			DurationString: "0:32",
			Language:       "en",
		},
		Source:              mediadl.TranscriptSourceSTT,
		TranscriptionPolicy: mediadl.TranscriptionPolicyAuto,
		Language:            "en",
		STTProvider:         "openai",
		STTModel:            "gpt-4o-transcribe",
	})
	want := map[string]any{
		"shortcode":            "C8Qh1Z6Iq3K",
		"uploader":             "nasa",
		"uploader_id":          "nasa",
		"like_count":           likeCount,
		"comment_count":        commentCount,
		"view_count":           viewCount,
		"upload_date":          nil,
		"duration":             int64(32),
		"duration_string":      "0:32",
		"language":             "en",
		"transcript_source":    "stt",
		"transcription_policy": "auto",
		"transcript_language":  "en",
		"stt_provider":         "openai",
		"stt_model":            "gpt-4o-transcribe",
	}
	if !reflect.DeepEqual(values, want) {
		t.Fatalf("frontmatter = %#v, want %#v", values, want)
	}
	if _, exists := values["tags"]; exists {
		t.Fatalf("frontmatter should not include reserved tags: %#v", values)
	}
}

func TestIngestInstagramCommandHonorsTranscribeFlag(t *testing.T) {
	t.Run("Should honor transcribe flag", func(t *testing.T) {
		restoreIngestGlobals(t)

		var gotPolicy mediadl.TranscriptionPolicy
		loadIngestConfig = func() (kconfig.Config, error) {
			return kconfig.Config{Instagram: kconfig.InstagramConfig{Transcription: "auto"}}, nil
		}
		newInstagramTranscriptExtractor = func(kconfig.Config) instagramTranscriptExtractor {
			return fakeInstagramExtractor{
				extract: func(_ context.Context, _ string, options mediadl.ExtractOptions) (*mediadl.Result, error) {
					gotPolicy = options.TranscriptionPolicy
					return &mediadl.Result{Markdown: "## Caption\nx", Metadata: mediadl.Metadata{Title: "x"}}, nil
				},
			}
		}
		runIngestTopicInfo = func(_, slug string) (models.TopicInfo, error) {
			return models.TopicInfo{Slug: slug, Domain: "space"}, nil
		}
		runIngest = func(_ context.Context, options kingest.Options) (models.IngestResult, error) {
			return models.IngestResult{Topic: options.Topic, SourceType: options.SourceKind}, nil
		}

		vaultPath := t.TempDir()
		command := newRootCommand()
		command.SetOut(new(bytes.Buffer))
		command.SetErr(new(bytes.Buffer))
		command.SetArgs([]string{
			"ingest", "instagram", "https://www.instagram.com/reel/C8Qh1Z6Iq3K/",
			"--topic", "space",
			"--vault", vaultPath,
			"--transcribe", "stt",
		})

		if err := command.ExecuteContext(context.Background()); err != nil {
			t.Fatalf("ExecuteContext returned error: %v", err)
		}
		if gotPolicy != mediadl.TranscriptionPolicySTT {
			t.Fatalf("transcription policy = %q, want stt", gotPolicy)
		}
	})
}
