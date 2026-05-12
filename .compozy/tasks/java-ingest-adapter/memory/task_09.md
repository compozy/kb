# Task Memory: task_09.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot
- Implement deterministic handling for ambiguous Java import targets (duplicate simple-name imports and static import conflicts) without emitting misleading `calls` relations.
- Add diagnostics coverage and regression tests (unit + integration) for ambiguity scenarios while preserving stable behavior for non-ambiguous paths.

## Important Decisions
- Added an explicit ambiguity index for imported class qualifiers in `buildJavaImportLookupIndexes` instead of last-write-wins map overwrite.
- Defined deterministic ambiguity precedence for deep call resolution:
  - Unqualified call with multiple static import targets -> unresolved with `ambiguous-static-call-target`.
  - Qualified call with ambiguous imported class qualifier -> unresolved with `ambiguous-import-class`.
- Applied the same ambiguity guard in fallback candidate selection so syntactic handoff does not recreate misleading `calls` edges.

## Learnings
- Existing static import handling already avoided emitting fallback edges when multiple targets existed, but deep resolution still preferred owner-method fallback; this needed an explicit ambiguity short-circuit.
- Duplicate explicit class imports with the same simple name were previously collapsed by map overwrite, which could create deterministic-but-misleading semantic edges.

## Files / Surfaces
- `internal/adapter/java_adapter.go`
- `internal/adapter/java_adapter_test.go`
- `internal/adapter/java_adapter_integration_test.go`
- Validation commands executed:
  - `go test ./internal/adapter -run "JavaAdapter|ResolveJava|JavaNested|JavaHelper|SortJavaDiagnostics"`
  - `go test -tags integration ./internal/adapter -run "JavaAdapter"`
  - `go test -tags integration ./internal/adapter -cover` (`80.4%`)
  - `make verify`

## Errors / Corrections
- No implementation blockers after code changes; all validation commands passed on first verification cycle.

## Ready for Next Run
- Task tracking files still need to be the source of truth for completion state.
- If follow-up tasks tighten diagnostics governance, reuse ambiguity reasons `ambiguous-import-class` and `ambiguous-static-call-target` as stable diagnostic signals.
