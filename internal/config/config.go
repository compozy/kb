// Package config loads TOML configuration and environment-backed runtime overrides for the kb CLI.
package config

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

const (
	defaultFirecrawlAPIURL        = "https://api.firecrawl.dev"
	defaultOpenAIAPIURL           = "https://api.openai.com"
	defaultOpenRouterAPIURL       = "https://openrouter.ai/api"
	defaultOpenRouterSTTModel     = "google/gemini-2.5-flash"
	defaultSTTProvider            = "openai"
	defaultSTTModel               = "gpt-4o-transcribe"
	defaultSTTLanguage            = "auto"
	defaultSTTAudioFormat         = "mp3"
	defaultSTTChunkDuration       = "10m"
	defaultSTTMaxChunkBytes       = 24_000_000
	defaultSTTConcurrency         = 2
	defaultSTTFFmpegPath          = "ffmpeg"
	defaultVaultRoot              = "."
	defaultTopicGlob              = "*"
	defaultYouTubeYTDLPPath       = "yt-dlp"
	defaultYouTubeTranscription   = "captions"
	defaultYouTubeRetryBackoff    = "1s"
	defaultYouTubeRetryAttempts   = 3
	defaultYouTubeCaptionLanguage = "orig"
	defaultYouTubeBulkConcurrency = 2
	defaultYouTubeBulkThrottle    = "2s"
	defaultYouTubeBulkBackoffMax  = "60s"
	defaultYouTubeBulkRetries     = 3
	defaultInstagramYTDLPPath     = "yt-dlp"
	defaultInstagramTranscription = "auto"
	defaultInstagramRetryBackoff  = "1s"
	defaultInstagramRetryAttempts = 3
)

// Config contains the complete TOML-backed runtime configuration.
type Config struct {
	App        AppConfig        `toml:"app"`
	Log        LogConfig        `toml:"log"`
	Vault      VaultConfig      `toml:"vault"`
	OKF        OKFConfig        `toml:"okf"`
	Firecrawl  FirecrawlConfig  `toml:"firecrawl"`
	OpenRouter OpenRouterConfig `toml:"openrouter"`
	STT        STTConfig        `toml:"stt"`
	YouTube    YouTubeConfig    `toml:"youtube"`
	Instagram  InstagramConfig  `toml:"instagram"`
}

// AppConfig contains the application identity and environment.
type AppConfig struct {
	Name string `toml:"name"`
	Env  string `toml:"env"`
}

// LogConfig controls structured logging output.
type LogConfig struct {
	Level string `toml:"level"`
}

// VaultConfig controls vault root discovery and topic enumeration.
type VaultConfig struct {
	Root       string   `toml:"root"`
	TopicGlobs []string `toml:"topic_globs"`
}

// OKFConfig controls local Open Knowledge Format producer standards.
type OKFConfig struct {
	Types []string `toml:"types"`
}

// FirecrawlConfig controls URL scraping API access.
type FirecrawlConfig struct {
	APIKey string `toml:"api_key"`
	APIURL string `toml:"api_url"`
}

// OpenRouterConfig controls the optional OpenRouter STT provider.
type OpenRouterConfig struct {
	APIKey   string `toml:"api_key"`
	APIURL   string `toml:"api_url"`
	STTModel string `toml:"stt_model"`
}

// STTConfig controls speech-to-text provider selection and audio chunking.
type STTConfig struct {
	Provider      string `toml:"provider"`
	APIKey        string `toml:"api_key"`
	APIURL        string `toml:"api_url"`
	Model         string `toml:"model"`
	Language      string `toml:"language"`
	Prompt        string `toml:"prompt"`
	AudioFormat   string `toml:"audio_format"`
	ChunkDuration string `toml:"chunk_duration"`
	MaxChunkBytes int64  `toml:"max_chunk_bytes"`
	Concurrency   int    `toml:"concurrency"`
	FFmpegPath    string `toml:"ffmpeg_path"`
}

