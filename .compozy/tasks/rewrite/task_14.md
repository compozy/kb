---
status: completed
title: Implement output formatter
type: backend
complexity: low
dependencies:
  - task_01
---

# Task 14: Implement output formatter

> **Reference implementation:** `~/dev/projects/kodebase` (original TypeScript kodebase) — use this repository as the behavioral and structural source when implementing this Go port.

## Overview

Port `output-formatter.ts` to Go — format tabular data as ASCII table, JSON, or TSV for CLI output. Used by all inspect subcommands to display metrics, symbols, and analysis results.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE IMPLEMENTATION — `~/dev/projects/kodebase` (original TypeScript kodebase) is the behavioral and structural source for this Go rewrite; align behavior and structure when porting
- REFERENCE TECHSPEC Phase 7.3 for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST support three output formats: table (ASCII), json, tsv
- MUST accept column headers and row data as generic tabular input
- MUST align ASCII table columns with padding
- MUST produce valid JSON array-of-objects for json format
- MUST produce tab-separated values with header row for tsv format
- MUST handle empty data gracefully (no panic, appropriate empty output)
- SHOULD truncate long cell values in table format with ellipsis
</requirements>

## Subtasks
- [x] 14.1 Define OutputFormat type and FormatOutput function signature
- [x] 14.2 Implement ASCII table formatter with column alignment
- [x] 14.3 Implement JSON formatter (array of objects)
- [x] 14.4 Implement TSV formatter with header row

## Implementation Details

Create `internal/output/formatter.go` and `internal/output/formatter_test.go`. Reference `~/dev/projects/kodebase/packages/cli/src/knowledge-base/output-formatter.ts` (84 lines).

### Relevant Files
- `~/dev/projects/kodebase/packages/cli/src/knowledge-base/output-formatter.ts` — TypeScript source (84 lines)

### Dependent Files
- `internal/cli/inspect_smells.go` — formats smell output
- `internal/cli/inspect_deadcode.go` — formats dead code output
- `internal/cli/inspect_complexity.go` — formats complexity output
- `internal/cli/inspect_blastradius.go` — formats blast radius output
- `internal/cli/inspect_coupling.go` — formats coupling output
- `internal/cli/inspect_symbol.go` — formats symbol details
- `internal/cli/inspect_file.go` — formats file details

## Deliverables
- `internal/output/formatter.go` with FormatOutput function
- `internal/output/formatter_test.go` with comprehensive tests
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests
- Unit tests:
  - [x] Table format aligns columns correctly with 3 columns and 5 rows
  - [x] Table format handles variable-width columns (short and long values)
  - [x] JSON format produces valid JSON with correct field names from headers
  - [x] TSV format produces tab-separated values with header row
  - [x] Empty data returns appropriate output (empty table, empty JSON array)
  - [x] Single row data formats correctly in all three formats
  - [x] Special characters in cell values are handled (quotes in JSON, tabs in TSV)
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- `make verify` passes
- JSON output is valid (parseable by encoding/json)
