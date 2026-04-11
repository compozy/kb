package convert

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/user/go-devstack/internal/models"
)

func TestNewRegistryRegistersDefaultConvertersInPriorityOrder(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()

	cases := []struct {
		name     string
		ext      string
		wantType any
		mimeType string
	}{
		{name: "text", ext: ".txt", wantType: TextConverter{}},
		{name: "html", ext: ".html", wantType: HTMLConverter{}},
		{name: "htm", ext: ".htm", wantType: HTMLConverter{}},
		{name: "pdf", ext: ".pdf", wantType: PDFConverter{}},
		{name: "csv", ext: ".csv", wantType: CSVConverter{}},
		{name: "json", ext: ".json", wantType: JSONConverter{}},
		{name: "xml", ext: ".xml", wantType: XMLConverter{}},
		{name: "markdown by mime", ext: "", mimeType: "text/markdown", wantType: TextConverter{}},
		{name: "html by mime", ext: "", mimeType: "text/html", wantType: HTMLConverter{}},
		{name: "pdf by mime", ext: "", mimeType: "application/pdf", wantType: PDFConverter{}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := registry.Match(tc.ext, tc.mimeType)
			if got == nil {
				t.Fatal("Match returned nil")
			}
			if reflect.TypeOf(got) != reflect.TypeOf(tc.wantType) {
				t.Fatalf("converter type = %T, want %T", got, tc.wantType)
			}
		})
	}
}

func TestRegistryMatchReturnsNilForUnsupportedInput(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()
	if got := registry.Match(".docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"); got != nil {
		t.Fatalf("Match returned %T, want nil", got)
	}
}

func TestRegistryConvertReturnsUnsupportedInputError(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()

	_, err := registry.Convert(context.Background(), models.ConvertInput{
		Reader:   strings.NewReader("binary"),
		FilePath: "report.docx",
		Options: map[string]any{
			"mimeType": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		},
	})
	if err == nil {
		t.Fatal("expected Convert to fail")
	}

	var unsupportedErr *UnsupportedInputError
	if !errors.As(err, &unsupportedErr) {
		t.Fatalf("expected UnsupportedInputError, got %T", err)
	}
	if !strings.Contains(err.Error(), ".docx") {
		t.Fatalf("error %q does not mention extension", err)
	}
}

func TestRegistryConvertDelegatesToFirstMatchingConverter(t *testing.T) {
	t.Parallel()

	first := &stubConverter{
		accepts: true,
		result:  &models.ConvertResult{Markdown: "first"},
	}
	second := &stubConverter{
		accepts: true,
		result:  &models.ConvertResult{Markdown: "second"},
	}

	registry := NewRegistry(first, second)

	result, err := registry.Convert(context.Background(), models.ConvertInput{
		Reader:   strings.NewReader("anything"),
		FilePath: "sample.txt",
	})
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}
	if result.Markdown != "first" {
		t.Fatalf("markdown = %q, want first", result.Markdown)
	}
	if first.convertCalls != 1 {
		t.Fatalf("first convert calls = %d, want 1", first.convertCalls)
	}
	if second.convertCalls != 0 {
		t.Fatalf("second convert calls = %d, want 0", second.convertCalls)
	}
}

func TestRegistryConvertMatchesByURLAndMIMEOption(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()

	result, err := registry.Convert(context.Background(), models.ConvertInput{
		Reader: strings.NewReader(`{"name":"URL Match"}`),
		URL:    "https://example.com/export",
		Options: map[string]any{
			"mime_type": "application/json; charset=utf-8",
		},
	})
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}
	if result.Title != "URL Match" {
		t.Fatalf("title = %q, want URL Match", result.Title)
	}
}

func TestNormalizeHelpers(t *testing.T) {
	t.Parallel()

	if got := normalizeExtension("README.MD"); got != ".md" {
		t.Fatalf("normalizeExtension README.MD = %q, want .md", got)
	}
	if got := normalizeExtension("json"); got != ".json" {
		t.Fatalf("normalizeExtension json = %q, want .json", got)
	}
	if got := normalizeMIMEType("text/markdown; charset=utf-8"); got != "text/markdown" {
		t.Fatalf("normalizeMIMEType = %q, want text/markdown", got)
	}
	if got := inputMIMEType(map[string]any{"contentType": "application/xml; charset=utf-8"}); got != "application/xml" {
		t.Fatalf("inputMIMEType = %q, want application/xml", got)
	}
}

func TestReadInputRequiresReader(t *testing.T) {
	t.Parallel()

	_, err := readInput(models.ConvertInput{})
	if err == nil {
		t.Fatal("expected readInput to fail")
	}
	if !strings.Contains(err.Error(), "reader is required") {
		t.Fatalf("error = %q", err)
	}
}

func TestUnsupportedInputErrorWithoutDetails(t *testing.T) {
	t.Parallel()

	err := &UnsupportedInputError{}
	if got := err.Error(); got != "convert: no converter accepts the provided input" {
		t.Fatalf("error = %q", got)
	}
}

func TestTextConverterConvertsTXTAndExtractsFirstLineTitle(t *testing.T) {
	t.Parallel()

	converter := TextConverter{}
	result, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   strings.NewReader("First title\nSecond line\n"),
		FilePath: "notes.txt",
	})
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}

	if result.Title != "First title" {
		t.Fatalf("title = %q, want First title", result.Title)
	}
	if result.Markdown != "First title\nSecond line\n" {
		t.Fatalf("markdown = %q", result.Markdown)
	}
}

