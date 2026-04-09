---
status: completed
title: Implement vault path and text utilities
type: backend
complexity: low
dependencies:
  - task_01
---

# Task 09: Implement vault path and text utilities

> **Reference implementation:** `~/dev/projects/kodebase` (original TypeScript kodebase) — use this repository as the behavioral and structural source when implementing this Go port.

## Overview

Port `path-utils.ts` and `text-utils.ts` to Go — path manipulation, vault path derivation, file/symbol ID generation, comment extraction, and text normalization. These utilities are shared by the renderer, writer, and reader.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE IMPLEMENTATION — `~/dev/projects/kodebase` (original TypeScript kodebase) is the behavioral and structural source for this Go rewrite; align behavior and structure when porting
- REFERENCE TECHSPEC Phases 5.1 and 5.2 for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST implement file ID generation from relative path (deterministic, URL-safe)
- MUST implement symbol ID generation from file path + symbol name
- MUST implement POSIX path normalization (forward slashes, no trailing slash)
- MUST implement relative-path-inside check (is path within root?)
- MUST implement vault path derivation (source path -> vault document path)
- MUST implement leading comment extraction from source text
- MUST implement quote stripping from strings
- All path functions MUST handle edge cases: empty strings, absolute paths, paths with spaces
</requirements>

## Subtasks
- [x] 9.1 Implement file and symbol ID generation functions
- [x] 9.2 Implement POSIX path normalization and relative-path-inside check
- [x] 9.3 Implement vault path derivation (source path to vault document path)
- [x] 9.4 Implement text utilities: leading comment extraction and quote stripping

## Implementation Details

Create `internal/vault/pathutils.go`, `internal/vault/pathutils_test.go`, `internal/vault/textutils.go`, and `internal/vault/textutils_test.go`. Reference:
- `~/dev/projects/kodebase/packages/cli/src/knowledge-base/path-utils.ts` (106 lines)
- `~/dev/projects/kodebase/packages/cli/src/knowledge-base/text-utils.ts` (36 lines)

### Relevant Files
- `~/dev/projects/kodebase/packages/cli/src/knowledge-base/path-utils.ts` — TypeScript source (106 lines)
- `~/dev/projects/kodebase/packages/cli/src/knowledge-base/text-utils.ts` — TypeScript source (36 lines)

### Dependent Files
- `internal/vault/render.go` — uses path derivation for document paths and ID generation
- `internal/vault/writer.go` — uses path derivation for file output paths
- `internal/vault/reader.go` — uses path derivation for reading vault documents

## Deliverables
- `internal/vault/pathutils.go` with path utility functions
- `internal/vault/pathutils_test.go` with tests
- `internal/vault/textutils.go` with text utility functions
- `internal/vault/textutils_test.go` with tests
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests
- Unit tests:
  - [x] File ID from `src/main.go` is deterministic and URL-safe
  - [x] Symbol ID from file path + symbol name is unique and deterministic
  - [x] POSIX normalization converts backslashes to forward slashes
  - [x] POSIX normalization removes trailing slashes
  - [x] Relative-path-inside returns true for child path, false for path outside root
  - [x] Vault path derivation maps source path to expected vault document path
  - [x] Leading comment extraction returns first comment block from Go source
  - [x] Leading comment extraction returns first comment block from TS source
  - [x] Quote stripping removes surrounding single and double quotes
  - [x] Edge case: empty string inputs don't panic
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- `make verify` passes
- ID generation is deterministic (same input always produces same output)
