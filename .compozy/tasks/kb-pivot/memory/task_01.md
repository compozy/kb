# Task Memory: task_01.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot
- Add `internal/frontmatter` with YAML frontmatter parse/generate helpers and typed accessors for KB markdown files.
- Keep scope to the new package plus unit tests; do not refactor `internal/vault` consumers in this task.

## Important Decisions
- Treat the PRD/techspec as the approved design baseline for this task and implement directly against that scope.
- Preserve deterministic generation order by sorting keys during YAML node construction instead of relying on Go map iteration.
- Encode `time.Time` frontmatter values as `YYYY-MM-DD` to match the KB schemas and allow parse/generate round-trips for date fields.

## Learnings
- `gopkg.in/yaml.v3` unmarshals plain `YYYY-MM-DD` scalars into `time.Time` when decoding into `map[string]any`.
- A line-based delimiter scan is simpler than the existing regex for supporting empty frontmatter blocks like `---` / `---` without false negatives.

## Files / Surfaces
- Added: `internal/frontmatter/frontmatter.go`
- Added: `internal/frontmatter/frontmatter_test.go`
- Reference only: `internal/vault/reader.go`, `internal/vault/render.go`, `.agents/skills/karpathy-kb/references/frontmatter-schemas.md`

## Errors / Corrections
- `make lint` initially failed on a De Morgan simplification in `TestGenerateProducesValidDelimitedYAMLWithSortedKeys`; corrected the condition and re-ran lint plus full verification.

## Ready for Next Run
- Frontmatter package is implemented and verified; future tasks can wire `internal/vault`, `internal/topic`, `internal/ingest`, and `internal/lint` to it without revisiting parse/generate behavior.
