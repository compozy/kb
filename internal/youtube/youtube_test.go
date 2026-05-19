// Package youtube extracts transcripts and metadata from YouTube videos.
package youtube

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/compozy/kb/internal/config"
)

func TestExtractorRequiresYTDLPBackend(t *testing.T) {
	t.Parallel()

	_, err := (&Extractor{}).Extract(context.Background(), "https://www.youtube.com/watch?v=dQw4w9WgXcQ", ExtractOptions{})
	if err == nil {
		t.Fatal("expected Extract to fail without yt-dlp backend")
	}
	if !strings.Contains(err.Error(), "yt-dlp backend is required") {
		t.Fatalf("error = %q, want yt-dlp requirement", err.Error())
	}
}

func TestParseVideoURLHandlesCommonYouTubeFormats(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		raw  string
	}{
		{name: "watch", raw: "https://www.youtube.com/watch?v=dQw4w9WgXcQ"},
		{name: "short", raw: "https://youtu.be/dQw4w9WgXcQ"},
		{name: "shorts", raw: "https://www.youtube.com/shorts/dQw4w9WgXcQ"},
		{name: "embed", raw: "https://www.youtube.com/embed/dQw4w9WgXcQ"},
		{name: "mobile", raw: "https://m.youtube.com/watch?v=dQw4w9WgXcQ"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			parsed, err := parseVideoURL(tc.raw)
			if err != nil {
				t.Fatalf("parseVideoURL returned error: %v", err)
			}
			if parsed.VideoID != "dQw4w9WgXcQ" {
				t.Fatalf("video ID = %q, want %q", parsed.VideoID, "dQw4w9WgXcQ")
			}
			if parsed.CanonicalURL != "https://www.youtube.com/watch?v=dQw4w9WgXcQ" {
				t.Fatalf("canonical url = %q", parsed.CanonicalURL)
			}
		})
	}
}

func TestParseVideoURLRejectsInvalidURLs(t *testing.T) {
	t.Parallel()

	for _, raw := range []string{
		"",
		"https://example.com/watch?v=dQw4w9WgXcQ",
		"https://www.youtube.com/watch?v=short",
		"://bad",
	} {
		raw := raw
		t.Run(raw, func(t *testing.T) {
			t.Parallel()

			_, err := parseVideoURL(raw)
			if err == nil {
				t.Fatal("expected parseVideoURL to fail")
			}
			var youtubeErr *Error
			if !errors.As(err, &youtubeErr) {
				t.Fatalf("expected structured error, got %T", err)
			}
			if youtubeErr.Kind != ErrorKindInvalidURL {
				t.Fatalf("kind = %q, want %q", youtubeErr.Kind, ErrorKindInvalidURL)
			}
		})
	}
}

func TestFormatTranscriptMarkdownUsesTimestampHeaders(t *testing.T) {
	t.Parallel()

	transcript := []transcriptSegment{
		{StartMs: 0, Text: " first   segment "},
		{StartMs: 3_723_000, Text: "later segment"},
		{StartMs: 4_000_000, Text: "   "},
	}

	got := formatTranscriptMarkdown(transcript)
	want := strings.Join([]string{
		"## 00:00",
		"first segment",
		"",
		"## 01:02:03",
		"later segment",
	}, "\n")
	if got != want {
		t.Fatalf("markdown = %q, want %q", got, want)
	}
}

func TestNewExtractorWithConfigConstructsYTDLPAndSTT(t *testing.T) {
	t.Parallel()

	youtubeConfig := config.YouTubeConfig{
		YTDLPPath:     "/opt/bin/yt-dlp",
		Proxy:         "http://proxy.internal:8080",
		CookiesFile:   "/tmp/youtube-cookies.txt",
		UserAgent:     "kb-test-agent",
		RetryAttempts: 4,
		RetryBackoff:  "250ms",
	}
	extractor := NewExtractorWithConfig(config.STTConfig{
		Provider:    "openai",
		APIKey:      "openai-key",
		APIURL:      "https://api.openai.test",
		Model:       "gpt-4o-transcribe",
		Language:    "auto",
		AudioFormat: "mp3",
	}, config.OpenRouterConfig{}, youtubeConfig)
	if extractor == nil {
		t.Fatal("expected extractor")
	}
	if extractor.stt == nil {
		t.Fatal("expected STT client")
	}
	if extractor.ytDLP == nil {
		t.Fatal("expected yt-dlp backend")
	}
	if extractor.ytDLP.binaryPath != "/opt/bin/yt-dlp" {
		t.Fatalf("yt-dlp path = %q, want /opt/bin/yt-dlp", extractor.ytDLP.binaryPath)
	}
	if extractor.ytDLP.cfg != youtubeConfig {
		t.Fatalf("yt-dlp config = %#v, want %#v", extractor.ytDLP.cfg, youtubeConfig)
	}
	if extractor.ytDLP.retry.Attempts != 4 {
		t.Fatalf("yt-dlp retry attempts = %d, want 4", extractor.ytDLP.retry.Attempts)
	}
	if extractor.ytDLP.retry.Backoff != 250*time.Millisecond {
		t.Fatalf("yt-dlp retry backoff = %v, want 250ms", extractor.ytDLP.retry.Backoff)
	}
	if extractor.setupErr != nil {
		t.Fatalf("setupErr = %v, want nil", extractor.setupErr)
	}
}

