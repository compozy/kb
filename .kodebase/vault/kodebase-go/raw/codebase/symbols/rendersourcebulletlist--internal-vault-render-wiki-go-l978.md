---
blast_radius: 5
centrality: 0.0732
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 989
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 12
outgoing_relation_count: 2
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/render_wiki.go"
stage: "raw"
start_line: 978
symbol_kind: "function"
symbol_name: "renderSourceBulletList"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: renderSourceBulletList"
type: "source"
---

# Codebase Symbol: renderSourceBulletList

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_wiki.go|internal/vault/render_wiki.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 5
- External references: 0
- Centrality: 0.0732
- LOC: 12
- Dead export: false
- Smells: None

## Signature
```text
func renderSourceBulletList(topic models.TopicMetadata, sources []string, emptyMessage string) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tosourcewikilink--internal-vault-render-go-l173]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/uniquestrings--internal-vault-render-wiki-go-l991]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/createdependencyhotspotsarticle--internal-vault-render-wiki-go-l471]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createdirectorymaparticle--internal-vault-render-wiki-go-l349]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createsymboltaxonomyarticle--internal-vault-render-wiki-go-l405]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_wiki.go|internal/vault/render_wiki.go]] via `contains` (syntactic)
