# Add yt-dlp Primary Backend For YouTube Captions

## Summary

- Root cause: `kb ingest youtube` still routes metadata and captions through `github.com/kkdai/youtube/v2`, so proxy/cookie/user-agent changes wrap the same protocol YouTube is rejecting.
- Fix: add a `yt-dlp` shell-out backend as the primary caption and metadata extractor, keep `kkdai` as a legacy fallback, and preserve OpenRouter STT behavior when captions are genuinely unavailable.
- This is a root-cause fix: proxy/cookies remain supported for auth/rate-limit cases, but normal public captions should use `yt-dlp`'s maintained YouTube extractor path.

## Key Changes

- Add `youtube.yt_dlp_path` TOML config and `YOUTUBE_YT_DLP_PATH` env override, defaulting to `yt-dlp`.
- Add an `internal/youtube` `yt-dlp` backend that resolves the executable with `exec.LookPath`, runs via `exec.CommandContext`, always passes `--ignore-config`, and forwards configured proxy/cookies/user-agent.
- Use a two-step flow: first `--dump-single-json --skip-download` for metadata and subtitle inventory, then download exactly one selected subtitle with `--write-subs` or `--write-auto-subs`, exact `--sub-langs`, `--sub-format json3/vtt/best`, a temp output directory, and `--skip-download`.
- Select captions deterministically: preferred languages first, manual before automatic, English before other languages when no preference exists, then lexical order.
- Parse `json3` captions into the existing timestamped Markdown format and add a minimal WebVTT fallback.
- Try `yt-dlp` first. Fall back to `kkdai` only when `yt-dlp` is unavailable or fails before proving captions are absent. If `yt-dlp` proves captions are absent, keep the current STT fallback eligibility.
- Update README and `skills/kb` docs so users install/update `yt-dlp` for reliable YouTube ingest and do not treat every `network_blocked` as only a proxy/cookie problem.

## Test Plan

- Unit tests in `internal/youtube`:
  - fake `yt-dlp` returns metadata JSON and writes `.json3`; assert metadata, selected language, source, and Markdown.
  - fake `yt-dlp` writes `.vtt`; assert VTT fallback parsing.
  - assert the backend runs a metadata command first and one exact-language caption command second.
  - assert argv includes `--ignore-config`, `--no-playlist`, exact `--sub-langs`, and the correct manual/automatic subtitle flag.
  - assert proxy, cookies file, and user-agent config are forwarded.
  - assert `yt-dlp` missing falls back to `kkdai`.
  - assert `yt-dlp` caption success prevents any `kkdai` call.
  - assert no-captions maps to `transcript_unavailable` and preserves STT fallback.
  - assert both-backend failure returns actionable diagnostics.
- Config tests: TOML and env loading for `youtube.yt_dlp_path`.
- CLI tests: `kb ingest youtube` passes complete config into the extractor while preserving JSON output and `raw/youtube/` ingest behavior.
- Verification: `rtk go test ./internal/youtube ./internal/config ./internal/cli`, `rtk make verify`, and a live `kb ingest youtube <failing-url> --topic go-best-practices` acceptance run when feasible.

## Assumptions And Defaults

- `yt-dlp` is an optional primary runtime dependency; if unavailable, `kb` still attempts the legacy `kkdai` fallback.
- Keep `kkdai/youtube` pinned at the current module-compatible version because `v2.10.6` declares `go 1.26` while this repo contract remains `go 1.24.0`.
- Do not add a user-facing `--lang` flag in this change.
- Do not rewrite OpenRouter STT audio download in this change; if captionless-video STT later proves blocked by the same legacy audio path, handle it as a separate root-cause fix.
