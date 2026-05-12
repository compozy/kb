---
status: completed
title: Register Java adapter in generate runner
type: backend
complexity: medium
dependencies:
  - task_03
---

# Task 04: Register Java adapter in generate runner

## Overview
Wire the Java adapter into the codebase generation pipeline so Java files discovered by scanner are parsed during ingest runs. This task connects the implemented adapter to orchestration without changing stage contracts or CLI command shape.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- `internal/generate/newRunner()` MUST include `adapter.JavaAdapter{}` in adapter registration.
- `runner.withDefaults()` MUST include `adapter.JavaAdapter{}` in fallback adapter list.
- Adapter selection behavior MUST remain deterministic for mixed-language workspaces.
- Existing generate and CLI help tests MUST be updated if language list outputs change.
</requirements>

## Subtasks
- [x] 4.1 Add Java adapter to runner adapter list in `newRunner`.
- [x] 4.2 Add Java adapter to fallback list in `withDefaults`.
- [x] 4.3 Update generate tests to validate Java adapter selection when Java language is present.
- [x] 4.4 Validate CLI language help snapshots that derive from supported language names.
- [x] 4.5 Run generate and CLI targeted test suites.

## Implementation Details
Modify only orchestration registration points and associated tests. Keep stage ordering, event reporting, and dry-run behavior unchanged while enabling Java parse participation through existing adapter selection logic.

### Relevant Files
- `internal/generate/generate.go` — adapter registration and default runner configuration.
- `internal/generate/generate_test.go` — adapter selection and runner orchestration tests.
- `internal/cli/generate.go` — supported language help rendering from model names.
- `internal/cli/ingest_codebase.go` — reuses supported language help in command description.
- `internal/cli/generate_test.go` — generate command help coverage.
- `internal/cli/ingest_test.go` — ingest command help coverage.

### Dependent Files
- `internal/cli/workflow_integration_test.go` — end-to-end ingest behavior depends on registration.
- `internal/generate/generate_integration_test.go` — integration results depend on selected adapters.

### Related ADRs
- [ADR-001: Adopt a balanced MVP strategy for Java codebase ingest](../adrs/adr-001.md) — Java must run through the same ingest orchestration path.

## Deliverables
- Java adapter registered in generate runner and runner defaults.
- Updated generate/CLI tests reflecting Java-aware adapter and language lists.
- Evidence that mixed-language adapter selection remains deterministic.
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for runner selection flow **(REQUIRED)**

## Tests
- Unit tests:
  - [x] `selectAdapters` includes Java adapter when `LangJava` appears in workspace languages.
  - [x] `newRunner` adapter list includes Java adapter in expected order.
  - [x] `withDefaults` populates Java adapter when custom list is empty.
  - [x] Command help output includes Java in supported language text.
- Integration tests:
  - [x] Generate runner dry-run with Java-scanned workspace reports Java in detected/selected outputs.
  - [x] Existing non-Java integration paths remain unchanged and passing.
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- Java files trigger Java adapter selection during generation
- CLI help and summary outputs consistently expose Java support
