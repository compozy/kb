// Package youtube extracts transcripts and metadata from YouTube videos by
// parsing YouTube URLs and delegating extraction to the shared mediadl core.
package youtube

import (
	"context"
	"errors"
	"net/url"
	"regexp"
	"strings"

	"github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/mediadl"
)

// Shared types re-exported from mediadl so existing callers (CLI, channel
// ingestion) keep compiling against the youtube package surface.
type (
	Result              = mediadl.Result
	Metadata            = mediadl.Metadata
	ExtractOptions      = mediadl.ExtractOptions
	TranscriptionPolicy = mediadl.TranscriptionPolicy
	TranscriptSource    = mediadl.TranscriptSource
	CaptionKind         = mediadl.CaptionKind
	Error               = mediadl.Error
	ErrorKind           = mediadl.ErrorKind
)

const (
	TranscriptionPolicyCaptions = mediadl.TranscriptionPolicyCaptions
	TranscriptionPolicyAuto     = mediadl.TranscriptionPolicyAuto
	TranscriptionPolicySTT      = mediadl.TranscriptionPolicySTT

	TranscriptSourceCaptions = mediadl.TranscriptSourceCaptions
	TranscriptSourceSTT      = mediadl.TranscriptSourceSTT

	CaptionKindManual    = mediadl.CaptionKindManual
	CaptionKindAutomatic = mediadl.CaptionKindAutomatic

	ErrorKindInvalidURL            = mediadl.ErrorKindInvalidURL
	ErrorKindUnavailable           = mediadl.ErrorKindUnavailable
	ErrorKindPrivate               = mediadl.ErrorKindPrivate
	ErrorKindAgeRestricted         = mediadl.ErrorKindAgeRestricted
	ErrorKindTranscriptUnavailable = mediadl.ErrorKindTranscriptUnavailable
	ErrorKindAudioUnavailable      = mediadl.ErrorKindAudioUnavailable
	ErrorKindRateLimited           = mediadl.ErrorKindRateLimited
	ErrorKindNetworkBlocked        = mediadl.ErrorKindNetworkBlocked
)

// ParseTranscriptionPolicy validates a transcription policy string.
var ParseTranscriptionPolicy = mediadl.ParseTranscriptionPolicy

var youtubeVideoIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)

// Extractor extracts YouTube transcripts by parsing YouTube URLs and delegating
// to the shared mediadl core.
type Extractor struct {
	core *mediadl.Extractor
	// extractFn overrides per-video extraction during bulk runs. It is a test
	// seam; when nil the exported Extract method is used.
	extractFn func(context.Context, string, ExtractOptions) (*Result, error)
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
	backendConfig := mediadl.BackendConfig{
		Platform:      "youtube",
		YTDLPPath:     youtubeConfig.YTDLPPath,
		Proxy:         youtubeConfig.Proxy,
		CookiesFile:   youtubeConfig.CookiesFile,
		UserAgent:     youtubeConfig.UserAgent,
		RetryAttempts: youtubeConfig.RetryAttempts,
		RetryBackoff:  youtubeConfig.RetryBackoff,
	}
	return &Extractor{core: mediadl.NewExtractor(backendConfig, sttConfig, openRouterConfig)}
}

// Extract fetches video metadata and transcript markdown from a YouTube URL.
func (extractor *Extractor) Extract(ctx context.Context, rawURL string, options ExtractOptions) (*Result, error) {
	if extractor == nil || extractor.core == nil {
		return nil, errors.New("youtube extract: extractor is not configured")
	}
	parsed, err := parseVideoURL(rawURL)
	if err != nil {
		return nil, err
	}
	return extractor.core.Extract(ctx, parsed, options)
}

func parseVideoURL(rawURL string) (mediadl.ParsedURL, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return mediadl.ParsedURL{}, &Error{
			Kind:    ErrorKindInvalidURL,
			URL:     rawURL,
			Message: "youtube url is required",
		}
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return mediadl.ParsedURL{}, &Error{
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
		return mediadl.ParsedURL{}, &Error{
			Kind:    ErrorKindInvalidURL,
			URL:     rawURL,
			Message: "unsupported YouTube host",
		}
	}

	if !youtubeVideoIDPattern.MatchString(videoID) {
		return mediadl.ParsedURL{}, &Error{
			Kind:    ErrorKindInvalidURL,
			URL:     rawURL,
			Message: "invalid YouTube video identifier",
		}
	}

	return mediadl.ParsedURL{
		CanonicalURL: "https://www.youtube.com/watch?v=" + videoID,
		VideoID:      videoID,
	}, nil
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
