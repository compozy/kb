---
status: completed
title: Expand rollout benchmark corpus and reproducible gate
type: test
complexity: high
dependencies:
  - task_11
---

# Task 13: Expand rollout benchmark corpus and reproducible gate

## Overview
Expand and standardize benchmark evidence so Java ingest rollout decisions remain reproducible across representative repository profiles. This task operationalizes the governance threshold with repeatable corpus and run policy.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- The benchmark corpus MUST cover the canonical repository profiles used by governance decisions.
- Runtime gate evaluation MUST be reproducible (same flags, repeated runs, median-based comparison).
- Benchmark outputs SHOULD be easy to archive and compare over time.
- The benchmark workflow MUST remain compatible with repository verification practices.
</requirements>

## Subtasks
- [x] 13.1 Define/curate benchmark fixtures for canonical Java profile coverage.
- [x] 13.2 Standardize benchmark execution policy (repeat count, flags, median extraction).
- [x] 13.3 Add or update benchmark tests/commands to enforce reproducible comparisons.
- [x] 13.4 Document benchmark evidence capture format for rollout governance.
- [x] 13.5 Run benchmark suite and capture baseline artifact for Phase 3.

## Implementation Details
Build on existing generate integration benchmarks and align them with governance corpus requirements for long-term rollout control. Reference TechSpec sections "Benchmark and E2E Validation" and "Technical Dependencies".

### Relevant Files
- `internal/generate/generate_integration_test.go` — benchmark and integration budget tests.
- `internal/generate/testdata/` — benchmark fixture definitions.
- `Makefile` — optional benchmark command wrappers for reproducibility.
- `.compozy/tasks/java-ingest-adapter/_rollout-mvp-signoff.md` — prior baseline evidence reference.

### Dependent Files
- `.compozy/tasks/java-ingest-adapter/adrs/adr-005.md` — governance threshold source.
- `.compozy/tasks/java-ingest-adapter/adrs/adr-006.md` — rollout closure context and deferred evidence.
- `.compozy/tasks/java-ingest-adapter/task_15.md` — adoption playbook should reference standardized benchmark flow.

### Related ADRs
- [ADR-003: Enforce 20% ingest performance budget with hybrid caching strategy](../adrs/adr-003.md) — benchmark gate definition.
- [ADR-005: Define MVP governance acceptance gates and pilot corpus](../adrs/adr-005.md) — canonical corpus and threshold policy.

## Deliverables
- Expanded benchmark corpus aligned with canonical profile coverage.
- Reproducible benchmark run policy and execution flow.
- Captured Phase 3 benchmark baseline evidence artifact.
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for benchmark gate behavior **(REQUIRED)**

## Tests
- Unit tests:
  - [x] Benchmark helper logic computes medians consistently for repeated runs.
  - [x] Benchmark fixture selection logic maps to canonical profile set.
- Integration tests:
  - [x] Java benchmark suite executes successfully on canonical fixtures.
  - [x] Gate comparison reports PASS/FAIL deterministically for threshold checks.
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- Benchmark governance evidence is reproducible and auditable across canonical profiles
- Performance gate tracking is operational for ongoing rollout decisions
