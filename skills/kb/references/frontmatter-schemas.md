# Frontmatter Schemas

All notes in the vault use YAML frontmatter for note metadata. Topic identity itself lives in `topic.yaml` at the topic root; the `domain` field on notes is a shortcut for Bases and qmd queries.

Conventions:

- `domain: <short-slug>` identifies the topic (e.g., `ai` for `ai-harness/`).
- `created` and `updated` use ISO date format `YYYY-MM-DD`.
- `tags` always include the domain plus the note type plus topic-specific tags.
- `sources` entries are wikilinks pointing at files in `raw/`.
- `kb` adds its own keys to sources and articles (see *kb-owned keys* below). Write only the keys in the schemas; leave the kb-owned ones to `kb classify`, `kb link` and the ingest gates.

---

## Topic metadata — `<topic>/topic.yaml`

```yaml
slug: harness/goclaw
title: Goclaw
domain: goclaw
mode: wiki   # wiki (default) | okf — omit/empty normalizes to wiki
contract:    # selection contract (wiki topics); see selection-contract.md
  purpose: ""
  core: []
  adjacent: []
  collected_on_purpose: []
  collected_on_purpose_paths: []   # optional topic-relative globs
  out_of_scope: []
# contract_draft:    written by `kb topic contract --draft|--import-claude`, activated by `--accept`
# vocabulary_draft:  written by `kb topic vocabulary --draft`, turned into stub articles by `--accept`
decisions:           # optional per-topic decision settings
  mode: shadow       # body-link insertion: shadow | apply
  gates: shadow      # relevance gates: shadow | apply (default shadow until calibrated)
  relevance: on      # off for topics whose membership is defined by construction
  thresholds: {}     # named threshold overrides; `kb review calibrate --write` fills them
  exclude: []        # globs never sent to any model
```

`contract`, `contract_draft`, `vocabulary_draft` and `decisions` are edited through `kb topic contract`, `kb topic vocabulary` and `kb review calibrate`, or by hand in `topic.yaml`; kb rewrites only the keys it touches and keeps comments and key order. The `## Selection contract` section of the topic `CLAUDE.md` is rendered from `contract:` and overwritten on every accept, so never edit it there.

`topic.yaml` is the primary source of truth for topic title/domain. `CLAUDE.md` remains the topic marker and schema document, but prose formatting in `CLAUDE.md` should not be treated as canonical metadata when `topic.yaml` exists.

`mode` selects the topic's lifecycle: `wiki` (the research lab schemas below) or `okf` (a portable OKF bundle — see *OKF concept* and `references/okf-mode.md`). An absent or empty `mode` normalizes to `wiki`.

---

## Wiki article — `<topic>/wiki/concepts/<Article Title>.md`

```yaml
---
title: Article Title
type: wiki
stage: compiled
domain: <topic-domain>
tags:
  - <topic-domain>
  - wiki
  - topic-specific-tag
  - another-topic-tag
created: YYYY-MM-DD
updated: YYYY-MM-DD
sources:
  - "[[Source File Name]]"
  - "[[Another Source]]"
aliases:            # optional; Obsidian aliases, merged by kb (never removed)
  - ATA
---
```

`stage` is `compiled` for agent-written articles and `stub` for articles created by `kb topic vocabulary --accept` or an accepted concept proposal; a stub carries `criterion` and an empty `sources: []` until the agent compiles it (then set `stage: compiled`). After `kb classify`, articles also carry `criterion`, `summary`, `genre`, `concepts` and `entities`, and after `kb link` the relation lists (see *kb-owned keys*).

## Raw article — `<topic>/raw/articles/<slug>.md`

```yaml
---
title: Descriptive Title
type: source
stage: raw
domain: <topic-domain>
source_kind: article
source_url: https://example.com/article
scraped: YYYY-MM-DD
tags:
  - <topic-domain>
  - raw
  - topic-specific-tag
---
```

`source_kind` values for current CLI ingests: `article`, `document`, `github-readme`, `youtube-transcript`, `instagram-video`, `bookmark-cluster`, `codebase-file`, `codebase-symbol`.