func TestErrorFormattingAndUnwrap(t *testing.T) {
	t.Parallel()

	cause := errors.New("boom")
	err := &Error{
		Kind:    ErrorKindPrivate,
		URL:     "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		Message: "video is private",
		Err:     cause,
	}

	if !strings.Contains(err.Error(), "youtube private") {
		t.Fatalf("error = %q", err.Error())
	}
	if !errors.Is(err, cause) {
		t.Fatal("expected wrapped error")
	}
}

func TestShouldAttemptSTT(t *testing.T) {
	t.Parallel()

	if (&Extractor{}).shouldAttemptSTT(TranscriptionPolicyAuto) {
		t.Fatal("expected false when STT client is missing")
	}

	extractor := &Extractor{stt: &stubSTTClient{}}
	if !extractor.shouldAttemptSTT(TranscriptionPolicyAuto) {
		t.Fatal("expected auto policy to request STT when a transcriber is present")
	}
	if !extractor.shouldAttemptSTT(TranscriptionPolicySTT) {
		t.Fatal("expected stt policy to request STT when a transcriber is present")
	}
	if extractor.shouldAttemptSTT(TranscriptionPolicyCaptions) {
		t.Fatal("expected captions policy to avoid STT")
	}
}

func TestParseTranscriptionPolicy(t *testing.T) {
	t.Parallel()

	for _, value := range []string{"", "captions", "auto", "stt", " STT "} {
		if _, err := ParseTranscriptionPolicy(value); err != nil {
			t.Fatalf("ParseTranscriptionPolicy(%q) returned error: %v", value, err)
		}
	}
	if _, err := ParseTranscriptionPolicy("maybe"); err == nil {
		t.Fatal("expected invalid transcription policy to fail")
	}
}

func TestLanguageHelpers(t *testing.T) {
	t.Parallel()

	got := normalizeLanguages([]string{" en ", "EN", "", "pt-BR"})
	if len(got) != 2 || got[0] != "en" || got[1] != "pt-br" {
		t.Fatalf("normalized languages = %#v", got)
	}
	if !languageMatches("en-US", "en") {
		t.Fatal("expected language prefix match")
	}
	if languageMatches("de", "en") {
		t.Fatal("did not expect unrelated language match")
	}
}

func TestIsYouTubeNetworkBlocked(t *testing.T) {
	t.Parallel()

	if !isYouTubeNetworkBlocked(errors.New("unexpected status code: 403")) {
		t.Fatal("expected forbidden status to be treated as network blocked")
	}
	if !isYouTubeNetworkBlocked(errors.New("unexpected status code: 429")) {
		t.Fatal("expected text status to be treated as network blocked")
	}
	if isYouTubeNetworkBlocked(errors.New("unexpected status code: 404")) {
		t.Fatal("did not expect 404 to be treated as network blocked")
	}
}

type stubSTTClient struct {
	called      bool
	audio       []byte
	format      string
	transcript  string
	transcripts []string
	err         error
}

func (client *stubSTTClient) Provider() string {
	return "stub"
}

func (client *stubSTTClient) Model() string {
	return "stub-model"
}

func (client *stubSTTClient) Transcribe(_ context.Context, audio []byte, format string) (string, error) {
	client.called = true
	client.audio = append([]byte(nil), audio...)
	client.format = format
	if client.err != nil {
		return "", client.err
	}
	if len(client.transcripts) > 0 {
		transcript := client.transcripts[0]
		client.transcripts = client.transcripts[1:]
		return transcript, nil
	}
	return client.transcript, nil
}
