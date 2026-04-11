---
blast_radius: 14
centrality: 0.1175
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 599
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/metrics/compute.go"
stage: "raw"
start_line: 597
symbol_kind: "function"
symbol_name: "createFileID"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: createFileID"
type: "source"
---

# Codebase Symbol: createFileID

Source file: [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 14
- External references: 0
- Centrality: 0.1175
- LOC: 3
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func createFileID(filePath string) string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/computemetrics--internal-metrics-compute-go-l23]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/getnodefileid--internal-metrics-compute-go-l519]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]] via `contains` (syntactic)
