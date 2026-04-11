---
status: completed
title: Implement CLI ingest commands
type: backend
complexity: high
dependencies:
  - task_10
  - task_11
  - task_12
  - task_15
---

# Task 16: Implement CLI ingest commands

## Overview

Implement the `kb ingest` parent command and its subcommands: `url`, `file`, `youtube`, `codebase`, and `bookmarks`. Each subcommand is a thin CLI adapter that parses flags, invokes the appropriate external client or converter, and delegates to the ingest orchestrator for vault writing.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST implement `kb ingest url <url> [--topic <slug>]` that uses the Firecrawl client to scrape and ingest
- MUST implement `kb ingest file <path> [--topic <slug>]` that uses the converter registry to convert and ingest
- MUST implement `kb ingest youtube <url> [--topic <slug>] [--stt]` that extracts transcripts and ingests
- MUST implement `kb ingest codebase <path> [--topic <slug>]` that runs the adapted generate pipeline
- MUST implement `kb ingest bookmarks <path> [--topic <slug>]` that ingests bookmark cluster markdown files
- All subcommands MUST accept `--topic` flag to specify target topic (required)
- All subcommands MUST validate that the target topic exists
- `ingest codebase` MUST accept all existing generate flags (--include, --exclude, --semantic, --progress)
- Each subcommand MUST print the IngestResult on success (file path, title, source type)
</requirements>

## Subtasks

- [x] 16.1 Create `internal/cli/ingest.go` parent command
- [x] 16.2 Implement `ingest url` subcommand wired to Firecrawl client → ingest orchestrator
- [x] 16.3 Implement `ingest file` subcommand wired to converter registry → ingest orchestrator
- [x] 16.4 Implement `ingest youtube` subcommand wired to YouTube extractor → ingest orchestrator
- [x] 16.5 Implement `ingest codebase` subcommand wired to adapted generate pipeline
- [x] 16.6 Implement `ingest bookmarks` subcommand wired to ingest orchestrator
- [x] 16.7 Write unit tests for flag parsing and command routing

## Implementation Details

Create `internal/cli/ingest.go` (parent), `internal/cli/ingest_url.go`, `internal/cli/ingest_file.go`, `internal/cli/ingest_youtube.go`, `internal/cli/ingest_codebase.go`, `internal/cli/ingest_bookmarks.go`.

Each subcommand follows the pattern: parse CLI flags → create client/load config → fetch/convert → call ingest orchestrator → print result. Keep commands thin — business logic lives in the respective packages.

`ingest codebase` largely replaces the existing `generate` command. Consider keeping `generate` as a hidden alias for backwards compatibility during transition.

### Relevant Files

- `internal/cli/root.go` (task_15) — register ingest parent command
- `internal/cli/generate.go` — existing generate command to adapt/replace for ingest codebase
- `internal/firecrawl/` (task_10) — Firecrawl client for URL ingest
- `internal/youtube/` (task_11) — YouTube extractor for youtube ingest
- `internal/ingest/` (task_12) — ingest orchestrator called by all subcommands
- `internal/generate/` (task_14) — adapted pipeline for codebase ingest

### Dependent Files

- `internal/cli/root.go` — registers ingest as a subcommand

### Related ADRs

- [ADR-001: Topic-Centric CLI Command Taxonomy](../adrs/adr-001.md) — ingest subcommand structure
- [ADR-004: Firecrawl REST API for URL Scraping](../adrs/adr-004.md) — URL ingest via REST API

## Deliverables

- `internal/cli/ingest.go` — parent command
- `internal/cli/ingest_url.go` — URL ingest command
- `internal/cli/ingest_file.go` — file ingest command
- `internal/cli/ingest_youtube.go` — YouTube ingest command
- `internal/cli/ingest_codebase.go` — codebase ingest command
- `internal/cli/ingest_bookmarks.go` — bookmarks ingest command
- Tests for all commands
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests

- Unit tests:
  - [x] `ingest` parent shows help listing all subcommands
  - [x] `ingest url` requires URL positional arg and --topic flag
  - [x] `ingest url` returns error when --topic is missing
  - [x] `ingest file` requires file path positional arg and --topic flag
  - [x] `ingest file` returns error for non-existent file
  - [x] `ingest youtube` requires YouTube URL positional arg and --topic flag
  - [x] `ingest youtube` accepts --stt flag
  - [x] `ingest codebase` requires path positional arg and --topic flag
  - [x] `ingest codebase` accepts --include, --exclude, --semantic flags
  - [x] `ingest bookmarks` requires path positional arg and --topic flag
- Test coverage target: >=80%
- All tests must pass

## Success Criteria

- All tests passing
- Test coverage >=80%
- All five ingest subcommands are discoverable via `kb ingest --help`
- Each subcommand correctly routes to its implementation package
- `make lint` reports zero findings
