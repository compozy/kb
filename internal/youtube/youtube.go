// Package youtube extracts transcripts and metadata from YouTube videos.
package youtube

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"

	ytdl "github.com/kkdai/youtube/v2"

	"github.com/compozy/kb/internal/config"
)

const (
	transcriptUnavailableMessage = "captions unavailable"
)

var youtubeVideoIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)

// TranscriptSource identifies how the transcript was produced.
type TranscriptSource string

const (
	// TranscriptSourceCaptions means YouTube captions were fetched directly.
	TranscriptSourceCaptions TranscriptSource = "captions"
	// TranscriptSourceSTT means audio was transcribed through OpenRouter.
	TranscriptSourceSTT TranscriptSource = "stt"
)

// ErrorKind categorizes user-facing YouTube extraction failures.
type ErrorKind string

const (
	// ErrorKindInvalidURL reports a malformed or unsupported YouTube URL.
	ErrorKindInvalidURL ErrorKind = "invalid_url"
	// ErrorKindUnavailable reports a video that cannot be accessed.
	ErrorKindUnavailable ErrorKind = "unavailable"
	// ErrorKindPrivate reports a private video.
	ErrorKindPrivate ErrorKind = "private"
	// ErrorKindAgeRestricted reports a video that requires login/age confirmation.
	ErrorKindAgeRestricted ErrorKind = "age_restricted"
	// ErrorKindTranscriptUnavailable reports missing or disabled captions.
	ErrorKindTranscriptUnavailable ErrorKind = "transcript_unavailable"
	// ErrorKindAudioUnavailable reports that no supported audio stream could be downloaded.
	ErrorKindAudioUnavailable ErrorKind = "audio_unavailable"
	// ErrorKindNetworkBlocked reports a YouTube network block or bot-detection response.
	ErrorKindNetworkBlocked ErrorKind = "network_blocked"
)

// Error carries structured failure details for callers that need to branch on
// YouTube-specific failure modes.
type Error struct {
	Kind    ErrorKind
	URL     string
	VideoID string
	Message string
	Err     error
}

// Error formats a human-readable error message.
func (err *Error) Error() string {
	if err == nil {
		return ""
	}

	target := strings.TrimSpace(err.URL)
	if target == "" {
		target = strings.TrimSpace(err.VideoID)
	}
	if target == "" {
		target = "video"
	}

	message := strings.TrimSpace(err.Message)
	if message == "" {
		message = string(err.Kind)
	}

	if err.Err != nil {
		return fmt.Sprintf("youtube %s for %q: %s: %v", err.Kind, target, message, err.Err)
	}

	return fmt.Sprintf("youtube %s for %q: %s", err.Kind, target, message)
}

// Unwrap returns the underlying cause.
func (err *Error) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.Err
}

// Metadata contains normalized video metadata.
type Metadata struct {
	VideoID     string
	URL         string
	Title       string
	Channel     string
	Duration    time.Duration
	PublishDate time.Time
}

// Result contains the extracted metadata and transcript markdown.
type Result struct {
	Metadata Metadata
	Markdown string
	Source   TranscriptSource
	Language string
}

// ExtractOptions controls transcript extraction behavior.
type ExtractOptions struct {
	EnableSTTFallback  bool
	PreferredLanguages []string
}

type parsedVideoURL struct {
	CanonicalURL string
	VideoID      string
}

type youtubeClient interface {
	GetVideoContext(ctx context.Context, rawURL string) (*ytdl.Video, error)
	GetTranscriptCtx(ctx context.Context, video *ytdl.Video, lang string) (ytdl.VideoTranscript, error)
	GetStreamContext(ctx context.Context, video *ytdl.Video, format *ytdl.Format) (io.ReadCloser, int64, error)
}

type sttClient interface {
	Configured() bool
	Transcribe(ctx context.Context, audio []byte, format string) (string, error)
}

type retryPolicy struct {
	Attempts int
	Backoff  time.Duration
}

// Extractor orchestrates transcript extraction and optional STT fallback.
type Extractor struct {
	youtube  youtubeClient
	stt      sttClient
	retry    retryPolicy
	setupErr error
}

