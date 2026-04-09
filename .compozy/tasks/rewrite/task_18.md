---
status: completed
title: Wire search and index-vault commands
type: backend
complexity: medium
dependencies:
  - task_13
  - task_17
---

# Task 18: Wire search and index-vault commands

> **Reference implementation:** `~/dev/projects/kodebase` (original TypeScript kodebase) — use this repository as the behavioral and structural source when implementing this Go port.

## Overview

Wire the search and index-vault CLI commands, connecting the QMD shell client (task_17) with the vault query resolver (task_13). The search command provides hybrid/lexical/vector search over vault contents; the index command creates or updates QMD collections for a vault topic.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE IMPLEMENTATION — `~/dev/projects/kodebase` (original TypeScript kodebase) is the behavioral and structural source for this Go rewrite; align behavior and structure when porting
- REFERENCE TECHSPEC Phases 8.2 and 8.3 for CLI flags and behavior — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST implement search command with flags: --lex, --vec, --limit, --min-score, --full, --all, --collection
- MUST default search to hybrid mode when neither --lex nor --vec is specified
- MUST implement index command with flags: --embed, --force-embed, --context
- MUST resolve vault/topic before invoking QMD client
- MUST display search results with source path, score, and preview
- MUST handle ErrQMDUnavailable with a user-friendly message suggesting QMD installation
</requirements>

## Subtasks
- [x] 18.1 Implement search command with search mode selection (hybrid/lex/vec)
- [x] 18.2 Wire search CLI flags and result display formatting
- [x] 18.3 Implement index command with collection creation/update
- [x] 18.4 Wire index CLI flags
- [x] 18.5 Handle QMD-not-installed scenario with helpful error message

## Implementation Details

Modify `internal/cli/search.go` and `internal/cli/index.go` to replace stubs with real implementations. Reference:
- `~/dev/projects/kodebase/packages/cli/src/commands/search.ts` (251 lines)
- `~/dev/projects/kodebase/packages/cli/src/commands/index-vault.ts` (135 lines)

### Relevant Files
- `~/dev/projects/kodebase/packages/cli/src/commands/search.ts` — TypeScript source (251 lines)
- `~/dev/projects/kodebase/packages/cli/src/commands/index-vault.ts` — TypeScript source (135 lines)
- `internal/cli/search.go` — stub from task_01, replace with real command
- `internal/cli/index.go` — stub from task_01, replace with real command
- `internal/qmd/client.go` — QMDClient (task_17)
- `internal/vault/query.go` — vault/topic resolution (task_13)

### Dependent Files
- `internal/cli/root.go` — commands already registered from task_01

## Deliverables
- Updated `internal/cli/search.go` with real search implementation
- Updated `internal/cli/index.go` with real index implementation
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests
- Unit tests:
  - [x] Search command defaults to hybrid mode when no flags specified
  - [x] Search with --lex flag selects lexical mode
  - [x] Search with --vec flag selects vector mode
  - [x] Search results display path, score, and preview text
  - [x] Index command resolves vault path before calling QMD
  - [x] QMD unavailable produces user-friendly error with installation hint
  - [x] Search with --limit flag passes limit to QMD client
- Integration tests:
  - [x] (Skip if QMD not installed) Search against an indexed vault returns results
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- `make verify` passes
- `kodebase search --help` and `kodebase index --help` show correct flags
- QMD-not-installed case handled gracefully
