---
status: completed
title: Update documentation and final verification
type: docs
complexity: low
dependencies:
  - task_19
---

# Task 20: Update documentation and final verification

> **Reference implementation:** `~/dev/projects/kodebase` (original TypeScript kodebase) — use this repository as the behavioral and structural source when implementing this Go port.

## Overview

Update project documentation (AGENTS.md, CLAUDE.md, config.example.toml) to reflect the completed kodebase implementation, then run the full verification pipeline to ensure everything passes cleanly.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE IMPLEMENTATION — `~/dev/projects/kodebase` (original TypeScript kodebase) is the behavioral and structural source for this Go rewrite; align behavior and structure when porting
- REFERENCE TECHSPEC Phase 9.2 and 9.3 for documentation requirements — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST update AGENTS.md with: project overview, complete package layout, all CLI commands
- MUST update CLAUDE.md with: any new conventions or patterns discovered during implementation
- MUST update config.example.toml with any kodebase-specific configuration options
- MUST run `make verify` and fix any remaining issues until it passes cleanly
- MUST ensure zero golangci-lint warnings
- MUST ensure all tests pass with -race flag
</requirements>

## Subtasks
- [x] 20.1 Update AGENTS.md with complete project documentation
- [x] 20.2 Update CLAUDE.md with any new conventions
- [x] 20.3 Update config.example.toml with kodebase configuration options
- [x] 20.4 Run `make verify` and fix any issues
- [x] 20.5 Final cleanup: remove any obsolete files or references

## Implementation Details

Update existing documentation files to accurately describe the completed project. Run `make verify` (fmt → lint → test → build → boundaries) and iterate until clean.

Reference TechSpec Phase 9.2 for documentation content requirements.

### Relevant Files
- `AGENTS.md` — project overview and agent instructions
- `CLAUDE.md` — project conventions and build commands
- `config.example.toml` — example configuration template
- `Makefile` — verify target
- `magefile.go` — build orchestration

### Dependent Files
- All source files — verify gate checks everything

## Deliverables
- Updated `AGENTS.md` with complete kodebase documentation
- Updated `CLAUDE.md` with any new conventions
- Updated `config.example.toml`
- Clean `make verify` output with zero issues
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests
- Unit tests:
  - [x] `make fmt` produces no changes (code already formatted)
  - [x] `make lint` reports zero issues
  - [x] `make test` passes all unit tests with -race flag
  - [x] `make build` produces bin/kodebase binary
  - [x] Package boundary check passes (no forbidden imports)
- Integration tests:
  - [x] `make test-integration` passes all integration tests
  - [x] `make verify` completes successfully end-to-end
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- `make verify` passes with zero warnings and zero errors
- Documentation accurately reflects the implemented project
- Binary builds and all commands respond to --help
