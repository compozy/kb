---
blast_radius: 3
centrality: 0.0661
cyclomatic_complexity: 20
domain: "kodebase-go"
end_line: 1331
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: true
language: "go"
loc: 42
outgoing_relation_count: 1
smells:
  - "feature-envy"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter.go"
stage: "raw"
start_line: 1290
symbol_kind: "function"
symbol_name: "matchCommonJSExportTarget"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: matchCommonJSExportTarget"
type: "source"
---

# Codebase Symbol: matchCommonJSExportTarget

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 20
- Long function: true
- Blast radius: 3
- External references: 0
- Centrality: 0.0661
- LOC: 42
- Dead export: false
- Smells: `feature-envy`, `long-function`

## Signature
```text
func matchCommonJSExportTarget(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/textof--internal-adapter-go-adapter-go-l588]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/extractcommonjsexports--internal-adapter-ts-adapter-go-l583]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] via `contains` (syntactic)
