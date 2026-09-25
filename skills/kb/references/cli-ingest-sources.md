# Ingest Sources Reference

Source-specific behaviour and configuration for `kb ingest url|file|youtube|channel|instagram|bookmarks` and `kb migrate`. Every one of these runs the ingest gates and then classifies and links the new sources (see `decision-workflow.md`); `kb ingest codebase` is separate (`cli-ingest-codebase.md`).

## Shared flags

```bash
kb ingest <kind> <input> --topic <topic-id>
  --batch <name>             # value of ingest_batch (default <command>-<YYYY-MM-DD>-<run id>)
  --budget <usd>             # spend ceiling for this run; 0 = cached answers only
  --decisions shadow|apply   # gate and body-link mode for this run
  --force                    # skip dedupe and gates (still classified and linked)
```

Use path-relative topic ids for nested topics (`--topic harness/goclaw`). stdout carries the JSON result (`filePath`, `triage`, `triage_reason`, `quality`, `duplicate_of`, `shadow`, `supersedes`, `refetched`, `review_item`); stderr carries the run summary.

## URLs (`kb ingest url`)

```bash
kb ingest url <url> --topic <topic-id>                 # one page via Firecrawl
kb ingest url <url1> <url2> --topic <topic-id>         # several URLs = bulk run
kb ingest url --from <file> --topic <topic-id>         # one URL per line, optional label after whitespace, # comments
kb ingest url --rescue <id> --topic <topic-id>         # fetch an item the pre-fetch gate skipped
```

Needs `FIRECRAWL_API_KEY` (or `[firecrawl].api_key`). Bulk runs judge each URL against the contract before fetching. A thin or broken capture is refetched once with Firecrawl's cache bypassed (`[firecrawl].refetch_wait_ms`, `refetch_only_main_content`).

## Local files (`kb ingest file`)

Converts PDF, DOCX, XLSX, PPTX, EPUB, HTML, CSV, JSON, XML, text and images (OCR) to markdown under `raw/articles/`. `ingest_query` records the file name.

## YouTube videos (`kb ingest youtube`)

Writes `raw/youtube/<slug>.md` with `yt-dlp` metadata (`view_count`, `like_count`, `comment_count`, `upload_date`, `duration`, `channel`, `channel_id`, `categories`, `youtube_tags`, `language`, `chapter_count`, ...; see `frontmatter-schemas.md`). Use `youtube_tags` for video keywords: `tags` is KB taxonomy.

```toml
[youtube]
yt_dlp_path = "yt-dlp"
proxy = ""
cookies_file = ""
user_agent = ""
retry_attempts = 3
retry_backoff = "1s"
transcription = "captions" # captions | auto | stt
```

`--transcribe auto` uses manual captions when present and STT otherwise; `--transcribe stt` forces STT (OpenAI by default: `OPENAI_API_KEY`; `[stt]` for provider, chunking and `ffmpeg`). `YOUTUBE_YT_DLP_PATH`, `YOUTUBE_PROXY`, `YOUTUBE_COOKIES_FILE` and `YOUTUBE_USER_AGENT` override the TOML values. Install or update `yt-dlp` before treating a caption failure as a network problem; `network_blocked` means YouTube refused the request (configure proxy/cookies or use a trusted network).

## YouTube channels and playlists (`kb ingest channel`)

Bulk-extracts every upload into `raw/youtube/`, one document per video. Flags: `--limit <n>` (newest n; `0` = all), `--all`, `--concurrency <n>`, `--throttle <dur>`, `--dry-run` (list without ingesting), `--transcribe`, `--rescue <id>`. Pool size, throttle, backoff and retries default from `[youtube].bulk_concurrency`, `bulk_throttle`, `bulk_backoff_max`, `bulk_retries`. Already-ingested videos are skipped; single-video URLs are rejected (use `kb ingest youtube`). Titles are judged against the contract before any transcript or STT cost.

## Instagram (`kb ingest instagram`)

Reels, video posts (`/p/`) and IGTV (`/tv/`) into `raw/instagram/`: the body is the caption under `## Caption` followed by the transcript under `## Transcript`; a reel without speech degrades to caption-only (`transcript_source: none`). `[instagram]` mirrors the `[youtube]` yt-dlp keys with **separate cookies** (Instagram needs its own session) and defaults to `transcription = "auto"`. Whole-profile ingestion is not supported.

## Bookmark clusters (`kb ingest bookmarks`)

Ingests one bookmark-cluster markdown file as one source under `raw/bookmarks/` (URLs extracted into `source_urls`). The file is gated as a single document (no pre-fetch stage, no `--rescue`).

## Layout migrations

```bash
kb migrate transcripts --topic <topic-id>   # move legacy raw/transcripts/*.md into raw/youtube/
```
