---
status: completed
title: Add best-effort enterprise module metadata hints
type: backend
complexity: medium
dependencies: []
---

# Task 10: Add best-effort enterprise module metadata hints

## Overview
Add best-effort module metadata hints for enterprise Java repositories to improve cross-module context without making ingest brittle. This task introduces optional metadata usage that enhances resolution quality while preserving non-blocking behavior.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- Module metadata parsing MUST be best-effort and MUST NOT fail ingest when metadata is missing or malformed.
- The resolver SHOULD use discovered module hints to improve relation consistency in multi-module repositories.
- Metadata handling MUST remain deterministic for identical repository snapshots.
- The implementation MUST preserve existing scanner and adapter contracts for repositories without module metadata.
</requirements>

## Subtasks
- [x] 10.1 Identify and parse minimal Maven/Gradle module metadata signals needed by resolver.
- [x] 10.2 Feed module hints into Java resolution context as optional inputs.
- [x] 10.3 Preserve non-blocking fallback path when metadata parsing is unavailable.
- [x] 10.4 Add unit tests for metadata present/missing/malformed scenarios.
- [x] 10.5 Add integration fixture for multi-module metadata-assisted resolution.

## Implementation Details
Implement module hint extraction in a minimal, optional path and consume it in Java resolution context without introducing hard dependencies on build-system fidelity. Reference TechSpec sections "Integration Points" and "Known Risks".

### Relevant Files
- `internal/adapter/java_adapter.go` — resolution context and optional metadata usage.
- `internal/scanner/scanner.go` — repository traversal context used by metadata discovery.
- `internal/adapter/java_adapter_test.go` — unit tests for optional metadata paths.
- `internal/adapter/java_adapter_integration_test.go` — integration fixture for multi-module metadata hints.

### Dependent Files
- `internal/generate/generate.go` — parse diagnostics and relation outputs affected by metadata hints.
- `internal/cli/workflow_integration_test.go` — E2E expectations in enterprise multi-module fixtures.

### Related ADRs
- [ADR-002: Use deep Java relation resolution with safe syntactic fallback](../adrs/adr-002.md) — metadata hints should improve deep resolution before fallback.
- [ADR-006: Close Java ingest MVP rollout using available pilot evidence](../adrs/adr-006.md) — Phase 2/3 follow-up calls for enterprise fidelity hardening.

## Deliverables
- Best-effort module metadata hint flow integrated with Java resolver context.
- Non-blocking behavior for missing or malformed metadata.
- Unit and integration tests for metadata-assisted multi-module behavior.
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for enterprise module hint behavior **(REQUIRED)**

## Tests
- Unit tests:
  - [x] Valid metadata source produces expected module hints.
  - [x] Missing metadata path leaves resolution behavior unchanged and successful.
  - [x] Malformed metadata emits diagnostics but does not fail parse stage.
- Integration tests:
  - [x] Multi-module fixture with metadata improves cross-module relation consistency.
  - [x] Equivalent fixture without metadata still completes ingest via fallback.
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- Enterprise module metadata improves relation consistency where available
- Ingest remains resilient and non-blocking when metadata quality is poor
