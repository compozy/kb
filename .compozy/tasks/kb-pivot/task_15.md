---
status: pending
title: Rename binary and rewrite CLI root and topic commands
type: refactor
complexity: medium
dependencies:
  - task_04
---

# Task 15: Rename binary and rewrite CLI root and topic commands

## Overview

Rename the binary from `kodebase` to `kb`, restructure the CLI root command, and implement the `topic` command group (`kb topic new`, `kb topic list`, `kb topic info`). This establishes the new command taxonomy foundation that all other CLI tasks build upon.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST rename `cmd/kodebase/` to `cmd/kb/` and update `main.go` accordingly
- MUST update the Cobra root command Use field from `kodebase` to `kb`
- MUST update the root command description to reflect the knowledge base focus
- MUST implement `kb topic new <slug> "<title>" <domain>` wired to `internal/topic.New()`
- MUST implement `kb topic list` wired to `internal/topic.List()`
- MUST implement `kb topic info <slug>` wired to `internal/topic.Info()`
- MUST add a `--vault` persistent flag to root command for specifying vault path (defaults to `.kodebase/vault/` auto-discovery)
- MUST preserve `kb version` command
</requirements>

## Subtasks

- [ ] 15.1 Rename `cmd/kodebase/` to `cmd/kb/` and update main.go
- [ ] 15.2 Update root command (Use, Short, Long descriptions)
- [ ] 15.3 Add `--vault` persistent flag with auto-discovery default
- [ ] 15.4 Implement `topic` parent command with `new`, `list`, `info` subcommands
- [ ] 15.5 Wire topic subcommands to `internal/topic/` functions
- [ ] 15.6 Write unit tests for command parsing and flag handling

## Implementation Details

Rename the directory `cmd/kodebase/` → `cmd/kb/`. Update `internal/cli/root.go` to change the root command. Create `internal/cli/topic.go` with the topic parent and its three subcommands.

The `--vault` persistent flag is added to the root command and inherited by all subcommands. It defaults to auto-discovery by walking up from CWD looking for `.kodebase/vault/` (reuse logic from `internal/vault/query.go`).

### Relevant Files

- `cmd/kodebase/main.go` — rename to `cmd/kb/main.go`
- `internal/cli/root.go` — root command definition, ExecuteContext()
- `internal/cli/version.go` — preserve as-is
- `internal/vault/query.go` — ResolveVaultQuery() for vault auto-discovery logic
- `internal/topic/` (task_04) — New, List, Info functions to wire

### Dependent Files

- `internal/cli/` (tasks 16, 17) — all CLI commands inherit from the updated root
- `Makefile` — build target needs updating (task_18)

### Related ADRs

- [ADR-001: Topic-Centric CLI Command Taxonomy](../adrs/adr-001.md) — topic-centric command structure
- [ADR-002: Rename Binary to `kb`](../adrs/adr-002.md) — binary rename decision

## Deliverables

- `cmd/kb/main.go` — renamed entry point
- Modified `internal/cli/root.go` — updated root command
- New `internal/cli/topic.go` — topic command group
- Tests for topic commands
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests

- Unit tests:
  - [ ] Root command has Use: "kb" and correct description
  - [ ] `--vault` flag is registered as persistent and accessible by subcommands
  - [ ] `topic new` requires exactly 3 positional args (slug, title, domain)
  - [ ] `topic new` returns error for missing args
  - [ ] `topic list` accepts `--vault` flag and returns formatted output
  - [ ] `topic info` requires exactly 1 positional arg (slug)
  - [ ] `version` command still works
- Test coverage target: >=80%
- All tests must pass

## Success Criteria

- All tests passing
- Test coverage >=80%
- `cmd/kb/` compiles and produces a `kb` binary
- `kb topic new`, `kb topic list`, `kb topic info` work end-to-end
- `kb version` continues to work
- `make lint` reports zero findings
