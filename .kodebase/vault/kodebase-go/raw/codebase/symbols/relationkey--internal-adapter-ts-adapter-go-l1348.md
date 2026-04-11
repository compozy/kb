---
blast_radius: 2
centrality: 0.1053
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 1355
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 8
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter.go"
stage: "raw"
start_line: 1348
symbol_kind: "function"
symbol_name: "relationKey"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: relationKey"
type: "source"
---

# Codebase Symbol: relationKey

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.1053
- LOC: 8
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func relationKey(relation models.RelationEdge) string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parsefileswithprogress--internal-adapter-ts-adapter-go-l93]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/pushuniquerelation--internal-adapter-ts-adapter-go-l1357]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] via `contains` (syntactic)
