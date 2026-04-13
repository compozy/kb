//go:build !ocr

package convert

import (
	"context"

	"github.com/compozy/kb/internal/models"
)

// ImageConverter preserves supported image metadata when OCR support is not
// compiled in.
type ImageConverter struct{}

// Accepts reports whether the input is an image handled by the OCR/no-op path.
func (ImageConverter) Accepts(ext string, mimeType string) bool {
	return acceptsImage(ext, mimeType)
}

// Convert returns an empty body and any extracted EXIF metadata.
func (ImageConverter) Convert(ctx context.Context, input models.ConvertInput) (*models.ConvertResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	data, err := readInput(input)
	if err != nil {
		return nil, err
	}

	return &models.ConvertResult{
		Markdown: "",
		Metadata: extractImageMetadata(data),
	}, nil
}
