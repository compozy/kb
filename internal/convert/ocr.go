//go:build ocr

package convert

import (
	"context"
	"fmt"
	"strings"

	"github.com/otiai10/gosseract/v2"
	"github.com/user/kb/internal/models"
)

// ImageConverter renders supported image inputs through OCR when the `ocr`
// build tag is enabled.
type ImageConverter struct{}

// Accepts reports whether the input is an OCR-capable image.
func (ImageConverter) Accepts(ext string, mimeType string) bool {
	return acceptsImage(ext, mimeType)
}

// Convert extracts OCR text and lightweight EXIF metadata from supported image
// formats.
func (ImageConverter) Convert(ctx context.Context, input models.ConvertInput) (*models.ConvertResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	data, err := readInput(input)
	if err != nil {
		return nil, err
	}

	client := gosseract.NewClient()
	defer client.Close()

	if languages := imageOCRLanguages(input.Options); len(languages) > 0 {
		if err := client.SetLanguage(languages...); err != nil {
			return nil, fmt.Errorf("convert: configure OCR languages: %w", err)
		}
	}

	if err := client.SetImageFromBytes(data); err != nil {
		return nil, fmt.Errorf("convert: load image for OCR: %w", err)
	}

	text, err := client.Text()
	if err != nil {
		return nil, fmt.Errorf("convert: OCR image: %w", err)
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	text = strings.TrimSpace(text)

	metadata := extractImageMetadata(data)
	if metadata == nil {
		metadata = map[string]any{}
	}
	if text == "" {
		metadata["warnings"] = []string{officeNoTextWarning}
	}

	return &models.ConvertResult{
		Markdown: text,
		Title:    firstNonEmptyLine(text),
		Metadata: metadata,
	}, nil
}

func imageOCRLanguages(options map[string]any) []string {
	if len(options) == 0 {
		return nil
	}

	for _, key := range []string{"ocrLanguages", "ocr_languages", "languages"} {
		raw, ok := options[key]
		if !ok {
			continue
		}

		switch typed := raw.(type) {
		case string:
			typed = strings.TrimSpace(typed)
			if typed == "" {
				return nil
			}
			return []string{typed}
		case []string:
			langs := make([]string, 0, len(typed))
			for _, lang := range typed {
				lang = strings.TrimSpace(lang)
				if lang != "" {
					langs = append(langs, lang)
				}
			}
			if len(langs) > 0 {
				return langs
			}
		case []any:
			langs := make([]string, 0, len(typed))
			for _, item := range typed {
				lang, ok := item.(string)
				if !ok {
					continue
				}
				lang = strings.TrimSpace(lang)
				if lang != "" {
					langs = append(langs, lang)
				}
			}
			if len(langs) > 0 {
				return langs
			}
		}
	}

	return nil
}
