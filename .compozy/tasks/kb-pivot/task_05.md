---
status: completed
title: Implement converter registry and simple converters
type: backend
complexity: medium
dependencies:
  - task_02
---

# Task 05: Implement converter registry and simple converters

## Overview

Create the `internal/convert/` package with the converter registry pattern and simple format converters (plain text, CSV, JSON, XML). The registry is the backbone of the ingest system — it matches files to converters by extension/MIME type and delegates conversion. Simple converters use only the Go standard library.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST implement the `Converter` interface as defined in TechSpec "Core Interfaces"
- MUST implement a `Registry` that stores converters in priority order and delegates to the first match
- MUST implement plain text converter (.txt, .md) that passes content through with title extraction
- MUST implement CSV converter (.csv) that converts to a Markdown table
- MUST implement JSON converter (.json) that converts to a fenced code block with optional key-value extraction
- MUST implement XML converter (.xml) that extracts text content
- Registry MUST return a clear error when no converter accepts the input
</requirements>

## Subtasks

- [x] 5.1 Create `internal/convert/` package with Converter interface and Registry type
- [x] 5.2 Implement Registry with Register, Match, and Convert methods
- [x] 5.3 Implement plain text/markdown pass-through converter
- [x] 5.4 Implement CSV-to-Markdown-table converter
- [x] 5.5 Implement JSON-to-fenced-code-block converter
- [x] 5.6 Implement XML text extraction converter
- [x] 5.7 Write unit tests for registry matching and each converter

## Implementation Details

Create `internal/convert/registry.go` for the registry, and one file per converter: `text.go`, `csv.go`, `json.go`, `xml.go`. Each converter is registered via an `init()` call or explicit registration in the registry constructor.

Reference TechSpec "Core Interfaces" section and ADR-003 for the converter pattern. Follow markitdown's architecture: iterate converters, first `Accepts()` match wins.

### Relevant Files

- `internal/models/kb_models.go` (task_02) — Converter interface, ConvertInput, ConvertResult types

### Dependent Files

- `internal/convert/html.go` (task_06) — follows same pattern
- `internal/convert/pdf.go` (task_07) — follows same pattern
- `internal/convert/docx.go`, `pptx.go`, `xlsx.go` (task_08) — follows same pattern
- `internal/convert/epub.go`, `ocr.go` (task_09) — follows same pattern
- `internal/ingest/` (task_12) — uses the Registry to convert files

### Related ADRs

- [ADR-003: Native Go Document Conversion with Converter Registry](../adrs/adr-003.md) — defines the converter registry pattern

## Deliverables

- `internal/convert/registry.go` — Registry type with Register, Match, Convert
- `internal/convert/text.go` — plain text/markdown converter
- `internal/convert/csv.go` — CSV-to-Markdown-table converter
- `internal/convert/json.go` — JSON-to-fenced-code-block converter
- `internal/convert/xml.go` — XML text extraction converter
- Tests for all files
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests

- Unit tests:
  - [x] Registry.Match returns correct converter for .txt, .csv, .json, .xml extensions
  - [x] Registry.Match returns nil for unregistered extension
  - [x] Registry.Convert returns error for unsupported extension
  - [x] Text converter passes through .txt content with first-line title extraction
  - [x] Text converter passes through .md content preserving markdown formatting
  - [x] CSV converter produces valid Markdown table from 3-column CSV with headers
  - [x] CSV converter handles empty CSV (headers only) gracefully
  - [x] CSV converter handles CSV with special characters (pipes, quotes)
  - [x] JSON converter wraps content in fenced code block with `json` language tag
  - [x] JSON converter extracts title from top-level "title" or "name" field if present
  - [x] XML converter extracts text nodes stripping tags
  - [x] XML converter handles empty/malformed XML with error
- Test coverage target: >=80%
- All tests must pass

## Success Criteria

- All tests passing
- Test coverage >=80%
- Registry correctly routes files to converters by extension
- Each simple converter produces clean, readable Markdown
- `make lint` reports zero findings
