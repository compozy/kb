---
afferent_coupling: 0
domain: "kodebase-go"
efferent_coupling: 1
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 1
is_god_file: false
is_orphan_file: true
language: "go"
outgoing_relation_count: 10
smells:
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/adapter/treesitter_test.go"
stage: "raw"
symbol_count: 4
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/adapter/treesitter_test.go"
type: "source"
---

# Codebase File: internal/adapter/treesitter_test.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 0
- Efferent coupling: 1
- Instability: 1
- Entry point: false
- Circular dependency: false
- Smells: `orphan-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/adapter--internal-adapter-treesitter-test-go-l1|adapter (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testlanguagesinitialize--internal-adapter-treesitter-test-go-l10|TestLanguagesInitialize (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testparsersparsetrivialsources--internal-adapter-treesitter-test-go-l52|TestParsersParseTrivialSources (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testnewparserrejectsnillanguage--internal-adapter-treesitter-test-go-l120|TestNewParserRejectsNilLanguage (function)]] · exported=true

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/adapter--internal-adapter-treesitter-test-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testlanguagesinitialize--internal-adapter-treesitter-test-go-l10]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnewparserrejectsnillanguage--internal-adapter-treesitter-test-go-l120]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testparsersparsetrivialsources--internal-adapter-treesitter-test-go-l52]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testlanguagesinitialize--internal-adapter-treesitter-test-go-l10]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnewparserrejectsnillanguage--internal-adapter-treesitter-test-go-l120]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testparsersparsetrivialsources--internal-adapter-treesitter-test-go-l52]]
- `imports` (syntactic) -> `errors`
- `imports` (syntactic) -> `tree_sitter (github.com/tree-sitter/go-tree-sitter)`
- `imports` (syntactic) -> `testing`

## Backlinks
None
