---
status: pending
title: Extend config with Firecrawl and OpenRouter
type: backend
complexity: low
dependencies: []
---

# Task 03: Extend config with Firecrawl and OpenRouter

## Overview

Extend the existing `internal/config/` package with new configuration sections for Firecrawl (URL scraping API) and OpenRouter (STT fallback). This enables the `ingest url` and `ingest youtube` commands to read API credentials from config or environment variables.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST add `FirecrawlConfig` struct with `api_key` and `api_url` TOML fields
- MUST add `OpenRouterConfig` struct with `api_key`, `api_url`, and `stt_model` TOML fields
- MUST support environment variable overrides: `FIRECRAWL_API_KEY`, `FIRECRAWL_API_URL`, `OPENROUTER_API_KEY`, `OPENROUTER_API_URL`
- MUST provide sensible defaults: Firecrawl API URL `https://api.firecrawl.dev`, OpenRouter API URL `https://openrouter.ai/api`, STT model `google/gemini-2.5-flash`
- MUST update `config.example.toml` with the new sections
</requirements>

## Subtasks

- [ ] 3.1 Add `FirecrawlConfig` and `OpenRouterConfig` structs to `Config`
- [ ] 3.2 Add environment variable loading for new API keys and URLs
- [ ] 3.3 Set defaults for API URLs and STT model
- [ ] 3.4 Update `config.example.toml` with documented new sections
- [ ] 3.5 Write unit tests for new config loading and env override behavior

## Implementation Details

Extend the existing `Config` struct in `internal/config/config.go` with two new embedded structs. Add new env constants and loading logic in `internal/config/env.go`.

Reference TechSpec "Data Models" Config model section for struct definitions.

### Relevant Files

- `internal/config/config.go` — Config struct, Load(), Validate() functions to extend
- `internal/config/env.go` — LoadSecretsFromEnv(), env constant definitions to extend
- `internal/config/config_test.go` — existing test patterns to follow
- `config.example.toml` — update with new sections

### Dependent Files

- `internal/firecrawl/` (task_10) — reads FirecrawlConfig for API key and URL
- `internal/youtube/` (task_11) — reads OpenRouterConfig for STT fallback

### Related ADRs

- [ADR-004: Firecrawl REST API for URL Scraping](../adrs/adr-004.md) — defines the config pattern for firecrawl integration

## Deliverables

- Modified `internal/config/config.go` with new config structs
- Modified `internal/config/env.go` with new env loading
- Updated `config.example.toml`
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests

- Unit tests:
  - [ ] Load TOML with `[firecrawl]` section populates FirecrawlConfig fields
  - [ ] Load TOML with `[openrouter]` section populates OpenRouterConfig fields
  - [ ] Missing `[firecrawl]` section uses default API URL
  - [ ] Missing `[openrouter]` section uses default API URL and STT model
  - [ ] FIRECRAWL_API_KEY env var overrides TOML `api_key`
  - [ ] OPENROUTER_API_KEY env var overrides TOML `api_key`
  - [ ] FIRECRAWL_API_URL env var overrides TOML `api_url`
- Test coverage target: >=80%
- All tests must pass

## Success Criteria

- All tests passing
- Test coverage >=80%
- `config.example.toml` documents all new keys
- `make lint` reports zero findings
