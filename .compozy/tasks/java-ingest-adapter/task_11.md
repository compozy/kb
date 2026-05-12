---
status: completed
title: Validate Phase 2 regression suite for Java fidelity
type: test
complexity: high
dependencies:
  - task_07
  - task_08
  - task_09
  - task_10
---

# Task 11: Validate Phase 2 regression suite for Java fidelity

## Overview
Consolidate Phase 2 improvements into a regression suite that proves relation fidelity gains and stability on enterprise-style inputs. This task closes Phase 2 with reproducible evidence across adapter tests, CLI E2E, and performance checks.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- Phase 2 regression coverage MUST include nested types, wildcard imports, ambiguity policy, and metadata-assisted multi-module scenarios.
- CLI E2E validation MUST confirm Java ingest outputs remain valid and lint-clean.
- Benchmark checks SHOULD verify no unexpected regression beyond accepted budget.
- Regression evidence MUST be reproducible and documented for transition to Phase 3.
</requirements>

## Subtasks
- [x] 11.1 Extend adapter integration fixture matrix to cover all Phase 2 behaviors.
- [x] 11.2 Add/adjust CLI E2E tests to assert stable Java summaries and artifacts.
- [x] 11.3 Re-run benchmark validation with updated Phase 2 behavior and capture deltas.
- [x] 11.4 Update task memory and rollout notes with Phase 2 validation outcomes.
- [x] 11.5 Run full repository verification gate after test updates.

## Implementation Details
Compose a cohesive regression pack using existing testing patterns in adapter, generate, and CLI integration surfaces. Reference TechSpec sections "Testing Approach", "Benchmark and E2E Validation", and "Monitoring and Observability".

### Relevant Files
- `internal/adapter/java_adapter_integration_test.go` — Phase 2 fidelity scenarios.
- `internal/adapter/java_adapter_test.go` — detailed unit assertions.
- `internal/cli/workflow_integration_test.go` — end-to-end ingest and lint validation.
- `internal/generate/generate_integration_test.go` — performance budget and integration evidence.
- `.compozy/tasks/java-ingest-adapter/memory/task_11.md` — evidence notes for handoff.

### Dependent Files
- `.compozy/tasks/java-ingest-adapter/_rollout-mvp-signoff.md` — baseline reference for comparing post-Phase 2 improvements.
- `.compozy/tasks/java-ingest-adapter/_tasks.md` — status updates after validation completion.

### Related ADRs
- [ADR-003: Enforce 20% ingest performance budget with hybrid caching strategy](../adrs/adr-003.md) — regression checks must retain budget compliance.
- [ADR-004: Require unit, integration, benchmark, and CLI E2E validation for Java ingest](../adrs/adr-004.md) — defines mandatory validation style.

## Deliverables
- Expanded Phase 2 regression coverage across adapter, generate, and CLI surfaces.
- Reproducible benchmark comparison for post-Phase 2 behavior.
- Documented validation evidence for Phase 3 entry.
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for full Phase 2 behavior matrix **(REQUIRED)**

## Tests
- Unit tests:
  - [x] Nested/wildcard/ambiguity metadata paths each have explicit assertion coverage.
  - [x] Java fallback diagnostics remain predictable under mixed complex scenarios.
- Integration tests:
  - [x] CLI ingest on enterprise-style fixture remains successful with expected summary values.
  - [x] `kb lint` on generated topic remains clean for tested fixtures.
  - [x] Performance test compares updated Java path against accepted budget baseline.
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- Phase 2 improvements are validated with reproducible and traceable regression evidence
- Phase 3 planning can rely on stable, measured Phase 2 outputs
