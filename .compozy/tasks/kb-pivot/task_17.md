---
status: pending
title: Implement CLI lint command and adapt existing commands
type: backend
complexity: medium
dependencies:
  - task_13
  - task_14
  - task_15
---

# Task 17: Implement CLI lint command and adapt existing commands

## Overview

Implement the `kb lint` command wired to the lint engine, and adapt the existing `inspect`, `search`, and `index` commands to work with the new topic-based vault structure. These existing commands need their vault resolution logic updated to find data under topic directories.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST implement `kb lint [<slug>] [--format table|json|tsv] [--save]` wired to `internal/lint.Lint()`
- MUST pass `--save` flag to write report to `outputs/reports/<date>-lint.md`
- MUST adapt `inspect` subcommands to resolve vault data from `<topic>/raw/codebase/` instead of top-level vault
- MUST adapt `search` command to use topic slug for QMD collection naming
- MUST adapt `index` command to use topic slug for collection naming and vault path
- MUST ensure all adapted commands accept `--topic` flag consistently
- MUST NOT break inspect subcommand functionality — only path resolution changes
</requirements>

## Subtasks

- [ ] 17.1 Implement `kb lint` command with --format and --save flags
- [ ] 17.2 Adapt inspect subcommands vault path resolution for topic-based layout
- [ ] 17.3 Adapt search command to use topic-scoped QMD collections
- [ ] 17.4 Adapt index command to use topic-scoped paths and collection names
- [ ] 17.5 Ensure --topic flag is consistently available on all adapted commands
- [ ] 17.6 Write unit tests for lint command and adapted command flag handling

## Implementation Details

Create `internal/cli/lint.go` for the lint command. Modify the existing `internal/cli/inspect.go`, `internal/cli/search.go`, and `internal/cli/index.go` to update their vault resolution logic.

The key change for inspect commands is in `inspectSharedOptions` and `runInspectCommand()` — the vault query resolver needs to look under `<topic>/raw/codebase/` for the codebase-specific data that inspect operates on.

### Relevant Files

- `internal/lint/` (task_13) — lint engine to wire
- `internal/cli/inspect.go` — shared inspect infrastructure, inspectSharedOptions
- `internal/cli/inspect_*.go` — all inspect subcommands (smells, dead-code, complexity, etc.)
- `internal/cli/search.go` — search command
- `internal/cli/index.go` — index command
- `internal/vault/query.go` — ResolveVaultQuery() may need topic-aware resolution
- `internal/output/` — FormatOutput() used by lint command

### Dependent Files

- `internal/vault/query.go` — may need modification for topic-aware vault resolution

### Related ADRs

- [ADR-001: Topic-Centric CLI Command Taxonomy](../adrs/adr-001.md) — command structure
- [ADR-005: Native Go Vault Lint Engine](../adrs/adr-005.md) — lint as first-class command

## Deliverables

- New `internal/cli/lint.go` — lint command
- Modified `internal/cli/inspect.go` — topic-aware vault resolution
- Modified `internal/cli/search.go` — topic-scoped search
- Modified `internal/cli/index.go` — topic-scoped indexing
- Tests for all changes
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests

- Unit tests:
  - [ ] `lint` command accepts optional slug positional arg
  - [ ] `lint` accepts --format flag with table, json, tsv values
  - [ ] `lint` accepts --save flag
  - [ ] `lint` prints formatted output matching lint engine results
  - [ ] Inspect commands accept --topic flag
  - [ ] Inspect commands resolve vault from `<topic>/raw/codebase/` subtree
  - [ ] Search command accepts --topic flag and scopes to topic collection
  - [ ] Index command accepts --topic flag and uses topic-scoped path
- Test coverage target: >=80%
- All tests must pass

## Success Criteria

- All tests passing
- Test coverage >=80%
- `kb lint` produces correct output for healthy and broken vaults
- All inspect/search/index commands work with topic-based vault layout
- `make lint` reports zero findings
