# Task Memory: task_07.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot
- Implement Task 07 on the current branch by adding the `internal/graph` normalizer, covering it with unit and integration tests, and validating with package coverage plus `make verify`.

## Important Decisions
- Use `internal/models` and `internal/graph` as the canonical package paths for this task; `_techspec.md` still points at stale `internal/kodebase/...` paths.
- Treat the existing task spec, workflow memory, and TypeScript `normalize-graph.ts` as the approved design baseline, so this run stays focused on implementation rather than creating a new design artifact.
- Return `models.GraphSnapshot` by value to match the task specification and the TypeScript reference object shape, despite the older `_techspec.md` pointer signature.
- Preserve the reference behavior that omits files which contain only diagnostics and no renderable graph data.

## Learnings
- The active branch already had the graph node and relation types restored by Tasks 05 and 06, so Task 07 could stay scoped to aggregation/ordering behavior.
- The task requirements are stricter than the reference implementation on collection ordering; files, symbols, and external nodes are now ordered by stable IDs while relations use the `(fromId, toId, type, confidence)` composite key.
- Package-level validation for this task is `go test ./internal/graph`, `go test -tags integration ./internal/graph`, and `go test -coverprofile=/tmp/kodebase-graph.cover ./internal/graph && go tool cover -func=/tmp/kodebase-graph.cover`, which produced 91.3% coverage before the repo-wide `make verify` gate.

## Files / Surfaces
- `internal/graph/normalize.go`
- `internal/graph/normalize_test.go`
- `internal/graph/normalize_integration_test.go`
- `.compozy/tasks/rewrite/task_07.md`
- `.compozy/tasks/rewrite/_tasks.md`

## Errors / Corrections
- Corrected the deterministic-order unit test to match the implemented composite-key relation ordering (`fromId`, then `toId`, then `type`, then `confidence`).

## Ready for Next Run
- Verification complete:
  - `go test ./internal/graph`
  - `go test -tags integration ./internal/graph`
  - `go test -coverprofile=/tmp/kodebase-graph.cover ./internal/graph && go tool cover -func=/tmp/kodebase-graph.cover` => 91.3%
  - `make verify`
- Task tracking files were updated locally and should remain unstaged for the automatic code commit.
