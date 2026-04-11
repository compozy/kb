---
blast_radius: 2
centrality: 0.0552
cyclomatic_complexity: 15
domain: "kodebase-go"
end_line: 690
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: true
language: "go"
loc: 69
outgoing_relation_count: 6
smells:
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/vault/render_wiki.go"
stage: "raw"
start_line: 622
symbol_kind: "function"
symbol_name: "createDeadCodeReportArticle"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: createDeadCodeReportArticle"
type: "source"
---

# Codebase Symbol: createDeadCodeReportArticle

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_wiki.go|internal/vault/render_wiki.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 15
- Long function: true
- Blast radius: 2
- External references: 0
- Centrality: 0.0552
- LOC: 69
- Dead export: false
- Smells: `long-function`

## Signature
```text
func createDeadCodeReportArticle(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getrawfiledocumentpath--internal-vault-pathutils-go-l181]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getrawsymboldocumentpath--internal-vault-pathutils-go-l186]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tosourcewikilink--internal-vault-render-go-l173]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/rendergroupedlinks--internal-vault-render-wiki-go-l964]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sortstrings--internal-vault-render-wiki-go-l1005]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/uniquestrings--internal-vault-render-wiki-go-l991]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/buildstarterwikiarticles--internal-vault-render-wiki-go-l19]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_wiki.go|internal/vault/render_wiki.go]] via `contains` (syntactic)
