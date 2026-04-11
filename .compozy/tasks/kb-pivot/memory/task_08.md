# Task Memory: task_08.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot

- Implement native Office converters for DOCX, PPTX, and XLSX in `internal/convert/` with metadata extraction, registry wiring, binary fixtures, and unit coverage above the task threshold.

## Important Decisions

- Added a shared `internal/convert/office.go` helper for OOXML ZIP validation, `docProps/core.xml` parsing, shared warnings, and slide-file ordering instead of duplicating archive logic across converters.
- DOCX output maps `HeadingN` paragraph styles to ATX Markdown headings and otherwise preserves paragraphs as blank-line-separated blocks.
- PPTX output renders one Markdown section per slide using `## Slide N` headings plus `---` separators so blank slides still preserve slide positions without becoming hard failures.
- XLSX output uses the first row of each sheet as the Markdown table header; multi-sheet workbooks add `## <sheet name>` headers and blank sheets render `_Empty sheet._`.

## Learnings

- `go get` before the new import exists leaves `github.com/xuri/excelize/v2` marked indirect in `go.mod`; `go mod tidy` is needed after the converter lands to normalize the module entry.
- Minimal OOXML fixtures are sufficient for converter tests: DOCX and PPTX can be lightweight ZIPs with only the XML parts the converter reads, while XLSX still needs a valid workbook skeleton for `excelize`.

## Files / Surfaces

- `go.mod`
- `go.sum`
- `internal/convert/registry.go`
- `internal/convert/registry_test.go`
- `internal/convert/office.go`
- `internal/convert/docx.go`
- `internal/convert/pptx.go`
- `internal/convert/xlsx.go`
- `internal/convert/office_test.go`
- `internal/convert/docx_test.go`
- `internal/convert/pptx_test.go`
- `internal/convert/xlsx_test.go`
- `internal/convert/testdata/sample.docx`
- `internal/convert/testdata/empty.docx`
- `internal/convert/testdata/sample.pptx`
- `internal/convert/testdata/image_only.pptx`
- `internal/convert/testdata/sample.xlsx`
- `internal/convert/testdata/multi_sheet.xlsx`

## Errors / Corrections

- `make verify` initially failed on staticcheck `QF1012` in `internal/convert/pptx.go`; replaced `WriteString(fmt.Sprintf(...))` with `fmt.Fprintf` and reran the full gate.

## Ready for Next Run

- Task implementation is verified; remaining closeout is tracking-file updates and the local commit.
