---
status: completed
title: Implement lint engine
type: backend
complexity: high
dependencies:
  - task_01
  - task_02
---

# Task 13: Implement lint engine

## Overview

Create the `internal/lint/` package that performs structural health checks on a knowledge base vault. The lint engine walks the vault directory tree, parses frontmatter, extracts wikilinks, and detects issues: dead wikilinks, orphan articles, missing source references, stale content, and frontmatter format violations. This replaces the Python `lint-wiki.py` script with native Go.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST detect dead wikilinks: `[[Target]]` links where the target file does not exist in the vault
- MUST detect orphan articles: wiki articles in `wiki/concepts/` with zero incoming wikilinks
- MUST detect missing sources: frontmatter `sources:` entries referencing absent `raw/` files
- MUST detect stale content: articles whose `updated` date is older than their source files' `scraped` date
- MUST detect format violations: files missing required frontmatter fields per their type/stage schema
- MUST return `[]LintIssue` sorted by severity (errors first, then warnings) and file path
- MUST support output as structured data compatible with `internal/output` formatter (table/json/tsv)
- MUST optionally save report to `outputs/reports/<date>-lint.md`
</requirements>

## Subtasks

- [x] 13.1 Create `internal/lint/` package with Lint function that accepts a topic path
- [x] 13.2 Implement vault walker that parses all markdown files and builds a wikilink graph
- [x] 13.3 Implement dead wikilink detection (link targets not found in vault)
- [x] 13.4 Implement orphan article detection (wiki concepts with zero incoming links)
- [x] 13.5 Implement missing source detection (frontmatter source refs not in raw/)
- [x] 13.6 Implement stale content and format violation detection
- [x] 13.7 Write unit tests with fixture vaults containing known issues

## Implementation Details

Create `internal/lint/lint.go` and `internal/lint/lint_test.go`. The lint engine builds an in-memory graph of `filePath → {incoming wikilinks, outgoing wikilinks, frontmatter}` by walking the topic directory.

Wikilink extraction regex: `\[\[([^\]|]+)(?:\|[^\]]+)?\]\]` — captures the target, ignoring display text after `|`.

Reference TechSpec "Testing Approach" lint tests section for the test strategy. Reference `.agents/skills/karpathy-kb/references/lint-procedure.md` for the full check list.

### Relevant Files

- `internal/frontmatter/` (task_01) — Parse function for reading frontmatter from vault files
- `internal/models/kb_models.go` (task_02) — LintIssue, LintIssueKind types
- `.agents/skills/karpathy-kb/references/lint-procedure.md` — full lint check list to implement
- `.agents/skills/karpathy-kb/scripts/lint-wiki.py` — Python reference implementation

### Dependent Files

- `internal/cli/` (task_17) — `lint` command wires to this package
- `internal/output/` — lint results formatted via existing formatter

### Related ADRs

- [ADR-005: Native Go Vault Lint Engine](../adrs/adr-005.md) — native Go chosen over keeping Python script

## Deliverables

- `internal/lint/lint.go` — lint engine with all check types
- `internal/lint/lint_test.go` — tests with fixture vaults
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests

- Unit tests:
  - [x] Detects dead wikilink when `[[Missing Article]]` target file does not exist
  - [x] Does NOT flag valid wikilinks to existing files
  - [x] Detects orphan article in wiki/concepts/ with zero incoming links
  - [x] Does NOT flag articles that have at least one incoming link
  - [x] Detects missing source when frontmatter `sources:` references absent raw/ file
  - [x] Detects stale article when `updated` < source `scraped` date
  - [x] Detects format violation: wiki article missing required `title` frontmatter
  - [x] Detects format violation: source file missing `source_kind` frontmatter
  - [x] Returns issues sorted by severity then file path
  - [x] Returns empty slice for a healthy vault with no issues
  - [x] Handles wikilinks with display text: `[[Target|Display]]` extracts "Target"
- Integration tests:
  - [x] Build fixture vault with mix of healthy and broken files → lint → verify all issues detected
- Test coverage target: >=80%
- All tests must pass

## Success Criteria

- All tests passing
- Test coverage >=80%
- All five check types (dead links, orphans, missing sources, stale, format) working
- Output compatible with `internal/output` formatter
- `make lint` reports zero findings
