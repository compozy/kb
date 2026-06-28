---
schema_version: 1
review_kind: implementation
round: 2
verdict: SHIP
reviewer_runtime: claude
reviewer_model: opus
generated_at: 2026-06-28T01:39:07Z
---

# Summary

The round-1 sole blocker (B-001, the ADR-004/TechSpec N-001 OKF render-smoke test) is now fixed by `TestRenderDocumentsUseOKFMarkdownLinkSyntax`, which drives `RenderDocuments` with `TopicModeOKF` and asserts a relative markdown link plus the absence of wikilinks through the live render branch; the three round-1 nits (dead `WriteMetadataFileWithMode`, whole-body description fallback, unreachable symlink branch) and the docs risk (R-002) are all correctly remediated and the full OKF MVP builds and passes 295 unit + 181 integration tests. No blockers remain and the change shape matches the accepted ADRs and TechSpec, so this ships; the two deferred risks (structured logging, external-bundle CLI checking) and minor test/code nits are acceptable follow-ups.

# Blockers

None.

# Risks

## R-001 — TechSpec observability contract still unimplemented

- File: internal/okf/okf.go
- Line: 1
- Issue: The TechSpec "Monitoring and Observability" section enumerates required `internal/logger` (slog) events for both verbs (`promote` source/target/type, link-transform count, broken-link count, `promote rejected: target not okf`, `description fallback used`; `okf check` files-scanned, `reserved file skipped`, errors/warnings/strict flag). `internal/okf`, `internal/cli/promote.go`, and `internal/cli/okf.go` still import no logger and emit none of these. Carried over from round-1 R-001 and explicitly deferred in the remediation. Operators get only the final JSON/table output with no debug trail for unresolved links, description fallbacks, or skipped files. This was classified as a non-blocking risk in round 1 and the success metrics (conformance on official bundles, end-to-end promote, zero regression) do not depend on it, so it does not block ship — but the spec'd contract remains unmet.
- Suggested fix: Thread `internal/logger` into `Promote`/`Check` (or log at the CLI layer from `ConceptResult`/issues) for the events listed in the TechSpec at debug/info level, in a follow-up.

## R-002 — `okf check` cannot validate arbitrary/external bundles the checker supports

