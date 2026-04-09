---
status: completed
title: Implement QMD shell client
type: backend
complexity: medium
dependencies:
  - task_01
---

# Task 17: Implement QMD shell client

> **Reference implementation:** `~/dev/projects/kodebase` (original TypeScript kodebase) — use this repository as the behavioral and structural source when implementing this Go port.

## Overview

Port `qmd-client.ts` to Go — execute QMD CLI commands via `os/exec` for vault indexing and hybrid search. Provides a graceful fallback when QMD is not installed on the system.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE IMPLEMENTATION — `~/dev/projects/kodebase` (original TypeScript kodebase) is the behavioral and structural source for this Go rewrite; align behavior and structure when porting
- REFERENCE TECHSPEC Phase 8.1 for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST implement QMDClient struct with configurable binary path (default: "qmd")
- MUST implement Index: `qmd collection add <path>` and `qmd collection update <path>`
- MUST implement Search: `qmd search --mode hybrid|lexical|vector <query>`
- MUST parse QMD JSON output into structured Go types (SearchResult, IndexStatus, etc.)
- MUST return `ErrQMDUnavailable` when QMD binary is not found on PATH
- MUST use context.Context for command execution (cancellation, timeout)
- MUST capture stderr for error diagnostics
- SHOULD support all QMD search options: mode, limit, all, minScore, full, collection
</requirements>

## Subtasks
- [x] 17.1 Implement QMDClient struct with binary path resolution
- [x] 17.2 Implement Index method (collection add/update via os/exec)
- [x] 17.3 Implement Search method (search with mode/options via os/exec)
- [x] 17.4 Implement QMD output parsing (JSON to Go structs)
- [x] 17.5 Implement graceful ErrQMDUnavailable fallback

## Implementation Details

Create `internal/qmd/client.go` and `internal/qmd/client_test.go`. Reference `~/dev/projects/kodebase/packages/cli/src/integrations/qmd-client.ts` (290 lines).

Use `exec.CommandContext` for all QMD invocations. Check for binary existence with `exec.LookPath` before attempting execution.

### Relevant Files
- `~/dev/projects/kodebase/packages/cli/src/integrations/qmd-client.ts` — TypeScript source (290 lines)

### Dependent Files
- `internal/cli/search.go` — will call QMDClient.Search
- `internal/cli/index.go` — will call QMDClient.Index

## Deliverables
- `internal/qmd/client.go` with QMDClient struct and methods
- `internal/qmd/client_test.go` with comprehensive tests
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests
- Unit tests:
  - [x] QMDClient with non-existent binary returns ErrQMDUnavailable
  - [x] Index method constructs correct command arguments for "add" operation
  - [x] Index method constructs correct command arguments for "update" operation
  - [x] Search method constructs correct command with --mode hybrid flag
  - [x] Search method constructs correct command with --limit and --min-score flags
  - [x] Search output JSON is parsed into SearchResult structs correctly
  - [x] Context cancellation stops running QMD command
  - [x] stderr capture provides diagnostic information on command failure
- Integration tests:
  - [x] (Skip if QMD not installed) Index and search a temp vault directory
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- `make verify` passes
- Graceful degradation when QMD is not installed
