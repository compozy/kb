---
blast_radius: 6
centrality: 0.114
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 780
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 6
is_dead_export: false
is_long_function: false
language: "go"
loc: 15
outgoing_relation_count: 3
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter.go"
stage: "raw"
start_line: 766
symbol_kind: "function"
symbol_name: "createTSSymbolMatch"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: createTSSymbolMatch"
type: "source"
---

# Codebase Symbol: createTSSymbolMatch

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 6
- External references: 0
- Centrality: 0.114
- LOC: 15
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func createTSSymbolMatch(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/clonestringset--internal-adapter-ts-adapter-go-l1380]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/collecttscalltargets--internal-adapter-ts-adapter-go-l1018]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createtssymbol--internal-adapter-ts-adapter-go-l826]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/extractclassmethodsymbols--internal-adapter-ts-adapter-go-l782]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/extractcommonjsexports--internal-adapter-ts-adapter-go-l583]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/extracttsexports--internal-adapter-ts-adapter-go-l484]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/extractvariablesymbols--internal-adapter-ts-adapter-go-l805]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/parsetsfile--internal-adapter-ts-adapter-go-l299]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] via `contains` (syntactic)
