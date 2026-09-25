# Codebase maintenance audit — 2026-09-24

## Scope and acceptance

Review the Go codebase for demonstrated bugs, unnecessary complexity and technical
debt; fix confirmed findings in small commits on `main`; update Go dependencies;
publish only after local checks pass and verify CI on the published commit.
`make verify` (including zero lint findings) and the existing integration suite
are delivery gates. Preserve user-owned content and unrelated worktree changes.

Baseline: `734b54b4912b79dfa538467eaf84333a015c955c`, clean `main`.
Baseline `make verify`: passed. GitHub CI run `36092228951`: successful.

## Coverage and findings (in progress)

| Area | Evidence / finding | Status |
| --- | --- | --- |
| Dependency graph and build tooling | Updated required modules, Go minimum/CI/docs to 1.26, Mage fallback, lint/modernize and release tooling. Removed the unused tree-sitter replacement. Applied the Go 1.26 analyzer's behavior-preserving changes. | Local gates passed; publishing |
| Package boundaries | `magefile.go` checks only four hard-coded directories, including removed `internal/kodebase`, and treats every grep failure as success. The stated rule applies to all internal packages. | Confirmed; fix pending |
| Corpus/frontmatter persistence, review and quarantine | Receipts and quarantine append without separating an unterminated last row, losing the new record on reload. State truncation uses offsets captured when the store opened. | Reproduction/fixes pending |
| Decision/generation/session boundaries | Budget, receipt validity, cache and cancellation review. | Pending |
| Ingestion, conversion and media | File/network/subprocess boundaries and cancellation review. | Pending |
| CLI, topics, contracts, OKF and review actions | Validation and mutation contracts. | Pending |
| Retrieval, links, lint and QMD | Ranking, paths and subprocess failure handling. | Pending |
| Codebase scan, adapters, graph, metrics and vault | Graph correctness, reproducibility and file ownership. | Pending |

Uninspected areas are pending, not evidence of an absence of defects. This audit
does not claim that tests prove the absence of every possible bug.

## Delivery evidence

The final commit list, local checks and exact published-head CI results will be
recorded as each portion is completed.

- Dependency packet: `make verify` and `make test-integration` passed with the
  updated Go 1.26 sources and dependency versions. Local `/dev/full` test skips on
  macOS; Linux CI exercises it. `go mod verify` passed.
- `go list -m -u -json all`, checked against all 46 `go.mod` requirements: no
  updates pending. The graph still lists nine older modules with no package in
  the module's dependency/test closure; `go mod tidy` removes their explicit
  upgrade constraints. No unused requirements or replacements were retained to
  make that graph-only inventory appear current.
- Release orchestrator v0.0.29: the existing `pr-release --force
  --enable-rollback --ci-output` flags verified against its help output.
