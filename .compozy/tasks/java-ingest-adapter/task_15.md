---
status: completed
title: Create Java portfolio adoption playbook
type: docs
complexity: medium
dependencies:
  - task_12
  - task_13
  - task_14
---

# Task 15: Create Java portfolio adoption playbook

## Overview
Create a practical playbook for operating Java ingest across large repository portfolios with governance, observability, and automation guidance. This task turns hardening outputs into repeatable adoption workflows for platform and modernization teams.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- The playbook MUST define recommended ingest flow for large Java portfolios.
- The playbook MUST reference performance gate policy, telemetry interpretation, and JSON automation contract.
- The playbook SHOULD include troubleshooting guidance for high fallback/unresolved scenarios.
- Documentation MUST be aligned with current CLI behavior and verified commands.
</requirements>

## Subtasks
- [x] 15.1 Draft portfolio-scale ingest workflow covering discovery, dry-run, full ingest, and post-checks.
- [x] 15.2 Document governance checkpoints and evidence collection templates.
- [x] 15.3 Document telemetry and diagnostics interpretation guidance for operators.
- [x] 15.4 Document automation contract usage patterns for external tooling.
- [x] 15.5 Validate all documented commands and references against current CLI behavior.

## Implementation Details
Author rollout/adoption guidance based on finalized observability telemetry, benchmark governance process, and stabilized JSON contract. Keep guidance actionable for recurring enterprise operations and Phase 3 long-term goals.

### Relevant Files
- `.compozy/tasks/java-ingest-adapter/_rollout-mvp-signoff.md` — rollout evidence baseline and governance framing.
- `.compozy/tasks/java-ingest-adapter/_prd.md` — Phase 3 goals and governance context.
- `.compozy/tasks/java-ingest-adapter/_techspec.md` — telemetry and benchmark policy references.
- `internal/cli/ingest_codebase.go` — command surface to document accurately.
- `internal/cli/lint.go` — post-ingest quality validation commands.

### Dependent Files
- `.compozy/tasks/java-ingest-adapter/adrs/adr-005.md` — governance criteria source.
- `.compozy/tasks/java-ingest-adapter/adrs/adr-006.md` — rollout closure context.
- Future Phase 3/4 planning artifacts — will depend on the playbook as operating baseline.

### Related ADRs
- [ADR-005: Define MVP governance acceptance gates and pilot corpus](../adrs/adr-005.md) — governance criteria to operationalize.
- [ADR-006: Close Java ingest MVP rollout using available pilot evidence](../adrs/adr-006.md) — transition to broad adoption guidance.

## Deliverables
- Java portfolio adoption playbook with governance and operations guidance.
- Command-validated examples for ingest, lint, and evidence collection.
- Troubleshooting matrix for fallback-heavy enterprise scenarios.
- Unit tests with 80%+ coverage **(REQUIRED)**
- Integration tests for documentation command correctness **(REQUIRED)**

## Tests
- Unit tests:
  - [x] Documentation references required governance thresholds and contract fields consistently.
  - [x] Playbook includes explicit handling for high `JAVA_RESOLUTION_FALLBACK` volumes.
- Integration tests:
  - [x] All documented commands execute successfully in a controlled test workflow.
  - [x] Playbook command outputs match expected CLI fields and semantics.
- Test coverage target: >=80%
- All tests must pass

## Success Criteria
- All tests passing
- Test coverage >=80%
- Teams can execute Java ingest governance workflow from a single operational playbook
- Adoption guidance is consistent with real CLI behavior and telemetry outputs
