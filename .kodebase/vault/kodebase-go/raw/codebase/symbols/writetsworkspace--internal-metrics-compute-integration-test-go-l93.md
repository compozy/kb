---
blast_radius: 1
centrality: 0.0651
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 120
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 28
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/metrics/compute_integration_test.go"
stage: "raw"
start_line: 93
symbol_kind: "function"
symbol_name: "writeTSWorkspace"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: writeTSWorkspace"
type: "source"
---

# Codebase Symbol: writeTSWorkspace

Source file: [[kodebase-go/raw/codebase/files/internal/metrics/compute_integration_test.go|internal/metrics/compute_integration_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0651
- LOC: 28
- Dead export: false
- Smells: None

## Signature
```text
func writeTSWorkspace(t *testing.T, rootDir string, sources map[string]string) []models.ScannedSourceFile {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testcomputemetricsintegrationonmultidirectoryproject--internal-metrics-compute-integration-test-go-l16]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/metrics/compute_integration_test.go|internal/metrics/compute_integration_test.go]] via `contains` (syntactic)
