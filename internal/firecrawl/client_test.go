package firecrawl

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/go-devstack/internal/config"
)

func TestScrapeSendsExpectedRequestAndReturnsContent(t *testing.T) {
	t.Parallel()

	sourceURL := "https://example.com/article"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %q, want %q", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/v2/scrape" {
			t.Fatalf("path = %q, want %q", r.URL.Path, "/v2/scrape")
		}
		if got := r.Header.Get("Authorization"); got != "Bearer firecrawl-key" {
			t.Fatalf("authorization = %q, want %q", got, "Bearer firecrawl-key")
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("content-type = %q, want %q", got, "application/json")
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request body: %v", err)
		}

		expected := map[string]any{
			"url": sourceURL,
			"formats": []any{
				"markdown",
			},
		}
		if !reflect.DeepEqual(body, expected) {
			t.Fatalf("request body = %#v, want %#v", body, expected)
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{
			"success": true,
			"data": {
				"markdown": "# Example",
				"metadata": {
					"title": "Example title",
					"sourceURL": "https://example.com/article"
				}
			}
		}`)); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient(config.FirecrawlConfig{
		APIKey: "firecrawl-key",
		APIURL: server.URL,
	})

	result, err := client.Scrape(context.Background(), sourceURL)
	if err != nil {
		t.Fatalf("Scrape returned error: %v", err)
	}
	if result.Markdown != "# Example" {
		t.Fatalf("markdown = %q, want %q", result.Markdown, "# Example")
	}
	if result.Title != "Example title" {
		t.Fatalf("title = %q, want %q", result.Title, "Example title")
	}
	if result.SourceURL != sourceURL {
		t.Fatalf("source url = %q, want %q", result.SourceURL, sourceURL)
	}
}

func TestScrapeRetriesOnTooManyRequests(t *testing.T) {
	t.Parallel()

	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt := attempts.Add(1)
		if attempt < 3 {
			writeJSONStatus(t, w, http.StatusTooManyRequests, `{"error":"rate limited"}`)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{
			"success": true,
			"data": {
				"markdown": "# Recovered",
				"metadata": {
					"title": "Recovered",
					"sourceURL": "https://example.com/recovered"
				}
			}
		}`)); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient(config.FirecrawlConfig{
		APIKey: "firecrawl-key",
		APIURL: server.URL,
	})
	client.retryBackoff = time.Millisecond

	var sleeps []time.Duration
	client.sleepWithBack = func(ctx context.Context, duration time.Duration) error {
		sleeps = append(sleeps, duration)
		return nil
	}

	result, err := client.Scrape(context.Background(), "https://example.com/recovered")
	if err != nil {
		t.Fatalf("Scrape returned error: %v", err)
	}
	if attempts.Load() != 3 {
		t.Fatalf("attempts = %d, want 3", attempts.Load())
	}
	if !reflect.DeepEqual(sleeps, []time.Duration{time.Millisecond, 2 * time.Millisecond}) {
		t.Fatalf("sleeps = %#v, want exponential backoff", sleeps)
	}
	if result.Title != "Recovered" {
		t.Fatalf("title = %q, want %q", result.Title, "Recovered")
	}
}

func TestScrapeRetriesOnServerErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		statusCode int
	}{
		{name: "500", statusCode: http.StatusInternalServerError},
		{name: "502", statusCode: http.StatusBadGateway},
		{name: "503", statusCode: http.StatusServiceUnavailable},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var attempts atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				attempt := attempts.Add(1)
				if attempt < 3 {
					writeJSONStatus(t, w, tc.statusCode, `{"error":"temporary upstream failure"}`)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				if _, err := w.Write([]byte(`{
					"success": true,
					"data": {
						"markdown": "# Success",
						"metadata": {
							"title": "Success",
							"sourceURL": "https://example.com/success"
						}
					}
				}`)); err != nil {
					t.Fatalf("write response: %v", err)
				}
			}))
			defer server.Close()

			client := NewClient(config.FirecrawlConfig{
				APIKey: "firecrawl-key",
				APIURL: server.URL,
			})
			client.retryBackoff = time.Millisecond
			client.sleepWithBack = func(ctx context.Context, duration time.Duration) error { return nil }

			if _, err := client.Scrape(context.Background(), "https://example.com/success"); err != nil {
				t.Fatalf("Scrape returned error: %v", err)
			}
			if attempts.Load() != 3 {
				t.Fatalf("attempts = %d, want 3", attempts.Load())
			}
		})
	}
}

func TestScrapeDoesNotRetryOnClientErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		statusCode int
	}{
		{name: "400", statusCode: http.StatusBadRequest},
		{name: "401", statusCode: http.StatusUnauthorized},
		{name: "404", statusCode: http.StatusNotFound},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var attempts atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				attempts.Add(1)
				writeJSONStatus(t, w, tc.statusCode, `{"error":"bad request"}`)
			}))
			defer server.Close()

			client := NewClient(config.FirecrawlConfig{
				APIKey: "firecrawl-key",
				APIURL: server.URL,
			})
			client.sleepWithBack = func(ctx context.Context, duration time.Duration) error {
				t.Fatalf("unexpected retry sleep for status %d", tc.statusCode)
				return nil
			}

			_, err := client.Scrape(context.Background(), "https://example.com/bad")
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if attempts.Load() != 1 {
				t.Fatalf("attempts = %d, want 1", attempts.Load())
			}
			if !strings.Contains(err.Error(), "status "+strconv.Itoa(tc.statusCode)) {
				t.Fatalf("error = %q, want status %d", err.Error(), tc.statusCode)
			}
		})
	}
}

func TestScrapeReturnsDescriptiveErrorWhenAPIKeyIsMissing(t *testing.T) {
	t.Parallel()

	client := NewClient(config.FirecrawlConfig{})

	_, err := client.Scrape(context.Background(), "https://example.com/article")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "missing API key") {
		t.Fatalf("error = %q, want missing API key message", err.Error())
	}
	if !strings.Contains(err.Error(), "FIRECRAWL_API_KEY") {
		t.Fatalf("error = %q, want env guidance", err.Error())
	}
}

func TestScrapeRejectsInvalidURL(t *testing.T) {
	t.Parallel()

	client := NewClient(config.FirecrawlConfig{
		APIKey: "firecrawl-key",
	})

	_, err := client.Scrape(context.Background(), "://bad-url")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid url") {
		t.Fatalf("error = %q, want invalid url message", err.Error())
	}
}

func TestScrapeReturnsAPIErrorFromSuccessfulHTTPResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"success": false, "error": "page blocked by upstream"}`)); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient(config.FirecrawlConfig{
		APIKey: "firecrawl-key",
		APIURL: server.URL,
	})
	client.sleepWithBack = func(ctx context.Context, duration time.Duration) error {
		t.Fatal("unexpected retry for API error response")
		return nil
	}

	_, err := client.Scrape(context.Background(), "https://example.com/blocked")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "page blocked by upstream") {
		t.Fatalf("error = %q, want API error detail", err.Error())
	}
}

func TestScrapeReturnsDescriptiveErrorWhenURLIsUnreachable(t *testing.T) {
	t.Parallel()

	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		writeJSONStatus(t, w, http.StatusServiceUnavailable, `{"error":"source URL is unreachable"}`)
	}))
	defer server.Close()

	client := NewClient(config.FirecrawlConfig{
		APIKey: "firecrawl-key",
		APIURL: server.URL,
	})
	client.retryBackoff = time.Millisecond
	client.sleepWithBack = func(ctx context.Context, duration time.Duration) error { return nil }

	_, err := client.Scrape(context.Background(), "https://example.com/unreachable")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if attempts.Load() != 3 {
		t.Fatalf("attempts = %d, want 3", attempts.Load())
	}
	if !strings.Contains(err.Error(), "after 3 attempts") {
		t.Fatalf("error = %q, want retry count", err.Error())
	}
	if !strings.Contains(err.Error(), "status 503") {
		t.Fatalf("error = %q, want status code", err.Error())
	}
	if !strings.Contains(err.Error(), "source URL is unreachable") {
		t.Fatalf("error = %q, want upstream failure detail", err.Error())
	}
}

func TestScrapeReturnsParseErrorForInvalidJSONResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte("not-json")); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient(config.FirecrawlConfig{
		APIKey: "firecrawl-key",
		APIURL: server.URL,
	})

	_, err := client.Scrape(context.Background(), "https://example.com/invalid-json")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "parse response") {
		t.Fatalf("error = %q, want parse response detail", err.Error())
	}
}

func TestScrapeRespectsContextCancellationDuringRetry(t *testing.T) {
	t.Parallel()

	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		writeJSONStatus(t, w, http.StatusTooManyRequests, `{"error":"rate limited"}`)
	}))
	defer server.Close()

	client := NewClient(config.FirecrawlConfig{
		APIKey: "firecrawl-key",
		APIURL: server.URL,
	})

	ctx, cancel := context.WithCancel(context.Background())
	client.sleepWithBack = func(ctx context.Context, duration time.Duration) error {
		cancel()
		<-ctx.Done()
		return ctx.Err()
	}

	_, err := client.Scrape(ctx, "https://example.com/retry")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if attempts.Load() != 1 {
		t.Fatalf("attempts = %d, want 1", attempts.Load())
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context canceled", err)
	}
	if !strings.Contains(err.Error(), "retry canceled") {
		t.Fatalf("error = %q, want retry canceled context", err.Error())
	}
}

func TestScrapeReturnsRequestFailure(t *testing.T) {
	t.Parallel()

	client := NewClient(config.FirecrawlConfig{
		APIKey: "firecrawl-key",
		APIURL: "https://firecrawl.internal",
	})
	client.httpClient = &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New("dial failed")
		}),
	}

	_, err := client.Scrape(context.Background(), "https://example.com/request-failure")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "request failed") {
		t.Fatalf("error = %q, want request failed detail", err.Error())
	}
	if !strings.Contains(err.Error(), "dial failed") {
		t.Fatalf("error = %q, want transport detail", err.Error())
	}
}

func TestScrapeUsesCustomAPIURL(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/self-hosted/v2/scrape" {
			t.Fatalf("path = %q, want %q", r.URL.Path, "/self-hosted/v2/scrape")
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{
			"success": true,
			"data": {
				"markdown": "# Self hosted",
				"metadata": {
					"title": "Self hosted"
				}
			}
		}`)); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient(config.FirecrawlConfig{
		APIKey: "firecrawl-key",
		APIURL: server.URL + "/self-hosted/",
	})

	result, err := client.Scrape(context.Background(), "https://example.com/self-hosted")
	if err != nil {
		t.Fatalf("Scrape returned error: %v", err)
	}
	if result.Markdown != "# Self hosted" {
		t.Fatalf("markdown = %q, want %q", result.Markdown, "# Self hosted")
	}
	if result.SourceURL != "https://example.com/self-hosted" {
		t.Fatalf("source url = %q, want %q", result.SourceURL, "https://example.com/self-hosted")
	}
}

