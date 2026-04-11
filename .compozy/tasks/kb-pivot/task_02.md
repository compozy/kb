---
status: completed
title: Extend domain models with KB types
type: backend
complexity: medium
dependencies: []
---

# Task 02: Extend domain models with KB types

## Overview

Extend `internal/models/` with new domain types required by the knowledge base pivot: converter interface, ingest result, lint issue, topic management types, and source kind constants. These types establish the shared vocabulary used across all new packages.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST define the `Converter` interface as specified in TechSpec "Core Interfaces"
- MUST define `ConvertInput`, `ConvertResult`, `IngestResult`, and `LintIssue` structs
- MUST define `SourceKind` string type with constants: article, github-readme, youtube-transcript, codebase-file, codebase-symbol, bookmark-cluster, document
- MUST define `LintIssueKind` string type with constants: dead-link, orphan, missing-source, stale, format
- MUST define `TopicInfo` struct for topic listing/info operations
- SHOULD keep new types in a separate file (`kb_models.go`) to minimize merge conflicts with existing models
</requirements>

## Subtasks

- [x] 2.1 Create `internal/models/kb_models.go` with all new types and constants
- [x] 2.2 Define the Converter interface and its input/output types
- [x] 2.3 Define ingest domain types (IngestResult, SourceKind constants)
- [x] 2.4 Define lint domain types (LintIssue, LintIssueKind constants)
- [x] 2.5 Define topic management types (TopicInfo)
- [x] 2.6 Write unit tests for constants and type helper functions

## Implementation Details

Add a new file `internal/models/kb_models.go` alongside the existing `models.go`. Keep existing types untouched.

Reference TechSpec "Core Interfaces" section for the Converter, ConvertInput, ConvertResult, IngestResult, and LintIssue type definitions.

### Relevant Files

- `internal/models/models.go` — existing domain types; new types must not conflict with existing names
- `internal/models/models_test.go` — existing test patterns to follow

### Dependent Files

- `internal/convert/` (task_05) — implements the Converter interface
- `internal/ingest/` (task_12) — uses IngestResult
- `internal/lint/` (task_13) — uses LintIssue, LintIssueKind
- `internal/topic/` (task_04) — uses TopicInfo

### Related ADRs

- [ADR-003: Native Go Document Conversion with Converter Registry](../adrs/adr-003.md) — defines the Converter interface pattern

## Deliverables

- `internal/models/kb_models.go` — all new types and constants
- `internal/models/kb_models_test.go` — tests
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests

- Unit tests:
  - [x] All SourceKind constants are unique non-empty strings
  - [x] All LintIssueKind constants are unique non-empty strings
  - [x] ConvertResult fields are accessible and zero-valued by default
  - [x] IngestResult fields are accessible and zero-valued by default
  - [x] LintIssue fields are accessible and zero-valued by default
- Test coverage target: >=80%
- All tests must pass

## Success Criteria

- All tests passing
- Test coverage >=80%
- New types compile cleanly alongside existing models
- `make lint` reports zero findings
