# Memory Ledger — kb-yt-channel skill rewrite

- Goal (incl. success criteria): Rewrite the `/kb-yt-channel` skill to adopt the new native `kb ingest channel` command introduced on branch `ytchan`. The skill must shrink: delegate all ingest mechanics to `kb ingest channel`, keeping only KB-side organization. Success = skill files updated, redundant script/shim/tests removed, new script syntactically valid, docs consistent with native flags.
- Constraints/Assumptions:
  - Native `kb ingest channel <url> --topic <slug>` handles: URL normalization, video enumeration, resume/dedup (existing video_id), bounded concurrency (`--concurrency`/`[youtube].bulk_concurrency`), throttle+jitter (`--throttle`/`bulk_throttle`), adaptive backoff (`bulk_backoff_max`), per-video retry (`bulk_retries`), native captions (`--sub-langs orig` / `caption_languages` / `allow_translated_captions`), transcription policy (captions|auto|stt), frontmatter, JSON summary (snake_case keys: channel_url, normalized_channel_url, selection, transcribe, resolved, videos[], ingested[], skipped[], failures[]).
  - Native requires the topic to ALREADY EXIST (resolveIngestTarget -> topic.Info), even for `--dry-run`. So script must scaffold first.
  - Native writes to `<topic>/raw/youtube/`. Topic id is category-qualified `yt-channels/<slug>`.
  - kb resolves vault from cwd (script runs kb with cwd=vault, no --vault flag), preserve that.
- Key decisions:
  - KEEP a thin Python orchestrator `ingest-channel.py` (topic scaffold under yt-channels/, topic.yaml, category docs, CLAUDE.md patch, wiki indexes, report, validation) — these are NOT done by native.
  - DELETE `scripts/ytdlp_native_shim.py` (native `--sub-langs orig` replaces it) and `scripts/test_skill.py` (tested classify_failure + shim, both removed).
  - DROP from script: resolve_videos, existing_video_ids, classify_failure, setup_shim_env, remove_existing_files, ingest_one_video, ThreadPoolExecutor, normalize_channel_url, ProxyExhaustedError, rate-limit flags (--sleep/--max-retries/--retry-backoff/--metadata-direct/--overwrite).
  - dry-run = scaffold topic skeleton + native `--dry-run` (lists videos, ingests nothing). Documented as creating empty topic.
  - Rate-limit/proxy/cookie config moves to kb config/env (`[youtube]`, YOUTUBE_PROXY, YOUTUBE_COOKIES_FILE) — not script flags.
- State: COMPLETE.
- Done: explored via subagents; rewrote ingest-channel.py (1027->747 lines, thin orchestrator over `kb ingest channel`); rewrote SKILL.md (95 lines) + troubleshooting.md (52 lines); updated report-template.md; deleted ytdlp_native_shim.py + test_skill.py. py_compile OK, --help OK, smoke tests pass (validate_inputs + kb command construction for --limit/--all/--dry-run). No stale refs (shim/--sleep/--overwrite/kind taxonomy). skills-lock.json left untouched (upstream sync artifact).
- Now: done.
- Next: (optional) push skill source to pedronauck/skills + re-sync to refresh skills-lock.json hash.
- Open questions: none.
- Working set: .agents/skills/kb-yt-channel/{SKILL.md,scripts/ingest-channel.py,references/troubleshooting.md,references/channel-topic-contract.md,assets/report-template.md}; internal/cli/ingest_channel.go (native ref).
