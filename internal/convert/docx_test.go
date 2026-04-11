package convert

import (
	"bytes"
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/user/go-devstack/internal/models"
)

func TestDOCXConverterAcceptsExtensionsAndMIMETypes(t *testing.T) {
	t.Parallel()

	converter := DOCXConverter{}

	cases := []struct {
		ext      string
		mimeType string
		want     bool
	}{
		{ext: ".docx", want: true},
		{mimeType: docxMIMEType, want: true},
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

func TestDOCXConverterConvertsParagraphsHeadingsAndMetadata(t *testing.T) {
	t.Parallel()

	result := convertDOCXFixture(t, "sample.docx")

	wantMarkdown := strings.Join([]string{
		"# Quarterly Review",
		"",
		"This paragraph survives DOCX conversion.",
		"",
		"## Next Steps",
		"",
		"Ship the Office converters.",
	}, "\n")
	if result.Markdown != wantMarkdown {
		t.Fatalf("markdown = %q, want %q", result.Markdown, wantMarkdown)
	}
	if result.Title != "DOCX Fixture Title" {
		t.Fatalf("title = %q, want DOCX Fixture Title", result.Title)
	}

	wantMetadata := map[string]any{
		"title":  "DOCX Fixture Title",
		"author": "DOCX Fixture Author",
	}
	for key, expected := range wantMetadata {
		if got, ok := result.Metadata[key]; !ok || !reflect.DeepEqual(got, expected) {
			t.Fatalf("metadata[%q] = %#v, want %#v", key, got, expected)
		}
	}
}

func TestDOCXConverterHandlesEmptyDocumentGracefully(t *testing.T) {
	t.Parallel()

	result := convertDOCXFixture(t, "empty.docx")

	if result.Markdown != "" {
		t.Fatalf("markdown = %q, want empty", result.Markdown)
	}
	if result.Title != "Empty DOCX Title" {
		t.Fatalf("title = %q, want Empty DOCX Title", result.Title)
	}

	warnings := warningList(t, result.Metadata)
	if !reflect.DeepEqual(warnings, []string{officeNoTextWarning}) {
		t.Fatalf("warnings = %#v, want %#v", warnings, []string{officeNoTextWarning})
	}
}

func TestDOCXConverterReturnsErrorForCorruptedZip(t *testing.T) {
	t.Parallel()

	converter := DOCXConverter{}
	_, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   bytes.NewReader([]byte("not a zip")),
		FilePath: "broken.docx",
	})
	if err == nil {
		t.Fatal("expected Convert to fail")
	}
	if !strings.Contains(err.Error(), "invalid DOCX content") {
		t.Fatalf("error = %q, want invalid DOCX content message", err)
	}
}

func TestDOCXHeadingLevelParsesHeadingStyles(t *testing.T) {
	t.Parallel()

	cases := map[string]int{
		"Heading1": 1,
		"heading2": 2,
		"Heading9": 6,
		"Title":    0,
		"":         0,
	}

	for style, want := range cases {
		if got := docxHeadingLevel(style); got != want {
			t.Fatalf("docxHeadingLevel(%q) = %d, want %d", style, got, want)
		}
	}
}

func convertDOCXFixture(t *testing.T, fixtureName string) *models.ConvertResult {
	t.Helper()

	converter := DOCXConverter{}
	result, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   bytes.NewReader(readConvertFixtureBytes(t, fixtureName)),
		FilePath: fixtureName,
	})
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}

	return result
}
