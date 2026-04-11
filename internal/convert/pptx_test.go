package convert

import (
	"bytes"
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/user/go-devstack/internal/models"
)

func TestPPTXConverterAcceptsExtensionsAndMIMETypes(t *testing.T) {
	t.Parallel()

	converter := PPTXConverter{}

	cases := []struct {
		ext      string
		mimeType string
		want     bool
	}{
		{ext: ".pptx", want: true},
		{mimeType: pptxMIMEType, want: true},
		{ext: ".txt", want: false},
		{mimeType: "text/plain", want: false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.ext+tc.mimeType, func(t *testing.T) {
			t.Parallel()

			if got := converter.Accepts(tc.ext, tc.mimeType); got != tc.want {
				t.Fatalf("Accepts(%q, %q) = %t, want %t", tc.ext, tc.mimeType, got, tc.want)
			}
		})
	}
}

func TestPPTXConverterConvertsSlidesWithSeparatorsAndMetadata(t *testing.T) {
	t.Parallel()

	result := convertPPTXFixture(t, "sample.pptx")

	for _, want := range []string{
		"## Slide 1",
		"Launch Plan",
		"Milestone alpha",
		pptxSlideSeparator,
		"## Slide 2",
		"Roadmap",
		"Milestone beta",
	} {
		if !strings.Contains(result.Markdown, want) {
			t.Fatalf("markdown %q does not contain %q", result.Markdown, want)
		}
	}
	if result.Title != "PPTX Fixture Title" {
		t.Fatalf("title = %q, want PPTX Fixture Title", result.Title)
	}

	wantMetadata := map[string]any{
		"title":      "PPTX Fixture Title",
		"author":     "PPTX Fixture Author",
		"slideCount": 2,
	}
	for key, expected := range wantMetadata {
		if got, ok := result.Metadata[key]; !ok || !reflect.DeepEqual(got, expected) {
			t.Fatalf("metadata[%q] = %#v, want %#v", key, got, expected)
		}
	}
}

func TestPPTXConverterHandlesImageOnlySlidesGracefully(t *testing.T) {
	t.Parallel()

	result := convertPPTXFixture(t, "image_only.pptx")

	for _, want := range []string{"## Slide 1", "## Slide 2"} {
		if !strings.Contains(result.Markdown, want) {
			t.Fatalf("markdown %q does not contain %q", result.Markdown, want)
		}
	}
	if result.Title != "Image Only Deck" {
		t.Fatalf("title = %q, want Image Only Deck", result.Title)
	}

	warnings := warningList(t, result.Metadata)
	if !reflect.DeepEqual(warnings, []string{officeNoTextWarning}) {
		t.Fatalf("warnings = %#v, want %#v", warnings, []string{officeNoTextWarning})
	}
}

func TestPPTXConverterReturnsErrorForCorruptedZip(t *testing.T) {
	t.Parallel()

	converter := PPTXConverter{}
	_, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   bytes.NewReader([]byte("not a zip")),
		FilePath: "broken.pptx",
	})
	if err == nil {
		t.Fatal("expected Convert to fail")
	}
	if !strings.Contains(err.Error(), "invalid PPTX content") {
		t.Fatalf("error = %q, want invalid PPTX content message", err)
	}
}

func convertPPTXFixture(t *testing.T, fixtureName string) *models.ConvertResult {
	t.Helper()

	converter := PPTXConverter{}
	result, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   bytes.NewReader(readConvertFixtureBytes(t, fixtureName)),
		FilePath: fixtureName,
	})
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}

	return result
}
