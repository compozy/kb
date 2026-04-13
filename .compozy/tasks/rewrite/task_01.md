---
status: completed
title: Rename module and bootstrap cobra CLI skeleton
type: backend
complexity: medium
dependencies: []
---

# Task 01: Rename module and bootstrap cobra CLI skeleton

> **Reference implementation:** `~/dev/projects/kodebase` (original TypeScript kodebase) — use this repository as the behavioral and structural source when implementing this Go port.

## Overview

Rename the Go module from `github.com/compozy/kb` to `github.com/pedronauck/kodebase`, install cobra, and create the CLI skeleton with root command and stub subcommands. This establishes the project identity and command-routing foundation that all subsequent tasks build upon.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE IMPLEMENTATION — `~/dev/projects/kodebase` (original TypeScript kodebase) is the behavioral and structural source for this Go rewrite; align behavior and structure when porting
- REFERENCE TECHSPEC Phase 0.1 for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST rename module to `github.com/pedronauck/kodebase` in go.mod and all import paths
- MUST install cobra via `go get github.com/spf13/cobra@latest`
- MUST create root command in `internal/cli/root.go` with project description
- MUST create stub subcommands: generate, inspect, search, index, version
- MUST rewrite `cmd/kodebase/main.go` to call `cli.Execute()`
- MUST preserve existing config, logger, and version packages (update import paths only)
- SHOULD wire version command to use existing `internal/version` package
- All stub subcommands MUST print "not implemented" and return nil
</requirements>

## Subtasks

- [ ] 1.1 Rename module in go.mod and update all import paths across existing files
- [ ] 1.2 Install cobra dependency
- [ ] 1.3 Create root command with PersistentPreRun hook for config/logger setup
- [ ] 1.4 Create stub subcommands (generate, inspect, search, index, version)
- [ ] 1.5 Rewrite cmd/kodebase/main.go to delegate to cli.Execute()
- [ ] 1.6 Verify build compiles and subcommands respond

## Implementation Details

Rename module path in go.mod and update all existing imports in `internal/config`, `internal/logger`, `internal/version`, and `cmd/kodebase/main.go`. The existing main.go has a complex signal-handling setup — replace it with a simple `cli.Execute()` call since cobra handles its own lifecycle.

Reference TechSpec Phase 0.1 for root command structure and stub subcommand patterns.

### Relevant Files

- `go.mod` — current module is `github.com/compozy/kb`, needs renaming
- `cmd/kodebase/main.go` — existing entry point (127 lines), needs rewrite to use cobra
- `internal/version/version.go` — existing version package with ldflags, wire into version subcommand
- `internal/config/config.go` — existing config loading, import paths need updating
- `internal/logger/logger.go` — existing slog logger, import paths need updating

### Dependent Files

- `internal/config/config_test.go` — import paths need updating
- `internal/logger/logger_test.go` — import paths need updating
- `magefile.go` — may reference module path in build commands

## Deliverables

- Renamed module with all import paths updated
- Cobra root command in `internal/cli/root.go`
- Five stub subcommands in `internal/cli/`
- Simplified `cmd/kodebase/main.go`
- Unit tests with 80%+ coverage **(REQUIRED)**
- Build compiles and `kodebase version` outputs version info

## Tests

- Unit tests:
  - [ ] Root command executes without error
  - [ ] Version subcommand prints version, commit, and date
  - [ ] Generate subcommand with path arg prints "not implemented"
  - [ ] Inspect subcommand without args shows help text
  - [ ] Search subcommand prints "not implemented"
  - [ ] Index subcommand prints "not implemented"
- Integration tests:
  - [ ] Binary builds and `./bin/kodebase version` returns 0 exit code
- Test coverage target: >=80%
- All tests must pass

## Success Criteria

- All tests passing
- Test coverage >=80%
- `make verify` passes with zero warnings
- `go build ./...` compiles clean
- All existing config and logger tests still pass with new import paths
