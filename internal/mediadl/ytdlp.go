package mediadl

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const ytdlpDefaultBinary = "yt-dlp"

var errYTDLPUnavailable = errors.New("yt-dlp unavailable")
var errYTDLPCaptionFetchFailed = errors.New("yt-dlp caption fetch failed")

type commandContextFunc func(context.Context, string, ...string) *exec.Cmd
type lookPathFunc func(string) (string, error)

// BackendConfig holds the platform-neutral yt-dlp network knobs. Platform labels
// the source ("youtube", "instagram") so diagnostics can point at the right TOML
// section. RetryAttempts and RetryBackoff are consumed by NewExtractor.
type BackendConfig struct {
	Platform      string
	YTDLPPath     string
	Proxy         string
	CookiesFile   string
	UserAgent     string
	RetryAttempts int
	RetryBackoff  string
}

type ytDLPBackend struct {
	binaryPath     string
	cfg            BackendConfig
	retry          retryPolicy
	lookPath       lookPathFunc
	commandContext commandContextFunc
}

type ytDLPInfo struct {
	ID                   string                           `json:"id"`
	Title                string                           `json:"title"`
	Description          string                           `json:"description"`
	Uploader             string                           `json:"uploader"`
	UploaderID           string                           `json:"uploader_id"`
	Channel              string                           `json:"channel"`
	ChannelID            string                           `json:"channel_id"`
	ChannelFollowerCount *int64                           `json:"channel_follower_count"`
	Duration             float64                          `json:"duration"`
	DurationString       string                           `json:"duration_string"`
	UploadDate           string                           `json:"upload_date"`
	WebpageURL           string                           `json:"webpage_url"`
	ViewCount            *int64                           `json:"view_count"`
	LikeCount            *int64                           `json:"like_count"`
	CommentCount         *int64                           `json:"comment_count"`
	Categories           []string                         `json:"categories"`
	Tags                 []string                         `json:"tags"`
	Language             string                           `json:"language"`
	LiveStatus           string                           `json:"live_status"`
	WasLive              *bool                            `json:"was_live"`
	Chapters             []ytDLPChapter                   `json:"chapters"`
	Subtitles            map[string][]ytDLPSubtitleFormat `json:"subtitles"`
	AutomaticCaptions    map[string][]ytDLPSubtitleFormat `json:"automatic_captions"`
}

type ytDLPSubtitleFormat struct {
	Ext string `json:"ext"`
}

type ytDLPChapter struct{}

type ytDLPCaptionCandidate struct {
	Language     string
	BaseLanguage string
	Automatic    bool
	Original     bool
	Translated   bool
	OriginalKey  bool
}

type ytDLPDownloadedAudio struct {
	Path    string
	Format  string
	cleanup func()
}

// PlaylistEntry is a single raw entry resolved from a channel or playlist URL.
// Callers apply platform-specific filtering and URL canonicalization.
type PlaylistEntry struct {
	ID         string
	Title      string
	URL        string
	WebpageURL string
}

// PlaylistListing is the top-level yt-dlp playlist/channel response plus entries.
type PlaylistListing struct {
	Title    string
	Channel  string
	Uploader string
	Entries  []PlaylistEntry
}

type ytDLPPlaylist struct {
	Title    string               `json:"title"`
	Channel  string               `json:"channel"`
	Uploader string               `json:"uploader"`
	Entries  []ytDLPPlaylistEntry `json:"entries"`
}

type ytDLPPlaylistEntry struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	URL        string `json:"url"`
	WebpageURL string `json:"webpage_url"`
}

func (audio ytDLPDownloadedAudio) Cleanup() {
	if audio.cleanup != nil {
		audio.cleanup()
	}
}

func newYTDLPBackend(cfg BackendConfig, retry retryPolicy) *ytDLPBackend {
	binaryPath := strings.TrimSpace(cfg.YTDLPPath)
	if binaryPath == "" {
		binaryPath = ytdlpDefaultBinary
	}
	return &ytDLPBackend{
		binaryPath:     binaryPath,
		cfg:            cfg,
		retry:          retry,
		lookPath:       exec.LookPath,
		commandContext: exec.CommandContext,
	}
}

