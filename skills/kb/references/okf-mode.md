# OKF Dual-Mode Reference

`kb` runs two distinct knowledge lifecycles, selected per topic by a `mode` field:

- **`wiki`** (default) — the Karpathy research **lab**: `ingest → compile → query → lint`, Obsidian `[[wikilinks]]`, the `raw/`+`wiki/`+`outputs/` pyramid. Documented across the rest of this skill.
- **`okf`** — a portable **catalog** conforming to the **Open Knowledge Format (OKF) v0.1** (Google Cloud, 2026): a flat bundle of typed concepts, plain relative markdown links, meant to be *declared and shared* so other people's agents and tools can consume it.

They are two lifecycles, not one derived from the other. The daily high-leverage move is **distilling research into the catalog** — `kb promote` turns a finished wiki concept into a typed OKF concept; `kb okf check` validates the bundle.

Mode lives in `topic.yaml` as `mode: wiki|okf`. Absent/empty `mode` normalizes to `wiki`, so every pre-existing topic is unchanged (zero regression).

> **Scope (MVP).** This reference covers the shipped MVP: `mode`, OKF scaffold, `kb promote`, `kb okf check`, and the `[okf].types` vocabulary. Out of scope for now: `kb okf export` (wiki→bundle), `kb okf ingest` (external bundle→research), and codebase→OKF emission. An OKF bundle is **one topic** with concepts **flat at the bundle root**; type-grouped or nested concept directories are a later-phase option.

---

## Bundle layout

`kb topic new <slug> <title> <domain> --mode okf` scaffolds a **flat** bundle (no wiki pyramid, no `.gitkeep` files):

```
<bundle>/
  CLAUDE.md      # OKF-flavored schema document + topic marker
  AGENTS.md      # symlink → CLAUDE.md (Codex parity)
  topic.yaml     # slug, title, domain, mode: okf
  index.md       # generated; frontmatter `okf_version: "0.1"`; body "# OKF Bundle Index"
  log.md         # "# Directory Update Log" + an Initialization entry under "## YYYY-MM-DD"
  <concept>.md   # typed concept files at the bundle root (added by promote / authoring)
```

The topic marker is still `CLAUDE.md` (as for wiki topics). `index.md` and `log.md` are **auto-maintained by `kb`** — regenerated on scaffold and on every `promote`. The OKF `CLAUDE.md` documents the bundle contract and the `kb promote` / `kb okf check` operations instead of the wiki pipeline.

---

## OKF concept frontmatter contract

Every concept `.md` carries the **four producer fields**. `kb` emits them **alphabetically** (`frontmatter.Generate` sorts keys), so a promoted concept looks like:

```yaml
---
description: One-line summary of the concept.
tags:                  # optional — only when the source carried tags
  - example
timestamp: 2026-06-27T14:03:11Z   # RFC3339, UTC
title: Concept Title
type: Voice Profile
---
```

- `type` is the **only** field OKF v0.1 strictly requires. The other three producer fields and the relative-link style are **two deliberate, documented deviations** (ADR-002) that match Google's reference tooling:
  1. emit the four producer fields (the spec mandates only `type`);
  2. emit **relative** markdown links (the spec recommends absolute `/path.md`, which breaks GitHub rendering).
- Do **not** assert a `type`-first field order in tooling — YAML maps are unordered and `kb` sorts keys alphabetically.
- Wiki stage markers (`stage`, wiki `type: wiki`, `domain`, `sources`, `created`/`updated`) are **dropped**; they do not belong in an OKF concept.

### `promote` frontmatter remap

`kb promote` derives each OKF field deterministically from the source wiki document:

| OKF field | Source | Fallback |
| --- | --- | --- |
| `type` | `--type` (validated against `[okf].types`) | none — **required**; error if absent |
| `title` | source `title` | humanized concept key (slug of the source base name) |
| `description` | `--description` flag | first non-empty body sentence (markdown-stripped) → else empty **+ warning** |
| `timestamp` | promote-time clock, RFC3339 / UTC | — (promotion *is* the last meaningful change; source `scraped`/`created` are date-only and not used) |
| `tags` | source `tags` when present | omitted |

`promote` is **mechanical, non-LLM, and non-destructive** — the source wiki document is never modified. Any intelligent rewriting/condensing is left to the operator or an external agent.

---

## Concept path & link model

A **single canonical key** drives both the concept filename and every inbound link, so they can never diverge.

- **Canonical key:** `vault.SlugifySegment(base(source-path))` — the existing slug primitive (lowercase `[a-z0-9-]`, drops accents so `Conteúdo` → `conte-do`, empty → `item`).
- **Filename:** a concept promoted from `SourceDocPath` is written to `<bundle-root>/<key>.md`. Filename collisions get a numeric suffix (`-2`, `-3`, …).
- **Wikilink → markdown transform:** each `[[target|label]]` in the body becomes a relative markdown link. The target is normalized (strip a leading `/`, strip the source topic slug prefix, slugify the base name) to `<key>.md`, i.e. the **same** key the target concept's own filename uses — so a link to an already-promoted concept resolves to its real file. With a flat root the link is `[label](<key>.md)`.
  - A `#anchor` fragment is preserved: `[label](<key>.md#anchor)`.
  - A missing label falls back to the humanized key.
