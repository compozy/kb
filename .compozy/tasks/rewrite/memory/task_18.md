# Task Memory: task_18.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot
- Wire top-level `search` and `index` CLI commands to the QMD shell client and vault query resolver, with task-required flags, user-friendly QMD-unavailable handling, unit coverage >=80% for `internal/cli`, a tagged CLI integration test, and a clean `make verify`.

## Important Decisions
- Keep the TypeScript-compatible `--format`, `--vault`, and `--topic` flags on `search` in addition to the task-required flags so the Go CLI matches the existing output formatter and vault-resolution pattern.
- Expose `QMDClient.Status` and choose `add` vs `update` in the CLI because the real `qmd collection add` command fails when the collection already exists.
- Use `--embed=false` plus lexical search in CLI integration coverage so the test does not depend on sqlite-vector support on the host.

## Learnings
- The live `qmd status` output for an empty index still contains enough structure to parse as a zero-value status; the parser only needed to accept the `No collections.` variant.
- An explicit `--collection` on `search` should bypass vault/topic resolution, matching the upstream TypeScript command behavior.

## Files / Surfaces
- `internal/cli/root.go`
- `internal/cli/search.go`
- `internal/cli/index.go`
- `internal/cli/search_test.go`
- `internal/cli/index_test.go`
- `internal/cli/search_index_integration_test.go`
- `internal/qmd/client.go`
- `internal/qmd/client_test.go`

## Errors / Corrections
- Initial CLI integration test compile failed because the new file missed the `internal/qmd` import; fixed before rerunning focused tests and the full verify gate.
- Self-review caught that whitespace-only search queries were accepted; trimmed and validated the query before the final verification run.

## Ready for Next Run
- Task is complete in code. Fresh evidence: `go test ./internal/qmd ./internal/cli`, `go test ./internal/cli -cover` (81.4%), `go test -tags integration ./internal/cli -run TestSearchCommandReturnsResultsAgainstIndexedVault -v`, and `make verify` all passed after the final code change.
