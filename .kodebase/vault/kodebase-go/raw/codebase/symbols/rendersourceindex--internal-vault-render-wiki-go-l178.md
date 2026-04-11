---
blast_radius: 1
centrality: 0.0524
cyclomatic_complexity: 6
domain: "kodebase-go"
end_line: 228
exported: false
external_reference_count: 1
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: true
language: "go"
loc: 51
outgoing_relation_count: 5
smells:
  - "feature-envy"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/vault/render_wiki.go"
stage: "raw"
start_line: 178
symbol_kind: "function"
symbol_name: "renderSourceIndex"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: renderSourceIndex"
type: "source"
---

# Codebase Symbol: renderSourceIndex

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_wiki.go|internal/vault/render_wiki.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 6
- Long function: true
- Blast radius: 1
- External references: 1
- Centrality: 0.0524
- LOC: 51
- Dead export: false
- Smells: `feature-envy`, `long-function`

## Signature
```text
func renderSourceIndex(topic models.TopicMetadata, articles []starterWikiArticle) models.RenderedDocument {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getwikiconceptpath--internal-vault-pathutils-go-l208]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getwikiindexpath--internal-vault-pathutils-go-l213]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/totopicwikilink--internal-vault-pathutils-go-l228]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sortedmapkeys--internal-vault-render-go-l798]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tosourcewikilink--internal-vault-render-go-l173]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/renderdocuments--internal-vault-render-go-l20]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_wiki.go|internal/vault/render_wiki.go]] via `contains` (syntactic)
