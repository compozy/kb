---
status: completed
title: Add diagnostics governance checks in lint workflow
type: backend
complexity: high
dependencies:
  - task_12
  - task_11
---

# Task 16: Add diagnostics governance checks in lint workflow

## Overview
Introduce governance-oriented quality checks for Java diagnostics so operators can enforce clearer acceptance criteria during broad rollout. This task connects Java ingest telemetry and diagnostics with actionable quality gates in existing lint-oriented workflows.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- Lint/governance checks MUST surface high-risk Java diagnostics patterns in a consistent, machine-readable form.
- Checks MUST distinguish parse errors from fallback warnings and avoid over-blocking normal fallback behavior.
- Governance checks SHOULD support threshold-based policies for rollout operations.
- Existing lint behavior for non-Java topics MUST remain backward-compatible.
</requirements>

## Subtasks
- [x] 16.1 Define Java diagnostics governance policy and threshold model for lint workflow.
- [x] 16.2 Implement diagnostics aggregation/check logic in lint-compatible surfaces.
- [x] 16.3 Add machine-readable output support for governance checks.
- [x] 16.4 Add unit tests for threshold pass/fail behavior and diagnostic categorization.
- [x] 16.5 Add integration tests with Java-generated topics containing controlled diagnostics.

## Implementation Details
Extend lint or adjacent quality-evaluation pathways with Java diagnostics governance checks that are strict enough for rollout control but compatible with expected fallback behavior. Reference TechSpec sections "Monitoring and Observability" and "Technical Considerations."

### Relevant Files
- `internal/lint/lint.go` — quality issue modeling and reporting path.
- `internal/lint/lint_test.go` — lint behavior and output assertions.
- `internal/models/models.go` — structured diagnostic definitions and severity/stage fields.
- `internal/vault/reader.go` — source data loading for lint/inspect quality checks.
- `internal/cli/lint.go` — command-level output and option handling.

### Dependent Files
- `internal/cli/workflow_integration_test.go` — end-to-end lint behavior validation with Java topics.
- `.compozy/tasks/java-ingest-adapter/_rollout-mvp-signoff.md` — future sign-offs should consume new governance check outcomes.

### Related ADRs
- [ADR-005: Define MVP governance acceptance gates and pilot corpus](../adrs/adr-005.md) — governance checks support objective rollout criteria.
- [ADR-006: Close Java ingest MVP rollout using available pilot evidence](../adrs/adr-006.md) — Phase 3 requires stronger quality governance.

## Deliverables
- Java diagnostics governance checks integrated into lint-compatible workflow.
- Threshold-based pass/fail reporting for governance operations.
- Unit and integration tests validating governance behavior.
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for governance check behavior on Java topics **(REQUIRED)**

## Tests
- Unit tests:
  - [x] `JAVA_PARSE_ERROR` contributes to blocking governance outcomes as defined.
  - [x] `JAVA_RESOLUTION_FALLBACK` is categorized and thresholded without false blocking defaults.
  - [x] Governance output includes machine-readable counts by diagnostic type.
- Integration tests:
  - [x] Java topic with controlled diagnostics yields expected governance check result.
  - [x] Non-Java and clean Java topics maintain existing lint compatibility.
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- Governance checks provide actionable Java diagnostics controls for rollout operations
- Existing lint workflow remains stable for current users
