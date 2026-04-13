package convert

import (
	"context"

	"github.com/compozy/kb/internal/models"
)

// TextConverter passes plain text and Markdown files through as Markdown.
type TextConverter struct{}

// Accepts reports whether the input is plain text or Markdown.
func (TextConverter) Accepts(ext string, mimeType string) bool {
	switch normalizeExtension(ext) {
	case ".txt", ".md":
		return true
	}

	switch normalizeMIMEType(mimeType) {
	case "text/plain", "text/markdown", "text/x-markdown":
		return true
	}

	return false
}

// Convert returns the input text unchanged and extracts the first title-like
// line for downstream metadata.
func (TextConverter) Convert(ctx context.Context, input models.ConvertInput) (*models.ConvertResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	data, err := readInput(input)
	if err != nil {
		return nil, err
	}

	markdown := string(data)
	return &models.ConvertResult{
		Markdown: markdown,
		Title:    firstNonEmptyLine(markdown),
	}, nil
}
