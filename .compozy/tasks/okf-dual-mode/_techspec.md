# TechSpec: OKF Dual-Mode for `kb`

> Implements `_prd.md`. Business context (WHAT/WHY/WHO) lives there; this document
> covers HOW. Phasing follows the PRD: **MVP** = `mode` + OKF scaffold +
> `kb promote` + `kb okf check`; Phase 2 = export/ingest; Phase 3 = codebase→OKF.

## Executive Summary

Add OKF as a first-class topic mode to `kb` by (1) introducing a `mode: wiki|okf`
field on the topic model, (2) abstracting link emission behind a mode-aware
`LinkFormatter`, (3) housing all OKF behavior in a new cohesive `internal/okf`
package (frontmatter contract, wikilink→markdown transform, `index.md`/`log.md`
generation, conformance checker, and the `promote` orchestration), and (4) exposing
a top-level `kb promote` plus a `kb okf` command group. OKF emission stays
**mechanical and non-LLM**, matching `kb`'s identity; conformance is **lenient per
OKF §9** with local-standard **warnings** and a `--strict` gate, emitting the four
producer fields and relative, GitHub-safe links to match Google's real tooling.

**Primary trade-off:** per the user's decision (ADR-004), the MVP introduces the
full `LinkFormatter` abstraction and migrates **all ~43 link call sites** in the
codebase→wiki render pipeline now, even though that pipeline produces wiki-mode
topics and its OKF branch stays **dormant until Phase 3**. We accept a large
mechanical refactor with no MVP-exercised OKF render path in exchange for a single,
finished mode-aware link seam — guarded against regression by byte-identical golden
tests of wiki output. The MVP's *active* OKF link consumer is `promote` (body
transform) and `index.md` generation.

## System Architecture

### Component Overview

| Component | Home | Responsibility |
| --- | --- | --- |
| `mode` field | `internal/models`, `internal/topic` | Mark a topic as `wiki` (default) or `okf`; normalize empty→`wiki`. |
| `LinkFormatter` | `internal/vault` | Mode-aware link rendering; `WikiLinkFormatter` + `OKFLinkFormatter`. |
| OKF core | `internal/okf` (new) | Frontmatter contract, link transform, `index.md`/`log.md` gen, conformance checker, `promote`. |
| Type vocabulary | `internal/config` | `[okf].types` canonical list for local-standard validation. |
| OKF scaffold | `internal/topic` | `--mode okf` builds a minimal bundle (no wiki pyramid). |
| CLI | `internal/cli` | `kb promote`, `kb okf check`, `kb topic new --mode`. |

### Data Flow

- **Scaffold:** `kb topic new --mode okf` → `internal/topic` writes `topic.yaml`
  (`mode: okf`), a root `index.md` (`okf_version: "0.1"`), `log.md`, an OKF-flavored
  `CLAUDE.md`, and the `AGENTS.md` symlink. No `raw/`/`wiki/`/`outputs/`.
- **Promote:** `kb promote <wiki-doc> --to <topic> --type <T>` → `internal/okf`
  reads the source via `vault` + `frontmatter.Parse`, remaps frontmatter to the OKF
  contract, transforms `[[wikilinks]]`→relative markdown links, writes the new
  concept to `<bundle-root>/<key(base(SourceDocPath))>.md` (see **Concept Path
  Model**) via `frontmatter.Generate` + `vault`, appends to `log.md`, and regenerates
  `index.md` (preserving the root `okf_version`). Source doc is untouched.
- **Check:** `kb okf check <topic>` → `internal/okf` walks the bundle, **skipping the
  reserved/excluded set** (`index.md`, `log.md`, `CLAUDE.md`, `AGENTS.md`, symlinks,
  `README`/license files), validates §9 leniently over the remaining concept files,
  emits `[]models.LintIssue` rendered by `internal/output`.

Dependency direction is one-way: `internal/okf` → {`frontmatter`, `vault`,
`config`, `models`}; nothing imports `okf` back (ADR-003).

## Implementation Design

### Core Interfaces

Mode-aware link rendering, defined in `internal/vault` (ADR-004):

