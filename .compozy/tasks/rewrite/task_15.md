---
status: completed
title: Wire inspect analysis subcommands
type: backend
complexity: high
dependencies:
  - task_13
  - task_14
---

# Task 15: Wire inspect analysis subcommands

> **Reference implementation:** `~/dev/projects/kodebase` (original TypeScript kodebase) — use this repository as the behavioral and structural source when implementing this Go port.

## Overview

Implement the first batch of inspect subcommands focused on code analysis: smells, dead-code, complexity, blast-radius, and coupling. Each subcommand resolves the vault/topic, reads the snapshot via the vault reader, queries relevant metrics, and formats output via the output formatter.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE IMPLEMENTATION — `~/dev/projects/kodebase` (original TypeScript kodebase) is the behavioral and structural source for this Go rewrite; align behavior and structure when porting
- REFERENCE TECHSPEC Phase 7.4 for subcommand definitions — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST implement shared inspect infrastructure: vault/topic resolution + snapshot reading + format flag
- MUST implement `inspect smells` — list code smells from metrics
- MUST implement `inspect dead-code` — list dead exports from metrics
- MUST implement `inspect complexity` — list symbols sorted by cyclomatic complexity
- MUST implement `inspect blast-radius` — list symbols sorted by blast radius
- MUST implement `inspect coupling` — list files sorted by coupling/instability
- All subcommands MUST support --format flag (table|json|tsv)
- All subcommands MUST support --vault and --topic flags
- All subcommands MUST register under the `inspect` parent command
</requirements>

## Subtasks
- [x] 15.1 Implement shared inspect infrastructure (resolve vault/topic, read snapshot, format dispatch)
- [x] 15.2 Implement `inspect smells` subcommand
- [x] 15.3 Implement `inspect dead-code` subcommand
- [x] 15.4 Implement `inspect complexity` subcommand
- [x] 15.5 Implement `inspect blast-radius` subcommand
- [x] 15.6 Implement `inspect coupling` subcommand

## Implementation Details

Modify `internal/cli/inspect.go` to add shared infrastructure. Create individual subcommand files: `inspect_smells.go`, `inspect_deadcode.go`, `inspect_complexity.go`, `inspect_blastradius.go`, `inspect_coupling.go`.

Reference `~/dev/projects/kodebase/packages/cli/src/commands/inspect/` (shared.ts 184 lines, metrics.ts 267 lines).

Each subcommand follows the same pattern: resolve vault → read snapshot → extract relevant data → sort → format output.

### Relevant Files
- `~/dev/projects/kodebase/packages/cli/src/commands/inspect/shared.ts` — TypeScript source (184 lines)
- `~/dev/projects/kodebase/packages/cli/src/commands/inspect/metrics.ts` — TypeScript source (267 lines)
- `internal/vault/reader.go` — reads vault documents (task_13)
- `internal/vault/query.go` — resolves vault/topic paths (task_13)
- `internal/output/formatter.go` — formats tabular output (task_14)
- `internal/models/models.go` — SymbolMetrics, FileMetrics types

### Dependent Files
- `internal/cli/root.go` — inspect subcommands register under inspect parent

## Deliverables
- Modified `internal/cli/inspect.go` with shared infrastructure and inspect parent command
- `internal/cli/inspect_smells.go`
- `internal/cli/inspect_deadcode.go`
- `internal/cli/inspect_complexity.go`
- `internal/cli/inspect_blastradius.go`
- `internal/cli/inspect_coupling.go`
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests
- Unit tests:
  - [x] Inspect smells lists symbols with code smell flags
  - [x] Inspect dead-code lists exported symbols with zero dependents
  - [x] Inspect complexity lists symbols sorted by descending complexity
  - [x] Inspect blast-radius lists symbols sorted by descending blast radius
  - [x] Inspect coupling lists files sorted by instability
  - [x] --format json produces valid JSON output
  - [x] --format tsv produces valid TSV output
  - [x] Missing vault returns descriptive error message
- Integration tests:
  - [x] Run inspect subcommands against a vault generated from a fixture project
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- `make verify` passes
- All 5 subcommands registered and respond to `kodebase inspect <subcommand> --help`