func TestStatusErrorFormatting(t *testing.T) {
	t.Parallel()

	var nilErr *statusError
	if got := nilErr.Error(); got != "" {
		t.Fatalf("nil error string = %q, want empty string", got)
	}
	if err := nilErr.afterAttempts(3); err != nil {
		t.Fatalf("nil error afterAttempts = %v, want nil", err)
	}

	reqErr := &statusError{
		sourceURL:  "https://example.com",
		statusCode: http.StatusServiceUnavailable,
		message:    "upstream unavailable",
	}

	if got := reqErr.Error(); got != `firecrawl scrape "https://example.com": request failed with status 503: upstream unavailable` {
		t.Fatalf("Error() = %q", got)
	}
	if got := reqErr.afterAttempts(3).Error(); got != `firecrawl scrape "https://example.com": request failed after 3 attempts with status 503: upstream unavailable` {
		t.Fatalf("afterAttempts() = %q", got)
	}
}

func TestBackoffDuration(t *testing.T) {
	t.Parallel()

	if got := backoffDuration(0, 0); got != defaultRetryBackoff {
		t.Fatalf("backoffDuration(0, 0) = %v, want %v", got, defaultRetryBackoff)
	}
	if got := backoffDuration(time.Millisecond, 2); got != 2*time.Millisecond {
		t.Fatalf("backoffDuration(1ms, 2) = %v, want %v", got, 2*time.Millisecond)
	}
}

func TestSleepContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := sleepContext(ctx, time.Millisecond); !errors.Is(err, context.Canceled) {
		t.Fatalf("sleepContext canceled error = %v, want context canceled", err)
	}

	if err := sleepContext(context.Background(), time.Millisecond); err != nil {
		t.Fatalf("sleepContext success error = %v, want nil", err)
	}
}

func TestParseErrorMessageFallbacks(t *testing.T) {
	t.Parallel()

	if got := parseErrorMessage(nil); got != "" {
		t.Fatalf("parseErrorMessage(nil) = %q, want empty string", got)
	}
	if got := parseErrorMessage([]byte("plain failure")); got != "plain failure" {
		t.Fatalf("parseErrorMessage(plain) = %q, want %q", got, "plain failure")
	}
	if got := parseErrorMessage([]byte(`{"message":"upstream failed"}`)); got != "upstream failed" {
		t.Fatalf("parseErrorMessage(json) = %q, want %q", got, "upstream failed")
	}
}

func writeJSONStatus(t *testing.T, w http.ResponseWriter, statusCode int, body string) {
	t.Helper()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write([]byte(body)); err != nil {
		t.Fatalf("write response: %v", err)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}
