# TechSpec: Java Codebase Ingest Adapter

## Executive Summary

This specification adds Java support to the existing codebase ingest pipeline by extending language discovery, parser adapter orchestration, and graph generation without changing the CLI workflow. The design follows current adapter patterns (TS/Go/Rust), introducing a Java adapter with deep relation resolution as the primary strategy and syntactic fallback as the resilience mechanism.

The core trade-off is higher relation quality vs. runtime cost. We accept additional complexity and bounded overhead to deliver stronger cross-file dependency accuracy for Java repositories, while enforcing a hard performance budget (<=20% runtime increase on representative fixtures) and preserving user-facing ingest reliability.

## System Architecture

### Component Overview

- **Language Registry (`internal/models`)**
  - Add `LangJava` to supported language lists.
  - Keeps language ordering deterministic for adapter selection and CLI language help.

- **Workspace Scanner (`internal/scanner`)**
  - Extend extension mapping to recognize `.java`.
  - Reuses existing ignore/include/exclude behavior and file grouping by language.

- **Tree-sitter Binding Layer (`internal/adapter/treesitter.go`)**
  - Add `javaLanguage()` provider and parser initialization plumbing.
  - Keeps one parser lifecycle per adapter run, consistent with existing adapters.

- **Java Adapter (`internal/adapter/java_adapter.go`)**
  - Parses Java files, extracts files/symbols/relations/diagnostics.
  - Primary mode: deep relation resolution (package/module/classpath-aware when available).
  - Fallback mode: syntactic relation extraction for unresolved deep paths.

- **Generation Orchestrator (`internal/generate`)**
  - Register `JavaAdapter{}` in adapter list and default runner fallback.
  - No stage contract changes (`scan -> select_adapters -> parse -> normalize -> metrics -> render -> write`).

- **Vault/Inspect Consumers**
  - Consume normalized graph-derived documents unchanged.
  - Java output remains language-agnostic at rendering and inspect layers.

Data flow remains unchanged: Java files enter at scan, route to Java adapter in parse stage, merge in graph normalization, and flow through existing render/write/inspect behavior.

## Implementation Design

### Core Interfaces

```go
type JavaAdapter struct{}

func (JavaAdapter) Supports(language models.SupportedLanguage) bool {
	return language == models.LangJava
}

func (adapter JavaAdapter) ParseFiles(
	files []models.ScannedSourceFile,
	rootPath string,
) ([]models.ParsedFile, error) {
	return adapter.ParseFilesWithProgress(files, rootPath, nil)
}
```

```go
type javaResolver interface {
	Resolve(
		file models.ScannedSourceFile,
		symbols []models.SymbolNode,
		imports []javaImportRef,
	) (resolvedRelations []models.RelationEdge, unresolved []javaUnresolvedRef)
}
```

### Data Models

- **New/Extended Existing Models**
  - `models.SupportedLanguage`: add `LangJava`.
  - `models.ScannedSourceFile`: no schema change.
  - `models.ParsedFile`: reused unchanged.

- **Java Adapter Internal Models (package-private)**
  - `javaParsedFile`:
    - `file models.GraphFile`
    - `symbols []models.SymbolNode`
    - `externalNodes map[string]models.ExternalNode`
    - `relations []models.RelationEdge`
    - `diagnostics []models.StructuredDiagnostic`
  - `javaImportRef`:
    - `importPath string`
    - `isStatic bool`
    - `isWildcard bool`
    - `alias string`
  - `javaUnresolvedRef`:
    - `sourceSymbolID string`
    - `targetHint string`
    - `reason string`

- **Diagnostic Code**
  - `JAVA_PARSE_ERROR` for parse failures.
  - `JAVA_RESOLUTION_FALLBACK` warning diagnostic when deep resolution falls back.

### API Endpoints

No API endpoint changes are required. This feature extends internal CLI pipeline behavior only.

## Integration Points

- **Tree-sitter Java grammar binding**
  - Add Go module dependency for Java grammar binding compatible with current Tree-sitter runtime.
  - Integration is internal and parser-scoped; no external service call is introduced.

- **Classpath/module metadata usage**
  - Deep resolver may consume repository-local metadata (e.g., module manifests) when present.
  - Failure to resolve metadata never blocks ingest; fallback is automatic.

## Impact Analysis

| Component | Impact Type | Description and Risk | Required Action |
|-----------|-------------|----------------------|-----------------|
| `internal/models/models.go` | modified | Adds `LangJava`; low risk, broad compile impact | Update constants and language lists |
| `internal/scanner/scanner.go` | modified | Adds `.java` mapping; low risk | Extend `supportedLanguage()` |
| `internal/adapter/treesitter.go` | modified | Adds Java language loader; medium risk due binding compatibility | Add `javaLanguage()` and tests |
| `internal/adapter/java_adapter.go` | new | Core parse/relation logic; high complexity risk | Implement adapter + fallback diagnostics |
| `internal/generate/generate.go` | modified | Registers Java adapter; low risk | Add adapter in runner defaults |
| `go.mod` / `go.sum` | modified | Adds tree-sitter Java dependency; medium risk | Add dependency via `go get` |
| `internal/vault/*` | unchanged | Consumes normalized outputs; low risk | No code change expected |
| `internal/cli/*` | unchanged/indirect | Language help updates via model list; low risk | Validate help text through tests |

## Testing Approach

### Unit Tests