```go
// LinkFormatter renders a link from the document being written (fromDir,
// bundle-relative, forward-slash) to a target document, with an optional label.
type LinkFormatter interface {
    Link(fromDir, targetPath, label string) string
}

// LinkFormatterFor selects the formatter from the topic's mode.
func LinkFormatterFor(topic models.TopicMetadata) LinkFormatter {
    if topic.Mode == models.TopicModeOKF {
        return OKFLinkFormatter{}
    }
    return WikiLinkFormatter{Slug: topic.Slug}
}
```

OKF promotion and conformance, in `internal/okf`:

```go
type PromoteInput struct {
    SourceDocPath string           // wiki doc to read (vault-relative)
    TargetTopic   models.TopicInfo // an OKF-mode topic (the bundle)
    Type          string           // OKF concept type (validated vs config.OKF.Types)
}

// ConceptResult reports what promote wrote; UnresolvedLinks tolerated (§9).
type ConceptResult struct {
    WrittenPath     string   // bundle-relative path of the new concept
    Type            string
    LinksRewritten  int      // [[wikilink]] occurrences converted
    UnresolvedLinks []string // targets with no promoted counterpart yet
}

type Promoter interface {
    Promote(ctx context.Context, in PromoteInput) (ConceptResult, error)
}

type Checker interface {
    // Check walks an OKF bundle; strict promotes warnings to errors.
    Check(ctx context.Context, bundlePath string, strict bool) ([]models.LintIssue, error)
}
```

### Data Models

Additions only; existing structs are unchanged except for one field each.

```go
// internal/models
type TopicMode string // defined type (not an alias) per ADR-003
const (
    TopicModeWiki TopicMode = "wiki"
    TopicModeOKF  TopicMode = "okf"
)
// TopicMetadata and TopicInfo each gain: Mode TopicMode `json:"mode"`

// internal/topic topicMetadataFile gains:
//   Mode string `yaml:"mode,omitempty"`   // empty normalized to "wiki" on read

// internal/config
type OKFConfig struct {
    Types []string `toml:"types"` // local vocabulary; ships EMPTY (no false warnings)
}
// Config gains: OKF OKFConfig `toml:"okf"`
// Empty Types ⇒ the type-vocabulary check is a no-op until the operator opts in.
```

**OKF frontmatter contract — source→OKF remap (B-201).** `promote` and authoring
write the four producer fields via this deterministic map; wiki stage markers
(`stage`, wiki `type`) are dropped:

| OKF field | Source | Fallback |
| --- | --- | --- |
| `type` | `--type` (validated vs `config.OKF.Types`) | none — required; error if absent |
| `title` | source `title` | humanized `base(SourceDocPath)` |
| `description` | `--description` flag | first non-empty body sentence (markdown-stripped) → else empty + warning |
| `timestamp` | promote-time clock, RFC3339 / UTC | — (promotion *is* the "last meaningful change"; source `scraped`/`created` are date-only `2006-01-02`, so they are not the timestamp source) |
| `tags` | source `tags` when present | omitted |

The clock is injected as a seam (a `Clock` function) so golden/integration tests
assert a fixed `timestamp` instead of freezing a wall-clock literal.

**Field order:** `frontmatter.Generate` sorts keys alphabetically
(`internal/frontmatter/frontmatter.go:396`), so concepts emit `description, tags,
timestamp, title, type`. This is conformant (YAML maps are unordered); tests MUST
NOT assert a `type`-first producer order. If byte-matching the official samples is
ever required, that needs a dedicated ordered encoder (out of scope).

### Concept Path Model (resolves B-001 / B-002 / B-203)

For the MVP an OKF bundle is **one topic**, and promoted concepts land **flat at the
bundle root** (`type`-grouped or nested subdirectories are deferred). A **single
canonical key** drives both the filename and every inbound link, so they can never
diverge (B-203):

- **Canonical key:** `vault.SlugifySegment(base(path))` — the existing slug primitive
  (`internal/vault/pathutils.go:98`; lowercase `[a-z0-9-]`, `"item"` empty fallback,
  drops accents so `Conteúdo`→`conte-do`). `internal/okf` reuses it (N-201) rather
  than re-implementing a slug.
- **Filename:** a concept from `SourceDocPath` is written to
  `<topic-root>/<key(base(SourceDocPath))>.md`. Collisions get a numeric suffix
  (`-2`, `-3`, …).
