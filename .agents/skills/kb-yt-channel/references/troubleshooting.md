# Troubleshooting YouTube Channel Ingest

## Channel Listing Fails

- Run yt-dlp --version and update yt-dlp if it is stale.
- Retry the channel videos URL directly, for example https://www.youtube.com/@aiDotEngineer/videos.
- Configure YOUTUBE_PROXY, YOUTUBE_COOKIES_FILE, or YOUTUBE_USER_AGENT when YouTube blocks the local network.

## Captions Fail

- Confirm the video has public captions or automatic captions.
- Keep --transcribe captions when the goal is zero-cost ingestion and report missing captions as a real blocker.
- Use --transcribe auto to try captions first and fall back to STT.
- Use --transcribe stt when captions are known to be unusable.

## STT Fails

- Confirm ffmpeg is installed and on PATH.
- Confirm OPENAI_API_KEY or OPENROUTER_API_KEY is configured for the selected provider.
- Check kb.toml [stt] settings for provider, model, language, audio_format, chunk_duration, max_chunk_bytes, concurrency, and ffmpeg_path.
- Reduce STT concurrency or chunk size if provider requests time out or exceed upload limits.

## Partial Runs

The ingest script writes reports after failures. Rerun the same command after fixing the blocker. Existing raw YouTube source URLs are detected and skipped, so successful ingests do not need to be repeated.
