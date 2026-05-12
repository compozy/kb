# Task Memory: task_12.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot
- Add Java parse-stage operational telemetry to generate events (parse duration, resolver/fallback usage, unresolved counters) without breaking non-Java event consumers.
- Validate telemetry in unit tests and CLI integration JSON log output for Java ingest.

## Important Decisions
- Kept telemetry in the existing `stage_completed` parse event `fields` map to preserve event-kind/stage contracts and avoid introducing a new event type.
- Derived Java fallback and unresolved telemetry from parse diagnostics (`JAVA_RESOLUTION_FALLBACK`) already emitted by the Java adapter, avoiding adapter API churn and preserving deterministic low-overhead behavior.
- Emitted Java telemetry fields only when Java files are parsed (`java_files_processed > 0`) so non-Java runs keep existing payload shape plus `parsed_files`.

## Learnings
- `generate` parse-stage events can safely carry additional machine-readable telemetry through `fields` without changes to CLI JSON observer wiring.
- Counting unresolved Java targets from fallback diagnostic detail (`";"`-delimited unresolved fragments) gives a stable fallback signal for rollout observability while reusing existing diagnostic sources.

## Files / Surfaces
- `internal/generate/generate.go`
- `internal/generate/generate_test.go`
- `internal/cli/workflow_integration_test.go`

## Errors / Corrections
- Initial unit assertion expected a fixed parse duration value; corrected to type-safe non-negative duration assertion because stage timing can be zero in deterministic test clocks.

## Ready for Next Run
- Java parse-stage telemetry keys are now emitted in JSON logs for Java ingest: `java_parse_duration_millis`, `java_files_processed`, `java_resolver_mode`, `java_fallback_count`, and `java_unresolved_count`.
- Full gate passed with `make verify`; task tracking can remain marked completed unless new telemetry contract requirements are added.
