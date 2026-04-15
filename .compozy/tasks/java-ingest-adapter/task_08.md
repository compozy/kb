---
status: completed
title: Add wildcard import deep-resolution support
type: backend
complexity: high
dependencies: []
---

# Task 08: Add wildcard import deep-resolution support

## Overview
Add deep-resolution support for Java wildcard imports (`import pkg.*`) to reduce fallback noise and improve enterprise repository fidelity. This task focuses on building reliable wildcard lookup behavior while preserving deterministic output.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- The deep resolver MUST handle wildcard import references for resolvable symbols in scanned source.
- Wildcard import handling MUST avoid non-deterministic target selection.
- Fallback behavior MUST continue when wildcard expansion cannot resolve symbols.
- The implementation SHOULD avoid significant regression in parse stage runtime.
</requirements>

## Subtasks
- [x] 8.1 Extend Java import indexing to include wildcard package candidates.
- [x] 8.2 Update deep resolver to resolve simple type names via wildcard import indexes.
- [x] 8.3 Keep fallback diagnostics for unresolved wildcard cases.
- [x] 8.4 Add unit tests for wildcard import success and unresolved branches.
- [x] 8.5 Add integration fixtures covering wildcard-heavy repository patterns.

## Implementation Details
Evolve Java import lookup structures and deep resolution flow to incorporate wildcard imports using deterministic candidate selection rules. Reference TechSpec sections "Integration Points", "Data Models", and "Benchmark and E2E Validation" for constraints.

### Relevant Files
- `internal/adapter/java_adapter.go` — import parsing, lookup indexes, and deep resolution implementation.
- `internal/adapter/java_adapter_test.go` — unit tests for wildcard resolution behavior.
- `internal/adapter/java_adapter_integration_test.go` — integration coverage for wildcard import scenarios.
- `internal/models/models.go` — diagnostic and relation structures.

### Dependent Files
- `internal/generate/generate.go` — parse diagnostics and relation counts influenced by wildcard resolution quality.
- `internal/graph/normalize.go` — downstream edge normalization depends on relation consistency.
- `internal/cli/workflow_integration_test.go` — E2E expectations may change as fallback volume decreases.

### Related ADRs
- [ADR-002: Use deep Java relation resolution with safe syntactic fallback](../adrs/adr-002.md) — wildcard support extends deep-first strategy.
- [ADR-003: Enforce 20% ingest performance budget with hybrid caching strategy](../adrs/adr-003.md) — wildcard expansion must remain budget-safe.

## Deliverables
- Wildcard import-aware deep resolution in Java adapter.
- Preserved deterministic output and fallback diagnostics when unresolved.
- Unit and integration tests for wildcard import behavior.
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for wildcard import cases **(REQUIRED)**

## Tests
- Unit tests:
  - [x] `import pkg.*` resolves known symbols when package candidates exist.
  - [x] Unresolvable wildcard imports emit fallback diagnostics without hard failure.
  - [x] Deterministic target selection for repeated wildcard parse runs.
- Integration tests:
  - [x] Multi-file fixture with wildcard imports emits expected `references` edges.
  - [x] Parse stage remains successful when wildcard targets are partially missing.
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- Wildcard imports improve deep resolution fidelity in common enterprise patterns
- Fallback remains safe and deterministic when wildcard candidates are unavailable
