// Package urlfetch downloads public HTTPS resources for URL ingestion.
package urlfetch

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultMaxBytes = 20 * 1024 * 1024
	defaultTimeout  = 30 * time.Second
)

// Result is a downloaded URL resource and its provenance metadata.
type Result struct {
	SourceURL   string
	FinalURL    string
	ContentType string
	FileName    string
	Body        []byte
	FetchedAt   time.Time
	ContentHash string
}

// Client downloads a public resource while enforcing URL and size limits.
type Client struct {
	httpClient *http.Client
	lookupIP   func(context.Context, string) ([]net.IP, error)
	now        func() time.Time
	maxBytes   int64
}

// NewClient returns a client with conservative defaults for URL ingestion.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: defaultTimeout},
		lookupIP: func(ctx context.Context, host string) ([]net.IP, error) {
			return net.DefaultResolver.LookupIP(ctx, "ip", host)
		},
		now:      time.Now,
		maxBytes: defaultMaxBytes,
	}
}

// ValidatePublicURL checks that rawURL is a public HTTPS URL suitable for a
// network-capable provider such as a headless browser.
func ValidatePublicURL(ctx context.Context, rawURL string) error {
	_, err := NewClient().validateURL(ctx, strings.TrimSpace(rawURL))
	return err
}

// Fetch downloads sourceURL, following only redirects to public HTTPS hosts.
func (client *Client) Fetch(ctx context.Context, sourceURL string) (*Result, error) {
	if client == nil {
		return nil, errors.New("url fetch: client is nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	sourceURL = strings.TrimSpace(sourceURL)
	if _, err := client.validateURL(ctx, sourceURL); err != nil {
		return nil, err
	}

	httpClient := client.httpClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}
	requestClient := *httpClient
	previousRedirect := requestClient.CheckRedirect
	requestClient.CheckRedirect = func(request *http.Request, via []*http.Request) error {
		if _, err := client.validateURL(request.Context(), request.URL.String()); err != nil {
			return err
		}
		if previousRedirect != nil {
			return previousRedirect(request, via)
		}
		return nil
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return nil, fmt.Errorf("url fetch %q: build request: %w", sourceURL, err)
	}
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/pdf,application/vnd.openxmlformats-officedocument.wordprocessingml.document;q=0.9,*/*;q=0.1")
	request.Header.Set("User-Agent", "kb/1.0 (+https://github.com/compozy/kb)")

	response, err := requestClient.Do(request)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, fmt.Errorf("url fetch %q: request canceled: %w", sourceURL, ctxErr)
		}
		return nil, fmt.Errorf("url fetch %q: request failed: %w", sourceURL, err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("url fetch %q: request failed with status %d", sourceURL, response.StatusCode)
	}

	maxBytes := client.maxBytes
	if maxBytes <= 0 {
		maxBytes = defaultMaxBytes
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, maxBytes+1))
	if err != nil {
		return nil, fmt.Errorf("url fetch %q: read response: %w", sourceURL, err)
	}
	if int64(len(body)) > maxBytes {
		return nil, fmt.Errorf("url fetch %q: response exceeds maximum size of %d bytes", sourceURL, maxBytes)
	}

	finalURL := response.Request.URL.String()
	contentType := normalizeContentType(response.Header.Get("Content-Type"))
	if contentType == "" {
		contentType = normalizeContentType(http.DetectContentType(body))
	}
	now := client.now
	if now == nil {
		now = time.Now
	}
	hash := sha256.Sum256(body)

	return &Result{
		SourceURL:   sourceURL,
		FinalURL:    finalURL,
		ContentType: contentType,
		FileName:    fileNameForURL(finalURL, contentType),
		Body:        body,
		FetchedAt:   now().UTC(),
		ContentHash: fmt.Sprintf("sha256:%x", hash),
	}, nil
}

func (client *Client) validateURL(ctx context.Context, rawURL string) (*url.URL, error) {
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil || parsed == nil || parsed.Host == "" {
		return nil, fmt.Errorf("url fetch: invalid url %q", rawURL)
	}
	if parsed.Scheme != "https" {
		return nil, fmt.Errorf("url fetch: only HTTPS URLs are allowed, got %q", parsed.Scheme)
	}
	if parsed.User != nil {
		return nil, errors.New("url fetch: URLs with user info are not allowed")
	}

	host := strings.TrimSuffix(strings.ToLower(parsed.Hostname()), ".")
	if host == "" || host == "localhost" || strings.HasSuffix(host, ".localhost") {
		return nil, fmt.Errorf("url fetch: host %q is not allowed", parsed.Hostname())
	}
	if address, err := netip.ParseAddr(host); err == nil {
		if isPrivateAddress(address) {
			return nil, fmt.Errorf("url fetch: host %q resolves to a non-public IP address", host)
		}
		return parsed, nil
	}

	lookupIP := client.lookupIP
	if lookupIP == nil {
		lookupIP = func(ctx context.Context, host string) ([]net.IP, error) {
			return net.DefaultResolver.LookupIP(ctx, "ip", host)
		}
	}
	addresses, err := lookupIP(ctx, host)
	if err != nil {
		return nil, fmt.Errorf("url fetch: resolve host %q: %w", host, err)
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf("url fetch: host %q did not resolve to an IP address", host)
	}
	for _, address := range addresses {
		if parsedAddress, ok := netip.AddrFromSlice(address); ok && isPrivateAddress(parsedAddress.Unmap()) {
			return nil, fmt.Errorf("url fetch: host %q resolves to a non-public IP address", host)
		}
	}

	return parsed, nil
}

func isPrivateAddress(address netip.Addr) bool {
	if !address.IsValid() {
		return true
	}
	if address.IsLoopback() || address.IsPrivate() || address.IsLinkLocalUnicast() || address.IsLinkLocalMulticast() || address.IsMulticast() || address.IsUnspecified() {
		return true
	}
	return netip.MustParsePrefix("100.64.0.0/10").Contains(address)
}

func normalizeContentType(value string) string {
	mediaType, _, err := mime.ParseMediaType(value)
	if err != nil {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(mediaType))
}

func fileNameForURL(rawURL string, contentType string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "source" + extensionForContentType(contentType)
	}
	name := path.Base(parsed.Path)
	if name == "." || name == "/" || name == "" {
		name = "source"
	}
	if filepath.Ext(name) == "" {
		name += extensionForContentType(contentType)
	}
	return name
}

func extensionForContentType(contentType string) string {
	switch contentType {
	case "text/html", "application/xhtml+xml":
		return ".html"
	case "application/pdf":
		return ".pdf"
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return ".docx"
	case "application/vnd.openxmlformats-officedocument.presentationml.presentation":
		return ".pptx"
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		return ".xlsx"
	case "application/json":
		return ".json"
	case "application/xml", "text/xml":
		return ".xml"
	default:
		return ".txt"
	}
}
