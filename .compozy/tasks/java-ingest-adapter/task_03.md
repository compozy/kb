---
status: completed
title: Implement Java adapter MVP parsing pipeline
type: backend
complexity: high
dependencies:
  - task_01
  - task_02
---

# Task 03: Implement Java adapter MVP parsing pipeline

## Overview
Create the first working Java adapter that parses Java source files into the graph model with symbols, imports, relations, and diagnostics. This task delivers the foundational parser behavior required before pipeline registration and deep resolution enhancements.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- A new `JavaAdapter` MUST implement `models.LanguageAdapter` and `ParseFilesWithProgress`.
- The adapter MUST produce deterministic `ParsedFile` outputs for Java files sorted by relative path.
- The adapter MUST emit structured parse diagnostics with Java-specific diagnostic codes on parse failure.
- The adapter MUST extract core Java symbols and imports sufficient for MVP ingest artifacts.
- Adapter tests MUST cover happy path parsing, error diagnostics, and relation emission basics.
</requirements>

## Subtasks
- [x] 3.1 Add `internal/adapter/java_adapter.go` with adapter contract implementation.
- [x] 3.2 Implement Java file parse flow to emit file, symbol, external node, and relation data.
- [x] 3.3 Add Java parse diagnostic behavior for syntax errors and nil tree/root edge cases.
- [x] 3.4 Add unit tests for Java adapter symbol and import extraction.
- [x] 3.5 Add integration test fixture coverage for multi-file Java relation behavior.
- [x] 3.6 Update existing non-Java adapter tests to assert Java is not supported by them.

## Implementation Details
Implement Java adapter behavior in the same style as existing adapters (`go_adapter`, `ts_adapter`, `rust_adapter`) with deterministic ordering, parser lifecycle management, and structured diagnostics. Keep the MVP scope aligned with TechSpec sections “Core Interfaces” and “Data Models.”

### Relevant Files
- `internal/adapter/go_adapter.go` — reference implementation for deterministic parse flow and diagnostics.
- `internal/adapter/ts_adapter.go` — reference for richer import/relation handling patterns.
- `internal/adapter/rust_adapter.go` — reference for multi-file resolution scaffolding pattern.
- `internal/models/models.go` — graph and diagnostic model contracts consumed by adapter output.
- `internal/adapter/java_adapter.go` — new Java adapter implementation target.
- `internal/adapter/java_adapter_test.go` — new unit test surface for Java adapter behavior.
- `internal/adapter/java_adapter_integration_test.go` — new integration test surface.

### Dependent Files
- `internal/generate/generate.go` — will register and invoke Java adapter in parse stage.
- `internal/graph/normalize.go` — consumes adapter `ParsedFile` output for merged graph.
- `internal/vault/render.go` — downstream rendering relies on adapter output shape.

### Related ADRs
- [ADR-001: Adopt a balanced MVP strategy for Java codebase ingest](../adrs/adr-001.md) — requires broad, practical MVP extraction quality.
- [ADR-002: Use deep Java relation resolution with safe syntactic fallback](../adrs/adr-002.md) — this task establishes base adapter behavior for later deep-resolution work.

## Deliverables
- New `internal/adapter/java_adapter.go` implementing MVP Java parse flow.
- New unit and integration tests for Java adapter behavior.
- Updated adapter support/rejection tests where language matrices include Java.
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for Java adapter file-to-file behavior **(REQUIRED)**

## Tests
- Unit tests:
  - [x] `JavaAdapter.Supports` accepts `models.LangJava` and rejects non-Java languages.
  - [x] Java class/method/package sources produce expected symbol kinds and signatures.
  - [x] Java imports produce expected `RelImports` edges and external nodes.
  - [x] Invalid Java source emits `JAVA_PARSE_ERROR` diagnostic with `StageParse`.
  - [x] Progress callback reports one tick per parsed Java file.
- Integration tests:
  - [x] Multi-file Java fixture emits cross-file relations for common call/import paths.
  - [x] Adapter output remains deterministic across repeated parse runs.
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- Java adapter produces valid `ParsedFile` output consumable by existing normalization/render stages
- Parse errors are represented as structured diagnostics rather than hard pipeline crashes