// YouTubeConfig controls YouTube network access and retry behavior.
type YouTubeConfig struct {
	YTDLPPath     string `toml:"yt_dlp_path"`
	Proxy         string `toml:"proxy"`
	CookiesFile   string `toml:"cookies_file"`
	UserAgent     string `toml:"user_agent"`
	Transcription string `toml:"transcription"`
	RetryAttempts int    `toml:"retry_attempts"`
	RetryBackoff  string `toml:"retry_backoff"`
	// CaptionLanguages controls caption track selection. "orig" means the
	// video's original language as reported by yt-dlp metadata.
	CaptionLanguages        []string `toml:"caption_languages"`
	AllowTranslatedCaptions bool     `toml:"allow_translated_captions"`
	// Bulk* control the kb ingest channel worker pool, inter-request throttle,
	// and adaptive backoff used when ingesting an entire channel or playlist.
	BulkConcurrency int    `toml:"bulk_concurrency"`
	BulkThrottle    string `toml:"bulk_throttle"`
	BulkBackoffMax  string `toml:"bulk_backoff_max"`
	BulkRetries     int    `toml:"bulk_retries"`
}

// InstagramConfig controls Instagram network access and retry behavior for
// `kb ingest instagram`. Cookies are kept separate from YouTube because the two
// platforms require distinct authenticated sessions.
type InstagramConfig struct {
	YTDLPPath     string `toml:"yt_dlp_path"`
	Proxy         string `toml:"proxy"`
	CookiesFile   string `toml:"cookies_file"`
	UserAgent     string `toml:"user_agent"`
	Transcription string `toml:"transcription"`
	RetryAttempts int    `toml:"retry_attempts"`
	RetryBackoff  string `toml:"retry_backoff"`
}

// Default returns a sane starting configuration.
func Default() Config {
	return Config{
		App: AppConfig{
			Name: "app",
			Env:  "development",
		},
		Log: LogConfig{
			Level: "info",
		},
		Vault: VaultConfig{
			Root:       defaultVaultRoot,
			TopicGlobs: []string{defaultTopicGlob},
		},
		Firecrawl: FirecrawlConfig{
			APIURL: defaultFirecrawlAPIURL,
		},
		OpenRouter: OpenRouterConfig{
			APIURL:   defaultOpenRouterAPIURL,
			STTModel: defaultOpenRouterSTTModel,
		},
		STT: STTConfig{
			Provider:      defaultSTTProvider,
			APIURL:        defaultOpenAIAPIURL,
			Model:         defaultSTTModel,
			Language:      defaultSTTLanguage,
			AudioFormat:   defaultSTTAudioFormat,
			ChunkDuration: defaultSTTChunkDuration,
			MaxChunkBytes: defaultSTTMaxChunkBytes,
			Concurrency:   defaultSTTConcurrency,
			FFmpegPath:    defaultSTTFFmpegPath,
		},
		YouTube: YouTubeConfig{
			YTDLPPath:     defaultYouTubeYTDLPPath,
			Transcription: defaultYouTubeTranscription,
			RetryAttempts: defaultYouTubeRetryAttempts,
			RetryBackoff:  defaultYouTubeRetryBackoff,
			CaptionLanguages: []string{
				defaultYouTubeCaptionLanguage,
			},
			BulkConcurrency: defaultYouTubeBulkConcurrency,
			BulkThrottle:    defaultYouTubeBulkThrottle,
			BulkBackoffMax:  defaultYouTubeBulkBackoffMax,
			BulkRetries:     defaultYouTubeBulkRetries,
		},
		Instagram: InstagramConfig{
			YTDLPPath:     defaultInstagramYTDLPPath,
			Transcription: defaultInstagramTranscription,
			RetryAttempts: defaultInstagramRetryAttempts,
			RetryBackoff:  defaultInstagramRetryBackoff,
		},
	}
}

// Load reads and validates the TOML config file, then overlays runtime secrets
// from the environment.
func Load(path string) (Config, error) {
	cfg := Default()
	if path != "" {
		if err := decodeFile(path, &cfg); err != nil {
			return Config{}, err
		}
	}
	cfg.applyDefaults()
	ApplyEnvOverrides(&cfg)
	return cfg, cfg.Validate()
}

