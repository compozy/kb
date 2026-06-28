# PRD: OKF Dual-Mode for `kb`

## Overview

`kb` today implements one knowledge lifecycle: the Karpathy **LLM-Wiki** — a
personal research *lab* where raw sources are ingested, compiled into a wiki,
queried, and linted, linked with Obsidian `[[wikilinks]]`.

The **Open Knowledge Format (OKF)** v0.1 (Google Cloud, 2026-06-12) standardizes
the same markdown + YAML-frontmatter substrate for a different purpose: a
portable, vendor-neutral **catalog** that other people's agents and tools consume.
It requires exactly one thing of every concept — a `type` field — uses plain
markdown links, and is meant to be *declared and shared*, not *explored*.

These are two lifecycles, not one derived from the other:

- **LLM-Wiki** = a lab: `ingest → compile → query → lint`; for the operator to
  research.
- **OKF** = a catalog/API: `declare → consume`; for others (humans and agents) to
  consume; portable.

This duality already exists by hand in the operator's workspace: a `research/`
area (LLM-wiki, where they research) and a `second-brain/` area (OKF, where they
declare "what is mine" for operations). The daily, highest-leverage move is
**distilling research into the operational catalog** — done manually today.

This feature makes both lifecycles first-class in `kb` via a per-topic
**`mode: wiki | okf`** (default `wiki`, so nothing regresses), and ships the
distill loop as a one-command **`promote`** plus an OKF **conformance check**. It
is valuable because it turns a hand-maintained workflow into a tool-enforced one,
and positions `kb` as an OKF-conformant producer in a brand-new, vendor-neutral
ecosystem.

## Goals

- Let an operator maintain LLM-Wiki research and OKF catalogs in the same tool,
  with the research → catalog distill loop tool-enforced.
- Ship OKF as a first-class authoring mode (not an export-only projection).
- Produce bundles that are conformant with OKF v0.1 and aligned with Google's real
  sample bundles and tooling.
- Preserve a consistent, drift-free local standard for concept types.
- Guarantee zero regression for existing wiki topics.
- Milestone: MVP delivers `mode` + OKF scaffold + `promote` + conformance check;
  Phase 2 adds publish/consume interop; Phase 3 adds codebase → OKF.

## User Stories

**Primary persona — the Knowledge Operator** (researches in wiki mode, declares
"what's mine" in OKF mode, distills daily):

- As an operator, I want to mark a topic as an OKF bundle so `kb` scaffolds and
  lints it as a portable catalog instead of a research lab.
- As an operator, I want to promote a finished wiki concept into my OKF bundle
  with one command so my research compounds into my operational catalog without
  manual reformatting or re-linking.
- As an operator, I want my source research to stay untouched when I promote, so
  the lab remains the immutable record and the catalog is the distilled view.
- As an operator, I want `kb` to validate my bundle against OKF conformance so my
  catalog stays portable and consumable by agents.
- As an operator, I want a canonical list of concept types so my bundle stays
  consistent and avoids type drift over time.
- As an operator, I want my existing wiki topics to keep working exactly as before
  when OKF mode lands.

**Secondary personas** (later phases):

- As a knowledge publisher, I want to export a wiki topic as a portable OKF bundle
  so others can consume it (Phase 2).
- As an AI agent / consuming tool, I want to read a `kb`-produced bundle as
  standard OKF so I can use it without a custom integration (enabled throughout).

## Core Features

Grouped by priority. MVP features first.

### 1. OKF mode for topics (MVP)

- A per-topic **`mode: wiki | okf`** setting, defaulting to `wiki`.
- `kb topic new --mode okf` scaffolds an **OKF bundle** (flat typed concepts,
  auto-maintained `index.md` and `log.md`) instead of the wiki pyramid
  (`raw/` + `wiki/` + `outputs/`).
- An OKF bundle is a `kb` topic; it reuses the topic, vault, lint, and index
  machinery. `mode` selects scaffold shape, link format, frontmatter contract, and
  lint ruleset.

### 2. `kb promote` — wiki → OKF distill (MVP)

- `kb promote <wiki-doc> --to <okf-topic> --type <ConceptType>` creates a new
  typed OKF concept in the target bundle from a compiled wiki document.
- **Mechanical and non-LLM**: converts `[[wikilinks]]` to relative markdown links,
  remaps wiki frontmatter to the OKF contract (`type`, `title`, `description`,
  `timestamp`), and appends a dated entry to the bundle's `log.md`.
- **Non-destructive**: the source wiki document stays in research; a new concept is
  born in the catalog.
