// Package youtube extracts transcripts and metadata from YouTube videos.
package youtube

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/compozy/kb/internal/config"
)

const (
	transcriptUnavailableMessage = "captions unavailable"
)

var youtubeVideoIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)

var errTranscriptDisabled = errors.New("transcript disabled")

// TranscriptSource identifies how the transcript was produced.
type TranscriptSource string

const (
	// TranscriptSourceCaptions means YouTube captions were fetched directly.
	TranscriptSourceCaptions TranscriptSource = "captions"
	// TranscriptSourceSTT means audio was transcribed through an STT provider.
	TranscriptSourceSTT TranscriptSource = "stt"
)

// CaptionKind identifies whether fetched YouTube captions were manual or ASR.
type CaptionKind string

const (
	CaptionKindManual    CaptionKind = "manual"
	CaptionKindAutomatic CaptionKind = "automatic"
)

// TranscriptionPolicy controls whether captions or STT are used.
type TranscriptionPolicy string

const (
	TranscriptionPolicyCaptions TranscriptionPolicy = "captions"
	TranscriptionPolicyAuto     TranscriptionPolicy = "auto"
	TranscriptionPolicySTT      TranscriptionPolicy = "stt"
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
	Metadata            Metadata
	Markdown            string
	Source              TranscriptSource
	Language            string
	CaptionKind         CaptionKind
	TranscriptionPolicy TranscriptionPolicy
	STTProvider         string
	STTModel            string
}

// ExtractOptions controls transcript extraction behavior.
type ExtractOptions struct {
	TranscriptionPolicy TranscriptionPolicy
	PreferredLanguages  []string
}

type parsedVideoURL struct {
	CanonicalURL string
	VideoID      string
}

type transcriptSegment struct {
	StartMs int
	Text    string
}

type sttClient interface {
	Provider() string
	Model() string
	Transcribe(ctx context.Context, audio []byte, format string) (string, error)
}

type retryPolicy struct {
	Attempts int
	Backoff  time.Duration
}

// Extractor orchestrates transcript extraction and optional STT fallback.
type Extractor struct {
	ytDLP     *ytDLPBackend
	stt       sttClient
	sttConfig config.STTConfig
	setupErr  error
}

// NewExtractorWithConfig constructs an extractor with explicit STT provider
// configuration.
func NewExtractorWithConfig(
	sttConfig config.STTConfig,
	openRouterConfig config.OpenRouterConfig,
	youtubeConfigs ...config.YouTubeConfig,
) *Extractor {
	youtubeConfig := config.Default().YouTube
	if len(youtubeConfigs) > 0 {
		youtubeConfig = youtubeConfigs[0]
	}
	backoff, err := youtubeConfig.RetryBackoffDuration()
	attempts := youtubeConfig.RetryAttempts
	if attempts < 1 {
		attempts = 1
	}

	retry := retryPolicy{
		Attempts: attempts,
		Backoff:  backoff,
	}

	return &Extractor{
		ytDLP:     newYTDLPBackend(youtubeConfig, retry),
		stt:       NewTranscriber(sttConfig, openRouterConfig),
		sttConfig: normalizeSTTConfig(sttConfig),
		setupErr:  err,
	}
}

// Extract fetches video metadata and transcript markdown from a YouTube URL.
func (extractor *Extractor) Extract(ctx context.Context, rawURL string, options ExtractOptions) (*Result, error) {
	if extractor == nil {
		return nil, errors.New("youtube extract: extractor is nil")
	}
	if extractor.ytDLP == nil {
		return nil, errors.New("youtube extract: yt-dlp backend is required")
	}
	if extractor.setupErr != nil {
		return nil, fmt.Errorf("youtube extract: configure yt-dlp backend: %w", extractor.setupErr)
	}
	if ctx == nil {
		ctx = context.Background()
	}

	parsed, err := parseVideoURL(rawURL)
	if err != nil {
		return nil, err
	}

	return extractor.extractWithYTDLPBackend(ctx, parsed, options)
}

func (extractor *Extractor) extractWithYTDLPBackend(
	ctx context.Context,
	parsed parsedVideoURL,
	options ExtractOptions,
) (*Result, error) {
	policy := normalizeTranscriptionPolicy(options)
	info, err := extractor.ytDLP.loadInfo(ctx, parsed.CanonicalURL)
	if err != nil {
		if isContextError(err) {
			return nil, err
		}
		return nil, fmt.Errorf("yt-dlp backend: %w", err)
	}

	result := &Result{
		Metadata:            metadataFromYTDLPInfo(parsed, info),
		TranscriptionPolicy: policy,
	}
	if policy == TranscriptionPolicySTT {
		return extractor.extractSTTFromYTDLP(ctx, parsed, info, result, nil)
	}

	allowAutomaticCaptions := policy == TranscriptionPolicyCaptions
	captionResult, err := extractor.ytDLP.extractCaptionsFromInfo(ctx, parsed, info, options.PreferredLanguages, allowAutomaticCaptions)
	if err == nil {
		captionResult.TranscriptionPolicy = policy
		return captionResult, nil
	}
	if captionResult != nil {
		result = captionResult
		result.TranscriptionPolicy = policy
	}
	if isContextError(err) {
		return result, err
	}
	if isTranscriptUnavailable(err) && policy == TranscriptionPolicyAuto {
		return extractor.extractSTTFromYTDLP(ctx, parsed, info, result, err)
	}
	return result, err
}

func (extractor *Extractor) shouldAttemptSTT(policy TranscriptionPolicy) bool {
	if extractor.stt == nil {
		return false
	}
	return policy == TranscriptionPolicyAuto || policy == TranscriptionPolicySTT
}

func normalizeTranscriptionPolicy(options ExtractOptions) TranscriptionPolicy {
	policy, err := ParseTranscriptionPolicy(string(options.TranscriptionPolicy))
	if err == nil {
		return policy
	}
	return TranscriptionPolicyCaptions
}

// ParseTranscriptionPolicy validates a transcription policy string.
func ParseTranscriptionPolicy(value string) (TranscriptionPolicy, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case string(TranscriptionPolicyAuto):
		return TranscriptionPolicyAuto, nil
	case string(TranscriptionPolicySTT):
		return TranscriptionPolicySTT, nil
	case "", string(TranscriptionPolicyCaptions):
		return TranscriptionPolicyCaptions, nil
	default:
		return "", fmt.Errorf("transcription policy must be captions, auto, or stt: %q", value)
	}
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

func formatTranscriptMarkdown(transcript []transcriptSegment) string {
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
