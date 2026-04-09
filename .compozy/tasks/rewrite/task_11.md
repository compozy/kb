---
status: completed
title: Implement vault writer
type: backend
complexity: medium
dependencies:
  - task_10
---

# Task 11: Implement vault writer

> **Reference implementation:** `~/dev/projects/kodebase` (original TypeScript kodebase) — use this repository as the behavioral and structural source when implementing this Go port.

## Overview

Port `write-vault.ts` to Go — write the full Obsidian vault to disk. Creates the topic directory skeleton, writes rendered documents in batches, manages CLAUDE.md (topic manifest), manages log.md (append-only audit log), and removes stale managed wiki concepts.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE IMPLEMENTATION — `~/dev/projects/kodebase` (original TypeScript kodebase) is the behavioral and structural source for this Go rewrite; align behavior and structure when porting
- REFERENCE TECHSPEC Phase 5.4 for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST create topic directory skeleton: raw/, wiki/, outputs/, bases/
- MUST write all RenderedDocuments to correct subdirectories based on ManagedArea
- MUST generate YAML frontmatter + body format for each document file
- MUST manage CLAUDE.md as a topic manifest (vault metadata)
- MUST manage log.md as an append-only audit log of generation runs
- MUST detect and remove stale managed wiki concepts that are no longer generated
- MUST create parent directories as needed (no errors on missing intermediate dirs)
- MUST handle concurrent writes safely (batch writes, not parallel file writes)
</requirements>

## Subtasks
- [x] 11.1 Implement directory skeleton creation (raw/, wiki/, outputs/, bases/)
- [x] 11.2 Implement document writing with frontmatter + body formatting
- [x] 11.3 Implement CLAUDE.md topic manifest generation
- [x] 11.4 Implement log.md append-only audit log
- [x] 11.5 Implement stale wiki concept cleanup

## Implementation Details

Create `internal/vault/writer.go` and `internal/vault/writer_test.go`. Reference `~/dev/projects/kodebase/packages/cli/src/knowledge-base/write-vault.ts` (265 lines).

The writer receives a slice of RenderedDocument from the renderer (task_10) and writes them to the output directory. Documents are routed to subdirectories based on their ManagedArea field.

### Relevant Files
- `~/dev/projects/kodebase/packages/cli/src/knowledge-base/write-vault.ts` — TypeScript source (265 lines)
- `internal/models/models.go` — RenderedDocument, ManagedArea, DocumentKind
- `internal/vault/render.go` — produces RenderedDocument slice that this writer consumes
- `internal/vault/pathutils.go` — vault path derivation for output file paths

### Dependent Files
- `internal/cli/generate.go` — calls WriteVault as the final pipeline step
- `internal/vault/reader.go` — reads back the vault documents written by this writer

## Deliverables
- `internal/vault/writer.go` with WriteVault function
- `internal/vault/writer_test.go` with comprehensive tests
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests
- Unit tests:
  - [x] WriteVault creates correct directory skeleton in t.TempDir()
  - [x] Raw documents are written to raw/ subdirectory
  - [x] Wiki documents are written to wiki/ subdirectory
  - [x] Base documents are written to bases/ subdirectory
  - [x] Written files contain valid YAML frontmatter + markdown body
  - [x] CLAUDE.md is created with topic metadata
  - [x] log.md is created and contains generation timestamp
  - [x] Second write appends to log.md, does not overwrite
  - [x] Stale wiki files are removed on re-generation
  - [x] Parent directories are created automatically
- Integration tests:
  - [x] Write a full document set and verify all files exist with correct content
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- `make verify` passes
- Vault directory structure matches expected layout
- Written documents can be read by a markdown parser
