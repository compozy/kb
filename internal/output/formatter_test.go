package output_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/user/go-devstack/internal/output"
)

func TestFormatOutputTableAlignsColumns(t *testing.T) {
	t.Parallel()

	rows := []map[string]any{
		{"complexity": 2, "file": "internal/cli/root.go", "name": "Execute"},
		{"complexity": 12, "file": "internal/output/formatter.go", "name": "FormatOutput"},
		{"complexity": 5, "file": "internal/vault/query.go", "name": "ResolveVault"},
		{"complexity": 3, "file": "internal/metrics/compute.go", "name": "ComputeMetrics"},
		{"complexity": 1, "file": "internal/scanner/scanner.go", "name": "ScanWorkspace"},
	}

	rendered := output.FormatOutput(output.FormatOptions{
		Format:  output.OutputFormatTable,
		Columns: []string{"name", "file", "complexity"},
		Data:    rows,
	})

	lines := strings.Split(strings.TrimSuffix(rendered, "\n"), "\n")
	if len(lines) != 7 {
		t.Fatalf("expected 7 lines, got %d: %q", len(lines), rendered)
	}

	header := lines[0]
	divider := lines[1]
	fileColumn := strings.Index(header, "file")
	complexityColumn := strings.Index(header, "complexity")

	if fileColumn <= 0 || complexityColumn <= fileColumn {
		t.Fatalf("unexpected header alignment %q", header)
	}

	dividerParts := strings.Split(divider, "  ")
	if len(dividerParts) != 3 {
		t.Fatalf("unexpected divider shape %q", divider)
	}

	for _, part := range dividerParts {
		if part == "" || strings.Trim(part, "-") != "" {
			t.Fatalf("unexpected divider segment %q in %q", part, divider)
		}
	}

	for rowIndex, row := range rows {
		line := lines[rowIndex+2]
		fileValue := row["file"].(string)
		complexityValue := fmt.Sprint(row["complexity"])

		if got := strings.Index(line, fileValue); got != fileColumn {
			t.Fatalf("row %d file column starts at %d, want %d in %q", rowIndex, got, fileColumn, line)
		}

		if got := strings.Index(line, complexityValue); got != complexityColumn {
			t.Fatalf("row %d complexity column starts at %d, want %d in %q", rowIndex, got, complexityColumn, line)
		}
	}
}

func TestFormatOutputTableHandlesVariableWidthsAndTruncation(t *testing.T) {
	t.Parallel()

	longCell := strings.Repeat("very-long-cell-", 6)
	rendered := output.FormatOutput(output.FormatOptions{
		Format:  output.OutputFormatTable,
		Columns: []string{"name", "summary"},
		Data: []map[string]any{
			{"name": "short", "summary": "tiny"},
			{"name": "long", "summary": longCell},
		},
	})

	lines := strings.Split(strings.TrimSuffix(rendered, "\n"), "\n")
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines, got %d: %q", len(lines), rendered)
	}

	secondRow := lines[3]
	if strings.Contains(secondRow, longCell) {
		t.Fatalf("expected long cell to be truncated, got %q", secondRow)
	}

	if !strings.Contains(secondRow, "...") {
		t.Fatalf("expected truncated row to contain ellipsis, got %q", secondRow)
	}
}

func TestFormatOutputJSONProducesValidProjectedObjects(t *testing.T) {
	t.Parallel()

	rendered := output.FormatOutput(output.FormatOptions{
		Format:  output.OutputFormatJSON,
		Columns: []string{"name", "file", "meta"},
		Data: []map[string]any{
			{
				"file":    "src/index.ts",
				"ignored": true,
				"meta": map[string]any{
					"owner": "cli",
					"score": 5,
				},
				"name": "greet",
			},
		},
	})

	var decoded []map[string]any
	if err := json.Unmarshal([]byte(rendered), &decoded); err != nil {
		t.Fatalf("expected valid JSON, got error: %v", err)
	}

	if len(decoded) != 1 {
		t.Fatalf("expected 1 row, got %d", len(decoded))
	}

	row := decoded[0]
	if _, exists := row["ignored"]; exists {
		t.Fatalf("expected ignored column to be omitted, got %v", row)
	}

	if row["name"] != "greet" || row["file"] != "src/index.ts" {
		t.Fatalf("unexpected projected row %v", row)
	}

	meta, ok := row["meta"].(map[string]any)
	if !ok {
		t.Fatalf("expected nested meta object, got %T", row["meta"])
	}

	if meta["owner"] != "cli" {
		t.Fatalf("unexpected nested meta object %v", meta)
	}
}