func (backend *ytDLPBackend) Extract(
	ctx context.Context,
	parsed ParsedURL,
	preferredLanguages []string,
) (*Result, error) {
	return backend.ExtractCaptions(ctx, parsed, preferredLanguages, true, false)
}

func (backend *ytDLPBackend) ExtractCaptions(
	ctx context.Context,
	parsed ParsedURL,
	preferredLanguages []string,
	allowAutomatic bool,
	allowTranslated bool,
) (*Result, error) {
	if backend == nil {
		return nil, fmt.Errorf("%w: backend is nil", errYTDLPUnavailable)
	}
	if ctx == nil {
		ctx = context.Background()
	}

	info, err := backend.loadInfo(ctx, parsed.CanonicalURL)
	if err != nil {
		return nil, err
	}

	return backend.extractCaptionsFromInfo(ctx, parsed, info, preferredLanguages, allowAutomatic, allowTranslated)
}

func (backend *ytDLPBackend) extractCaptionsFromInfo(
	ctx context.Context,
	parsed ParsedURL,
	info ytDLPInfo,
	preferredLanguages []string,
	allowAutomatic bool,
	allowTranslated bool,
) (*Result, error) {
	result := &Result{Metadata: metadataFromYTDLPInfo(parsed, info)}
	candidate, ok := selectYTDLPCaption(info, preferredLanguages, allowAutomatic, allowTranslated)
	if !ok {
		return result, &Error{
			Kind:    ErrorKindTranscriptUnavailable,
			URL:     parsed.CanonicalURL,
			VideoID: parsed.VideoID,
			Message: ytDLPCaptionUnavailableMessage(preferredLanguages, allowTranslated),
			Err:     errTranscriptDisabled,
		}
	}

	markdown, err := backend.downloadCaption(ctx, parsed.CanonicalURL, candidate)
	if err != nil {
		return result, err
	}
	if strings.TrimSpace(markdown) == "" {
		return result, fmt.Errorf("%w: empty transcript for %q", errYTDLPCaptionFetchFailed, candidate.Language)
	}

	result.Markdown = markdown
	result.Source = TranscriptSourceCaptions
	result.Language = candidate.Language
	if candidate.Automatic {
		result.CaptionKind = CaptionKindAutomatic
	} else {
		result.CaptionKind = CaptionKindManual
	}
	return result, nil
}

func (backend *ytDLPBackend) loadInfo(ctx context.Context, rawURL string) (ytDLPInfo, error) {
	stdout, _, err := backend.run(ctx, "metadata", append(backend.baseArgs(),
		"--dump-single-json",
		"--skip-download",
		rawURL,
	)...)
	if err != nil {
		return ytDLPInfo{}, backend.wrapMetadataFetchError(err)
	}

	var info ytDLPInfo
	if err := json.Unmarshal([]byte(stdout), &info); err != nil {
		return ytDLPInfo{}, fmt.Errorf("yt-dlp metadata: parse JSON: %w", err)
	}
	return info, nil
}

// ListPlaylist resolves a channel or playlist URL. A limit greater than zero caps
// the result to the newest uploads. Callers apply platform-specific filtering and
// URL canonicalization on the returned entries.
func (backend *ytDLPBackend) ListPlaylist(ctx context.Context, channelURL string, limit int) (PlaylistListing, error) {
	if backend == nil {
		return PlaylistListing{}, fmt.Errorf("%w: backend is nil", errYTDLPUnavailable)
	}
	if ctx == nil {
		ctx = context.Background()
	}

	args := append(backend.playlistArgs(),
		"--flat-playlist",
		"--dump-single-json",
		"--skip-download",
	)
	if limit > 0 {
		args = append(args, "--playlist-end", strconv.Itoa(limit))
	}
	args = append(args, channelURL)

	stdout, _, err := backend.run(ctx, "channel", args...)
	if err != nil {
		return PlaylistListing{}, backend.wrapMetadataFetchError(err)
	}

	var playlist ytDLPPlaylist
	if err := json.Unmarshal([]byte(stdout), &playlist); err != nil {
		return PlaylistListing{}, fmt.Errorf("yt-dlp channel: parse JSON: %w", err)
	}

	entries := make([]PlaylistEntry, 0, len(playlist.Entries))
	for _, entry := range playlist.Entries {
		entries = append(entries, PlaylistEntry{
			ID:         strings.TrimSpace(entry.ID),
			Title:      strings.TrimSpace(entry.Title),
			URL:        strings.TrimSpace(entry.URL),
			WebpageURL: strings.TrimSpace(entry.WebpageURL),
		})
	}
	return PlaylistListing{
		Title:    strings.TrimSpace(playlist.Title),
		Channel:  strings.TrimSpace(playlist.Channel),
		Uploader: strings.TrimSpace(playlist.Uploader),
		Entries:  entries,
	}, nil
}

