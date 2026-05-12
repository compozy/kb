# Task Memory: task_14.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot
- Stabilize the automation-facing JSON contract for `kb ingest codebase` by defining required keys, enforcing them in CLI unit/integration tests, and publishing compatibility guidance for future evolution.

## Important Decisions
- Reused existing payload shapes (`codebaseIngestResult` + `models.GenerationSummary`) without introducing new top-level schema fields, to keep current automation consumers backward-compatible.
- Defined contract stability through test-enforced required key sets (top-level result keys, summary keys, and timings keys) and mode semantics (`dryRun` write counters must remain zero).
- Published dedicated contract documentation in `.compozy/tasks/java-ingest-adapter/_automation-json-contract.md` and linked it from rollout signoff docs for adoption visibility.

## Learnings
- Existing CLI tests validated selected values but did not enforce a full required-key contract across dry-run/full-run modes.
- Contract validation is more maintainable when shared helper assertions live in `internal/cli/workflow_test_helpers_test.go` and are reused across both unit and integration suites.

## Files / Surfaces
- `internal/cli/workflow_test_helpers_test.go`
- `internal/cli/ingest_test.go`
- `internal/cli/workflow_integration_test.go`
- `.compozy/tasks/java-ingest-adapter/_automation-json-contract.md`
- `.compozy/tasks/java-ingest-adapter/_rollout-mvp-signoff.md`

## Errors / Corrections
- Corrected contract documentation to use the actual `sourceType` value (`codebase-file`) from `models.SourceKindCodebaseFile`.

## Ready for Next Run
- Validation evidence:
  - `go test ./internal/cli` (PASS)
  - `go test ./internal/cli -cover` -> `coverage: 80.6%` (PASS, >=80%)
  - `go test -tags integration ./internal/cli -run "TestCLIIntegrationJavaIngestJSONContractStableAcrossModes|TestCLIIntegrationScaffoldIngestJavaWorkspaceCodebase"` (PASS)
  - `make verify` (PASS)
