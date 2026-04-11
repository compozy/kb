---
blast_radius: 13
centrality: 0.0635
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 604
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 4
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/metrics/compute.go"
stage: "raw"
start_line: 601
symbol_kind: "function"
symbol_name: "containsKey"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: containsKey"
type: "source"
---

# Codebase Symbol: containsKey

Source file: [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 13
- External references: 0
- Centrality: 0.0635
- LOC: 4
- Dead export: false
- Smells: None

## Signature
```text
func containsKey[K comparable](items map[K]struct{}, key K) bool {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/computemetrics--internal-metrics-compute-go-l23]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]] via `contains` (syntactic)
