---
schema_version: 1
review_kind: techspec
round: 2
readiness: NEEDS_REWORK
reviewer_runtime: claude
reviewer_model: opus
generated_at: 2026-06-27T23:58:00Z
---

# Summary

Both round-1 blockers are genuinely closed: the spec now defines a concrete Concept
Path Model (flat bundle root, `slug(title)` filenames with `-N` collision suffix), a
full `ConceptResult` struct, and a deterministic `[[wikilink]]`→relative-markdown
mapping rule, resolving the PRD catalog-granularity open question. It is held back
from approval because three contracts on the `promote`/`check` critical path are
still undefined or self-contradictory: the four-field frontmatter remap has no source
mapping for `description`/`timestamp` (neither field exists in any `kb`-produced doc),
the conformance checker would hard-error on the `CLAUDE.md`/`AGENTS.md` files its own
OKF scaffold writes, and the filename-vs-link key derivation is asymmetric, silently
breaking links to already-promoted concepts.

# Blockers

## B-201 — Four-field frontmatter remap is undefined for `description` and `timestamp`

- Section: `_techspec.md` "Data Models → OKF frontmatter contract" (lines 131-135), "Data Flow → Promote", "Testing → Frontmatter remap"
- Issue: The contract names four producer fields (`type`, `title`, `description`,
  `timestamp`) but only defines two sources: `type` ← `--type` and (implicitly)
  `title` ← source `title`. There is no source mapping for `description` or
  `timestamp`, and verification against the code shows neither field is ever produced
  by `kb`: every ingest/source/wiki document carries `title, type, stage, domain,
  source_kind, scraped/created, tags` and nothing else (`internal/ingest/ingest.go:168-176`,
  `internal/vault/render_wiki.go:53-60`, `internal/vault/render.go:392-452`). So
  `description` has no source at all, and the spec does not say whether `timestamp`
  maps from `scraped`/`created`, from the promote-time clock, or something else.
  Compounding this, `timestamp` is specified as "ISO 8601" but the only available
  source timestamps use `frontmatter.DateLayout = "2006-01-02"` (date-only;
  `internal/frontmatter/frontmatter.go:18`), so even a `scraped→timestamp` mapping
  would not satisfy a full ISO-8601/RFC3339 expectation without a defined formatting
  rule.
- Rationale: This is an interface/contract gap on the primary MVP path — the
  scoped-write contract names "payloads, types, or schemas waved at instead of
  defined" as a blocker, and ADR-002 promises `promote` is a "deterministic,
  testable" mechanical transform. "Remaps wiki frontmatter to the OKF contract" is
  not implementable as written when two of the four target fields have no defined
  origin. It also makes the "Frontmatter remap" unit test and the end-to-end
  "concept written with four fields" assertion impossible to specify, and renders the
  local-standard "all four producer fields present" warning permanently true-fires
  (description always empty) with no design answer.
- Suggested fix: Add an explicit remap table — source frontmatter key → OKF key — and
  define fallbacks: `title` ← source `title` (fallback to humanized source basename
  when absent), `description` ← a named source field or an explicit "empty + warning"
  policy, `timestamp` ← a named source date field (`scraped`/`created`) normalized to
  RFC3339, with the promote-time clock as the documented fallback. Inject the clock as
  a seam so golden/integration tests assert a fixed value instead of freezing a
  wall-clock literal.

## B-202 — OKF conformance checker hard-errors on the scaffold's own `CLAUDE.md`/`AGENTS.md`

