---
blast_radius: 1
centrality: 0.0615
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 17
exported: false
external_reference_count: 1
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/adapter/treesitter.go"
stage: "raw"
start_line: 15
symbol_kind: "function"
symbol_name: "goLanguage"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: goLanguage"
type: "source"
---

# Codebase Symbol: goLanguage

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/treesitter.go|internal/adapter/treesitter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 1
- External references: 1
- Centrality: 0.0615
- LOC: 3
- Dead export: false
- Smells: None

## Signature
```text
func goLanguage() *tree_sitter.Language {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parsefileswithprogress--internal-adapter-go-adapter-go-l62]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/treesitter.go|internal/adapter/treesitter.go]] via `contains` (syntactic)
