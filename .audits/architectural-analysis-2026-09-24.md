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
| Dependency graph and build tooling | Updated required modules, Go minimum/CI/docs to 1.26, Mage fallback, lint/modernize and release tooling. Removed the unused tree-sitter replacement. Applied the Go 1.26 analyzer's behavior-preserving changes. | Published `d9ecec9`; CI and Release workflow successful |
| Package boundaries | Replaced the four-directory grep check with parsed imports across internal packages, including tests and platform files. Missing/unreadable sources fail; fixtures and the CLI layer are excluded. | Published `18044ed`; CI and Release workflow successful |
| Append-only decision logs | Receipt and quarantine recovery tests reproduced lost rows after an interrupted append. Reused review's tail-separation algorithm in `internal/jsonl` for receipts, review, quarantine, skips and inserted links. Writes sync before returning. Link and review actions now share the insertion record writer. Failed receipt loads no longer poison the cache with an empty map. | Published `93611eb`; CI and Release workflow successful |
| State persistence | Independent writers lost acknowledged rows during stale-tail repair and compaction. State mutations now share an OS file lock; repair reads the current tail and compaction reloads the current log. Normal appends read one tail byte, avoiding a full-log scan per write. | Fixed; race suite, cross-compilation and full gates passed; broader frontmatter ownership review ongoing |
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
- `d9ecec9`: GitHub CI `36093489452` and Release `36093489476` both succeeded.
- Boundary regression: compiled the old Mage command and ran it on an isolated
  `internal/ingest` package importing `internal/cli`: it incorrectly exited 0.
  The fixed command exits 1 and identifies that source file. The owning
  `internal/repohealth` race suite and `make verify` passed (1,857 tests, one
  platform skip).
- `18044ed`: GitHub CI `36093735496` and Release `36093735509` succeeded.
- Record recovery regressions failed before the production fixes: cache-only
  receipt reloads missed persisted results and quarantine listing lost the move
  record. A failed receipt read also prevented loading a subsequently repaired
  file. All now pass; existing skip rescue, link insertion and concurrent review
  append suites protect their record contracts without duplicating low-level
  helper tests. `make verify`: 1,864 tests; `make test-integration`: 1,980 tests;
  both passed with the existing single macOS `/dev/full` skip.
- `93611eb`: GitHub CI `36094298588` and Release `36094298551` succeeded.
- State regression: stale repair reverted an updated row and removed another;
  concurrent compaction lost all 40 acknowledged rows in the reproduction.
  Both deterministic stale-store cases and concurrent independent appends plus
  compaction now pass with `-race`. Existing malformed-middle-line rejection
  remains intact. Corpus test binaries compile for Linux and Windows; execution
  was local macOS (Linux execution also runs in CI).
- The CLI lint integration fixture itself creates the state lock during setup.
  Replaced its incorrect fixed filename allowlist with a before/after comparison
  of all record names and bytes, strengthening the read-only assertion. Final
  `make verify`: 1,868 tests; integration: 1,984 tests; both passed with one
  existing macOS platform skip.
