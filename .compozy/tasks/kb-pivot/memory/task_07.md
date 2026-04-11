# Task Memory: task_07.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot

- Implement `internal/convert/pdf.go` plus tests and fixtures for task 07, keeping scope to PDF conversion, registry registration, and required validation coverage.

## Important Decisions

- Use `pdfcpu` for metadata, validation, encryption detection, and per-page content stream access; build Markdown text extraction on top of the extracted page content because `pdfcpu` does not expose a high-level plain-text extraction API.
- Keep page concatenation explicit with Markdown page separators so multi-page output is readable and testable.
- Disable the pdfcpu config directory for this package at init time because PDF conversion only needs built-in defaults and pdfcpu's lazy config initialization is not safe under concurrent first use.

## Learnings

- `api.PDFInfo` returns title, author, page count, and encrypted state, while `pdfcpu.ExtractPageContent` returns decoded page content in PDF syntax rather than ready-made plain text.
- `make verify` runs package tests with `-race`, so pdfcpu integration has to be safe under concurrent test execution, not just functionally correct.

## Files / Surfaces

- `go.mod`
- `go.sum`
- `internal/convert/pdf.go`
- `internal/convert/pdf_test.go`
- `internal/convert/registry.go`
- `internal/convert/registry_test.go`
- `internal/convert/testdata/`
- `.compozy/tasks/kb-pivot/task_07.md`
- `.compozy/tasks/kb-pivot/_tasks.md`

## Errors / Corrections

- Fixed a parser bug where `Tm` line-break detection used the X coordinate instead of the Y coordinate, which collapsed headings and paragraph boundaries.
- Fixed a concurrent first-use race from pdfcpu's config initialization instead of weakening test parallelism.

## Ready for Next Run

- Task implementation, targeted converter coverage, and repo-wide `make verify` are complete; only commit state and any downstream review feedback would remain.
