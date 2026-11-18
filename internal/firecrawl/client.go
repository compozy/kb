// Package firecrawl provides a client for the Firecrawl REST API.
package firecrawl

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/compozy/kb/internal/config"
)

const (
	defaultMaxAttempts   = 3
	defaultRetryBackoff  = 250 * time.Millisecond
	scrapeEndpointSuffix = "/v2/scrape"
	cloudHost            = "api.firecrawl.dev"
	cloudSource          = "firecrawl-cloud"
)

type sleepFunc func(context.Context, time.Duration) error

// Client calls the Firecrawl scrape API.
type Client struct {
	apiKey string
	apiURL string
	host   string

	httpClient    *http.Client
	maxAttempts   int
	retryBackoff  time.Duration
	sleepWithBack sleepFunc
}

// ScrapeResult contains the normalized Firecrawl scrape payload used by ingest.
type ScrapeResult struct {
	Markdown  string
	Title     string
	SourceURL string
	// StatusCode is the HTTP status Firecrawl reports for the scraped page
	// (metadata.statusCode); 0 when absent.
	StatusCode int
	// FinalURL is the page URL after redirects (metadata.url) when it differs
	// from SourceURL; empty otherwise.
	FinalURL string
	// SiteName is the site name the page declares (metadata.ogSiteName,
	// og:site_name or siteName); empty when absent. The ingest gate uses it
	// for the "title equals the site name" rule.
	SiteName string
	// Blocked, Source and Warning are optional Firecrawl Local fields.
	Blocked bool
	Source  string
	Warning string
	// FetchedVia identifies the backend that supplied Markdown.
	FetchedVia string
}

// ScrapeOptions are the optional freshness knobs of a scrape request (spec §7
// stage 4). A nil pointer or zero WaitFor leaves the field out of the request
// so Firecrawl applies its default.
type ScrapeOptions struct {
	// MaxAge is the maximum cache age in milliseconds; 0 bypasses the cache.
	MaxAge *int64
	// WaitFor is the delay in milliseconds before capture (client-side rendering).
	WaitFor int
	// OnlyMainContent toggles Firecrawl's main-content filter.
	OnlyMainContent *bool
}

// RefetchOptions returns the fresh-scrape options used before a quality gate
// judges a thin or broken capture: maxAge 0, waitFor refetch_wait_ms and
// onlyMainContent refetch_only_main_content.
func RefetchOptions(cfg config.FirecrawlConfig) ScrapeOptions {
	maxAge := int64(0)
	onlyMain := cfg.RefetchOnlyMainContent
	return ScrapeOptions{MaxAge: &maxAge, WaitFor: cfg.RefetchWaitMS, OnlyMainContent: &onlyMain}
}

type scrapeRequest struct {
	URL             string   `json:"url"`
	Formats         []string `json:"formats"`
	MaxAge          *int64   `json:"maxAge,omitempty"`
	WaitFor         int      `json:"waitFor,omitempty"`
	OnlyMainContent *bool    `json:"onlyMainContent,omitempty"`
}

type scrapeResponse struct {
	Success bool               `json:"success"`
	Error   string             `json:"error"`
	Message string             `json:"message"`
	Data    scrapeResponseData `json:"data"`
}

type scrapeResponseData struct {
	Markdown string                 `json:"markdown"`
	Metadata scrapeResponseMetadata `json:"metadata"`
	Blocked  bool                   `json:"blocked"`
	Source   string                 `json:"source"`
	Warning  string                 `json:"warning"`
}

type scrapeResponseMetadata struct {
	Title       string   `json:"title"`
	SourceURL   string   `json:"sourceURL"`
	URL         string   `json:"url"`
	StatusCode  int      `json:"statusCode"`
	Error       string   `json:"error"`
	OGSiteName  metaText `json:"ogSiteName"`
	OGSiteName2 metaText `json:"og:site_name"`
	SiteName    metaText `json:"siteName"`
}

// metaText is a Firecrawl metadata value that is a string or, when a page
// repeats the meta tag, an array of strings (the first non-empty one wins).
type metaText string

func (m *metaText) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*m = metaText(strings.TrimSpace(single))
		return nil
	}
	var many []string
	if err := json.Unmarshal(data, &many); err == nil {
		for _, value := range many {
			if trimmed := strings.TrimSpace(value); trimmed != "" {
				*m = metaText(trimmed)
				return nil
			}
		}
	}
	*m = ""
	return nil
}

type scrapeErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Data    struct {
		Metadata scrapeResponseMetadata `json:"metadata"`
	} `json:"data"`
}

type statusError struct {
	sourceURL  string
	statusCode int
	message    string
}

func (err *statusError) Error() string {
	if err == nil {
		return ""
	}
	if err.message != "" {
		return fmt.Sprintf("firecrawl scrape %q: request failed with status %d: %s", err.sourceURL, err.statusCode, err.message)
	}
	return fmt.Sprintf("firecrawl scrape %q: request failed with status %d", err.sourceURL, err.statusCode)
}

func (err *statusError) afterAttempts(attempts int) error {
	if err == nil {
		return nil
	}
	if err.message != "" {
		return fmt.Errorf("firecrawl scrape %q: request failed after %d attempts with status %d: %s", err.sourceURL, attempts, err.statusCode, err.message)
	}
	return fmt.Errorf("firecrawl scrape %q: request failed after %d attempts with status %d", err.sourceURL, attempts, err.statusCode)
}

// NewClient constructs a Firecrawl client from runtime configuration.
func NewClient(cfg config.FirecrawlConfig) *Client {
	defaults := config.Default().Firecrawl

	apiURL := strings.TrimSpace(cfg.APIURL)
	if apiURL == "" {
		apiURL = defaults.APIURL
	}
	var host string
	if parsed, err := url.Parse(apiURL); err == nil {
		host = strings.TrimSuffix(strings.ToLower(parsed.Hostname()), ".")
	}

	return &Client{
		apiKey:        strings.TrimSpace(cfg.APIKey),
		apiURL:        strings.TrimRight(apiURL, "/"),
		host:          host,
		httpClient:    http.DefaultClient,
		maxAttempts:   defaultMaxAttempts,
		retryBackoff:  defaultRetryBackoff,
		sleepWithBack: sleepContext,
	}
}

// NewRefetchClient returns nil when stage 4 should use the primary client.
// An error explains why a configured cloud refetch cannot be used.
func NewRefetchClient(cfg config.FirecrawlConfig) (*Client, error) {
	apiURL := strings.TrimSpace(cfg.RefetchAPIURL)
	if apiURL == "" {
		return nil, nil
	}
	apiKey := cmp.Or(strings.TrimSpace(cfg.RefetchAPIKey), strings.TrimSpace(cfg.APIKey))
	client := NewClient(config.FirecrawlConfig{APIURL: apiURL, APIKey: apiKey})
	if client.missingAPIKey() {
		return nil, errors.New("missing API key; set firecrawl.refetch_api_key or FIRECRAWL_REFETCH_API_KEY")
	}
	return client, nil
}

func (client *Client) missingAPIKey() bool {
	return client.host == cloudHost && client.apiKey == ""
}

func (client *Client) fetchedVia(source string) string {
	switch {
	case source == cloudSource || client.host == cloudHost:
		return cloudSource
	case source == "local" || (source == "" && (client.host == "localhost" || net.ParseIP(client.host).IsLoopback())):
		return "local"
	default:
		return "self-hosted"
	}
}

// Scrape converts a source URL into markdown through Firecrawl with
// Firecrawl's default freshness options.
func (client *Client) Scrape(ctx context.Context, sourceURL string) (*ScrapeResult, error) {
	return client.ScrapeWithOptions(ctx, sourceURL, ScrapeOptions{})
}

