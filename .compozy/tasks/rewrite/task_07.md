---
status: completed
title: Implement graph normalizer
type: backend
complexity: medium
dependencies:
  - task_02
---

# Task 07: Implement graph normalizer

> **Reference implementation:** `~/dev/projects/kodebase` (original TypeScript kodebase) — use this repository as the behavioral and structural source when implementing this Go port.

## Overview

Port `normalize-graph.ts` to Go — merge multiple `ParsedFile` results from language adapters into a single deduplicated, deterministically ordered `GraphSnapshot`. This is the data fusion step between parsing and metrics computation.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE IMPLEMENTATION — `~/dev/projects/kodebase` (original TypeScript kodebase) is the behavioral and structural source for this Go rewrite; align behavior and structure when porting
- REFERENCE TECHSPEC Phase 4.1 for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST accept a root path and slice of ParsedFile, return a GraphSnapshot
- MUST deduplicate files by ID
- MUST deduplicate symbols by ID
- MUST deduplicate external nodes by ID
- MUST deduplicate relations by composite key (fromId + toId + type + confidence)
- MUST attach symbol IDs to their parent GraphFile's SymbolIds field
- MUST order all collections deterministically (sorted by ID or composite key)
- MUST order diagnostics by stage, then filePath, then message
</requirements>

## Subtasks
- [x] 7.1 Implement NormalizeGraph function signature and basic aggregation
- [x] 7.2 Implement deduplication logic for files, symbols, external nodes, and relations
- [x] 7.3 Implement symbol-to-file attachment (populate GraphFile.SymbolIds)
- [x] 7.4 Implement deterministic ordering for all collections
- [x] 7.5 Implement diagnostic ordering by stage/filePath/message

## Implementation Details

Create `internal/graph/normalize.go` and `internal/graph/normalize_test.go`. Reference `~/dev/projects/kodebase/packages/cli/src/knowledge-base/normalize-graph.ts` (114 lines) — this is a compact but critical transformation step.

Use maps for deduplication, then sort the results into slices for deterministic output.

### Relevant Files
- `~/dev/projects/kodebase/packages/cli/src/knowledge-base/normalize-graph.ts` — TypeScript source (114 lines)
- `internal/models/models.go` — ParsedFile, GraphSnapshot, GraphFile, SymbolNode, ExternalNode, RelationEdge

### Dependent Files
- `internal/metrics/compute.go` — receives GraphSnapshot as input
- `internal/vault/render.go` — receives GraphSnapshot as input
- `internal/cli/generate.go` — calls NormalizeGraph in the pipeline

## Deliverables
- `internal/graph/normalize.go` with NormalizeGraph function
- `internal/graph/normalize_test.go` with comprehensive tests
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests
- Unit tests:
  - [x] Empty ParsedFile slice returns empty GraphSnapshot with rootPath set
  - [x] Single ParsedFile passes through all nodes and relations unchanged
  - [x] Duplicate files (same ID) are merged into one entry
  - [x] Duplicate symbols (same ID) are merged into one entry
  - [x] Duplicate relations (same composite key) are merged into one entry
  - [x] Symbols are attached to correct parent file via SymbolIds
  - [x] Output files are sorted deterministically by ID
  - [x] Output relations are sorted deterministically by composite key
  - [x] Diagnostics are ordered by stage, then filePath, then message
- Integration tests:
  - [x] Normalize output from two parsed files with overlapping imports produces correct merged graph
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- `make verify` passes
- Output is fully deterministic (same input always produces same output byte-for-byte)
