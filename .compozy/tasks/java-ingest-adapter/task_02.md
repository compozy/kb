---
status: completed
title: Integrate Tree-sitter Java language binding
type: backend
complexity: medium
dependencies: []
---

# Task 02: Integrate Tree-sitter Java language binding

## Overview
Add Java grammar binding support to the adapter tree-sitter language registry so a Java parser can be created by the new adapter. This task provides parser infrastructure only and does not implement Java domain extraction yet.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- The codebase MUST add the official Tree-sitter Java Go binding dependency via `go get`.
- `internal/adapter/treesitter.go` MUST expose a `javaLanguage()` helper matching existing language helper patterns.
- Parser initialization tests MUST validate Java language loading and trivial Java source parsing.
- Existing tree-sitter language initialization tests MUST continue passing for all current languages.
</requirements>

## Subtasks
- [x] 2.1 Add Tree-sitter Java dependency to module manifests using repository dependency workflow.
- [x] 2.2 Register Java language loader in adapter tree-sitter helpers.
- [x] 2.3 Extend `treesitter_test` language initialization matrix with Java.
- [x] 2.4 Extend trivial parser test matrix with Java source fixture.
- [x] 2.5 Run targeted adapter package tests for tree-sitter helpers.

## Implementation Details
Follow the existing language loader pattern in `internal/adapter/treesitter.go` for Go/TS/JS/Rust and add Java in the same style. Ensure tests validate language ABI and parser correctness for basic Java source.

### Relevant Files
- `go.mod` — direct dependency declaration for tree-sitter Java binding.
- `go.sum` — checksum updates from dependency resolution.
- `internal/adapter/treesitter.go` — language helper functions and parser setup.
- `internal/adapter/treesitter_test.go` — parser/language initialization test matrix.

### Dependent Files
- `internal/adapter/java_adapter.go` — consumes `javaLanguage()` for parser creation.
- `internal/adapter/java_adapter_test.go` — relies on Java parser availability.

### Related ADRs
- [ADR-001: Adopt a balanced MVP strategy for Java codebase ingest](../adrs/adr-001.md) — Java support must be native within current adapter architecture.

## Deliverables
- Tree-sitter Java dependency added and resolved in module files.
- `javaLanguage()` helper added to adapter tree-sitter layer.
- Updated tree-sitter tests covering Java initialization and parse sanity.
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for parser initialization matrix **(REQUIRED)**

## Tests
- Unit tests:
  - [x] Java language helper returns non-nil language with valid ABI version.
  - [x] `newParser(javaLanguage())` creates parser without errors.
  - [x] Trivial Java source parses without `root.HasError()`.
  - [x] Existing nil-language parser rejection behavior remains unchanged.
- Integration tests:
  - [x] Full `internal/adapter` tree-sitter test suite passes with Java included.
  - [x] Existing language initialization matrix remains green after dependency addition.
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- Java grammar is loadable through the shared tree-sitter helper layer
- Parser sanity tests prove Java parse capability for downstream adapter work
