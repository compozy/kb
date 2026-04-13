// Package youtube extracts transcripts and metadata from YouTube videos.
package youtube

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/url"
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

// Extractor orchestrates transcript extraction and optional STT fallback.
type Extractor struct {
	youtube youtubeClient
	stt     sttClient
}

// NewExtractor constructs a default extractor backed by kkdai/youtube and the
// OpenRouter STT client.
func NewExtractor(cfg config.OpenRouterConfig) *Extractor {
	return &Extractor{
		youtube: &ytdl.Client{},
		stt:     NewOpenRouterClient(cfg),
	}
}

// Extract fetches video metadata and transcript markdown from a YouTube URL.
func (extractor *Extractor) Extract(ctx context.Context, rawURL string, options ExtractOptions) (*Result, error) {
	if extractor == nil {
		return nil, errors.New("youtube extract: extractor is nil")
	}
	if extractor.youtube == nil {
		return nil, errors.New("youtube extract: client is nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	parsed, err := parseVideoURL(rawURL)
	if err != nil {
		return nil, err
	}

	video, err := extractor.youtube.GetVideoContext(ctx, parsed.CanonicalURL)
	if err != nil {
		return nil, wrapVideoError(parsed, err)
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

	if !extractor.shouldAttemptSTT(options) {
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
		transcript, err := extractor.youtube.GetTranscriptCtx(ctx, video, track.LanguageCode)
		if err != nil {
			if isContextError(err) {
				return "", "", err
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

	stream, _, err := extractor.youtube.GetStreamContext(ctx, video, format)
	if err != nil {
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