func (backend *ytDLPBackend) downloadCaption(
	ctx context.Context,
	rawURL string,
	candidate ytDLPCaptionCandidate,
) (string, error) {
	tempDir, err := os.MkdirTemp("", "kb-ytdlp-*")
	if err != nil {
		return "", fmt.Errorf("yt-dlp captions: create temp dir: %w", err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	args := append(backend.baseArgs(),
		"--skip-download",
		"--sub-langs", candidate.Language,
		"--sub-format", "json3/vtt/best",
		"--paths", "home:"+tempDir,
		"--output", "%(id)s.%(ext)s",
	)
	if candidate.Automatic {
		args = append(args, "--write-auto-subs")
	} else {
		args = append(args, "--write-subs")
	}
	args = append(args, rawURL)

	if _, _, err := backend.run(ctx, "captions", args...); err != nil {
		return "", backend.wrapCaptionFetchError(err, candidate.Language)
	}

	data, ext, err := readYTDLPSubtitleFile(tempDir)
	if err != nil {
		return "", backend.wrapCaptionFetchError(err, candidate.Language)
	}
	switch ext {
	case ".json3":
		markdown, err := formatYTDLPJSON3Transcript(data)
		if err != nil {
			return "", backend.wrapCaptionFetchError(err, candidate.Language)
		}
		return markdown, nil
	case ".vtt":
		markdown, err := formatYTDLPVTTTranscript(data)
		if err != nil {
			return "", backend.wrapCaptionFetchError(err, candidate.Language)
		}
		return markdown, nil
	default:
		return "", backend.wrapCaptionFetchError(
			fmt.Errorf("yt-dlp captions: unsupported subtitle extension %q", ext),
			candidate.Language,
		)
	}
}

func (backend *ytDLPBackend) downloadAudio(
	ctx context.Context,
	rawURL string,
	format string,
) (ytDLPDownloadedAudio, error) {
	if backend == nil {
		return ytDLPDownloadedAudio{}, fmt.Errorf("%w: backend is nil", errYTDLPUnavailable)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	format = strings.TrimSpace(strings.ToLower(format))
	if format == "" {
		format = "mp3"
	}
	tempDir, err := os.MkdirTemp("", "kb-ytdlp-audio-*")
	if err != nil {
		return ytDLPDownloadedAudio{}, fmt.Errorf("yt-dlp audio: create temp dir: %w", err)
	}
	cleanup := func() {
		_ = os.RemoveAll(tempDir)
	}

	args := append(backend.baseArgs(),
		"--extract-audio",
		"--audio-format", format,
		"--paths", "home:"+tempDir,
		"--output", "%(id)s.%(ext)s",
		rawURL,
	)
	if _, _, err := backend.run(ctx, "audio", args...); err != nil {
		cleanup()
		return ytDLPDownloadedAudio{}, backend.wrapAudioFetchError(err)
	}

	path, detectedFormat, err := findYTDLPAudioFile(tempDir)
	if err != nil {
		cleanup()
		return ytDLPDownloadedAudio{}, backend.wrapAudioFetchError(err)
	}
	return ytDLPDownloadedAudio{Path: path, Format: detectedFormat, cleanup: cleanup}, nil
}

// networkBlockHint renders a remediation hint that points at the right TOML
// section for the configured platform.
func (backend *ytDLPBackend) networkBlockHint() string {
	platform := strings.TrimSpace(backend.cfg.Platform)
	if platform == "" {
		platform = "youtube"
	}
	return fmt.Sprintf(
		"update yt-dlp, configure [%s].proxy or [%s].cookies_file, or run from a trusted network",
		platform, platform,
	)
}

func (backend *ytDLPBackend) wrapAudioFetchError(err error) error {
	if err == nil {
		return nil
	}
	if isYTDLPRateLimitedError(err) {
		return newRateLimitedError(
			"",
			"audio download was rate limited through yt-dlp; "+backend.networkBlockHint(),
			err,
		)
	}
	if isYTDLPNetworkBlockedError(err) {
		return newNetworkBlockedError(
			"",
			"audio download was blocked through yt-dlp; "+backend.networkBlockHint(),
			err,
		)
	}
	return fmt.Errorf("yt-dlp audio: %w", err)
}

func findYTDLPAudioFile(directory string) (string, string, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return "", "", fmt.Errorf("read audio output directory: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(entry.Name())), ".")
		if ext == "" {
			continue
		}
		switch ext {
		case "mp3", "mp4", "mpeg", "mpga", "m4a", "wav", "webm":
			return filepath.Join(directory, entry.Name()), ext, nil
		}
	}
	return "", "", errors.New("yt-dlp audio: no audio file was produced")
}

func (backend *ytDLPBackend) wrapMetadataFetchError(err error) error {
	if err == nil {
		return nil
	}
	if isYTDLPRateLimitedError(err) {
		return newRateLimitedError("", "metadata request was rate limited through yt-dlp; "+backend.networkBlockHint(), err)
	}
	if isYTDLPNetworkBlockedError(err) {
		return newNetworkBlockedError("", "metadata request was blocked through yt-dlp; "+backend.networkBlockHint(), err)
	}
	return err
}

func (backend *ytDLPBackend) wrapCaptionFetchError(err error, language string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, errYTDLPCaptionFetchFailed) {
		return err
	}
	if isYTDLPRateLimitedError(err) {
		return fmt.Errorf("%w: %w", errYTDLPCaptionFetchFailed, newRateLimitedError(
			"",
			fmt.Sprintf("caption language %q was rate limited through yt-dlp; %s", language, backend.networkBlockHint()),
			err,
		))
	}
	if isYTDLPNetworkBlockedError(err) {
		return fmt.Errorf("%w: %w", errYTDLPCaptionFetchFailed, newNetworkBlockedError(
			"",
			"captions request was blocked through yt-dlp; "+backend.networkBlockHint(),
			err,
		))
	}
	return fmt.Errorf("%w: %w", errYTDLPCaptionFetchFailed, err)
}

// baseArgs returns the shared yt-dlp flags for single-video operations. It pins
// --no-playlist so video URLs that also carry a playlist reference resolve to the
// single video.
func (backend *ytDLPBackend) baseArgs() []string {
	return backend.commonArgs(false)
}

// playlistArgs returns the shared yt-dlp flags for channel/playlist enumeration.
// It omits --no-playlist so the full uploads/playlist list is resolved.
func (backend *ytDLPBackend) playlistArgs() []string {
	return backend.commonArgs(true)
}

func (backend *ytDLPBackend) commonArgs(allowPlaylist bool) []string {
	args := []string{"--ignore-config"}
	if !allowPlaylist {
		args = append(args, "--no-playlist")
	}
	args = append(args,
		"--retries", strconv.Itoa(normalizedRetryAttempts(backend.retry)),
		"--fragment-retries", strconv.Itoa(normalizedRetryAttempts(backend.retry)),
	)
	if sleep := normalizedRetrySleep(backend.retry); sleep != "" {
		args = append(args, "--retry-sleep", sleep)
	}
	if proxy := strings.TrimSpace(backend.cfg.Proxy); proxy != "" {
		args = append(args, "--proxy", proxy)
	}
	if cookiesFile := strings.TrimSpace(backend.cfg.CookiesFile); cookiesFile != "" {
		args = append(args, "--cookies", cookiesFile)
	}
	if userAgent := strings.TrimSpace(backend.cfg.UserAgent); userAgent != "" {
		args = append(args, "--user-agent", userAgent)
	}
	return args
}

func (backend *ytDLPBackend) run(ctx context.Context, label string, args ...string) (string, string, error) {
	binaryPath, err := backend.resolveBinary()
	if err != nil {
		return "", "", err
	}

	command := backend.commandContext(ctx, binaryPath, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err = command.Run()
	stdoutText := strings.TrimSpace(stdout.String())
	stderrText := strings.TrimSpace(stderr.String())
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return stdoutText, stderrText, fmt.Errorf("yt-dlp %s: %w", label, ctxErr)
		}
		diagnostics := cleanYTDLPDiagnostics(stderrText)
		if diagnostics == "" {
			diagnostics = cleanYTDLPDiagnostics(stdoutText)
		}
		if diagnostics != "" {
			return stdoutText, stderrText, fmt.Errorf("yt-dlp %s: %w: %s", label, err, diagnostics)
		}
		return stdoutText, stderrText, fmt.Errorf("yt-dlp %s: %w", label, err)
	}
	return stdoutText, stderrText, nil
}

func (backend *ytDLPBackend) resolveBinary() (string, error) {
	binaryPath := strings.TrimSpace(backend.binaryPath)
	if binaryPath == "" {
		binaryPath = ytdlpDefaultBinary
	}

	resolvedPath, err := backend.lookPath(binaryPath)
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) || errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("%w: install yt-dlp or configure yt_dlp_path: %w", errYTDLPUnavailable, err)
		}
		return "", fmt.Errorf("resolve yt-dlp binary %q: %w", binaryPath, err)
	}
	return resolvedPath, nil
}

