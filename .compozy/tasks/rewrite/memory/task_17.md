# Task Memory: task_17.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot
- Completed Task 17 by implementing `internal/qmd/client.go` plus unit/integration tests so later CLI work can search/index vault topics through the real QMD shell commands with graceful missing-binary handling.

## Important Decisions
- Treat the task's hybrid/lex/vector API as a client-level abstraction and map it onto the current QMD CLI commands (`query`, `search`, `vsearch`) because the live binary does not provide a documented `search --mode ...` flag.
- Parse QMD JSON search output directly, and parse the human-readable `collection add`, `update`, and `status` summaries into structured Go results for index/status reporting.
- Keep QMD test isolation entirely inside temporary XDG cache/config roots so the user's existing QMD collections do not bleed into integration runs.

## Learnings
- `npx @tobilu/qmd query --json "QMD test"` returns a JSON array with fields like `docid`, `score`, `file`, `title`, and `snippet`; `--full` swaps `snippet` for `body`.
- `npx @tobilu/qmd collection add <path> --name <collection>` both creates the collection and performs the initial index sync.
- `npx @tobilu/qmd update <collection>` refreshes a single collection and prints an `Indexed: <new> new, <updated> updated, <unchanged> unchanged, <removed> removed` summary.
- `npx @tobilu/qmd status --json` and `collection list --json` still emit human-readable text, so index/status parsing must not rely on JSON there.
- Search command JSON is emitted on stdout, while warnings/progress (including hybrid query expansion traces) stay on stderr. The client can safely unmarshal stdout while preserving stderr for failures.
- `qmd embed` can fail on this host with `sqlite-vec is not available`, so Task 17 verification kept integration coverage on add/search without forcing embeddings.

## Files / Surfaces
- `internal/qmd/client.go`
- `internal/qmd/client_test.go`
- `internal/qmd/client_integration_test.go`
- `/Users/pedronauck/dev/projects/kodebase/packages/cli/src/integrations/qmd-client.ts`
- `/Users/pedronauck/dev/projects/kodebase/packages/cli/src/commands/search.ts`
- `/Users/pedronauck/dev/projects/kodebase/packages/cli/src/commands/index-vault.ts`

## Errors / Corrections
- Corrected the task/techspec assumption that QMD exposes `search --mode`; the live CLI requires command translation instead.
- Corrected the integration test isolation after QMD reused collection state from outside the temporary cache root; setting `XDG_CONFIG_HOME` alongside `XDG_CACHE_HOME` and `HOME` fixed the leak.

## Ready for Next Run
- Task 17 is complete. Verification evidence:
  - `go test ./internal/qmd`
  - `go test -coverprofile=/tmp/kodebase-qmd.cover ./internal/qmd && go tool cover -func=/tmp/kodebase-qmd.cover` => 82.4%
  - `go test -tags integration ./internal/qmd`
  - `make verify`
- Task 18 can consume `internal/qmd.QMDClient` for real CLI wiring.