Every `kb ingest` (except `codebase`) also writes `ingest_batch` and `ingest_query`, and the gates and `kb classify` add the facet keys of *kb-owned keys* below (`triage`, `relevance`, `genre`, `summary`, ...). Codebase snapshot documents (`codebase-file`, `codebase-symbol`) never get kb-owned keys.

## GitHub README — `<topic>/raw/github/<slug>.md`

```yaml
---
title: Repository or Doc Title
type: source
stage: raw
domain: <topic-domain>
source_kind: github-readme
source_url: https://github.com/owner/repo
scraped: YYYY-MM-DD
tags:
  - <topic-domain>
  - raw
  - github
  - topic-specific-tag
---
```

## YouTube transcript — `<topic>/raw/youtube/<slug>.md`

```yaml
---
title: Video Title
type: source
stage: raw
domain: <topic-domain>
source_kind: youtube-transcript
source_url: https://www.youtube.com/watch?v=<video-id>
scraped: YYYY-MM-DD
view_count: 3271
like_count: 77
comment_count: 11
upload_date: YYYY-MM-DD
duration: 6441
duration_string: 1:47:21
channel: Changelog
channel_id: UCZb...
uploader_id: "@Changelog"
channel_follower_count: 20000
categories:
  - Science & Technology
youtube_tags:
  - go
language: en
live_status: not_live
was_live: false
chapter_count: 17
transcript_source: captions
transcription_policy: captions
transcript_language: en
caption_kind: manual
tags:
  - <topic-domain>
  - raw
  - youtube-transcript
---
```

The YouTube metric fields come from `yt-dlp --dump-single-json`. Engagement counts, `channel_follower_count`, `was_live`, and other scalar source-controlled fields may be `null` when YouTube hides or omits them. `categories` and `youtube_tags` are always lists and are empty when absent. `chapter_count` is always an integer count; the CLI stores only that count, not full chapter content. `youtube_tags` stores video keywords because `tags` is reserved for KB taxonomy.

`raw/youtube/` is the canonical transcript directory. Legacy `raw/transcripts/` content should be moved with `kb migrate transcripts --topic <topic-id>`, not treated as a second valid layout.

`kb ingest channel` writes the same `youtube-transcript` schema into `raw/youtube/` — one document per upload — so it shares this section's layout and fields.

## Instagram video — `<topic>/raw/instagram/<slug>.md`

```yaml
---
title: Reel Title
type: source
stage: raw
domain: <topic-domain>
source_kind: instagram-video
source_url: https://www.instagram.com/reel/<shortcode>/
scraped: YYYY-MM-DD
shortcode: <shortcode>
uploader: account-name
uploader_id: account-name
like_count: 4200
comment_count: 88
view_count: 150000
upload_date: YYYY-MM-DD
duration: 32
duration_string: "0:32"
language: en
transcript_source: stt
transcription_policy: auto
transcript_language: en
caption_kind: manual
stt_provider: openai
stt_model: gpt-4o-transcribe
tags:
  - <topic-domain>
  - raw
  - instagram-video
---
```

Instagram documents (`kb ingest instagram`) come from the same `yt-dlp` + STT engine as YouTube. The **body** is the post caption under a `## Caption` heading followed by the spoken transcript under `## Transcript`. Engagement scalars (`like_count`, `comment_count`, `view_count`) may be `null` when Instagram omits them. `transcription_policy` defaults to `auto`; `transcript_source` is `captions`, `stt`, or `none` (caption-only fallback for a music-only reel, in which case there is no `## Transcript` section). `caption_kind`, `stt_provider`, `stt_model`, and `transcript_language` appear only when applicable. `shortcode` is the Instagram media shortcode; `uploader`/`uploader_id` are the account.

## Bookmark cluster — `<topic>/raw/bookmarks/<Topic> Bookmarks <Subtopic>.md`

