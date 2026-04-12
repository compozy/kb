package convert

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/user/kb/internal/models"
)

func TestHTMLConverterAcceptsExtensionsAndMIMETypes(t *testing.T) {
	t.Parallel()

	converter := HTMLConverter{}

	cases := []struct {
		ext      string
		mimeType string
		want     bool
	}{
		{ext: ".html", want: true},
		{ext: ".htm", want: true},
		{mimeType: "text/html", want: true},
		{mimeType: "application/xhtml+xml", want: true},
		{ext: ".xml", want: false},
		{mimeType: "application/xml", want: false},
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

func TestHTMLConverterConvertSimpleHTML(t *testing.T) {
	t.Parallel()

	result := convertHTMLFixture(t, "simple.html")

	if result.Title != "Fixture Title" {
		t.Fatalf("title = %q, want Fixture Title", result.Title)
	}

	for _, want := range []string{
		"# Welcome to Kodebase",
		"This is a [reference link](https://example.com/docs).",
	} {
		if !strings.Contains(result.Markdown, want) {
			t.Fatalf("markdown %q does not contain %q", result.Markdown, want)
		}
	}
}

func TestHTMLConverterConvertTable(t *testing.T) {
	t.Parallel()

	result := convertHTMLFixture(t, "table.html")
	rows := markdownTableRows(result.Markdown)

	if len(rows) < 4 {
		t.Fatalf("table rows = %v, want at least 4 rows", rows)
	}
	if !reflect.DeepEqual(rows[0], []string{"Name", "Value"}) {
		t.Fatalf("header row = %v", rows[0])
	}
	if !strings.Contains(strings.Join(rows[1], " "), "---") {
		t.Fatalf("separator row = %v, want markdown separator", rows[1])
	}
	if !reflect.DeepEqual(rows[2], []string{"Alpha", "1"}) {
		t.Fatalf("first data row = %v", rows[2])
	}
	if !reflect.DeepEqual(rows[3], []string{"Beta", "2"}) {
		t.Fatalf("second data row = %v", rows[3])
	}
}

func TestHTMLConverterConvertCodeBlocksAndLists(t *testing.T) {
	t.Parallel()

	result := convertHTMLFixture(t, "code_lists.htm")

	for _, want := range []string{
		"```go",
		`fmt.Println("hello")`,
		"```",
		"1. First step",
		"2. Second step",
		"- Apples",
		"- Oranges",
	} {
		if !strings.Contains(result.Markdown, want) {
			t.Fatalf("markdown %q does not contain %q", result.Markdown, want)
		}
	}
}

func TestHTMLConverterStripsScriptAndStyleAndExtractsTitle(t *testing.T) {
	t.Parallel()

	result := convertHTMLFixture(t, "script_style.html")

	if result.Title != "Sanitized Doc" {
		t.Fatalf("title = %q, want Sanitized Doc", result.Title)
	}
	for _, unwanted := range []string{"console.log", "display:none", "body {", "alert("} {
		if strings.Contains(result.Markdown, unwanted) {
			t.Fatalf("markdown %q unexpectedly contains %q", result.Markdown, unwanted)
		}
	}
	if !strings.Contains(result.Markdown, "# Clean Output") {
		t.Fatalf("markdown %q does not contain cleaned heading", result.Markdown)
	}
}

func TestHTMLConverterFallsBackToFirstHeadingWhenTitleMissing(t *testing.T) {
	t.Parallel()

	result := convertHTMLFixture(t, "no_title.html")

	if result.Title != "Fallback Heading" {
		t.Fatalf("title = %q, want Fallback Heading", result.Title)
	}
	if !strings.Contains(result.Markdown, "# Fallback Heading") {
		t.Fatalf("markdown %q does not contain h1 heading", result.Markdown)
	}
}

func TestHTMLConverterHandlesEmptyInputGracefully(t *testing.T) {
	t.Parallel()

	result := convertHTMLFixture(t, "empty.html")

	if result.Title != "" {
		t.Fatalf("title = %q, want empty", result.Title)
	}
	if result.Markdown != "" {
		t.Fatalf("markdown = %q, want empty", result.Markdown)
	}
}

func TestHTMLConverterHandlesMalformedHTMLGracefully(t *testing.T) {
	t.Parallel()

	result := convertHTMLFixture(t, "malformed.html")

	if result.Title != "Broken Heading" {
		t.Fatalf("title = %q, want Broken Heading", result.Title)
	}
	for _, want := range []string{"# Broken Heading", "Paragraph without closing tags."} {
		if !strings.Contains(result.Markdown, want) {
			t.Fatalf("markdown %q does not contain %q", result.Markdown, want)
		}
	}
}

func TestHTMLToMarkdownReusableHelper(t *testing.T) {
	t.Parallel()

	markdown, err := HTMLToMarkdown(readHTMLFixture(t, "table.html"))
	if err != nil {
		t.Fatalf("HTMLToMarkdown returned error: %v", err)
	}

	rows := markdownTableRows(markdown)
	if len(rows) < 4 {
		t.Fatalf("table rows = %v, want at least 4 rows", rows)
	}
}

func convertHTMLFixture(t *testing.T, fixtureName string) *models.ConvertResult {
	t.Helper()

	converter := HTMLConverter{}
	content := readHTMLFixture(t, fixtureName)

	result, err := converter.Convert(context.Background(), models.ConvertInput{
		Reader:   strings.NewReader(content),
		FilePath: fixtureName,
	})
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}

	return result
}

func readHTMLFixture(t *testing.T, fixtureName string) string {
	t.Helper()

	path := filepath.Join("testdata", "html", fixtureName)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) returned error: %v", path, err)
	}

	return string(data)
}

func markdownTableRows(markdown string) [][]string {
	lines := strings.Split(strings.TrimSpace(markdown), "\n")
	rows := make([][]string, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "|") || !strings.HasSuffix(line, "|") {
			continue
		}

		parts := strings.Split(line, "|")
		cells := make([]string, 0, len(parts)-2)
		for _, part := range parts[1 : len(parts)-1] {
			cells = append(cells, strings.TrimSpace(part))
		}
		rows = append(rows, cells)
	}

	return rows
}
