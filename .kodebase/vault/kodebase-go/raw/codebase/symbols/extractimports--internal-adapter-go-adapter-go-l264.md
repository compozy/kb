---
blast_radius: 2
centrality: 0.0573
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 300
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 37
outgoing_relation_count: 4
smells:
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter.go"
stage: "raw"
start_line: 264
symbol_kind: "function"
symbol_name: "extractImports"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: extractImports"
type: "source"
---

# Codebase Symbol: extractImports

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.0573
- LOC: 37
- Dead export: false
- Smells: None

## Signature
```text
func extractImports(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/collectnodesbykind--internal-adapter-go-adapter-go-l577]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createexternalid--internal-adapter-go-adapter-go-l673]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/stripquotes--internal-adapter-go-adapter-go-l648]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/textof--internal-adapter-go-adapter-go-l588]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parsegofile--internal-adapter-go-adapter-go-l161]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] via `contains` (syntactic)
