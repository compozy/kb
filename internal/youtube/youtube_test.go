// Package youtube extracts transcripts and metadata from YouTube videos.
package youtube

import (
	"errors"
	"testing"
)

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
