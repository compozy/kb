---
status: pending
title: Implement topic scaffolding
type: backend
complexity: medium
dependencies:
  - task_01
  - task_02
---

# Task 04: Implement topic scaffolding

## Overview

Create the `internal/topic/` package that manages knowledge base topic lifecycle: scaffolding new topics (directory tree, template files, CLAUDE.md, log.md), listing existing topics, and retrieving topic info. This replaces the shell script `new-topic.sh` from the karpathy-kb skill with native Go.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST create the full topic directory tree as specified in TechSpec "Data Models" section (raw/{articles,bookmarks,codebase,github,youtube}, wiki/{concepts,index}, outputs/{queries,briefings,diagrams,reports}, bases/)
- MUST install template files (Dashboard, Concept Index, Source Index, log.md, CLAUDE.md) from `.agents/skills/karpathy-kb/assets/` with domain/title/slug substituted
- MUST create AGENTS.md as symlink to CLAUDE.md
- MUST implement List operation that discovers topics by scanning vault root for directories with the expected skeleton
- MUST implement Info operation that reads topic metadata (article count, source count, last log entry)
- MUST use the frontmatter package (task_01) for template frontmatter generation
- MUST append a `## [YYYY-MM-DD] scaffold | <topic-slug>` entry to log.md on creation
</requirements>

## Subtasks

- [ ] 4.1 Create `internal/topic/` package with New, List, and Info functions
- [ ] 4.2 Implement directory tree creation with all required subdirectories
- [ ] 4.3 Implement template installation with variable substitution (slug, title, domain, date)
- [ ] 4.4 Implement topic listing by scanning vault root
- [ ] 4.5 Implement topic info by reading metadata from the topic skeleton
- [ ] 4.6 Write unit tests using `t.TempDir()` for filesystem operations

## Implementation Details

Create `internal/topic/topic.go` and `internal/topic/topic_test.go`. The template files are read from `.agents/skills/karpathy-kb/assets/` at build time or embedded via `//go:embed`.

The existing vault writer (`internal/vault/writer.go`) already has `ensureTopicSkeleton()` — study its approach but implement independently in the topic package since the vault writer's version is codebase-specific.

### Relevant Files

- `.agents/skills/karpathy-kb/assets/` — 6 template files (dashboard, concept-index, source-index, log, topic-claude, wiki-article)
- `.agents/skills/karpathy-kb/scripts/new-topic.sh` — current shell implementation to replicate behavior from
- `internal/vault/writer.go` — has `ensureTopicSkeleton()` which creates a subset of the directory tree
- `internal/frontmatter/` (task_01) — used for template frontmatter generation

### Dependent Files

- `internal/ingest/` (task_12) — needs topic to exist before writing to `raw/`
- `internal/cli/` (task_15) — `topic new`, `topic list`, `topic info` commands wire to this package

### Related ADRs

- [ADR-001: Topic-Centric CLI Command Taxonomy](../adrs/adr-001.md) — topic is the primary organizing concept

## Deliverables

- `internal/topic/topic.go` — New, List, Info functions
- `internal/topic/topic_test.go` — comprehensive tests
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests

- Unit tests:
  - [ ] New creates all expected subdirectories (raw/articles, raw/bookmarks, raw/codebase, raw/github, raw/youtube, wiki/concepts, wiki/index, outputs/queries, outputs/briefings, outputs/diagrams, outputs/reports, bases/)
  - [ ] New installs Dashboard, Concept Index, Source Index templates with substituted variables
  - [ ] New creates CLAUDE.md with correct topic metadata
  - [ ] New creates AGENTS.md as symlink pointing to CLAUDE.md
  - [ ] New appends scaffold entry to log.md
  - [ ] New returns error if topic directory already exists
  - [ ] List returns empty slice for vault with no topics
  - [ ] List returns topic slugs for vault with multiple topics
  - [ ] Info returns correct article/source counts for a populated topic
- Test coverage target: >=80%
- All tests must pass

## Success Criteria

- All tests passing
- Test coverage >=80%
- Scaffolded topic directory matches the structure defined in TechSpec
- Template files contain correct substituted values
- `make lint` reports zero findings
