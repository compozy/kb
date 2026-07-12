// Package browser renders JavaScript pages through a configured Chromium command.
package browser

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/compozy/kb/internal/urlfetch"
)

const (
	defaultTimeout  = 30 * time.Second
	defaultMaxBytes = 20 * 1024 * 1024
)

// Client renders a page's DOM using Chrome or Chromium in headless mode.
type Client struct {
	command  string
	now      func() time.Time
	timeout  time.Duration
	maxBytes int64
}

// NewClient constructs a browser client. command must resolve to a Chromium
// executable that supports --headless=new and --dump-dom.
func NewClient(command string) *Client {
	return &Client{command: strings.TrimSpace(command), now: time.Now, timeout: defaultTimeout, maxBytes: defaultMaxBytes}
}

// Fetch renders sourceURL and returns its serialized DOM as HTML.
func (client *Client) Fetch(ctx context.Context, sourceURL string) (*urlfetch.Result, error) {
	if client == nil {
		return nil, errors.New("browser fetch: client is nil")
	}
	if client.command == "" {
		return nil, errors.New("browser fetch: missing browser command; set browser.command or use --provider http-local")
	}
	if err := urlfetch.ValidatePublicURL(ctx, sourceURL); err != nil {
		return nil, err
	}
	if ctx == nil {
		ctx = context.Background()
	}
	timeout := client.timeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	command := exec.CommandContext(ctx, client.command, "--headless=new", "--disable-gpu", "--dump-dom", sourceURL)
	stdout, err := command.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("browser fetch %q: open renderer output: %w", sourceURL, err)
	}
	if err := command.Start(); err != nil {
		return nil, fmt.Errorf("browser fetch %q: start renderer: %w", sourceURL, err)
	}
	maxBytes := client.maxBytes
	if maxBytes <= 0 {
		maxBytes = defaultMaxBytes
	}
	body, readErr := io.ReadAll(io.LimitReader(stdout, maxBytes+1))
	if int64(len(body)) > maxBytes {
		_ = command.Process.Kill()
		_ = command.Wait()
		return nil, fmt.Errorf("browser fetch %q: rendered page exceeds maximum size of %d bytes", sourceURL, maxBytes)
	}
	if waitErr := command.Wait(); waitErr != nil || readErr != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, fmt.Errorf("browser fetch %q: render canceled: %w", sourceURL, ctxErr)
		}
		if readErr != nil {
			return nil, fmt.Errorf("browser fetch %q: read rendered page: %w", sourceURL, readErr)
		}
		return nil, fmt.Errorf("browser fetch %q: render page: %w", sourceURL, waitErr)
	}
	if len(strings.TrimSpace(string(body))) == 0 {
		return nil, fmt.Errorf("browser fetch %q: rendered page is empty", sourceURL)
	}
	now := client.now
	if now == nil {
		now = time.Now
	}
	hash := sha256.Sum256(body)
	return &urlfetch.Result{
		SourceURL:   sourceURL,
		FinalURL:    sourceURL,
		ContentType: "text/html",
		FileName:    "rendered.html",
		Body:        body,
		FetchedAt:   now().UTC(),
		ContentHash: fmt.Sprintf("sha256:%x", hash),
	}, nil
}
