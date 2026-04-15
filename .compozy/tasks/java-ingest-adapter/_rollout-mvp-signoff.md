# Java Ingest Adapter — MVP Rollout Sign-off

## Decision

MVP rollout is **closed** for the Java ingest adapter initiative based on available execution evidence, technical verification, and stakeholder directive to finalize using the current pilot run.

Date: 2026-04-15

## Evidence Sources

1. Pilot ingest execution on legacy Java repository:
   - Command:
     - `<KB_BINARY_PATH>/kb ingest codebase <LEGACY_REPO_PATH> --topic java-lang`
   - Working directory:
     - `<PROJECT_WORKDIR>`
   - Captured output artifact:
     - `<EXECUTION_LOG_ARTIFACT>`
2. Integration benchmark evidence from implementation memory:
   - `memory/task_06.md`
3. Post-ingest structural lint:
   - `kb lint java-lang --vault <PILOT_VAULT_PATH> --format json`
4. Repository verification gate:
   - `make verify`

## Pilot Execution Summary (Legacy Repo)

- `detectedLanguages`: `["java"]`
- `selectedAdapters`: `["adapter.JavaAdapter"]`
- `filesScanned`: `2075`
- `filesParsed`: `2075`
- `symbolsExtracted`: `39362`
- `relationsEmitted`: `109932`
- `rawDocumentsWritten`: `41487`
- `wikiDocumentsWritten`: `10`
- `indexDocumentsWritten`: `3`
- `totalMillis`: `12635`

## Diagnostic Summary

- Parse-stage warning diagnostics were emitted with `JAVA_RESOLUTION_FALLBACK` as expected by design.
- Count in captured run: `1817` fallback warnings.
- No `JAVA_PARSE_ERROR` entries were found.
- No diagnostics with `severity: "error"` were found in the captured run output.

Interpretation: deep resolver fallback occurred in many files (notably external/classpath-unresolved references), but ingest completed successfully and produced the full output corpus.

## Governance Gate Assessment

### Gate 1 — Performance threshold (`<=20%`)

- Benchmark evidence (from `memory/task_06.md`):
  - Java: `3793232 ns/op`
  - Go baseline: `3388442 ns/op`
  - Overhead: `~11.95%`
- Result: **PASS** (within threshold).

### Gate 2 — Canonical pilot corpus (3 profiles)

- Evidence available in this sign-off:
  - Executed pilot: enterprise-style legacy Java repository (`<LEGACY_REPO_PATH>`).
- Additional profile evidence (single-module and Spring-style pilot repositories) is not attached in this sign-off packet.
- Result: **WAIVED BY ROLLOUT DECISION CONTEXT** (see ADR-006).

### Gate 3 — Confidence target (`>=80%` with `>=4/5`)

- Formal survey dataset is not attached in this sign-off packet.
- Result: **WAIVED BY ROLLOUT DECISION CONTEXT** (see ADR-006).

## Operational Validation

- Topic lint after ingest:
  - Output: `[]`
  - Result: **PASS**
- Project verification gate:
  - `make verify` result: **PASS**
  - Includes fmt/lint/tests/build/boundaries success.

## Phase 2 Regression Validation (Task 11)

- Adapter integration regression:
  - `go test -tags integration ./internal/adapter -run "TestJavaAdapterPhase2EnterpriseScenarioRegression"`
  - Result: **PASS**
  - Validates combined nested + wildcard + ambiguity + metadata-assisted multi-module behavior, with predictable fallback detail for unresolved wildcard packages.
- CLI E2E regression:
  - `go test -tags integration ./internal/cli -run "TestCLIIntegrationScaffoldIngestJavaWorkspaceCodebase"`
  - Result: **PASS**
  - Confirms enterprise Java ingest summary stability (`FilesScanned=6`, `FilesParsed=6`, `SelectedAdapters=[adapter.JavaAdapter]`), generated artifacts, and clean `kb lint`.
- Generate integration regression:
  - `go test -tags integration ./internal/generate -run "TestGenerateIntegrationBuildsVaultFromJavaPhase2Workspace"`
  - Result: **PASS**
  - Confirms enterprise tri-module fixture output creation in `raw/codebase/files/*`.
- Performance budget rerun:
  - `go test -tags integration ./internal/generate -run "TestGenerateIntegrationJavaIngestPerformanceBudget" -v`
  - Latest sampled output: `baseline=4.840792ms java=4.17975ms overhead=-13.66% budget=20.00%`
  - Result: **PASS** (within ADR-003 budget).
- Coverage check:
  - `go test -tags integration ./internal/adapter -cover`
  - Result: `coverage: 80.7% of statements` (**PASS**, >= 80%).

## Residual Risks Accepted at Rollout Closure

- High fallback volume indicates unresolved classpath/external symbol scenarios remain common in large enterprise repositories.
- Relationship fidelity for advanced Java patterns should continue to improve in Phase 2.

## Follow-up Actions (Phase 2 Planning Inputs)

1. Reduce fallback warning volume for enterprise classpath patterns.
2. Execute explicit pilot runs for single-module and Spring-style repositories and attach evidence to Phase 2 baseline.
3. Introduce structured confidence collection for broader rollout governance.

## Automation Contract Reference

- Java ingest automation consumers should use the stabilized contract notes in:
  - `.compozy/tasks/java-ingest-adapter/_automation-json-contract.md`
- That document defines required keys, dry-run/full-run semantics, compatibility policy, and non-guaranteed fields.

## Approval Notes

Rollout closure was finalized under explicit request to complete sign-off using available evidence from the previous execution and current validation artifacts.
