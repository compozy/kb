# Task Memory: task_16.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot
- Implement inspect lookup/navigation subcommands: `symbol`, `file`, `backlinks`, `deps`, and `circular-deps`, all on top of the shared inspect infrastructure from task 15.
- Finish with unit coverage >=80%, integration coverage for the new lookup commands, and a clean `make verify`.

## Important Decisions
- `backlinks` and `deps` resolve an exact file-path match first, then fall back to a single symbol-name match. This satisfies the task requirement to navigate either symbols or files while remaining compatible with the raw vault data model.
- `circular-deps` reconstructs cycle lists from file-document `imports` relations because the read-side snapshot only stores per-file circular flags.
- Symbol detail extracts the signature from the stored markdown body, while file detail derives contained symbols from the symbol documents that share the same `source_path`.

## Learnings
- The generated fixture used by the inspect integration test has an ambiguous `main` symbol lookup, so integration coverage for `deps` uses an exact file path instead of a symbol fragment.
- The `Hello` symbol in the fixture does not expose outgoing dependency rows, but its backlinks are present via file containment/export relations.

## Files / Surfaces
- `internal/cli/inspect.go`
- `internal/cli/inspect_symbol.go`
- `internal/cli/inspect_file.go`
- `internal/cli/inspect_backlinks.go`
- `internal/cli/inspect_deps.go`
- `internal/cli/inspect_circulardeps.go`
- `internal/cli/inspect_test.go`
- `internal/cli/inspect_helpers_test.go`
- `internal/cli/inspect_integration_test.go`
- `internal/vault/reader.go`

## Errors / Corrections
- `make verify` initially failed on staticcheck because the new error strings were capitalized and ended with periods; the inspect errors were normalized and verification rerun.
- Self-review caught an outdated `strings.HasPrefix` check in the generic entity resolver after the error messages were normalized; it was fixed before tracking updates.

## Ready for Next Run
- Task 16 implementation is complete and verified. `make verify` passes, CLI package coverage is 82.1%, and the new inspect lookup commands are covered by unit and integration tests.