func TestFormatOutputTSVRendersHeaderAndSanitizesCells(t *testing.T) {
	t.Parallel()

	rendered := output.FormatOutput(output.FormatOptions{
		Format:  output.OutputFormatTSV,
		Columns: []string{"name", "note"},
		Data: []map[string]any{
			{"name": "greet", "note": "line1\tline2\nline3"},
		},
	})

	expected := "name\tnote\ngreet\tline1 line2 line3\n"
	if rendered != expected {
		t.Fatalf("unexpected TSV output:\n%s\nwant:\n%s", rendered, expected)
	}
}

func TestFormatOutputEmptyData(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		format output.OutputFormat
		want   string
	}{
		{name: "table", format: output.OutputFormatTable, want: "No results.\n"},
		{name: "json", format: output.OutputFormatJSON, want: "[]\n"},
		{name: "tsv", format: output.OutputFormatTSV, want: "name\n"},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := output.FormatOutput(output.FormatOptions{
				Format:  testCase.format,
				Columns: []string{"name"},
				Data:    nil,
			})

			if got != testCase.want {
				t.Fatalf("FormatOutput(%s) = %q, want %q", testCase.format, got, testCase.want)
			}
		})
	}
}

func TestFormatOutputSingleRowInAllFormats(t *testing.T) {
	t.Parallel()

	row := []map[string]any{{
		"file": "src/index.ts",
		"name": "greet",
	}}

	tableRendered := output.FormatOutput(output.FormatOptions{
		Format:  output.OutputFormatTable,
		Columns: []string{"name", "file"},
		Data:    row,
	})
	tableLines := strings.Split(strings.TrimSuffix(tableRendered, "\n"), "\n")
	if len(tableLines) != 3 {
		t.Fatalf("expected 3 table lines, got %d: %q", len(tableLines), tableRendered)
	}

	jsonRendered := output.FormatOutput(output.FormatOptions{
		Format:  output.OutputFormatJSON,
		Columns: []string{"name", "file"},
		Data:    row,
	})
	if !strings.Contains(jsonRendered, "\"name\": \"greet\"") || !strings.Contains(jsonRendered, "\"file\": \"src/index.ts\"") {
		t.Fatalf("unexpected JSON output %q", jsonRendered)
	}

	tsvRendered := output.FormatOutput(output.FormatOptions{
		Format:  output.OutputFormatTSV,
		Columns: []string{"name", "file"},
		Data:    row,
	})
	if tsvRendered != "name\tfile\ngreet\tsrc/index.ts\n" {
		t.Fatalf("unexpected TSV output %q", tsvRendered)
	}
}

func TestFormatOutputSpecialCharactersRemainValidJSON(t *testing.T) {
	t.Parallel()

	rendered := output.FormatOutput(output.FormatOptions{
		Format:  output.OutputFormatJSON,
		Columns: []string{"name", "quote", "values"},
		Data: []map[string]any{
			{
				"name":   "greet",
				"quote":  `say "hi"`,
				"values": []string{"a", "b"},
			},
		},
	})

	var decoded []struct {
		Name   string   `json:"name"`
		Quote  string   `json:"quote"`
		Values []string `json:"values"`
	}
	if err := json.Unmarshal([]byte(rendered), &decoded); err != nil {
		t.Fatalf("expected valid JSON, got error: %v", err)
	}

	if len(decoded) != 1 || decoded[0].Quote != `say "hi"` {
		t.Fatalf("unexpected decoded payload %+v", decoded)
	}

	if strings.Join(decoded[0].Values, ",") != "a,b" {
		t.Fatalf("unexpected array payload %+v", decoded[0].Values)
	}
}

func TestFormatOutputDefaultsUnsupportedFormatsToTable(t *testing.T) {
	t.Parallel()

	rendered := output.FormatOutput(output.FormatOptions{
		Format:  output.OutputFormat("unsupported"),
		Columns: []string{"name"},
		Data:    []map[string]any{{"name": "greet"}},
	})

	if !strings.HasPrefix(rendered, "name \n-----") {
		t.Fatalf("expected fallback table output, got %q", rendered)
	}
}
