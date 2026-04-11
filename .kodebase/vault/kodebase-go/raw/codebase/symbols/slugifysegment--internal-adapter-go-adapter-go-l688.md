---
blast_radius: 12
centrality: 0.1108
cyclomatic_complexity: 9
domain: "kodebase-go"
end_line: 718
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 31
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter.go"
stage: "raw"
start_line: 688
symbol_kind: "function"
symbol_name: "slugifySegment"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: slugifySegment"
type: "source"
---

# Codebase Symbol: slugifySegment

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 9
- Long function: false
- Blast radius: 12
- External references: 0
- Centrality: 0.1108
- LOC: 31
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func slugifySegment(value string) string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/createsymbolid--internal-adapter-go-adapter-go-l677]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] via `contains` (syntactic)
