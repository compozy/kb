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
| State persistence and frontmatter ownership | Independent writers lost acknowledged rows during stale-tail repair and compaction. State mutations now share an OS file lock; repair/compaction read current data. Alias merging also deleted existing values and overwrote unrecognized shapes; it now preserves existing entries and skips unknown formats. | State fix published `784a10f`, CI/Release successful; ownership fix verified locally |
| Decision/generation/session boundaries | Queue admission now covers accounting/fatal state and cache reuse. Invalid output envelopes lost billed usage; authentication/cancellation lost entire call receipts. Both clients now account for usage before decoding output and retain interrupted attempts. | Queue fix published `5c08099`, CI/Release successful; accounting fix verified locally |
| Ingestion, conversion and media | Reviewed exclusive raw writes, gate orchestration, Firecrawl and media subprocess/cancellation paths. JSON conversion rounded large integers/precise decimals, underflowed small numbers and rejected valid large exponents through float64 decoding; numeric literals now remain exact. | Converter fix verified; CI exposed an executable-fixture race in the media suite |
| CLI, topics, contracts, OKF and review actions | Validation and mutation contracts. | Pending |
| Retrieval, links, lint and QMD | Ranking, paths and subprocess failure handling. | Pending |
| Codebase scan, adapters, graph, metrics and vault | Invalid rendered input deleted the previous codebase output before validation; a duplicate symlink writer overwrote manual `AGENTS.md`. Output validation now runs before mutations and the existing preserving scaffold owns AGENTS creation. Scanner, pipeline cancellation boundaries, graph normalization and metric entrypoints inspected. | Vault fixes verified; adapter/inspection review ongoing |

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
- JSON conversion: regression cases reproduced `9007199254740993` rounding
  down, precision loss in decimal metadata, underflow to zero, and rejection
  of valid `1e400`. The existing converter test now checks the serialized
  numeric value instead of requiring a float64 representation. JSON numbers
  remain exact via the standard decoder's `UseNumber`. `make verify` (1,893
  tests) and integration (2,009 tests) passed with the existing platform skip.
- CI `36096191996` on `75c9ab3` failed in an unchanged media assertion: the
  shell executable created by `writeFakeYTDLP` returned Linux `ETXTBSY` before
  metadata parsing. The fixture now uses the same `/bin/sh <script>` command
  boundary as the existing QMD suite, keeping the real backend/subprocess and
  every assertion. No production retry, test retry or skip was added.
  `go test -race ./internal/mediadl -count=20`, `make verify` (1,893 tests) and
  integration (2,009 tests) passed locally; Linux confirmation is via CI.
- CI also reported deprecated Node 20 runtimes in checkout/setup-go. Verified
  the official action releases and Node 24 manifests, then updated both to v7
  using the repository's existing major-version pinning convention.
- `784a10f`: GitHub CI `36094879117` and Release `36094879126` succeeded.
- Decision queue regressions reproduced two provider calls where the second
  should have stopped after budget exhaustion or authentication failure.
  Deterministic `testing/synctest` cases exercise the HTTP boundary without
  sleeps and also cover queued cache reuse. The owning engine race suite,
  `make verify` (1,872 tests), and integration (1,988 tests) passed; one existing
  macOS platform skip in each full suite.
- `5c08099`: GitHub CI `36095329968` and Release `36095329974` succeeded.
- Accounting regressions reproduced lost usage for null/malformed decision
  answers and malformed generation choices, plus missing receipts for
  authentication failures and cancellations on both routes. The client suites
  now retain billed cost, stop at the shared budget and preserve failed calls
  with unknown cost; existing error returns and invalid-output rejection stay
  intact. The owning race suites, `make verify` (1,879 tests), and integration
  (1,995 tests) passed, with the existing macOS platform skip.
- Vault regressions reproduced deletion of a previously generated document
  on invalid input and replacement of a manual `AGENTS.md`. Both now pass
  through the real filesystem writer. `make verify` (1,880 tests) and
  integration (1,996 tests) passed with the existing macOS platform skip.
- Provider metadata boundaries: a null token count became a reported zero.
  Large Retry-After values overflowed when converted before capping (the old
  expression produced a negative duration in an executed amd64 binary on this
  Mac; arm64 saturates differently). The fixes retain unknown tokens and cap
  numeric retry delays before conversion. Existing client suites own these
  regressions. `make verify` (1,882 tests) and integration (1,998 tests) passed
  with the existing macOS platform skip.
- `b565b17`: Release `36095700897` succeeded; CI `36095700917` was cancelled
  by the next main push under the workflow's concurrency policy. Its changes
  are included in `fca713c`, whose CI `36095873675` and Release `36095873664`
  both succeeded. Subsequent pushes wait for the preceding CI to finish.
- Go adapter: two `package main` directories incorrectly shared the function
  lookup and produced a call from one command to the other's helper. Package
  lookup now includes directory and package name. The owning adapter test
  covers local cross-file resolution, repeated names and external test
  packages. Full verification and integration passed.
- Alias regressions reproduced removal of duplicates/spacing and replacement
  of mixed lists, maps and numeric values with generated aliases. The writer
  now preserves all existing string entries and reports unknown shapes as
  user-owned keys. The existing append/deduplication test remains intact.
  `make verify` (1,888 tests) and integration (2,004 tests) passed with the
  existing macOS platform skip.
