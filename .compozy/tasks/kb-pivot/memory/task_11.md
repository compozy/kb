# Task Memory: task_11.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot

- Build `internal/youtube` with transcript-first extraction, metadata mapping, and OpenRouter STT fallback behind an explicit/configured gate.

## Important Decisions

- Pinned `github.com/kkdai/youtube/v2` to `v2.10.5` instead of `v2.10.6` because `v2.10.6` forces the repo `go` directive to `1.26`; `v2.10.5` still exposes the needed transcript APIs while keeping the module compatible with `go 1.24.0`.
- Added a strict YouTube URL parser that accepts supported YouTube hosts and URL shapes (`watch`, `youtu.be`, `shorts`, `embed`) instead of trusting generic ID extraction from arbitrary hosts.
- STT fallback returns normalized Markdown with a `00:00` header because the OpenRouter chat-completions audio flow in scope here does not provide deterministic timestamp segments like YouTube captions do.

## Learnings

- `kkdai/youtube` transcript extraction requires a concrete caption language code; passing an empty language is not a safe default, so caption-track ordering and language selection need to happen in the package before calling `GetTranscriptCtx`.
- `kkdai/youtube` audio-only format sorting already prefers `audio/mp4` ahead of `opus/webm`, which makes it a good fit for OpenRouter audio input without adding transcoding.

## Files / Surfaces

- `internal/youtube/youtube.go`
- `internal/youtube/openrouter.go`
- `internal/youtube/youtube_test.go`
- `internal/youtube/openrouter_test.go`
- `go.mod`
- `go.sum`

## Errors / Corrections

- Initial `go get github.com/kkdai/youtube/v2@v2.10.6` upgraded the module `go` directive to `1.26`; corrected by downgrading to `v2.10.5` and resetting the repo directive with `go get ... go@1.24.0`.

## Ready for Next Run

- `go test ./internal/youtube -cover` passes with `80.0%` coverage and `make verify` passes cleanly.
