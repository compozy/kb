// Package output renders structured data as table, JSON, or TSV for CLI display.
package output

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"
)

const (
	maxTableCellWidth = 60
	tableEllipsis     = "..."
)

// OutputFormat controls how tabular CLI data is rendered.
type OutputFormat string

const (
	// OutputFormatTable renders aligned ASCII columns.
	OutputFormatTable OutputFormat = "table"
	// OutputFormatJSON renders rows as a JSON array of objects.
	OutputFormatJSON OutputFormat = "json"
	// OutputFormatTSV renders rows as tab-separated values with a header row.
	OutputFormatTSV OutputFormat = "tsv"
)

// FormatOptions describes the tabular payload to render for CLI output.
type FormatOptions struct {
	Format  OutputFormat
	Columns []string
	Data    []map[string]any
}

// FormatOutput renders tabular data using the requested output format.
func FormatOutput(options FormatOptions) string {
	switch options.Format {
	case OutputFormatJSON:
		return formatJSON(options.Columns, options.Data)
	case OutputFormatTSV:
		return formatTSV(options.Columns, options.Data)
	case OutputFormatTable:
		fallthrough
	default:
		return formatTable(options.Columns, options.Data)
	}
}

func formatTable(columns []string, data []map[string]any) string {
	if len(data) == 0 {
		return "No results.\n"
	}

	if len(columns) == 0 {
		return "\n"
	}

	rows := projectStringRows(columns, data, true)
	widths := make([]int, len(columns))

	for columnIndex, column := range columns {
		widths[columnIndex] = runeCount(column)
		for _, row := range rows {
			if width := runeCount(row[columnIndex]); width > widths[columnIndex] {
				widths[columnIndex] = width
			}
		}
	}

	dividerParts := make([]string, len(widths))
	for index, width := range widths {
		dividerParts[index] = strings.Repeat("-", width)
	}

	lines := make([]string, 0, len(rows)+3)
	lines = append(lines, formatStringRow(columns, widths))
	lines = append(lines, strings.Join(dividerParts, "  "))

	for _, row := range rows {
		lines = append(lines, formatStringRow(row, widths))
	}

	lines = append(lines, "")
	return strings.Join(lines, "\n")
}

func formatJSON(columns []string, data []map[string]any) string {
	if len(data) == 0 {
		return "[]\n"
	}

	var builder strings.Builder
	builder.WriteString("[\n")

	for rowIndex, row := range data {
		builder.WriteString("  {")
		wroteField := false

		for _, column := range columns {
			value, ok := row[column]
			if !ok {
				continue
			}

			if wroteField {
				builder.WriteString(",\n")
			} else {
				builder.WriteString("\n")
			}

			keyJSON, _ := json.Marshal(column)
			valueJSON, err := json.Marshal(normalizeJSONValue(value))
			if err != nil {
				valueJSON, _ = json.Marshal(fmt.Sprint(value))
			}

			builder.WriteString("    ")
			builder.Write(keyJSON)
			builder.WriteString(": ")
			builder.Write(valueJSON)
			wroteField = true
		}

		if wroteField {
			builder.WriteString("\n  ")
		}

		builder.WriteString("}")
		if rowIndex < len(data)-1 {
			builder.WriteString(",\n")
			continue
		}

		builder.WriteString("\n")
	}

	builder.WriteString("]\n")
	return builder.String()
}

func formatTSV(columns []string, data []map[string]any) string {
	if len(columns) == 0 {
		return "\n"
	}

	rows := projectStringRows(columns, data, false)
	lines := make([]string, 0, len(rows)+1)
	lines = append(lines, strings.Join(columns, "\t"))

	for _, row := range rows {
		lines = append(lines, strings.Join(row, "\t"))
	}

	return strings.Join(lines, "\n") + "\n"
}

