# Plan: stabilize skeleton compatibility, `ingest codebase`, and `kb index`

## Summary

- Fix 3 root causes already reproduced:
  - legacy topics disappear from `topic list/info` because read-time validation requires the newest scaffold
  - `ingest codebase` can leave a topic invalid when generation writes zero raw codebase notes
  - `kb index` fails on hosts without `sqlite-vec` even when lexical indexing is otherwise usable
- Use the approved product behavior:
  - compatibility on read, normalization on write
  - graceful `index` fallback when vector support is unavailable

## Implementation Changes

- Unify topic scaffold expectations between `internal/topic` and `internal/vault`.
- Define a compatible topic contract for read/discovery and keep the full current scaffold as the write-time invariant.
- Update `topic.List` and `topic.Info` to accept legacy topics that still have the topic root structure and required marker files.
- Normalize the full current scaffold idempotently in write paths so `topic.New`, generic ingest, and codebase ingest all preserve the latest directory layout.
- Fix `vault.WriteVault` subtree reset so `raw/codebase/files` and `raw/codebase/symbols` always exist after codebase writes, including zero-file codebase runs.
- Extend `qmd.Index` and the `kb index` JSON payload to expose explicit embed outcome:
  - `completed`
  - `skipped_unavailable`
  - `not_requested`
- Keep `--force-embed` strict even when default embedding falls back.

## Public Interface Changes

- `kb index` JSON output gains:
  - `embedStatus`
  - `embedWarning` when embeddings are skipped because vector support is unavailable
- Existing `embedResult` remains present only when embedding actually runs.

## Test Plan

- Unit tests:
  - legacy topic skeleton remains visible to `topic list/info`
  - write-time skeleton normalization is idempotent and preserves topic validity
  - `WriteVault` preserves `raw/codebase/files` and `raw/codebase/symbols` after empty codebase generation
  - `qmd.Index` skips embeddings only for the vector-unavailable case
  - `qmd.Index` still fails when `ForceEmbed=true`
  - `kb index` payload reports the new embed fields correctly
- Integration and QA:
  - run codebase ingest on a normal fixture and on an input that yields zero parsed files
  - run `topic list/info`, `ingest file`, `ingest codebase`, `lint`, `inspect`, `index`, and lexical search in temp vaults
  - finish with `make verify`

## Assumptions

- No new repair/migration command in this change.
- `ingest url` without credentials and vault auto-discovery outside `.kb/vault` remain configuration limits, not bugs.
