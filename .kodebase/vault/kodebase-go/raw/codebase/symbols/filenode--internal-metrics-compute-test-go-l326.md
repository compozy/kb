---
blast_radius: 10
centrality: 0.1521
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 334
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 11
is_dead_export: false
is_long_function: false
language: "go"
loc: 9
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/metrics/compute_test.go"
stage: "raw"
start_line: 326
symbol_kind: "function"
symbol_name: "fileNode"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: fileNode"
type: "source"
---

# Codebase Symbol: fileNode

Source file: [[kodebase-go/raw/codebase/files/internal/metrics/compute_test.go|internal/metrics/compute_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 10
- External references: 0
- Centrality: 0.1521
- LOC: 9
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func fileNode(filePath string, symbolIDs ...string) models.GraphFile {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testcomputemetricsaggregatesdirectorymetrics--internal-metrics-compute-test-go-l276]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricscomputesbalancedinstability--internal-metrics-compute-test-go-l129]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricscountsblastradiusacrossdependents--internal-metrics-compute-test-go-l49]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricscountsfileafferentcoupling--internal-metrics-compute-test-go-l103]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricscountsfileefferentcoupling--internal-metrics-compute-test-go-l81]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricsdetectscirculardependencies--internal-metrics-compute-test-go-l191]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricsflagsdeadexports--internal-metrics-compute-test-go-l151]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricsflagslongfunctions--internal-metrics-compute-test-go-l171]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricsreturnsnocirculardependenciesforacyclicgraph--internal-metrics-compute-test-go-l255]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricssinglesymbolhaszeroblastradius--internal-metrics-compute-test-go-l32]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/metrics/compute_test.go|internal/metrics/compute_test.go]] via `contains` (syntactic)
