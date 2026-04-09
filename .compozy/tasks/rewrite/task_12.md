---
status: completed
title: Wire generate command end-to-end
type: backend
complexity: high
dependencies:
  - task_03
  - task_05
  - task_06
  - task_08
  - task_11
---

# Task 12: Wire generate command end-to-end

> **Reference implementation:** `~/dev/projects/kodebase` (original TypeScript kodebase) — use this repository as the behavioral and structural source when implementing this Go port.

## Overview

Connect all pipeline components — scanner, adapters, graph normalizer, metrics engine, document renderer, and vault writer — into the `generate` CLI command. This is the primary user-facing command that transforms a source repository into an Obsidian knowledge vault.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE IMPLEMENTATION — `~/dev/projects/kodebase` (original TypeScript kodebase) is the behavioral and structural source for this Go rewrite; align behavior and structure when porting
- REFERENCE TECHSPEC Phase 6.1 for pipeline sequence and CLI flags — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST implement Generate orchestrator function: scan → select adapters → parse → normalize → compute metrics → render → write
- MUST accept GenerateOptions with: source path, output path, topic name, title, domain, include/exclude patterns, semantic flag
- MUST select correct adapters based on languages found in scan results
- MUST return GenerationSummary with file/symbol counts and timing
- MUST wire CLI flags: --output, --topic, --title, --domain, --include, --exclude, --semantic
- MUST use context.Context for cancellation support
- MUST log pipeline progress via slog at key stages
</requirements>

## Subtasks
- [x] 12.1 Implement Generate orchestrator function with full pipeline sequence
- [x] 12.2 Implement adapter selection based on scanned languages
- [x] 12.3 Wire CLI flags on the generate cobra command
- [x] 12.4 Implement GenerationSummary reporting
- [x] 12.5 Add progress logging at each pipeline stage
- [x] 12.6 Write integration test against fixture project

## Implementation Details

Modify `internal/cli/generate.go` to replace the stub with a real implementation. Create a pipeline orchestrator — this can live in the cli package or a separate package depending on complexity. Reference `~/dev/projects/kodebase/packages/cli/src/knowledge-base/generate-knowledge-base.ts` (114 lines).

The pipeline sequence: ScanWorkspace → select adapters → ParseFiles (per adapter) → NormalizeGraph → ComputeMetrics → RenderDocuments → WriteVault.

### Relevant Files
- `~/dev/projects/kodebase/packages/cli/src/knowledge-base/generate-knowledge-base.ts` — TypeScript source (114 lines)
- `internal/cli/generate.go` — stub from task_01, replace with real command
- `internal/scanner/scanner.go` — ScanWorkspace (step 1)
- `internal/adapter/go_adapter.go` — GoAdapter (step 2-3)
- `internal/adapter/ts_adapter.go` — TSAdapter (step 2-3)
- `internal/graph/normalize.go` — NormalizeGraph (step 4)
- `internal/metrics/compute.go` — ComputeMetrics (step 5)
- `internal/vault/render.go` — RenderDocuments (step 6)
- `internal/vault/writer.go` — WriteVault (step 7)

### Dependent Files
- `internal/models/models.go` — GenerateOptions, GenerationSummary types

## Deliverables
- Updated `internal/cli/generate.go` with real command implementation
- Generate orchestrator function (location TBD based on CLAUDE.md conventions)
- Integration test against a fixture Go project
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests
- Unit tests:
  - [x] Generate orchestrator calls pipeline stages in correct order (mock adapters)
  - [x] Adapter selection picks GoAdapter for Go-only projects
  - [x] Adapter selection picks both adapters for mixed Go + TS projects
  - [x] GenerationSummary reports correct file and symbol counts
  - [x] Missing source path returns descriptive error
- Integration tests:
  - [x] Generate vault from a small fixture Go project in testdata/
  - [x] Verify vault directory structure exists after generation
  - [x] Verify raw documents, wiki articles, and base files are created
  - [x] Verify GenerationSummary counts match expected fixture totals
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- `make verify` passes
- `kodebase generate <path>` produces a valid Obsidian vault
- Pipeline is cancellable via context
