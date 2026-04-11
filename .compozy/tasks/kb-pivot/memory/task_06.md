# Task Memory: task_06.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot

- Add `internal/convert/html.go` plus tests for `.html`/`.htm` conversion, title extraction, script/style stripping, and registry registration per task_06.

## Important Decisions

- Implement HTML conversion around a parsed DOM so the task can explicitly strip `script` and `style` nodes, extract `<title>` metadata, and then hand the cleaned document to a configured html-to-markdown converter.
- Use a custom html-to-markdown v2 converter with base, commonmark, and table plugins because the package-level helpers do not enable table rendering.

## Learnings

- The html-to-markdown table support is not part of the package-level helpers, so the reusable helper must instantiate its own converter with the table plugin to satisfy both direct HTML ingest and future EPUB chapter conversion.
- Falling back to raw `<h1>` text is fragile for malformed documents because broken nesting can pull paragraph text into the heading; extracting heading text while skipping nested block nodes keeps the fallback title stable.

## Files / Surfaces

- `internal/convert/registry.go`
- `internal/convert/registry_test.go`
- `internal/convert/html.go`
- `internal/convert/html_test.go`
- `internal/convert/testdata/`
- `go.mod`
- `go.sum`

## Errors / Corrections

- Initial malformed-HTML title fallback used full descendant text from the first `<h1>`, which captured nested paragraph text in broken markup. Reworked the fallback to ignore nested block descendants so the title remains the heading text.

## Ready for Next Run

- Package test evidence: `go test ./internal/convert -cover` passed at 84.8% coverage.
- Repository gates passed once on the code changes: `make lint` and `make verify`.
- Remaining closeout after this note: update task tracking, rerun the final verification gate, and create the local commit without staging tracking-only files.