```yaml
---
title: <Topic> Bookmarks <Subtopic>
type: source
stage: raw
domain: <topic-domain>
source_kind: bookmark-cluster
status: seeded
created: YYYY-MM-DD
updated: YYYY-MM-DD
source_urls:
  - https://twitter.com/user/status/123
  - https://twitter.com/user/status/456
tags:
  - <topic-domain>
  - bookmarks
  - raw
  - topic-specific-tag
---
```

`status` values: `seeded`, `enriched`, `archived`.

## Research output — `<topic>/outputs/queries/<YYYY-MM-DD> <slug>.md`

```yaml
---
title: Output Title
type: output
stage: query
domain: <topic-domain>
tags:
  - <topic-domain>
  - output
  - query
  - topic-specific-tag
created: YYYY-MM-DD
updated: YYYY-MM-DD
informed_by:
  - "[[Wiki Article 1]]"
  - "[[Wiki Article 2]]"
---
```

`stage` values for outputs: `briefing`, `query`, `diagram`, `lint-report`.

## Lint report — `<topic>/outputs/reports/<YYYY-MM-DD>-lint.md`

```yaml
---
title: Lint Report YYYY-MM-DD
type: output
stage: lint-report
domain: <topic-domain>
tags:
  - <topic-domain>
  - output
  - lint-report
created: YYYY-MM-DD
issues_found: N
issues_fixed: M
---
```

## Topic index — Dashboard / Concept Index / Source Index

These files are human-browsed hubs, not research notes. Keep frontmatter minimal:

```yaml
---
title: Dashboard
type: index
domain: <topic-domain>
updated: YYYY-MM-DD
---
```

## OKF concept — `<okf-bundle>/<concept>.md` (mode: okf)

OKF bundles use a different contract from wiki topics. Each concept carries the **four producer fields**, emitted **alphabetically** (`kb` sorts frontmatter keys), with optional `tags`:

```yaml
---
description: One-line summary of the concept.
tags:
  - example   # optional — only when the source carried tags
timestamp: 2026-06-27T14:03:11Z   # RFC3339, UTC
title: Concept Title
type: Voice Profile               # OKF concept type (the only field OKF v0.1 requires)
---
```

Concept bodies use **relative markdown links** (`[label](other-concept.md)`), never `[[wikilinks]]`. Wiki stage markers (`stage`, wiki `type`, `domain`, `sources`, `created`/`updated`) are dropped on promote. The bundle-root `index.md` carries only `okf_version: "0.1"` frontmatter and is regenerated by `kb`. See `references/okf-mode.md` for the full contract, the promote remap table, and the conformance ruleset.

## kb-owned keys

`kb` writes these keys on sources and articles, plain and flat at the top level so Obsidian, Dataview and Bases read them directly. The list lives once in code (`decisions.OwnedKeys`).

| Key | On | Shape | Written by |
|-----|----|-------|------------|
| `triage` | source | `kept` \| `review` \| `quarantined` | ingest gates |
| `triage_reason` | source (not kept) | `off_topic`, `paywall`, `error_page`, `thin`, `not_an_article`, `no_speech`, `duplicate` | ingest gates |
| `genre` | source, article | `paper`, `article_or_essay`, `tutorial_or_guide`, `reference_docs`, `announcement_or_news`, `opinion_or_discussion`, `talk_or_interview`, `repository_or_code`, `dataset_or_benchmark`, `other` | `kb classify` |
| `depth` | source | 0–3, one decimal (0 mention, 1 overview, 2 detailed, 3 primary) | `kb classify` |
| `relevance` | source | `core` \| `adjacent` \| `collected_on_purpose` \| `general` \| `off_topic` \| `unknown` | `kb classify` (`collected_on_purpose` from `contract.collected_on_purpose_paths`) |
| `concepts` | source, article | list of quoted wikilinks to articles | `kb classify` |
| `summary` | source, article | ≤ 400 characters | `kb classify` (generation model) |
| `criterion` | article | ≤ 300 characters, definitional: what a document must discuss to cover the concept; no dates, no market claims, no "this article" | `kb classify` (generation model) |
| `entities` | source, article | ≤ 15 names copied verbatim from the body | `kb classify` (generation model, verified by code) |
| `questions` | source | 3–6 questions the document answers | `kb classify` (generation model) |
| `quality` | source | `thin`, `error_page`, `paywall`, `not_an_article`, `no_speech` | gates, `kb classify` |
| `ingest_batch`, `ingest_query` | source | ingest run (`<command>-<YYYY-MM-DD>-<run id>` or `--batch`) and the query, URL or file that produced the item | `kb ingest` (once, at creation) |
| `related`, `extends`, `prerequisite`, `example_of`, `contradicts` | source, article | list of quoted wikilinks (`"[[File name]]"`) | `kb link` |
| `affects` | source | list of quoted wikilinks to articles the source should change | `kb link` |
| `supersedes` | source | list of quoted wikilinks to older sources | near-duplicate gate |
| `aliases` | article | list of strings (acronym, synonyms, other-language names) | `kb classify` when an article has none; merged, never removed |
| `recaptured` | source | date the body was replaced by a fresh recapture | `kb review accept --queue recapture` |

