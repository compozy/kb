// Package instagram extracts captions and transcripts from Instagram video URLs
// (reels, video posts, and IGTV) by parsing Instagram URLs and delegating audio
// transcription to the shared mediadl core.
package instagram

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/compozy/kb/internal/config"
	"github.com/compozy/kb/internal/mediadl"
)

// transcriptSourceNone marks a document whose body is the caption only because no
// transcript could be produced (for example a music-only reel).
const transcriptSourceNone = mediadl.TranscriptSource("none")

var instagramShortcodePattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

// transcriptExtractor is the subset of mediadl.Extractor used by this package. It
// is an interface so tests can inject a stub without invoking yt-dlp.
type transcriptExtractor interface {
	Extract(ctx context.Context, parsed mediadl.ParsedURL, options mediadl.ExtractOptions) (*mediadl.Result, error)
}

// Extractor extracts Instagram video transcripts by parsing Instagram URLs and
// delegating to the shared mediadl core.
type Extractor struct {
	core transcriptExtractor
}

// NewExtractorWithConfig constructs an extractor with explicit STT provider and
// Instagram network configuration.
func NewExtractorWithConfig(
	sttConfig config.STTConfig,
	openRouterConfig config.OpenRouterConfig,
	whisperCPPConfig config.WhisperCPPConfig,
	instagramConfig config.InstagramConfig,
) *Extractor {
	backendConfig := mediadl.BackendConfig{
		Platform:      "instagram",
		YTDLPPath:     instagramConfig.YTDLPPath,
		Proxy:         instagramConfig.Proxy,
		CookiesFile:   instagramConfig.CookiesFile,
		UserAgent:     instagramConfig.UserAgent,
		RetryAttempts: instagramConfig.RetryAttempts,
		RetryBackoff:  instagramConfig.RetryBackoff,
	}
	return &Extractor{core: mediadl.NewExtractor(backendConfig, sttConfig, openRouterConfig, whisperCPPConfig)}
}

// Extract fetches metadata, caption, and transcript for an Instagram video URL.
// The document body is the post caption followed by the spoken transcript. When
// no transcript can be produced but a caption exists, it degrades to a
// caption-only document instead of failing.
func (extractor *Extractor) Extract(ctx context.Context, rawURL string, options mediadl.ExtractOptions) (*mediadl.Result, error) {
	if extractor == nil || extractor.core == nil {
		return nil, errors.New("instagram extract: extractor is not configured")
	}
	parsed, err := parseInstagramURL(rawURL)
	if err != nil {
		return nil, err
	}

	result, err := extractor.core.Extract(ctx, parsed, options)
	if err != nil {
		if mediadl.IsTranscriptUnavailable(err) && result != nil && captionOf(result) != "" {
			backfillMetadataURL(result, parsed.CanonicalURL)
			result.Markdown = composeBody(captionOf(result), "")
			result.Source = transcriptSourceNone
			result.Language = ""
			result.CaptionKind = ""
			return result, nil
		}
		return nil, err
	}

	backfillMetadataURL(result, parsed.CanonicalURL)
	result.Markdown = composeBody(captionOf(result), result.Markdown)
	return result, nil
}

func backfillMetadataURL(result *mediadl.Result, canonicalURL string) {
	if result == nil {
		return
	}
	if strings.TrimSpace(result.Metadata.URL) == "" {
		result.Metadata.URL = canonicalURL
	}
}

func captionOf(result *mediadl.Result) string {
	if result == nil {
		return ""
	}
	return strings.TrimSpace(result.Metadata.Description)
}

func composeBody(caption string, transcript string) string {
	caption = strings.TrimSpace(caption)
	transcript = strings.TrimSpace(transcript)
	sections := make([]string, 0, 2)
	if caption != "" {
		sections = append(sections, "## Caption\n"+caption)
	}
	if transcript != "" {
		sections = append(sections, "## Transcript\n"+transcript)
	}
	return strings.Join(sections, "\n\n")
}

func parseInstagramURL(rawURL string) (mediadl.ParsedURL, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return mediadl.ParsedURL{}, invalidInstagramURL(rawURL, "instagram url is required", nil)
	}

	value := rawURL
	if !strings.Contains(value, "://") {
		value = "https://" + value
	}
	parsedURL, err := url.Parse(value)
	if err != nil || parsedURL.Host == "" {
		return mediadl.ParsedURL{}, invalidInstagramURL(rawURL, "invalid Instagram URL", err)
	}

	switch strings.ToLower(parsedURL.Host) {
	case "instagram.com", "www.instagram.com", "m.instagram.com":
	default:
		return mediadl.ParsedURL{}, invalidInstagramURL(rawURL, "unsupported Instagram host", nil)
	}

	kind, shortcode, ok := instagramKindAndShortcode(parsedURL.Path)
	if !ok {
		return mediadl.ParsedURL{}, invalidInstagramURL(rawURL, "expected an Instagram reel, post, or tv URL", nil)
	}

	return mediadl.ParsedURL{
		CanonicalURL: fmt.Sprintf("https://www.instagram.com/%s/%s/", kind, shortcode),
		VideoID:      shortcode,
	}, nil
}

// instagramKindAndShortcode locates the content kind (reel/p/tv) and the shortcode
// in a path, supporting both bare (/reel/<code>) and account-prefixed
// (/<user>/reel/<code>) forms. The "reels" alias is normalized to "reel".
func instagramKindAndShortcode(path string) (string, string, bool) {
	segments := pathSegments(path)
	for index := 0; index+1 < len(segments); index++ {
		kind := strings.ToLower(segments[index])
		switch kind {
		case "reel", "reels", "p", "tv":
			shortcode := segments[index+1]
			if !instagramShortcodePattern.MatchString(shortcode) {
				return "", "", false
			}
			if kind == "reels" {
				kind = "reel"
			}
			return kind, shortcode, true
		}
	}
	return "", "", false
}

func pathSegments(value string) []string {
	trimmed := strings.Trim(strings.TrimSpace(value), "/")
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "/")
}

func invalidInstagramURL(rawURL string, message string, err error) error {
	return &mediadl.Error{
		Kind:    mediadl.ErrorKindInvalidURL,
		URL:     rawURL,
		Message: message,
		Err:     err,
	}
}
