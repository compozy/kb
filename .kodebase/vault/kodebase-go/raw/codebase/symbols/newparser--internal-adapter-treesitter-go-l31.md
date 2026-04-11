---
blast_radius: 4
centrality: 0.1081
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 43
exported: false
external_reference_count: 3
has_smells: true
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 13
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/adapter/treesitter.go"
stage: "raw"
start_line: 31
symbol_kind: "function"
symbol_name: "newParser"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: newParser"
type: "source"
---

# Codebase Symbol: newParser

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/treesitter.go|internal/adapter/treesitter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 4
- External references: 3
- Centrality: 0.1081
- LOC: 13
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func newParser(language *tree_sitter.Language) (*tree_sitter.Parser, error) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parsefileswithprogress--internal-adapter-go-adapter-go-l62]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testnewparserrejectsnillanguage--internal-adapter-treesitter-test-go-l120]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/parsetsfile--internal-adapter-ts-adapter-go-l299]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/treesitter.go|internal/adapter/treesitter.go]] via `contains` (syntactic)