User-owned control key:

| Key | On | Shape | Effect |
|-----|----|-------|--------|
| `locked` | any | `true` | kb writes nothing to this file (no owned keys, no body links); reported as `skipped:locked` |

Rules:

- **Ownership, not a prefix.** kb writes an owned key only when it is absent, or when its current value still hashes to what kb last wrote (value hashes live in `<topic>/.decisions/state.jsonl`, not in frontmatter). A key you edit, or one your vault already used for something else, becomes yours: kb leaves it, reports `skipped:user-key` in the run summary, and `kb lint` reports `key-conflict`. This is how an edited `criterion` or `summary` is kept. Delete the key to hand it back to kb.
- **Relation lists merge.** In a `related`, `extends`, `prerequisite`, `example_of`, `contradicts`, `affects` or `supersedes` list you wrote, kb only appends new targets (deduplicated by link target); it never removes or reorders your items, and lint does not call the list a `key-conflict`. `concepts` keeps the ownership rule above.
- **`aliases` merge.** kb adds validated aliases (1–6 words, ≤ 8, no collision with another article's title or alias) and never removes one you wrote. Aliases are display text for `[[target|alias]]`; like titles, they never resolve a link on their own.
- **No renamed user keys.** kb uses `triage` and `genre` because `status` and `kind` are common user keys; `status` (bookmarks) and `source_kind` keep their meaning.
- **Byte-preserving writes.** kb replaces or appends only its own keys, never re-serializes the rest of the YAML, and refuses to write when the body or a key it does not own changed since it read the file (`skipped:changed`).
- **Probabilities stay out.** Frontmatter holds banded outcomes and the `depth` expectation only; probabilities live in `<topic>/.decisions/receipts.jsonl`.
- **Ingest extras cannot set owned keys** (or `locked`).
- Quoted wikilinks use the file stem (`"[[File name]]"`); when the stem is ambiguous in the vault, the vault-relative path without `.md`.

## Quick reference

| File type | Path | type | stage |
|-----------|------|------|-------|
| Wiki article | `wiki/concepts/` | `wiki` | `compiled` (or `stub` before compilation) |
| Raw article | `raw/articles/` | `source` | `raw` |
| Raw GitHub | `raw/github/` | `source` | `raw` |
| Raw YouTube | `raw/youtube/` | `source` | `raw` |
| Raw Instagram | `raw/instagram/` | `source` | `raw` |
| Raw bookmarks | `raw/bookmarks/` | `source` | `raw` |
| Quarantined source | `raw/_quarantine/<path under raw/>` | `source` | `raw` (with `triage: quarantined`) |
| Briefing | `outputs/briefings/` | `output` | `briefing` |
| Query result | `outputs/queries/` | `output` | `query` |
| Diagram | `outputs/diagrams/` | `output` | `diagram` |
| Lint report | `outputs/reports/` | `output` | `lint-report` |
| Index | `wiki/index/` | `index` | — |
| OKF concept | `<okf-bundle>/` | `<concept type>` | — (uses `timestamp`, not `stage`) |
| OKF bundle index | `<okf-bundle>/index.md` | — (`okf_version` only) | — |
