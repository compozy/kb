---
name: kb-yt-channel
description: Creates and maintains Knowledge Base topics from YouTube channels by resolving recent or full uploads, scaffolding yt-channels topics, ingesting transcripts through kb ingest youtube, and validating/indexing the result. Use when turning a YouTube channel into a Karpathy KB topic. Do not use for single-video ingestion, general video summaries, or non-YouTube sources.
---

# KB YouTube Channel

## Required Inputs

- Channel URL, such as https://www.youtube.com/@aiDotEngineer.
- Topic slug, such as ai-dot-engineer.
- Topic title and domain.
- Video selection: either --limit N for the newest uploads or --all for the full uploads list.
- Transcript policy: --transcribe captions, auto, or stt.
- Caption language policy: defaults to `orig`; optionally pass --sub-langs with a comma-separated list such as `orig,pt`.
- Indexing policy: lexical QMD indexing by default; pass --embed when vector embeddings are explicitly needed.

## Procedure

**Step 1: Confirm Prerequisites**
1. Run kb version from the vault root.
2. Run yt-dlp --version from the vault root.
3. Run qmd --version when final indexing is required.
4. If --transcribe stt is selected, confirm ffmpeg -version and the configured STT provider credentials.

**Step 2: Run Channel Ingest**
1. Read references/channel-topic-contract.md when topic metadata or output validation is unclear.
2. Execute python3 scripts/ingest-channel.py from this skill directory, passing:
   - --vault /path/to/vault
   - --channel-url <youtube-channel-url>
   - --topic-slug <slug>
   - --title <title>
   - --domain <domain>
	   - either --limit <n> or --all
	   - --transcribe captions|auto|stt
	   - optionally --sub-langs orig,<lang>
	   - optionally --embed for vector embedding after collection sync
3. Treat the script's JSON stdout as the run summary. Treat stderr as progress and diagnostics.
4. The script delegates video enumeration, resume, and per-video transcript ingestion to `kb ingest channel`, which applies low-concurrency throttling and adaptive backoff to avoid YouTube rate limiting and skips videos already present in raw/youtube. Caption selection defaults to the video's original language (`orig`); translated YouTube automatic captions require `[youtube].allow_translated_captions = true`. Tune `[youtube].bulk_concurrency`, `[youtube].bulk_throttle`, `[youtube].bulk_backoff_max`, and `[youtube].bulk_retries` (or the `--concurrency`/`--throttle` flags) when a large channel still trips rate limits.

Example:

    python3 .agents/skills/kb-yt-channel/scripts/ingest-channel.py --vault . --channel-url https://www.youtube.com/@aiDotEngineer --topic-slug ai-dot-engineer --title "AI Engineer Channel" --domain youtube-channel --limit 10 --transcribe captions

**Step 3: Verify The Topic**
1. Confirm kb topic info yt-channels/<slug> returns the expected topic path and source count.
2. Confirm <topic>/outputs/reports/ contains the channel ingest report for the run.
3. Confirm kb search "<channel topic>" --topic yt-channels/<slug> --collection <slug> --lex --format json returns raw transcript or index content after indexing.
4. Leave wiki article compilation to the normal kb compile workflow unless the user explicitly asks for article synthesis.

## Error Handling

- If channel resolution fails, read references/troubleshooting.md and retry after updating yt-dlp or YouTube proxy/cookie settings.
- If captions are unavailable and --transcribe captions was selected, report the failed video and rerun only when the user chooses --transcribe auto or stt.
- If STT fails, read references/troubleshooting.md to check ffmpeg, provider credentials, model, chunk size, and network settings.
- If a partial topic exists after failure, rerun the same command. The script reuses the existing topic and skips videos whose source URLs are already present in raw/youtube.