func (c *Config) applyDefaults() {
	if c == nil {
		return
	}
	if strings.TrimSpace(c.Vault.Root) == "" {
		c.Vault.Root = defaultVaultRoot
	}
	if len(c.Vault.TopicGlobs) == 0 {
		c.Vault.TopicGlobs = []string{defaultTopicGlob}
	}
	c.OKF.Types = normalizeConfigStringList(c.OKF.Types)
	if strings.TrimSpace(c.YouTube.YTDLPPath) == "" {
		c.YouTube.YTDLPPath = defaultYouTubeYTDLPPath
	}
	if strings.TrimSpace(c.YouTube.Transcription) == "" {
		c.YouTube.Transcription = defaultYouTubeTranscription
	}
	c.YouTube.CaptionLanguages = normalizeConfigStringList(c.YouTube.CaptionLanguages)
	if len(c.YouTube.CaptionLanguages) == 0 {
		c.YouTube.CaptionLanguages = []string{defaultYouTubeCaptionLanguage}
	}
	if strings.TrimSpace(c.STT.Provider) == "" {
		c.STT.Provider = defaultSTTProvider
	}
	if strings.TrimSpace(c.STT.APIURL) == "" {
		c.STT.APIURL = defaultOpenAIAPIURL
	}
	if strings.TrimSpace(c.STT.Model) == "" {
		c.STT.Model = defaultSTTModel
	}
	if strings.TrimSpace(c.STT.Language) == "" {
		c.STT.Language = defaultSTTLanguage
	}
	if strings.TrimSpace(c.STT.AudioFormat) == "" {
		c.STT.AudioFormat = defaultSTTAudioFormat
	}
	if strings.TrimSpace(c.STT.ChunkDuration) == "" {
		c.STT.ChunkDuration = defaultSTTChunkDuration
	}
	if c.STT.MaxChunkBytes == 0 {
		c.STT.MaxChunkBytes = defaultSTTMaxChunkBytes
	}
	if c.STT.Concurrency == 0 {
		c.STT.Concurrency = defaultSTTConcurrency
	}
	if strings.TrimSpace(c.STT.FFmpegPath) == "" {
		c.STT.FFmpegPath = defaultSTTFFmpegPath
	}
	if c.YouTube.RetryAttempts == 0 {
		c.YouTube.RetryAttempts = defaultYouTubeRetryAttempts
	}
	if strings.TrimSpace(c.YouTube.RetryBackoff) == "" {
		c.YouTube.RetryBackoff = defaultYouTubeRetryBackoff
	}
	if c.YouTube.BulkConcurrency == 0 {
		c.YouTube.BulkConcurrency = defaultYouTubeBulkConcurrency
	}
	if strings.TrimSpace(c.YouTube.BulkThrottle) == "" {
		c.YouTube.BulkThrottle = defaultYouTubeBulkThrottle
	}
	if strings.TrimSpace(c.YouTube.BulkBackoffMax) == "" {
		c.YouTube.BulkBackoffMax = defaultYouTubeBulkBackoffMax
	}
	if c.YouTube.BulkRetries == 0 {
		c.YouTube.BulkRetries = defaultYouTubeBulkRetries
	}
	if strings.TrimSpace(c.Instagram.YTDLPPath) == "" {
		c.Instagram.YTDLPPath = defaultInstagramYTDLPPath
	}
	if strings.TrimSpace(c.Instagram.Transcription) == "" {
		c.Instagram.Transcription = defaultInstagramTranscription
	}
	if c.Instagram.RetryAttempts == 0 {
		c.Instagram.RetryAttempts = defaultInstagramRetryAttempts
	}
	if strings.TrimSpace(c.Instagram.RetryBackoff) == "" {
		c.Instagram.RetryBackoff = defaultInstagramRetryBackoff
	}
}

func decodeFile(path string, cfg *Config) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("stat config %q: %w", path, err)
	}

	meta, err := toml.DecodeFile(path, cfg)
	if err != nil {
		return fmt.Errorf("decode config %q: %w", path, err)
	}

	if undecoded := meta.Undecoded(); len(undecoded) > 0 {
		keys := make([]string, 0, len(undecoded))
		for _, key := range undecoded {
			keys = append(keys, key.String())
		}
		sort.Strings(keys)
		return fmt.Errorf("unknown config keys: %s", strings.Join(keys, ", "))
	}

	return nil
}

