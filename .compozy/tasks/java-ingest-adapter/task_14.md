---
status: completed
title: Stabilize JSON contract for automation consumers
type: backend
complexity: medium
dependencies:
  - task_11
---

# Task 14: Stabilize JSON contract for automation consumers

## Overview
Stabilize the Java ingest JSON contract used by automation and platform workflows to reduce integration fragility at scale. This task formalizes expected fields and compatibility boundaries for CLI outputs and summary payloads.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- The codebase ingest JSON output contract MUST define stable required fields for automation.
- Contract changes SHOULD remain backward-compatible or explicitly versioned.
- CLI tests MUST assert required contract keys and value semantics for Java ingest outputs.
- Documentation MUST clearly state contract guarantees and non-guaranteed fields.
</requirements>

## Subtasks
- [x] 14.1 Define required JSON contract surface for Java ingest summary/result payloads.
- [x] 14.2 Add/update CLI tests that enforce required output contract keys.
- [x] 14.3 Add compatibility guidance for future contract evolution.
- [x] 14.4 Ensure contract behavior holds for dry-run and full ingest modes.
- [x] 14.5 Publish contract notes in initiative documentation.

## Implementation Details
Leverage existing `codebaseIngestResult` and `GenerationSummary` payloads and lock minimum contract expectations for external consumers. Reference TechSpec sections "Impact Analysis" and "Monitoring and Observability".

### Relevant Files
- `internal/cli/ingest_codebase.go` — JSON result payload shape.
- `internal/models/models.go` — `GenerationSummary` contract surface.
- `internal/cli/ingest_test.go` — command output and help behavior tests.
- `internal/cli/workflow_integration_test.go` — end-to-end JSON payload assertions.

### Dependent Files
- `.compozy/tasks/java-ingest-adapter/task_15.md` — adoption playbook should reference the stabilized contract.
- `.compozy/tasks/java-ingest-adapter/_rollout-mvp-signoff.md` — evidence sections may consume contract fields.

### Related ADRs
- [ADR-004: Require unit, integration, benchmark, and CLI E2E validation for Java ingest](../adrs/adr-004.md) — contract must be protected by tests.
- [ADR-006: Close Java ingest MVP rollout using available pilot evidence](../adrs/adr-006.md) — Phase 3 hardening includes automation stability.

## Deliverables
- Defined and documented stable JSON contract for Java ingest automation.
- Updated CLI test coverage enforcing contract keys and modes.
- Backward-compatibility guidance for future contract evolution.
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for JSON contract stability **(REQUIRED)**

## Tests
- Unit tests:
  - [x] Java ingest JSON includes required identity fields (`topic`, `sourceType`, summary core fields).
  - [x] Dry-run and full-run payloads maintain documented required key set.
- Integration tests:
  - [x] E2E CLI ingest JSON remains parseable with expected key/value schema.
  - [x] Existing automation-facing output assertions remain green after updates.
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- Automation consumers can rely on a documented stable JSON contract
- Contract behavior remains consistent across dry-run and full ingest workflows
