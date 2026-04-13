// Package firecrawl provides a client for the Firecrawl REST API.
package firecrawl

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
)

type sleepFunc func(context.Context, time.Duration) error

// Client calls the Firecrawl scrape API.
type Client struct {
	apiKey string
	apiURL string

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
}

type scrapeRequest struct {
	URL     string   `json:"url"`
	Formats []string `json:"formats"`
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
}

type scrapeResponseMetadata struct {
	Title     string `json:"title"`
	SourceURL string `json:"sourceURL"`
	URL       string `json:"url"`
	Error     string `json:"error"`
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

	return &Client{
		apiKey:        strings.TrimSpace(cfg.APIKey),
		apiURL:        strings.TrimRight(apiURL, "/"),
		httpClient:    http.DefaultClient,
		maxAttempts:   defaultMaxAttempts,
		retryBackoff:  defaultRetryBackoff,
		sleepWithBack: sleepContext,
	}
}

// Scrape converts a source URL into markdown through Firecrawl.
func (client *Client) Scrape(ctx context.Context, sourceURL string) (*ScrapeResult, error) {
	if client == nil {
		return nil, errors.New("firecrawl scrape: client is nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	if client.apiKey == "" {
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
		URL:     sourceURL,
		Formats: []string{"markdown"},
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
			var reqErr *statusError
			if errors.As(err, &reqErr) {
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
	req.Header.Set("Authorization", "Bearer "+client.apiKey)
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

	return &ScrapeResult{
		Markdown:  payload.Data.Markdown,
		Title:     strings.TrimSpace(payload.Data.Metadata.Title),
		SourceURL: firstNonEmpty(strings.TrimSpace(payload.Data.Metadata.SourceURL), strings.TrimSpace(payload.Data.Metadata.URL), sourceURL),
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
