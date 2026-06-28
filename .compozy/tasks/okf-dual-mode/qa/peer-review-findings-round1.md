---
schema_version: 1
review_kind: techspec
round: 1
readiness: NEEDS_REWORK
reviewer_runtime: claude
reviewer_model: opus
generated_at: 2026-06-27T22:20:41Z
---

# Summary

The design is sound and well-bounded — the `mode` field is a clean additive,
backward-compatible change; the `internal/okf` package boundary is one-way and
cycle-free; and the conformance posture is correctly reasoned. It is held back
from approval because the MVP's central verb, `promote`, has two undefined
contracts on its critical path: where a promoted concept file lands in the bundle,
and how a `[[wikilink]]` target (a *source-topic* path) is resolved to a relative
path *inside the OKF bundle*. Both must be specified before implementation.

# Blockers

## B-001 — `promote` output contract is undefined (write location, filename, `ConceptResult`)

- Section: `_techspec.md` "Data Flow → Promote", "Core Interfaces" (`PromoteInput`/`Promoter`/`ConceptResult`)
- Issue: `PromoteInput` carries `SourceDocPath`, `TargetTopic`, and `Type` but no
  output location. The data flow says promote "writes the new concept via
  `frontmatter.Generate` + `vault`" without defining the target path or filename
  inside the bundle: is it the source basename, a slug of `title`, or a
  `type`-grouped subdirectory? `Promoter.Promote` returns `ConceptResult`, which is
  referenced but never defined. The PRD's own open question — "Catalog granularity …
  sub-areas map to one OKF topic with nested concept directories, or multiple OKF
  topics?" — is explicitly deferred "to TechSpec" and remains unanswered here, yet
  it directly determines the concept path model.
- Rationale: This is an interface/contract gap on the primary MVP path. The
  scoped-write contract for this review calls out "payloads, types, or schemas waved
  at instead of defined" as a blocker, and `CLAUDE.md` requires business logic
  (concept placement) to be concretely specified before Cobra wiring. Without a
  defined path model, the end-to-end promote integration test ("concept written
  with four fields + relative links … `index.md` updated") cannot assert a concrete
  location, and `index.md` type-grouping (which must enumerate concept files) has no
  defined tree to walk.
- Suggested fix: Define the bundle's internal concept-path model explicitly: the
  concept filename derivation (recommend deterministic slug of `title`, collision
  rule stated), whether concepts sit at bundle root or under a `type`/area
  subdirectory, and the full `ConceptResult` struct (at minimum: written path,
  transformed-link count, broken/unresolved-link list). Resolve the PRD "catalog
  granularity" open question in the spec (flat root for MVP is the simplest,
  deletable choice) rather than leaving it to implementation.

## B-002 — wikilink→relative-markdown target resolution is unspecified

- Section: `_techspec.md` "Data Flow → Promote", "Implementation Design → Wikilink→markdown body transform"; `adr-004.md` "Decision"
- Issue: `OKFLinkFormatter.Link(fromDir, targetPath, label)` computes a relative
  path from the new concept's directory to `targetPath`. During promote, the source
  body contains `[[slug/path|label]]`, where `slug/path` is a target in the *source
  wiki topic's* namespace (verified: `ToTopicWikiLink` emits `[[topicSlug/docPath]]`
  at `internal/vault/pathutils.go:245`). The promoted concept lives in a *different*
  topic (the OKF bundle) and the target concept usually is not promoted yet. The
  spec never defines the rule that maps a wiki-namespace target to a bundle-relative
  `targetPath` — so `OKFLinkFormatter.Link` cannot be called deterministically, and
  "converts `[[wikilinks]]` to relative markdown links" is not implementable as
  written.
- Rationale: This is a second interface gap on the critical path, and it is the
  mechanical core that ADR-002 promises is "deterministic, testable." The PRD
  success metric "links resolved and the four producer fields present" and the unit
  test "links to not-yet-promoted concepts (tolerated/flagged)" both presuppose a
  defined target-derivation rule that does not exist. It is coupled to B-001: the
  link target is whatever path model B-001 settles.
