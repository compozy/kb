package youtube

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/compozy/kb/internal/mediadl"
)

func TestNormalizeChannelURL(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name            string
		input           string
		want            string
		wantErrContains string
	}{
		{name: "handle gets videos tab", input: "https://www.youtube.com/@aiDotEngineer", want: "https://www.youtube.com/@aiDotEngineer/videos"},
		{name: "videos tab unchanged", input: "https://www.youtube.com/@chan/videos", want: "https://www.youtube.com/@chan/videos"},
		{name: "trailing slash trimmed", input: "https://www.youtube.com/@chan/videos/", want: "https://www.youtube.com/@chan/videos"},
		{name: "shorts rewritten to videos", input: "https://www.youtube.com/@chan/shorts", want: "https://www.youtube.com/@chan/videos"},
		{name: "streams rewritten to videos", input: "https://www.youtube.com/@chan/streams", want: "https://www.youtube.com/@chan/videos"},
		{name: "scheme inferred", input: "www.youtube.com/@chan", want: "https://www.youtube.com/@chan/videos"},
		{name: "channel id path", input: "https://www.youtube.com/channel/UCabcdefghijklmnopqrstuv", want: "https://www.youtube.com/channel/UCabcdefghijklmnopqrstuv/videos"},
		{name: "playlist preserved", input: "https://www.youtube.com/playlist?list=PLabc123", want: "https://www.youtube.com/playlist?list=PLabc123"},
		{
			name:            "watch url rejected",
			input:           "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			wantErrContains: "expected a channel or playlist URL, not a video URL",
		},
		{
			name:            "watch url with list rejected",
			input:           "https://www.youtube.com/watch?v=dQw4w9WgXcQ&list=PLabc",
			wantErrContains: "expected a channel or playlist URL, not a video URL",
		},
		{
			name:            "youtu.be rejected",
			input:           "https://youtu.be/dQw4w9WgXcQ",
			wantErrContains: "expected a channel or playlist URL, not a video URL",
		},
		{
			name:            "non youtube host rejected",
			input:           "https://example.com/@chan",
			wantErrContains: "expected a youtube.com channel or playlist URL",
		},
		{
			name:            "youtube lookalike host rejected",
			input:           "https://notyoutube.com/@chan",
			wantErrContains: "expected a youtube.com channel or playlist URL",
		},
		{
			name:            "youtube suffix lookalike host rejected",
			input:           "https://youtube.com.evil.test/@chan",
			wantErrContains: "expected a youtube.com channel or playlist URL",
		},
		{name: "empty rejected", input: "   ", wantErrContains: "channel url is required"},
	}

	for _, testCase := range cases {
		t.Run("Should normalize "+testCase.name, func(t *testing.T) {
			t.Parallel()
			got, err := NormalizeChannelURL(testCase.input)
			if testCase.wantErrContains != "" {
				if err == nil {
					t.Fatalf("expected error for %q, got %q", testCase.input, got)
				}
				var youtubeErr *Error
				if !errors.As(err, &youtubeErr) {
					t.Fatalf("expected structured error for %q, got %T", testCase.input, err)
				}
				if youtubeErr.Kind != ErrorKindInvalidURL {
					t.Fatalf("kind = %q, want %q for %q", youtubeErr.Kind, ErrorKindInvalidURL, testCase.input)
				}
				if !strings.Contains(err.Error(), testCase.wantErrContains) {
					t.Fatalf("error = %q, want it to contain %q", err.Error(), testCase.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", testCase.input, err)
			}
			if got != testCase.want {
				t.Fatalf("NormalizeChannelURL(%q) = %q, want %q", testCase.input, got, testCase.want)
			}
		})
	}
}

func TestChannelVideosFromEntries(t *testing.T) {
	t.Parallel()

	t.Run("Should filter and synthesize URLs", func(t *testing.T) {
		t.Parallel()

		entries := []mediadl.PlaylistEntry{
			{ID: "aaaaaaaaaaa", Title: "First", URL: "https://www.youtube.com/watch?v=aaaaaaaaaaa"},
			{ID: "bbbbbbbbbbb", Title: "Second"},
			{ID: "aaaaaaaaaaa", Title: "Duplicate"},
			{ID: "bad", Title: "Bad ID"},
		}
		videos := channelVideosFromEntries(entries)
		if len(videos) != 2 {
			t.Fatalf("want 2 videos, got %d: %+v", len(videos), videos)
		}
		if videos[0].VideoID != "aaaaaaaaaaa" || videos[0].Title != "First" {
			t.Fatalf("unexpected first video: %+v", videos[0])
		}
		if videos[0].URL != "https://www.youtube.com/watch?v=aaaaaaaaaaa" {
			t.Fatalf("first video url = %q", videos[0].URL)
		}
		if videos[1].URL != "https://www.youtube.com/watch?v=bbbbbbbbbbb" {
			t.Fatalf("second video url should be synthesized from id, got %q", videos[1].URL)
		}
	})

	t.Run("Should return empty when no IDs are valid", func(t *testing.T) {
		t.Parallel()

		if got := channelVideosFromEntries([]mediadl.PlaylistEntry{{ID: "bad"}, {ID: ""}}); len(got) != 0 {
			t.Fatalf("expected no videos for invalid IDs, got %+v", got)
		}
	})
}

func TestChannelListingFromPlaylist(t *testing.T) {
	t.Parallel()

	listing := channelListingFromPlaylist(mediadl.PlaylistListing{
		Title:    "Asimov Academy Videos",
		Channel:  "Asimov Academy",
		Uploader: "Asimov",
		Entries: []mediadl.PlaylistEntry{
			{ID: "aaaaaaaaaaa", Title: "First", URL: "https://www.youtube.com/watch?v=aaaaaaaaaaa"},
			{ID: "bbbbbbbbbbb", Title: "Second"},
		},
	})

	if listing.Title != "Asimov Academy Videos" || listing.Channel != "Asimov Academy" || listing.Uploader != "Asimov" {
		t.Fatalf("listing metadata = %+v, want playlist metadata preserved", listing)
	}
	if len(listing.Videos) != 2 {
		t.Fatalf("videos length = %d, want 2", len(listing.Videos))
	}
	if listing.Videos[1].URL != "https://www.youtube.com/watch?v=bbbbbbbbbbb" {
		t.Fatalf("second video url = %q, want synthesized watch URL", listing.Videos[1].URL)
	}
}

func bulkVideos(ids ...string) []ChannelVideo {
	videos := make([]ChannelVideo, 0, len(ids))
	for _, id := range ids {
		videos = append(videos, ChannelVideo{VideoID: id, Title: id, URL: "https://www.youtube.com/watch?v=" + id})
	}
	return videos
}

func collectOutcomes(
	t *testing.T,
	extractor *Extractor,
	videos []ChannelVideo,
	options BulkOptions,
) []VideoOutcome {
	t.Helper()
	var outcomes []VideoOutcome
	err := extractor.BulkExtract(context.Background(), videos, options, func(outcome VideoOutcome) {
		outcomes = append(outcomes, outcome)
	})
	if err != nil {
		t.Fatalf("BulkExtract returned error: %v", err)
	}
	return outcomes
}

func TestBulkExtract(t *testing.T) {
	t.Run("Should isolate failures", func(t *testing.T) {
		extractor := &Extractor{extractFn: func(_ context.Context, rawURL string, _ ExtractOptions) (*Result, error) {
			if strings.Contains(rawURL, "fail") {
				return nil, &Error{Kind: ErrorKindTranscriptUnavailable, Message: "no captions"}
			}
			return &Result{Markdown: "transcript"}, nil
		}}

		videos := []ChannelVideo{
			{VideoID: "vid00000001", URL: "https://youtu.be/ok1"},
			{VideoID: "vid00000002", URL: "https://youtu.be/fail"},
			{VideoID: "vid00000003", URL: "https://youtu.be/ok2"},
		}
		outcomes := collectOutcomes(t, extractor, videos, BulkOptions{Concurrency: 2, MaxRetries: 1})

		var success, failure int
		for _, outcome := range outcomes {
			if outcome.Err != nil {
				failure++
			} else {
				success++
			}
		}
		if success != 2 || failure != 1 {
			t.Fatalf("want 2 success / 1 failure, got %d / %d", success, failure)
		}
	})

	t.Run("Should retry after block then succeed", func(t *testing.T) {
		var mu sync.Mutex
		attempts := map[string]int{}
		extractor := &Extractor{extractFn: func(_ context.Context, rawURL string, _ ExtractOptions) (*Result, error) {
			mu.Lock()
			attempts[rawURL]++
			count := attempts[rawURL]
			mu.Unlock()
			if count == 1 {
				return nil, &Error{Kind: ErrorKindRateLimited, Message: "HTTP Error 429"}
			}
			return &Result{Markdown: "transcript"}, nil
		}}

		videos := bulkVideos("vid00000001")
		outcomes := collectOutcomes(t, extractor, videos, BulkOptions{
			Concurrency: 1,
			MaxRetries:  3,
			BackoffMax:  10 * time.Millisecond,
		})

		if len(outcomes) != 1 || outcomes[0].Result == nil {
			t.Fatalf("expected one successful outcome after retry, got %+v", outcomes)
		}
		if attempts["https://www.youtube.com/watch?v=vid00000001"] != 2 {
			t.Fatalf("expected 2 attempts (block then success), got %d", attempts["https://www.youtube.com/watch?v=vid00000001"])
		}
	})

	t.Run("Should surface persistent block", func(t *testing.T) {
		var calls atomic.Int32
		extractor := &Extractor{extractFn: func(_ context.Context, _ string, _ ExtractOptions) (*Result, error) {
			calls.Add(1)
			return nil, &Error{Kind: ErrorKindNetworkBlocked, Message: "HTTP Error 403"}
		}}

		outcomes := collectOutcomes(t, extractor, bulkVideos("vid00000001"), BulkOptions{
			Concurrency: 1,
			MaxRetries:  2,
			BackoffMax:  5 * time.Millisecond,
		})

		if len(outcomes) != 1 || outcomes[0].Err == nil {
			t.Fatalf("expected one failed outcome, got %+v", outcomes)
		}
		if !isNetworkBlocked(outcomes[0].Err) {
			t.Fatalf("expected a network-block error, got %v", outcomes[0].Err)
		}
		if got := calls.Load(); got != 2 {
			t.Fatalf("expected 2 attempts before giving up, got %d", got)
		}
	})

	t.Run("Should respect concurrency limit", func(t *testing.T) {
		var inFlight, maxInFlight atomic.Int32
		extractor := &Extractor{extractFn: func(_ context.Context, _ string, _ ExtractOptions) (*Result, error) {
			current := inFlight.Add(1)
			for {
				observed := maxInFlight.Load()
				if current <= observed || maxInFlight.CompareAndSwap(observed, current) {
					break
				}
			}
			time.Sleep(15 * time.Millisecond)
			inFlight.Add(-1)
			return &Result{Markdown: "transcript"}, nil
		}}

		videos := bulkVideos("vid00000001", "vid00000002", "vid00000003", "vid00000004", "vid00000005", "vid00000006")
		outcomes := collectOutcomes(t, extractor, videos, BulkOptions{Concurrency: 3, MaxRetries: 1})

		if len(outcomes) != len(videos) {
			t.Fatalf("want %d outcomes, got %d", len(videos), len(outcomes))
		}
		if got := maxInFlight.Load(); got > 3 {
			t.Fatalf("concurrency exceeded limit: observed %d in flight, limit 3", got)
		}
	})

	t.Run("Should apply throttle", func(t *testing.T) {
		extractor := &Extractor{extractFn: func(_ context.Context, _ string, _ ExtractOptions) (*Result, error) {
			return &Result{Markdown: "transcript"}, nil
		}}

		videos := bulkVideos("vid00000001", "vid00000002", "vid00000003", "vid00000004")
		start := time.Now()
		collectOutcomes(t, extractor, videos, BulkOptions{Concurrency: 1, Throttle: 15 * time.Millisecond, MaxRetries: 1})
		elapsed := time.Since(start)

		// Four serial videos each wait at least the base throttle before extraction.
		if elapsed < 45*time.Millisecond {
			t.Fatalf("expected throttled run to take >= 45ms, took %s", elapsed)
		}
	})

	t.Run("Should stop on canceled context", func(t *testing.T) {
		extractor := &Extractor{extractFn: func(_ context.Context, _ string, _ ExtractOptions) (*Result, error) {
			return &Result{Markdown: "transcript"}, nil
		}}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := extractor.BulkExtract(ctx, bulkVideos("vid00000001", "vid00000002"), BulkOptions{Concurrency: 1, MaxRetries: 1}, func(VideoOutcome) {})
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})
}
