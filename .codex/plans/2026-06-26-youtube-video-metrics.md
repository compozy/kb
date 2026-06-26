# YouTube Video Metrics In Frontmatter

## Summary

Enrich new YouTube ingests so `raw/youtube/*.md` receives video metrics and context directly in YAML frontmatter, using data already returned by `yt-dlp --dump-single-json`. Keep the current pipeline (`yt-dlp -> youtube.Result -> ingest.ExtraFrontmatter`) and do not add flags, dependencies, backfill, or a generated Base/list.

## Public Schema And Type Changes

Add the following optional flat fields to `youtube-transcript` frontmatter:

| Field | Value |
| --- | --- |
| `view_count` | integer or `null` |
| `like_count` | integer or `null` |
| `comment_count` | integer or `null` |
| `upload_date` | `YYYY-MM-DD` or `null` |
| `duration` | duration in seconds or `null` |
| `duration_string` | readable string, for example `1:47:21`, or `null` |
| `channel` | channel name or `null` |
| `channel_id` | stable channel ID or `null` |
| `uploader_id` | uploader handle/id or `null` |
| `channel_follower_count` | approximate integer or `null` |
| `categories` | string list, empty when absent |
| `youtube_tags` | YouTube video keywords/tags; use this name because `tags` is reserved by KB |
| `language` | video language or `null` |
| `live_status` | raw `yt-dlp` status, for example `not_live`, or `null` |
| `was_live` | boolean when known, otherwise `null` |
| `chapter_count` | chapter count; do not persist the full chapter payload |

Keep existing transcript provenance fields: `transcript_source`, `transcription_policy`, `transcript_language`, `caption_kind`, `stt_provider`, and `stt_model`.

## Implementation Changes

- Expand `internal/youtube`:
  - Add corresponding fields to `ytDLPInfo` with `yt-dlp` JSON tags.
  - Extend `youtube.Metadata` with normalized values.
  - Use pointers for optional counts and `was_live`, preserving real `0`/`false` values and serializing `null` when unknown.
  - Convert `upload_date` from `YYYYMMDD` to `YYYY-MM-DD`.
  - Persist `duration` as integer seconds and use `duration_string` from `yt-dlp`, with a calculated fallback from seconds when empty.
  - Calculate only `chapter_count = len(chapters)`.

- Update `internal/cli/ingest_youtube.go`:
  - Make `youtubeFrontmatter` include video metrics plus transcript provenance.
  - Keep `tags` untouched; video keywords go only into `youtube_tags`.
  - Do not change error behavior: if the `yt-dlp` metadata pass fails, ingest still fails before writing a partial file.

- Update validation and documentation:
  - In `internal/lint`, add `upload_date` as an optional date field for `raw/youtube/` and `categories`/`youtube_tags` as optional string lists.
  - Update `skills/kb/references/frontmatter-schemas.md` with the full YouTube schema.
  - Update `README.md` and `.compozy/tasks/kb-pivot/_techspec.md` to document engagement, publication, and classification metadata on `ingest youtube`.

## Test Plan

- `internal/youtube`:
  - Update the fake `yt-dlp` metadata JSON to return all tier 1-3 fields.
  - Assert normalization of counts, `upload_date`, `duration`, `duration_string`, channel fields, classification, live status, and `chapter_count`.
  - Cover `null`/missing values for `like_count`, `comment_count`, `channel_follower_count`, and `was_live`.

- `internal/cli`:
  - Update `TestYouTubeFrontmatterIncludesTranscriptProvenance` or create an adjacent test to prove `youtubeFrontmatter` produces metrics plus provenance.
  - Assert `youtube_tags` is used and reserved `tags` is not overwritten.
  - Update the `ingest youtube` command test to verify `ExtraFrontmatter` receives metrics.

- `internal/ingest` / `internal/lint`:
  - Expand the `ExtraFrontmatter` test to prove serialization of integers, bools, lists, and `null`.
  - Add lint regression for valid/invalid `upload_date` and optional `categories`/`youtube_tags` lists.

- Final verification:
  - `rtk go test ./internal/youtube ./internal/cli ./internal/ingest ./internal/lint`
  - `rtk make verify`
  - Run `$cy-impl-peer-review` until SHIP.

## Assumptions

- The improvement applies to new ingests only; no automatic migration/backfill for old transcripts.
- No new Obsidian Base or generated list will be created; new fields are available for existing lists, Dataview, QMD, Bases, or manual filters.
- The schema is flat, with the only exception that YouTube keywords are named `youtube_tags` because `tags` is reserved.
- `channel_follower_count` is approximate and preserves the raw value returned by `yt-dlp`.
