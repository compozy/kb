---
status: completed
title: Implement Firecrawl REST client
type: backend
complexity: medium
dependencies:
  - task_03
---

# Task 10: Implement Firecrawl REST client

## Overview

Create the `internal/firecrawl/` package that provides a Go client for the Firecrawl REST API (`POST /v2/scrape`). This client powers the `ingest url` command, converting web URLs to clean Markdown via firecrawl's scraping infrastructure.

<critical>
- ALWAYS READ the PRD and TechSpec before starting
- REFERENCE TECHSPEC for implementation details — do not duplicate here
- FOCUS ON "WHAT" — describe what needs to be accomplished, not how
- MINIMIZE CODE — show code only to illustrate current structure or problem areas
- TESTS REQUIRED — every task MUST include tests in deliverables
</critical>

<requirements>
- MUST call `POST /v2/scrape` with `{"url": "...", "formats": ["markdown"]}` and Bearer token auth
- MUST read API key and URL from `FirecrawlConfig` (task_03)
- MUST return scraped Markdown content, extracted title, and source URL metadata
- MUST implement exponential backoff retry for transient failures (429, 5xx) with max 3 attempts
- MUST return actionable error messages when API key is missing, request fails, or URL is unreachable
- MUST support custom API URL for self-hosted firecrawl instances
- MUST use `context.Context` for cancellation support
</requirements>

## Subtasks

- [x] 10.1 Create `internal/firecrawl/` package with Client struct and Scrape method
- [x] 10.2 Implement HTTP request construction with Bearer auth and JSON body
- [x] 10.3 Implement response parsing for markdown content and metadata
- [x] 10.4 Implement exponential backoff retry logic for transient errors
- [x] 10.5 Write unit tests using `httptest.NewServer` mock

## Implementation Details

Create `internal/firecrawl/client.go` and `internal/firecrawl/client_test.go`. Use Go's `net/http` standard library — no third-party HTTP client needed.

Reference TechSpec "Integration Points" Firecrawl section for the API contract and error handling strategy.

### Relevant Files

- `internal/config/config.go` (task_03) — FirecrawlConfig struct with APIKey and APIURL
- `internal/qmd/client.go` — existing pattern for external tool integration (subprocess-based, but similar error handling patterns)

### Dependent Files

- `internal/cli/` (task_16) — `ingest url` command wires to this client

### Related ADRs

- [ADR-004: Firecrawl REST API for URL Scraping](../adrs/adr-004.md) — REST API chosen over CLI shell-out

## Deliverables

- `internal/firecrawl/client.go` — Firecrawl REST client
- `internal/firecrawl/client_test.go` — tests with HTTP mock
- Unit tests with 80%+ coverage **(REQUIRED)**

## Tests

- Unit tests:
  - [x] Scrape sends correct POST request with Bearer auth and JSON body
  - [x] Scrape returns markdown content and title from successful response
  - [x] Scrape retries on 429 (rate limit) up to 3 times with backoff
  - [x] Scrape retries on 500/502/503 up to 3 times
  - [x] Scrape does NOT retry on 400/401/404 (client errors)
  - [x] Scrape returns descriptive error when API key is empty
  - [x] Scrape returns descriptive error when URL is unreachable (non-2xx after retries)
  - [x] Scrape respects context cancellation during retry
  - [x] Scrape uses custom API URL when configured
- Test coverage target: >=80%
- All tests must pass

## Success Criteria

- All tests passing
- Test coverage >=80%
- Client handles all error scenarios gracefully with actionable messages
- `make lint` reports zero findings
