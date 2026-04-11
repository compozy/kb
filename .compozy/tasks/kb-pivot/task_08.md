---
status: pending
title: Implement Office format converters (DOCX, PPTX, XLSX)
type: backend
complexity: high
dependencies:
  - task_05
---

# Task 08: Implement Office format converters (DOCX, PPTX, XLSX)

## Overview

Add converters for Microsoft Office formats to the `internal/convert/` package. DOCX and PPTX are handled via ZIP+XML parsing (the same approach markitdown uses), while XLSX uses the `excelize` library. These are among the most commonly shared document formats for knowledge base ingestion.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST convert .docx files by unzipping and parsing `word/document.xml` to extract text with paragraph/heading structure
- MUST convert .pptx files by unzipping and parsing `ppt/slides/slide*.xml` to extract text per slide
- MUST convert .xlsx files using `excelize` to extract sheets as Markdown tables
- MUST extract document metadata (title, author) from Office XML properties where available
- MUST handle malformed or empty Office files gracefully with structured errors
- MUST add `github.com/xuri/excelize/v2` dependency via `go get`
</requirements>

## Subtasks

- [ ] 8.1 Add excelize dependency
- [ ] 8.2 Implement DOCX converter with ZIP+XML paragraph/heading extraction
- [ ] 8.3 Implement PPTX converter with ZIP+XML per-slide text extraction
- [ ] 8.4 Implement XLSX converter using excelize for sheet-to-table conversion
- [ ] 8.5 Extract document properties (title, author) from docProps/core.xml
- [ ] 8.6 Register all three converters in the registry
- [ ] 8.7 Write unit tests with fixture Office files

## Implementation Details

Create `internal/convert/docx.go`, `internal/convert/pptx.go`, `internal/convert/xlsx.go` and corresponding test files. DOCX and PPTX share a common ZIP+XML parsing approach — extract a shared helper if the overlap is significant.

For DOCX: the main content is in `word/document.xml` under `<w:body>` → `<w:p>` (paragraphs) → `<w:r>` (runs) → `<w:t>` (text). Heading styles are in `<w:pStyle w:val="Heading1">`.

For PPTX: slides are in `ppt/slides/slide1.xml` etc. Text is in `<a:t>` elements under `<p:txBody>`.

For XLSX: use excelize's `GetRows()` per sheet and format as Markdown tables.

### Relevant Files

- `internal/convert/registry.go` (task_05) — register converters
- `internal/models/kb_models.go` (task_02) — Converter interface

### Dependent Files

- None — Office converters are leaf dependencies

### Related ADRs

- [ADR-003: Native Go Document Conversion with Converter Registry](../adrs/adr-003.md) — ZIP+XML approach chosen over heavy Office libraries

## Deliverables

- `internal/convert/docx.go` — DOCX converter
- `internal/convert/pptx.go` — PPTX converter
- `internal/convert/xlsx.go` — XLSX converter
- Test files for all three
- `internal/convert/testdata/` — small fixture .docx, .pptx, .xlsx files
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests

- Unit tests:
  - [ ] DOCX converter extracts paragraphs as Markdown text
  - [ ] DOCX converter converts heading styles to Markdown headings (# H1, ## H2)
  - [ ] DOCX converter extracts title from document properties
  - [ ] DOCX converter handles empty document gracefully
  - [ ] PPTX converter extracts text from multiple slides with slide separators
  - [ ] PPTX converter handles slides with no text (image-only) gracefully
  - [ ] XLSX converter converts single sheet to Markdown table with headers
  - [ ] XLSX converter converts multiple sheets with sheet name headers
  - [ ] XLSX converter handles empty sheet gracefully
  - [ ] All converters return error for corrupted ZIP files
- Test coverage target: >=80%
- All tests must pass

## Success Criteria

- All tests passing
- Test coverage >=80%
- Extracted text preserves document structure (headings, paragraphs, slides, tables)
- `make lint` reports zero findings
