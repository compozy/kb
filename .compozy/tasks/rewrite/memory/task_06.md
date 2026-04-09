# Task Memory: task_06.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot
- Implement Task 06 end-to-end on the current branch by adding the TypeScript/JavaScript adapter under `internal/adapter`, covering it with unit and integration tests, and validating with coverage plus `make verify`.

## Important Decisions
- Use `internal/models` and `internal/adapter` as the canonical package paths for this task. `_techspec.md` still mentions `internal/kodebase/...`, but the active repository layout and prior adapter work use the shorter paths.
- Treat the existing PRD, techspec, workflow memory, and reference TypeScript adapter as the approved design baseline for this implementation task instead of creating a separate design artifact mid-run.
- Reuse the existing helper patterns from `go_adapter.go` for IDs, comment normalization, and import-resolution logic where practical so Task 06 stays scoped to the adapter rather than creating unrelated utility packages early.
- Resolve CommonJS identifier requires as both default and namespace bindings so `helper()` and `dep.work()` can both link syntactically without needing semantic module analysis.
- Resolve `export { x } from "./file"` and `export * from "./file"` against the parsed export maps so downstream tasks can treat barrel files as normal symbol-exporting files.

## Learnings
- The current branch already contains the graph/model surface Task 06 needs from Task 05, so this task can stay focused on adapter behavior and tests.
- Tree-sitter TypeScript/TSX exposes top-level exports as `export_statement` wrapping nodes such as `function_declaration`, `class_declaration`, `interface_declaration`, `type_alias_declaration`, and `lexical_declaration`.
- Tree-sitter JavaScript exposes CommonJS requires as `call_expression` under variable declarators, and `module.exports = ...` as an `assignment_expression` inside an `expression_statement`.
- Tree-sitter TSX requires a separate `LanguageTSX()` loader even though JSX still produces a `program` root node like TypeScript and JavaScript in the Go bindings.

## Files / Surfaces
- `internal/adapter/treesitter.go`
- `internal/adapter/treesitter_test.go`
- `internal/adapter/ts_adapter.go`
- `internal/adapter/ts_adapter_test.go`
- `internal/adapter/ts_adapter_integration_test.go`
- `.compozy/tasks/rewrite/task_06.md`
- `.compozy/tasks/rewrite/_tasks.md`

## Errors / Corrections
- Corrected `formatTSReturnType` to trim the post-colon whitespace so function signatures render as `: number` instead of `:  number`.
- Corrected re-export relation emission so `export * from "./dep"` produces concrete `exports` edges to the underlying target symbols.

## Ready for Next Run
- Verification complete:
  - `go test ./internal/models ./internal/adapter`
  - `go test -tags integration ./internal/adapter`
  - `go test -coverprofile=/tmp/kodebase-ts-adapter.cover ./internal/adapter && go tool cover -func=/tmp/kodebase-ts-adapter.cover` => 82.0%
  - `make verify`
- Task tracking files were updated locally and left unstaged by design.
- Local code commit created: `bf9bca2` (`feat: implement typescript adapter`).
