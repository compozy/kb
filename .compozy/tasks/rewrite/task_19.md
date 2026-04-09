---
status: completed
title: End-to-end integration test
type: test
complexity: medium
dependencies:
    - task_12
    - task_16
    - task_18
---

# Task 19: End-to-end integration test

> **Reference implementation:** `~/dev/projects/kodebase` (original TypeScript kodebase) — use this repository as the behavioral and structural source when implementing this Go port.

## Overview

Create an end-to-end integration test that exercises the complete pipeline: generate a vault from a fixture project, then run inspect commands against it. This validates that all components work together correctly and catches integration issues that unit tests miss.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE IMPLEMENTATION — `~/dev/projects/kodebase` (original TypeScript kodebase) is the behavioral and structural source for this Go rewrite; align behavior and structure when porting
- REFERENCE TECHSPEC Phase 9.1 for test fixture requirements — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST create a fixture Go project in testdata/ with: multiple packages, circular imports, dead exports, complex functions
- MUST use `//go:build integration` build tag on test files
- MUST generate a vault from the fixture project
- MUST verify vault directory structure and document counts
- MUST run inspect subcommands against the generated vault and verify output
- MUST co-locate integration tests with the package they test
- MUST complete within 30 seconds
</requirements>

## Subtasks
- [ ] 19.1 Create fixture Go project in testdata/ with required characteristics
- [ ] 19.2 Write integration test that generates a vault from the fixture
- [ ] 19.3 Verify vault structure: directories, file counts, document types
- [ ] 19.4 Run inspect commands against the vault and verify output correctness
- [ ] 19.5 Ensure test runs within 30-second time limit

## Implementation Details

Create fixture project under `testdata/fixture-project/` with multiple Go packages demonstrating various code patterns. Write integration test in the appropriate package (likely alongside the generate orchestrator).

Reference TechSpec Phase 9.1 for fixture requirements.

### Relevant Files
- `internal/cli/generate.go` — generate command to invoke
- `internal/cli/inspect.go` — inspect commands to verify
- `internal/vault/writer.go` — writes vault that will be inspected
- `internal/vault/reader.go` — reads vault documents for verification
- `Makefile` — `make test-integration` target runs with `-tags integration`

### Dependent Files
- `magefile.go` — test-integration target must include this test

## Deliverables
- `testdata/fixture-project/` with multi-package Go fixture
- Integration test file with `//go:build integration` tag
- Tests verify end-to-end pipeline correctness **(REQUIRED)**

## Tests
- Integration tests:
  - [ ] Generate vault from fixture project completes without error
  - [ ] Vault contains raw/ directory with file and symbol documents
  - [ ] Vault contains wiki/ directory with 10 concept articles and 3 index pages
  - [ ] Vault contains bases/ directory with 12 base definition files
  - [ ] Inspect smells against generated vault returns expected code smells
  - [ ] Inspect dead-code against generated vault detects dead exports in fixture
  - [ ] Inspect circular-deps against generated vault detects the circular import
  - [ ] Inspect complexity against generated vault lists complex functions
  - [ ] Inspect blast-radius returns non-empty results for core symbols
  - [ ] Full pipeline completes within 30 seconds
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- `make test-integration` passes
- `make verify` passes
- Pipeline processes fixture in under 30 seconds
- All major pipeline features are exercised (scan, parse, normalize, metrics, render, write, read, inspect)
