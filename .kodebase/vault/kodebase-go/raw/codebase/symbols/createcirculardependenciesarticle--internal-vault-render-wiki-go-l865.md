---
blast_radius: 2
centrality: 0.0552
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 890
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 26
outgoing_relation_count: 3
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/render_wiki.go"
stage: "raw"
start_line: 865
symbol_kind: "function"
symbol_name: "createCircularDependenciesArticle"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: createCircularDependenciesArticle"
type: "source"
---

# Codebase Symbol: createCircularDependenciesArticle

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_wiki.go|internal/vault/render_wiki.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.0552
- LOC: 26
- Dead export: false
- Smells: None

## Signature
```text
func createCircularDependenciesArticle(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getrawfiledocumentpath--internal-vault-pathutils-go-l181]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tosourcewikilink--internal-vault-render-go-l173]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/uniquestrings--internal-vault-render-wiki-go-l991]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/buildstarterwikiarticles--internal-vault-render-wiki-go-l19]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_wiki.go|internal/vault/render_wiki.go]] via `contains` (syntactic)