// NewExtractor constructs a default extractor backed by kkdai/youtube and the
// OpenRouter STT client.
func NewExtractor(cfg config.OpenRouterConfig, youtubeConfigs ...config.YouTubeConfig) *Extractor {
	youtubeConfig := config.Default().YouTube
	if len(youtubeConfigs) > 0 {
		youtubeConfig = youtubeConfigs[0]
	}
	httpClient, err := newHTTPClient(youtubeConfig)
	backoff, backoffErr := youtubeConfig.RetryBackoffDuration()
	if err == nil {
		err = backoffErr
	}
	attempts := youtubeConfig.RetryAttempts
	if attempts < 1 {
		attempts = 1
	}

	return &Extractor{
		youtube: &ytdl.Client{HTTPClient: httpClient},
		stt:     NewOpenRouterClient(cfg),
		retry: retryPolicy{
			Attempts: attempts,
			Backoff:  backoff,
		},
		setupErr: err,
	}
}

func newHTTPClient(cfg config.YouTubeConfig) (*http.Client, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	proxyValue := strings.TrimSpace(cfg.Proxy)
	if proxyValue != "" {
		proxyURL, err := url.Parse(proxyValue)
		if err != nil {
			return nil, fmt.Errorf("parse youtube proxy: %w", err)
		}
		switch strings.ToLower(proxyURL.Scheme) {
		case "http", "https", "socks5":
		default:
			return nil, fmt.Errorf("youtube proxy scheme must be http, https, or socks5: %q", proxyURL.Scheme)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	} else {
		transport.Proxy = http.ProxyFromEnvironment
	}

	cookies, err := loadCookiesFile(cfg.CookiesFile)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: &requestDecoratingTransport{
			base:      transport,
			cookies:   cookies,
			userAgent: strings.TrimSpace(cfg.UserAgent),
		},
	}, nil
}

type requestDecoratingTransport struct {
	base      http.RoundTripper
	cookies   []*http.Cookie
	userAgent string
}

func (transport *requestDecoratingTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if transport == nil {
		return http.DefaultTransport.RoundTrip(request)
	}
	clonedRequest := request.Clone(request.Context())
	for _, cookie := range transport.cookies {
		clonedRequest.AddCookie(cookie)
	}
	if transport.userAgent != "" {
		clonedRequest.Header.Set("User-Agent", transport.userAgent)
	}
	base := transport.base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(clonedRequest)
}

func loadCookiesFile(path string) ([]*http.Cookie, error) {
	cleanPath := strings.TrimSpace(path)
	if cleanPath == "" {
		return nil, nil
	}

	file, err := os.Open(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("open youtube cookies file %q: %w", cleanPath, err)
	}
	defer func() {
		_ = file.Close()
	}()

	cookies, err := parseCookies(file)
	if err != nil {
		return nil, fmt.Errorf("parse youtube cookies file %q: %w", cleanPath, err)
	}
	return cookies, nil
}