- **Wikilink target mapping:** a body link `[[srcSlug/docPath|label]]` strips the
  source slug and maps to `<topic-root>/<key(base(docPath))>.md` — the **same** key
  the target's own filename uses, so a link to an already-promoted concept always
  resolves to its real file. `OKFLinkFormatter.Link(fromDir, target, label)` renders
  it; with a flat root `fromDir == targetDir`, the link is `[label](<key>.md)`. A
  `#anchor` fragment is carried through: `[label](<key>.md#anchor)` (N-203).
- **Unresolved targets:** when no promoted counterpart exists yet, the link is still
  emitted to the would-be path (tolerated per §9) and recorded in
  `ConceptResult.UnresolvedLinks` (logged; a warning under `--strict`).
- **`--to`** names the target OKF topic and **must resolve to a `mode: okf` topic** —
  `promote` errors otherwise (N-202), so wiki-formatter output never leaks into a
  bundle. `--to <topic>/<subdir>` writes under that subdir (`fromDir` updates
  accordingly); the default is the bundle root.

> Filename trade-off: keying off the source path (not `title`) guarantees
> symmetry/correctness at the cost of less pretty names for ingested docs (a slug of
> a source path/URL). Title-named files with an index/manifest-based link resolver is
> a deliberate Phase-2 option, not MVP.

> Resolves the PRD open question "catalog granularity": for the MVP a bundle is one
> OKF topic with concepts at its root. Splitting `second-brain/` sub-areas
> (`identidade`, `estudos`, …) into separate OKF topics, or introducing nested
> concept directories, is a deliberate Phase-2 follow-up, not MVP behavior.

### Command Surface (CLI)

| Command | Args / Flags | Behavior |
| --- | --- | --- |
| `kb topic new <slug> <title> <domain>` | `--mode wiki\|okf` (default `wiki`) | Wiki scaffold (unchanged) or minimal OKF scaffold. |
| `kb promote <wiki-doc>` | `--to <topic>`, `--type <T>`, `--description <text>` | Mechanical, non-destructive wiki→OKF concept; `--to` must be a `mode: okf` topic. |
| `kb okf check <topic>` | `--strict`, `--format table\|json\|tsv` | OKF v0.1 conformance + local-standard warnings. |

Phase 2 grows the group: `kb okf export`, `kb okf ingest` (out of MVP scope).

## Integration Points

No external services. The only "integration" is the **OKF v0.1 file format**
(`GoogleCloudPlatform/knowledge-catalog` → `okf/SPEC.md`). Two deliberate,
documented deviations from the written spec, matching Google's reference tooling
(ADR-002): emit the four producer fields (spec mandates only `type`) and emit
relative links (spec recommends absolute `/path.md`, which breaks GitHub
rendering). The official sample bundles are vendored as test fixtures (ADR-005),
not a runtime dependency.

## Impact Analysis

| Component | Impact Type | Description and Risk | Required Action |
| --- | --- | --- | --- |
| `internal/okf` | new | Cohesive OKF package; isolated, low blast radius. | Build per ADR-003/004/005. |
| `internal/vault` (pathutils, render, render_wiki) | modified | `LinkFormatter` + route ~43 call sites; **regression risk** in wiki output. | Golden byte-identical wiki tests; keep suite green. |
| `internal/models` | modified | Add `Mode` to `TopicMetadata`/`TopicInfo` + mode consts; low risk. | Add fields/consts. |
| `internal/topic` | modified | `Mode` on `topicMetadataFile`; mode-conditional scaffold + OKF templates; medium risk (scaffold branch). | Normalize empty→wiki; add OKF assets. |
| `internal/config` | modified | Add `[okf].types`; low risk (additive). | Add struct + example + defaults. |
| `internal/cli` | modified | `kb promote`, `kb okf` group, `topic new --mode`; low risk. | Wire commands per existing patterns. |
| `internal/lint` | unchanged | OKF uses a separate checker (ADR-005); wiki lint untouched. | None. |
| existing wiki topics | unchanged | Default mode wiki ⇒ zero behavior change. | Verify via existing tests. |

## Testing Approach

### Unit Tests

