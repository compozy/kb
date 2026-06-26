# Frontmatter Schemas

All notes in the vault use YAML frontmatter for note metadata. Topic identity itself lives in `topic.yaml` at the topic root; the `domain` field on notes is a shortcut for Bases and qmd queries.

Conventions:

- `domain: <short-slug>` identifies the topic (e.g., `ai` for `ai-harness/`).
- `created` and `updated` use ISO date format `YYYY-MM-DD`.
- `tags` always include the domain plus the note type plus topic-specific tags.
- `sources` entries are wikilinks pointing at files in `raw/`.

---

## Topic metadata — `<topic>/topic.yaml`

```yaml
slug: harness/goclaw
title: Goclaw
domain: goclaw
```

`topic.yaml` is the primary source of truth for topic title/domain. `CLAUDE.md` remains the topic marker and schema document, but prose formatting in `CLAUDE.md` should not be treated as canonical metadata when `topic.yaml` exists.

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
---
```

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

`source_kind` values for current CLI ingests: `article`, `document`, `github-readme`, `youtube-transcript`, `bookmark-cluster`, `codebase-file`, `codebase-symbol`.

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

## Quick reference

| File type | Path | type | stage |
|-----------|------|------|-------|
| Wiki article | `wiki/concepts/` | `wiki` | `compiled` |
| Raw article | `raw/articles/` | `source` | `raw` |
| Raw GitHub | `raw/github/` | `source` | `raw` |
| Raw YouTube | `raw/youtube/` | `source` | `raw` |
| Raw bookmarks | `raw/bookmarks/` | `source` | `raw` |
| Briefing | `outputs/briefings/` | `output` | `briefing` |
| Query result | `outputs/queries/` | `output` | `query` |
| Diagram | `outputs/diagrams/` | `output` | `diagram` |
| Lint report | `outputs/reports/` | `output` | `lint-report` |
| Index | `wiki/index/` | `index` | — |
