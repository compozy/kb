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
		t.Run("Should parse "+tc.name, func(t *testing.T) {
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

func TestExtract(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		rawURL           string
		coreResult       *mediadl.Result
		coreErr          error
		wantMarkdown     string
		wantSource       mediadl.TranscriptSource
		wantMetadataURL  string
		wantDelegatedURL string
		wantErrKind      mediadl.ErrorKind
	}{
		{
			name:   "Should compose caption and transcript",
			rawURL: "https://www.instagram.com/reel/C8Qh1Z6Iq3K/",
			coreResult: &mediadl.Result{
				Metadata: mediadl.Metadata{Description: "  Great reel  ", Title: "Reel"},
				Markdown: "## 00:00\nspoken words",
				Source:   mediadl.TranscriptSourceSTT,
			},
			wantMarkdown:     "## Caption\nGreat reel\n\n## Transcript\n## 00:00\nspoken words",
			wantSource:       mediadl.TranscriptSourceSTT,
			wantMetadataURL:  "https://www.instagram.com/reel/C8Qh1Z6Iq3K/",
			wantDelegatedURL: "https://www.instagram.com/reel/C8Qh1Z6Iq3K/",
		},
		{
			name:   "Should keep transcript when caption is absent",
			rawURL: "https://www.instagram.com/reel/C8Qh1Z6Iq3K/",
			coreResult: &mediadl.Result{
				Markdown: "## 00:00\nwords",
				Source:   mediadl.TranscriptSourceSTT,
			},
			wantMarkdown:     "## Transcript\n## 00:00\nwords",
			wantSource:       mediadl.TranscriptSourceSTT,
			wantMetadataURL:  "https://www.instagram.com/reel/C8Qh1Z6Iq3K/",
			wantDelegatedURL: "https://www.instagram.com/reel/C8Qh1Z6Iq3K/",
		},
		{
			name:             "Should fall back to caption only when transcript is unavailable",
			rawURL:           "instagram.com/reels/C8Qh1Z6Iq3K",
			coreResult:       &mediadl.Result{Metadata: mediadl.Metadata{Description: "Music only reel"}},
			coreErr:          &mediadl.Error{Kind: mediadl.ErrorKindTranscriptUnavailable, Message: "captions unavailable"},
			wantMarkdown:     "## Caption\nMusic only reel",
			wantSource:       transcriptSourceNone,
			wantMetadataURL:  "https://www.instagram.com/reel/C8Qh1Z6Iq3K/",
			wantDelegatedURL: "https://www.instagram.com/reel/C8Qh1Z6Iq3K/",
		},
		{
			name:             "Should propagate transcript unavailable when no caption exists",
			rawURL:           "https://www.instagram.com/reel/C8Qh1Z6Iq3K/",
			coreResult:       &mediadl.Result{Metadata: mediadl.Metadata{}},
			coreErr:          &mediadl.Error{Kind: mediadl.ErrorKindTranscriptUnavailable},
			wantDelegatedURL: "https://www.instagram.com/reel/C8Qh1Z6Iq3K/",
			wantErrKind:      mediadl.ErrorKindTranscriptUnavailable,
		},
		{
			name:        "Should reject invalid URL before delegation",
			rawURL:      "https://example.com/reel/abc/",
			coreResult:  &mediadl.Result{},
			wantErrKind: mediadl.ErrorKindInvalidURL,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			core := &stubCore{result: tc.coreResult, err: tc.coreErr}
			extractor := &Extractor{core: core}
			result, err := extractor.Extract(context.Background(), tc.rawURL, mediadl.ExtractOptions{})
			if tc.wantErrKind != "" {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				var mediaErr *mediadl.Error
				if !errors.As(err, &mediaErr) || mediaErr.Kind != tc.wantErrKind {
					t.Fatalf("error = %v, want kind %s", err, tc.wantErrKind)
				}
				if tc.wantDelegatedURL == "" && core.gotURL != "" {
					t.Fatalf("delegated URL = %q, want no delegation", core.gotURL)
				}
				return
			}
			if err != nil {
				t.Fatalf("Extract returned error: %v", err)
			}
			if core.gotURL != tc.wantDelegatedURL {
				t.Fatalf("delegated url = %q, want %q", core.gotURL, tc.wantDelegatedURL)
			}
			if result.Markdown != tc.wantMarkdown {
				t.Fatalf("markdown = %q, want %q", result.Markdown, tc.wantMarkdown)
			}
			if result.Source != tc.wantSource {
				t.Fatalf("source = %q, want %q", result.Source, tc.wantSource)
			}
			if result.Metadata.URL != tc.wantMetadataURL {
				t.Fatalf("metadata url = %q, want %q", result.Metadata.URL, tc.wantMetadataURL)
			}
		})
	}
}
