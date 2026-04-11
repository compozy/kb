---
blast_radius: 7
centrality: 0.0831
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 857
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 32
outgoing_relation_count: 6
smells:
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter.go"
stage: "raw"
start_line: 826
symbol_kind: "function"
symbol_name: "createTSSymbol"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: createTSSymbol"
type: "source"
---

# Codebase Symbol: createTSSymbol

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 7
- External references: 0
- Centrality: 0.0831
- LOC: 32
- Dead export: false
- Smells: None

## Signature
```text
func createTSSymbol(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createsymbolid--internal-adapter-go-adapter-go-l677]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/computetscyclomaticcomplexity--internal-adapter-ts-adapter-go-l1069]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extracttsdoccomment--internal-adapter-ts-adapter-go-l1004]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/formattssignature--internal-adapter-ts-adapter-go-l929]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/gettssymbolkind--internal-adapter-ts-adapter-go-l859]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolvetssymbolname--internal-adapter-ts-adapter-go-l880]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/createtssymbolmatch--internal-adapter-ts-adapter-go-l766]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] via `contains` (syntactic)
