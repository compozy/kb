# Task Memory: task_20.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot

- Update `AGENTS.md`, `CLAUDE.md`, and `config.example.toml` so they describe the active kodebase Go implementation, then finish with clean verification and tracking updates.

## Important Decisions

- Treated the live repository state, not stale techspec path examples, as the source for package layout and CLI command documentation.
- Kept `config.example.toml` limited to keys actually accepted by `internal/config` and documented that generate/search/index tuning still lives on CLI flags.
- Treated the macOS linker warning emitted by `go run` tool bootstrapping as a real verification blocker and fixed Mage to prefer installed `golangci-lint` / `gotestsum` binaries before falling back to pinned module execution.

## Learnings

- The verification stack already passes on this branch: `make fmt`, `make lint`, `make test`, `make build`, `make test-integration`, and `make verify` all completed successfully without production code changes.
- Aggregate unit-test coverage measured with `go test -coverprofile` is `82.8%`.
- The built binary at `bin/kodebase` responds successfully to `--help` for the root command and every documented subcommand.
- Preferring installed verification binaries removes the macOS `ld: warning: -bind_at_load is deprecated` noise from the lint step while preserving pinned `go run` fallbacks for environments that do not have those tools on `PATH`.

## Files / Surfaces

- `AGENTS.md`
- `CLAUDE.md`
- `config.example.toml`
- `Makefile`
- `magefile.go`
- `.compozy/tasks/rewrite/task_20.md`
- `.compozy/tasks/rewrite/_tasks.md`

## Errors / Corrections

- Initial `--help` sweep failed because the shell loop used `commands` as a variable name in `zsh`; reran with a non-reserved variable name and the CLI checks passed.

## Ready for Next Run

- Task complete.
- Local commits created:
  - `661126f` (`docs: refresh project docs and config example`)
  - `8483c67` (`build: stabilize verification tooling`)
- Tracking files were updated locally but intentionally left out of the commit because `.compozy/` is tracking-only workspace state.
