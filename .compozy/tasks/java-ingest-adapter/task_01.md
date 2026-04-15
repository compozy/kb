---
status: completed
title: Add Java language support to models and scanner
type: backend
complexity: medium
dependencies: []
---

# Task 01: Add Java language support to models and scanner

## Overview
Add Java as a first-class supported language in the domain model and workspace scanner so Java files can enter the ingest pipeline. This task establishes the minimum language registration surface required by downstream adapter and generator tasks.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- The system MUST add `LangJava` to `models.SupportedLanguage` constants and expose it through `SupportedLanguages()` and `SupportedLanguageNames()`.
- The scanner MUST classify `.java` files as `models.LangJava` while preserving current behavior for existing languages.
- Existing deterministic ordering assumptions for supported languages MUST remain stable after adding Java.
- Unit tests MUST cover model language lists and scanner language detection for Java.
</requirements>

## Subtasks
- [x] 1.1 Add Java language constant and include it in supported language lists.
- [x] 1.2 Extend scanner extension mapping to recognize `.java` files.
- [x] 1.3 Update model tests to assert Java appears in supported language outputs.
- [x] 1.4 Update scanner tests to assert Java file discovery and grouping behavior.
- [x] 1.5 Run targeted package tests for `internal/models` and `internal/scanner`.

## Implementation Details
Update the language registry in models and the file extension mapping in scanner so Java files are scanned and grouped correctly. Keep compatibility with current code paths that rely on ordered language slices for reporting and adapter selection.

### Relevant Files
- `internal/models/models.go` — source of supported language constants and list builders.
- `internal/models/models_test.go` — assertions for supported language lists.
- `internal/scanner/scanner.go` — extension-to-language mapping in `supportedLanguage`.
- `internal/scanner/scanner_test.go` — scanner behavior and grouped language test patterns.

### Dependent Files
- `internal/generate/generate.go` — consumes model language lists during adapter selection.
- `internal/cli/generate.go` — renders supported language help from model names.
- `internal/cli/ingest_codebase.go` — includes generated supported language help text.

### Related ADRs
- [ADR-001: Adopt a balanced MVP strategy for Java codebase ingest](../adrs/adr-001.md) — Java must be recognized as a standard language in the same workflow.

## Deliverables
- Java language constant and list registration in `internal/models`.
- Java extension mapping in `internal/scanner`.
- Updated unit tests in models and scanner packages.
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for scanner language grouping behavior **(REQUIRED)**

## Tests
- Unit tests:
  - [x] `SupportedLanguages()` includes `LangJava` in deterministic order.
  - [x] `SupportedLanguageNames()` includes `java` and preserves stable ordering.
  - [x] `supportedLanguage("src/Foo.java")` returns `models.LangJava`.
  - [x] Existing extension mappings (`.go`, `.rs`, `.ts`, `.tsx`, `.js`, `.jsx`) remain unchanged.
- Integration tests:
  - [x] Workspace scan with mixed file types includes Java files in `Files`.
  - [x] Workspace scan groups Java files under `FilesByLanguage[models.LangJava]`.
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- Java files are discoverable and grouped correctly by scanner output
- Model language registries expose Java consistently for downstream components
