# Task Memory: task_05.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot

- Implement Task 05 end-to-end on the current branch by adding the Go adapter and the missing `internal/models` types/interfaces it depends on, then validate with task-specific tests, coverage, and `make verify`.

## Important Decisions

- Use `internal/models` and `internal/adapter` as the canonical package paths for this task. `_techspec.md` still mentions `internal/kodebase/...`, but the active repository layout and downstream task files use the shorter paths.
- Keep scope tight to Task 05: add only the model surface needed for adapter parsing and tests instead of silently finishing unrelated Task 02 deliverables.
- Keep call-graph extraction syntactic and package-scoped: resolve only direct identifier calls against parsed package-level functions, including cross-file calls within the same `ParseFiles` invocation, and ignore selector calls such as `fmt.Println`.

## Learnings

- The current branch only contains the scanner-facing subset of `internal/models`; Task 05 cannot compile without restoring graph, relation, diagnostic, and adapter interface types locally.
- The TypeScript Go adapter only emits `calls` relations for direct identifier calls inside function or method bodies; selector calls such as `fmt.Println` are intentionally ignored for call-graph edges.
- The current branch still uses `module github.com/compozy/kb`; new Task 05 code must import `internal/models` through that module path until Task 01 lands.
- Tree-sitter Go exposes aliases as `type_alias`, not `type_spec`, so alias extraction must handle both node kinds.

## Files / Surfaces

- `internal/models/models.go`
- `internal/models/models_test.go`
- `internal/adapter/treesitter.go`
- `internal/adapter/go_adapter.go`
- `internal/adapter/go_adapter_test.go`
- `internal/adapter/go_adapter_integration_test.go`
- `.compozy/tasks/rewrite/task_05.md`
- `.compozy/tasks/rewrite/_tasks.md`

## Errors / Corrections

- Corrected the package-path ambiguity by following the active `internal/models` layout instead of the stale `internal/kodebase/...` paths in `_techspec.md`.
- Corrected the first implementation pass to treat `type_alias` as a valid type declaration node.
- Fixed `extractLeadingComment` so module docs only use the file's true leading comment block instead of matching later line comments.

## Ready for Next Run

- Verification complete:
  - `go test ./internal/models ./internal/adapter`
  - `go test -tags integration ./internal/adapter`
  - `go test -coverprofile=/tmp/kodebase-go-adapter.cover ./internal/adapter && go tool cover -func=/tmp/kodebase-go-adapter.cover` => 86.3%
  - `make verify`
- Update tracking files and create the local code commit without staging workflow memory or task-tracking artifacts.
