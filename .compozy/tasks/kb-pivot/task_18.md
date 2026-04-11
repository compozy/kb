---
status: completed
title: Integration tests and Makefile update
type: chore
complexity: high
dependencies:
  - task_16
  - task_17
---

# Task 18: Integration tests and Makefile update

## Overview

Write end-to-end integration tests that exercise the full CLI and update the Makefile/Mage build targets for the renamed binary. This is the verification gate — the project must pass `make verify` (fmt → lint → test → build) before the pivot is considered complete.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST update Makefile `build` target to compile `cmd/kb/` and produce `kb` binary
- MUST update any Mage build logic referencing the old binary path
- MUST write integration test: scaffold topic → ingest file → verify vault output
- MUST write integration test: scaffold topic → ingest codebase → run inspect → verify output
- MUST write integration test: scaffold topic → ingest multiple files → lint → verify issues detected
- MUST ensure `make verify` passes: fmt → lint → test → build
- MUST update `config.example.toml` with all new configuration sections
- MUST update CLAUDE.md with the new CLI surface, package layout, and conventions
</requirements>

## Subtasks

- [x] 18.1 Update Makefile build target for `cmd/kb/`
- [x] 18.2 Update Mage build configuration if needed
- [x] 18.3 Write ingest end-to-end integration test (scaffold → ingest file → verify vault)
- [x] 18.4 Write codebase ingest integration test (scaffold → ingest codebase → inspect → verify)
- [x] 18.5 Write lint integration test (scaffold → ingest → lint → verify issues)
- [x] 18.6 Update CLAUDE.md with new CLI surface and package layout
- [x] 18.7 Run `make verify` and fix any issues until it passes clean

## Implementation Details

Integration tests go in the package they exercise, behind `//go:build integration` tags. Use `t.TempDir()` for vault isolation and `os.Setenv` for config overrides.

Update CLAUDE.md to reflect:
- New CLI surface (`kb topic`, `kb ingest`, `kb lint`, etc.)
- New package layout (internal/convert/, internal/ingest/, internal/lint/, internal/firecrawl/, internal/youtube/, internal/topic/, internal/frontmatter/)
- Updated conventions for the KB pivot

### Relevant Files

- `Makefile` — build targets to update
- `magefiles/` — Mage build configuration (if exists)
- `CLAUDE.md` — project documentation to update
- `config.example.toml` — add firecrawl and openrouter sections
- `internal/cli/` — integration tests for full command flows
- `internal/generate/generate_integration_test.go` — existing integration test pattern to follow

### Dependent Files

- All packages — this is the final verification gate

## Deliverables

- Updated `Makefile` with `cmd/kb/` build target
- Integration test files (co-located with packages, `//go:build integration`)
- Updated `CLAUDE.md` with new project documentation
- Updated `config.example.toml`
- `make verify` passing clean **(REQUIRED)**

## Tests

- Integration tests:
  - [x] Scaffold topic → ingest .txt file → verify file in raw/articles/ with correct frontmatter
  - [x] Scaffold topic → ingest .csv file → verify Markdown table in raw/articles/
  - [x] Scaffold topic → ingest codebase (Go fixture) → verify raw/codebase/files/ and raw/codebase/symbols/ populated
  - [x] Scaffold topic → ingest codebase → run inspect smells → verify structured output
  - [x] Scaffold topic → create vault with dead wikilink → lint → verify dead-link issue detected
  - [x] Scaffold topic → create vault with orphan article → lint → verify orphan issue detected
  - [x] Full `make verify` pipeline passes (fmt → lint → test → build)
- Test coverage target: >=80%
- All tests must pass

## Success Criteria

- All tests passing
- Test coverage >=80%
- `make verify` passes clean with zero lint findings
- `kb` binary compiles and all commands are accessible
- CLAUDE.md accurately reflects the current project state
- `config.example.toml` documents all configuration keys
