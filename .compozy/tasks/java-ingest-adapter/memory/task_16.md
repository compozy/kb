# Task Memory: task_16.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot
- Implement diagnostics governance checks in `kb lint` for Java ingest telemetry, with threshold controls and machine-readable outcomes.

## Important Decisions
- Added lint-side Java governance policy with explicit thresholds:
  - `java-max-parse-errors` default `0` (blocking on any parse errors).
  - `java-max-fallback-warnings` default `-1` (fallback governance disabled by default to avoid over-blocking).
- Persisted Java diagnostic counters into `raw/codebase/index/java.md` frontmatter during render so lint can evaluate governance from topic artifacts without rerunning ingest.
- Represented governance outcomes as lint issues (`java-diagnostic-governance`) with JSON payload in `message` for machine-readable count/threshold/status data.

## Learnings
- Existing Java integration fixtures remain lint-clean under default policy because fallback governance is opt-in (`-1` default), preserving prior workflow compatibility.
- Emitting deterministic scalar counters in frontmatter is safer than nested objects for current markdown frontmatter renderer behavior.

## Files / Surfaces
- `internal/models/kb_models.go`
- `internal/models/kb_models_test.go`
- `internal/vault/render.go`
- `internal/lint/lint.go`
- `internal/lint/lint_test.go`
- `internal/cli/lint.go`
- `internal/cli/lint_test.go`
- `internal/cli/workflow_integration_test.go`

## Errors / Corrections
- Initial fallback-governance logic emitted warnings even with disabled threshold, which risked breaking existing Java lint expectations; corrected by skipping governance emission when threshold is negative.

## Ready for Next Run
- Governance signal is now available in lint outputs for Java parse errors by default and fallback diagnostics when threshold is enabled.
- Task tracking files were updated (`task_16.md` + `_tasks.md`) and verification evidence is green (`make verify` + targeted integration run).