// ScrapeWithOptions is Scrape with explicit freshness options; only the
// options that are set are sent.
func (client *Client) ScrapeWithOptions(ctx context.Context, sourceURL string, opts ScrapeOptions) (*ScrapeResult, error) {
	if client == nil {
		return nil, errors.New("firecrawl scrape: client is nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	if client.missingAPIKey() {
		return nil, errors.New("firecrawl scrape: missing API key; set firecrawl.api_key or FIRECRAWL_API_KEY")
	}

	sourceURL = strings.TrimSpace(sourceURL)
	if sourceURL == "" {
		return nil, errors.New("firecrawl scrape: url is required")
	}
	if _, err := url.ParseRequestURI(sourceURL); err != nil {
		return nil, fmt.Errorf("firecrawl scrape: invalid url %q: %w", sourceURL, err)
	}

	body, err := json.Marshal(scrapeRequest{
		URL:             sourceURL,
		Formats:         []string{"markdown"},
		MaxAge:          opts.MaxAge,
		WaitFor:         max(opts.WaitFor, 0),
		OnlyMainContent: opts.OnlyMainContent,
	})
	if err != nil {
		return nil, fmt.Errorf("firecrawl scrape %q: encode request: %w", sourceURL, err)
	}

	maxAttempts := client.maxAttempts
	if maxAttempts <= 0 {
		maxAttempts = defaultMaxAttempts
	}

	retryBackoff := client.retryBackoff
	if retryBackoff <= 0 {
		retryBackoff = defaultRetryBackoff
	}

	sleepWithBack := client.sleepWithBack
	if sleepWithBack == nil {
		sleepWithBack = sleepContext
	}

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		result, retry, err := client.scrapeOnce(ctx, sourceURL, body)
		if err == nil {
			return result, nil
		}
		if !retry {
			return nil, err
		}
		if attempt == maxAttempts {
			if reqErr, ok := errors.AsType[*statusError](err); ok {
				return nil, reqErr.afterAttempts(attempt)
			}
			return nil, err
		}
		if err := sleepWithBack(ctx, backoffDuration(retryBackoff, attempt)); err != nil {
			return nil, fmt.Errorf("firecrawl scrape %q: retry canceled: %w", sourceURL, err)
		}
	}

	return nil, fmt.Errorf("firecrawl scrape %q: exhausted retries", sourceURL)
}

func (client *Client) scrapeOnce(ctx context.Context, sourceURL string, body []byte) (*ScrapeResult, bool, error) {
	httpClient := client.httpClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, client.endpointURL(), bytes.NewReader(body))
	if err != nil {
		return nil, false, fmt.Errorf("firecrawl scrape %q: build request: %w", sourceURL, err)
	}
	if client.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+client.apiKey)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, false, fmt.Errorf("firecrawl scrape %q: request canceled: %w", sourceURL, ctxErr)
		}
		return nil, false, fmt.Errorf("firecrawl scrape %q: request failed: %w", sourceURL, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false, fmt.Errorf("firecrawl scrape %q: read response: %w", sourceURL, err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		reqErr := &statusError{
			sourceURL:  sourceURL,
			statusCode: resp.StatusCode,
			message:    parseErrorMessage(responseBody),
		}
		return nil, isRetryableStatus(resp.StatusCode), reqErr
	}

	var payload scrapeResponse
	if err := json.Unmarshal(responseBody, &payload); err != nil {
		return nil, false, fmt.Errorf("firecrawl scrape %q: parse response: %w", sourceURL, err)
	}

	if apiError := firstNonEmpty(payload.Error, payload.Message, payload.Data.Metadata.Error); !payload.Success && apiError != "" {
		return nil, false, fmt.Errorf("firecrawl scrape %q: api error: %s", sourceURL, apiError)
	}

	meta := payload.Data.Metadata
	resolvedSource := firstNonEmpty(strings.TrimSpace(meta.SourceURL), strings.TrimSpace(meta.URL), sourceURL)
	finalURL := strings.TrimSpace(meta.URL)
	if finalURL == resolvedSource {
		finalURL = ""
	}
	source := strings.TrimSpace(payload.Data.Source)
	return &ScrapeResult{
		Markdown:   payload.Data.Markdown,
		Title:      strings.TrimSpace(meta.Title),
		SourceURL:  resolvedSource,
		StatusCode: meta.StatusCode,
		FinalURL:   finalURL,
		SiteName:   firstNonEmpty(string(meta.OGSiteName), string(meta.OGSiteName2), string(meta.SiteName)),
		Blocked:    payload.Data.Blocked,
		Source:     source,
		Warning:    strings.TrimSpace(payload.Data.Warning),
		FetchedVia: client.fetchedVia(source),
	}, false, nil
}

func (client *Client) endpointURL() string {
	return client.apiURL + scrapeEndpointSuffix
}

func backoffDuration(base time.Duration, attempt int) time.Duration {
	if base <= 0 {
		base = defaultRetryBackoff
	}
	if attempt <= 0 {
		return base
	}
	return base * time.Duration(1<<(attempt-1))
}

func sleepContext(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func isRetryableStatus(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests || statusCode >= http.StatusInternalServerError
}

func parseErrorMessage(body []byte) string {
	if len(body) == 0 {
		return ""
	}

	var payload scrapeErrorResponse
	if err := json.Unmarshal(body, &payload); err == nil {
		if message := firstNonEmpty(payload.Error, payload.Message, payload.Data.Metadata.Error); message != "" {
			return message
		}
	}

	return strings.TrimSpace(string(body))
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
