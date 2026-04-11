---
blast_radius: 1
centrality: 0.0524
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 73
exported: false
external_reference_count: 1
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 9
outgoing_relation_count: 2
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/render_wiki.go"
stage: "raw"
start_line: 65
symbol_kind: "function"
symbol_name: "renderWikiArticle"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: renderWikiArticle"
type: "source"
---

# Codebase Symbol: renderWikiArticle

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_wiki.go|internal/vault/render_wiki.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 1
- External references: 1
- Centrality: 0.0524
- LOC: 9
- Dead export: false
- Smells: None

## Signature
```text
func renderWikiArticle(topic models.TopicMetadata, article starterWikiArticle) models.RenderedDocument {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getwikiconceptpath--internal-vault-pathutils-go-l208]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/makewikifrontmatter--internal-vault-render-wiki-go-l43]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/renderdocuments--internal-vault-render-go-l20]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_wiki.go|internal/vault/render_wiki.go]] via `contains` (syntactic)
