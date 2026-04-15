# PRD: Java Codebase Ingest Adapter

## Overview

The product will add Java support to `kb ingest codebase` so teams can generate usable knowledge artifacts from Java systems with the same workflow already available for TypeScript/JavaScript, Go, and Rust.

It is designed for platform architects, product engineers, and modernization teams who need to understand structure, dependencies, and change impact in Java codebases without manual mapping. The value is faster repository comprehension, better decision confidence, and reduced time to start modernization work.

## Goals

- Enable Java repositories to be ingested and indexed through the existing `kb ingest codebase` command flow.
- Reduce time-to-understanding for architecture/dependency analysis in Java systems.
- Improve confidence in change planning by exposing meaningful relationship data in generated knowledge artifacts.
- Support first-cycle adoption across common Java repository shapes, including standard Java, Spring-oriented projects, and multi-module layouts.
- Preserve user-perceived ingest performance for large repositories (no significant regression versus current non-Java runs).

## User Stories

- As a platform architect, I want to ingest a Java codebase and see structural artifacts quickly so that I can map architecture and dependencies early.
- As a product engineer, I want to inspect relationships between Java symbols and files so that I can estimate change impact with more confidence.
- As a modernization lead, I want consistent Java ingest outputs across multi-module repositories so that migration planning starts from a reliable baseline.
- As an engineering manager, I want Java ingest to be predictable in runtime so that teams can include it in routine analysis workflows.

## Core Features

- **Java language detection in codebase ingest**
  - The system recognizes Java repositories as supported ingest input.
  - Java files are included in scan and ingest summaries using the same reporting conventions as other supported languages.

- **Java structural extraction for knowledge artifacts**
  - The system produces knowledge artifacts from Java source structures that are useful for architecture and dependency exploration.
  - Extracted structures are visible through existing generated outputs and inspect flows.

- **Relationship visibility for decision support**
  - The system provides relationship data sufficient for practical dependency and impact exploration in the first release.
  - Relationship output prioritizes high-value common cases used during architecture review and change planning.

- **Performance-safe ingest experience**
  - Java support is delivered without materially degrading ingest experience in large repositories.
  - Users retain predictable progress and completion behavior in existing codebase ingest workflows.

## User Experience

- A user runs the same command they already use for codebase ingest and points it to a Java repository.
- The command reports Java as a detected language in the run summary.
- The generated topic content includes Java-backed artifacts in the same navigational structure users already understand.
- The user can inspect generated outputs to review architecture shape, dependencies, and likely change impact.
- The experience remains familiar: no new conceptual workflow is required to benefit from Java support.

## High-Level Technical Constraints

- Must integrate into existing codebase ingest command behavior and output conventions.
- Must preserve current user-facing reliability and performance expectations for large repositories.
- Must maintain compatibility with current topic artifact structure and inspection UX.
- Must avoid introducing workflow fragmentation (Java support should feel native, not parallel).

## Non-Goals (Out of Scope)

- Full semantic precision for every advanced Java edge case in the first release.
- Expanding ingest to non-Java new languages in this initiative.
- Redesigning inspect UX or topic artifact taxonomy as part of Java support.
- Solving organization-specific governance or migration strategy decisions beyond ingest output.

## Phased Rollout Plan

### MVP (Phase 1)

- Deliver Java ingest support with broad structural coverage and practical relation quality for high-frequency use cases.
- Support standard Java repositories, common Spring-oriented project shapes, and typical multi-module layouts.
- Preserve ingest runtime expectations on representative medium/large repositories.

Success criteria to proceed to Phase 2:
- Teams can complete architecture and dependency first-pass analysis from generated artifacts without manual reconstruction.
- Users report improved confidence in change planning for Java code.
- No significant ingest performance regression is observed in pilot usage.

### Phase 2

- Improve relationship fidelity for more complex Java usage patterns that appear in pilot feedback.
- Strengthen consistency of outputs across diverse multi-module enterprise repository structures.
- Expand fit for additional real-world repository conventions.

Success criteria to proceed to Phase 3:
- Reduction in unresolved or ambiguous analysis outcomes reported by pilot teams.
- Repeat usage by all three target personas in routine workflows.

### Phase 3

- Mature Java ingest quality to support broader organizational rollout and governance use cases.
- Optimize adoption enablement with clearer guidance for large-scale Java portfolio ingestion.

Long-term success criteria:
- Java ingest becomes a default discovery step in architecture review, change assessment, and modernization planning workflows.

## Success Metrics

- **Time-to-understanding:** measurable reduction in time spent to produce an initial architecture/dependency map for Java repositories.
- **Change confidence:** increase in self-reported confidence for planning Java code changes based on ingest outputs.
- **Modernization acceleration:** shorter lead time from repository handoff to first modernization plan draft.
- **Performance quality:** stable ingest completion behavior on large Java repositories, without meaningful user-perceived slowdown.

## MVP Governance Decisions

- **No significant regression threshold:** Java ingest is considered within budget when median total runtime increase is <=20% versus agreed baseline runs on the canonical pilot set, measured with identical command flags over 3 repeated runs.
- **Canonical pilot repository set:** MVP validation uses three repository profiles: (1) single-module Java library, (2) Spring-style service repository, and (3) multi-module enterprise-style repository. These profiles are the mandatory acceptance corpus for rollout decisions.
- **Minimum confidence target for rollout:** Proceed beyond MVP only when at least 80% of pilot participants report confidence >=4/5 for Java change-impact and dependency analysis workflows, with no unresolved critical workflow blockers.

## Risks and Mitigations

- **Adoption risk:** users may perceive early Java output as insufficient for complex cases.
  - **Mitigation:** communicate MVP scope clearly and prioritize high-value improvements in Phase 2 using pilot feedback.

- **Expectation mismatch risk:** stakeholders may assume full semantic coverage from day one.
  - **Mitigation:** position release as balanced MVP focused on practical value and progressive fidelity.

- **Rollout risk across varied repos:** diverse enterprise repo structures may reveal edge gaps early.
  - **Mitigation:** include representative multi-module repositories in pilot and use phased expansion criteria.

- **Timeline risk:** broad persona demands can pull scope beyond MVP.
  - **Mitigation:** enforce phased boundaries and gate expansion on agreed success metrics.

## Architecture Decision Records

- [ADR-001: Adopt a balanced MVP strategy for Java codebase ingest](adrs/adr-001.md) — Choose broad early coverage with practical relation quality and performance-safe rollout over precision-first or coverage-first extremes.
- [ADR-005: Define MVP governance acceptance gates and pilot corpus](adrs/adr-005.md) — Formalizes performance threshold, pilot validation corpus, and confidence gate for rollout readiness.
- [ADR-006: Close Java ingest MVP rollout using available pilot evidence](adrs/adr-006.md) — Records MVP rollout closure decision and deferred governance evidence handling into Phase 2.

## Open Questions

- None at MVP governance level. Additional open items should be tracked in Phase 2 planning artifacts.