func parseCookies(reader io.Reader) ([]*http.Cookie, error) {
	scanner := bufio.NewScanner(reader)
	cookies := []*http.Cookie{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#HttpOnly_") {
			line = strings.TrimPrefix(line, "#HttpOnly_")
		} else if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Count(line, "	") >= 6 {
			cookie, err := parseNetscapeCookie(line)
			if err != nil {
				return nil, err
			}
			cookies = append(cookies, cookie)
			continue
		}

		for _, part := range strings.Split(line, ";") {
			name, value, ok := strings.Cut(strings.TrimSpace(part), "=")
			if !ok || strings.TrimSpace(name) == "" {
				continue
			}
			cookies = append(cookies, &http.Cookie{Name: strings.TrimSpace(name), Value: strings.TrimSpace(value)})
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cookies, nil
}

func parseNetscapeCookie(line string) (*http.Cookie, error) {
	fields := strings.Split(line, "	")
	if len(fields) < 7 {
		return nil, fmt.Errorf("invalid Netscape cookie line")
	}
	name := strings.TrimSpace(fields[5])
	if name == "" {
		return nil, fmt.Errorf("invalid Netscape cookie line: missing name")
	}

	return &http.Cookie{
		Domain: strings.TrimSpace(fields[0]),
		Path:   strings.TrimSpace(fields[2]),
		Secure: strings.EqualFold(strings.TrimSpace(fields[3]), "TRUE"),
		Name:   name,
		Value:  strings.TrimSpace(fields[6]),
	}, nil
}

// Extract fetches video metadata and transcript markdown from a YouTube URL.
func (extractor *Extractor) Extract(ctx context.Context, rawURL string, options ExtractOptions) (*Result, error) {
	if extractor == nil {
		return nil, errors.New("youtube extract: extractor is nil")
	}
	if extractor.youtube == nil {
		return nil, errors.New("youtube extract: client is nil")
	}
	if extractor.setupErr != nil {
		return nil, fmt.Errorf("youtube extract: configure network client: %w", extractor.setupErr)
	}
	if ctx == nil {
		ctx = context.Background()
	}

	parsed, err := parseVideoURL(rawURL)
	if err != nil {
		return nil, err
	}

	video, err := extractor.getVideo(ctx, parsed)
	if err != nil {
		return nil, err
	}

	result := &Result{
		Metadata: Metadata{
			VideoID:     parsed.VideoID,
			URL:         parsed.CanonicalURL,
			Title:       strings.TrimSpace(video.Title),
			Channel:     strings.TrimSpace(video.Author),
			Duration:    video.Duration,
			PublishDate: video.PublishDate.UTC(),
		},
	}

	markdown, language, err := extractor.extractTranscript(ctx, video, options.PreferredLanguages)
	if err == nil {
		result.Markdown = markdown
		result.Source = TranscriptSourceCaptions
		result.Language = language
		return result, nil
	}
	if isContextError(err) {
		return result, err
	}

	if !extractor.shouldAttemptSTT(options) || !isTranscriptUnavailable(err) {
		return result, err
	}

	audio, format, audioErr := extractor.downloadAudio(ctx, video)
	if audioErr != nil {
		return result, errors.Join(err, fmt.Errorf("youtube stt fallback: %w", audioErr))
	}

	transcript, sttErr := extractor.stt.Transcribe(ctx, audio, format)
	if sttErr != nil {
		return result, errors.Join(err, fmt.Errorf("youtube stt fallback: %w", sttErr))
	}

	result.Markdown = formatSTTMarkdown(transcript)
	result.Source = TranscriptSourceSTT

	return result, nil
}

func (extractor *Extractor) shouldAttemptSTT(options ExtractOptions) bool {
	if extractor.stt == nil {
		return false
	}
	return options.EnableSTTFallback || extractor.stt.Configured()
}

func (extractor *Extractor) getVideo(ctx context.Context, parsed parsedVideoURL) (*ytdl.Video, error) {
	var video *ytdl.Video
	err := extractor.retryOperation(ctx, func() error {
		loadedVideo, err := extractor.youtube.GetVideoContext(ctx, parsed.CanonicalURL)
		if err != nil {
			return wrapVideoError(parsed, err)
		}
		video = loadedVideo
		return nil
	})
	if err != nil {
		return nil, err
	}
	return video, nil
}

func (extractor *Extractor) getTranscript(ctx context.Context, video *ytdl.Video, language string) (ytdl.VideoTranscript, error) {
	var transcript ytdl.VideoTranscript
	err := extractor.retryOperation(ctx, func() error {
		loadedTranscript, err := extractor.youtube.GetTranscriptCtx(ctx, video, language)
		if err != nil {
			return err
		}
		transcript = loadedTranscript
		return nil
	})
	if err != nil {
		return nil, err
	}
	return transcript, nil
}

func (extractor *Extractor) getStream(ctx context.Context, video *ytdl.Video, format *ytdl.Format) (io.ReadCloser, int64, error) {
	var stream io.ReadCloser
	var size int64
	err := extractor.retryOperation(ctx, func() error {
		loadedStream, loadedSize, err := extractor.youtube.GetStreamContext(ctx, video, format)
		if err != nil {
			return err
		}
		stream = loadedStream
		size = loadedSize
		return nil
	})
	if err != nil {
		return nil, 0, err
	}
	return stream, size, nil
}

func (extractor *Extractor) retryOperation(ctx context.Context, operation func() error) error {
	attempts := extractor.retry.Attempts
	if attempts < 1 {
		attempts = 1
	}
	backoff := extractor.retry.Backoff

	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		err := operation()
		if err == nil {
			return nil
		}
		if isContextError(err) {
			return err
		}
		lastErr = err
		if attempt == attempts || !isRetryableYouTubeError(err) {
			return err
		}
		if backoff <= 0 {
			continue
		}
		timer := time.NewTimer(backoff)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
		backoff *= 2
	}
	return lastErr
}