- Any intelligent rewriting/condensing is left to the operator or an external
  agent — `kb` does the structural work only.

### 3. OKF conformance check (MVP)

- A conformance command for an OKF bundle that validates OKF v0.1:
  - every non-reserved `.md` has parseable frontmatter with a non-empty `type`;
  - `index.md` / `log.md` follow their required shapes when present.
- **Lenient by default** (per spec §9): tolerates broken cross-links, unknown
  `type` values, and missing optional fields (so externally produced bundles pass).
- **Local standards as warnings** on the operator's own bundles: all four producer
  fields present, and `type` within the configured vocabulary.
- A `--strict` option promotes warnings to errors for CI gates.

### 4. Concept-type governance (MVP)

- A canonical **type vocabulary** in `kb.toml` (e.g. `[okf].types`).
- `promote` and authoring take `type` explicitly; a value outside the vocabulary
  produces a lint warning. The vocabulary is extended by editing `kb.toml`.
- Prevents type drift (`Voice Profile` vs `voice-profile`) and keeps a consistent
  local standard, while staying compatible with OKF's "no global registry" model.

### 5. OKF-native emission (MVP, underpins 1–4)

- Relative, GitHub-safe markdown links (not `[[wikilinks]]`, not absolute paths).
- The four producer fields on every concept; auto-maintained `index.md`
  (type-grouped bullet lists) and `log.md` (ISO-date entries, newest first).

### Later phases (not MVP)

- **`kb export okf`** — publish an existing wiki topic as a portable bundle
  (Phase 2).
- **`kb ingest okf`** — consume an external OKF bundle into research (Phase 2).
- **Codebase → OKF** — emit the typed symbol/relation graph as a bundle
  (Phase 3, differentiator demo).

## User Experience

**Operator journey (MVP):**

1. **Create a catalog.** `kb topic new <name> --mode okf` scaffolds an OKF bundle.
   Existing topics are unaffected (they remain `mode: wiki`).
2. **Research as usual.** Ingest and compile in wiki-mode topics — unchanged.
3. **Distill.** When a wiki concept is ready, `kb promote <wiki-doc> --to <bundle>
   --type <ConceptType>`. The concept appears in the catalog with relative links,
   the four producer fields, and a `log.md` entry; the source stays put.
4. **Verify.** Run the conformance check on the bundle; fix any warnings (missing
   fields, off-vocabulary types). Use `--strict` in CI.
5. **Consume / share.** The bundle is plain markdown + frontmatter — readable in
   Obsidian/GitHub and parseable by any OKF-aware agent.

**UX considerations:**

- `mode` is invisible to existing users; OKF behavior is opt-in per topic.
- `promote` is one command with a predictable, mechanical result (no surprises
  from an LLM in the loop).
- Conformance output is actionable (what failed, where, how to fix) and
  distinguishes hard errors from local-standard warnings.
- Discoverability: OKF options surface under existing `kb topic` / `kb lint`-style
  help; `kb.toml` documents the type vocabulary.

## High-Level Technical Constraints

Product boundaries (not implementation):

- **Conform to OKF v0.1**, and remain aligned with Google's reference samples
  (GA4, Stack Overflow, Bitcoin) and producer behavior.
- **Two deliberate, documented deviations** from the written spec, matching real
  tooling: emit the four producer fields (not just `type`), and emit relative
  links (not the spec-recommended absolute links, which break GitHub rendering).
- **Stay non-LLM.** The core workflow, including `promote`, performs no model
  inference; `kb` remains "the non-LLM workflow."
- **Zero regression.** Default `mode: wiki`; existing topics scaffold, lint, and
  build identically; no migration required.
- **Portable output.** Bundles are usable without `kb`: plain files renderable in
  Obsidian/GitHub and parseable by any OKF consumer.
- OKF is a v0.1 draft; conformance rules must be pinned to v0.1 and isolated so the
  spec can evolve without disrupting the rest of the tool.

## Non-Goals (Out of Scope)

- **Export and ingest interop** (`kb export okf`, `kb ingest okf`) — Phase 2.
- **Codebase → OKF** emission — Phase 3.
- **LLM-assisted distillation** inside `promote` (no `--distill`) — possibly a
  later flag, not now.
- **A global/shared type registry** — only a local per-vault vocabulary.
- **Whole-vault or cross-topic OKF operations** — OKF is per-topic/bundle.
- **An HTML/graph visualizer** for bundles.
- **Migrating existing wiki topics** to OKF automatically.
- **Changing the LLM-Wiki lifecycle** — wiki mode behavior is unchanged.

## Phased Rollout Plan

