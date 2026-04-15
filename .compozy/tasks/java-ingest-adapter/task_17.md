---
status: completed
title: Harden large-scale Java ingest operational behavior
type: backend
complexity: high
dependencies:
  - task_12
  - task_11
---

# Task 17: Harden large-scale Java ingest operational behavior

## Overview
Harden Java ingest behavior for large-scale enterprise operation by improving runtime resilience, resource safety, and operational predictability. This task focuses on production-grade robustness without changing the user-facing workflow model.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- Large Java ingest runs MUST remain operationally safe under high file counts and complex relation workloads.
- Runtime behavior SHOULD include predictable handling for high diagnostic/fallback volume scenarios.
- Hardening changes MUST preserve deterministic outputs and existing command contracts.
- The implementation MUST stay within agreed performance guardrails for representative workloads.
</requirements>

## Subtasks
- [x] 17.1 Identify and address high-scale operational bottlenecks in Java ingest path.
- [x] 17.2 Add safeguards for high-volume fallback/diagnostic scenarios.
- [x] 17.3 Improve operational predictability for long-running Java ingest executions.
- [x] 17.4 Add unit tests for hardening logic and deterministic safeguards.
- [x] 17.5 Add integration and benchmark checks for large-scale fixture behavior.

## Implementation Details
Apply hardening improvements in Java ingest execution paths with emphasis on stability and predictable operation in enterprise-scale repositories. Reference TechSpec sections "Known Risks", "Monitoring and Observability", and "Build Order" constraints.

### Relevant Files
- `internal/adapter/java_adapter.go` — core Java parse and relation workload behavior.
- `internal/generate/generate.go` — stage orchestration, progress, and summary timings.
- `internal/generate/generate_integration_test.go` — performance and integration behavior under realistic workload.
- `internal/cli/workflow_integration_test.go` — end-to-end operational flow validation.
- `internal/models/models.go` — summary/diagnostic model fields used in operational evidence.

### Dependent Files
- `.compozy/tasks/java-ingest-adapter/task_13.md` — benchmark gate and corpus use hardened behavior.
- `.compozy/tasks/java-ingest-adapter/task_15.md` — adoption playbook should reference hardening guidance.
- `.compozy/tasks/java-ingest-adapter/_rollout-mvp-signoff.md` — future rollout evidence should reflect improved operational stability.

### Related ADRs
- [ADR-003: Enforce 20% ingest performance budget with hybrid caching strategy](../adrs/adr-003.md) — hardening must remain budget-compliant.
- [ADR-006: Close Java ingest MVP rollout using available pilot evidence](../adrs/adr-006.md) — Phase 3 requires scale-focused follow-up.

## Deliverables
- Hardened Java ingest behavior for large-scale operational scenarios.
- Safeguards and deterministic handling for high diagnostic volume conditions.
- Updated integration/benchmark evidence demonstrating improved operational stability.
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for large-scale Java ingest behavior **(REQUIRED)**

## Tests
- Unit tests:
  - [x] High-volume diagnostic input paths remain deterministic and non-blocking.
  - [x] Hardening safeguards trigger expected behavior under stress-like conditions.
  - [x] Existing Java ingest contract fields remain unchanged.
- Integration tests:
  - [x] Large Java fixture ingest completes successfully with predictable timings and outputs.
  - [x] Benchmark comparison confirms no unacceptable regression from hardening changes.
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- Java ingest remains robust and predictable for large enterprise-scale repositories
- Hardening improvements preserve existing workflow compatibility and performance guardrails