func normalizedRetryAttempts(retry retryPolicy) int {
	if retry.Attempts < 1 {
		return 1
	}
	return retry.Attempts
}

func normalizedRetrySleep(retry retryPolicy) string {
	if retry.Backoff <= 0 {
		return ""
	}
	return strconv.FormatFloat(retry.Backoff.Seconds(), 'f', -1, 64)
}

func cleanYTDLPDiagnostics(text string) string {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	cleaned := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		cleaned = append(cleaned, line)
	}
	return strings.Join(cleaned, "; ")
}

func isYTDLPRateLimitedError(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	for _, marker := range []string{
		"http error 429",
		"too many requests",
		"unexpected status code: 429",
	} {
		if strings.Contains(message, marker) {
			return true
		}
	}
	return false
}

func isYTDLPNetworkBlockedError(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	for _, marker := range []string{
		"http error 400",
		"http error 403",
		"http error 500",
		"http error 502",
		"http error 503",
		"http error 504",
		"bot",
		"sign in to confirm",
		"not a bot",
	} {
		if strings.Contains(message, marker) {
			return true
		}
	}
	return false
}

func metadataFromYTDLPInfo(parsed ParsedURL, info ytDLPInfo) Metadata {
	videoID := strings.TrimSpace(info.ID)
	if videoID == "" {
		videoID = parsed.VideoID
	}
	videoURL := strings.TrimSpace(info.WebpageURL)
	if videoURL == "" {
		videoURL = parsed.CanonicalURL
	}
	channel := strings.TrimSpace(info.Channel)
	if channel == "" {
		channel = strings.TrimSpace(info.Uploader)
	}
	duration := durationFromYTDLPSeconds(info.Duration)
	return Metadata{
		VideoID:              videoID,
		URL:                  videoURL,
		Title:                strings.TrimSpace(info.Title),
		Description:          strings.TrimSpace(info.Description),
		Channel:              channel,
		ChannelID:            strings.TrimSpace(info.ChannelID),
		UploaderID:           strings.TrimSpace(info.UploaderID),
		Duration:             duration,
		DurationString:       durationStringFromYTDLP(info.DurationString, duration),
		PublishDate:          parseYTDLPUploadDate(info.UploadDate),
		ViewCount:            cloneInt64(info.ViewCount),
		LikeCount:            cloneInt64(info.LikeCount),
		CommentCount:         cloneInt64(info.CommentCount),
		ChannelFollowerCount: cloneInt64(info.ChannelFollowerCount),
		Categories:           normalizeYTDLPStringSlice(info.Categories),
		VideoTags:            normalizeYTDLPStringSlice(info.Tags),
		Language:             strings.TrimSpace(info.Language),
		LiveStatus:           strings.TrimSpace(info.LiveStatus),
		WasLive:              cloneBool(info.WasLive),
		ChapterCount:         len(info.Chapters),
	}
}

