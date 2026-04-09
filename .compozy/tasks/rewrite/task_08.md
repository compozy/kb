---
status: completed
title: Implement metrics engine
type: backend
complexity: high
dependencies:
  - task_02
  - task_07
---

# Task 08: Implement metrics engine

> **Reference implementation:** `~/dev/projects/kodebase` (original TypeScript kodebase) — use this repository as the behavioral and structural source when implementing this Go port.

## Overview

Port `compute-metrics.ts` to Go — compute symbol-level, file-level, and directory-level code metrics from a GraphSnapshot. Includes blast radius (BFS), centrality (PageRank-like), coupling, instability, dead export detection, code smell identification, and circular dependency detection (DFS).

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE IMPLEMENTATION — `~/dev/projects/kodebase` (original TypeScript kodebase) is the behavioral and structural source for this Go rewrite; align behavior and structure when porting
- REFERENCE TECHSPEC Phase 4.2 for the full metric catalog — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST compute per-symbol metrics: blast radius (BFS traversal), centrality, direct dependents count, external reference count, dead export detection, long function detection, code smell flags
- MUST compute per-file metrics: afferent coupling (Ca), efferent coupling (Ce), instability ratio (Ce/(Ca+Ce)), orphan file detection, god file detection, circular dependency participation
- MUST compute per-directory metrics: aggregated coupling and instability
- MUST detect circular dependencies using DFS cycle detection on the file-level import graph
- MUST return MetricsResult with maps keyed by symbol/file/directory ID
- MUST handle empty graphs gracefully (no panics, return empty MetricsResult)
- SHOULD detect entry point heuristics (main functions, command files)
</requirements>

## Subtasks
- [x] 8.1 Implement blast radius computation via BFS traversal of the relation graph
- [x] 8.2 Implement centrality and coupling metrics (afferent, efferent, instability)
- [x] 8.3 Implement dead export detection and code smell identification
- [x] 8.4 Implement circular dependency detection via DFS cycle finding
- [x] 8.5 Implement directory-level metric aggregation
- [x] 8.6 Assemble MetricsResult from all computed metrics

## Implementation Details

Create `internal/metrics/compute.go` and `internal/metrics/compute_test.go`. Reference `~/dev/projects/kodebase/packages/cli/src/knowledge-base/compute-metrics.ts` (480 lines) for algorithm details.

Build adjacency lists from GraphSnapshot relations, then traverse for BFS (blast radius) and DFS (cycles). Coupling metrics derive from relation counts between files.

### Relevant Files
- `~/dev/projects/kodebase/packages/cli/src/knowledge-base/compute-metrics.ts` — TypeScript source (480 lines)
- `internal/models/models.go` — GraphSnapshot, SymbolMetrics, FileMetrics, DirectoryMetrics, MetricsResult
- `internal/graph/normalize.go` — produces the GraphSnapshot input

### Dependent Files
- `internal/vault/render.go` — uses MetricsResult to generate wiki articles about hotspots, smells, etc.
- `internal/cli/generate.go` — calls ComputeMetrics in the pipeline
- `internal/vault/reader.go` — reads back metrics from vault documents

## Deliverables
- `internal/metrics/compute.go` with ComputeMetrics function
- `internal/metrics/compute_test.go` with comprehensive tests
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests
- Unit tests:
  - [x] Empty graph returns empty MetricsResult with no errors
  - [x] Single symbol with no relations has blast radius of 0
  - [x] Symbol imported by 3 files has blast radius >= 3
  - [x] File with 2 imports has efferent coupling of 2
  - [x] File imported by 4 files has afferent coupling of 4
  - [x] Instability ratio is Ce/(Ca+Ce) = 0.5 for balanced file
  - [x] Exported symbol with no dependents is flagged as dead export
  - [x] Function exceeding LOC threshold is flagged as long function
  - [x] A -> B -> C -> A cycle is detected by DFS
  - [x] Acyclic graph returns empty circular dependency list
  - [x] Directory metrics aggregate child file coupling correctly
- Integration tests:
  - [x] Compute metrics on a realistic multi-package GraphSnapshot and verify key metrics
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- `make verify` passes
- All metric types from TechSpec Phase 4.2 are computed
- Circular dependency detection handles both simple and transitive cycles
