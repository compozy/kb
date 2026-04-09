# Task Memory: task_08.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot
- Implement Task 08 on the current branch by adding the `internal/metrics` package, extending `internal/models` with the metrics result types it depends on, covering the engine with unit and integration tests, and validating with package coverage plus `make verify`.

## Important Decisions
- Use `internal/models` and `internal/metrics` as the canonical package paths for this task; `_techspec.md` still points at stale `internal/kodebase/...` paths.
- Treat the existing task spec, workflow memory, and TypeScript `compute-metrics.ts` as the approved design baseline, so this run stays focused on implementation rather than creating a new design artifact.
- Reuse the branch's existing `file:` / `symbol:` / `external:` ID conventions from the adapters instead of introducing Task 09 path-utility scope.
- Treat the missing `.compozy/tasks/rewrite/adrs/` directory as absent context, not a blocker, because the caller only provided the PRD directory and there are no ADR files to reconcile.

## Learnings
- The current branch has `internal/graph` and adapter outputs in place, but `internal/models` still lacks `SymbolMetrics`, `FileMetrics`, `DirectoryMetrics`, and `MetricsResult`, so Task 08 must supply that model surface before the engine can compile.
- The repo does not yet contain `internal/cli` or `internal/vault`, so this task should stay scoped to the metrics engine and its direct model/test surfaces.
- Task-specific validation for this run is `go test ./internal/metrics`, `go test -tags integration ./internal/metrics`, package coverage for `internal/metrics`, and the repo-wide `make verify` gate.
- The TypeScript reference behavior ports cleanly to Go with local helpers only: BFS for blast radius, DFS plus canonical rotation for cycles, and normalized PageRank-style centrality over `calls` and `references` edges.
- An integration-tagged TS workspace fixture is sufficient to validate the real adapter-plus-normalizer-to-metrics pipeline without depending on later CLI or vault tasks.

## Files / Surfaces
- `internal/models/models.go`
- `internal/metrics/compute.go`
- `internal/metrics/compute_test.go`
- `internal/metrics/compute_integration_test.go`
- `.compozy/tasks/rewrite/task_08.md`
- `.compozy/tasks/rewrite/_tasks.md`

## Errors / Corrections
- The PRD task workflow says to read ADRs, but `.compozy/tasks/rewrite/adrs/` does not exist on this branch; proceeding with the task spec, `_techspec.md`, and workflow memory as the authoritative context.
- No implementation corrections were required after verification; `go test ./internal/metrics`, `go test -tags integration ./internal/metrics`, coverage, and `make verify` all passed on the first validation run.

## Ready for Next Run
- Verification complete:
  - `go test ./internal/metrics`
  - `go test -tags integration ./internal/metrics`
  - `go test -coverprofile=/tmp/kodebase-metrics.cover ./internal/metrics && go tool cover -func=/tmp/kodebase-metrics.cover` => 89.7%
  - `make verify`
- Local code commit created: `368bc29` (`feat: implement metrics engine`)
- Tracking files were updated locally and should remain unstaged for the automatic code commit.
