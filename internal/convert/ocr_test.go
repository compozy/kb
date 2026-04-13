//go:build ocr

package convert

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/compozy/kb/internal/models"
)

func TestImageConverterWithOCRExtractsText(t *testing.T) {
	t.Parallel()

	converter := ImageConverter{}
	data := makeOCRPNG(t, "KODEBASE")

	result, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   bytes.NewReader(data),
		FilePath: "fixture.png",
		Options: map[string]any{
			"ocrLanguages": []string{"eng"},
		},
	})
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}

	if !strings.Contains(normalizeOCRText(result.Markdown), "KODEBASE") {
		t.Fatalf("markdown = %q, want OCR text to contain KODEBASE", result.Markdown)
	}
}