func (extractor *Extractor) extractTranscript(ctx context.Context, video *ytdl.Video, preferredLanguages []string) (string, string, error) {
	tracks := orderedCaptionTracks(video.CaptionTracks, preferredLanguages)
	if len(tracks) == 0 {
		return "", "", &Error{
			Kind:    ErrorKindTranscriptUnavailable,
			VideoID: video.ID,
			Message: transcriptUnavailableMessage,
			Err:     ytdl.ErrTranscriptDisabled,
		}
	}

	var transcriptErrors []error

	for _, track := range tracks {
		transcript, err := extractor.getTranscript(ctx, video, track.LanguageCode)
		if err != nil {
			if isContextError(err) {
				return "", "", err
			}
			if isYouTubeNetworkBlocked(err) {
				return "", "", newNetworkBlockedError(video.ID, "captions request was blocked by YouTube; configure [youtube].proxy or [youtube].cookies_file, or run from a trusted network", err)
			}
			transcriptErrors = append(transcriptErrors, fmt.Errorf("%s: %w", track.LanguageCode, err))
			continue
		}

		markdown := formatTranscriptMarkdown(transcript)
		if markdown == "" {
			transcriptErrors = append(transcriptErrors, fmt.Errorf("%s: empty transcript", track.LanguageCode))
			continue
		}

		return markdown, track.LanguageCode, nil
	}

	return "", "", &Error{
		Kind:    ErrorKindTranscriptUnavailable,
		VideoID: video.ID,
		Message: transcriptUnavailableMessage,
		Err:     errors.Join(transcriptErrors...),
	}
}

func (extractor *Extractor) downloadAudio(ctx context.Context, video *ytdl.Video) ([]byte, string, error) {
	format, normalizedFormat, err := pickAudioFormat(video.Formats)
	if err != nil {
		return nil, "", err
	}

	stream, _, err := extractor.getStream(ctx, video, format)
	if err != nil {
		if isYouTubeNetworkBlocked(err) {
			return nil, "", &Error{
				Kind:    ErrorKindAudioUnavailable,
				VideoID: video.ID,
				Message: "audio download was blocked by YouTube; configure [youtube].proxy or [youtube].cookies_file, or run from a trusted network",
				Err:     err,
			}
		}
		return nil, "", fmt.Errorf("download audio stream: %w", err)
	}
	defer func() {
		_ = stream.Close()
	}()

	audio, err := io.ReadAll(stream)
	if err != nil {
		return nil, "", fmt.Errorf("read audio stream: %w", err)
	}
	if len(audio) == 0 {
		return nil, "", &Error{
			Kind:    ErrorKindAudioUnavailable,
			VideoID: video.ID,
			Message: "audio stream is empty",
		}
	}

	return audio, normalizedFormat, nil
}

func parseVideoURL(rawURL string) (parsedVideoURL, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return parsedVideoURL{}, &Error{
			Kind:    ErrorKindInvalidURL,
			URL:     rawURL,
			Message: "youtube url is required",
		}
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return parsedVideoURL{}, &Error{
			Kind:    ErrorKindInvalidURL,
			URL:     rawURL,
			Message: "invalid YouTube URL",
			Err:     err,
		}
	}

	var videoID string

	switch strings.ToLower(parsedURL.Host) {
	case "youtu.be":
		videoID = firstPathSegment(parsedURL.Path)
	case "www.youtube.com", "youtube.com", "m.youtube.com", "music.youtube.com":
		switch strings.Trim(parsedURL.Path, "/") {
		case "watch":
			videoID = strings.TrimSpace(parsedURL.Query().Get("v"))
		default:
			segments := pathSegments(parsedURL.Path)
			if len(segments) >= 2 && (segments[0] == "shorts" || segments[0] == "embed") {
				videoID = segments[1]
			}
		}
	default:
		return parsedVideoURL{}, &Error{
			Kind:    ErrorKindInvalidURL,
			URL:     rawURL,
			Message: "unsupported YouTube host",
		}
	}

	if !youtubeVideoIDPattern.MatchString(videoID) {
		return parsedVideoURL{}, &Error{
			Kind:    ErrorKindInvalidURL,
			URL:     rawURL,
			Message: "invalid YouTube video identifier",
		}
	}

	return parsedVideoURL{
		CanonicalURL: "https://www.youtube.com/watch?v=" + videoID,
		VideoID:      videoID,
	}, nil
}

