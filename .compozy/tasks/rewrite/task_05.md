---
status: completed
title: Implement Go language adapter
type: backend
complexity: high
dependencies:
  - task_02
  - task_04
---

# Task 05: Implement Go language adapter

> **Reference implementation:** `~/dev/projects/kodebase` (original TypeScript kodebase) — use this repository as the behavioral and structural source when implementing this Go port.

## Overview

Port `go-adapter.ts` to Go — parse Go source files using tree-sitter-go to extract symbols (functions, methods, types, interfaces, structs), imports, call relations, doc comments, and compute cyclomatic complexity. This adapter powers the analysis of Go source code in the pipeline.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE IMPLEMENTATION — `~/dev/projects/kodebase` (original TypeScript kodebase) is the behavioral and structural source for this Go rewrite; align behavior and structure when porting
- REFERENCE TECHSPEC Phase 3.2 for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST implement `LanguageAdapter` interface from `internal/models`
- MUST `Supports()` return true only for `LangGo`
- MUST extract package declarations from each file
- MUST extract function and method declarations with: name, signature, doc comment, exported flag, start/end lines
- MUST extract type declarations: struct, interface, type alias
- MUST extract import statements and build import relations
- MUST extract function call expressions and build call relations
- MUST compute cyclomatic complexity per function (count decision points: if, switch, case, for, &&, ||)
- MUST detect exported symbols via capitalized first letter
- MUST produce diagnostics for parse errors
- MUST use compile-time interface verification: `var _ models.LanguageAdapter = (*GoAdapter)(nil)`
</requirements>

## Subtasks
- [x] 5.1 Implement GoAdapter struct with Supports() and ParseFiles() methods
- [x] 5.2 Implement symbol extraction (functions, methods, types, interfaces, structs)
- [x] 5.3 Implement import extraction and import relation building
- [x] 5.4 Implement call relation extraction from function bodies
- [x] 5.5 Implement cyclomatic complexity computation
- [x] 5.6 Implement doc comment extraction and export detection

## Implementation Details

Create `internal/adapter/go_adapter.go` and `internal/adapter/go_adapter_test.go`. Reference `~/dev/projects/kodebase/packages/cli/src/knowledge-base/adapters/go-adapter.ts` (317 lines) for the parsing logic and tree-sitter query patterns.

Tree-sitter AST traversal extracts nodes by type (function_declaration, method_declaration, type_declaration, import_declaration, call_expression). Cyclomatic complexity counts branching nodes within function bodies.

### Relevant Files
- `~/dev/projects/kodebase/packages/cli/src/knowledge-base/adapters/go-adapter.ts` — TypeScript source (317 lines)
- `internal/models/models.go` — LanguageAdapter interface, ParsedFile, GraphFile, SymbolNode, RelationEdge
- `internal/adapter/treesitter_test.go` — smoke test from task_04 confirms bindings work

### Dependent Files
- `internal/cli/generate.go` — will select GoAdapter for .go files
- `internal/graph/normalize.go` — receives ParsedFile output from this adapter

## Deliverables
- `internal/adapter/go_adapter.go` implementing LanguageAdapter
- `internal/adapter/go_adapter_test.go` with comprehensive tests
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests
- Unit tests:
  - [x] Parse simple Go file with one function — returns correct GraphFile and SymbolNode
  - [x] Parse file with exported and unexported functions — exported flag set correctly
  - [x] Parse file with struct and interface declarations — correct symbol kinds
  - [x] Parse file with method (receiver function) — method name and receiver extracted
  - [x] Extract imports from `import (...)` block — correct RelationEdge entries
  - [x] Extract function calls — call relations link caller to callee
  - [x] Cyclomatic complexity of linear function equals 1
  - [x] Cyclomatic complexity of function with if/else/for equals expected count
  - [x] Parse error in source produces diagnostic, not a panic
  - [x] Doc comments above functions are extracted correctly
- Integration tests:
  - [x] Parse a multi-file Go package and verify cross-file relations
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- `make verify` passes
- Compile-time interface verification present
- All symbol kinds from TechSpec are handled
