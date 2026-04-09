---
status: completed
title: Define domain models
type: backend
complexity: medium
dependencies:
    - task_01
---

# Task 02: Define domain models

> **Reference implementation:** `~/dev/projects/kodebase` (original TypeScript kodebase) — use this repository as the behavioral and structural source when implementing this Go port.

## Overview

Port all TypeScript interfaces and types from the reference `models.ts` to Go structs and type constants in `internal/models/`. This package defines the shared vocabulary (graph nodes, relations, metrics, documents) used by every other package in the pipeline.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE IMPLEMENTATION — `~/dev/projects/kodebase` (original TypeScript kodebase) is the behavioral and structural source for this Go rewrite; align behavior and structure when porting
- REFERENCE TECHSPEC Phase 1.1 for the full type list — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST port all types listed in TechSpec Phase 1.1 to Go equivalents
- MUST use string-typed constants (not iota) for enum-like types: SupportedLanguage, RelationType, RelationConfidence, DiagnosticSeverity, DiagnosticStage, DocumentKind, ManagedArea, BaseViewType
- MUST define LanguageAdapter interface: `Supports(lang SupportedLanguage) bool`, `ParseFiles(files []ScannedSourceFile, rootPath string) ([]ParsedFile, error)`
- MUST provide `SupportedLanguages()` helper returning all language constants
- MUST use compile-time interface verification pattern: `var _ Interface = (*Type)(nil)` where applicable
- MUST NOT use `interface{}` or `any` when concrete types are known
- SHOULD use `json` struct tags for serialization compatibility
</requirements>

## Subtasks
- [ ] 2.1 Define string-typed constants for all enum types (SupportedLanguage, RelationType, etc.)
- [ ] 2.2 Define graph node structs (GraphFile, SymbolNode, ExternalNode, RelationEdge)
- [ ] 2.3 Define container structs (ParsedFile, GraphSnapshot, ScannedSourceFile, ScannedWorkspace)
- [ ] 2.4 Define output structs (RenderedDocument, TopicMetadata, GenerateOptions, GenerationSummary)
- [ ] 2.5 Define metrics structs (SymbolMetrics, FileMetrics, DirectoryMetrics, MetricsResult)
- [ ] 2.6 Define LanguageAdapter interface and Base/filter structs

## Implementation Details

Single file `internal/models/models.go` containing all types. Reference the TypeScript source at `~/dev/projects/kodebase/packages/cli/src/knowledge-base/models.ts` (220 lines) for exact field names and relationships.

The LanguageAdapter interface is consumed by adapter implementations (task_05, task_06) and the generate orchestrator (task_12).

### Relevant Files
- `~/dev/projects/kodebase/packages/cli/src/knowledge-base/models.ts` — TypeScript source (220 lines), authoritative type definitions

### Dependent Files
- `internal/scanner/scanner.go` — uses ScannedSourceFile, ScannedWorkspace, SupportedLanguage
- `internal/adapter/go_adapter.go` — implements LanguageAdapter, uses ParsedFile, GraphFile, SymbolNode
- `internal/adapter/ts_adapter.go` — implements LanguageAdapter, uses ParsedFile, GraphFile, SymbolNode
- `internal/graph/normalize.go` — uses ParsedFile, GraphSnapshot
- `internal/metrics/compute.go` — uses GraphSnapshot, SymbolMetrics, FileMetrics, MetricsResult
- `internal/vault/render.go` — uses GraphSnapshot, MetricsResult, RenderedDocument

## Deliverables
- `internal/models/models.go` with all domain types
- `internal/models/models_test.go` with type verification tests
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests
- Unit tests:
  - [ ] `SupportedLanguages()` returns all 5 language constants (ts, tsx, js, jsx, go)
  - [ ] All enum-type constants have non-empty string values
  - [ ] GraphFile struct can be instantiated with all required fields
  - [ ] SymbolNode struct correctly stores cyclomatic complexity and export flag
  - [ ] ParsedFile aggregates file, symbols, external nodes, relations, and diagnostics
  - [ ] GraphSnapshot aggregates all node types with root path
  - [ ] MetricsResult maps are non-nil after zero-value construction helper
- Integration tests:
  - [ ] JSON round-trip serialization/deserialization of GraphSnapshot preserves all fields
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- `make verify` passes
- All types from TechSpec Phase 1.1 are present
- No `interface{}` or `any` usage
