---
status: pending
title: Implement HTML-to-Markdown converter
type: backend
complexity: medium
dependencies:
  - task_05
---

# Task 06: Implement HTML-to-Markdown converter

## Overview

Add an HTML-to-Markdown converter to the `internal/convert/` package using the `JohannesKaufmann/html-to-markdown` v2 library. This converter is critical — it handles direct HTML files and is reused by the EPUB converter (task_09) for XHTML chapter conversion.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST convert .html and .htm files to Markdown preserving headings, lists, tables, code blocks, and links
- MUST extract `<title>` tag content as the document title
- MUST strip script and style tags before conversion
- MUST register with the converter registry from task_05
- MUST add `github.com/JohannesKaufmann/html-to-markdown/v2` dependency via `go get`
- SHOULD expose a reusable `HTMLToMarkdown(htmlContent string) (string, error)` function for use by other converters (EPUB)
</requirements>

## Subtasks

- [ ] 6.1 Add html-to-markdown v2 dependency
- [ ] 6.2 Implement HTML converter with title extraction and tag stripping
- [ ] 6.3 Expose reusable HTMLToMarkdown function for other converters
- [ ] 6.4 Register converter in the registry
- [ ] 6.5 Write unit tests with fixture HTML files

## Implementation Details

Create `internal/convert/html.go` and `internal/convert/html_test.go`. Add the dependency via `go get github.com/JohannesKaufmann/html-to-markdown/v2`.

### Relevant Files

- `internal/convert/registry.go` (task_05) — register this converter
- `internal/models/kb_models.go` (task_02) — Converter interface to implement

### Dependent Files

- `internal/convert/epub.go` (task_09) — reuses HTMLToMarkdown for XHTML chapter conversion

### Related ADRs

- [ADR-003: Native Go Document Conversion with Converter Registry](../adrs/adr-003.md) — converter selection rationale

## Deliverables

- `internal/convert/html.go` — HTML converter + reusable HTMLToMarkdown function
- `internal/convert/html_test.go` — tests
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests

- Unit tests:
  - [ ] Converts simple HTML with headings, paragraphs, and links to Markdown
  - [ ] Converts HTML table to Markdown table
  - [ ] Converts HTML code blocks (pre/code) to fenced Markdown code blocks
  - [ ] Converts HTML ordered and unordered lists to Markdown lists
  - [ ] Extracts `<title>` content as document title
  - [ ] Strips `<script>` and `<style>` tags before conversion
  - [ ] Handles HTML with no `<title>` tag (falls back to first `<h1>` or empty)
  - [ ] Handles empty/malformed HTML gracefully
  - [ ] Accepts .html and .htm extensions
- Test coverage target: >=80%
- All tests must pass

## Success Criteria

- All tests passing
- Test coverage >=80%
- HTML-to-Markdown output is clean and readable in Obsidian
- `make lint` reports zero findings
