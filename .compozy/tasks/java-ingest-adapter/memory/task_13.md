# Task Memory: task_13.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot
- Expand benchmark coverage to the ADR-005 canonical Java corpus and enforce a reproducible runtime gate policy for rollout governance.
- Produce archive-friendly Phase 3 baseline evidence with deterministic command and output format.

## Important Decisions
- Introduced shared benchmark policy helpers in `internal/generate/benchmark_policy.go` so canonical profiles, repeat count (`3`), overhead budget (`20%`), and dry-run flags are defined in one place.
- Updated the Java ingest performance integration gate to iterate all canonical profiles (`single-module-library`, `spring-service`, `multi-module-enterprise`) instead of a single synthetic Java fixture.
- Standardized reproducible execution via `make benchmark-java-rollout` and documented archive format in a dedicated Phase 3 baseline artifact.

## Learnings
- Keeping policy constants behind a non-tagged helper (`canonicalJavaBenchmarkPolicy`) avoids lint issues when integration-only tests consume governance values.
- The canonical profile loop provides deterministic PASS/FAIL threshold behavior while preserving lightweight runtime by using generated fixtures and dry-run options.
- Storing both median gate table and benchmark `ns/op` snapshot in one artifact makes historical comparison straightforward without replaying raw logs.

## Files / Surfaces
- `internal/generate/benchmark_policy.go`
- `internal/generate/benchmark_policy_test.go`
- `internal/generate/generate_integration_test.go`
- `internal/generate/testdata/java-benchmark-corpus/README.md`
- `Makefile`
- `.compozy/tasks/java-ingest-adapter/_phase3-benchmark-baseline.md`

## Errors / Corrections
- `make verify` initially failed on unused benchmark policy constants; corrected by introducing `canonicalJavaBenchmarkPolicy()` and consuming it from integration tests.

## Ready for Next Run
- Canonical benchmark gate is now reproducible with `make benchmark-java-rollout`, and Task 13 baseline evidence is captured in `_phase3-benchmark-baseline.md`.
- `go test -tags integration ./internal/generate -cover` reports `86.8%` coverage for `internal/generate`, satisfying the task coverage target.
