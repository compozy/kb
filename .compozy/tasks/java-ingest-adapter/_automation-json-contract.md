# Java Ingest Automation JSON Contract

This document defines the minimum stable JSON contract for automation consumers of:

- `kb ingest codebase ...` stdout payload (`codebaseIngestResult`)
- Embedded `summary` payload (`models.GenerationSummary`)

The contract applies to Java ingest workflows and remains compatible with other language ingest flows that use the same payload shape.

## Stability Guarantee

Automation consumers may rely on the presence and meaning of the required keys below.

### Required top-level keys (`codebaseIngestResult`)

- `topic` (`string`): resolved topic slug.
- `sourceType` (`string`): must be `codebase-file`.
- `filePath` (`string`): topic-relative path for codebase raw area.
- `title` (`string`): resolved topic title.
- `summary` (`object`): generation summary.

### Required summary keys (`GenerationSummary`)

- `command` (`string`)
- `rootPath` (`string`)
- `vaultPath` (`string`)
- `topicPath` (`string`)
- `topicSlug` (`string`)
- `dryRun` (`bool`)
- `detectedLanguages` (`string[]`)
- `selectedAdapters` (`string[]`)
- `filesScanned` (`number`)
- `filesParsed` (`number`)
- `filesSkipped` (`number`)
- `symbolsExtracted` (`number`)
- `relationsEmitted` (`number`)
- `rawDocumentsWritten` (`number`)
- `wikiDocumentsWritten` (`number`)
- `indexDocumentsWritten` (`number`)
- `timings` (`object`)
- `diagnostics` (`array`)

### Required timings keys (`GenerationTimings`)

- `scanMillis`
- `selectAdaptersMillis`
- `parseMillis`
- `normalizeMillis`
- `metricsMillis`
- `renderMillis`
- `writeMillis`
- `totalMillis`

## Value Semantics

- `topic` and `summary.topicSlug` represent the same topic identity and must match.
- `sourceType` must remain `codebase-file`.
- For full ingest (`summary.dryRun=false`), `rawDocumentsWritten` should be `> 0` when files are discovered.
- For dry-run (`summary.dryRun=true`), write counters must remain `0`:
  - `rawDocumentsWritten`
  - `wikiDocumentsWritten`
  - `indexDocumentsWritten`

## Compatibility Policy

- Backward compatibility is the default.
- Existing required keys must not be removed or renamed without an explicit versioning plan.
- Additive changes are allowed (new optional keys) if existing required keys and semantics remain unchanged.
- Breaking changes require:
  1. explicit contract versioning,
  2. migration guidance for automation consumers,
  3. updated CLI unit + integration coverage for both old/new behaviors during transition.

## Non-Guaranteed Surfaces

Consumers should not hard-couple to:

- exact ordering of arrays (`detectedLanguages`, `selectedAdapters`, `diagnostics`),
- exact diagnostic counts/messages/detail text,
- absolute timing values or performance ratios,
- incidental extra keys introduced in future additive releases.

## Verification Coverage

Contract enforcement is covered by:

- unit tests in `internal/cli/ingest_test.go` for required keys and mode semantics;
- integration tests in `internal/cli/workflow_integration_test.go` for Java full-run and dry-run contract stability.
