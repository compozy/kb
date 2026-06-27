package youtube

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/compozy/kb/internal/mediadl"
)

// ChannelVideo identifies a single video resolved from a channel or playlist.
type ChannelVideo struct {
	VideoID string
	Title   string
	URL     string
}

// ChannelListing is a normalized channel or playlist listing plus its videos.
type ChannelListing struct {
	Title    string
	Channel  string
	Uploader string
	Videos   []ChannelVideo
}

// BulkOptions controls channel/playlist bulk transcript extraction. Concurrency,
// Throttle, BackoffMax, and MaxRetries form the rate-limit defense: a small worker
// pool, an inter-request throttle with jitter, and adaptive exponential backoff
// when YouTube starts blocking requests.
type BulkOptions struct {
	TranscriptionPolicy     TranscriptionPolicy
	PreferredLanguages      []string
	AllowTranslatedCaptions bool
	Concurrency             int
	Throttle                time.Duration
	BackoffMax              time.Duration
	MaxRetries              int
}

// VideoOutcome reports the result of extracting a single video during a bulk run.
// Exactly one of Result or Err is set.
type VideoOutcome struct {
	Video  ChannelVideo
	Result *Result
	Err    error
}

// NormalizeChannelURL canonicalizes a channel or playlist URL to the form yt-dlp
// enumerates cleanly. Channel handles and the /shorts and /streams tabs are
// rewritten to /videos; playlist URLs keep their list query. Single-video URLs
// are rejected so callers route them to single-video ingestion instead.
func NormalizeChannelURL(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", &Error{Kind: ErrorKindInvalidURL, URL: raw, Message: "channel url is required"}
	}
	if !strings.Contains(value, "://") {
		value = "https://" + value
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Host == "" {
		return "", &Error{Kind: ErrorKindInvalidURL, URL: raw, Message: "invalid channel URL", Err: err}
	}

	host := strings.ToLower(parsed.Hostname())
	if host == "youtu.be" || host == "www.youtu.be" {
		return "", videoURLNotAllowed(raw)
	}
	if !isYouTubeHost(host) {
		return "", &Error{
			Kind:    ErrorKindInvalidURL,
			URL:     raw,
			Message: fmt.Sprintf("expected a youtube.com channel or playlist URL: %q", raw),
		}
	}

	path := strings.TrimRight(parsed.Path, "/")
	if path == "/watch" || strings.TrimSpace(parsed.Query().Get("v")) != "" {
		return "", videoURLNotAllowed(raw)
	}

	if listID := strings.TrimSpace(parsed.Query().Get("list")); listID != "" {
		normalized := url.URL{Scheme: schemeOrHTTPS(parsed.Scheme), Host: parsed.Host, Path: "/playlist", RawQuery: "list=" + listID}
		return normalized.String(), nil
	}

	if path == "" {
		return "", &Error{Kind: ErrorKindInvalidURL, URL: raw, Message: "channel URL path is required"}
	}

	var normalizedPath string
	switch {
	case strings.HasSuffix(path, "/videos"):
		normalizedPath = path
	case strings.HasSuffix(path, "/shorts"), strings.HasSuffix(path, "/streams"):
		normalizedPath = path[:strings.LastIndex(path, "/")] + "/videos"
	default:
		normalizedPath = path + "/videos"
	}

	normalized := url.URL{Scheme: schemeOrHTTPS(parsed.Scheme), Host: parsed.Host, Path: normalizedPath}
	return normalized.String(), nil
}

func isYouTubeHost(host string) bool {
	return host == "youtube.com" || strings.HasSuffix(host, ".youtube.com")
}

func videoURLNotAllowed(raw string) error {
	return &Error{
		Kind:    ErrorKindInvalidURL,
		URL:     raw,
		Message: "expected a channel or playlist URL, not a video URL; use kb ingest youtube for single videos",
	}
}

func schemeOrHTTPS(scheme string) string {
	if strings.TrimSpace(scheme) == "" {
		return "https"
	}
	return scheme
}

// ListChannel resolves the videos and top-level metadata for an already-normalized
// channel or playlist URL. A limit greater than zero caps the result to the newest uploads.
func (extractor *Extractor) ListChannel(ctx context.Context, normalizedURL string, limit int) (ChannelListing, error) {
	if extractor == nil || extractor.core == nil {
		return ChannelListing{}, errors.New("youtube channel: extractor is not configured")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	playlist, err := extractor.core.ListPlaylist(ctx, normalizedURL, limit)
	if err != nil {
		return ChannelListing{}, err
	}
	listing := channelListingFromPlaylist(playlist)
	if len(listing.Videos) == 0 {
		return ChannelListing{}, fmt.Errorf("yt-dlp channel: no videos resolved for %q", normalizedURL)
	}
	return listing, nil
}

func channelListingFromPlaylist(playlist mediadl.PlaylistListing) ChannelListing {
	return ChannelListing{
		Title:    strings.TrimSpace(playlist.Title),
		Channel:  strings.TrimSpace(playlist.Channel),
		Uploader: strings.TrimSpace(playlist.Uploader),
		Videos:   channelVideosFromEntries(playlist.Entries),
	}
}

// channelVideosFromEntries filters raw playlist entries down to valid YouTube
// videos, deduplicates by ID, and synthesizes watch URLs when needed.
func channelVideosFromEntries(entries []mediadl.PlaylistEntry) []ChannelVideo {
	videos := make([]ChannelVideo, 0, len(entries))
	seen := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		id := strings.TrimSpace(entry.ID)
		if !youtubeVideoIDPattern.MatchString(id) {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}

		videoURL := strings.TrimSpace(entry.WebpageURL)
		if videoURL == "" {
			videoURL = strings.TrimSpace(entry.URL)
		}
		if !strings.HasPrefix(strings.ToLower(videoURL), "http") {
			videoURL = "https://www.youtube.com/watch?v=" + id
		}

		videos = append(videos, ChannelVideo{
			VideoID: id,
			Title:   strings.TrimSpace(entry.Title),
			URL:     videoURL,
		})
	}
	return videos
}