func durationFromYTDLPSeconds(seconds float64) time.Duration {
	if seconds <= 0 || math.IsNaN(seconds) || math.IsInf(seconds, 0) {
		return 0
	}
	return time.Duration(seconds * float64(time.Second))
}

func parseYTDLPUploadDate(value string) time.Time {
	parsed, err := time.Parse("20060102", strings.TrimSpace(value))
	if err != nil {
		return time.Time{}
	}
	return parsed.UTC()
}

func durationStringFromYTDLP(value string, duration time.Duration) string {
	if trimmed := strings.TrimSpace(value); trimmed != "" {
		return trimmed
	}
	if duration <= 0 {
		return ""
	}
	totalSeconds := int64(duration / time.Second)
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60
	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

func cloneInt64(value *int64) *int64 {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func cloneBool(value *bool) *bool {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func normalizeYTDLPStringSlice(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}
	normalized := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		normalized = append(normalized, trimmed)
	}
	return normalized
}

func selectYTDLPCaption(
	info ytDLPInfo,
	preferredLanguages []string,
	allowAutomatic bool,
	allowTranslated bool,
) (ytDLPCaptionCandidate, bool) {
	candidates := ytDLPCaptionCandidates(info, allowAutomatic)
	if len(candidates) == 0 {
		return ytDLPCaptionCandidate{}, false
	}
	for _, preferredLanguage := range normalizeYTDLPCaptionPreferences(preferredLanguages) {
		matches := matchingYTDLPCaptionCandidates(candidates, preferredLanguage, allowTranslated)
		if len(matches) == 0 {
			continue
		}
		sort.SliceStable(matches, func(i int, j int) bool {
			return ytDLPCaptionLess(matches[i], matches[j], preferredLanguage)
		})
		return matches[0], true
	}
	return ytDLPCaptionCandidate{}, false
}

func ytDLPCaptionCandidates(info ytDLPInfo, allowAutomatic bool) []ytDLPCaptionCandidate {
	originalLanguage := detectYTDLPOriginalLanguage(info)
	candidates := make([]ytDLPCaptionCandidate, 0, len(info.Subtitles)+len(info.AutomaticCaptions))
	for language, formats := range info.Subtitles {
		if isUsableYTDLPCaption(language, formats) {
			candidates = append(candidates, newYTDLPCaptionCandidate(language, false, originalLanguage))
		}
	}
	if !allowAutomatic {
		return candidates
	}
	for language, formats := range info.AutomaticCaptions {
		if isUsableYTDLPCaption(language, formats) {
			candidates = append(candidates, newYTDLPCaptionCandidate(language, true, originalLanguage))
		}
	}
	return candidates
}

func newYTDLPCaptionCandidate(language string, automatic bool, originalLanguage string) ytDLPCaptionCandidate {
	trimmed := strings.TrimSpace(language)
	baseLanguage := ytDLPCaptionBaseLanguage(trimmed)
	originalKey := isYTDLPOriginalCaptionKey(trimmed)
	original := originalKey
	if originalLanguage != "" && languageMatches(baseLanguage, originalLanguage) {
		original = true
	}
	if originalLanguage == "" && !automatic {
		original = true
	}
	return ytDLPCaptionCandidate{
		Language:     trimmed,
		BaseLanguage: baseLanguage,
		Automatic:    automatic,
		Original:     original,
		Translated:   automatic && !original,
		OriginalKey:  originalKey,
	}
}

func detectYTDLPOriginalLanguage(info ytDLPInfo) string {
	if language := normalizeYTDLPCaptionLanguage(info.Language); language != "" {
		return language
	}
	languages := make([]string, 0, len(info.AutomaticCaptions))
	for language, formats := range info.AutomaticCaptions {
		if !isUsableYTDLPCaption(language, formats) || !isYTDLPOriginalCaptionKey(language) {
			continue
		}
		if baseLanguage := ytDLPCaptionBaseLanguage(language); baseLanguage != "" {
			languages = append(languages, baseLanguage)
		}
	}
	sort.Strings(languages)
	if len(languages) == 0 {
		return ""
	}
	return languages[0]
}

func normalizeYTDLPCaptionPreferences(preferredLanguages []string) []string {
	preferred := normalizeLanguages(preferredLanguages)
	if len(preferred) == 0 {
		return []string{"orig"}
	}
	return preferred
}

func ytDLPCaptionUnavailableMessage(preferredLanguages []string, allowTranslated bool) string {
	preferred := normalizeYTDLPCaptionPreferences(preferredLanguages)
	if len(preferred) == 1 && preferred[0] == "orig" {
		return "no original-language caption available"
	}
	message := fmt.Sprintf("no eligible caption available for languages %s", strings.Join(preferred, ", "))
	if !allowTranslated {
		message += "; translated captions are disabled"
	}
	return message
}

func matchingYTDLPCaptionCandidates(
	candidates []ytDLPCaptionCandidate,
	preferredLanguage string,
	allowTranslated bool,
) []ytDLPCaptionCandidate {
	matches := make([]ytDLPCaptionCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		if preferredLanguage == "orig" {
			if candidate.Original {
				matches = append(matches, candidate)
			}
			continue
		}
		if !languageMatches(candidate.Language, preferredLanguage) {
			continue
		}
		if candidate.Translated && (!allowTranslated || !captionLanguageExactMatch(candidate.Language, preferredLanguage)) {
			continue
		}
		matches = append(matches, candidate)
	}
	return matches
}

func ytDLPCaptionLess(left ytDLPCaptionCandidate, right ytDLPCaptionCandidate, preferredLanguage string) bool {
	leftRank := ytDLPCaptionRank(left, preferredLanguage)
	rightRank := ytDLPCaptionRank(right, preferredLanguage)
	if leftRank != rightRank {
		return leftRank < rightRank
	}
	if left.Original && right.Original {
		leftEnglish := 1
		if languageMatches(left.BaseLanguage, "en") {
			leftEnglish = 0
		}
		rightEnglish := 1
		if languageMatches(right.BaseLanguage, "en") {
			rightEnglish = 0
		}
		if leftEnglish != rightEnglish {
			return leftEnglish < rightEnglish
		}
	}
	return left.Language < right.Language
}

func ytDLPCaptionRank(candidate ytDLPCaptionCandidate, preferredLanguage string) int {
	if preferredLanguage == "orig" {
		if !candidate.Automatic {
			return 0
		}
		switch {
		case candidate.OriginalKey:
			return 1
		case captionLanguageExactMatch(candidate.Language, candidate.BaseLanguage):
			return 2
		default:
			return 3
		}
	}

	matchRank := 2
	switch {
	case captionLanguageExactMatch(candidate.Language, preferredLanguage):
		matchRank = 0
	case candidate.BaseLanguage == preferredLanguage:
		matchRank = 1
	}
	automaticRank := 0
	if candidate.Automatic {
		automaticRank = 1
	}
	translatedRank := 0
	if candidate.Translated {
		translatedRank = 1
	}
	return matchRank*10 + automaticRank*2 + translatedRank
}

func isYTDLPOriginalCaptionKey(language string) bool {
	return strings.HasSuffix(normalizeYTDLPCaptionLanguage(language), "-orig")
}

func ytDLPCaptionBaseLanguage(language string) string {
	return strings.TrimSuffix(normalizeYTDLPCaptionLanguage(language), "-orig")
}

func normalizeYTDLPCaptionLanguage(language string) string {
	return strings.ToLower(strings.TrimSpace(language))
}

func captionLanguageExactMatch(trackLanguage string, preferredLanguage string) bool {
	return normalizeYTDLPCaptionLanguage(trackLanguage) == normalizeYTDLPCaptionLanguage(preferredLanguage)
}

func isUsableYTDLPCaption(language string, formats []ytDLPSubtitleFormat) bool {
	language = strings.TrimSpace(language)
	if language == "" || strings.EqualFold(language, "live_chat") {
		return false
	}
	for _, format := range formats {
		switch strings.ToLower(strings.TrimSpace(format.Ext)) {
		case "json3", "vtt":
			return true
		}
	}
	return false
}

func readYTDLPSubtitleFile(dir string) ([]byte, string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, "", fmt.Errorf("yt-dlp captions: read temp dir: %w", err)
	}

	var selected fs.DirEntry
	selectedExt := ""
	for _, ext := range []string{".json3", ".vtt"} {
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ext) {
				continue
			}
			selected = entry
			selectedExt = ext
			break
		}
		if selected != nil {
			break
		}
	}
	if selected == nil {
		return nil, "", fmt.Errorf("yt-dlp captions: no supported caption file was written")
	}

	path := filepath.Join(dir, selected.Name())
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", fmt.Errorf("yt-dlp captions: read %q: %w", path, err)
	}
	return data, selectedExt, nil
}

