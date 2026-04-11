---
status: completed
title: Adapt codebase generate pipeline
type: refactor
complexity: high
dependencies:
  - task_04
  - task_12
---

# Task 14: Adapt codebase generate pipeline

## Overview

Modify the existing `internal/generate/` pipeline so that its output writes to `<topic>/raw/codebase/` instead of the top-level vault. The `generate` command becomes the implementation behind `ingest codebase`, preserving the full analysis pipeline (scan → parse → graph → metrics → render) while integrating with the topic-based vault structure.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST modify `GenerateOptions` to accept a topic slug and vault path instead of a standalone output path
- MUST write rendered documents (raw file/symbol pages) to `<topic>/raw/codebase/{files,symbols}/`
- MUST write wiki concept pages to `<topic>/wiki/concepts/` (codebase-specific wiki articles)
- MUST write index pages to `<topic>/wiki/index/`
- MUST write base files to `<topic>/bases/`
- MUST append `## [YYYY-MM-DD] ingest | codebase (<file_count> files, <symbol_count> symbols)` to topic log.md
- MUST preserve the existing `GenerationSummary` output structure
- MUST NOT break the existing vault writer's core logic — adapt paths, don't rewrite the writer
- SHOULD update `TopicMetadata` to be derived from topic + vault path instead of generate options
</requirements>

## Subtasks

- [x] 14.1 Modify `GenerateOptions` and `TopicMetadata` to work with topic-based paths
- [x] 14.2 Update vault writer output paths to write under `<topic>/raw/codebase/`
- [x] 14.3 Update rendered document paths (raw documents go to `raw/codebase/`, wiki stays in `wiki/`)
- [x] 14.4 Add log.md append for codebase ingest completion
- [x] 14.5 Ensure inspect commands still work against the new vault layout
- [x] 14.6 Update existing generate tests to verify new path structure

## Implementation Details

The key change is in how `TopicMetadata` paths are constructed and how `WriteVault` places files. The existing vault writer already manages `raw/codebase/` as a managed subtree — the adaptation is making it write there as the primary output location for codebase analysis.

The generate pipeline itself (scan → parse → graph → metrics → render) is unchanged. Only the write phase and the options/metadata construction change.

### Relevant Files

- `internal/generate/generate.go` — Generate(), GenerateWithObserver(), runner struct
- `internal/generate/events.go` — Observer pattern (unchanged)
- `internal/vault/writer.go` — WriteVault(), WriteVaultOptions, ensureTopicSkeleton()
- `internal/vault/render.go` — RenderDocuments() output paths
- `internal/models/models.go` — GenerateOptions, TopicMetadata, GenerationSummary
- `internal/topic/` (task_04) — topic path resolution

### Dependent Files

- `internal/cli/` (task_16) — `ingest codebase` command wires to this adapted pipeline
- `internal/vault/reader.go` — vault reader may need path adjustments for the new layout
- `internal/cli/inspect_*.go` — inspect commands must resolve vault paths under the topic structure

### Related ADRs

- [ADR-001: Topic-Centric CLI Command Taxonomy](../adrs/adr-001.md) — codebase as ingest source, not top-level command

## Deliverables

- Modified `internal/generate/generate.go` — topic-aware generation
- Modified `internal/vault/writer.go` — topic-based output paths
- Modified `internal/models/models.go` — updated GenerateOptions/TopicMetadata if needed
- Updated test files
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests

- Unit tests:
  - [x] Generate with topic slug writes raw documents to `<topic>/raw/codebase/files/`
  - [x] Generate with topic slug writes symbol documents to `<topic>/raw/codebase/symbols/`
  - [x] Generate with topic slug writes wiki concepts to `<topic>/wiki/concepts/`
  - [x] Generate with topic slug writes index pages to `<topic>/wiki/index/`
  - [x] Generate with topic slug writes base files to `<topic>/bases/`
  - [x] Generate appends ingest log entry to `<topic>/log.md`
  - [x] GenerationSummary structure is preserved (file count, symbol count, timings)
- Integration tests:
  - [x] Full generate pipeline against fixture Go project → verify all files in correct topic subtree
- Test coverage target: >=80%
- All tests must pass

## Success Criteria

- All tests passing
- Test coverage >=80%
- Codebase generation output is correctly placed under topic directory
- Existing inspect commands can find and read the new vault layout
- `make verify` passes
