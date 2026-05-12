---
status: completed
title: Add Java operational observability telemetry
type: backend
complexity: high
dependencies:
  - task_11
---

# Task 12: Add Java operational observability telemetry

## Overview
Strengthen production visibility for Java ingest operations by exposing structured telemetry for parse duration, fallback usage, and unresolved relation signals. This task supports broad rollout governance and faster triage in large enterprise repositories.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- Generate-stage observability MUST include Java parse/fallback signals in structured event fields.
- Telemetry output MUST remain machine-readable for `--log-format json` consumers.
- Added observability SHOULD not break existing event contracts for non-Java ingest flows.
- The implementation MUST preserve deterministic and low-overhead event emission.
</requirements>

## Subtasks
- [x] 12.1 Add Java-focused structured telemetry fields to generate parse-stage events.
- [x] 12.2 Ensure fallback/unresolved counters are emitted in stable JSON-compatible form.
- [x] 12.3 Keep compatibility for existing event consumers and non-Java flows.
- [x] 12.4 Add unit tests for event field presence and stability.
- [x] 12.5 Add integration validation for JSON log output during Java ingest.

## Implementation Details
Extend existing generate event emission with Java-centric observability fields while preserving current event model compatibility. Reference TechSpec sections "Monitoring and Observability" and "System Architecture".

### Relevant Files
- `internal/generate/generate.go` — event emission and stage lifecycle.
- `internal/generate/events.go` — event shape and field map conventions.
- `internal/generate/generate_test.go` — event assertions and stage progress tests.
- `internal/cli/generate.go` — wiring for log format selection and observer usage.
- `internal/adapter/java_adapter.go` — source diagnostics and counters for observability data.

### Dependent Files
- `internal/cli/workflow_integration_test.go` — can validate observable behavior in E2E runs.
- `.compozy/tasks/java-ingest-adapter/_rollout-mvp-signoff.md` — future sign-off should reference observability metrics.

### Related ADRs
- [ADR-005: Define MVP governance acceptance gates and pilot corpus](../adrs/adr-005.md) — observability supports governance evidence.
- [ADR-006: Close Java ingest MVP rollout using available pilot evidence](../adrs/adr-006.md) — Phase 3 follow-up requires stronger operational telemetry.

## Deliverables
- Structured Java ingest observability fields in generate events.
- Compatibility-safe telemetry behavior across log formats and language mixes.
- Unit and integration tests for telemetry correctness.
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for Java telemetry in event output **(REQUIRED)**

## Tests
- Unit tests:
  - [x] Parse-stage event includes Java-specific counters/fields when Java files are processed.
  - [x] Non-Java runs do not emit malformed Java telemetry fields.
  - [x] Event JSON shape remains stable and parseable.
- Integration tests:
  - [x] Java ingest with `--log-format json` includes expected structured telemetry.
  - [x] Existing parse/write progress event tests remain passing.
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- Java ingest operational metrics are visible and machine-readable in standard run output
- Observability additions do not regress existing event consumers
