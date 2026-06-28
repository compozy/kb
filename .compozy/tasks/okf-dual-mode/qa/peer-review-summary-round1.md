# Spec Peer Review — Round 1 Summary

- **Spec:** `.compozy/tasks/okf-dual-mode/_techspec.md`
- **Reviewer runtime:** `compozy exec` · claude / opus / xhigh (cross-LLM, independent)
- **Verdict:** **NEEDS_REWORK** — 2 blockers, 6 nits
- **Validator:** `validate-findings.sh` passed (exit 0); status before==after (scoped-write honored)

## Blockers (must resolve before approval)

- **B-001 — `promote` output contract undefined.** No bundle write-location /
  filename derivation, and `ConceptResult` is referenced but never defined. Tied to
  the PRD's deferred "catalog granularity" open question. *Rationale:* interface gap
  on the primary MVP path; the end-to-end promote test can't assert a concrete
  location and `index.md` type-grouping has no defined tree to walk.
- **B-002 — wikilink→relative target resolution unspecified.** `[[slug/path]]` lives
  in the *source* topic's namespace; the promoted concept lives in a *different*
  topic (the OKF bundle). The mapping rule from wiki-namespace target → bundle-
  relative `targetPath` is undefined, so `OKFLinkFormatter.Link` can't be called
  deterministically. *Rationale:* the mechanical core ADR-002 calls "deterministic,
  testable"; couples to B-001.

> Root cause: both blockers reduce to **the bundle's internal concept-path model**.
> Reviewer recommends resolving the "catalog granularity" question once (flat
> bundle root for the MVP = simplest, deletable) — and both blockers follow.

## Nits (non-blocking)

- **N-001** — dormant OKF render branch ships untested in MVP (re-raises ADR-004 Alt 1: body-transform-only is simpler).
- **N-002** — `frontmatter.Generate` sorts keys alphabetically → emitted order is `description,tags,timestamp,title,type`, not `type`-first; tests must not freeze producer order.
- **N-003** — `index.md` regeneration must preserve the root `okf_version`; state hand-edit ownership.
- **N-004** — new wiki scaffolds will gain a visible `mode: wiki` line; separate "old files unchanged" from "new scaffold output" in the zero-regression story.
- **N-005** — default `[okf].types` vocabulary contents unspecified (enumerate or ship empty).
- **N-006** — `topic.New(vaultPath, slug, title, domain)` signature change for `--mode` not called out in the Build Order.

## Sections / ADRs likely affected by incorporation

- `_techspec.md`: Data Flow → Promote; Core Interfaces (`PromoteInput`/`ConceptResult`); Data Models; Testing; Build Order steps 1, 5, 7, 9.
- `adr-004.md` (link target resolution + migration-effort estimate ~12 vs ~34 lines).
- `adr-005.md` (index.md regeneration / okf_version).
- PRD open questions: "catalog granularity" and "index.md/log.md ownership" should be resolved into the TechSpec.

## Reviewer limitation noted

Could not read OKF `SPEC.md` / sample bundles (not present in repo): the §9/§6/§7
claims and "four producer fields" / field-order behavior were taken from the ADRs,
not independently verified against upstream.

## Operational artifacts (round 1)

- Findings: `qa/peer-review-findings-round1.md`
- Prompt: `qa/peer-review-prompt-round1.md`
- Events: `qa/peer-review-events-round1.jsonl`
- Stderr: `qa/peer-review-result-round1.err`
- Status snapshots: `qa/peer-review-status-{before,after}-round1.txt`