- **`OKFLinkFormatter.Link`** — table-driven relative-path cases: sibling, child,
  parent (`../`), `fromDir == targetDir`, bundle root; never a leading `/`.
- **`WikiLinkFormatter.Link`** — asserts byte-identical output to today's
  `ToTopicWikiLink` (with and without label).
- **Wikilink→markdown body transform** — `[[a]]`, `[[a|label]]`, `[[a#anchor]]`,
  links to not-yet-promoted concepts (tolerated/flagged).
- **Frontmatter remap (B-201)** — source→OKF per the remap table; `description`
  fallback chain (flag → first body sentence → empty + warning); `timestamp` from the
  injected clock, asserted as a fixed RFC3339 value; stage markers dropped.
- **Reserved/excluded files (B-202)** — `CLAUDE.md`, `AGENTS.md` (and the symlink),
  `index.md`, `log.md`, `README`/license are NOT flagged for a missing `type`; a
  freshly scaffolded-then-promoted bundle passes `kb okf check`.
- **Reject non-okf target (N-202)** — `promote --to <wiki-topic>` errors before writing.
- **`--type` vocabulary** — in-vocabulary (ok), out-of-vocabulary (warning),
  empty (error).
- **mode normalization** — empty/absent `mode` → `wiki`; invalid explicit value
  rejected at scaffold.
- **Conformance rules** — missing `type` (error), reserved files excluded,
  `index.md`/`log.md` shape, `--strict` promotion.
- **Cross-concept wikilink mapping** — `[[srcSlug/Foo Bar|Foo]]` → `[Foo](foo-bar.md)`;
  assert the exact emitted relative path; filename collision → `-2` suffix;
  unresolved target recorded in `ConceptResult.UnresolvedLinks`.
- **OKF render smoke (N-001)** — render one document for a `mode: okf` topic and
  assert the OKF link branch and its `fromDir` plumbing execute and emit a relative
  markdown link, so the wired render branch is live, not shipped dead.

### Integration Tests (`//go:build integration`)

- **Conformance against vendored official bundles** (GA4, Stack Overflow,
  crypto_bitcoin) — all pass leniently (ADR-005).
- **Negative fixtures** — missing `type`, unterminated frontmatter, list-not-
  mapping frontmatter → expected errors.
- **Must-tolerate fixtures** — broken link, unknown `type`, frontmatter-less
  `index.md` → no error under lenient mode.
- **End-to-end promote** — scaffold an OKF topic, ingest+compile a small wiki
  topic, `kb promote` a concept, assert: source untouched, concept written with
  four fields + relative links, `log.md` appended, `index.md` updated, and
  `kb okf check` passes.
- **Zero-regression** — an existing wiki topic renders byte-identically and
  `kb lint` output is unchanged.

Use `t.TempDir()` for filesystem isolation; co-locate tests with their package.

## Development Sequencing

### Build Order

1. **`mode` data model** — add `TopicMode` consts + `Mode` on
   `models.TopicMetadata`/`TopicInfo`; add `Mode` to `topicMetadataFile` with
   empty→wiki normalization in `readTopicYAMLMetadata`. *No dependencies.*
2. **`internal/config` `[okf].types`** — add `OKFConfig`, defaults, env/TOML load,
   `config.example.toml`. *No dependencies.*
3. **`LinkFormatter` in `internal/vault`** — interface + `WikiLinkFormatter`
   (wraps current behavior) + `OKFLinkFormatter` + `LinkFormatterFor`. *Depends on
   step 1 (reads `topic.Mode`).*
4. **Migrate the ~43 call sites** through `linkFor(topic, fromDir, target, label)`;
   add golden wiki-output tests. *Depends on step 3.*
5. **OKF scaffold** — `internal/topic` mode-conditional dirs + OKF `CLAUDE.md`/
   `index.md`(`okf_version`)/`log.md` templates. Extend `topic.New`/`newWithDate`
   with a `mode` parameter (or a `NewWithMode` variant) and write an explicit `mode`
   value into `topic.yaml`. *Depends on step 1.*
