# Fix YouTube Caption Language Selection

## Summary

Fix the root cause in the active shared yt-dlp implementation under `internal/mediadl`: caption selection currently flattens original and translated tracks, then globally biases English. Preserve the current in-progress unification where `internal/youtube` delegates to `internal/mediadl`; do not resurrect the deleted duplicate `internal/youtube/ytdlp.go`.

The new default behavior is native/original captions via `["orig"]`. Explicit language lists are authoritative in order, but machine-translated tracks are selected only when the requested language is exact and `[youtube].allow_translated_captions = true`.

## Public Interfaces

- Add `[youtube]` config:
  - `caption_languages = ["orig"]`
  - `allow_translated_captions = false`
- Add env override:
  - `YOUTUBE_CAPTION_LANGUAGES=orig,pt,en`
  - Precedence: CLI flag > env > TOML > default `["orig"]`.
- Add CLI aliases on both `kb ingest youtube` and `kb ingest channel`:
  - `--sub-langs <comma-list>`
  - `--lang <comma-list>`
  - Both accept `orig`; if both aliases are supplied with different normalized values, return a clear error.
- Extend shared types:
  - `mediadl.ExtractOptions` / re-exported `youtube.ExtractOptions`: add `AllowTranslatedCaptions bool`.
  - `youtube.BulkOptions`: add `AllowTranslatedCaptions bool` and continue passing `PreferredLanguages`.
  - Add `mediadl.ErrorKindRateLimited = "rate_limited"` and re-export `youtube.ErrorKindRateLimited`.

## Implementation Changes

- Config and CLI wiring:
  - Add `CaptionLanguages []string` and `AllowTranslatedCaptions bool` to `config.YouTubeConfig`.
  - Apply defaults before env overrides; parse `YOUTUBE_CAPTION_LANGUAGES` as a trimmed comma list.
  - Validate caption language lists contain at least one non-empty token after defaults/overrides.
  - Pass resolved caption languages and `AllowTranslatedCaptions` into single-video and channel extraction options.
- Caption selection:
  - Extend `ytDLPCaptionCandidate` with `Original`, `Translated`, and base-language metadata while preserving the exact yt-dlp key in `Language`.
  - Detect the original language from `info.language`; if absent, infer from usable automatic caption keys ending in `-orig`.
  - Treat manual captions matching the selected token as preferred over automatic captions.
  - For `orig`, select original manual captions first, then automatic `<orig>-orig`, then bare/region-compatible original automatic tracks.
  - Remove the global English bias. English may only break ties among already-original candidates at the same priority; otherwise use deterministic exact-match then lexical ordering.
  - Classify automatic captions whose base language differs from the detected original as translated.
  - Reject translated candidates unless the requested token exactly names that candidate language and `AllowTranslatedCaptions` is true.
  - When no requested/original track is eligible, return `transcript_unavailable` with a message naming the requested languages and, when relevant, explaining that translated captions are disabled.
- 429 handling:
  - Add a rate-limit classifier for yt-dlp diagnostics containing HTTP 429 / Too Many Requests / unexpected status code 429.
  - Check rate-limit before generic network-block classification in metadata, caption, and audio wrappers.
  - Include requested caption language in caption rate-limit messages, e.g. `caption language "en" was rate limited through yt-dlp`.
  - Keep bulk retry/backoff behavior for rate limits by making channel retry logic treat `rate_limited` as retryable, without labeling it as `network_blocked`.

## Docs And Skill

- Update `config.example.toml` and `README.md`:
  - Document default native-caption behavior.
  - Document `caption_languages`, `YOUTUBE_CAPTION_LANGUAGES`, `allow_translated_captions`, and `--sub-langs` / `--lang`.
  - Explain that translated YouTube captions are off by default because they use YouTube's translation endpoint.
- Update `.agents/skills/kb-yt-channel/SKILL.md` and `scripts/ingest-channel.py`:
  - Add optional `--sub-langs` passthrough to `kb ingest channel`.
  - Reflect caption language selection in the generated command line/topic ingest metadata where the script already records ingest options.
  - Note that translated caption languages require `[youtube].allow_translated_captions = true`.

## Test Plan

- Add failing mediadl selector tests before implementation:
  - Portuguese metadata with `info.language="pt"` and automatic `{pt-orig, pt, en, es, fr}` selects `pt-orig` by default, never `en`.
  - Manual original `pt` beats automatic `pt-orig`.
  - English metadata with `en-orig` selects original English.
  - `PreferredLanguages: ["orig"]` resolves to the detected original language.
  - `PreferredLanguages: ["pt"]` selects the Portuguese original track.
  - `PreferredLanguages: ["es"]` fails while `AllowTranslatedCaptions=false`.
  - `PreferredLanguages: ["es"]` selects `es` when `AllowTranslatedCaptions=true`.
  - No original track and no eligible preferred track returns a clear `transcript_unavailable` error.
  - Caption fetch 429 returns `ErrorKindRateLimited` and includes the selected language.
- Update CLI/config/channel tests:
  - Defaults expose `CaptionLanguages == []string{"orig"}` and `AllowTranslatedCaptions == false`.
  - TOML and `YOUTUBE_CAPTION_LANGUAGES` are honored with CLI override precedence.
  - `kb ingest youtube --sub-langs pt` and `--lang orig` pass the expected `PreferredLanguages`.
  - `kb ingest channel --sub-langs pt` passes the same preferences through `BulkOptions`.
  - Rate-limited channel extraction still retries/backoffs, but the structured error kind is `rate_limited`.
- Verification commands:
  - `rtk go test ./internal/config ./internal/mediadl ./internal/youtube ./internal/cli`
  - `rtk make verify`

## Assumptions

- The current untracked `internal/mediadl` extraction is intentional working state; implementation should update it and keep `internal/youtube` as the shim.
- `--sub-langs` / `--lang` order controls selection. The default list is `["orig"]`; users who want fallback behavior should configure lists like `["orig", "en"]`.
- No retry, sleep, proxy, or masking workaround is part of this fix. The root fix is correct caption selection plus distinct rate-limit classification.
