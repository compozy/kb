# Java Ingest Adapter — Phase 3 Benchmark Baseline

Date: 2026-04-15

## Reproducible Run Policy

- Canonical corpus profiles:
  - `single-module-library`
  - `spring-service`
  - `multi-module-enterprise`
- Runtime gate:
  - median runtime overhead `<=20%` versus Go baseline
  - `3` repeated dry-run samples per profile
- Fixed execution command:

```bash
make benchmark-java-rollout
```

## Governance Gate Evidence (Median-Based)

Command:

```bash
go test -tags integration ./internal/generate -run "TestGenerateIntegrationJavaIngestPerformanceBudget" -v
```

Result: PASS

| Profile | Baseline Median | Java Median | Overhead | Budget |
| --- | --- | --- | --- | --- |
| `single-module-library` | `5.373458ms` | `3.317875ms` | `-38.25%` | `20.00%` |
| `spring-service` | `5.373458ms` | `4.238666ms` | `-21.12%` | `20.00%` |
| `multi-module-enterprise` | `5.373458ms` | `4.034ms` | `-24.93%` | `20.00%` |

## Benchmark Snapshot (Archive-Friendly)

Command:

```bash
make benchmark-java-rollout
```

```
BenchmarkGenerateIntegrationGoBaselineDryRun-10                                  368   3291831 ns/op  2902660 B/op  38368 allocs/op
BenchmarkGenerateIntegrationJavaCanonicalDryRun/single-module-library-10         453   2641736 ns/op  3206188 B/op  36649 allocs/op
BenchmarkGenerateIntegrationJavaCanonicalDryRun/spring-service-10                310   4035944 ns/op  4310281 B/op  46691 allocs/op
BenchmarkGenerateIntegrationJavaCanonicalDryRun/multi-module-enterprise-10       308   3846089 ns/op  4116548 B/op  46696 allocs/op
```

## Evidence Capture Format

For future governance comparisons, archive each run with:

1. command line used (`make benchmark-java-rollout`)
2. commit SHA
3. machine profile (`goos`, `goarch`, CPU)
4. median gate table (baseline/java/overhead/budget per profile)
5. benchmark snapshot (`ns/op`, `B/op`, `allocs/op` per profile)