func TestTextConverterPreservesMarkdownFormatting(t *testing.T) {
	t.Parallel()

	converter := TextConverter{}
	input := "# Heading\n\n- item\n"

	result, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   strings.NewReader(input),
		FilePath: "notes.md",
	})
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}

	if result.Markdown != input {
		t.Fatalf("markdown = %q, want %q", result.Markdown, input)
	}
	if result.Title != "Heading" {
		t.Fatalf("title = %q, want Heading", result.Title)
	}
}

func TestCSVConverterProducesMarkdownTable(t *testing.T) {
	t.Parallel()

	converter := CSVConverter{}
	result, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   strings.NewReader("name,kind,count\nalpha,service,3\nbeta,worker,5\n"),
		FilePath: "inventory.csv",
	})
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}

	want := strings.Join([]string{
		"| name | kind | count |",
		"| --- | --- | --- |",
		"| alpha | service | 3 |",
		"| beta | worker | 5 |",
		"",
	}, "\n")
	if result.Markdown != want {
		t.Fatalf("markdown = %q, want %q", result.Markdown, want)
	}
}

func TestCSVConverterHandlesHeaderOnlyAndEscapesSpecialCharacters(t *testing.T) {
	t.Parallel()

	converter := CSVConverter{}

	headerOnly, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   strings.NewReader("name,notes\n"),
		FilePath: "header-only.csv",
	})
	if err != nil {
		t.Fatalf("Convert returned error for header-only CSV: %v", err)
	}
	wantHeaderOnly := strings.Join([]string{
		"| name | notes |",
		"| --- | --- |",
		"",
	}, "\n")
	if headerOnly.Markdown != wantHeaderOnly {
		t.Fatalf("header-only markdown = %q, want %q", headerOnly.Markdown, wantHeaderOnly)
	}

	special, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   strings.NewReader("name,notes\nalpha,\"has | pipe and \"\"quotes\"\"\"\n"),
		FilePath: "special.csv",
	})
	if err != nil {
		t.Fatalf("Convert returned error for special CSV: %v", err)
	}
	if !strings.Contains(special.Markdown, `has \| pipe and "quotes"`) {
		t.Fatalf("special markdown = %q", special.Markdown)
	}
}

func TestJSONConverterWrapsContentAndExtractsTitle(t *testing.T) {
	t.Parallel()

	converter := JSONConverter{}
	result, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   strings.NewReader(`{"title":"KB Pivot","count":2}`),
		FilePath: "doc.json",
	})
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}

	if result.Title != "KB Pivot" {
		t.Fatalf("title = %q, want KB Pivot", result.Title)
	}
	if !strings.HasPrefix(result.Markdown, "```json\n{\n") {
		t.Fatalf("markdown = %q", result.Markdown)
	}
	if !strings.Contains(result.Markdown, `"count": 2`) {
		t.Fatalf("markdown = %q", result.Markdown)
	}
	if !strings.HasSuffix(result.Markdown, "\n```") {
		t.Fatalf("markdown = %q", result.Markdown)
	}
	if got := result.Metadata["count"]; got != float64(2) {
		t.Fatalf("metadata count = %#v, want 2", got)
	}
}

func TestJSONConverterFallsBackToNameField(t *testing.T) {
	t.Parallel()

	converter := JSONConverter{}
	result, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   strings.NewReader(`{"name":"Fallback Name"}`),
		FilePath: "named.json",
	})
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}
	if result.Title != "Fallback Name" {
		t.Fatalf("title = %q, want Fallback Name", result.Title)
	}
}

func TestJSONConverterRejectsEmptyInput(t *testing.T) {
	t.Parallel()

	converter := JSONConverter{}
	_, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   strings.NewReader("  "),
		FilePath: "empty.json",
	})
	if err == nil {
		t.Fatal("expected empty JSON to fail")
	}
}

func TestXMLConverterExtractsTextAndRejectsInvalidInput(t *testing.T) {
	t.Parallel()

	converter := XMLConverter{}

	result, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   strings.NewReader("<root><title>Spec</title><body>Hello <b>world</b></body></root>"),
		FilePath: "doc.xml",
	})
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}
	if result.Title != "Spec" {
		t.Fatalf("title = %q, want Spec", result.Title)
	}
	if result.Markdown != "Spec Hello world" {
		t.Fatalf("markdown = %q, want %q", result.Markdown, "Spec Hello world")
	}

	_, err = converter.Convert(context.Background(), models.ConvertInput{
		Reader:   strings.NewReader("<root>"),
		FilePath: "broken.xml",
	})
	if err == nil {
		t.Fatal("expected malformed XML to fail")
	}
}

func TestXMLConverterRejectsEmptyInput(t *testing.T) {
	t.Parallel()

	converter := XMLConverter{}
	_, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   strings.NewReader(" "),
		FilePath: "empty.xml",
	})
	if err == nil {
		t.Fatal("expected empty XML to fail")
	}
}

type stubConverter struct {
	accepts      bool
	result       *models.ConvertResult
	err          error
	convertCalls int
}

func (s *stubConverter) Accepts(string, string) bool {
	return s.accepts
}

func (s *stubConverter) Convert(context.Context, models.ConvertInput) (*models.ConvertResult, error) {
	s.convertCalls++
	return s.result, s.err
}
