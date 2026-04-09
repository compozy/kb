---
status: completed
title: Set up tree-sitter infrastructure
type: infra
complexity: low
dependencies:
  - task_01
---

# Task 04: Set up tree-sitter infrastructure

> **Reference implementation:** `~/dev/projects/kodebase` (original TypeScript kodebase) — use this repository as the behavioral and structural source when implementing this Go port.

## Overview

Install the Go tree-sitter bindings and grammar packages (Go, TypeScript, JavaScript) as project dependencies. This provides the parsing foundation for the language adapters in tasks 05 and 06.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE IMPLEMENTATION — `~/dev/projects/kodebase` (original TypeScript kodebase) is the behavioral and structural source for this Go rewrite; align behavior and structure when porting
- REFERENCE TECHSPEC Phase 3.1 for dependency list — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST install `github.com/tree-sitter/go-tree-sitter` (Go bindings)
- MUST install `github.com/tree-sitter/tree-sitter-go/bindings/go` (Go grammar)
- MUST install `github.com/tree-sitter/tree-sitter-typescript/bindings/go` (TypeScript grammar)
- MUST install `github.com/tree-sitter/tree-sitter-javascript/bindings/go` (JavaScript grammar)
- MUST verify `go build ./...` compiles cleanly after installation
- MUST verify tree-sitter can parse a trivial Go source string in a smoke test
</requirements>

## Subtasks
- [x] 4.1 Install all four tree-sitter packages via `go get`
- [x] 4.2 Run `go mod tidy` to clean up dependencies
- [x] 4.3 Write smoke test that parses a trivial Go source string
- [x] 4.4 Verify `go build ./...` compiles cleanly

## Implementation Details

Use `go get` to install all four packages. Write a minimal smoke test in `internal/adapter/treesitter_test.go` that creates a parser, sets the Go language, and parses a simple `package main` source to verify the bindings work.

Reference TechSpec Phase 3.1 for exact package paths.

### Relevant Files
- `go.mod` — will have new tree-sitter dependencies added

### Dependent Files
- `internal/adapter/go_adapter.go` — task_05 will use tree-sitter-go
- `internal/adapter/ts_adapter.go` — task_06 will use tree-sitter-typescript and tree-sitter-javascript

## Deliverables
- Updated `go.mod` with tree-sitter dependencies
- `internal/adapter/treesitter_test.go` smoke test
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests
- Unit tests:
  - [x] Tree-sitter parser initializes without error for Go grammar
  - [x] Tree-sitter parser initializes without error for TypeScript grammar
  - [x] Tree-sitter parser initializes without error for JavaScript grammar
  - [x] Parsing `package main` with Go grammar produces a valid root node
  - [x] Parsing `const x = 1;` with TypeScript grammar produces a valid root node
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- `go build ./...` compiles with no errors
- `make verify` passes
