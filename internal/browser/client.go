// Package browser renders JavaScript pages through a configured Chromium command.
package browser

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/compozy/kb/internal/urlfetch"
)

// Client renders a page's DOM using Chrome or Chromium in headless mode.
type Client struct {
	command string
	now     func() time.Time
}

// NewClient constructs a browser client. command must resolve to a Chromium
// executable that supports --headless=new and --dump-dom.
func NewClient(command string) *Client {
	return &Client{command: strings.TrimSpace(command), now: time.Now}
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

	command := exec.CommandContext(ctx, client.command, "--headless=new", "--disable-gpu", "--dump-dom", sourceURL)
	body, err := command.Output()
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, fmt.Errorf("browser fetch %q: render canceled: %w", sourceURL, ctxErr)
		}
		return nil, fmt.Errorf("browser fetch %q: render page: %w", sourceURL, err)
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
