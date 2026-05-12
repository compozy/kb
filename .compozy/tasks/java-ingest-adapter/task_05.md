---
status: completed
title: Add deep Java relation resolution with fallback
type: backend
complexity: high
dependencies:
  - task_03
---

# Task 05: Add deep Java relation resolution with fallback

## Overview
Enhance the Java adapter with deep relation resolution to improve cross-file dependency and call accuracy while retaining automatic syntactic fallback for unresolved cases. This task operationalizes the technical direction from ADR-002 and the performance guardrails from ADR-003.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- Deep relation resolution MUST be attempted first for Java symbols where repository metadata allows.
- The adapter MUST automatically fallback to syntactic relation resolution without failing ingest.
- Fallback paths MUST emit structured diagnostics that distinguish fallback from parse failures.
- Relation resolution behavior MUST preserve deterministic output ordering.
- Implementation SHOULD avoid runtime regressions that violate the 20% overhead budget target.
</requirements>

## Subtasks
- [x] 5.1 Add deep resolver abstraction and package/classpath-aware resolution pass.
- [x] 5.2 Add automatic fallback path to syntactic relation resolution for unresolved targets.
- [x] 5.3 Emit `JAVA_RESOLUTION_FALLBACK` diagnostics for fallback scenarios.
- [x] 5.4 Add unit tests for deep-resolution success and fallback behavior.
- [x] 5.5 Add integration fixtures covering multi-package and partial-metadata Java repositories.
- [x] 5.6 Validate deterministic ordering and relation consistency across repeated runs.

## Implementation Details
Extend Java adapter internals (or companion helper files within `internal/adapter`) to support deep-first resolution and graceful fallback. Keep behavior aligned with TechSpec sections “Core Interfaces,” “Data Models,” and “Technical Considerations,” and avoid introducing cross-package API churn.

### Relevant Files
- `internal/adapter/java_adapter.go` — primary implementation location for resolver flow.
- `internal/models/models.go` — diagnostic stage/severity and relation types used by adapter output.
- `internal/adapter/rust_adapter.go` — reference for multi-file resolution and relation dedup patterns.
- `internal/adapter/ts_adapter.go` — reference for relation dedup and import binding flow.
- `internal/adapter/java_adapter_test.go` — unit tests for resolution branches.
- `internal/adapter/java_adapter_integration_test.go` — integration tests for cross-file accuracy.

### Dependent Files
- `internal/graph/normalize.go` — consumes additional relations and diagnostics from Java adapter.
- `internal/generate/generate.go` — parse-stage output volume and diagnostics reporting depend on resolver behavior.
- `internal/generate/generate_test.go` — summaries may reflect changed diagnostics/relation counts.

### Related ADRs
- [ADR-002: Use deep Java relation resolution with safe syntactic fallback](../adrs/adr-002.md) — establishes deep-first fallback strategy.
- [ADR-003: Enforce 20% ingest performance budget with hybrid caching strategy](../adrs/adr-003.md) — requires bounded runtime impact.

## Deliverables
- Deep Java relation resolver integrated into adapter parse flow.
- Automatic syntactic fallback path with explicit fallback diagnostics.
- Unit and integration tests covering deep success, unresolved fallback, and deterministic behavior.
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for deep-resolution and fallback scenarios **(REQUIRED)**

## Tests
- Unit tests:
  - [x] Deep resolver maps imports/calls to cross-file symbols in a multi-package fixture.
  - [x] Unresolvable deep target triggers fallback diagnostic and still emits syntactic relations.
  - [x] Resolver output ordering remains stable across identical repeated runs.
  - [x] Adapter does not return hard error when deep resolution cannot fully resolve metadata.
- Integration tests:
  - [x] Multi-module Java fixture demonstrates improved relation quality versus MVP-only path.
  - [x] Partial classpath fixture completes ingest and records fallback diagnostics.
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- Deep-first resolution improves cross-file Java relation quality while preserving ingest completion
- Fallback behavior is visible through diagnostics and never blocks parsing flow