// Validate ensures the config is internally consistent before runtime startup.
func (c Config) Validate() error {
	if err := c.App.Validate(); err != nil {
		return err
	}
	if err := c.Log.Validate(); err != nil {
		return err
	}
	if err := c.Vault.Validate(); err != nil {
		return err
	}
	if err := c.STT.Validate(); err != nil {
		return err
	}
	if err := c.YouTube.Validate(); err != nil {
		return err
	}
	if err := c.Instagram.Validate(); err != nil {
		return err
	}
	return nil
}

// Validate ensures application identity settings are usable.
func (c AppConfig) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("app.name is required")
	}
	switch strings.ToLower(strings.TrimSpace(c.Env)) {
	case "development", "staging", "production":
	default:
		return fmt.Errorf("app.env must be development, staging, or production: %q", c.Env)
	}
	return nil
}

// Validate ensures the log level is supported.
func (c LogConfig) Validate() error {
	switch strings.ToLower(strings.TrimSpace(c.Level)) {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("log.level must be debug, info, warn, or error: %q", c.Level)
	}
	return nil
}

// Validate ensures vault discovery settings are usable.
func (c VaultConfig) Validate() error {
	if strings.TrimSpace(c.Root) == "" {
		return errors.New("vault.root is required")
	}
	if len(c.TopicGlobs) == 0 {
		return errors.New("vault.topic_globs must contain at least one glob")
	}
	for _, glob := range c.TopicGlobs {
		cleanGlob := strings.TrimSpace(glob)
		switch {
		case cleanGlob == "":
			return errors.New("vault.topic_globs cannot contain an empty glob")
		case strings.HasPrefix(cleanGlob, "/"):
			return fmt.Errorf("vault.topic_globs must be relative: %q", glob)
		case cleanGlob == "." || cleanGlob == ".." || strings.Contains(cleanGlob, "../") || strings.Contains(cleanGlob, "..\\"):
			return fmt.Errorf("vault.topic_globs cannot escape the vault root: %q", glob)
		}
	}
	return nil
}

// Validate ensures STT provider and chunking settings are usable.
func (c STTConfig) Validate() error {
	switch strings.ToLower(strings.TrimSpace(c.Provider)) {
	case "openai", "openrouter":
	default:
		return fmt.Errorf("stt.provider must be openai or openrouter: %q", c.Provider)
	}
	if strings.TrimSpace(c.APIURL) == "" {
		return errors.New("stt.api_url is required")
	}
	if strings.TrimSpace(c.Model) == "" {
		return errors.New("stt.model is required")
	}
	if strings.TrimSpace(c.Language) == "" {
		return errors.New("stt.language is required")
	}
	switch strings.ToLower(strings.TrimSpace(c.AudioFormat)) {
	case "mp3", "mp4", "mpeg", "mpga", "m4a", "wav", "webm":
	default:
		return fmt.Errorf("stt.audio_format is unsupported: %q", c.AudioFormat)
	}
	if _, err := c.ChunkDurationValue(); err != nil {
		return err
	}
	if c.MaxChunkBytes < 1 {
		return fmt.Errorf("stt.max_chunk_bytes must be at least 1: %d", c.MaxChunkBytes)
	}
	if c.Concurrency < 1 {
		return fmt.Errorf("stt.concurrency must be at least 1: %d", c.Concurrency)
	}
	if strings.TrimSpace(c.FFmpegPath) == "" {
		return errors.New("stt.ffmpeg_path is required")
	}
	return nil
}

// ChunkDurationValue parses the configured STT chunk duration.
func (c STTConfig) ChunkDurationValue() (time.Duration, error) {
	value := strings.TrimSpace(c.ChunkDuration)
	if value == "" {
		value = defaultSTTChunkDuration
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("stt.chunk_duration must be a Go duration: %q", c.ChunkDuration)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("stt.chunk_duration must be positive: %q", c.ChunkDuration)
	}
	return duration, nil
}