func projectStringRows(columns []string, data []map[string]any, truncate bool) [][]string {
	rows := make([][]string, 0, len(data))

	for _, row := range data {
		projected := make([]string, len(columns))
		for columnIndex, column := range columns {
			cell := sanitizeInlineValue(normalizeCellValue(row[column]))
			if truncate {
				cell = truncateTableCell(cell)
			}

			projected[columnIndex] = cell
		}

		rows = append(rows, projected)
	}

	return rows
}

func normalizeCellValue(value any) string {
	if value == nil {
		return ""
	}

	if stringer, ok := value.(fmt.Stringer); ok {
		return stringer.String()
	}

	if text, ok := value.(string); ok {
		return text
	}

	resolved := dereferenceValue(value)
	if resolved == nil {
		return ""
	}

	reflected := reflect.ValueOf(resolved)
	switch reflected.Kind() {
	case reflect.Slice, reflect.Array:
		parts := make([]string, 0, reflected.Len())
		for index := 0; index < reflected.Len(); index++ {
			parts = append(parts, normalizeCellValue(reflected.Index(index).Interface()))
		}

		return strings.Join(parts, ", ")
	case reflect.Map, reflect.Struct:
		encoded, err := json.Marshal(normalizeJSONValue(resolved))
		if err == nil {
			return string(encoded)
		}
	}

	return fmt.Sprint(resolved)
}

func normalizeJSONValue(value any) any {
	resolved := dereferenceValue(value)
	if resolved == nil {
		return nil
	}

	reflected := reflect.ValueOf(resolved)
	switch reflected.Kind() {
	case reflect.Map:
		normalized := make(map[string]any, reflected.Len())
		iterator := reflected.MapRange()
		for iterator.Next() {
			normalized[fmt.Sprint(iterator.Key().Interface())] = normalizeJSONValue(iterator.Value().Interface())
		}

		return normalized
	case reflect.Slice, reflect.Array:
		normalized := make([]any, reflected.Len())
		for index := 0; index < reflected.Len(); index++ {
			normalized[index] = normalizeJSONValue(reflected.Index(index).Interface())
		}

		return normalized
	case reflect.Func, reflect.Chan, reflect.Complex64, reflect.Complex128, reflect.UnsafePointer:
		return fmt.Sprint(resolved)
	default:
		if _, err := json.Marshal(resolved); err == nil {
			return resolved
		}

		return fmt.Sprint(resolved)
	}
}

func dereferenceValue(value any) any {
	reflected := reflect.ValueOf(value)
	for reflected.IsValid() && reflected.Kind() == reflect.Pointer {
		if reflected.IsNil() {
			return nil
		}

		reflected = reflected.Elem()
	}

	if !reflected.IsValid() {
		return nil
	}

	return reflected.Interface()
}

func sanitizeInlineValue(value string) string {
	if value == "" {
		return ""
	}

	var builder strings.Builder
	builder.Grow(len(value))
	lastWasReplacement := false

	for _, r := range value {
		if r == '\r' || r == '\n' || r == '\t' {
			if !lastWasReplacement {
				builder.WriteByte(' ')
				lastWasReplacement = true
			}

			continue
		}

		builder.WriteRune(r)
		lastWasReplacement = false
	}

	return strings.TrimSpace(builder.String())
}

func truncateTableCell(value string) string {
	if runeCount(value) <= maxTableCellWidth {
		return value
	}

	limit := maxTableCellWidth - runeCount(tableEllipsis)
	if limit <= 0 {
		ellipsisRunes := []rune(tableEllipsis)
		if maxTableCellWidth >= len(ellipsisRunes) {
			return tableEllipsis
		}

		return string(ellipsisRunes[:maxTableCellWidth])
	}

	return string([]rune(value)[:limit]) + tableEllipsis
}

func formatStringRow(values []string, widths []int) string {
	parts := make([]string, len(values))
	for index, value := range values {
		parts[index] = padRight(value, widths[index])
	}

	return strings.Join(parts, "  ")
}

func padRight(value string, width int) string {
	padding := width - runeCount(value)
	if padding <= 0 {
		return value
	}

	return value + strings.Repeat(" ", padding)
}

func runeCount(value string) int {
	return utf8.RuneCountInString(value)
}
