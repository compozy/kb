---
schema_version: 1
review_kind: implementation
round: 1
verdict: SHIP
reviewer_runtime: claude
reviewer_model: opus
generated_at: 2026-06-27T01:35:57Z
---

# Summary

The change correctly replaces the global English caption bias with original-language (`orig`) defaulting, wires caption-language preferences through config/env/CLI/channel ingestion, gates translated automatic captions behind `allow_translated_captions`, and classifies yt-dlp HTTP 429 as `rate_limited` before generic network-block classification. Every spec acceptance criterion is implemented with matching tests; build, `golangci-lint`, and the targeted suites all pass, so there are no blockers — only a couple of latent edge/maintenance follow-ups.

# Blockers

None.

# Risks

## R-001 — `orig` default fails for manual-only videos with no language metadata

- File: internal/mediadl/ytdlp.go
- Line: 767
- Issue: `detectYTDLPOriginalLanguage` derives the original language only from `info.language` or from automatic caption keys ending in `-orig`. When yt-dlp reports an empty `info.language` AND the video exposes no automatic captions (only a manual subtitle track), detection returns `""`, so `newYTDLPCaptionCandidate` marks the manual track `Original=false`, `matchingYTDLPCaptionCandidates` finds no `orig` candidate, and extraction returns `transcript_unavailable` ("no original-language caption available") for a video that demonstrably has a usable caption. This matches the spec's stated detection algorithm (info.language or `-orig` keys), and the case is rare because YouTube usually populates `info.language` and/or emits a `*-orig` ASR track — but the failure is silent-looking and counter-intuitive when it hits.
- Suggested fix: As a follow-up, when no original language is detectable and exactly one usable manual caption track exists, treat that single manual track as the original (or document that users must pass `--sub-langs <lang>` for such videos). No code change required to ship.

## R-002 — Production-unused HTTP-status classifier kept alive only by a test

- File: internal/mediadl/mediadl.go
- Line: 398
- Issue: `unexpectedStatusCode`, `isStatusCodeRateLimited`, and `isStatusCodeNetworkBlocked` are never called by production code — the real rate-limit/network-block classification (`IsRateLimited`/`IsNetworkBlocked`) routes through the string-marker helpers `isYTDLPRateLimitedError`/`isYTDLPNetworkBlockedError` in `ytdlp.go`. The only caller of the status-code family is `TestIsStatusCodeNetworkBlocked`, so the test freezes a parallel code path that never runs (the `"unexpected status code: NNN"` format is produced by the Firecrawl client, not yt-dlp). This is a migration leftover from the deleted `internal/youtube/ytdlp.go` and a test-shape smell (CLAUDE.md: tests should not freeze unused implementation details). `golangci-lint`'s `unused` does not flag it because the test references it.
- Suggested fix: Either delete the status-code helpers plus their test, or wire `isStatusCode*` into the `wrap*FetchError` paths if classifying `unexpected status code` diagnostics is actually intended.

# Nits

## N-001 — Dead constant `transcriptUnavailableMessage`

- File: internal/mediadl/mediadl.go
- Line: 19
- Issue: `transcriptUnavailableMessage = "captions unavailable"` is unused across the package (the unavailable message is now produced by `ytDLPCaptionUnavailableMessage`); it's a leftover from the old youtube package and is not caught by `unused` because it is a constant.
- Suggested fix: Delete the constant.

## N-002 — New `video_id` frontmatter field is undocumented

- File: internal/cli/ingest_youtube.go
- Line: 172
- Issue: `youtubeFrontmatter` now persists `video_id`, but README's enumeration of YouTube transcript frontmatter fields (engagement/publication/classification) does not mention it.
- Suggested fix: Add `video_id` to the README list of persisted YouTube metadata fields.

# Evidence

- Read in full: `.codex/plans/20260626221122-youtube-caption-language.md`, `.agents/skills/kb-yt-channel/scripts/ingest-channel.py`, `CLAUDE.md` (project + global), and the full review patch.
- Read final-state source in full: `internal/mediadl/mediadl.go`, `internal/mediadl/ytdlp.go`, `internal/youtube/youtube.go`, `internal/youtube/channel.go`, `internal/youtube/channel_test.go`, `internal/cli/ingest_channel.go`, plus the patch hunks for `internal/config/config.go`, `internal/config/env.go`, `internal/config/config_test.go`, `internal/cli/ingest_youtube.go`, `internal/cli/ingest_test.go`, `internal/mediadl/mediadl_test.go`, `internal/mediadl/ytdlp_test.go`, README.md, config.example.toml, SKILL.md.
- Verified helper/wiring existence: `optionalString`, `existingYouTubeVideoIDs`, `runIngestTopicInfo`, `newInstagramTranscriptExtractor`, and channel/instagram command registration in `internal/cli/ingest.go`.
- `go build ./...` — success.
- `go test ./internal/config ./internal/mediadl ./internal/youtube ./internal/cli` — 239 passed.
- `golangci-lint run` over config/youtube/cli/mediadl/instagram — 0 issues (`staticcheck` + `unused` enabled per `.golangci.yml`).
- Spec cross-check: confirmed default `["orig"]`, env `YOUTUBE_CAPTION_LANGUAGES` with CLI>env>TOML>default precedence, `--sub-langs`/`--lang` on both `ingest youtube` and `ingest channel` with disagreement error, `AllowTranslatedCaptions` plumbed through `ExtractOptions`/`BulkOptions`, `ErrorKindRateLimited` added + re-exported, manual-over-automatic and `<lang>-orig`-over-bare ranking, translated gating (exact-token + allow flag), English used only as a tie-break among originals, 429 classified before network-block in metadata/caption/audio wrappers with the requested language in the caption message, and bulk retry treating `rate_limited` as retryable.
- Limitation: did not execute against a live yt-dlp/YouTube; correctness verified through the fake-yt-dlp harness tests and source reading. `make verify` (full gate including integration) was not run end-to-end; build+lint+targeted unit tests were.

# Deferred Or Follow-Up

- R-001: improve `orig` fallback for manual-only/no-language-metadata videos, or document the `--sub-langs` workaround.
- R-002: remove or wire-in the unused HTTP-status classifier and reconcile its test.
- N-001 / N-002: delete the dead constant; document the `video_id` frontmatter field.
