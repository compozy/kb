---
blast_radius: 3
centrality: 0.0577
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 824
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 20
outgoing_relation_count: 3
smells:
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter.go"
stage: "raw"
start_line: 805
symbol_kind: "function"
symbol_name: "extractVariableSymbols"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: extractVariableSymbols"
type: "source"
---

# Codebase Symbol: extractVariableSymbols

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 3
- External references: 0
- Centrality: 0.0577
- LOC: 20
- Dead export: false
- Smells: None

## Signature
```text
func extractVariableSymbols(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/collectnodesbykind--internal-adapter-go-adapter-go-l577]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createtssymbolmatch--internal-adapter-ts-adapter-go-l766]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvetsvariablename--internal-adapter-ts-adapter-go-l916]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/extracttsexports--internal-adapter-ts-adapter-go-l484]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/parsetsfile--internal-adapter-ts-adapter-go-l299]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] via `contains` (syntactic)
