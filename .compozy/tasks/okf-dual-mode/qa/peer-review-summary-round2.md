# Spec Peer Review — Round 2 Summary

- **Spec:** `.compozy/tasks/okf-dual-mode/_techspec.md` (post round-1 incorporation)
- **Reviewer runtime:** `compozy exec` · claude / opus / xhigh (cross-LLM)
- **Verdict:** **NEEDS_REWORK** — 3 blockers, 5 nits
- **Round-1 closure:** B-001 and B-002 **confirmed closed** against current spec text.
- **Validator:** passed (exit 0); status before==after (scoped-write honored).

## Blockers (new, field/file-level contract gaps)

- **B-201 — four-field remap undefined for `description`/`timestamp`.** No `kb`-produced
  doc carries `description` or `timestamp`; the only date source is `DateLayout`
  (`2006-01-02`, date-only), so the "ISO 8601" `timestamp` and the `description` field
  have no defined origin. *Fix:* explicit source→OKF remap table with fallbacks; inject
  the clock as a seam for tests.
- **B-202 — checker hard-errors on the scaffold's own `CLAUDE.md`/`AGENTS.md`.** Only
  `index.md`/`log.md` are reserved, so a freshly scaffolded bundle fails `kb okf check`
  — contradicting the e2e test and the PRD success metric; likely also breaks the
  official-bundle pass if they carry `README.md`. *Fix:* define the full
  reserved/excluded set (`index.md`, `log.md`, `CLAUDE.md`, `AGENTS.md`, skip symlinks,
  `README`/license), scope `type` to concept files, add assertions.
- **B-203 — filename vs inbound-link key are asymmetric.** Filename = `slug(title)`,
  link target = `slug(base(docPath))`; they diverge for ingested/authored docs, so a
  link to an already-promoted concept resolves to a phantom path and is falsely logged
  as unresolved. *Fix:* derive both from one canonical key (recommend keying the
  filename off `slug(base(sourceDocPath))` — the same input links use — or resolve via
  a maintained index/manifest).

## Nits

- **N-201** — reuse `vault.SlugifySegment` (not a fresh `slug()`); note `"item"` empty
  fallback and accent-drop (`Conteúdo`→`conte-do`) for PT-BR names.
- **N-202** — `promote --to` must reject a non-`okf` target (else `[[wikilinks]]` leak
  into an OKF concept).
- **N-203** — specify `#anchor` fragment handling in the wikilink mapping (`[label](slug.md#anchor)`).
- **N-204** — `type TopicMode = string` is an alias; use a defined type per ADR-003.
- **N-205** — add observability events for "promote rejected: target not okf" and
  "reserved file skipped".

## Sections / ADRs likely affected by incorporation

- `_techspec.md`: Data Models (frontmatter remap table, `TopicMode` defined type),
  Concept Path Model (single canonical key, anchor handling), Data Flow → Scaffold +
  Promote, Testing, Monitoring/Observability.
- `adr-005.md` (reserved/excluded file set), `adr-003.md` (defined `TopicMode`),
  `adr-002.md` (description/timestamp sourcing if it touches the product contract).

## Reviewer limitation noted

OKF `SPEC.md` / sample bundles not in repo: the §9.2 reserved-file exemptions (B-202)
and the exact `timestamp` ISO-8601 shape (B-201) are reasoned from ADRs, to be
re-confirmed at implementation time against the vendored bundles.

## Operational artifacts (round 2)

- Findings: `qa/peer-review-findings-round2.md`
- Prompt: `qa/peer-review-prompt-round2.md`
- Events: `qa/peer-review-events-round2.jsonl` · Stderr: `qa/peer-review-result-round2.err`
- Status snapshots: `qa/peer-review-status-{before,after}-round2.txt`
