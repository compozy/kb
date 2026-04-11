---
blast_radius: 3
centrality: 0.0577
cyclomatic_complexity: 25
domain: "kodebase-go"
end_line: 764
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: true
language: "go"
loc: 121
outgoing_relation_count: 7
smells:
  - "feature-envy"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter.go"
stage: "raw"
start_line: 644
symbol_kind: "function"
symbol_name: "extractRequireBindings"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: extractRequireBindings"
type: "source"
---

# Codebase Symbol: extractRequireBindings

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 25
- Long function: true
- Blast radius: 3
- External references: 0
- Centrality: 0.0577
- LOC: 121
- Dead export: false
- Smells: `feature-envy`, `long-function`

## Signature
```text
func extractRequireBindings(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/collectnodesbykind--internal-adapter-go-adapter-go-l577]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createexternalid--internal-adapter-go-adapter-go-l673]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createfileid--internal-adapter-go-adapter-go-l669]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/namedchildren--internal-adapter-go-adapter-go-l566]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/stripquotes--internal-adapter-go-adapter-go-l648]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/textof--internal-adapter-go-adapter-go-l588]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolverelativeimportfile--internal-adapter-ts-adapter-go-l1207]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/extracttsexports--internal-adapter-ts-adapter-go-l484]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/parsetsfile--internal-adapter-ts-adapter-go-l299]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] via `contains` (syntactic)
