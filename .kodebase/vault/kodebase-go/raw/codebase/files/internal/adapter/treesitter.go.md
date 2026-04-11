---
afferent_coupling: 3
domain: "kodebase-go"
efferent_coupling: 0
has_circular_dependency: false
has_smells: false
incoming_relation_count: 0
instability: 0
is_god_file: false
is_orphan_file: false
language: "go"
outgoing_relation_count: 12
smells:
source_kind: "codebase-file"
source_path: "internal/adapter/treesitter.go"
stage: "raw"
symbol_count: 6
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/adapter/treesitter.go"
type: "source"
---

# Codebase File: internal/adapter/treesitter.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 3
- Efferent coupling: 0
- Instability: 0
- Entry point: false
- Circular dependency: false
- Smells: None

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/adapter--internal-adapter-treesitter-go-l1|adapter (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/golanguage--internal-adapter-treesitter-go-l15|goLanguage (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/typescriptlanguage--internal-adapter-treesitter-go-l19|typeScriptLanguage (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/tsxlanguage--internal-adapter-treesitter-go-l23|tsxLanguage (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/javascriptlanguage--internal-adapter-treesitter-go-l27|javaScriptLanguage (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/newparser--internal-adapter-treesitter-go-l31|newParser (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/adapter--internal-adapter-treesitter-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/golanguage--internal-adapter-treesitter-go-l15]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/javascriptlanguage--internal-adapter-treesitter-go-l27]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newparser--internal-adapter-treesitter-go-l31]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tsxlanguage--internal-adapter-treesitter-go-l23]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/typescriptlanguage--internal-adapter-treesitter-go-l19]]
- `imports` (syntactic) -> `errors`
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `tree_sitter (github.com/tree-sitter/go-tree-sitter)`
- `imports` (syntactic) -> `tree_sitter_go (github.com/tree-sitter/tree-sitter-go/bindings/go)`
- `imports` (syntactic) -> `tree_sitter_javascript (github.com/tree-sitter/tree-sitter-javascript/bindings/go)`
- `imports` (syntactic) -> `tree_sitter_typescript (github.com/tree-sitter/tree-sitter-typescript/bindings/go)`

## Backlinks
None
