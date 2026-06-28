---
schema_version: 1
review_kind: implementation
round: 1
verdict: FIX_BEFORE_SHIP
reviewer_runtime: claude
reviewer_model: opus
generated_at: 2026-06-28T01:28:15Z
---

# Summary

The OKF dual-mode MVP is implemented coherently and matches the accepted ADRs: the `mode` field, OKF scaffold, `[okf].types` config, mode-aware `LinkFormatter` with the full ~43-site migration, the mechanical/non-destructive `promote`, the lenient `kb okf check`, both documented deviations (four producer fields + relative links), and the vendored official fixtures all build, lint clean, and pass unit + integration tests (including the three official-bundle conformance tests and the end-to-end promote — the PRD's explicit success metrics). The one blocker is the omission of the ADR-004/TechSpec-mandated N-001 "OKF render smoke" test, which was the accepted condition for migrating the dormant render-pipeline OKF branch now rather than in Phase 3, leaving that branch shipped without any test exercising it through a render function.

# Blockers

## B-001 — N-001 OKF render-smoke test mandated by ADR-004 is missing

- File: internal/vault/render.go
- Line: 191
- Issue: ADR-004 (Risks → "Dormant OKF render branch ships as dead code in the MVP (N-001)") and the TechSpec Testing Approach ("OKF render smoke (N-001) — render one document for a `mode: okf` topic and assert the OKF link branch and its `fromDir` plumbing execute and emit a relative markdown link, so the wired render branch is live, not shipped dead") make this smoke test the explicit, accepted justification for migrating all ~43 link call sites now. The migration was done (every site routes through `linkFor(topic, "", target, label)`), but no test ever invokes a render function (`RenderDocuments`/`renderDashboard`/`renderConceptIndex`/etc.) with `TopicMetadata{Mode: TopicModeOKF}`. The only OKF-mode coverage is `internal/vault/pathutils_test.go:232`, which calls `vault.LinkFormatterFor(...).Link(...)` directly — it exercises the formatter selection but not the render-site `linkFor` plumbing with `fromDir`. The render-site OKF branch therefore ships effectively untested, which is precisely the dead-code outcome N-001 was added to prevent. The codex plan's Test Plan silently dropped N-001 from the TechSpec.
- Rationale: This is an architectural decision recorded in an accepted ADR that is "implemented differently than specified" — the migration landed but its required guard test did not. Project rules (CLAUDE.md / AGENTS.md "Treat test failures as behavior bugs first"; the global "no partial deliverables" rule) and the review's dead-code / test-shape criteria make shipping a wired-but-unexercised branch without the spec-mandated smoke test a fix-before-ship gap. Remediation is local and the change shape is correct.
- Suggested fix: Add one table-friendly test in `internal/vault` that builds a minimal graph/metrics fixture (reuse `testTopicFixture()`), sets `topic.Mode = models.TopicModeOKF`, calls `vault.RenderDocuments(...)` (or a single render helper), and asserts at least one emitted body contains a relative markdown link (`[Label](...md)`) and no `[[wikilink]]`, confirming the OKF branch and `fromDir` plumbing execute through the render path.

# Risks

## R-001 — No structured logging despite the TechSpec observability contract

- File: internal/okf/okf.go
- Line: 1
- Issue: The TechSpec "Monitoring and Observability" section enumerates required `internal/logger` (slog) events for both verbs — `promote` source/target/type, link-transform count, broken-link count, `promote rejected: target not okf` (N-202), `description fallback used`; and `okf check` files-scanned, `reserved file skipped` (N-205), errors/warnings/strict flag. The `internal/okf` package and the new `internal/cli/promote.go` / `internal/cli/okf.go` import no logger and emit none of these. Operators get only the final JSON/table output and have no debug trail for unresolved links, description fallbacks, or skipped files.
- Suggested fix: Thread `internal/logger` into `Promote`/`Check` (or log at the CLI layer from `ConceptResult`/issues) for the events listed in the TechSpec, at debug/info level.

## R-002 — New user-facing commands are not documented in the project CLI surface

- File: CLAUDE.md
- Line: null
- Issue: `kb promote`, the `kb okf` group / `kb okf check`, and `topic new --mode` are absent from the CLI command tables in CLAUDE.md ("CLI Surface"), AGENTS.md ("CLI Commands"), and README. `config.example.toml` was correctly updated with `[okf]`, but the command docs now drift from the shipped surface, and these files are the canonical agent/operator references for what `kb` exposes. A hard-cut change should touch code, config, and docs together.
- Suggested fix: Add `kb promote`, `kb okf {check}`, and the `topic new --mode wiki|okf` flag to the CLI tables/notes in CLAUDE.md and AGENTS.md (and README if it lists commands), in the same change.

## R-003 — `kb okf check` cannot validate arbitrary/external bundles the checker supports

- File: internal/cli/okf.go
- Line: 66
- Issue: ADR-005 states the checker "takes the resolved bundle (topic) path; it does not require `mode: okf` so it can validate arbitrary directories and external bundles." The underlying `okf.Check` honors this (the integration test runs it on raw `testdata/official/...` paths). The CLI, however, resolves the argument through `ktopic.Info(vaultPath, topicSlug)`, which requires a valid kb topic (e.g. a readable `CLAUDE.md`/`topic.yaml`), so `kb okf check <some-external-dir>` cannot be run from the CLI. This is acceptable for the MVP (the operator's own bundle is a kb topic and the success metric is met via the direct `Check` integration test), but it narrows the documented capability and will need revisiting for Phase-2 ingest.
- Suggested fix: Note the MVP limitation, or allow `okf check` to fall back to treating the argument as a direct bundle path when topic resolution fails, so it matches the `Check` contract.

# Nits

## N-001 — Dead exported function `WriteMetadataFileWithMode`

- File: internal/topic/topic.go
- Line: 561
- Issue: `WriteMetadataFileWithMode` is exported but has zero callers (production or test); `internal/vault/writer.go` still uses the mode-less `WriteMetadataFile`. It is speculative dead API surface (the linter won't catch it because it is exported).
- Suggested fix: Delete `WriteMetadataFileWithMode`, or route the one real caller (`writer.go:106`) through it with an explicit wiki mode if a written `mode:` line is desired there.

## N-002 — `firstBodySentence` returns the entire body when no sentence delimiter exists

- File: internal/okf/okf.go
- Line: 446
- Issue: The TechSpec description fallback is "first non-empty body sentence (markdown-stripped) → else empty + warning." When the body has text but no `. `/`! `/`? ` delimiter, `firstBodySentence` returns the whole normalized body, producing a potentially very long `description` with no warning instead of a single sentence.
- Suggested fix: Cap the no-delimiter fallback (e.g. first line, or first N chars) or emit the empty-description warning when no sentence boundary is found.

## N-003 — Unreachable symlink-directory branch in the conformance walk

- File: internal/okf/okf.go
- Line: 177
- Issue: In `Check` (and identically in `loadConcepts`), `if entry.Type()&fs.ModeSymlink != 0 { if entry.IsDir() { return filepath.SkipDir } ... }` — for a symlink, `fs.DirEntry.IsDir()` is always false (it reflects `Type().IsDir()`), and `filepath.WalkDir` never descends into symlinked directories anyway, so the inner `SkipDir` branch is dead.
- Suggested fix: Drop the `entry.IsDir()` sub-branch for symlinks (always `return nil`), or add a comment noting it is defensive.

# Evidence

- Read in full: PRD (`_prd.md`), TechSpec (`_techspec.md`), ADR-001..006, codex plan `20260627-220335-okf-dual-mode.md`; project rules CLAUDE.md, AGENTS.md, CONTRIBUTING.md.
- Read in full: `internal/okf/okf.go`, `internal/okf/okf_test.go`, `internal/okf/official_integration_test.go`, `internal/cli/promote.go`, `internal/cli/okf.go`, `internal/cli/okf_test.go`, `internal/cli/okf_integration_test.go`, `internal/cli/topic.go` (diff), `internal/topic/topic.go` (diff + on-disk), `internal/vault/pathutils.go`, `internal/frontmatter/frontmatter.go`, `internal/models/kb_models.go`; reviewed the raw patch hunks for `models.go`, `config.go`, `config_test.go`, `generate.go`, `render.go`, `render_wiki.go`, `writer.go`, `topic_test.go`, `pathutils_test.go`, and the `okf-claude-template.md` asset.
- `go build ./...` — Success.
- `go test ./internal/okf/... ./internal/cli/... ./internal/topic/... ./internal/vault/... ./internal/config/...` — all `ok`.
- `go test -tags integration ./internal/okf/... ./internal/cli/...` — all `ok` (official ga4/stackoverflow/crypto_bitcoin lenient conformance + end-to-end promote/check pass).
- `golangci-lint run` (v2.11.4) over the changed packages — `0 issues.`
- Inspected `internal/okf/testdata/official`: ga4 (17 md), stackoverflow (53 md), crypto_bitcoin (8 md); root and nested `index.md` files carry no frontmatter; every concept `.md` has a `type:`; no `log.md`/README/LICENSE inside the individual bundles; only `viz.html` non-md files — confirming the lenient checker classifies the corpus correctly.
- Verified zero-regression coverage: existing `internal/vault/render_test.go` asserts exact `[[demo-repo/...|...]]` wikilink output through `RenderDocuments`, so the wiki branch of the migration is byte-verified by the existing suite (ADR-004 mitigation satisfied for the wiki side only).
- Confirmed no logger usage in `internal/okf` / new CLI files, and no mention of the new commands in CLAUDE.md/AGENTS.md/README.
- Limitation: did not run the full `make verify` (mage/mise gating noted in the plan); evidence above used direct `go`/`golangci-lint` invocations on the affected packages.

# Deferred Or Follow-Up

- `--to <topic>/<subdir>` promotion is intentionally deferred (codex plan Assumptions); flat-root concepts only for MVP — no action needed now, but track for Phase 2.
- Phase 2/3 surfaces (`kb okf export`, `kb okf ingest`, codebase→OKF) are out of scope and correctly absent.
- Consider consolidating the duplicated bundle-walk logic shared by `Check` and `loadConcepts` (entry filtering, exclusion set, posix-rel derivation) into one helper before Phase 2 grows the package.