func wrapVideoError(parsed parsedVideoURL, err error) error {
	switch {
	case errors.Is(err, ytdl.ErrVideoPrivate):
		return &Error{
			Kind:    ErrorKindPrivate,
			URL:     parsed.CanonicalURL,
			VideoID: parsed.VideoID,
			Message: "video is private",
			Err:     err,
		}
	case errors.Is(err, ytdl.ErrLoginRequired):
		return &Error{
			Kind:    ErrorKindAgeRestricted,
			URL:     parsed.CanonicalURL,
			VideoID: parsed.VideoID,
			Message: "video is age restricted",
			Err:     err,
		}
	default:
		if isYouTubeNetworkBlocked(err) {
			return &Error{
				Kind:    ErrorKindNetworkBlocked,
				URL:     parsed.CanonicalURL,
				VideoID: parsed.VideoID,
				Message: "video request was blocked by YouTube; configure [youtube].proxy or [youtube].cookies_file, or run from a trusted network",
				Err:     err,
			}
		}

		var statusErr *ytdl.ErrPlayabiltyStatus
		if errors.As(err, &statusErr) {
			return &Error{
				Kind:    ErrorKindUnavailable,
				URL:     parsed.CanonicalURL,
				VideoID: parsed.VideoID,
				Message: "video is unavailable",
				Err:     err,
			}
		}

		return fmt.Errorf("load YouTube video %q: %w", parsed.CanonicalURL, err)
	}
}

func newNetworkBlockedError(videoID string, message string, err error) error {
	return &Error{
		Kind:    ErrorKindNetworkBlocked,
		VideoID: videoID,
		Message: message,
		Err:     err,
	}
}

func isTranscriptUnavailable(err error) bool {
	var youtubeErr *Error
	return errors.As(err, &youtubeErr) && youtubeErr.Kind == ErrorKindTranscriptUnavailable
}

func isRetryableYouTubeError(err error) bool {
	if isContextError(err) {
		return false
	}
	if isYouTubeNetworkBlocked(err) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr)
}

func isYouTubeNetworkBlocked(err error) bool {
	status, ok := unexpectedStatusCode(err)
	if !ok {
		return false
	}
	return status == http.StatusBadRequest ||
		status == http.StatusForbidden ||
		status == http.StatusTooManyRequests ||
		status >= http.StatusInternalServerError
}

func unexpectedStatusCode(err error) (int, bool) {
	if err == nil {
		return 0, false
	}
	var statusErr ytdl.ErrUnexpectedStatusCode
	if errors.As(err, &statusErr) {
		return int(statusErr), true
	}

	message := err.Error()
	for _, status := range []int{
		http.StatusBadRequest,
		http.StatusForbidden,
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
	} {
		if strings.Contains(message, fmt.Sprintf("unexpected status code: %d", status)) {
			return status, true
		}
	}
	return 0, false
}

func orderedCaptionTracks(tracks []ytdl.CaptionTrack, preferredLanguages []string) []ytdl.CaptionTrack {
	ordered := append([]ytdl.CaptionTrack(nil), tracks...)
	if len(ordered) == 0 {
		return nil
	}

	preferredLanguages = normalizeLanguages(preferredLanguages)

	sort.SliceStable(ordered, func(i int, j int) bool {
		left := captionTrackPriority(ordered[i], preferredLanguages)
		right := captionTrackPriority(ordered[j], preferredLanguages)
		return left < right
	})

	return ordered
}

func captionTrackPriority(track ytdl.CaptionTrack, preferredLanguages []string) int {
	preferredRank := len(preferredLanguages) + 1
	for index, language := range preferredLanguages {
		if languageMatches(track.LanguageCode, language) {
			preferredRank = index
			break
		}
	}

	manualRank := 1
	if strings.TrimSpace(strings.ToLower(track.Kind)) != "asr" {
		manualRank = 0
	}

	return preferredRank*2 + manualRank
}

