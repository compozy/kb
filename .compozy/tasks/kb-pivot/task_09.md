---
status: completed
title: Implement EPUB and image OCR converters
type: backend
complexity: medium
dependencies:
    - task_05
    - task_06
---

# Task 09: Implement EPUB and image OCR converters

## Overview

Add EPUB and image OCR converters to `internal/convert/`. EPUB uses the ZIP+XHTML approach (unzip, read spine order from content.opf, convert each XHTML chapter via the HTML converter from task_06). Image OCR uses `gosseract` (Tesseract bindings) gated behind a build tag so the Tesseract system dependency is optional.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST convert .epub files by unzipping, reading `content.opf` for spine order, and converting each XHTML chapter via the HTMLToMarkdown function from task_06
- MUST extract EPUB metadata (title, author, language) from `content.opf`
- MUST implement image OCR converter for .png, .jpg, .jpeg, .tiff, .bmp using `gosseract`
- MUST gate OCR behind a Go build tag (`//go:build ocr`) so Tesseract is not a mandatory dependency
- MUST provide a no-op OCR fallback (extracts EXIF metadata only) when built without the `ocr` tag
- MUST add `github.com/otiai10/gosseract/v2` as an optional dependency
</requirements>

## Subtasks

- [ ] 9.1 Implement EPUB converter with ZIP+XHTML parsing and spine-ordered chapter conversion
- [ ] 9.2 Implement image OCR converter with gosseract behind build tag
- [ ] 9.3 Implement no-op image fallback (EXIF metadata extraction) for non-OCR builds
- [ ] 9.4 Register both converters in the registry
- [ ] 9.5 Write unit tests (EPUB with fixture, OCR with build-tag-aware test)

## Implementation Details

Create `internal/convert/epub.go`, `internal/convert/ocr.go` (build tag: `ocr`), and `internal/convert/ocr_noop.go` (build tag: `!ocr`).

EPUB structure: unzip → read `META-INF/container.xml` for rootfile path → read `content.opf` for `<spine>` element → read each XHTML itemref in order → convert via `HTMLToMarkdown()`.

### Relevant Files

- `internal/convert/html.go` (task_06) — HTMLToMarkdown reusable function for EPUB chapter conversion
- `internal/convert/registry.go` (task_05) — register converters
- `internal/models/kb_models.go` (task_02) — Converter interface

### Dependent Files

- None — these are leaf converters

### Related ADRs

- [ADR-003: Native Go Document Conversion with Converter Registry](../adrs/adr-003.md) — EPUB as ZIP+XHTML, OCR as optional

## Deliverables

- `internal/convert/epub.go` — EPUB converter
- `internal/convert/ocr.go` — OCR converter (build tag: ocr)
- `internal/convert/ocr_noop.go` — no-op fallback (build tag: !ocr)
- Test files
- `internal/convert/testdata/sample.epub` — small test fixture
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests

- Unit tests:
  - [ ] EPUB converter extracts chapters in spine order as concatenated Markdown
  - [ ] EPUB converter extracts title and author from content.opf metadata
  - [ ] EPUB converter handles EPUB with missing content.opf gracefully
  - [ ] EPUB converter handles EPUB with empty chapters
  - [ ] OCR converter (with tag) extracts text from a simple image with clear text
  - [ ] OCR no-op fallback returns empty body with EXIF metadata only
  - [ ] Image converter accepts .png, .jpg, .jpeg, .tiff, .bmp extensions
  - [ ] Image converter rejects non-image extensions
- Test coverage target: >=80%
- All tests must pass

## Success Criteria

- All tests passing
- Test coverage >=80%
- EPUB output preserves chapter order and structure
- Build succeeds both with and without the `ocr` build tag
- `make lint` reports zero findings
