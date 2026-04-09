---
status: completed
title: Wire inspect lookup subcommands
type: backend
complexity: medium
dependencies:
  - task_15
---

# Task 16: Wire inspect lookup subcommands

> **Reference implementation:** `~/dev/projects/kodebase` (original TypeScript kodebase) — use this repository as the behavioral and structural source when implementing this Go port.

## Overview

Implement the second batch of inspect subcommands focused on lookup and navigation: symbol, file, backlinks, deps, and circular-deps. These subcommands provide detailed views of individual symbols/files and their relationship networks.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE IMPLEMENTATION — `~/dev/projects/kodebase` (original TypeScript kodebase) is the behavioral and structural source for this Go rewrite; align behavior and structure when porting
- REFERENCE TECHSPEC Phase 7.4 for subcommand definitions — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST implement `inspect symbol <name>` — show detailed symbol information (signature, metrics, relations)
- MUST implement `inspect file <path>` — show detailed file information (symbols, metrics, coupling)
- MUST implement `inspect backlinks <name>` — list all symbols/files that reference a given entity
- MUST implement `inspect deps <name>` — list all dependencies of a given symbol/file
- MUST implement `inspect circular-deps` — list all detected circular dependency cycles
- All subcommands MUST support --format flag (table|json|tsv)
- All subcommands MUST support --vault and --topic flags
- All subcommands MUST reuse shared infrastructure from task_15
</requirements>

## Subtasks
- [x] 16.1 Implement `inspect symbol` subcommand with detail view
- [x] 16.2 Implement `inspect file` subcommand with detail view
- [x] 16.3 Implement `inspect backlinks` subcommand
- [x] 16.4 Implement `inspect deps` subcommand
- [x] 16.5 Implement `inspect circular-deps` subcommand

## Implementation Details

Create: `internal/cli/inspect_symbol.go`, `inspect_file.go`, `inspect_backlinks.go`, `inspect_deps.go`, `inspect_circulardeps.go`.

Reference `~/dev/projects/kodebase/packages/cli/src/commands/inspect/` (lookup.ts 168 lines, graph.ts 101 lines).

Symbol and file lookups use the vault reader (task_13) for name resolution. Backlinks and deps traverse the relation graph stored in vault documents.

### Relevant Files
- `~/dev/projects/kodebase/packages/cli/src/commands/inspect/lookup.ts` — TypeScript source (168 lines)
- `~/dev/projects/kodebase/packages/cli/src/commands/inspect/graph.ts` — TypeScript source (101 lines)
- `internal/cli/inspect.go` — shared infrastructure from task_15
- `internal/vault/reader.go` — vault document reading (task_13)
- `internal/vault/query.go` — symbol/file lookup (task_13)
- `internal/output/formatter.go` — output formatting (task_14)

### Dependent Files
- `internal/cli/root.go` — subcommands register under inspect parent

## Deliverables
- `internal/cli/inspect_symbol.go`
- `internal/cli/inspect_file.go`
- `internal/cli/inspect_backlinks.go`
- `internal/cli/inspect_deps.go`
- `internal/cli/inspect_circulardeps.go`
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests
- Unit tests:
  - [x] Inspect symbol shows signature, metrics, and outgoing relations
  - [x] Inspect symbol with unknown name returns descriptive error
  - [x] Inspect file shows contained symbols and file-level metrics
  - [x] Inspect file with unknown path returns descriptive error
  - [x] Inspect backlinks lists all incoming references for a symbol
  - [x] Inspect deps lists all outgoing dependencies for a symbol
  - [x] Inspect circular-deps lists all detected cycles
  - [x] Inspect circular-deps with no cycles shows "no circular dependencies found"
- Integration tests:
  - [x] Run lookup subcommands against a vault generated from a fixture project
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- `make verify` passes
- All 5 subcommands registered and respond to `kodebase inspect <subcommand> --help`
- All 10 inspect subcommands (task_15 + task_16) functional
