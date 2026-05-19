# Superseded: yt-dlp Primary Backend For YouTube Captions

This plan has been superseded by `.codex/plans/2026-05-18-advanced-stt.md` after the accepted greenfield STT decision.

Current YouTube/STT contract:
- `yt-dlp` is mandatory for YouTube metadata, captions, and audio extraction.
- There is no alternate YouTube extractor or legacy audio path.
- `kb ingest youtube` exposes `--transcribe captions|auto|stt`; the previous shortcut flag is removed.
- OpenAI `/v1/audio/transcriptions` is the default STT provider; OpenRouter is an optional provider.
- Long STT audio is chunked with `ffmpeg`.

Keep this file only as a historical pointer; do not use it as implementation guidance.
