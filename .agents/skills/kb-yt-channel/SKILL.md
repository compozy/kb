---
name: kb-yt-channel
description: "Create and maintain Karpathy KB topics from YouTube channels through scaffolding, bulk transcript ingestion, validation, and indexing. Excludes single videos, general summaries, and other sources."
disable-model-invocation: true
---
# KB YouTube Channel

Use the native `kb ingest channel` for channel resolution, transcripts, dedup/resume, retries, and pacing. The bundled `scripts/ingest-channel.py` adds topic organization under `yt-channels/`, metadata, indexes, validation, and a run report.

1. Resolve the vault, channel/playlist URL, topic slug/title/domain, video selection (`--limit N` or `--all`), and transcript policy (`captions`, `auto`, or `stt`) from the request. Ask only for choices that cannot reasonably be inferred.
2. Check required tools once, reusing current evidence: `kb`, `yt-dlp`, `qmd` when indexing, and `ffmpeg` plus the configured provider when STT is requested. Do not install optional networking dependencies as a default preflight.
3. Run `python3 <kb-yt-channel-dir>/scripts/ingest-channel.py` with `--vault`, `--channel-url`, `--topic-slug`, `--title`, `--domain`, selection, and `--transcribe`. Resolve the helper from this skill directory. Native-language captions are the default; override `--sub-langs` only when needed. `--no-index` defers indexing; `--embed` requests embeddings.
4. A preview is optional. The existing `--dry-run` **creates/updates the topic skeleton and metadata**, then lists the ingest plan without transcripts; it is not a read-only preview. Use it only within authorized topic creation/update. Do not require preview followed by a repeated full preflight.
5. Inspect the helper's JSON summary and failures. Reuse its completed topic-info/lint/index validation; run additional checks only for missing evidence. When searchability is part of the task, verify one representative `kb search` in the topic/collection after indexing.
6. Report actual ingested/skipped/failed counts and the report path. Compile wiki articles only if the task includes synthesis; transcript ingestion alone does not imply compilation.

## Failures and references

Start with the existing network configuration and native bounded retries/backoff. On rate limits, honor server guidance, lower concurrency, or increase pacing. Configure authenticated cookies or a permitted proxy only when an observed access/network problem requires it; never expose credentials or mandate a residential proxy for large channels.

Read `references/channel-topic-contract.md` for metadata/output details and `references/troubleshooting.md` for an actual failure. Resume the same request after repair; native dedup preserves successful ingests. Missing captions can justify `auto`/STT if that scope and provider cost are authorized. Preserve partial outputs; use the supported refresh path for a specific transcript rather than deleting user data by default.
