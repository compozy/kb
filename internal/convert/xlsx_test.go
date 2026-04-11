package convert

import (
	"bytes"
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/user/go-devstack/internal/models"
)

func TestXLSXConverterAcceptsExtensionsAndMIMETypes(t *testing.T) {
	t.Parallel()

	converter := XLSXConverter{}

	cases := []struct {
		ext      string
		mimeType string
		want     bool
	}{
		{ext: ".xlsx", want: true},
		{mimeType: xlsxMIMEType, want: true},
		{ext: ".csv", want: false},
		{mimeType: "text/csv", want: false},
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

func TestXLSXConverterConvertsSingleSheetToMarkdownTable(t *testing.T) {
	t.Parallel()

	result := convertXLSXFixture(t, "sample.xlsx")

	wantMarkdown := strings.Join([]string{
		"| Name | Value |",
		"| --- | --- |",
		"| Pages | 12 |",
		"| Slides | 4 |",
	}, "\n")
	if result.Markdown != wantMarkdown {
		t.Fatalf("markdown = %q, want %q", result.Markdown, wantMarkdown)
	}
	if result.Title != "XLSX Fixture Title" {
		t.Fatalf("title = %q, want XLSX Fixture Title", result.Title)
	}

	wantMetadata := map[string]any{
		"title":      "XLSX Fixture Title",
		"author":     "XLSX Fixture Author",
		"sheetCount": 1,
		"sheetNames": []string{"Metrics"},
	}
	for key, expected := range wantMetadata {
		if got, ok := result.Metadata[key]; !ok || !reflect.DeepEqual(got, expected) {
			t.Fatalf("metadata[%q] = %#v, want %#v", key, got, expected)
		}
	}
}

func TestXLSXConverterConvertsMultipleSheetsWithHeaders(t *testing.T) {
	t.Parallel()

	result := convertXLSXFixture(t, "multi_sheet.xlsx")

	for _, want := range []string{
		"## Summary",
		"| Metric | Count |",
		"| Files | 8 |",
		"## Empty Sheet",
	} {
		if !strings.Contains(result.Markdown, want) {
			t.Fatalf("markdown %q does not contain %q", result.Markdown, want)
		}
	}
}

func TestXLSXConverterHandlesEmptySheetGracefully(t *testing.T) {
	t.Parallel()

	result := convertXLSXFixture(t, "multi_sheet.xlsx")

	if !strings.Contains(result.Markdown, "_Empty sheet._") {
		t.Fatalf("markdown %q does not contain empty sheet marker", result.Markdown)
	}
	if got, ok := result.Metadata["emptySheets"]; !ok || !reflect.DeepEqual(got, []string{"Empty Sheet"}) {
		t.Fatalf("emptySheets = %#v, want %#v", got, []string{"Empty Sheet"})
	}
}

func TestXLSXConverterReturnsErrorForCorruptedZip(t *testing.T) {
	t.Parallel()

	converter := XLSXConverter{}
	_, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   bytes.NewReader([]byte("not a zip")),
		FilePath: "broken.xlsx",
	})
	if err == nil {
		t.Fatal("expected Convert to fail")
	}
	if !strings.Contains(err.Error(), "invalid XLSX content") {
		t.Fatalf("error = %q, want invalid XLSX content message", err)
	}
}

func TestRenderXLSXTableReturnsEmptyForBlankRows(t *testing.T) {
	t.Parallel()

	if got := renderXLSXTable([][]string{{"", ""}, {}}); got != "" {
		t.Fatalf("renderXLSXTable(blank) = %q, want empty", got)
	}
}

func convertXLSXFixture(t *testing.T, fixtureName string) *models.ConvertResult {
	t.Helper()

	converter := XLSXConverter{}
	result, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   bytes.NewReader(readConvertFixtureBytes(t, fixtureName)),
		FilePath: fixtureName,
	})
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}

	return result
}
