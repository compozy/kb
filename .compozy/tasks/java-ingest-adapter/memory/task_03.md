# Task Memory: task_03.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot
- Deliver Java adapter MVP parse pipeline with deterministic output, structured parse diagnostics, and baseline cross-file relation emission.

## Important Decisions
- Implemented Java MVP adapter with tree-sitter parse validation plus syntactic extraction for package/import metadata and call-target hints.
- Kept parse failures non-fatal per file (`JAVA_PARSE_ERROR`, `StageParse`) including nil-tree and nil-root defensive handling.
- Added baseline cross-file resolution as syntactic mapping (`import` -> class symbol reference, `Class.method()` -> method call relation) without deep classpath resolver (deferred to later task).

## Learnings
- Java package and import extraction is stable for MVP via source-pattern parsing while symbol/method boundaries are still anchored to tree-sitter declaration nodes.
- Existing adapter package coverage remains below 80% globally even with strong Java coverage; this is a package-level baseline issue, not specific to this task.

## Files / Surfaces
- `internal/adapter/java_adapter.go`
- `internal/adapter/java_adapter_test.go`
- `internal/adapter/java_adapter_integration_test.go`
- `internal/adapter/go_adapter_test.go`
- `internal/adapter/ts_adapter_test.go`
- `internal/adapter/rust_adapter_test.go`

## Errors / Corrections
- Initial shell filtering used `rg` in shell pipeline, but environment lacked shell `rg`; switched to direct `go tool cover -func` output inspection.

## Ready for Next Run
- Task deliverables for Java adapter MVP are implemented and verified with `make verify`.
- Next task can register `JavaAdapter` into generate runner and then deepen relation resolution in task 05.
