# Java Portfolio Adoption Playbook

## Purpose

Run Java ingest safely across large repository portfolios with repeatable governance, observable execution, and automation-safe outputs.

## Operating Scope

- Portfolio scale: many repositories, mixed module layouts, recurring ingest cadence.
- Personas: platform engineering, architecture modernization, and governance operators.
- Baseline assumptions: Java ingest adapter is enabled and teams already use `kb ingest codebase`.

## Recommended Portfolio Flow

### 1) Discover and prepare

Goal: establish inventory, pilot profile mapping, and execution plan before writing artifacts.

Required checks:

- Classify each target repository into one canonical profile:
  - single-module Java library
  - Spring-style service repository
  - multi-module enterprise-style repository
- Define topic slug conventions (`<portfolio>-<repo>` or equivalent) to keep outputs deterministic.
- Confirm vault location strategy (`--vault`) for central governance runs.

Suggested command (one-time topic bootstrap per repository):

```bash
kb topic new <topic-slug> "<topic-title>" <domain> --vault <portfolio-vault>
```

### 2) Dry-run ingest gate

Goal: validate scan/adapters/summary contract before writing KB documents.

Command:

```bash
kb ingest codebase <repo-path> \
  --topic <topic-slug> \
  --vault <portfolio-vault> \
  --progress never \
  --log-format json \
  --dry-run
```

Expected dry-run semantics:

- `summary.dryRun = true`
- `summary.rawDocumentsWritten = 0`
- `summary.wikiDocumentsWritten = 0`
- `summary.indexDocumentsWritten = 0`
- `summary.detectedLanguages` includes `java`
- parse-stage telemetry appears on stderr JSON events (when Java files are processed):
  - `java_parse_duration_millis`
  - `java_files_processed`
  - `java_resolver_mode`
  - `java_fallback_count`
  - `java_unresolved_count`

### 3) Full ingest execution

Goal: persist codebase artifacts and topic indexes for operational use.

Command:

```bash
kb ingest codebase <repo-path> \
  --topic <topic-slug> \
  --vault <portfolio-vault> \
  --progress never \
  --log-format json
```

Expected full-run semantics:

- `summary.dryRun = false`
- `summary.rawDocumentsWritten > 0` when source files are discovered
- `sourceType = "codebase-file"`
- `summary.selectedAdapters` includes `adapter.JavaAdapter` for Java-only repos

### 4) Post-ingest quality checks

Goal: ensure generated topic quality before rollout sign-off.

Command:

```bash
kb lint <topic-slug> --vault <portfolio-vault> --format json
```

Pass condition:

- JSON output is `[]` for a clean topic, or all returned issues are triaged and tracked before broader adoption.

## Governance Checkpoints

## Gate A - Performance budget

Reference: ADR-005, ADR-006, rollout sign-off.

- Threshold: median ingest runtime overhead `<= 20%` vs baseline.
- Measurement policy:
  - same flags for baseline and Java runs
  - 3 repeated runs
  - canonical profiles covered in aggregate

Evidence template:

| Field | Value |
| --- | --- |
| profile | `<single-module|spring-service|multi-module-enterprise>` |
| baseline_median_ms | `<number>` |
| java_median_ms | `<number>` |
| overhead_percent | `<number>` |
| within_budget | `<true|false>` |

## Gate B - Canonical pilot coverage

- Required repository profiles:
  - single-module library
  - Spring-style service
  - multi-module enterprise
- Each profile requires successful dry-run + full-run + lint evidence.

## Gate C - Confidence readiness

- Rollout readiness target: `>= 80%` of participants report confidence `>= 4/5`.
- No unresolved critical workflow blockers at sign-off.

## Telemetry and Diagnostics Interpretation

Use parse-stage JSON events (`--log-format json`) and summary diagnostics to assess ingest health:

- `java_files_processed`: confirms Java files were parsed.
- `java_parse_duration_millis`: parse stage cost for Java workload.
- `java_resolver_mode`: `deep` or `fallback`; monitor shifts over time.
- `java_fallback_count`: number of fallback situations. Rising trends signal metadata/classpath gaps.
- `java_unresolved_count`: unresolved relation pressure after fallback handling.

Diagnostic code guidance:

- `JAVA_PARSE_ERROR`:
  - severity error; parse failures that reduce usable graph coverage.
  - action: block rollout for affected repositories until parser issues are triaged.
- `JAVA_RESOLUTION_FALLBACK`:
  - warning path from deep resolver to syntactic fallback.
  - action: ingestion can proceed, but high volume requires governance attention and follow-up backlog.

## Troubleshooting Matrix

| Scenario | Signal | Likely cause | Operator action |
| --- | --- | --- | --- |
| High fallback volume | `java_fallback_count` high, many `JAVA_RESOLUTION_FALLBACK` warnings | Incomplete module/classpath hints or ambiguous imports | Continue ingest, classify pattern, prioritize metadata and resolver tuning in next cycle |
| High unresolved count | `java_unresolved_count` rising vs baseline | Enterprise dependency topology not fully represented | Compare with previous baseline, add profile-specific fixture, schedule resolver hardening |
| Parse errors present | `JAVA_PARSE_ERROR` diagnostics | Unsupported syntax or parser mismatch in repository subset | Treat as blocking for affected repo, isolate failing files, add regression fixture |
| Budget breach | overhead > 20% median | Large parse/resolution cost on portfolio subset | Pause broad rollout for that profile, capture benchmark evidence, optimize before expanding |

Escalation trigger recommendation:

- Trigger governance review when `java_fallback_count` or `java_unresolved_count` trends materially upward for the same repository profile across repeated runs.

## Automation Contract Usage

Reference: `_automation-json-contract.md`.

Automation consumers must rely on these required keys:

- Top-level (`codebaseIngestResult`):
  - `topic`
  - `sourceType`
  - `filePath`
  - `title`
  - `summary`
- Summary (`GenerationSummary`):
  - `command`
  - `rootPath`
  - `vaultPath`
  - `topicPath`
  - `topicSlug`
  - `dryRun`
  - `detectedLanguages`
  - `selectedAdapters`
  - `filesScanned`
  - `filesParsed`
  - `filesSkipped`
  - `symbolsExtracted`
  - `relationsEmitted`
  - `rawDocumentsWritten`
  - `wikiDocumentsWritten`
  - `indexDocumentsWritten`
  - `timings`
  - `diagnostics`
- Timings (`GenerationTimings`):
  - `scanMillis`
  - `selectAdaptersMillis`
  - `parseMillis`
  - `normalizeMillis`
  - `metricsMillis`
  - `renderMillis`
  - `writeMillis`
  - `totalMillis`

Automation policy:

- Treat required keys and documented semantics as stable.
- Allow additive optional fields without breaking pipelines.
- Do not hard-couple to array order, exact diagnostic text/count, or absolute timing values.

## Evidence Collection Checklist

Use this minimal packet for each portfolio wave:

1. Dry-run JSON output archived.
2. Full-run JSON output archived.
3. Parse-stage telemetry event sample archived.
4. Lint JSON output archived.
5. Performance comparison sheet (`<= 20%` gate).
6. Confidence summary (`>= 80%` at `>= 4/5`) and blocker status.

## Operational Notes for Phase 3+

- Keep this playbook as the single operating baseline for Java ingest portfolio rollouts.
- Feed recurring fallback/unresolved hotspots into Phase 3/4 hardening backlogs.
- Re-run governance gates after significant resolver, benchmark-policy, or contract changes.
