# Workflow Memory

Keep only durable, cross-task context here. Do not duplicate facts that are obvious from the repository, PRD documents, or git history.

## Current State

## Shared Decisions
- `models.SupportedLanguages()` order now includes Java appended at the end: `ts, tsx, js, jsx, go, rust, java` to keep prior ordering stable for existing consumers while exposing Java for downstream adapter selection/help text.
- Java MVP adapter (`internal/adapter/java_adapter.go`) now emits deterministic file outputs, `JAVA_PARSE_ERROR` diagnostics, Java imports as external nodes, and baseline cross-file `references`/`calls` via syntactic import + class/method name matching.
- `internal/generate` default adapter registration order is now `TS -> Go -> Rust -> Java` in both `newRunner()` and `runner.withDefaults()`, preserving deterministic mixed-language adapter selection across normal and fallback paths.
- Java adapter now runs deep-first relation resolution with `semantic` confidence and applies syntactic fallback only for unresolved deep targets, emitting per-file `JAVA_RESOLUTION_FALLBACK` warnings at parse stage.
- Java performance budget enforcement now lives in `internal/generate/generate_integration_test.go` via `TestGenerateIntegrationJavaIngestPerformanceBudget` plus paired benchmarks `BenchmarkGenerateIntegrationGoBaselineDryRun` and `BenchmarkGenerateIntegrationJavaCanonicalDryRun`.
- Java benchmark governance policy is now centralized in `internal/generate/benchmark_policy.go` with canonical profile ordering (`single-module-library`, `spring-service`, `multi-module-enterprise`), fixed repeat count (`3`), and budget (`20%`) shared by integration gate and unit tests.
- Reproducible rollout benchmark execution is now exposed as `make benchmark-java-rollout`, and Phase 3 baseline evidence should be archived in `.compozy/tasks/java-ingest-adapter/_phase3-benchmark-baseline.md`.
- Nested Java types are now represented with qualified symbol names (for example `Outer.Inner`), and Java qualifier parsing preserves dotted ownership chains so deep/fallback resolution can map nested references deterministically.
- Java wildcard imports (`import pkg.*`) now deep-resolve deterministically to package-local top-level class symbols and feed wildcard-aware simple-name lookup for deep call resolution.
- Java import resolution now tracks ambiguous explicit import qualifiers (same simple/type qualifier mapping to multiple FQNs) and treats those call targets as unresolved (`ambiguous-import-class`) in both deep and fallback paths to prevent misleading `calls` edges.
- Java deep call resolution now treats multiple static import candidates for the same unqualified method call as unresolved ambiguity (`ambiguous-static-call-target`) instead of falling back to owner-method resolution.
- Java adapter now parses Gradle/Maven module hints in best-effort mode and narrows ambiguous class-target selection by current module + declared module dependencies, while keeping metadata optional and non-fatal.
- Phase 2 regression gate now standardizes an enterprise tri-module fixture (`shared-a`, `shared-b`, `app` with `app -> shared-b`) across CLI and generate integration tests to assert deterministic nested/wildcard/ambiguity behavior under module metadata.
- Generate parse-stage completion events now emit Java telemetry fields only when Java files are parsed: `java_parse_duration_millis`, `java_files_processed`, `java_resolver_mode`, `java_fallback_count`, and `java_unresolved_count`, keeping non-Java parse payloads contract-compatible.
- Java ingest automation JSON contract is now explicitly documented in `.compozy/tasks/java-ingest-adapter/_automation-json-contract.md`, defining required `codebaseIngestResult`/`GenerationSummary`/`GenerationTimings` keys plus dry-run/full-run value semantics for external consumers.
- Portfolio-scale Java operations now use `.compozy/tasks/java-ingest-adapter/_java-portfolio-adoption-playbook.md` as the baseline runbook, with governance gates, telemetry interpretation, troubleshooting guidance, and CLI-validated command flow.
- Lint workflow now supports Java diagnostics governance via `java-diagnostic-governance` issues driven by `raw/codebase/index/java.md` counters (`java_parse_error_count`, `java_resolution_fallback_count`), with CLI thresholds `--java-max-parse-errors` (default `0`) and `--java-max-fallback-warnings` (default `-1`, disabled).
- Java fallback/module-hint diagnostics are now payload-bounded for scale safety (`entry` + `byte` caps) and emit deterministic `meta:truncated (...)` markers instead of unbounded detail growth.

## Shared Learnings
- Scanner language routing for Java is keyed by `.java` in `supportedLanguage()`; downstream Java adapter tasks should rely on `models.LangJava` coming from scanner grouping (`FilesByLanguage`).
- Deterministic ordering for Java outputs now requires sorting both relation edges and fallback diagnostics to keep repeated runs byte-stable in integration fixtures.
- For local-class lookup, mapping simple names should only be added when the simple name resolves to exactly one FQN in the file; qualified names remain the safe default for nested types.
- Unresolved wildcard package imports now surface as `missing-wildcard-package` in fallback diagnostics, preserving ingest success while making unresolved package scope explicit for follow-up tuning.
- Deterministic behavior alone is insufficient for ambiguity safety; explicit ambiguity classification is required so resolver precedence does not accidentally emit stable-but-incorrect semantic relations.
- Module metadata warnings should be emitted as parse-stage warnings (`JAVA_MODULE_HINT_WARNING`) and must not block relation context indexing; only error diagnostics should suppress file participation in context maps.
- For reproducible performance evidence in regression tasks, the Java budget test should log baseline/java median durations and computed overhead even on PASS runs.
- Java unresolved telemetry can be deterministically derived from `JAVA_RESOLUTION_FALLBACK` diagnostic detail segments, avoiding parser/adapter contract changes while still exposing machine-readable fallback pressure signals.
- Policy values consumed by integration-only tests should be accessed through non-tagged helper functions (for example `canonicalJavaBenchmarkPolicy()`) to avoid lint failures from build-tag-specific constant usage.
- Contract stability for automation is now regression-protected by shared CLI helper assertions (`internal/cli/workflow_test_helpers_test.go`) reused in unit and integration tests, so future payload evolution should update helper key lists and docs together.
- Operational documentation drift is best controlled with paired tests: unit checks for required playbook content plus integration checks that execute the same documented commands (`topic new`, `ingest codebase` dry/full, `lint`).
- Java parse telemetry unresolved counting must ignore `meta:truncated` diagnostic segments so fallback counters remain semantically accurate when diagnostic payload capping is active.

## Open Risks
- `internal/adapter` coverage currently sits just above the threshold (`80.6%` with `go test -tags integration ./internal/adapter -cover`), so unrelated coverage regressions in other adapters can still break template-level coverage expectations.
- Java deep resolution currently relies on repository-local package/import metadata; advanced enterprise classpath scenarios still depend on fallback and may need richer metadata ingestion in future tasks.

## Handoffs
