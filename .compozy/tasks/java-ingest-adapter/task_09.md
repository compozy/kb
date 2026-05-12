---
status: completed
title: Add deterministic policy for ambiguous import targets
type: backend
complexity: medium
dependencies:
  - task_08
---

# Task 09: Add deterministic policy for ambiguous import targets

## Overview
Define and implement deterministic handling for ambiguous Java import targets, including duplicate simple names and static import conflicts. This task reduces false positives and makes resolution outcomes predictable for large codebases.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- Ambiguous import targets MUST follow a deterministic resolution policy.
- The resolver MUST avoid emitting misleading semantic relations in unresolved ambiguity cases.
- Ambiguous cases SHOULD produce structured diagnostics to aid governance and debugging.
- Existing non-ambiguous import resolution behavior MUST remain stable.
</requirements>

## Subtasks
- [x] 9.1 Define deterministic precedence and ambiguity handling rules for import conflicts.
- [x] 9.2 Apply ambiguity policy across deep-resolution and fallback handoff.
- [x] 9.3 Emit clear diagnostics for unresolved ambiguous targets.
- [x] 9.4 Add unit tests for duplicate simple-name imports and static import collisions.
- [x] 9.5 Add integration fixture asserting no unstable relation drift in ambiguous scenarios.

## Implementation Details
Use the wildcard-aware import model from Task 08 and layer deterministic ambiguity policy on top. Align decisions with TechSpec sections "Key Decisions", "Known Risks", and "Monitoring and Observability".

### Relevant Files
- `internal/adapter/java_adapter.go` — ambiguous candidate handling and diagnostic emission.
- `internal/adapter/java_adapter_test.go` — unit-level ambiguity scenarios.
- `internal/adapter/java_adapter_integration_test.go` — integration-level ambiguity regression checks.
- `internal/models/models.go` — diagnostic structure and severity/stage definitions.

### Dependent Files
- `internal/generate/generate.go` — summary diagnostics and counts may change with ambiguity handling.
- `.compozy/tasks/java-ingest-adapter/_rollout-mvp-signoff.md` — future sign-off updates should reference lower ambiguity rates.

### Related ADRs
- [ADR-002: Use deep Java relation resolution with safe syntactic fallback](../adrs/adr-002.md) — ambiguity policy governs deep/fallback boundary.
- [ADR-005: Define MVP governance acceptance gates and pilot corpus](../adrs/adr-005.md) — deterministic behavior supports reliable governance evidence.

## Deliverables
- Deterministic ambiguity resolution policy implemented in Java adapter.
- Structured diagnostics for ambiguous unresolved targets.
- Unit and integration tests for ambiguity scenarios.
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for ambiguous import behavior **(REQUIRED)**

## Tests
- Unit tests:
  - [x] Duplicate simple-name imports trigger deterministic ambiguity handling.
  - [x] Static import conflicts avoid incorrect semantic relation emission.
  - [x] Ambiguous cases emit expected warning diagnostics.
- Integration tests:
  - [x] Ambiguous fixture produces stable relation output across repeated runs.
  - [x] Non-ambiguous fixtures keep existing expected outputs unchanged.
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- Ambiguous import scenarios no longer produce non-deterministic relation results
- Diagnostics clearly expose ambiguity without breaking ingest completion
