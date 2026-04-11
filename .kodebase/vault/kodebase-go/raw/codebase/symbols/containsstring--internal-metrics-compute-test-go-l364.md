---
blast_radius: 2
centrality: 0.068
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 372
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 9
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/metrics/compute_test.go"
stage: "raw"
start_line: 364
symbol_kind: "function"
symbol_name: "containsString"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: containsString"
type: "source"
---

# Codebase Symbol: containsString

Source file: [[kodebase-go/raw/codebase/files/internal/metrics/compute_test.go|internal/metrics/compute_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.068
- LOC: 9
- Dead export: false
- Smells: None

## Signature
```text
func containsString(values []string, target string) bool {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testcomputemetricsflagsdeadexports--internal-metrics-compute-test-go-l151]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testcomputemetricsflagslongfunctions--internal-metrics-compute-test-go-l171]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/metrics/compute_test.go|internal/metrics/compute_test.go]] via `contains` (syntactic)
