# Plan: Native `kb ingest channel` — bulk YouTube transcript ingestion with rate-limit handling

## Context

`kb ingest youtube` is strictly single-video today (`cobra.ExactArgs(1)` in
`internal/cli/ingest_youtube.go:32`, and `--no-playlist` hardcoded in
`internal/youtube/ytdlp.go:335`). To transcribe an entire channel the user runs
the `kb-yt-channel` skill, whose Python script (`.agents/skills/kb-yt-channel/scripts/ingest-channel.py`)
enumerates videos with `yt-dlp --flat-playlist` and then loops calling
`kb ingest youtube` once per video — with **no throttle, no concurrency cap, no
adaptive backoff**, which is exactly what triggers YouTube rate-limiting (429/403)
on large channels.

We will move the *mechanics* (enumeration + resumable, throttled, back-off-aware
bulk ingest) into the Go binary as a dedicated `kb ingest channel <url>` command
that ingests into an **existing** topic and prints a JSON summary. The skill keeps
owning all *organization* (creating `yt-channels/<slug>`, `topic.yaml`, CLAUDE.md
patching, wiki indexes, report, validation) and simply swaps its per-video loop for
one call to `kb ingest channel`.

Reference (ideas borrowed, NOT the insecure parts): `roundyyy/yt-bulk-subtitles-downloader`
— take channel enumeration, bounded worker pool, resume, per-item failure isolation,
backoff. Explicitly reject its free-proxy scraping (MITM + the 403/409 source) and
1000-thread model.

### Decisions (confirmed with user)
- **Scope:** mechanics only — dump into an existing `--topic`; no topic creation / wiki / report in Go.
- **CLI surface:** dedicated `kb ingest channel <url>`; `kb ingest youtube` stays single-video.
- **Transcription:** default `captions` (recommended for bulk); `--transcribe` still accepts `captions|auto|stt` for parity with the skill's existing interface.
- **Proxy:** keep the single trusted `[youtube].proxy`; no proxy pool. Rate-limit is handled by throttle + low concurrency + adaptive backoff.

## Architecture

Keep Cobra thin (per CLAUDE.md). New logic lives in packages; the command wires them.

```
kb ingest channel <url> --topic <existing> [--limit N|--all] [--transcribe captions]
                        [--concurrency N] [--throttle 2s] [--dry-run]
   │
   ├─ youtube.NormalizeChannelURL(url)              # @handle//shorts//streams → /videos; reject single-video URLs; accept playlist ?list=
   ├─ youtube.ListChannelVideos(ctx, url, limit)    # yt-dlp --flat-playlist --dump-single-json (NO --no-playlist)
   ├─ ingest.ExistingYouTubeVideoIDs(vault, topic)  # resume: read raw/youtube/*.md frontmatter (video_id + source_url)
   ├─ Extractor.BulkExtract(...)                    # bounded worker pool + throttle/jitter + adaptive backoff
   │      └─ per video → Extractor.Extract (reused as-is)
   └─ runIngest per successful Result               # reuse kingest.Ingest + youtubeFrontmatter
   → JSON summary {resolved, ingested[], skipped[], failures[]}
```

## Changes

### 1. `internal/youtube/channel.go` (new) — enumeration + bulk orchestrator
- `ChannelVideo{ VideoID, Title, URL string }`.
- `NormalizeChannelURL(raw string) (string, error)` — port `normalize_channel_url` from the skill (ingest-channel.py:60): handle/`/shorts`/`/streams` → `/videos`, reject `watch?v=`, allow `?list=` playlists.
- `(*ytDLPBackend) listEntries(ctx, url, limit)` — new method using a **playlist-allowed** arg set. Refactor `baseArgs()` (ytdlp.go:332) to split the shared flags (`--ignore-config`, retries, proxy, cookies, user-agent) from `--no-playlist`; add `playlistArgs()` that omits `--no-playlist` and adds `--flat-playlist --dump-single-json` (+ `--playlist-end N` when limit>0). Parse `entries[]` (id/title/url), dedupe by id. Reuse existing `run()` / `commandContext` / `lookPath` seams.
- `BulkOptions{ TranscriptionPolicy, PreferredLanguages, Concurrency, Throttle, BackoffMax, MaxRetries }`.
- `VideoOutcome{ Video, Result *Result, Err error, Skipped bool }`.
- `(*Extractor) BulkExtract(ctx, videos []ChannelVideo, opts BulkOptions, skip func(ChannelVideo) bool, sink func(VideoOutcome) error) error`:
  - bounded worker pool (`Concurrency`, default 2);
  - per-task pre-wait `Throttle` + jitter (`math/rand`);
  - **adaptive backoff:** shared penalty duration; when `Extract` returns a `*youtube.Error` with `ErrorKindNetworkBlocked` (set by `wrapYTDLPCaptionFetchError`, ytdlp.go:322), sleep `penalty` (exponential, capped at `BackoffMax`) and retry that video up to `MaxRetries`; decay penalty on success;
  - calls `sink` for each completed video so the caller persists incrementally; `skip` lets the caller drop already-ingested videos.
  - `Extractor` is stateless/concurrency-safe (confirmed: youtube.go:179) so concurrent `Extract` is fine.

