# Task Memory: task_10.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot

- Implement `internal/firecrawl` with a config-backed Firecrawl REST client for `POST /v2/scrape`, including retry/cancellation behavior and unit coverage above the task threshold.

## Important Decisions

- `NewClient` accepts `config.FirecrawlConfig` directly and falls back to `config.Default().Firecrawl.APIURL` when the configured API URL is empty.
- Retry scope stays tight to transient Firecrawl failures only: `429` and `5xx`, with 3 total attempts and exponential backoff from a `250ms` base delay.
- `ScrapeResult.SourceURL` prefers Firecrawl metadata in this order: `sourceURL`, `url`, then the input URL as a final fallback.

## Learnings

- The repository module path is still `github.com/compozy/kb`; new internal packages must import that path until the module rename task happens.
- `make lint` enforces explicit handling of `resp.Body.Close()` through `errcheck`, so deferred closes in new HTTP clients must use a closure or another checked path.
- Focused Firecrawl coverage started at `69.3%`; adding direct helper/error-path tests raised the package to `85.1%`.

## Files / Surfaces

- `internal/firecrawl/client.go`
- `internal/firecrawl/client_test.go`

## Errors / Corrections

- Corrected the initial import path from the repo name to the actual module path in `go.mod`.
- Replaced a bare `defer resp.Body.Close()` with a deferred closure to satisfy repo lint rules.
- Added focused tests for invalid URLs, API error payloads, transport failures, invalid JSON, helper functions, and retry timing to meet the coverage gate.

## Ready for Next Run

- Task 10 implementation is in place and repo verification is green; follow-on ingest work can wire `internal/firecrawl.Client` into the orchestrator and CLI without revisiting the REST contract.
