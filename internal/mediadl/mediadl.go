// Package mediadl extracts transcripts and metadata from media URLs (YouTube,
// Instagram, and other yt-dlp supported sources) and transcribes audio through an
// STT provider. Platform-specific URL parsing lives in the caller; this package
// operates on an already-parsed ParsedURL.
package mediadl

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/compozy/kb/internal/config"
)

const (
	transcriptUnavailableMessage = "captions unavailable"
	defaultRetryBackoff          = "1s"
)

var errTranscriptDisabled = errors.New("transcript disabled")

// TranscriptSource identifies how the transcript was produced.
type TranscriptSource string

const (
	// TranscriptSourceCaptions means platform captions were fetched directly.
	TranscriptSourceCaptions TranscriptSource = "captions"
	// TranscriptSourceSTT means audio was transcribed through an STT provider.
	TranscriptSourceSTT TranscriptSource = "stt"
)

// CaptionKind identifies whether fetched captions were manual or ASR.
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

// ErrorKind categorizes user-facing media extraction failures.
type ErrorKind string

const (
	// ErrorKindInvalidURL reports a malformed or unsupported media URL.
	ErrorKindInvalidURL ErrorKind = "invalid_url"
	// ErrorKindUnavailable reports media that cannot be accessed.
	ErrorKindUnavailable ErrorKind = "unavailable"
	// ErrorKindPrivate reports private media.
	ErrorKindPrivate ErrorKind = "private"
	// ErrorKindAgeRestricted reports media that requires login/age confirmation.
	ErrorKindAgeRestricted ErrorKind = "age_restricted"
	// ErrorKindTranscriptUnavailable reports missing or disabled captions/transcript.
	ErrorKindTranscriptUnavailable ErrorKind = "transcript_unavailable"
	// ErrorKindAudioUnavailable reports that no supported audio stream could be downloaded.
	ErrorKindAudioUnavailable ErrorKind = "audio_unavailable"
	// ErrorKindRateLimited reports a rate-limit response such as HTTP 429.
	ErrorKindRateLimited ErrorKind = "rate_limited"
	// ErrorKindNetworkBlocked reports a network block or bot-detection response.
	ErrorKindNetworkBlocked ErrorKind = "network_blocked"
)

// Error carries structured failure details for callers that need to branch on
// platform-specific failure modes.
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
		target = "media"
	}

	message := strings.TrimSpace(err.Message)
	if message == "" {
		message = string(err.Kind)
	}

	if err.Err != nil {
		return fmt.Sprintf("media %s for %q: %s: %v", err.Kind, target, message, err.Err)
	}

	return fmt.Sprintf("media %s for %q: %s", err.Kind, target, message)
}

// Unwrap returns the underlying cause.
func (err *Error) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.Err
}

// Metadata contains normalized media metadata.
type Metadata struct {
	VideoID              string
	URL                  string
	Title                string
	Description          string
	Channel              string
	ChannelID            string
	UploaderID           string
	Duration             time.Duration
	DurationString       string
	PublishDate          time.Time
	ViewCount            *int64
	LikeCount            *int64
	CommentCount         *int64
	ChannelFollowerCount *int64
	Categories           []string
	VideoTags            []string
	Language             string
	LiveStatus           string
	WasLive              *bool
	ChapterCount         int
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
	TranscriptionPolicy     TranscriptionPolicy
	PreferredLanguages      []string
	AllowTranslatedCaptions bool
}