- Section: `_techspec.md` "Data Flow → Scaffold" (lines 43-45), "Testing → End-to-end promote"; `adr-005.md` "Decision" / "Implementation Notes"
- Issue: The OKF scaffold writes `topic.yaml`, `index.md`, `log.md`, "an OKF-flavored
  `CLAUDE.md`, and the `AGENTS.md` symlink." `CLAUDE.md` and `AGENTS.md` are `.md`
  files with no `type` frontmatter (`internal/topic/topic.go:27` `topicMarkerFile =
  "CLAUDE.md"`; `:607-613` `AGENTS.md` is a symlink to `CLAUDE.md`). The checker
  "recursively walks the bundle; for every non-reserved `.md` … require a non-empty
  `type` (§9.2)" and ADR-005 declares only `index.md`/`log.md` reserved. Therefore a
  freshly scaffolded bundle fails `kb okf check` with hard errors on `CLAUDE.md` and
  `AGENTS.md` — directly contradicting the end-to-end integration test ("scaffold an
  OKF topic … `kb okf check` passes") and the PRD success metric "the operator's
  catalog passes the check." The same gap likely affects the vendored official
  bundles (GA4/Stack Overflow/Bitcoin) if any carries a `README.md` or other
  non-concept `.md`, which would break the headline "conformance passes on all three
  official sample bundles" criterion.
- Rationale: This is a designed-in conformance contradiction (scaffold output vs.
  checker rule) and a missing-verification surface: the spec's own success criteria
  cannot pass as specified. The reserved/excluded-file set is a checker contract that
  is currently incomplete.
- Suggested fix: Define the full reserved/excluded set the walker skips —
  `index.md`, `log.md`, `CLAUDE.md`, `AGENTS.md` (and skip symlinks so `AGENTS.md`
  is not followed into `CLAUDE.md`), plus whatever the official bundles legitimately
  carry (`README.md`, license/attribution files) — or scope the `type` requirement to
  concept files only. State the rule explicitly and add a unit case asserting these
  files are not flagged, plus the integration assertion that a scaffolded-then-promoted
  bundle passes `kb okf check`.

## B-203 — Concept filename and inbound-link key derive from different inputs (silent broken links to promoted concepts)

- Section: `_techspec.md` "Concept Path Model" (lines 149-160)
- Issue: A concept's filename is `slug(title)` (keyed on the source doc's frontmatter
  `title`), but an inbound cross-concept link's target is computed as
  `slug(base(docPath))` (keyed on the wikilink's path basename). These coincide only
  when `base(docPath) == title`. For codebase-rendered wiki concepts that holds
  (`GetWikiConceptPath` embeds the title verbatim as the basename,
  `internal/vault/pathutils.go:219-222`), but for ingested/authored research docs the
  basename is a slugified source path/URL (`GetRawFileDocumentPath`/
  `GetRawSymbolDocumentPath`, `pathutils.go:192-200`) or an arbitrary Obsidian
  note name that need not equal the frontmatter `title`. When they differ, a body
  link to an *already-promoted* concept resolves to a phantom `slug(base(path)).md`
  while the real file is `slug(title).md` — and the spec records this in
  `UnresolvedLinks` as "no promoted counterpart exists yet," which is false. The same
  asymmetry makes `index.md` (links keyed on each concept's own `title`) and body
  cross-links (keyed on `base(path)`) point at two different addresses for one
  concept.
- Rationale: This defeats the determinism that round-1 B-002's fix was supposed to
  deliver and undermines the PRD success metric "links resolved." It is a correctness
  gap on the headline `promote` guarantee, not merely a tolerated-broken-link case.
- Suggested fix: Derive both the filename and the link target from a single canonical
  key. Simplest deletable option: key the concept filename off `slug(base(sourceDocPath))`
  (the same input links use) so a target's would-be path always equals its eventual
  filename; or, if title-named files are required, build the bundle's path→file
  resolution from the maintained `index.md`/a promote manifest and resolve links
  against it. State which, and add a unit case that promotes two linked concepts and
  asserts the inbound link resolves to the real file (not a phantom).

# Nits

## N-201 — Re-implements slugging instead of reusing `vault.SlugifySegment`

- Section: `_techspec.md` "Concept Path Model" (`slug(title)`)
- Issue: A canonical slug primitive already exists (`vault.SlugifySegment`, `internal/vault/pathutils.go:98-127`: lowercase, `[a-z0-9]`, single-dash, `"item"` empty fallback); the spec describes a fresh `slug()` ("lowercase, ASCII, kebab-case") without naming it, and `internal/okf` may legally import `vault` per ADR-003.
- Suggested fix: Specify reuse of `vault.SlugifySegment` for filename and link derivation, and note its `"item"` empty fallback and accent-dropping behavior (e.g. `Conteúdo`→`conte-do`) so the operator's Portuguese catalog names are accounted for.

## N-202 — `promote --to` does not validate the target is `mode: okf`

- Section: `_techspec.md` "Concept Path Model → `--to`" (line 160), "Core Interfaces" (`PromoteInput.TargetTopic`)
- Issue: `LinkFormatterFor(topic)` returns `WikiLinkFormatter` for any non-okf topic, so promoting into a wiki-mode `--to` target would silently emit `[[wikilinks]]` into an OKF concept; the spec assumes an OKF target but states no guard.
- Suggested fix: Require `promote` to error when `--to` resolves to a non-`okf` topic, and add a unit/integration case asserting that rejection.

## N-203 — Wikilink `#anchor` fragment handling unspecified in the mapping rule

- Section: `_techspec.md` "Concept Path Model → Wikilink target mapping"; "Testing → Wikilink→markdown body transform" (`[[a#anchor]]`)
- Issue: A `[[a#anchor]]` test case is listed, but the strip→`base`→`slug` rule never says whether the `#anchor` fragment is preserved as `[label](slug.md#anchor)` or dropped.
- Suggested fix: State that the fragment is carried through to the emitted relative link and add the asserted-output case.

## N-204 — `TopicMode = string` is a type alias, not a defined type

- Section: `_techspec.md` "Data Models" (lines 113-118); `adr-003.md` "typed constants"
- Issue: `type TopicMode = string` is an alias with zero type distinction, contradicting ADR-003's "typed constants" and giving the scaffold/normalization no compiler help for exhaustiveness.
- Suggested fix: Use `type TopicMode string` (defined type) or justify the alias explicitly; keep the `topicMetadataFile.Mode` field's type consistent with it.

## N-205 — Observability omits the new failure/skip events

- Section: `_techspec.md` "Monitoring and Observability"
- Issue: The logged events cover transform/broken-link counts but not "promote rejected: target not okf" (N-202) or "reserved file skipped" (B-202), which are the new actionable failure surfaces.
- Suggested fix: Add those events to the `promote`/`okf check` debug/info logging list.

# Evidence

- Read in full: `_techspec.md`, `_prd.md`, `adrs/adr-001.md`..`adr-006.md`,
  `qa/peer-review-findings-round1.md`, `qa/peer-review-incorporation-round1.md`,
  `CLAUDE.md`, `AGENTS.md`, `CONTRIBUTING.md`.
- Verified against source: `internal/vault/pathutils.go:98-127` (`SlugifySegment`),
  `:192-200` (`GetRawFileDocumentPath`/`GetRawSymbolDocumentPath`), `:219-222`
  (`GetWikiConceptPath` embeds title verbatim, no slug), `:245-252` (`ToTopicWikiLink`
  emits `[[slug/path|label]]`); `internal/frontmatter/frontmatter.go:18` (`DateLayout`
  = `2006-01-02`), `:71` (`Parse`), `:102-118` (`Generate`), `:390-410`
  (`buildMappingNode` alphabetical `sort.Strings`, confirmed line 396);
  `internal/ingest/ingest.go:160-205` (`buildFrontmatter` field set — no
  `description`/`timestamp`); `internal/vault/render_wiki.go:53-60` and
  `internal/vault/render.go:392-452` (wiki/source frontmatter field sets);
  `internal/topic/topic.go:27` (`topicMarkerFile = "CLAUDE.md"`), `:78-79`
  (scaffold writes `CLAUDE.md`/`log.md`), `:114`/`:186` (`New`/`newWithDate`
  signatures), `:607-628` (`AGENTS.md` symlink to `CLAUDE.md`).
- Confirmed round-1 closures in the current spec text: B-001 — "Concept Path Model"
  section (lines 144-165), `ConceptResult` struct (lines 89-95), `PromoteInput`
  (lines 83-87), catalog-granularity resolved (lines 162-165); B-002 — deterministic
  `[[wikilink]]`→relative mapping (lines 153-159). Incorporated nits N-001..N-006
  verified present (smoke test line 221-223; field-order note lines 137-140; empty
  vocabulary lines 124-128; `mode` normalization lines 120/246-248; `New`/`newWithDate`
  mode parameter lines 257-259; index.md `okf_version` preservation lines 49-53).
- Limitation: the OKF v0.1 `SPEC.md` and the vendored sample bundles are not present
  in this repo, so the §9.2 reserved-file exemptions (B-202) and the exact ISO-8601
  shape OKF expects for `timestamp` (B-201) are reasoned from the ADRs, not verified
  against the upstream spec/bundles.

# Deferred Or Follow-Up

- B-201/B-202/B-203 are independent contract gaps but share a theme: the `promote`
  output and the `okf check` input contracts are under-specified at the field/file
  level. Resolving them needs a remap table, a reserved-file set, and a single
  filename/link key — all spec-text edits, not architecture changes.
- Re-confirm at implementation time (against the vendored bundles) whether official
  OKF bundles carry non-concept `.md` files that the reserved/excluded set must cover
  (B-202) and the precise `timestamp` format the samples use (B-201).