// BulkExtract extracts transcripts for the supplied videos using a bounded worker
// pool with an inter-request throttle and adaptive backoff on YouTube rate limits.
// sink is invoked exactly once per video (serialized, so it need not be
// thread-safe) so callers can persist incrementally. It returns the context error
// if the run was canceled.
func (extractor *Extractor) BulkExtract(
	ctx context.Context,
	videos []ChannelVideo,
	options BulkOptions,
	sink func(VideoOutcome),
) error {
	if extractor == nil {
		return errors.New("youtube bulk: extractor is nil")
	}
	if sink == nil {
		return errors.New("youtube bulk: sink is required")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if len(videos) == 0 {
		return nil
	}

	concurrency := max(options.Concurrency, 1)
	maxRetries := max(options.MaxRetries, 1)
	backoffMax := options.BackoffMax
	if backoffMax <= 0 {
		backoffMax = 60 * time.Second
	}
	throttle := max(options.Throttle, 0)

	extractOptions := ExtractOptions{
		TranscriptionPolicy:     options.TranscriptionPolicy,
		PreferredLanguages:      options.PreferredLanguages,
		AllowTranslatedCaptions: options.AllowTranslatedCaptions,
	}
	pacer := &bulkPacer{throttle: throttle, backoffMax: backoffMax}

	var sinkMu sync.Mutex
	emit := func(outcome VideoOutcome) {
		sinkMu.Lock()
		defer sinkMu.Unlock()
		sink(outcome)
	}

	jobs := make(chan ChannelVideo)
	var wg sync.WaitGroup
	for range concurrency {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for video := range jobs {
				if ctx.Err() != nil {
					emit(VideoOutcome{Video: video, Err: ctx.Err()})
					continue
				}
				emit(extractor.extractOneVideo(ctx, video, extractOptions, pacer, maxRetries))
			}
		}()
	}

feed:
	for _, video := range videos {
		select {
		case jobs <- video:
		case <-ctx.Done():
			break feed
		}
	}
	close(jobs)
	wg.Wait()

	return ctx.Err()
}

func (extractor *Extractor) extractOneVideo(
	ctx context.Context,
	video ChannelVideo,
	options ExtractOptions,
	pacer *bulkPacer,
	maxRetries int,
) VideoOutcome {
	extract := extractor.extractFn
	if extract == nil {
		extract = extractor.Extract
	}

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if err := pacer.wait(ctx); err != nil {
			return VideoOutcome{Video: video, Err: err}
		}

		result, err := extract(ctx, video.URL, options)
		if err == nil {
			pacer.onSuccess()
			return VideoOutcome{Video: video, Result: result}
		}

		lastErr = err
		if mediadl.IsContextError(err) {
			return VideoOutcome{Video: video, Err: err}
		}
		if isNetworkBlocked(err) && attempt < maxRetries {
			pacer.onBlocked()
			continue
		}
		break
	}
	return VideoOutcome{Video: video, Err: lastErr}
}

// bulkPacer spreads requests across workers with a base throttle plus jitter and
// raises a shared exponential penalty when YouTube starts blocking requests.
type bulkPacer struct {
	throttle   time.Duration
	backoffMax time.Duration

	mu      sync.Mutex
	penalty time.Duration
}

func (pacer *bulkPacer) wait(ctx context.Context) error {
	pacer.mu.Lock()
	delay := pacer.penalty
	pacer.mu.Unlock()

	if pacer.throttle > 0 {
		jitter := time.Duration(rand.Int63n(int64(pacer.throttle/2) + 1))
		delay += pacer.throttle + jitter
	}
	return bulkSleep(ctx, delay)
}

func (pacer *bulkPacer) onBlocked() {
	pacer.mu.Lock()
	defer pacer.mu.Unlock()
	if pacer.penalty <= 0 {
		base := pacer.throttle
		if base <= 0 {
			base = time.Second
		}
		pacer.penalty = base
	} else {
		pacer.penalty *= 2
	}
	if pacer.penalty > pacer.backoffMax {
		pacer.penalty = pacer.backoffMax
	}
}

func (pacer *bulkPacer) onSuccess() {
	pacer.mu.Lock()
	defer pacer.mu.Unlock()
	pacer.penalty /= 2
	if pacer.penalty < 0 {
		pacer.penalty = 0
	}
}

func bulkSleep(ctx context.Context, delay time.Duration) error {
	if delay <= 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func isNetworkBlocked(err error) bool {
	return mediadl.IsNetworkBlocked(err) || mediadl.IsRateLimited(err)
}
