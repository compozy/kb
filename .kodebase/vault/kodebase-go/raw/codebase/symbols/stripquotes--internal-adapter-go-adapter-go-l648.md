---
blast_radius: 8
centrality: 0.0812
cyclomatic_complexity: 6
domain: "kodebase-go"
end_line: 660
exported: false
external_reference_count: 3
has_smells: false
incoming_relation_count: 5
is_dead_export: false
is_long_function: false
language: "go"
loc: 13
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter.go"
stage: "raw"
start_line: 648
symbol_kind: "function"
symbol_name: "stripQuotes"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: stripQuotes"
type: "source"
---

# Codebase Symbol: stripQuotes

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 6
- Long function: false
- Blast radius: 8
- External references: 3
- Centrality: 0.0812
- LOC: 13
- Dead export: false
- Smells: None

## Signature
```text
func stripQuotes(value string) string {
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