- Suggested fix: Specify the deterministic mapping from a `[[slug/path|label]]`
  target to a bundle-relative path (e.g., strip the source slug, map to the same
  concept-filename derivation chosen in B-001, emit `[label](relpath)` from the
  concept's `fromDir`), and state explicitly that targets with no promoted
  counterpart still emit a relative link to the *would-be* path and are reported as
  tolerated/unresolved (per §9). Add a unit case asserting the exact emitted
  relative path for a cross-concept link.

# Nits

## N-001 — Dormant OKF render branch ships untested in the MVP

- Section: `_techspec.md` "Executive Summary"/"Primary trade-off"; `adr-004.md` Decision §3
- Issue: The ~40+ render call sites (verified: 7 in `render.go`, 34 lines in
  `render_wiki.go`) gain a mode-aware seam whose OKF branch is "wired but dormant
  until Phase 3" and never exercised at render sites in the MVP, i.e. shipped dead
  code; ADR-004 Alternative 1 (body-transform only now) is the simpler, deletable
  option the project's bias favors.
- Suggested fix: Either defer the 43-site migration per ADR-004 Alt 1, or add one
  render-site-in-OKF-mode smoke test so the shipped branch (and its `fromDir`
  plumbing) is not dead and untested.

## N-002 — `frontmatter.Generate` sorts keys; producer field order is alphabetical

- Section: `_techspec.md` "OKF frontmatter contract"; `internal/frontmatter/frontmatter.go:396`
- Issue: `buildMappingNode` does `sort.Strings(keys)`, so the four fields emit as
  `description, timestamp, title, type` (plus `tags`), not Google's `type`-first
  producer order — harmless for conformance (YAML maps are unordered) but worth
  stating so no test freezes a producer order the helper can't honor.
- Suggested fix: Note in the spec that emitted field order is alphabetical and that
  conformance/tests must not assert producer ordering (or add a dedicated ordered
  encoder if byte-matching the samples is actually required).

## N-003 — `index.md` regeneration may drop root `okf_version` and clobber hand edits

- Section: `_techspec.md` "Data Flow → Promote" ("regenerates `index.md`"); `adr-005.md`
- Issue: The scaffold writes a root `index.md` with `okf_version: "0.1"`
  frontmatter, but promote "regenerates `index.md`" as type-grouped bullet lists;
  the spec does not say the regeneration preserves the root `okf_version`, and the
  PRD open question "index.md / log.md ownership … whether operators ever hand-edit
  them" is unresolved.
- Suggested fix: Specify that regeneration preserves the bundle-root `okf_version`
  and state the hand-edit ownership contract (fully `kb`-owned vs merge-preserving).

## N-004 — Scaffold `mode: wiki` vs `omitempty` "byte-stable" claim

- Section: `_techspec.md` Build Order step 1; `adr-003.md` Implementation Notes
- Issue: ADR-003 argues `omitempty` keeps files "byte-stable when mode is wiki," but
  that only covers *not rewriting* existing files; a newly scaffolded wiki topic
  will gain a visible `mode: wiki` line unless wiki is written as empty — a
  scaffold-output change the zero-regression story should acknowledge.
- Suggested fix: State whether wiki scaffolds emit `mode: wiki` (recommended for
  clarity) and update the scaffold golden expectation, separating "old files
  unchanged" from "new scaffold output."

## N-005 — `[okf].types` default vocabulary contents unspecified

- Section: `_techspec.md` "Data Models" (`OKFConfig`); Build Order step 2
- Issue: The spec adds `[okf].types` "defaults" and a `config.example.toml` entry but
  never lists the default concept-type vocabulary.
- Suggested fix: Enumerate the default `types` list (or state it ships empty) so the
  example file and out-of-vocabulary warning tests have a concrete baseline.

## N-006 — `topic.New` signature change for `--mode` not called out

- Section: `_techspec.md` Build Order steps 5/9; `internal/topic/topic.go:114`
- Issue: `New(vaultPath, slug, title, domain string)` has no `mode` parameter; the
  spec implies but never states the signature/constructor change needed to thread
  mode from CLI into scaffold.
- Suggested fix: Name the `New`/`newWithDate` signature change (or a `NewWithMode`
  variant) in the build order so the CLI→topic plumbing is explicit.

# Evidence

- Read in full: `_techspec.md`, `_prd.md`, `adrs/adr-001.md`..`adr-006.md`,
  `CLAUDE.md`, `AGENTS.md`, `CONTRIBUTING.md`.
- Verified against source: `internal/vault/pathutils.go:239-252` (`ToTopicWikiLink`
  emits `[[slug/path|label]]`); `internal/vault/render.go:186-187` and call sites at
  220/339/473/555/566/632/643; `internal/vault/render_wiki.go` (34 matching lines);
  `internal/frontmatter/frontmatter.go:102-118` (`Generate`) and `:390-410`
  (`buildMappingNode` alphabetical `sort.Strings`); `internal/models/models.go:271`
  (`TopicMetadata`, has `Slug`, no `Mode`); `internal/models/kb_models.go:115-123`
  (`TopicInfo` exists, no `Mode`); `internal/topic/topic.go:104-110`
  (`topicMetadataFile`, no `Mode`), `:114` (`New` signature), `:607-628`
  (`ensureAgentsSymlink` — AGENTS.md is a symlink to CLAUDE.md), `:770-781`
  (`readTopicYAMLMetadata` returns empty struct when file absent);
  `internal/lint/lint.go:25/703/904` (`wikilinkPattern`, `schemaForPath`).
- Limitation: I could not read the OKF v0.1 `SPEC.md` or the vendored sample bundles
  (not present in this repo); §9/§6/§7 claims and the "four producer fields" /
  field-order behavior are taken from the ADRs, not independently verified against
  the upstream spec.

# Deferred Or Follow-Up

- B-001/B-002 share a root cause (the bundle's internal concept-path model); resolve
  the PRD "catalog granularity" open question once and both follow.
- Confirm whether matching Google's sample bundles requires byte-identical
  producer-field ordering (N-002) or only structural conformance.
- ADR-004 estimates "~12 direct `ToTopicWikiLink` calls in `render_wiki.go`" but the
  file has ~34 matching lines; reconcile the migration-effort estimate (non-blocking).
