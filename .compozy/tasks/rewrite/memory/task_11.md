# Task Memory: task_11.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot
- Implement the vault writer for Task 11 in `internal/vault` so rendered markdown documents and base definitions are persisted under the topic directory, with manifest/log management and stale generated wiki concept cleanup.
- Finish with task-specific tests, coverage evidence, clean `make verify`, tracking updates, and one local commit.

## Important Decisions
- Follow the upstream writer contract from `write-vault.ts`: `RenderedDocument.Body` is already the final markdown payload, so the writer must persist `document.Body` directly instead of re-rendering frontmatter.
- Mirror the reference topic skeleton, including `raw/articles`, `raw/bookmarks`, `raw/github`, `raw/codebase`, `wiki/concepts`, `wiki/index`, `outputs/briefings`, `outputs/queries`, `outputs/diagrams`, `outputs/reports`, and `bases`.
- Remove stale generated wiki concepts by deleting concept markdown files whose YAML frontmatter marks `generator: kodebase`, while preserving unmanaged/manual concept files.
- Validate managed document contracts before writing: kind-to-area mapping, expected subdirectory prefixes, and required YAML frontmatter/body shape are enforced at the writer boundary.

## Learnings
- The shared workflow memory already settled the renderer/writer contract: markdown documents carry their own YAML frontmatter, while `.base` files must still be rendered from `BaseDefinition`.
- The TypeScript reference also rewrites `raw/codebase/` and `wiki/index/` wholesale on each run, which prevents stale raw/index notes from surviving regeneration.
- The topic manifest remains the root topic marker (`CLAUDE.md`), and the Go writer now mirrors the reference behavior of refreshing `AGENTS.md` as a symlink to that manifest.
- Task-specific validation completed cleanly with `go test ./internal/vault`, `go test -tags integration ./internal/vault`, package coverage `85.9%`, and `make verify`.

## Files / Surfaces
- `internal/vault/writer.go`
- `internal/vault/writer_test.go`
- `internal/vault/writer_integration_test.go`
- `.compozy/tasks/rewrite/task_11.md`
- `.compozy/tasks/rewrite/_tasks.md`
- `.compozy/tasks/rewrite/memory/task_11.md`

## Errors / Corrections
- The repo-local required skills live under `.agents/skills/`, not `~/.agents/skills/`; use the repo copies for this task.
- `_techspec.md` still points at stale `internal/kodebase/...` paths. The active branch layout for this task is `internal/models` and `internal/vault`.
- The first `make verify` run failed on `errcheck` because `appendLog` did not handle `Close()` on the opened log file. Fixed by closing explicitly on both the success and write-error paths, then reran verification cleanly.

## Ready for Next Run
- Task complete. Task tracking files are updated in the workspace and should remain out of the code commit unless the repo later requires tracking artifacts to be staged.
