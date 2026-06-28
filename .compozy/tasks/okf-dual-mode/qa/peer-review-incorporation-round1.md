# Peer Review — Round 1 Incorporation

- **Decision:** incorporate **both blockers + all six nits** (N-001 handled as a
  smoke test, preserving the ADR-004 "migrate now" decision rather than deferring).
- **Date:** 2026-06-27

## Incorporated

| Finding | Resolution | File(s) changed |
| --- | --- | --- |
| B-001 — promote output contract undefined | Added a **Concept Path Model** (flat bundle root for MVP; filename = `slug(title)` with `-N` collision suffix) and defined the full `ConceptResult` struct; resolved the PRD "catalog granularity" open question into the spec. | `_techspec.md` (Core Interfaces, Data Flow, new Concept Path Model section) |
| B-002 — wikilink→relative target resolution unspecified | Specified the deterministic mapping `[[srcSlug/docPath]]` → strip slug → `base(docPath)` → `slug(...)` → bundle-relative path; flat root ⇒ `[label](<slug>.md)`; unresolved targets emit a would-be link and land in `ConceptResult.UnresolvedLinks`. | `_techspec.md` (Concept Path Model, Data Flow) |
| N-001 — dormant OKF render branch (dead code) | Added an **OKF render smoke test** (render-site-in-`okf`-mode) to keep the wired branch live; recorded as a risk mitigation. | `_techspec.md` (Testing), `adr-004.md` (Risks) |
| N-002 — alphabetical key order | Documented that emitted field order is `description, tags, timestamp, title, type`; tests must not assert `type`-first ordering. | `_techspec.md` (frontmatter contract) |
| N-003 — index.md regeneration / okf_version | Stated regeneration preserves the root `okf_version`; defined `index.md` (kb-owned, regenerated) vs `log.md` (append-only) ownership. | `_techspec.md` (Data Flow), `adr-005.md` (Implementation Notes) |
| N-004 — scaffold `mode: wiki` line | Separated "old files unchanged" from "new scaffold output gains an explicit `mode` line"; clarified `omitempty` scope. | `adr-003.md` (Implementation Notes) |
| N-005 — default `[okf].types` vocabulary | Specified the vocabulary ships **empty** (check is a no-op until the operator opts in). | `_techspec.md` (Data Models) |
| N-006 — `topic.New` signature change | Called out the `New`/`newWithDate` `mode` parameter (or `NewWithMode`) in Build Order step 5. | `_techspec.md` (Build Order) |
| Follow-up — ~12 vs ~34 migration estimate | Reconciled to "~7 in render.go + ~34 lines in render_wiki.go (~43 working estimate)". | `adr-004.md` (Implementation Notes) |

## Deferred

- **None.** All round-1 findings were incorporated. The reviewer's note that it
  could not independently verify the OKF spec/sample bundles (not in repo) is a
  review-environment limitation, not a spec gap — the vendored fixtures (ADR-005)
  will provide that verification at implementation time.

## Files changed

- `_techspec.md`
- `adrs/adr-003.md`
- `adrs/adr-004.md`
- `adrs/adr-005.md`
