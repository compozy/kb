---
status: pending
title: Implement YouTube transcript extractor
type: backend
complexity: high
dependencies:
  - task_03
---

# Task 11: Implement YouTube transcript extractor

## Overview

Create the `internal/youtube/` package that extracts video transcripts from YouTube URLs. The primary path uses `kkdai/youtube` to fetch auto-generated captions programmatically. When captions are unavailable, a fallback path downloads the audio track and transcribes it via the OpenRouter API (chat completion with audio input).

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST extract video metadata (title, channel, duration, publish date) from YouTube URLs
- MUST extract auto-generated or manual captions/transcripts as the primary method using `kkdai/youtube`
- MUST format transcript as timestamped Markdown (timestamp headers with text paragraphs)
- MUST implement OpenRouter STT fallback: download audio → base64 encode → send to chat completion with `input_audio` → receive transcription
- MUST only invoke STT fallback when transcript extraction fails AND `--stt` flag is set or OpenRouter API key is configured
- MUST add `github.com/kkdai/youtube/v2` dependency
- MUST return structured error when video is unavailable, private, or age-restricted
</requirements>

## Subtasks

- [ ] 11.1 Add kkdai/youtube dependency
- [ ] 11.2 Implement YouTube video metadata extraction (title, channel, duration)
- [ ] 11.3 Implement caption/transcript extraction and formatting as timestamped Markdown
- [ ] 11.4 Implement OpenRouter STT fallback (audio download → transcription)
- [ ] 11.5 Implement fallback routing logic (try transcript first, STT if configured and needed)
- [ ] 11.6 Write unit tests with mocked YouTube and OpenRouter responses

## Implementation Details

Create `internal/youtube/youtube.go`, `internal/youtube/openrouter.go`, and corresponding test files. The YouTube client from `kkdai/youtube` provides video info and transcript access. The OpenRouter client sends audio as base64 in the `input_audio` content block of a chat completion request.

Reference TechSpec "Integration Points" sections for OpenRouter API and YouTube.

### Relevant Files

- `internal/config/config.go` (task_03) — OpenRouterConfig with APIKey, APIURL, STTModel

### Dependent Files

- `internal/cli/` (task_16) — `ingest youtube` command wires to this package

### Related ADRs

- [ADR-003: Native Go Document Conversion with Converter Registry](../adrs/adr-003.md) — YouTube transcript as a conversion source

## Deliverables

- `internal/youtube/youtube.go` — transcript extraction + metadata
- `internal/youtube/openrouter.go` — OpenRouter STT client
- Test files for both
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests

- Unit tests:
  - [ ] Extracts video metadata (title, channel, duration) from mock video info
  - [ ] Extracts transcript text from mock caption data
  - [ ] Formats transcript with timestamps as Markdown headers
  - [ ] Falls back to STT when transcript extraction returns no captions
  - [ ] Does NOT invoke STT when API key is not configured and --stt is not set
  - [ ] OpenRouter STT sends correct chat completion request with base64 audio
  - [ ] OpenRouter STT parses transcription from response
  - [ ] Returns error for private/unavailable video
  - [ ] Returns error for invalid YouTube URL format
  - [ ] Handles YouTube URL variations (youtube.com/watch?v=, youtu.be/, youtube.com/shorts/)
- Test coverage target: >=80%
- All tests must pass

## Success Criteria

- All tests passing
- Test coverage >=80%
- Transcript output is readable and well-formatted in Markdown
- STT fallback works end-to-end with mocked API
- `make lint` reports zero findings