// Validate ensures YouTube runtime settings are usable.
func (c YouTubeConfig) Validate() error {
	if strings.TrimSpace(c.YTDLPPath) == "" {
		return errors.New("youtube.yt_dlp_path is required")
	}
	switch strings.ToLower(strings.TrimSpace(c.Transcription)) {
	case "captions", "auto", "stt":
	default:
		return fmt.Errorf("youtube.transcription must be captions, auto, or stt: %q", c.Transcription)
	}
	if c.RetryAttempts < 1 {
		return fmt.Errorf("youtube.retry_attempts must be at least 1: %d", c.RetryAttempts)
	}
	if _, err := c.RetryBackoffDuration(); err != nil {
		return err
	}
	if len(normalizeConfigStringList(c.CaptionLanguages)) == 0 {
		return errors.New("youtube.caption_languages must contain at least one language or orig")
	}
	if c.BulkConcurrency < 1 {
		return fmt.Errorf("youtube.bulk_concurrency must be at least 1: %d", c.BulkConcurrency)
	}
	if c.BulkRetries < 1 {
		return fmt.Errorf("youtube.bulk_retries must be at least 1: %d", c.BulkRetries)
	}
	if _, err := c.BulkThrottleDuration(); err != nil {
		return err
	}
	if _, err := c.BulkBackoffMaxDuration(); err != nil {
		return err
	}
	return nil
}

// RetryBackoffDuration parses the configured retry backoff.
func (c YouTubeConfig) RetryBackoffDuration() (time.Duration, error) {
	backoff := strings.TrimSpace(c.RetryBackoff)
	if backoff == "" {
		backoff = defaultYouTubeRetryBackoff
	}
	duration, err := time.ParseDuration(backoff)
	if err != nil {
		return 0, fmt.Errorf("youtube.retry_backoff must be a Go duration: %q", c.RetryBackoff)
	}
	if duration < 0 {
		return 0, fmt.Errorf("youtube.retry_backoff cannot be negative: %q", c.RetryBackoff)
	}
	return duration, nil
}

func normalizeConfigStringList(values []string) []string {
	normalized := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		normalized = append(normalized, value)
	}
	return normalized
}

// BulkThrottleDuration parses the configured inter-request throttle for bulk
// channel ingestion.
func (c YouTubeConfig) BulkThrottleDuration() (time.Duration, error) {
	value := strings.TrimSpace(c.BulkThrottle)
	if value == "" {
		value = defaultYouTubeBulkThrottle
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("youtube.bulk_throttle must be a Go duration: %q", c.BulkThrottle)
	}
	if duration < 0 {
		return 0, fmt.Errorf("youtube.bulk_throttle cannot be negative: %q", c.BulkThrottle)
	}
	return duration, nil
}

// BulkBackoffMaxDuration parses the configured adaptive backoff ceiling for bulk
// channel ingestion.
func (c YouTubeConfig) BulkBackoffMaxDuration() (time.Duration, error) {
	value := strings.TrimSpace(c.BulkBackoffMax)
	if value == "" {
		value = defaultYouTubeBulkBackoffMax
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("youtube.bulk_backoff_max must be a Go duration: %q", c.BulkBackoffMax)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("youtube.bulk_backoff_max must be positive: %q", c.BulkBackoffMax)
	}
	return duration, nil
}

// Validate ensures Instagram runtime settings are usable.
func (c InstagramConfig) Validate() error {
	if strings.TrimSpace(c.YTDLPPath) == "" {
		return errors.New("instagram.yt_dlp_path is required")
	}
	switch strings.ToLower(strings.TrimSpace(c.Transcription)) {
	case "captions", "auto", "stt":
	default:
		return fmt.Errorf("instagram.transcription must be captions, auto, or stt: %q", c.Transcription)
	}
	if c.RetryAttempts < 1 {
		return fmt.Errorf("instagram.retry_attempts must be at least 1: %d", c.RetryAttempts)
	}
	if _, err := c.RetryBackoffDuration(); err != nil {
		return err
	}
	return nil
}

// RetryBackoffDuration parses the configured Instagram retry backoff.
func (c InstagramConfig) RetryBackoffDuration() (time.Duration, error) {
	backoff := strings.TrimSpace(c.RetryBackoff)
	if backoff == "" {
		backoff = defaultInstagramRetryBackoff
	}
	duration, err := time.ParseDuration(backoff)
	if err != nil {
		return 0, fmt.Errorf("instagram.retry_backoff must be a Go duration: %q", c.RetryBackoff)
	}
	if duration < 0 {
		return 0, fmt.Errorf("instagram.retry_backoff cannot be negative: %q", c.RetryBackoff)
	}
	return duration, nil
}
