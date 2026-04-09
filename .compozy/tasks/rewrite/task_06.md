---
status: completed
title: Implement TypeScript/JavaScript adapter
type: backend
complexity: high
dependencies:
  - task_02
  - task_04
---

# Task 06: Implement TypeScript/JavaScript adapter

> **Reference implementation:** `~/dev/projects/kodebase` (original TypeScript kodebase) — use this repository as the behavioral and structural source when implementing this Go port.

## Overview

Port `oxc-typescript-adapter.ts` to Go — parse TypeScript, TSX, JavaScript, and JSX files using tree-sitter-typescript and tree-sitter-javascript to extract top-level declarations, imports, exports, call relations, and compute cyclomatic complexity. This is the largest adapter by reference LOC (936 lines).

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE IMPLEMENTATION — `~/dev/projects/kodebase` (original TypeScript kodebase) is the behavioral and structural source for this Go rewrite; align behavior and structure when porting
- REFERENCE TECHSPEC Phase 3.3 for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST implement `LanguageAdapter` interface from `internal/models`
- MUST `Supports()` return true for LangTS, LangTSX, LangJS, LangJSX
- MUST select correct tree-sitter grammar: typescript for .ts, tsx for .tsx, javascript for .js/.jsx
- MUST extract top-level declarations: functions, classes, interfaces, type aliases, variables, export statements
- MUST build import bindings from import/require statements
- MUST resolve call targets and construct cross-file call relations
- MUST compute cyclomatic complexity per function/method
- MUST handle both named and default exports
- MUST produce diagnostics for parse errors
- MUST use compile-time interface verification: `var _ models.LanguageAdapter = (*TSAdapter)(nil)`
</requirements>

## Subtasks
- [x] 6.1 Implement TSAdapter struct with grammar selection logic per language
- [x] 6.2 Implement top-level declaration extraction (functions, classes, interfaces, type aliases, variables)
- [x] 6.3 Implement import/require binding extraction and import relation building
- [x] 6.4 Implement export handling (named exports, default exports, re-exports)
- [x] 6.5 Implement call relation resolution with cross-file target linking
- [x] 6.6 Implement cyclomatic complexity computation for TS/JS

## Implementation Details

Create `internal/adapter/ts_adapter.go` and `internal/adapter/ts_adapter_test.go`. Reference `~/dev/projects/kodebase/packages/cli/src/knowledge-base/adapters/oxc-typescript-adapter.ts` (936 lines) for comprehensive parsing logic.

The TSAdapter must handle four file extensions, each potentially using a different tree-sitter grammar. TypeScript and TSX use `tree-sitter-typescript`'s typescript and tsx languages respectively; JS and JSX use `tree-sitter-javascript`.

### Relevant Files
- `~/dev/projects/kodebase/packages/cli/src/knowledge-base/adapters/oxc-typescript-adapter.ts` — TypeScript source (936 lines)
- `~/dev/projects/kodebase/packages/cli/src/knowledge-base/adapters/typescript-adapter.ts` — alternative TS adapter (538 lines)
- `internal/models/models.go` — LanguageAdapter interface, ParsedFile, GraphFile, SymbolNode, RelationEdge

### Dependent Files
- `internal/cli/generate.go` — will select TSAdapter for .ts/.tsx/.js/.jsx files
- `internal/graph/normalize.go` — receives ParsedFile output from this adapter

## Deliverables
- `internal/adapter/ts_adapter.go` implementing LanguageAdapter
- `internal/adapter/ts_adapter_test.go` with comprehensive tests
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests
- Unit tests:
  - [x] Parse simple TypeScript file with function and interface — correct symbols extracted
  - [x] Parse TSX file with React component — JSX component detected as symbol
  - [x] Parse JavaScript file with require() imports — import relations built correctly
  - [x] Parse file with class declaration — class and method symbols extracted
  - [x] Named exports produce correct export relations
  - [x] Default export produces correct export relation
  - [x] Import statement produces correct import bindings (named, default, namespace)
  - [x] Call expressions inside function bodies produce call relations
  - [x] Cyclomatic complexity of simple function equals 1
  - [x] Cyclomatic complexity with ternary and logical operators computed correctly
  - [x] Parse error produces diagnostic, not a panic
- Integration tests:
  - [x] Parse a multi-file TS project with cross-file imports and verify relation linking
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- `make verify` passes
- Compile-time interface verification present
- Correct grammar selected for each file type
