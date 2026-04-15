---
status: completed
title: Improve nested and inner Java type resolution
type: backend
complexity: high
dependencies: []
---

# Task 07: Improve nested and inner Java type resolution

## Overview
Improve Java relation fidelity for nested and inner class patterns that are common in enterprise repositories. This task reduces unresolved relationships by teaching the adapter to model ownership and qualified names for nested types consistently.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- The Java adapter MUST represent nested and inner type ownership in a deterministic way.
- The Java resolver MUST resolve common `Outer.Inner` references without degrading existing top-level resolution.
- Symbol and relation IDs MUST remain deterministic across repeated runs.
- The implementation SHOULD preserve current parse performance characteristics for non-nested code.
</requirements>

## Subtasks
- [x] 7.1 Extend Java symbol modeling to include nested/inner type ownership context.
- [x] 7.2 Update deep and syntactic resolution paths for qualified nested type usage.
- [x] 7.3 Add deterministic ordering checks for nested-type symbols and relations.
- [x] 7.4 Add unit coverage for nested declarations and cross-file usage patterns.
- [x] 7.5 Add integration fixture coverage for nested classes across multiple files.

## Implementation Details
Implement nested-type awareness in the Java adapter resolution context and ensure compatibility with existing output contracts. Use TechSpec sections "Core Interfaces", "Data Models", and "Known Risks" as implementation guidance.

### Relevant Files
- `internal/adapter/java_adapter.go` — primary parser/resolver logic for Java symbols and relations.
- `internal/adapter/java_adapter_test.go` — unit coverage for symbol extraction and relation behavior.
- `internal/adapter/java_adapter_integration_test.go` — integration fixtures and cross-file assertions.
- `internal/models/models.go` — relation/symbol structures consumed by adapter output.

### Dependent Files
- `internal/graph/normalize.go` — consumes adapter output and depends on deterministic IDs.
- `internal/generate/generate.go` — summary counts and diagnostics depend on adapter relation quality.
- `internal/vault/render.go` — downstream rendering quality depends on improved relation graph.

### Related ADRs
- [ADR-002: Use deep Java relation resolution with safe syntactic fallback](../adrs/adr-002.md) — nested resolution must still preserve fallback behavior.

## Deliverables
- Nested/inner type-aware symbol and relation extraction in Java adapter.
- Deterministic relation emission for nested type scenarios.
- Unit and integration tests for nested type resolution behavior.
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for nested cross-file resolution **(REQUIRED)**

## Tests
- Unit tests:
  - [x] Nested class declaration emits expected symbol ownership metadata.
  - [x] `Outer.Inner` references resolve to the correct target symbol.
  - [x] Nested-type relation output is stable across repeated parse runs.
- Integration tests:
  - [x] Multi-file fixture with nested classes resolves expected `references/calls` edges.
  - [x] Existing top-level type resolution assertions remain unchanged.
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- Nested and inner Java types resolve with fewer ambiguous/unresolved relations
- Adapter output remains deterministic and compatible with existing pipeline consumers