### MVP (Phase 1) — the distill loop

- OKF mode (`mode: wiki | okf`) + `kb topic new --mode okf` scaffold.
- `kb promote` (mechanical, non-destructive, non-LLM).
- OKF conformance check (lenient + local warnings + `--strict`).
- Concept-type vocabulary in `kb.toml`.
- OKF-native emission (relative links, four producer fields, auto `index.md` /
  `log.md`).
- **Success criteria to proceed:** the conformance check passes on all three
  official sample bundles; the operator's catalog passes the check; a real wiki
  concept promotes end-to-end and the result conforms; existing wiki topics show
  zero behavior change.

### Phase 2 — publish & consume interop

- `kb export okf` (wiki topic → portable bundle).
- `kb ingest okf` (external bundle → research), accepting real-world bundles under
  lenient conformance.
- **Success criteria to proceed:** a wiki topic round-trips through export and back
  through ingest without loss of conceptual content; third-party sample bundles
  ingest cleanly.

### Phase 3 — codebase → OKF

- Emit the existing typed symbol/relation graph as an OKF bundle (concepts =
  files/symbols, links = relations).
- **Long-term success:** a real codebase produces a conformant, navigable OKF
  bundle usable as agent context.

## Success Metrics

- **Spec correctness:** the conformance check passes on all three official OKF
  sample bundles (GA4, Stack Overflow, Bitcoin).
- **Loop works end-to-end:** a real wiki concept promotes into a bundle and the
  result conforms, with links resolved and the four producer fields present.
- **Operator catalog conformant:** the existing `second-brain/` catalog passes the
  conformance check (with local standards) once adopted.
- **Zero regression:** existing wiki topics scaffold, lint, and build identically
  to before (verified by the existing suite).
- **No type drift:** every concept type in the operator's bundle is within the
  `kb.toml` vocabulary.
- **Adoption:** the operator stops distilling by hand — `promote` replaces the
  manual copy/reformat step.

## Risks and Mitigations

- **Adoption risk** — if `promote` is not clearly lower-friction than copy/paste,
  it won't be used. *Mitigation:* one command, mechanical, predictable output;
  non-destructive so it's safe to run freely.
- **Spec-churn risk** — OKF is a v0.1 draft and may change. *Mitigation:* pin to
  v0.1; isolate conformance rules; track the spec repo.
- **Scope risk** — the MVP already touches several surfaces. *Mitigation:* strict
  phasing; export/ingest/codebase are explicit non-goals for the MVP.
- **Standard-drift risk** — free-text types fragment the catalog. *Mitigation:*
  `kb.toml` vocabulary + lint warning + `--strict` in CI.
- **Identity-drift risk** — pressure to add LLM "magic" to `promote`.
  *Mitigation:* keep the core non-LLM by decision (ADR-002); defer any
  `--distill` flag.
- **Interop-expectation risk** — bundles must be consumable outside `kb`.
  *Mitigation:* validate against official samples; emit what real tooling expects
  (four fields, relative links).

## Architecture Decision Records

- [ADR-001: OKF as a first-class authoring mode, not an export-only projection](adrs/adr-001.md)
  — OKF is a first-class per-topic `mode`, an OKF bundle is a `kb` topic, and the
  MVP optimizes the research → operations distill loop with phased delivery.
- [ADR-002: promote is a mechanical non-LLM transform; type governance and a spec-compatible conformance posture](adrs/adr-002.md)
  — `promote` is mechanical, non-destructive, and non-LLM; `type` is explicit and
  governed by a local `kb.toml` vocabulary; conformance emits four fields + relative
  links and validates lenient with local warnings and `--strict`.

## Open Questions

- **Conformance command surface** — `kb okf check <bundle>` vs `kb lint --okf`
  (naming/placement). Defer to TechSpec.
- **Catalog granularity** — within the operator's `second-brain/`, do sub-areas
  (`identidade/`, `estudos/`, `conteudo/`, `vendas/`) map to one OKF topic with
  nested concept directories, or to multiple OKF topics? Needs confirmation; both
  are valid OKF.
- **index.md / log.md ownership** — assumed auto-maintained by `kb`
  (generated/updated on scaffold + promote). Confirm whether operators ever
  hand-edit them.
- **Off-vocabulary type handling** — assumed "warn + allow + extend via
  `kb.toml`". Confirm whether `promote` should suggest the closest existing type.
- **Broken-link reporting on promote** — promoting a concept whose wikilinks point
  to not-yet-promoted concepts yields tolerated broken links (per §9). Confirm
  whether the operator wants a warning listing them.
