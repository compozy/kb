package convert

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/user/kb/internal/models"
	"github.com/xuri/excelize/v2"
)

const xlsxMIMEType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

// XLSXConverter renders workbook sheets as Markdown tables.
type XLSXConverter struct{}

// Accepts reports whether the input is XLSX content.
func (XLSXConverter) Accepts(ext string, mimeType string) bool {
	return normalizeExtension(ext) == ".xlsx" || normalizeMIMEType(mimeType) == xlsxMIMEType
}

// Convert transforms workbook sheets into Markdown tables while surfacing core
// document metadata when present.
func (XLSXConverter) Convert(ctx context.Context, input models.ConvertInput) (*models.ConvertResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	archive, err := openOfficeArchive(input, "XLSX")
	if err != nil {
		return nil, err
	}

	props, err := archive.readCoreProperties("XLSX")
	if err != nil {
		return nil, err
	}

	workbook, err := excelize.OpenReader(bytes.NewReader(archive.data))
	if err != nil {
		return nil, fmt.Errorf("convert: invalid XLSX content: %w", err)
	}
	defer func() {
		_ = workbook.Close()
	}()

	sheetNames := workbook.GetSheetList()
	if len(sheetNames) == 0 {
		metadata := officeMetadata(props)
		addOfficeWarning(metadata, officeNoTextWarning)

		return &models.ConvertResult{
			Title:    officeDocumentTitle(props, ""),
			Metadata: metadata,
		}, nil
	}

	sections := make([]string, 0, len(sheetNames))
	emptySheets := make([]string, 0)
	includeSheetHeadings := len(sheetNames) > 1
	hasTable := false

	for _, sheetName := range sheetNames {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		rows, err := workbook.GetRows(sheetName)
		if err != nil {
			return nil, fmt.Errorf("convert: read XLSX sheet %q: %w", sheetName, err)
		}

		table := renderXLSXTable(rows)
		if table != "" {
			hasTable = true
		} else {
			emptySheets = append(emptySheets, sheetName)
			table = "_Empty sheet._"
		}

		if includeSheetHeadings {
			table = "## " + sheetName + "\n\n" + table
		}

		sections = append(sections, table)
	}

	metadata := officeMetadata(props)
	metadata["sheetCount"] = len(sheetNames)
	metadata["sheetNames"] = append([]string(nil), sheetNames...)
	if len(emptySheets) > 0 {
		metadata["emptySheets"] = emptySheets
	}
	if !hasTable {
		addOfficeWarning(metadata, officeNoTextWarning)
	}

	return &models.ConvertResult{
		Markdown: strings.Join(sections, "\n\n"),
		Title:    officeDocumentTitle(props, sheetNames[0]),
		Metadata: metadata,
	}, nil
}

func renderXLSXTable(rows [][]string) string {
	if sheetRowsEmpty(rows) {
		return ""
	}

	width := maxCSVWidth(rows)
	normalized := make([][]string, len(rows))
	for i, row := range rows {
		normalized[i] = padCSVRecord(row, width)
	}

	var builder strings.Builder
	writeMarkdownTableRow(&builder, normalized[0])
	writeMarkdownTableSeparator(&builder, width)
	for _, row := range normalized[1:] {
		writeMarkdownTableRow(&builder, row)
	}

	return strings.TrimRight(builder.String(), "\n")
}

func sheetRowsEmpty(rows [][]string) bool {
	for _, row := range rows {
		for _, cell := range row {
			if strings.TrimSpace(cell) != "" {
				return false
			}
		}
	}

	return true
}