- Add/update unit tests for:
  - `SupportedLanguages()` and `SupportedLanguageNames()` include `java`.
  - Scanner maps `.java` correctly.
  - Tree-sitter Java language initializes and parses trivial source.
  - Java adapter symbol extraction for package/class/interface/enum/record/method.
  - Fallback diagnostic emission when deep resolver cannot resolve.

- Edge cases:
  - Static/wildcard imports.
  - Nested classes and overloaded methods.
  - Missing module/classpath metadata.

### Integration Tests

- Add `java_adapter_integration_test.go` with representative fixtures:
  - Single-module Java project.
  - Spring-style package layout.
  - Multi-module repository with cross-module references.
- Validate:
  - symbol count and kinds,
  - imports/external nodes,
  - `calls/references` quality in common cross-file paths,
  - fallback behavior for unresolved deep relations.

### Benchmark and E2E Validation

- Add benchmark scenario for Java ingest on large fixture(s) and compare with baseline.
- Enforce gate: Java-enabled ingest runtime must stay within <=20% overhead.
- Add CLI E2E integration:
  - `kb ingest codebase` against Java multi-module fixture,
  - verify summary language detection, output artifact presence, and inspect compatibility.
- Run acceptance benchmarks over the canonical pilot profiles: single-module library, Spring-style service, and multi-module enterprise-style repository.

### Verification Gate

- Mandatory final validation: `make verify`.

## Development Sequencing

### Build Order

1. **Language model extension** (`internal/models`) - no dependencies.
2. **Scanner extension** (`internal/scanner`) - depends on step 1.
3. **Tree-sitter Java dependency and binding function** (`go.mod`, `internal/adapter/treesitter.go`) - depends on step 1.
4. **Java adapter skeleton with parse diagnostics** (`internal/adapter/java_adapter.go`) - depends on steps 1 and 3.
5. **Deep resolver + syntactic fallback path** (inside Java adapter) - depends on step 4.
6. **Generator registration** (`internal/generate/generate.go`) - depends on step 4.
7. **Unit test updates for models/scanner/treesitter** - depends on steps 1-3.
8. **Java adapter integration tests + CLI E2E + benchmark** - depends on steps 5 and 6.
9. **Full verification and performance gate validation** - depends on steps 7 and 8.

### Technical Dependencies

- Compatible Tree-sitter Java Go binding version.
- Stable Java test fixtures (single-module and multi-module) under `testdata`.
- Benchmark baseline definition (command flags and fixture set) agreed before final acceptance.
- Pilot feedback collection mechanism to measure confidence score and rollout readiness.

## Monitoring and Observability

- **Key metrics**
  - parse stage duration for Java files,
  - total ingest duration delta vs baseline,
  - unresolved deep-resolution count and ratio.

- **Log events / structured fields**
  - `stage=parse`, `language=java`, `files_processed`,
  - `resolver_mode=deep|fallback`,
  - `fallback_count`, `unresolved_count`.

- **Alerting thresholds**
  - performance budget breach (>20% over baseline),
  - unresolved ratio spikes above expected fixture thresholds.
  - pilot confidence readiness breach (<80% responses at >=4/5).

## Technical Considerations

### Key Decisions

- **Decision**: Deep relation resolution with automatic syntactic fallback.  
  **Rationale**: Balances relation quality and ingest resilience.  
  **Trade-off**: More adapter complexity and diagnostics handling.  
  **Alternatives rejected**: syntactic-only, strict fail-on-unresolved.

- **Decision**: 20% performance overhead budget with hybrid cache strategy.  
  **Rationale**: Gives explicit non-functional acceptance while keeping MVP feasible.  
  **Trade-off**: Persistent cache benefits deferred post-MVP.  
  **Alternatives rejected**: no numeric gate, immediate persistent cache.

- **Decision**: Require benchmark + CLI E2E in addition to unit/integration tests.  
  **Rationale**: Feature risk spans correctness, performance, and UX flow.  
  **Trade-off**: Longer test cycle and maintenance cost.  
  **Alternatives rejected**: unit/integration only, benchmark without E2E.

### Known Risks

- **Resolver drift risk**: Deep resolver misses enterprise-specific layouts.
  - Mitigation: fallback + diagnostics + fixture expansion in Phase 2.
- **Performance risk**: large repositories cause parser/resolver overhead spikes.
  - Mitigation: bounded resolver passes, benchmark gate, staged optimization.
- **Dependency compatibility risk**: Tree-sitter Java binding version mismatch.
  - Mitigation: lock compatible versions and include parser initialization tests.

## Architecture Decision Records

- [ADR-001: Adopt a balanced MVP strategy for Java codebase ingest](adrs/adr-001.md) — Product-level direction favors balanced early value over extremes.
- [ADR-002: Use deep Java relation resolution with safe syntactic fallback](adrs/adr-002.md) — Deep resolution is primary; fallback preserves ingest reliability.
- [ADR-003: Enforce 20% ingest performance budget with hybrid caching strategy](adrs/adr-003.md) — Numeric performance gate plus in-memory-first caching design.
- [ADR-004: Require unit, integration, benchmark, and CLI E2E validation for Java ingest](adrs/adr-004.md) — Release quality requires correctness, runtime, and workflow evidence.
- [ADR-005: Define MVP governance acceptance gates and pilot corpus](adrs/adr-005.md) — Formalizes performance threshold, pilot set, and confidence gate for rollout.
- [ADR-006: Close Java ingest MVP rollout using available pilot evidence](adrs/adr-006.md) — Captures MVP rollout closure decision and deferred non-blocking governance evidence.
