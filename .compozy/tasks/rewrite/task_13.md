---
status: completed
title: Implement vault reader and query resolver
type: backend
complexity: medium
dependencies:
  - task_02
  - task_09
---

# Task 13: Implement vault reader and query resolver

> **Reference implementation:** `~/dev/projects/kodebase` (original TypeScript kodebase) — use this repository as the behavioral and structural source when implementing this Go port.

## Overview

Port `vault-reader.ts` and `vault-query.ts` to Go — read generated vault documents back into structured data and resolve vault/topic paths for inspection commands. The reader parses YAML frontmatter and markdown bodies; the query resolver locates vaults and topics on disk.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE IMPLEMENTATION — `~/dev/projects/kodebase` (original TypeScript kodebase) is the behavioral and structural source for this Go rewrite; align behavior and structure when porting
- REFERENCE TECHSPEC Phases 7.1 and 7.2 for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST parse YAML frontmatter from vault markdown files
- MUST extract relation links from `- [[...]]` patterns in document bodies
- MUST resolve backlinks between documents
- MUST classify documents as symbols, files, directories, or wikis
- MUST extract specific sections from documents by heading name
- MUST implement vault path resolution: find vault root from current directory
- MUST implement topic resolution: find topic directory within vault
- MUST support symbol lookup by name (case-insensitive partial match)
- MUST handle missing vault/topic gracefully with descriptive errors
</requirements>

## Subtasks
- [x] 13.1 Implement vault document reader with frontmatter parsing
- [x] 13.2 Implement relation link extraction from markdown body
- [x] 13.3 Implement document classification (symbol, file, directory, wiki)
- [x] 13.4 Implement section extraction by heading name
- [x] 13.5 Implement vault and topic path resolution from working directory
- [x] 13.6 Implement symbol lookup by name with partial matching

## Implementation Details

Create `internal/vault/reader.go`, `internal/vault/reader_test.go`, `internal/vault/query.go`, and `internal/vault/query_test.go`. Reference:
- `~/dev/projects/kodebase/packages/cli/src/knowledge-base/vault-reader.ts` (207 lines)
- `~/dev/projects/kodebase/packages/cli/src/knowledge-base/vault-query.ts` (109 lines)

The reader uses regex-based parsing for frontmatter delimiters (`---`) and wiki-link patterns (`[[...]]`). The query resolver walks up the directory tree to find `.kodebase/` vault roots.

### Relevant Files
- `~/dev/projects/kodebase/packages/cli/src/knowledge-base/vault-reader.ts` — TypeScript source (207 lines)
- `~/dev/projects/kodebase/packages/cli/src/knowledge-base/vault-query.ts` — TypeScript source (109 lines)
- `internal/models/models.go` — RenderedDocument, DocumentKind types
- `internal/vault/pathutils.go` — path derivation and normalization (from task_09)

### Dependent Files
- `internal/cli/inspect.go` — all inspect subcommands use the reader and query resolver
- `internal/cli/search.go` — search command uses query resolver to find vault
- `internal/cli/index.go` — index command uses query resolver to find vault

## Deliverables
- `internal/vault/reader.go` with vault document reading functions
- `internal/vault/reader_test.go` with tests
- `internal/vault/query.go` with vault/topic resolution functions
- `internal/vault/query_test.go` with tests
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests
- Unit tests:
  - [x] Parse frontmatter from a vault markdown file returns correct metadata
  - [x] Extract `[[wiki-link]]` patterns from document body
  - [x] Classify symbol document by frontmatter fields
  - [x] Classify file document by frontmatter fields
  - [x] Extract named section returns correct content between headings
  - [x] Extract missing section returns empty string, not error
  - [x] Vault resolution finds .kodebase/ directory walking up from nested path
  - [x] Vault resolution returns error when no vault found
  - [x] Topic resolution finds topic directory within vault
  - [x] Symbol lookup matches partial name case-insensitively
- Integration tests:
  - [x] Read a vault directory written by WriteVault and verify all documents parse correctly
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- `make verify` passes
- Reader correctly round-trips documents written by the vault writer
