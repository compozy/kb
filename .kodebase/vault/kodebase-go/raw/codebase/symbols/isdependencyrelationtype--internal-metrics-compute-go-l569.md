---
blast_radius: 14
centrality: 0.0905
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 572
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 4
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/metrics/compute.go"
stage: "raw"
start_line: 569
symbol_kind: "function"
symbol_name: "isDependencyRelationType"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: isDependencyRelationType"
type: "source"
---

# Codebase Symbol: isDependencyRelationType

Source file: [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 14
- External references: 0
- Centrality: 0.0905
- LOC: 4
- Dead export: false
- Smells: None

## Signature
```text
func isDependencyRelationType(relationType models.RelationType) bool {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/computeapproxcentrality--internal-metrics-compute-go-l261]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/computemetrics--internal-metrics-compute-go-l23]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]] via `contains` (syntactic)
