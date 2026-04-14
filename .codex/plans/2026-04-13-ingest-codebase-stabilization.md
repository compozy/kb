# Plan: stabilize `kb ingest codebase`

## Summary

- Preserve the already-fixed Go scan/parse path and topic-skeleton compatibility in the current `HEAD`.
- Fix the remaining root causes:
  - empty codebase ingests still succeed and generate misleading starter wiki output
  - `CLAUDE.md` is still rewritten wholesale
  - manual `wiki/index/*` pages are still overwritten
  - generated codebase wiki still mixes with curated `wiki/concepts/*`
- Final ownership model:
  - `raw/codebase/` and `wiki/codebase/` are managed and fully regenerable
  - `wiki/concepts/*`, `wiki/index/*`, and manual `CLAUDE.md` content are preserved
  - empty ingests fail before any writes

## Public Interface Changes

- `kb ingest codebase` gains `--dry-run`.
- `kb generate` gains `--dry-run`.
- `models.GenerateOptions` gains `DryRun bool`.
- `models.GenerationSummary` gains:
  - `dryRun`
  - `detectedLanguages`
  - `selectedAdapters`
- `kb ingest codebase --help` and `kb generate --help` list supported languages: `go`, `ts`, `tsx`, `js`, `jsx`.
- Generated codebase wiki moves to:
  - `wiki/codebase/concepts/*.md`
  - `wiki/codebase/index/Codebase Dashboard.md`
  - `wiki/codebase/index/Codebase Concept Index.md`
  - `wiki/codebase/index/Codebase Source Index.md`
- `CLAUDE.md` uses a managed block delimited by:
  - `<!-- kb:codebase:start -->`
  - `<!-- kb:codebase:end -->`

## Implementation Changes

- Generator:
  - fail before render/write when `filesScanned == 0`
  - fail before render/write when scan succeeded but `filesParsed == 0`
  - in `--dry-run`, run through renderable planning but skip all writes
- Vault writer:
  - reset only `raw/codebase/` and `wiki/codebase/`
  - move generated articles/indexes out of `wiki/concepts/` and `wiki/index/`
  - drop the `Kodebase - ` title/file prefix
  - keep generated index names prefixed with `Codebase` to avoid collisions with manual indexes
  - update only the managed block in `CLAUDE.md`; append the block if absent
  - keep `log.md` append-only, but only on successful non-dry-run ingests
- Topic/lint compatibility:
  - extend the current skeleton with `wiki/codebase/concepts/` and `wiki/codebase/index/`
  - keep read-time compatibility for older topics
  - lint `wiki/codebase/concepts/*` as compiled wiki and `wiki/codebase/index/*` as index pages
  - keep `topic info` article counts limited to `wiki/concepts/*`

## Test Plan

- Unit:
  - path helpers for `wiki/codebase/*`
  - empty corpus errors and dry-run summaries
  - writer preserves manual `wiki/index/*` and manual concept pages
  - managed-block updates in `CLAUDE.md`
  - lint validates the new generated wiki paths
  - topic counts exclude `wiki/codebase/*`
- Integration:
  - Go fixture still generates raw codebase snapshots and generated codebase wiki
  - empty ingest fails with no writes
  - manual `CLAUDE.md` and manual `wiki/index/*` survive ingest unchanged outside managed areas
  - re-ingest is idempotent over managed areas only
- Verification:
  - focused package tests
  - `make verify`

## Assumptions

- No `kb adapters list` command in this change.
- No new `topic.yaml` or `.kb/kodebase-state.json`.
- No diff-style re-ingest delta reporting in the JSON summary for this change.