func normalizeLanguages(languages []string) []string {
	if len(languages) == 0 {
		return nil
	}

	normalized := make([]string, 0, len(languages))
	seen := make(map[string]struct{}, len(languages))

	for _, language := range languages {
		language = strings.ToLower(strings.TrimSpace(language))
		if language == "" {
			continue
		}
		if _, ok := seen[language]; ok {
			continue
		}
		seen[language] = struct{}{}
		normalized = append(normalized, language)
	}

	return normalized
}

func languageMatches(trackLanguage string, preferredLanguage string) bool {
	trackLanguage = strings.ToLower(strings.TrimSpace(trackLanguage))
	preferredLanguage = strings.ToLower(strings.TrimSpace(preferredLanguage))

	if trackLanguage == preferredLanguage {
		return true
	}
	if strings.HasPrefix(trackLanguage, preferredLanguage+"-") {
		return true
	}
	if strings.HasPrefix(preferredLanguage, trackLanguage+"-") {
		return true
	}

	return false
}

func formatTranscriptMarkdown(transcript ytdl.VideoTranscript) string {
	var builder strings.Builder
	wroteSegment := false

	for _, segment := range transcript {
		text := normalizeTranscriptText(segment.Text)
		if text == "" {
			continue
		}

		if wroteSegment {
			builder.WriteString("\n\n")
		}
		builder.WriteString("## ")
		builder.WriteString(formatTimestamp(segment.StartMs))
		builder.WriteString("\n")
		builder.WriteString(text)
		wroteSegment = true
	}

	return builder.String()
}

func formatSTTMarkdown(text string) string {
	text = normalizeTranscriptText(text)
	if text == "" {
		return ""
	}

	return "## 00:00\n" + text
}

func normalizeTranscriptText(text string) string {
	words := strings.Fields(text)
	return strings.TrimSpace(strings.Join(words, " "))
}

func formatTimestamp(startMs int) string {
	if startMs < 0 {
		startMs = 0
	}

	totalSeconds := startMs / 1000
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func pickAudioFormat(formats ytdl.FormatList) (*ytdl.Format, string, error) {
	audioFormats := append(ytdl.FormatList(nil), formats.Type("audio")...)
	if len(audioFormats) == 0 {
		return nil, "", &Error{
			Kind:    ErrorKindAudioUnavailable,
			Message: "no audio-only format available",
		}
	}

	audioFormats.Sort()

	for index := range audioFormats {
		formatName := normalizeAudioFormat(audioFormats[index].MimeType, audioFormats[index].URL)
		if formatName == "" {
			continue
		}
		return &audioFormats[index], formatName, nil
	}

	return nil, "", &Error{
		Kind:    ErrorKindAudioUnavailable,
		Message: "no supported audio format available",
	}
}

func normalizeAudioFormat(mimeType string, rawURL string) string {
	if mediaType, _, err := mime.ParseMediaType(mimeType); err == nil {
		switch strings.ToLower(strings.TrimSpace(mediaType)) {
		case "audio/mp4":
			return "m4a"
		case "audio/mpeg":
			return "mp3"
		case "audio/wav", "audio/x-wav":
			return "wav"
		case "audio/flac", "audio/x-flac":
			return "flac"
		case "audio/aac", "audio/x-aac":
			return "aac"
		case "audio/ogg", "audio/opus":
			return "ogg"
		case "audio/webm":
			return "webm"
		}
	}

	if parsedURL, err := url.Parse(rawURL); err == nil {
		switch strings.TrimPrefix(strings.ToLower(path.Ext(parsedURL.Path)), ".") {
		case "m4a", "mp3", "wav", "flac", "aac", "ogg", "webm":
			return strings.TrimPrefix(strings.ToLower(path.Ext(parsedURL.Path)), ".")
		}
	}

	return ""
}

func firstPathSegment(value string) string {
	segments := pathSegments(value)
	if len(segments) == 0 {
		return ""
	}
	return segments[0]
}

func pathSegments(value string) []string {
	trimmed := strings.Trim(strings.TrimSpace(value), "/")
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "/")
}

func isContextError(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}
