## Java Canonical Benchmark Corpus

This corpus defines the canonical Java repository profiles used by the rollout runtime governance gate (`<=20%` overhead, median over `3` repeated runs):

1. `single-module-library`
2. `spring-service`
3. `multi-module-enterprise`

The integration suite materializes these fixtures with deterministic generators in `internal/generate/generate_integration_test.go`.

Run the reproducible gate command:

```bash
make benchmark-java-rollout
```

This command executes the deterministic performance-budget integration test plus the canonical dry-run benchmarks for archive/compare workflows.
