# Task Memory: task_17.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot
- Harden Java ingest behavior for high-scale fallback/diagnostic scenarios without changing ingest command contracts.
- Keep runtime behavior deterministic under stress-like unresolved relation volume.

## Important Decisions
- Added deterministic caps for Java diagnostic payload construction in `internal/adapter/java_adapter.go`:
  - fallback detail capped by entry count and byte budget,
  - module hint warning detail capped by warning count and byte budget.
- Standardized truncation metadata marker as `meta:truncated (...)` so downstream telemetry parsing remains stable.

## Learnings
- Unbounded fallback diagnostic detail is the main scale risk path because high unresolved counts can inflate memory/string payloads even when ingest still succeeds.
- `countFallbackUnresolvedReferences` must ignore truncation metadata segments to avoid over-counting unresolved references in parse telemetry.

## Files / Surfaces
- `internal/adapter/java_adapter.go`
- `internal/adapter/java_adapter_test.go`
- `internal/generate/generate.go`
- `internal/generate/generate_test.go`
- `internal/generate/generate_integration_test.go`

## Errors / Corrections
- No implementation blockers found; targeted tests and full `make verify` passed after hardening changes.

## Ready for Next Run
- If future tasks add new diagnostic segment types, keep `meta:*`-style markers excluded from unresolved telemetry counts.
- If large-enterprise fixtures increase unresolved pressure, tune cap constants with benchmark evidence rather than removing bounds.
