# Peer Review — Round 2 Incorporation

- **Decision:** incorporate **all 3 blockers + all 5 nits**.
- **Date:** 2026-06-27
- **Round-1 closure:** B-001/B-002 confirmed closed by the reviewer; no action needed.

## Incorporated

| Finding | Resolution | File(s) |
| --- | --- | --- |
| B-201 — `description`/`timestamp` had no source | Added a **source→OKF remap table**: `type`←`--type`; `title`←source title (fallback humanized basename); `description`←`--description` flag → first body sentence → empty+warning; `timestamp`←promote-time clock (RFC3339), injected as a `Clock` seam for tests; `tags`←source when present. | `_techspec.md` (Data Models), `adr-002.md` |
| B-202 — checker rejects scaffold's own `CLAUDE.md`/`AGENTS.md` | Defined the full **reserved/excluded set** the walker skips: `index.md`, `log.md`, `CLAUDE.md`, `AGENTS.md`, symlinks, `README`/license; `type` requirement scoped to concept files. | `_techspec.md` (Data Flow→Check), `adr-005.md` |
| B-203 — filename vs link key asymmetric | Switched to a **single canonical key** `vault.SlugifySegment(base(path))` for both filename and link target, so a link to an already-promoted concept always resolves. Documented the filename trade-off (path-key now; title+manifest is a Phase-2 option). | `_techspec.md` (Concept Path Model, Data Flow) |
| N-201 — reuse slug primitive | Named `vault.SlugifySegment` as the key (with `"item"` fallback + accent-drop `Conteúdo`→`conte-do`). | `_techspec.md` (Concept Path Model) |
| N-202 — `--to` not validated as okf | `promote` errors when `--to` is not a `mode: okf` topic; added test + log event. | `_techspec.md` (Concept Path Model, Testing, Monitoring), `adr-006.md` |
| N-203 — `#anchor` handling | Fragment carried through as `[label](<key>.md#anchor)`. | `_techspec.md` (Concept Path Model) |
| N-204 — `TopicMode` alias | Changed to a defined type `type TopicMode string`. | `_techspec.md` (Data Models) |
| N-205 — missing failure/skip events | Added `promote rejected: target not okf`, `description fallback used`, `reserved file skipped` to observability. | `_techspec.md` (Monitoring) |

## Deferred

- **None.** All round-2 findings incorporated. The reviewer's environment limitation
  (OKF `SPEC.md` / sample bundles not in repo) is re-confirmed at implementation time
  against the vendored fixtures — already an ADR-005 task.

## Files changed

- `_techspec.md`, `adrs/adr-002.md`, `adrs/adr-005.md`, `adrs/adr-006.md`
