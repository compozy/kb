# Task Memory: task_17.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot

- Add `kb lint` to the Cobra root, wired to `internal/lint`, with formatter output and optional report persistence.
- Narrow inspect snapshot reads to topic-local codebase content under `raw/codebase` without changing inspect-facing relative paths or relation lookups.
- Confirm `search` and `index` continue to resolve topic-scoped collection names and paths through the topic resolver, then extend tests where behavior was previously implicit.

## Important Decisions

- Keep `vault.ResolveVaultQuery` as the topic resolver for CLI read commands; apply inspect-specific codebase narrowing in the CLI layer instead of changing the general vault query contract.
- Preserve inspect document `RelativePath` values as `raw/codebase/...` even when reading from the `raw/codebase` subtree so existing inspect outputs and circular-dependency link matching stay stable.
- Keep `search` and `index` production code unchanged because they already derive topic-scoped collection names and topic paths through `ResolveVaultQuery`; add explicit regression tests instead of refactoring for its own sake.

## Learnings

- `internal/vault.ReadVaultSnapshot` computes document relative paths from `ResolvedVault.TopicPath`, so pointing it directly at `raw/codebase` would otherwise strip the `raw/codebase/` prefix from inspect outputs.
- `internal/cli` package coverage reached 80.2% after adding lint command coverage plus topic-resolution assertions for inspect/search/index.

## Files / Surfaces

- `internal/cli/root.go`
- `internal/cli/lint.go` (new)
- `internal/cli/lint_test.go` (new)
- `internal/cli/inspect.go`
- `internal/cli/inspect_test.go`
- `internal/cli/search_test.go`
- `internal/cli/index_test.go`
- `internal/cli/search.go`
- `internal/cli/index.go`
- `internal/cli/*_test.go`

## Errors / Corrections

- No blocking issues found during self-review; `go test ./internal/cli`, `go test -cover ./internal/cli`, and `make verify` all passed after the first implementation pass.

## Ready for Next Run

- Tracking files still need the final status/checkbox update and the local commit for task completion.