- File: internal/cli/okf.go
- Line: 70
- Issue: ADR-005 states the checker "takes the resolved bundle (topic) path; it does not require `mode: okf` so it can validate arbitrary directories and external bundles," and `okf.Check` honors this (the integration test runs it directly on raw `testdata/official/...` paths). The CLI, however, resolves the argument through `ktopic.Info(vaultPath, topicSlug)`, which requires a valid kb topic, so `kb okf check <external-dir>` cannot be run from the CLI. Carried over from round-1 R-003 and intentionally deferred. Acceptable for the MVP (the operator's own bundle is a kb topic and the conformance success metric is met via the direct `Check` integration test), but it narrows the documented capability and will need revisiting for Phase-2 ingest. The shipped docs do not overclaim external-dir support, so there is no truthful-UI violation.
- Suggested fix: In a follow-up, fall back to treating the argument as a direct bundle path when topic resolution fails, matching the `Check` contract.

## R-003 — Lenient `okf check` hard-errors non-date `##` headings in `log.md`

- File: internal/okf/okf.go
- Line: 655
- Issue: `checkLogFile` emits `models.SeverityError` for every `## ` heading in `log.md` that does not parse as `YYYY-MM-DD`, and this is not gated by `--strict` (the `strict` parameter is discarded). The PRD/TechSpec posture is lenient-by-default so externally produced bundles pass; a real-world `log.md` containing a prose `## Section` heading would fail conformance hard even without `--strict`. This does not affect the MVP success metric (the three official fixtures ship no `log.md`, confirmed by round-1 evidence and the green integration run), so it is not a blocker, but it is a latent leniency gap for Phase-2 external-bundle ingest.
- Suggested fix: Treat malformed log headings as a `warningSeverity(strict)` diagnostic (warning by default, error under `--strict`), consistent with the rest of the lenient ruleset.

# Nits

## N-001 — Residual no-op discards and unused parameters in `internal/okf`

- File: internal/okf/okf.go
- Line: 83
- Issue: `_ = sourceTopicRoot` (line 83), `_ = bundlePath` / `_ = strict` (lines 650-651), and `_ = strict` (line 671) are dead no-op statements left after the round-2 dead-code sweep; `checkIndexFile`/`checkLogFile` carry parameters they never use, and `resolveSourceDocument` returns a `topicRoot` its only caller discards.
- Suggested fix: Drop the unused `topicRoot` return from `resolveSourceDocument` and remove the unused `bundlePath`/`strict` parameters (and their `_ =` discards) from `checkIndexFile`/`checkLogFile`.

## N-002 — Weak/incorrect log assertion in the OKF CLI integration test

- File: internal/cli/okf_integration_test.go
- Line: 87
- Issue: The promotion-log assertion is `if !strings.Contains(logContent, "## "+frontmatter.DateLayout[:4]) && !strings.Contains(logContent, "**Creation**")` — `frontmatter.DateLayout[:4]` is `"2006"`, so the first operand checks for the literal `"## 2006"`, which can never match a real date heading, and the `&&` makes the whole check pass on the `**Creation**` term alone. The date-heading half of the assertion provides false confidence and never actually validates the `## <date>` heading.
- Suggested fix: Assert both markers with `||` (fail if either is missing) and compare the heading against the formatted promotion date rather than `DateLayout[:4]`.

# Evidence

- Read in full: `_prd.md`, `_techspec.md`, ADR-004, round-1 findings (`impl-review-findings-round1.md`), round-1 remediation (`impl-review-remediation-round1.md`); project rules `CLAUDE.md`, `AGENTS.md`, `CONTRIBUTING.md`.
- Read in full: the round-2 raw patch (`impl-review-diff-round2.patch`, all 2997 lines), on-disk `internal/okf/okf.go`, `internal/vault/render_test.go`, and the relevant span of `internal/topic/topic.go`.
- Verified B-001 fix: `TestRenderDocumentsUseOKFMarkdownLinkSyntax` (internal/vault/render_test.go:266) drives `vault.RenderDocuments(...)` with `topic.Mode = models.TopicModeOKF` and asserts a relative markdown link (`[main (function)](raw/codebase/symbols/main--commands-run-ts-l1.md)`, `[Codebase Overview](wiki/codebase/concepts/Codebase%20Overview.md)`) plus no `[[` — exercising the previously-dead OKF render branch and `fromDir` plumbing through the render path, as ADR-004/TechSpec N-001 require.
- Verified N-001 (round 1) removal: `grep -rn "WriteMetadataFileWithMode"` matches only round-1 review artifacts, not production code; `WriteMetadataFile` now routes through `topicMetadataForRef(..., "")` (mode omitted via `omitempty`).
- Verified N-002 (round 1) tightening: `firstBodySentence` returns `""` when no `. `/`! `/`? ` boundary exists, `resolveDescription` then emits the empty-description warning; covered by `TestPromoteWarnsWhenSourceBodyHasNoSentenceFallback` (real assertions on empty description + warning).
- Verified N-003 (round 1) removal: the symlink branch in both `Check` (okf.go:177) and `loadConcepts` (okf.go:553) is now `if entry.Type()&fs.ModeSymlink != 0 { return nil }` with no dead `IsDir`/`SkipDir` sub-branch (behavior-preserving — `DirEntry.IsDir()` is always false for a symlink and `WalkDir` does not follow symlinks).
- Verified R-002 (round 1) docs: README.md, AGENTS.md, CLAUDE.md all updated with `topic new --mode wiki|okf`, `kb promote`, `kb okf check`, and `[okf].types`; docs do not overclaim external-dir support.
- Confirmed non-destructive promote and reserved-file exclusion remain correct (source-unchanged assertions and `kb okf check` clean-on-fresh-bundle assertions pass in unit + integration tests).
- `go build ./...` — Success.
- `go test ./internal/okf/... ./internal/vault/... ./internal/topic/... ./internal/cli/... ./internal/config/...` — 295 passed.
- `go test -tags integration ./internal/okf/... ./internal/cli/...` — 181 passed (official ga4/stackoverflow/crypto_bitcoin lenient conformance + end-to-end promote/check).
- Limitation: did not run the full `make verify` (mage/lint/boundaries) in this pass; the round-1 remediation reports `make verify` green (860 tests, zero lint), and direct `go build`/`go test` (unit + integration) over the affected packages pass here.

# Deferred Or Follow-Up

- R-001 structured logging and R-003 external-bundle CLI checking remain deferred follow-ups (non-blocking; confirmed acceptable for the MVP).
- The OKF render branch is now smoke-tested but still emits bundle-root-relative links (`fromDir == ""` at every render site); ADR-004 defers correct `fromDir` plumbing for codebase→OKF to Phase 3. No MVP path renders a codebase topic in OKF mode, so this is the accepted dormant-branch behavior — track for Phase 3.
- `--to <topic>/<subdir>` promotion and nested/`type`-grouped concept directories remain intentional Phase-2 deferrals.
- Consider consolidating the duplicated bundle-walk logic shared by `Check` and `loadConcepts` (entry filtering, exclusion set, posix-rel derivation) before Phase 2 grows the package.