- **Unresolved targets:** when no promoted counterpart exists yet, the link is still emitted to the would-be path (tolerated per OKF §9) and recorded in `ConceptResult.UnresolvedLinks` for follow-up.

### `kb promote` JSON result (`ConceptResult`)

```json
{
  "writtenPath": "voice-profile.md",
  "type": "Voice Profile",
  "linksRewritten": 4,
  "unresolvedLinks": ["offer-ladder.md"],
  "warnings": ["type \"Voice Profile\" is outside the configured OKF vocabulary"]
}
```

---

## Auto-maintained `index.md` and `log.md`

**`index.md`** is regenerated on every `promote`. It carries only `okf_version` frontmatter (preserved from the existing root index, defaulting to `0.1`) and a body that groups concepts by `type`, sorted by type then title then path:

```markdown
---
okf_version: "0.1"
---
# OKF Bundle Index

## Offer

* [Founder Offer](founder-offer.md) - The flagship offer for founders.

## Voice Profile

* [Brand Voice](brand-voice.md) - How the brand sounds across channels.
```

Bullets use relative markdown links; the trailing `- <description>` is included when the concept has one. **Do not hand-edit `index.md`** — your edits are overwritten on the next promote.

**`log.md`** uses ISO-date `## YYYY-MM-DD` headings, newest first, with bullet entries:

```markdown
# Directory Update Log

## 2026-06-27

* **Creation**: Promoted [Brand Voice](brand-voice.md) from `research/wiki/concepts/Brand Voice.md`.

## 2026-06-20

* **Initialization**: Created OKF bundle `second-brain` for `operations`.
```

`kb promote` inserts the entry under the current date heading (creating it if needed). No manual log entry is required after `promote`.

---

## Conformance checking (`kb okf check`)

```bash
kb okf check <okf-topic> [--strict] [--format table|json|tsv]
```

`kb okf check` walks the bundle and validates OKF v0.1 **leniently** (§9): it tolerates broken cross-links, unknown `type` values, and missing optional fields, so externally produced bundles pass. Diagnostics are emitted as columns `severity, kind, filePath, target, message`. The command **exits non-zero** when any `error` is present (and, under `--strict`, when any `warning` is present too) — suitable as a CI gate.

### Excluded files (never require a `type`)

`index.md`, `log.md`, `CLAUDE.md`, `AGENTS.md`, symlinks, and `README` / `LICENSE` / `NOTICE` / `ATTRIBUTION` markdown files, plus dot-directories. Everything else under the bundle is a concept file.

### Errors (always fail; `severity=error`)

- A concept `.md` with **missing or empty `type`**.
- A concept (or `index.md`) with **unparseable frontmatter**, or a concept with **no leading YAML frontmatter**.
- A non-root file named `index.md` declaring `okf_version`, or the root `index.md` declaring any key other than `okf_version`.
- A `log.md` `## ` heading that is not a valid `YYYY-MM-DD` date.

### Warnings (lenient by default; become errors under `--strict`)

- A concept missing a producer field: `title`, `description`, or `timestamp`.
- A concept `type` **outside** the configured `[okf].types` vocabulary (only checked when the vocabulary is non-empty).

---

## Type vocabulary (`[okf].types`)

The local concept-type standard, configured in `kb.toml`:

```toml
[okf]
# Local OKF concept type vocabulary. Empty (the default) means neither
# `kb promote` nor `kb okf check` warns about unknown types.
types = ["Voice Profile", "Offer", "Playbook"]
```

OKF has no global type registry by design; this is a **local** standard that prevents drift (`Voice Profile` vs `voice-profile`). When the list is empty, the type check is a no-op (no false warnings). When non-empty, an off-vocabulary `--type` warns at `promote` time and at `okf check` time; `--strict` makes it a CI-failing error. Extend the standard by editing `kb.toml`, not by inventing drifted variants.

---

## CLI surface

| Command | Args / Flags | Behavior |
| --- | --- | --- |
| `kb topic new <slug> <title> <domain>` | `--mode wiki\|okf` (default `wiki`) | Wiki scaffold (unchanged) or flat OKF bundle. |
| `kb promote <wiki-doc>` | `--to <okf-topic>` (req), `--type <T>` (req), `--description <text>` | Mechanical, non-destructive wiki→OKF concept; `--to` must be a `mode: okf` topic. Emits `ConceptResult` JSON. |
| `kb okf check <okf-topic>` | `--strict`, `--format table\|json\|tsv` | OKF v0.1 conformance + local-standard warnings; non-zero exit on errors (and warnings under `--strict`). |

> The `kb okf` group grows in later phases (`kb okf export`, `kb okf ingest`), which are out of MVP scope.
