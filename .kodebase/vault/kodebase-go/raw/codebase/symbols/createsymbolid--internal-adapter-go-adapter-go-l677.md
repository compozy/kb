---
blast_radius: 11
centrality: 0.0706
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 686
exported: false
external_reference_count: 1
has_smells: false
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 10
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter.go"
stage: "raw"
start_line: 677
symbol_kind: "function"
symbol_name: "createSymbolID"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: createSymbolID"
type: "source"
---

# Codebase Symbol: createSymbolID

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 11
- External references: 1
- Centrality: 0.0706
- LOC: 10
- Dead export: false
- Smells: None

## Signature
```text
func createSymbolID(symbol models.SymbolNode) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/slugifysegment--internal-adapter-go-adapter-go-l688]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/creategosymbol--internal-adapter-go-adapter-go-l318]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createtssymbol--internal-adapter-ts-adapter-go-l826]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] via `contains` (syntactic)
