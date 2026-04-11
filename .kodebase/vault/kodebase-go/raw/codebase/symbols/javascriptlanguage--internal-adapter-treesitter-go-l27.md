---
blast_radius: 3
centrality: 0.0661
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 29
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
start_line: 27
symbol_kind: "function"
symbol_name: "javaScriptLanguage"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: javaScriptLanguage"
type: "source"
---

# Codebase Symbol: javaScriptLanguage

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/treesitter.go|internal/adapter/treesitter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 3
- External references: 1
- Centrality: 0.0661
- LOC: 3
- Dead export: false
- Smells: None

## Signature
```text
func javaScriptLanguage() *tree_sitter.Language {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/selecttslanguage--internal-adapter-ts-adapter-go-l378]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/treesitter.go|internal/adapter/treesitter.go]] via `contains` (syntactic)
