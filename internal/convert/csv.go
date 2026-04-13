package convert

import (
	"context"
	"encoding/csv"
	"errors"
	"strings"

	"github.com/compozy/kb/internal/models"
)

// CSVConverter renders CSV content as a Markdown table.
type CSVConverter struct{}

// Accepts reports whether the input is CSV.
func (CSVConverter) Accepts(ext string, mimeType string) bool {
	return normalizeExtension(ext) == ".csv" ||
		normalizeMIMEType(mimeType) == "text/csv" ||
		normalizeMIMEType(mimeType) == "application/csv"
}

// Convert parses CSV rows and writes a Markdown table.
func (CSVConverter) Convert(ctx context.Context, input models.ConvertInput) (*models.ConvertResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	data, err := readInput(input)
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, errors.New("convert: csv input is empty")
	}

	reader := csv.NewReader(strings.NewReader(string(data)))
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, errors.New("convert: csv input is empty")
	}

	width := maxCSVWidth(records)
	normalized := make([][]string, len(records))
	for i, record := range records {
		normalized[i] = padCSVRecord(record, width)
	}

	header := normalized[0]
	rows := normalized[1:]

	var builder strings.Builder
	writeMarkdownTableRow(&builder, header)
	writeMarkdownTableSeparator(&builder, width)
	for _, row := range rows {
		writeMarkdownTableRow(&builder, row)
	}

	return &models.ConvertResult{
		Markdown: builder.String(),
	}, nil
}

func maxCSVWidth(records [][]string) int {
	width := 0
	for _, record := range records {
		if len(record) > width {
			width = len(record)
		}
	}
	if width == 0 {
		return 1
	}

	return width
}

func padCSVRecord(record []string, width int) []string {
	if len(record) >= width {
		return record
	}

	padded := make([]string, width)
	copy(padded, record)
	return padded
}

func writeMarkdownTableRow(builder *strings.Builder, row []string) {
	builder.WriteString("|")
	for _, cell := range row {
		builder.WriteString(" ")
		builder.WriteString(escapeMarkdownTableCell(cell))
		builder.WriteString(" |")
	}
	builder.WriteString("\n")
}

func writeMarkdownTableSeparator(builder *strings.Builder, width int) {
	builder.WriteString("|")
	for i := 0; i < width; i++ {
		builder.WriteString(" --- |")
	}
	builder.WriteString("\n")
}

func escapeMarkdownTableCell(value string) string {
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")
	value = strings.ReplaceAll(value, "\n", "<br>")
	value = strings.ReplaceAll(value, "|", `\|`)
	return value
}