// ParsedURL is a platform-neutral, already-validated media reference. Callers
// (the youtube/instagram shims) parse raw URLs into this form before extraction.
type ParsedURL struct {
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

// NewExtractor constructs an extractor from a platform-neutral backend config and
// STT provider configuration.
func NewExtractor(
	backendConfig BackendConfig,
	sttConfig config.STTConfig,
	openRouterConfig config.OpenRouterConfig,
	whisperCPPConfigs ...config.WhisperCPPConfig,
) *Extractor {
	backoff, err := parseRetryBackoff(backendConfig.RetryBackoff)
	attempts := max(backendConfig.RetryAttempts, 1)

	retry := retryPolicy{
		Attempts: attempts,
		Backoff:  backoff,
	}

	return &Extractor{
		ytDLP:     newYTDLPBackend(backendConfig, retry),
		stt:       NewTranscriber(sttConfig, openRouterConfig, whisperCPPConfigs...),
		sttConfig: normalizeSTTConfig(sttConfig),
		setupErr:  err,
	}
}

func parseRetryBackoff(value string) (time.Duration, error) {
	backoff := strings.TrimSpace(value)
	if backoff == "" {
		backoff = defaultRetryBackoff
	}
	duration, err := time.ParseDuration(backoff)
	if err != nil {
		return 0, fmt.Errorf("retry_backoff must be a Go duration: %q", value)
	}
	if duration < 0 {
		return 0, fmt.Errorf("retry_backoff cannot be negative: %q", value)
	}
	return duration, nil
}

// Extract fetches media metadata and transcript markdown for an already-parsed
// media reference.
func (extractor *Extractor) Extract(ctx context.Context, parsed ParsedURL, options ExtractOptions) (*Result, error) {
	if extractor == nil {
		return nil, errors.New("media extract: extractor is nil")
	}
	if extractor.ytDLP == nil {
		return nil, errors.New("media extract: yt-dlp backend is required")
	}
	if extractor.setupErr != nil {
		return nil, fmt.Errorf("media extract: configure yt-dlp backend: %w", extractor.setupErr)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if strings.TrimSpace(parsed.CanonicalURL) == "" {
		return nil, &Error{Kind: ErrorKindInvalidURL, Message: "media url is required"}
	}

	return extractor.extractWithYTDLPBackend(ctx, parsed, options)
}

// ListPlaylist resolves a channel or playlist URL with top-level metadata.
func (extractor *Extractor) ListPlaylist(ctx context.Context, channelURL string, limit int) (PlaylistListing, error) {
	if extractor == nil || extractor.ytDLP == nil {
		return PlaylistListing{}, errors.New("media channel: yt-dlp backend is required")
	}
	if extractor.setupErr != nil {
		return PlaylistListing{}, fmt.Errorf("media channel: configure yt-dlp backend: %w", extractor.setupErr)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return extractor.ytDLP.ListPlaylist(ctx, channelURL, limit)
}

func (extractor *Extractor) extractWithYTDLPBackend(
	ctx context.Context,
	parsed ParsedURL,
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
	captionResult, err := extractor.ytDLP.extractCaptionsFromInfo(
		ctx,
		parsed,
		info,
		options.PreferredLanguages,
		allowAutomaticCaptions,
		options.AllowTranslatedCaptions,
	)
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

func newNetworkBlockedError(videoID string, message string, err error) error {
	return &Error{
		Kind:    ErrorKindNetworkBlocked,
		VideoID: videoID,
		Message: message,
		Err:     err,
	}
}

func newRateLimitedError(videoID string, message string, err error) error {
	return &Error{
		Kind:    ErrorKindRateLimited,
		VideoID: videoID,
		Message: message,
		Err:     err,
	}
}

// IsTranscriptUnavailable reports whether err signals that no transcript could be
// produced (no captions and no usable STT output).
func IsTranscriptUnavailable(err error) bool {
	return isTranscriptUnavailable(err)
}

func isTranscriptUnavailable(err error) bool {
	var mediaErr *Error
	return errors.As(err, &mediaErr) && mediaErr.Kind == ErrorKindTranscriptUnavailable
}

// IsNetworkBlocked reports whether err is a structured network-blocked error or a
// raw yt-dlp network-block diagnostic.
func IsNetworkBlocked(err error) bool {
	var mediaErr *Error
	if errors.As(err, &mediaErr) && mediaErr.Kind == ErrorKindNetworkBlocked {
		return true
	}
	return isYTDLPNetworkBlockedError(err)
}

// IsRateLimited reports whether err is a structured rate-limit error or a raw
// yt-dlp rate-limit diagnostic.
func IsRateLimited(err error) bool {
	var mediaErr *Error
	if errors.As(err, &mediaErr) && mediaErr.Kind == ErrorKindRateLimited {
		return true
	}
	return isYTDLPRateLimitedError(err)
}

// IsContextError reports whether err is a context cancellation or deadline error.
func IsContextError(err error) bool {
	return isContextError(err)
}

func isStatusCodeRateLimited(err error) bool {
	status, ok := unexpectedStatusCode(err)
	return ok && status == http.StatusTooManyRequests
}

func isStatusCodeNetworkBlocked(err error) bool {
	status, ok := unexpectedStatusCode(err)
	if !ok {
		return false
	}
	return status == http.StatusBadRequest ||
		status == http.StatusForbidden ||
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

func isContextError(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}
