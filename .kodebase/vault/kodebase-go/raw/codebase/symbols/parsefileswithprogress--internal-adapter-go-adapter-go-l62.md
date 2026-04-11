---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 20
domain: "kodebase-go"
end_line: 159
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 98
outgoing_relation_count: 4
smells:
  - "dead-export"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter.go"
stage: "raw"
start_line: 62
symbol_kind: "method"
symbol_name: "ParseFilesWithProgress"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "method"
title: "Codebase Symbol: ParseFilesWithProgress"
type: "source"
---

# Codebase Symbol: ParseFilesWithProgress

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]]

## Kind
`method`

## Static Analysis
- Cyclomatic complexity: 20
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 98
- Dead export: true
- Smells: `dead-export`, `long-function`

## Signature
```text
func (adapter GoAdapter) ParseFilesWithProgress(
```

## Documentation
ParseFilesWithProgress parses Go files and reports one progress tick per file.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsegofile--internal-adapter-go-adapter-go-l161]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sortedexternalnodes--internal-adapter-go-adapter-go-l532]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/golanguage--internal-adapter-treesitter-go-l15]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newparser--internal-adapter-treesitter-go-l31]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] via `exports` (syntactic)
