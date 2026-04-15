---
status: completed
title: Validate Java ingest end-to-end with CLI and benchmark
type: test
complexity: high
dependencies:
  - task_04
  - task_05
---

# Task 06: Validate Java ingest end-to-end with CLI and benchmark

## Overview
Validate Java ingest behavior end-to-end through CLI integration and benchmark evidence aligned to the performance budget. This task proves that Java support is production-usable across workflow, artifact generation, and non-functional constraints.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- The codebase MUST include CLI E2E integration coverage for `kb ingest codebase` on Java multi-module fixtures.
- Validation MUST confirm Java appears in detected language summaries and output artifacts are written in expected codebase paths.
- Benchmark coverage MUST evaluate Java ingest runtime against agreed baselines and enforce <=20% overhead target.
- The workflow MUST continue to pass lint/inspect compatibility checks on generated Java ingest content.
</requirements>

## Subtasks
- [x] 6.1 Add Java fixture builder helpers for CLI integration tests.
- [x] 6.2 Add/extend CLI integration test to run Java multi-module ingest end-to-end.
- [x] 6.3 Validate generated artifacts and summary fields include Java-specific expectations.
- [x] 6.4 Add benchmark scenario for Java ingest performance budget tracking.
- [x] 6.5 Run full verification and benchmark commands, documenting budget compliance evidence.

## Implementation Details
Extend existing integration patterns used for Rust and Go workflows in CLI tests, and add benchmark validation aligned with TechSpec “Benchmark and E2E Validation” and ADR-004 acceptance criteria. Keep fixtures deterministic and minimal while still representing multi-module Java structure.

### Relevant Files
- `internal/cli/workflow_integration_test.go` — existing E2E ingest workflow tests for Go/Rust patterns.
- `internal/generate/generate_integration_test.go` — integration patterns for fixture-based ingest validation.
- `internal/adapter/java_adapter_integration_test.go` — relation correctness fixtures reused by E2E assertions.
- `internal/generate/generate.go` — source of summary fields validated in E2E outputs.
- `internal/models/models.go` — summary model fields referenced by validation assertions.

### Dependent Files
- `.compozy/tasks/java-ingest-adapter/_techspec.md` — acceptance thresholds and test strategy source of truth.
- `internal/lint/lint.go` — generated output lint compatibility is validated in E2E flow.
- `internal/cli/ingest_codebase.go` — command behavior exercised by new integration test.

### Related ADRs
- [ADR-003: Enforce 20% ingest performance budget with hybrid caching strategy](../adrs/adr-003.md) — defines runtime acceptance gate.
- [ADR-004: Require unit, integration, benchmark, and CLI E2E validation for Java ingest](../adrs/adr-004.md) — mandates this task’s validation scope.

## Deliverables
- New or updated CLI integration test covering Java multi-module ingest workflow.
- Java fixture generation helpers for deterministic E2E scenarios.
- Benchmark scenario and evidence for Java ingest runtime budget compliance.
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for end-to-end Java ingest workflow **(REQUIRED)**

## Tests
- Unit tests:
  - [x] Java fixture helper creates deterministic module/package structure expected by tests.
  - [x] Summary assertion helpers validate Java language presence and artifact counts correctly.
- Integration tests:
  - [x] `kb ingest codebase` on Java multi-module fixture returns success and includes `java` in detected languages.
  - [x] Generated topic contains expected Java file and symbol markdown artifacts.
  - [x] `kb lint` on generated output reports zero blocking issues for Java ingest content.
  - [x] Benchmark run verifies Java ingest stays within <=20% overhead against baseline fixture.
- Test coverage target: >=80%
- All tests must pass

## Validation Evidence
- `go test ./internal/cli -run "TestWriteJavaMultiModuleCodebaseFixtureCreatesDeterministicLayout|TestValidateJavaCodebaseSummary|TestAssertJavaCodebaseSummaryPassesForValidInput" -count=1` passed.
- `go test -tags integration ./internal/cli -run "TestCLIIntegrationScaffoldIngestJavaWorkspaceCodebase" -count=1` passed.
- `go test -tags integration ./internal/generate -run "TestGenerateIntegrationJavaIngestPerformanceBudget" -count=1` passed.
- `go test -tags integration ./internal/generate -run "^$" -bench "BenchmarkGenerateIntegration(GoBaselineDryRun|JavaDryRun)" -benchmem -count=1` passed with Java `3793232 ns/op` vs baseline Go `3388442 ns/op` (~11.95% overhead, within <=20% budget).
- `make verify` passed (fmt, lint, test, build, boundaries).

## Success Criteria
- All tests passing
- Test coverage >=80%
- Java ingest workflow is validated end-to-end through CLI integration tests
- Benchmark evidence confirms performance budget compliance for Java ingest
