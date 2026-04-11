---
blast_radius: 8
centrality: 0.0812
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 675
exported: false
external_reference_count: 3
has_smells: false
incoming_relation_count: 5
is_dead_export: false
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter.go"
stage: "raw"
start_line: 673
symbol_kind: "function"
symbol_name: "createExternalID"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: createExternalID"
type: "source"
---

# Codebase Symbol: createExternalID

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 8
- External references: 3
- Centrality: 0.0812
- LOC: 3
- Dead export: false
- Smells: None

## Signature
```text
func createExternalID(source string) string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/extractimports--internal-adapter-go-adapter-go-l264]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/extractrequirebindings--internal-adapter-ts-adapter-go-l644]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/extracttsexports--internal-adapter-ts-adapter-go-l484]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/extracttsimports--internal-adapter-ts-adapter-go-l391]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] via `contains` (syntactic)
