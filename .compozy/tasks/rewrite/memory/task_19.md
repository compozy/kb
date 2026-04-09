# Task Memory: task_19.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot
- Add end-to-end integration coverage for the full generate + inspect pipeline using a richer fixture repository and keep validation under 30 seconds.

## Important Decisions
- No implementation decision yet. Work is blocked on the task-doc mismatch around expected base-file counts.

## Learnings
- The current Go renderer integration test expects 11 base files, and that matches the upstream TypeScript renderer and README.
- Task 19 and `_techspec.md` still say 12 base files, while `task_10.md` already records this as stale task text rather than a real renderer requirement.
- `.compozy/tasks/rewrite/adrs/` is absent on this branch, so there is no ADR content to reconcile for this task.

## Files / Surfaces
- `internal/generate/generate_integration_test.go`
- `internal/cli/inspect_integration_test.go`
- `internal/vault/render_integration_test.go`
- `.compozy/tasks/rewrite/task_19.md`
- `.compozy/tasks/rewrite/_techspec.md`

## Errors / Corrections
- Blocker: governing sources disagree on base-file count. Current implementation and reference behavior are 11; task text still says 12.

## Ready for Next Run
- Await confirmation on whether Task 19 should follow the authoritative 11-base behavior or the stale 12-base wording before implementation starts.
