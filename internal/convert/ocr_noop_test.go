//go:build !ocr

package convert

import (
	"bytes"
	"context"
	"testing"

	"github.com/user/go-devstack/internal/models"
)

func TestImageConverterWithoutOCRReturnsMetadataOnly(t *testing.T) {
	t.Parallel()

	converter := ImageConverter{}
	data := makeJPEGWithEXIF(t, map[uint16]string{
		0x010E: "Fixture Image",
		0x013B: "Fixture Artist",
	})

	result, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   bytes.NewReader(data),
		FilePath: "fixture.jpg",
	})
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}

	if result.Markdown != "" {
		t.Fatalf("markdown = %q, want empty", result.Markdown)
	}
	if result.Title != "" {
		t.Fatalf("title = %q, want empty", result.Title)
	}

	want := map[string]any{
		"imageDescription": "Fixture Image",
		"artist":           "Fixture Artist",
	}
	for key, expected := range want {
		if got, ok := result.Metadata[key]; !ok || got != expected {
			t.Fatalf("metadata[%q] = %#v, want %#v", key, got, expected)
		}
	}
	if len(result.Metadata) != len(want) {
		t.Fatalf("metadata = %#v, want only EXIF keys", result.Metadata)
	}
}
