---
status: completed
title: Implement frontmatter package
type: backend
complexity: medium
dependencies: []
---

# Task 01: Implement frontmatter package

## Overview

Create the `internal/frontmatter/` package that provides YAML frontmatter parsing and generation for all vault markdown files. This is a foundational package used by topic scaffolding, ingest, lint, and vault reading — every component that touches markdown files with metadata depends on it.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST parse YAML frontmatter delimited by `---` from markdown files into `map[string]any`
- MUST generate YAML frontmatter from `map[string]any` and prepend to markdown body
- MUST handle all frontmatter schema variants defined in TechSpec "Data Models" section (source, wiki, output, index types)
- MUST round-trip parse/generate without data loss for supported YAML types (strings, ints, floats, bools, string slices, dates)
- MUST return structured errors for malformed frontmatter (missing delimiters, invalid YAML)
- SHOULD preserve field ordering when generating (use a deterministic key sort)
</requirements>

## Subtasks

- [x] 1.1 Create `internal/frontmatter/` package with Parse and Generate functions
- [x] 1.2 Implement frontmatter extraction from markdown content (split on `---` delimiters)
- [x] 1.3 Implement frontmatter generation that produces `---`-delimited YAML prepended to body
- [x] 1.4 Add typed helper accessors for common fields (GetString, GetStringSlice, GetTime, GetBool)
- [x] 1.5 Write comprehensive unit tests covering all schema variants and error paths

## Implementation Details

Create new package `internal/frontmatter/` with two files: `frontmatter.go` and `frontmatter_test.go`.

The existing codebase already uses `gopkg.in/yaml.v3` as a transitive dependency — use it directly for YAML marshaling/unmarshaling.

Reference the frontmatter schemas defined in `.agents/skills/karpathy-kb/references/frontmatter-schemas.md` for the field sets to support. Reference TechSpec "Data Models" section for the frontmatter schema definition.

### Relevant Files

- `internal/vault/reader.go` — currently parses frontmatter inline via string splitting; this package extracts and replaces that logic
- `internal/vault/render.go` — generates frontmatter as part of document rendering; should migrate to use this package
- `.agents/skills/karpathy-kb/references/frontmatter-schemas.md` — canonical frontmatter field definitions

### Dependent Files

- `internal/topic/` (task_04) — will use Generate for template frontmatter
- `internal/ingest/` (task_12) — will use Generate for ingested file frontmatter
- `internal/lint/` (task_13) — will use Parse for frontmatter validation
- `internal/vault/reader.go` — can be refactored to use this package in a future pass

### Related ADRs

- [ADR-005: Native Go Vault Lint Engine](../adrs/adr-005.md) — lint depends on frontmatter parsing for format violation checks

## Deliverables

- `internal/frontmatter/frontmatter.go` — Parse, Generate, and typed accessor functions
- `internal/frontmatter/frontmatter_test.go` — comprehensive tests
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests

- Unit tests:
  - [x] Parse valid frontmatter with all supported YAML types (string, int, float, bool, string slice, date)
  - [x] Parse markdown with no frontmatter returns empty map and full body
  - [x] Parse markdown with only `---` delimiters and empty YAML returns empty map
  - [x] Parse malformed frontmatter (missing closing delimiter) returns structured error
  - [x] Parse invalid YAML content returns structured error with position info
  - [x] Generate frontmatter from map produces valid `---`-delimited YAML
  - [x] Round-trip: Parse(Generate(map, body)) == (map, body) for all schema variants
  - [x] GetString, GetStringSlice, GetTime, GetBool return correct types and zero values for missing keys
  - [x] Generate with empty map produces no frontmatter prefix
- Test coverage target: >=80%
- All tests must pass

## Success Criteria

- All tests passing
- Test coverage >=80%
- Parse and Generate handle all frontmatter schemas from the karpathy-kb skill
- `make lint` reports zero findings for the new package
