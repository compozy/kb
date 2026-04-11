---
blast_radius: 2
centrality: 0.0542
cyclomatic_complexity: 18
domain: "kodebase-go"
end_line: 581
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: true
language: "go"
loc: 98
outgoing_relation_count: 13
smells:
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter.go"
stage: "raw"
start_line: 484
symbol_kind: "function"
symbol_name: "extractTSExports"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: extractTSExports"
type: "source"
---

# Codebase Symbol: extractTSExports

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 18
- Long function: true
- Blast radius: 2
- External references: 0
- Centrality: 0.0542
- LOC: 98
- Dead export: false
- Smells: `long-function`

## Signature
```text
func extractTSExports(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createexternalid--internal-adapter-go-adapter-go-l673]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createfileid--internal-adapter-go-adapter-go-l669]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/namedchildren--internal-adapter-go-adapter-go-l566]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/stripquotes--internal-adapter-go-adapter-go-l648]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/textof--internal-adapter-go-adapter-go-l588]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createtssymbolmatch--internal-adapter-ts-adapter-go-l766]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/declarationexportnames--internal-adapter-ts-adapter-go-l1265]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractclassmethodsymbols--internal-adapter-ts-adapter-go-l782]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractrequirebindings--internal-adapter-ts-adapter-go-l644]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractvariablesymbols--internal-adapter-ts-adapter-go-l805]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/findchildbykind--internal-adapter-ts-adapter-go-l1333]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isdefaultexport--internal-adapter-ts-adapter-go-l1285]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolverelativeimportfile--internal-adapter-ts-adapter-go-l1207]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parsetsfile--internal-adapter-ts-adapter-go-l299]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] via `contains` (syntactic)
