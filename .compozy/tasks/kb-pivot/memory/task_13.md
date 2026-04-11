# Task Memory: task_13.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot

- Implement `internal/lint` for topic-local KB structural checks: dead wikilinks, orphan wiki concepts, missing raw sources, stale wiki content, and frontmatter/schema violations.
- Expose lint results in a formatter-friendly shape for task 17 and support optional markdown report writing to `outputs/reports/<date>-lint.md`.
- Close the task only after targeted coverage hits the required 80%+ bar and `make verify` passes cleanly.

## Important Decisions

- Kept the package API small: `Lint(topicPath)`, `Columns()`, `Rows()`, and `SaveReport(topicPath, issues, now)`.
- Dead-link/orphan graphing uses markdown body wikilinks only; frontmatter `sources:` is handled separately as missing-source/stale logic to avoid double-reporting metadata links as generic dead links.
- Format validation is path-driven from the Karpathy KB schemas: managed files under `raw/`, `wiki/concepts/`, `wiki/index/`, and `outputs/` are checked for required fields and expected `type`/`stage` values.
- The vault walker skips the root `AGENTS.md` symlink so topic scaffolding does not create duplicate lint issues from the same `CLAUDE.md` content.

## Learnings

- Supporting both title-based wikilinks (for human-authored KB notes) and topic-prefixed path links (from generated codebase docs) is necessary for the lint engine to match existing vault patterns.
- A separate `SaveReport` helper is cleaner than baking report persistence into `Lint`, and still satisfies task 13 while leaving task 17 free to gate report writes on `--save`.
- Package coverage reached 80.2% with the required checks and report-writing path covered by fixture tests.

## Files / Surfaces

- `internal/lint/lint.go`
- `internal/lint/lint_test.go`
- `.compozy/tasks/kb-pivot/memory/task_13.md`
- `.compozy/tasks/kb-pivot/memory/MEMORY.md`

## Errors / Corrections

- First `make verify` run failed on three `staticcheck` findings inside the new lint package: one regexp literal style warning and two `fmt.Sprintf`/`WriteString` inefficiencies in report rendering.
- Fixed the warnings, reran package tests, and reran `make verify` from scratch to get a clean completion signal.

## Ready for Next Run

- Task 13 code is verified with `go test ./internal/lint/...`, `go test ./internal/lint/... -coverprofile=/tmp/lint.cover && go tool cover -func=/tmp/lint.cover` (80.2%), and full `make verify`.
- Task 17 can wire `kb lint` directly to `lint.Lint`, `lint.Columns`, `lint.Rows`, and call `lint.SaveReport` when `--save` is set.
