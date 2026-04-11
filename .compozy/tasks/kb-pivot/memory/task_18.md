# Task Memory: task_18.md

Keep only task-local execution context here. Do not duplicate facts that are obvious from the repository, task file, PRD documents, or git history.

## Objective Snapshot

- Add end-to-end CLI integration coverage for `kb topic`, `kb ingest file`, `kb ingest codebase`, `kb inspect`, and `kb lint`.
- Refresh repo docs/config copy to the KB pivot surface.
- Re-verify the branch with `make test-integration` and `make verify` after the new coverage and doc updates land.

## Important Decisions

- Treat Makefile and Mage build-target work as already satisfied unless verification uncovers a real mismatch; current repo state already builds `kb` from `cmd/kb`.
- Keep the new integration coverage in `internal/cli` because the task exercises full Cobra command flows rather than package-level helpers.
- Cover the file-ingest workflow with real `.txt` and `.csv` inputs so the task validates both frontmatter/path writing and CSV-to-Markdown conversion through the CLI.
- Fix the pre-existing `internal/qmd` test flake in this task because `make test-integration` is part of the verification evidence for the PRD gate.

## Learnings

- `make verify` already passes on the pre-change branch, so the missing acceptance criteria are additive workflow tests and stale documentation, not an existing red baseline.
- `CLAUDE.md` and `config.example.toml` still describe the old `kodebase` surface even though the binary and root command are already `kb`.
- `make test-integration` initially failed in `internal/qmd` with intermittent `fork/exec .../qmd: text file busy`; the stable fix was to run fake QMD shell scripts via `/bin/sh` in the test client rather than execing the generated script path directly.

## Files / Surfaces

- `internal/cli/` integration tests
- `internal/qmd/client_test.go`
- `CLAUDE.md`
- `config.example.toml`

## Errors / Corrections

- Writing the fake QMD script to a temp path and renaming it into place did not remove the `ETXTBSY` flake by itself; the final correction was changing the fake client execution path to `/bin/sh <script> ...`.

## Ready for Next Run

- Validation evidence for this task is `go test -tags integration ./internal/cli`, `make test-integration`, and `make verify`.
