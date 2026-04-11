---
blast_radius: 2
centrality: 0.0542
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 389
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 12
outgoing_relation_count: 3
smells:
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter.go"
stage: "raw"
start_line: 378
symbol_kind: "function"
symbol_name: "selectTSLanguage"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: selectTSLanguage"
type: "source"
---

# Codebase Symbol: selectTSLanguage

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.0542
- LOC: 12
- Dead export: false
- Smells: `feature-envy`

## Signature
```text
func selectTSLanguage(language models.SupportedLanguage) *tree_sitter.Language {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/javascriptlanguage--internal-adapter-treesitter-go-l27]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tsxlanguage--internal-adapter-treesitter-go-l23]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/typescriptlanguage--internal-adapter-treesitter-go-l19]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parsetsfile--internal-adapter-ts-adapter-go-l299]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] via `contains` (syntactic)
