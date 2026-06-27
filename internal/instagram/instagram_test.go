package instagram

import (
	"context"
	"errors"
	"testing"

	"github.com/compozy/kb/internal/mediadl"
)

func TestParseInstagramURL(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		raw       string
		wantCanon string
		wantID    string
		wantErr   bool
	}{
		{name: "reel", raw: "https://www.instagram.com/reel/C8Qh1Z6Iq3K/", wantCanon: "https://www.instagram.com/reel/C8Qh1Z6Iq3K/", wantID: "C8Qh1Z6Iq3K"},
		{name: "reels alias", raw: "https://www.instagram.com/reels/C8Qh1Z6Iq3K/", wantCanon: "https://www.instagram.com/reel/C8Qh1Z6Iq3K/", wantID: "C8Qh1Z6Iq3K"},
		{name: "post", raw: "https://instagram.com/p/ABC-123_/", wantCanon: "https://www.instagram.com/p/ABC-123_/", wantID: "ABC-123_"},
		{name: "tv", raw: "https://www.instagram.com/tv/XYZ/", wantCanon: "https://www.instagram.com/tv/XYZ/", wantID: "XYZ"},
		{name: "account prefixed reel", raw: "https://www.instagram.com/nasa/reel/C8Qh1Z6Iq3K/", wantCanon: "https://www.instagram.com/reel/C8Qh1Z6Iq3K/", wantID: "C8Qh1Z6Iq3K"},
		{name: "query stripped", raw: "https://www.instagram.com/reel/C8Qh1Z6Iq3K/?igsh=abc123", wantCanon: "https://www.instagram.com/reel/C8Qh1Z6Iq3K/", wantID: "C8Qh1Z6Iq3K"},
		{name: "mobile host", raw: "https://m.instagram.com/reel/C8Qh1Z6Iq3K/", wantCanon: "https://www.instagram.com/reel/C8Qh1Z6Iq3K/", wantID: "C8Qh1Z6Iq3K"},
		{name: "scheme inferred", raw: "instagram.com/reel/C8Qh1Z6Iq3K", wantCanon: "https://www.instagram.com/reel/C8Qh1Z6Iq3K/", wantID: "C8Qh1Z6Iq3K"},
		{name: "empty", raw: "", wantErr: true},
		{name: "non instagram host", raw: "https://example.com/reel/abc/", wantErr: true},
		{name: "profile only", raw: "https://www.instagram.com/nasa/", wantErr: true},
		{name: "missing shortcode", raw: "https://www.instagram.com/reel/", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			parsed, err := parseInstagramURL(tc.raw)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q, got %+v", tc.raw, parsed)
				}
				var mediaErr *mediadl.Error
				if !errors.As(err, &mediaErr) || mediaErr.Kind != mediadl.ErrorKindInvalidURL {
					t.Fatalf("expected invalid_url error, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseInstagramURL(%q) returned error: %v", tc.raw, err)
			}
			if parsed.CanonicalURL != tc.wantCanon {
				t.Fatalf("canonical url = %q, want %q", parsed.CanonicalURL, tc.wantCanon)
			}
			if parsed.VideoID != tc.wantID {
				t.Fatalf("video id = %q, want %q", parsed.VideoID, tc.wantID)
			}
		})
	}
}

type stubCore struct {
	result *mediadl.Result
	err    error
	gotURL string
}

func (core *stubCore) Extract(_ context.Context, parsed mediadl.ParsedURL, _ mediadl.ExtractOptions) (*mediadl.Result, error) {
	core.gotURL = parsed.CanonicalURL
	return core.result, core.err
}

func TestExtractComposesCaptionAndTranscript(t *testing.T) {
	t.Parallel()

	core := &stubCore{result: &mediadl.Result{
		Metadata: mediadl.Metadata{Description: "  Great reel  ", Title: "Reel"},
		Markdown: "## 00:00\nspoken words",
		Source:   mediadl.TranscriptSourceSTT,
	}}
	extractor := &Extractor{core: core}

	result, err := extractor.Extract(context.Background(), "https://www.instagram.com/reel/C8Qh1Z6Iq3K/", mediadl.ExtractOptions{})
	if err != nil {
		t.Fatalf("Extract returned error: %v", err)
	}
	if core.gotURL != "https://www.instagram.com/reel/C8Qh1Z6Iq3K/" {
		t.Fatalf("delegated url = %q", core.gotURL)
	}
	want := "## Caption\nGreat reel\n\n## Transcript\n## 00:00\nspoken words"
	if result.Markdown != want {
		t.Fatalf("markdown = %q, want %q", result.Markdown, want)
	}
	if result.Source != mediadl.TranscriptSourceSTT {
		t.Fatalf("source = %q, want stt", result.Source)
	}
}

func TestExtractTranscriptOnlyWhenNoCaption(t *testing.T) {
	t.Parallel()

	core := &stubCore{result: &mediadl.Result{
		Markdown: "## 00:00\nwords",
		Source:   mediadl.TranscriptSourceSTT,
	}}
	extractor := &Extractor{core: core}

	result, err := extractor.Extract(context.Background(), "https://www.instagram.com/reel/C8Qh1Z6Iq3K/", mediadl.ExtractOptions{})
	if err != nil {
		t.Fatalf("Extract returned error: %v", err)
	}
	if result.Markdown != "## Transcript\n## 00:00\nwords" {
		t.Fatalf("markdown = %q", result.Markdown)
	}
}

func TestExtractCaptionOnlyFallbackWhenTranscriptUnavailable(t *testing.T) {
	t.Parallel()

	core := &stubCore{
		result: &mediadl.Result{Metadata: mediadl.Metadata{Description: "Music only reel"}},
		err:    &mediadl.Error{Kind: mediadl.ErrorKindTranscriptUnavailable, Message: "captions unavailable"},
	}
	extractor := &Extractor{core: core}

	result, err := extractor.Extract(context.Background(), "https://www.instagram.com/reel/C8Qh1Z6Iq3K/", mediadl.ExtractOptions{})
	if err != nil {
		t.Fatalf("expected caption-only fallback, got error: %v", err)
	}
	if result.Markdown != "## Caption\nMusic only reel" {
		t.Fatalf("markdown = %q, want caption-only body", result.Markdown)
	}
	if result.Source != transcriptSourceNone {
		t.Fatalf("source = %q, want none", result.Source)
	}
}

func TestExtractPropagatesErrorWhenNoCaptionAvailable(t *testing.T) {
	t.Parallel()

	core := &stubCore{
		result: &mediadl.Result{Metadata: mediadl.Metadata{}},
		err:    &mediadl.Error{Kind: mediadl.ErrorKindTranscriptUnavailable},
	}
	extractor := &Extractor{core: core}

	if _, err := extractor.Extract(context.Background(), "https://www.instagram.com/reel/C8Qh1Z6Iq3K/", mediadl.ExtractOptions{}); err == nil {
		t.Fatal("expected error when neither transcript nor caption is available")
	}
}

func TestExtractRejectsInvalidURL(t *testing.T) {
	t.Parallel()

	extractor := &Extractor{core: &stubCore{}}
	if _, err := extractor.Extract(context.Background(), "https://example.com/reel/abc/", mediadl.ExtractOptions{}); err == nil {
		t.Fatal("expected invalid URL to fail before delegation")
	}
}
