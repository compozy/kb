---
blast_radius: 1
centrality: 0.0569
cyclomatic_complexity: 16
domain: "kodebase-go"
end_line: 376
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: true
language: "go"
loc: 78
outgoing_relation_count: 14
smells:
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter.go"
stage: "raw"
start_line: 299
symbol_kind: "function"
symbol_name: "parseTSFile"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: parseTSFile"
type: "source"
---

# Codebase Symbol: parseTSFile

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 16
- Long function: true
- Blast radius: 1
- External references: 0
- Centrality: 0.0569
- LOC: 78
- Dead export: false
- Smells: `long-function`

## Signature
```text
func parseTSFile(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createfileid--internal-adapter-go-adapter-go-l669]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractleadingcomment--internal-adapter-go-adapter-go-l595]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/namedchildren--internal-adapter-go-adapter-go-l566]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newparser--internal-adapter-treesitter-go-l31]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/applylocalexportstate--internal-adapter-ts-adapter-go-l1116]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createtsparsediagnostic--internal-adapter-ts-adapter-go-l1253]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createtssymbolmatch--internal-adapter-ts-adapter-go-l766]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractclassmethodsymbols--internal-adapter-ts-adapter-go-l782]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractcommonjsexports--internal-adapter-ts-adapter-go-l583]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractrequirebindings--internal-adapter-ts-adapter-go-l644]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extracttsexports--internal-adapter-ts-adapter-go-l484]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extracttsimports--internal-adapter-ts-adapter-go-l391]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractvariablesymbols--internal-adapter-ts-adapter-go-l805]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/selecttslanguage--internal-adapter-ts-adapter-go-l378]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parsefileswithprogress--internal-adapter-ts-adapter-go-l93]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] via `contains` (syntactic)
