---
status: completed
title: Implement workspace scanner
type: backend
complexity: medium
dependencies:
  - task_02
---

# Task 03: Implement workspace scanner

> **Reference implementation:** `~/dev/projects/kodebase` (original TypeScript kodebase) — use this repository as the behavioral and structural source when implementing this Go port.

## Overview

Port `scan-workspace.ts` to Go — discover source files in a repository, apply gitignore rules and user include/exclude patterns, classify files by language, and return a `ScannedWorkspace`. This is the first stage of the analysis pipeline.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE IMPLEMENTATION — `~/dev/projects/kodebase` (original TypeScript kodebase) is the behavioral and structural source for this Go rewrite; align behavior and structure when porting
- REFERENCE TECHSPEC Phase 2.1 for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST load and respect `.gitignore` files from the root path
- MUST apply default exclusions: node_modules, .git, dist, vendor, .kodebase, .turbo, .next, build, coverage
- MUST support user-provided include and exclude glob patterns via ScanOptions
- MUST classify files by extension to SupportedLanguage (.ts, .tsx, .js, .jsx, .go)
- MUST return ScannedWorkspace with files grouped by language
- MUST use `github.com/sabhiram/go-gitignore` for gitignore parsing
- MUST use functional options pattern: `NewScanner(opts ...Option)`
- SHOULD skip hidden directories (.git, .hg, .svn) early in traversal
</requirements>

## Subtasks

- [x] 3.1 Install go-gitignore dependency
- [x] 3.2 Implement ScanWorkspace function with directory traversal and gitignore loading
- [x] 3.3 Implement file classification by extension to SupportedLanguage
- [x] 3.4 Implement include/exclude pattern filtering
- [x] 3.5 Return ScannedWorkspace with files grouped by language

## Implementation Details

Create `internal/scanner/scanner.go` and `internal/scanner/scanner_test.go`. Reference `~/dev/projects/kodebase/packages/cli/src/knowledge-base/scan-workspace.ts` (313 lines) for logic.

The scanner walks the directory tree using `filepath.WalkDir`, loads `.gitignore` from the root, and filters files. Files that match a supported extension are classified and grouped.

### Relevant Files

- `~/dev/projects/kodebase/packages/cli/src/knowledge-base/scan-workspace.ts` — TypeScript source (313 lines)
- `internal/models/models.go` — ScannedSourceFile, ScannedWorkspace, SupportedLanguage types

### Dependent Files

- `internal/cli/generate.go` — will call ScanWorkspace as pipeline step 1
- `internal/adapter/go_adapter.go` — receives ScannedSourceFile slice from scanner output
- `internal/adapter/ts_adapter.go` — receives ScannedSourceFile slice from scanner output

## Deliverables

- `internal/scanner/scanner.go` with ScanWorkspace implementation
- `internal/scanner/scanner_test.go` with comprehensive tests
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests

- Unit tests:
  - [x] Scan temp directory with .go, .ts, .js files returns correct file count and languages
  - [x] Files in node_modules are excluded by default
  - [x] Files in .git directory are excluded by default
  - [x] .gitignore patterns are respected (create temp .gitignore with patterns)
  - [x] User include pattern restricts results to matching files only
  - [x] User exclude pattern removes matching files from results
  - [x] Unsupported file extensions (.py, .rb, .md) are not included
  - [x] Empty directory returns empty ScannedWorkspace with no error
  - [x] Files are correctly grouped by language in ScannedWorkspace.FilesByLanguage
- Integration tests:
  - [x] Scan a real project directory structure with nested packages
- Test coverage target: >=80%
- All tests must pass

## Success Criteria

- All tests passing
- Test coverage >=80%
- `make verify` passes
- Scanner correctly handles gitignore rules, default exclusions, and user patterns