6. **`internal/okf` emission primitives** — frontmatter contract, wikilink→markdown
   transform (reusing step 3's `OKFLinkFormatter`), `index.md`/`log.md` generation.
   *Depends on steps 2, 3.*
7. **`promote` orchestration** — read source, remap, transform, write concept,
   append log, regenerate index. *Depends on steps 5, 6.*
8. **Conformance checker** — recursive walk, §9 lenient + local warnings +
   `--strict`, emitting `models.LintIssue`. *Depends on step 2 (vocabulary).*
9. **CLI wiring** — `topic new --mode`, `kb promote`, `kb okf check`. *Depends on
   steps 5, 7, 8.*
10. **Fixtures + integration tests** — vendor official bundles, negatives,
    must-tolerate, end-to-end promote, zero-regression. *Depends on steps 7, 8, 9.*

### Technical Dependencies

- None external. Vendoring the official bundles requires keeping their Apache-2.0
  license/attribution under `testdata/`.

## Monitoring and Observability

CLI tool — observability is structured logging via `internal/logger` (slog) and
command output:

- **Log events** (debug/info): `promote` source/target/type, link-transform count,
  broken-link count, `promote rejected: target not okf` (N-202), `description fallback
  used`; `okf check` files-scanned, `reserved file skipped` (N-205), errors, warnings,
  strict flag.
- **Exit codes**: `okf check` returns non-zero on errors (and on warnings under
  `--strict`) for CI gating.
- **Report output**: `okf check` renders via `internal/output` (table/json/tsv),
  consistent with `kb lint`.

## Technical Considerations

### Key Decisions

- **Cohesive `internal/okf` package** (ADR-003). *Rationale:* one testable home,
  clean phase growth. *Trade-off:* a new package vs maximal reuse.
- **Full `LinkFormatter` migration in the MVP** (ADR-004). *Rationale:* one
  finished mode-aware seam. *Trade-off:* large refactor whose OKF render branch is
  dormant until Phase 3; mitigated by golden tests.
- **Dedicated conformance checker, lenient + `--strict`** (ADR-002/005).
  *Rationale:* OKF is structurally unlike the wiki pyramid and is permissive by
  spec. *Trade-off:* two checking engines coexist.
- **Non-LLM mechanical `promote`** (ADR-002). *Rationale:* preserves `kb`'s
  identity and determinism. *Trade-off:* distillation stays a human/agent job.
- **`kb promote` top-level + `kb okf` group** (ADR-006). *Rationale:* prominence
  for the daily verb, cohesive namespace for the ecosystem.

### Known Risks

- **Wiki render regression** during the 43-site migration (medium likelihood,
  high impact) → byte-identical golden tests + keep the `vault` suite green;
  default formatter is wiki when `mode != okf`.
- **OKF relative-path correctness** (medium) → exhaustive table-driven tests
  including parent/sibling/root.
- **OKF v0.1 is a draft** (medium) → isolate the ruleset behind a versioned module;
  track the spec repo.
- **Vocabulary ignored in practice** (low) → `--strict` in CI enforces it where it
  matters.
- **Vendored-fixture drift / license** (low) → pin a commit, retain attribution.

## Architecture Decision Records

- [ADR-001: OKF as a first-class authoring mode, not an export-only projection](adrs/adr-001.md) — OKF is a first-class per-topic mode with phased delivery (product).
- [ADR-002: promote is a mechanical non-LLM transform; type governance and a spec-compatible conformance posture](adrs/adr-002.md) — mechanical/non-destructive promote, `--type` + `kb.toml` vocabulary, emit-four + relative links, lenient validation (product).
- [ADR-003: A cohesive `internal/okf` package and a `mode` field on the topic model](adrs/adr-003.md) — one-way-dependency OKF package + `mode` on the topic model (default wiki).
- [ADR-004: A `LinkFormatter` abstraction and full migration of the link call sites](adrs/adr-004.md) — mode-aware `LinkFormatter` in `vault`; migrate all ~43 sites in the MVP; OKF render branch dormant until Phase 3.
- [ADR-005: A dedicated OKF conformance checker with vendored official fixtures](adrs/adr-005.md) — separate checker in `internal/okf`, lenient §9 + local warnings + `--strict`, vendored official bundles as fixtures.
- [ADR-006: CLI surface — `kb promote` top-level and a `kb okf` command group](adrs/adr-006.md) — top-level `kb promote`, `kb okf check` group, `topic new --mode`.
