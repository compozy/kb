---
status: pending
title: Implement PDF converter
type: backend
complexity: medium
dependencies:
  - task_05
---

# Task 07: Implement PDF converter

## Overview

Add a PDF-to-Markdown converter to the `internal/convert/` package using `pdfcpu` for text extraction. PDF ingestion is a core use case for knowledge base building — research papers, whitepapers, and documentation are commonly distributed as PDFs.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST convert .pdf files to Markdown text using `pdfcpu` library
- MUST extract document metadata (title, author, page count) into ConvertResult.Metadata
- MUST preserve paragraph structure and heading hierarchy where detectable
- MUST handle multi-page documents by concatenating page text with page separators
- MUST return a clear error for encrypted/password-protected PDFs
- MUST add `github.com/pdfcpu/pdfcpu` dependency via `go get`
</requirements>

## Subtasks

- [ ] 7.1 Add pdfcpu dependency
- [ ] 7.2 Implement PDF converter with text extraction and metadata reading
- [ ] 7.3 Handle page-by-page extraction with separators
- [ ] 7.4 Register converter in the registry
- [ ] 7.5 Write unit tests with small fixture PDF files

## Implementation Details

Create `internal/convert/pdf.go` and `internal/convert/pdf_test.go`. Add fixture PDF files under `internal/convert/testdata/`.

Study markitdown's PDF converter approach (uses pdfminer.six in Python) for structure preservation strategies.

### Relevant Files

- `internal/convert/registry.go` (task_05) — register this converter
- `internal/models/kb_models.go` (task_02) — Converter interface to implement

### Dependent Files

- None — PDF converter is a leaf dependency

### Related ADRs

- [ADR-003: Native Go Document Conversion with Converter Registry](../adrs/adr-003.md) — library selection rationale (pdfcpu chosen for Apache-2.0 license)

## Deliverables

- `internal/convert/pdf.go` — PDF converter
- `internal/convert/pdf_test.go` — tests
- `internal/convert/testdata/sample.pdf` — small test fixture
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests

- Unit tests:
  - [ ] Converts single-page PDF with text to Markdown
  - [ ] Converts multi-page PDF with page separators in output
  - [ ] Extracts PDF metadata (title, author) into ConvertResult.Metadata
  - [ ] Reports page count in metadata
  - [ ] Returns error for encrypted PDF
  - [ ] Returns error for non-PDF file with .pdf extension
  - [ ] Handles PDF with no extractable text (scanned image) — returns empty body with warning metadata
  - [ ] Accepts .pdf extension only
- Test coverage target: >=80%
- All tests must pass

## Success Criteria

- All tests passing
- Test coverage >=80%
- Extracted text is readable and preserves paragraph structure
- `make lint` reports zero findings
