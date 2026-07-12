package urlfetch

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestFetchDownloadsHTMLTracksRedirectAndHashesContent(t *testing.T) {
	t.Parallel()

	server := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/start":
			http.Redirect(writer, request, "/article", http.StatusFound)
		case "/article":
			writer.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = writer.Write([]byte("<html><title>Article</title><body>Hello</body></html>"))
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	client := testClient(server)
	client.now = func() time.Time { return time.Date(2026, 7, 12, 15, 0, 0, 0, time.UTC) }

	result, err := client.Fetch(context.Background(), serverURL(server, "/start"))
	if err != nil {
		t.Fatalf("fetch: %v", err)
	}
	if result.FinalURL != serverURL(server, "/article") {
		t.Fatalf("final URL = %q", result.FinalURL)
	}
	if result.ContentType != "text/html" || result.FileName != "article.html" {
		t.Fatalf("content metadata = %#v", result)
	}
	if result.ContentHash != "sha256:e329e0eeffd1f2763140d6a776b6056d08fd21e43f486749dd945ea55a787f5d" {
		t.Fatalf("content hash = %q", result.ContentHash)
	}
	if !result.FetchedAt.Equal(time.Date(2026, 7, 12, 15, 0, 0, 0, time.UTC)) {
		t.Fatalf("fetched at = %s", result.FetchedAt)
	}
}

func TestFetchRejectsNonHTTPSAndPrivateAddresses(t *testing.T) {
	t.Parallel()

	client := NewClient()
	for _, rawURL := range []string{
		"http://example.com/article",
		"https://127.0.0.1/article",
		"https://[::1]/article",
		"https://localhost/article",
	} {
		if _, err := client.Fetch(context.Background(), rawURL); err == nil {
			t.Fatalf("expected %q to be rejected", rawURL)
		}
	}
}

func TestFetchRejectsRedirectToPrivateAddress(t *testing.T) {
	t.Parallel()

	server := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, "https://internal.example/private", http.StatusFound)
	}))
	defer server.Close()

	client := testClient(server)
	client.lookupIP = func(_ context.Context, host string) ([]net.IP, error) {
		if host == "internal.example" {
			return []net.IP{net.ParseIP("10.0.0.7")}, nil
		}
		return []net.IP{net.ParseIP("93.184.216.34")}, nil
	}

	_, err := client.Fetch(context.Background(), serverURL(server, "/"))
	if err == nil || !strings.Contains(err.Error(), "non-public IP") {
		t.Fatalf("error = %v, want private redirect rejection", err)
	}
}

func TestFetchRejectsStatusAndOversizedResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/missing":
			http.NotFound(writer, request)
		case "/rate-limited":
			writer.Header().Set("Retry-After", "1")
			http.Error(writer, "slow down", http.StatusTooManyRequests)
		case "/blocked":
			http.Error(writer, "blocked", http.StatusForbidden)
		default:
			_, _ = writer.Write([]byte("too large"))
		}
	}))
	defer server.Close()

	client := testClient(server)
	if _, err := client.Fetch(context.Background(), serverURL(server, "/missing")); err == nil || !strings.Contains(err.Error(), "status 404") {
		t.Fatalf("error = %v, want 404", err)
	}
	if _, err := client.Fetch(context.Background(), serverURL(server, "/rate-limited")); err == nil || !strings.Contains(err.Error(), "status 429") {
		t.Fatalf("error = %v, want 429", err)
	}
	if _, err := client.Fetch(context.Background(), serverURL(server, "/blocked")); err == nil || !strings.Contains(err.Error(), "status 403") {
		t.Fatalf("error = %v, want 403", err)
	}
	client.maxBytes = 3
	if _, err := client.Fetch(context.Background(), serverURL(server, "/large")); err == nil || !strings.Contains(err.Error(), "maximum size") {
		t.Fatalf("error = %v, want size limit", err)
	}
}

func testClient(server *httptest.Server) *Client {
	serverURL, err := url.Parse(server.URL)
	if err != nil {
		panic(err)
	}
	transport := server.Client().Transport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec // test server only
	transport.DialContext = func(ctx context.Context, network, _ string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, network, serverURL.Host)
	}
	client := NewClient()
	client.httpClient = &http.Client{Transport: transport}
	client.lookupIP = func(_ context.Context, _ string) ([]net.IP, error) {
		return []net.IP{net.ParseIP("93.184.216.34")}, nil
	}
	return client
}

func serverURL(server *httptest.Server, suffix string) string {
	return "https://public.example" + suffix
}
