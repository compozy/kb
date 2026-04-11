---
status: completed
title: Implement ingest orchestrator
type: backend
complexity: high
dependencies:
  - task_01
  - task_02
  - task_04
  - task_05
---

# Task 12: Implement ingest orchestrator

## Overview

Create the `internal/ingest/` package that orchestrates the full ingestion pipeline: select converter → convert source → generate frontmatter → write to vault → append log entry. This is the central integration hub connecting converters, frontmatter, topic management, and vault writing.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST accept a source (file path or converted content) and target topic, and produce a fully ingested vault document
- MUST select the appropriate converter from the registry based on file extension or explicit source type
- MUST generate correct frontmatter per the schema in TechSpec "Data Models" section (type: source, stage: raw, domain, source_kind, tags, scraped date)
- MUST write output to `<topic>/raw/<source-type>/<slug>.md` using a sanitized slug from the title
- MUST append `## [YYYY-MM-DD] ingest | <slug>.md (<source_kind>)` to `<topic>/log.md`
- MUST return `IngestResult` with topic, source type, file path, and title
- MUST validate that the target topic exists before writing
- SHOULD generate a unique slug when a file with the same name already exists (append `-N` suffix)
</requirements>

## Subtasks

- [x] 12.1 Create `internal/ingest/` package with Ingest function
- [x] 12.2 Implement converter selection and invocation via the registry
- [x] 12.3 Implement frontmatter generation with correct schema fields per source kind
- [x] 12.4 Implement vault file writing to the correct `raw/<source-type>/` subdirectory
- [x] 12.5 Implement log.md appending with the standard format
- [x] 12.6 Implement slug generation with deduplication
- [x] 12.7 Write unit tests using `t.TempDir()` with a scaffolded topic

## Implementation Details

Create `internal/ingest/ingest.go` and `internal/ingest/ingest_test.go`. The orchestrator does NOT import specific converters or external clients — it receives converted content or delegates to the registry. The CLI layer (task_16) handles fetching from Firecrawl/YouTube before calling the ingest orchestrator.

Reference TechSpec "System Architecture" data flow diagram for the pipeline sequence.

### Relevant Files

- `internal/frontmatter/` (task_01) — Generate function for frontmatter creation
- `internal/topic/` (task_04) — topic validation and path resolution
- `internal/convert/registry.go` (task_05) — converter selection for file-based ingest
- `internal/models/kb_models.go` (task_02) — IngestResult, SourceKind types

### Dependent Files

- `internal/cli/` (task_16) — ingest commands call this orchestrator
- `internal/generate/` (task_14) — codebase ingest adapter calls this for output writing

### Related ADRs

- [ADR-001: Topic-Centric CLI Command Taxonomy](../adrs/adr-001.md) — ingest as the primary ingestion verb
- [ADR-003: Native Go Document Conversion with Converter Registry](../adrs/adr-003.md) — converter selection pattern

## Deliverables

- `internal/ingest/ingest.go` — ingest orchestrator
- `internal/ingest/ingest_test.go` — tests
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests

- Unit tests:
  - [x] Ingest file selects correct converter and writes to raw/articles/<slug>.md
  - [x] Ingest generates correct frontmatter with type: source, stage: raw, domain, source_kind
  - [x] Ingest appends log entry to log.md in correct format
  - [x] Ingest returns IngestResult with all fields populated
  - [x] Ingest returns error when target topic does not exist
  - [x] Ingest generates unique slug when file already exists (appends -2, -3, etc.)
  - [x] Ingest writes to correct subdirectory per source kind (articles/, youtube/, github/, bookmarks/)
  - [x] Ingest with pre-converted content (no converter needed) writes directly
- Integration tests:
  - [x] Scaffold topic → ingest file → verify vault output + log entry + frontmatter
- Test coverage target: >=80%
- All tests must pass

## Success Criteria

- All tests passing
- Test coverage >=80%
- Ingested files appear in correct vault location with valid frontmatter
- log.md reflects every ingest operation
- `make lint` reports zero findings
