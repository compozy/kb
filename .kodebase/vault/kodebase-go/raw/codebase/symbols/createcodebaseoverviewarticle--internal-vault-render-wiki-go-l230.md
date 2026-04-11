---
blast_radius: 2
centrality: 0.0552
cyclomatic_complexity: 13
domain: "kodebase-go"
end_line: 347
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: true
language: "go"
loc: 118
outgoing_relation_count: 8
smells:
  - "feature-envy"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/vault/render_wiki.go"
stage: "raw"
start_line: 230
symbol_kind: "function"
symbol_name: "createCodebaseOverviewArticle"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: createCodebaseOverviewArticle"
type: "source"
---

# Codebase Symbol: createCodebaseOverviewArticle

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_wiki.go|internal/vault/render_wiki.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 13
- Long function: true
- Blast radius: 2
- External references: 0
- Centrality: 0.0552
- LOC: 118
- Dead export: false
- Smells: `feature-envy`, `long-function`

## Signature
```text
func createCodebaseOverviewArticle(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getrawdirectoryindexpath--internal-vault-pathutils-go-l193]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getrawfiledocumentpath--internal-vault-pathutils-go-l181]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getrawlanguageindexpath--internal-vault-pathutils-go-l203]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getwikiconceptpath--internal-vault-pathutils-go-l208]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/totopicwikilink--internal-vault-pathutils-go-l228]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sortedmapkeys--internal-vault-render-go-l798]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tosourcewikilink--internal-vault-render-go-l173]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/uniquestrings--internal-vault-render-wiki-go-l991]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/buildstarterwikiarticles--internal-vault-render-wiki-go-l19]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_wiki.go|internal/vault/render_wiki.go]] via `contains` (syntactic)