func formatYTDLPJSON3Transcript(data []byte) (string, error) {
	var payload struct {
		Events []struct {
			StartMs  int `json:"tStartMs"`
			Segments []struct {
				Text string `json:"utf8"`
			} `json:"segs"`
		} `json:"events"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return "", fmt.Errorf("yt-dlp captions: parse json3: %w", err)
	}

	transcript := make([]transcriptSegment, 0, len(payload.Events))
	for _, event := range payload.Events {
		if len(event.Segments) == 0 {
			continue
		}
		var builder strings.Builder
		for _, segment := range event.Segments {
			builder.WriteString(segment.Text)
		}
		text := normalizeTranscriptText(builder.String())
		if text == "" {
			continue
		}
		transcript = append(transcript, transcriptSegment{
			StartMs: event.StartMs,
			Text:    text,
		})
	}
	return formatTranscriptMarkdown(transcript), nil
}

func formatYTDLPVTTTranscript(data []byte) (string, error) {
	lines := strings.Split(string(data), "\n")
	transcript := make([]transcriptSegment, 0, len(lines)/3)
	for index := 0; index < len(lines); index++ {
		line := strings.TrimSpace(strings.TrimPrefix(lines[index], "\ufeff"))
		if !strings.Contains(line, "-->") {
			continue
		}
		startText, _, _ := strings.Cut(line, "-->")
		startMs, ok := parseVTTTimestampMs(startText)
		if !ok {
			continue
		}
		textLines := make([]string, 0, 2)
		for index++; index < len(lines); index++ {
			textLine := strings.TrimSpace(lines[index])
			if textLine == "" {
				break
			}
			textLines = append(textLines, textLine)
		}
		text := normalizeTranscriptText(strings.Join(textLines, " "))
		if text == "" {
			continue
		}
		transcript = append(transcript, transcriptSegment{StartMs: startMs, Text: text})
	}
	return formatTranscriptMarkdown(transcript), nil
}

func parseVTTTimestampMs(value string) (int, bool) {
	value = strings.TrimSpace(value)
	parts := strings.Split(value, ":")
	if len(parts) < 2 || len(parts) > 3 {
		return 0, false
	}

	secondsPart := parts[len(parts)-1]
	seconds, err := strconv.ParseFloat(strings.Replace(secondsPart, ",", ".", 1), 64)
	if err != nil {
		return 0, false
	}
	minutes, err := strconv.Atoi(parts[len(parts)-2])
	if err != nil {
		return 0, false
	}
	hours := 0
	if len(parts) == 3 {
		hours, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0, false
		}
	}

	return int((float64(hours*3600+minutes*60) + seconds) * 1000), true
}
