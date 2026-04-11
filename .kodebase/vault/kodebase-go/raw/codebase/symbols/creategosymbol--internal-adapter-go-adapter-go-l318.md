---
blast_radius: 2
centrality: 0.0573
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 352
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 35
outgoing_relation_count: 6
smells:
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter.go"
stage: "raw"
start_line: 318
symbol_kind: "function"
symbol_name: "createGoSymbol"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: createGoSymbol"
type: "source"
---

# Codebase Symbol: createGoSymbol

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.0573
- LOC: 35
- Dead export: false
- Smells: None

## Signature
```text
func createGoSymbol(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/computecyclomaticcomplexity--internal-adapter-go-adapter-go-l465]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createsymbolid--internal-adapter-go-adapter-go-l677]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractgodoccomment--internal-adapter-go-adapter-go-l381]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/formatgosignature--internal-adapter-go-adapter-go-l372]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isgoexported--internal-adapter-go-adapter-go-l662]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvegosymbolname--internal-adapter-go-adapter-go-l354]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parsegofile--internal-adapter-go-adapter-go-l161]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] via `contains` (syntactic)