### 2. `internal/ingest/ingest.go` — resume helper + `video_id` frontmatter
- Add `ExistingYouTubeVideoIDs(vaultPath, topicRef string) (map[string]struct{}, error)`: resolve topic, walk `raw/youtube/*.md` (dir from `rawDirectoryForSourceKind`, ingest.go:257), parse frontmatter via existing `frontmatter.Parse`, collect `video_id` field and the `v=` id parsed from `source_url`.
- In `internal/cli/ingest_youtube.go:youtubeFrontmatter` add `"video_id": optionalString(metadata.VideoID)` so resume is reliable (today only `source_url` carries the id). `Metadata.VideoID` already exists (youtube.go:118).

### 3. `internal/cli/ingest_channel.go` (new) — thin command
- `newIngestChannelCommand()` registered in `newIngestCommand()` (cli/ingest.go:38).
- Flags: `--topic` (required, reuse `requireTopicFlag`), `--transcribe` (default captions), `--limit int` (0/unset = all), `--all` (alias for limit 0), `--concurrency int`, `--throttle duration`, `--dry-run`.
- Seam for tests: `var newYouTubeChannelLister = func(cfg) channelLister` mirroring `newYouTubeTranscriptExtractor` (ingest_youtube.go:21).
- Flow: `resolveIngestTarget` (existing topic) → normalize+enumerate → `ExistingYouTubeVideoIDs` → `BulkExtract` with `sink` = `runIngest` (reuse `kingest.Options` + `youtubeFrontmatter`) → aggregate → `writeJSON` summary.
- `--dry-run`: enumerate + print summary (`resolved`, `videos[]`) without ingesting (mirror skill dry-run).
- Summary shape consumable by the skill: `{ topic, channelUrl, normalizedChannelUrl, selection, transcribe, resolved, ingested[], skipped[], failures[] }` with `{video_id,title,url[,filePath|error]}` entries.

### 4. `internal/config` — bulk knobs
- Add to `YouTubeConfig` (config.go:93): `BulkConcurrency int` (default 2), `BulkThrottle string` (default `"2s"`), `BulkBackoffMax string` (default `"60s"`), `BulkRetries int` (default 3), with parse helpers like `RetryBackoffDuration`. Wire defaults (config.go:137), `Validate` (concurrency ≥1, durations parseable), and document in `config.example.toml` + the `[youtube]` notes in `CLAUDE.md`.

### 5. Skill delegation (follow-up, edits skill files)
- Update `.agents/skills/kb-yt-channel/scripts/ingest-channel.py`: replace `resolve_videos` + the per-video loop + `existing_video_ids` (ingest-channel.py:94/641) with a single `kb ingest channel <url> --topic yt-channels/<slug> --transcribe <p> [--limit N]` call, consuming its JSON summary to build `ingested/skipped/failures`. Keep all organization (scaffold, topic.yaml, CLAUDE patch, wiki indexes, report, validation) unchanged. Update `SKILL.md` Step 2 accordingly.

## Rate-limit strategy (the actual fix)
1. Low concurrency (default 2, configurable) — never the repo1 1000-thread storm.
2. Throttle + jitter between requests — the highest-impact lever.
3. Adaptive exponential backoff on `ErrorKindNetworkBlocked` (429/403), with per-video bounded retry.
4. Reuse authenticated session (`[youtube].cookies_file` + `user_agent`) and single `[youtube].proxy` — already supported in `baseArgs()`.

## Testing (canonical suites, table-driven, `t.TempDir()`)
- `internal/youtube/channel_test.go`: `NormalizeChannelURL` table; `listEntries` via the existing fake yt-dlp seam (`newFakeYTDLPBackend`, ytdlp_test.go:601) — extend the fake script to emit flat-playlist JSON and assert `--flat-playlist`/no `--no-playlist`/`--playlist-end`; `BulkExtract` concurrency cap, throttle ordering, skip, and **backoff-then-success on a simulated block** (fake extractor returning `ErrorKindNetworkBlocked` once). Assert no weakened behavior — a block that never clears must surface as a failure, not a hang.
- `internal/ingest/ingest_test.go`: `ExistingYouTubeVideoIDs` over fixture `raw/youtube` files (via `video_id` and via `source_url`).
- `internal/cli/ingest_test.go`: channel command happy-path (fake lister + fake extractor + `runIngest` seam via `restoreIngestGlobals`, ingest_test.go:948) asserting JSON summary, resume-skip, failure isolation, and `--dry-run`.

## Verification (end-to-end)
1. `make verify` (fmt → lint(0) → test → build → boundaries) — non-negotiable gate.
2. Real run against a small channel into a throwaway topic:
   - `bin/kb topic new demo-chan "Demo" youtube-channel`
   - `bin/kb ingest channel https://www.youtube.com/@aiDotEngineer/videos --topic demo-chan --limit 3 --transcribe captions --throttle 2s --concurrency 2`
   - assert 3 files under `<vault>/demo-chan/raw/youtube/`, JSON summary lists them.
   - rerun the same command → all 3 in `skipped[]` (resume), 0 new writes.
   - `--dry-run` lists videos without writing.
3. Skill regression: run `python3 .agents/skills/kb-yt-channel/scripts/ingest-channel.py --vault . --channel-url <url> --topic-slug demo --title "Demo" --domain youtube-channel --limit 3 --transcribe captions` and confirm it now delegates to `kb ingest channel`, still producing the `yt-channels/demo` topic + report.

## Out of scope (stays in the skill)
Topic creation under `yt-channels/`, `topic.yaml`, CLAUDE.md/category-doc maintenance,
`wiki/index` Dashboard + Source Index, ingest report, and lint/index/search validation.
