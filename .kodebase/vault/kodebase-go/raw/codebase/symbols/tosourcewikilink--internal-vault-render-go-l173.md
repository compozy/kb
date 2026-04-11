---
blast_radius: 23
centrality: 0.4056
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 175
exported: false
external_reference_count: 13
has_smells: true
incoming_relation_count: 19
is_dead_export: false
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 1
smells:
  - "bottleneck"
  - "feature-envy"
  - "high-blast-radius"
source_kind: "codebase-symbol"
source_path: "internal/vault/render.go"
stage: "raw"
start_line: 173
symbol_kind: "function"
symbol_name: "toSourceWikiLink"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: toSourceWikiLink"
type: "source"
---

# Codebase Symbol: toSourceWikiLink

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 23
- External references: 13
- Centrality: 0.4056
- LOC: 3
- Dead export: false
- Smells: `bottleneck`, `feature-envy`, `high-blast-radius`

## Signature
```text
func toSourceWikiLink(topic models.TopicMetadata, relativePath, label string) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/totopicwikilink--internal-vault-pathutils-go-l228]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/linkfornode--internal-vault-render-go-l199]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderrawdirectoryindex--internal-vault-render-go-l524]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderrawfiledocument--internal-vault-render-go-l312]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderrawlanguageindex--internal-vault-render-go-l600]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderrawsymboldocument--internal-vault-render-go-l449]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createcirculardependenciesarticle--internal-vault-render-wiki-go-l865]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createcodebaseoverviewarticle--internal-vault-render-wiki-go-l230]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createcodesmellsarticle--internal-vault-render-wiki-go-l784]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createcomplexityhotspotsarticle--internal-vault-render-wiki-go-l547]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createdeadcodereportarticle--internal-vault-render-wiki-go-l622]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createdependencyhotspotsarticle--internal-vault-render-wiki-go-l471]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createdirectorymaparticle--internal-vault-render-wiki-go-l349]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createhighimpactsymbolsarticle--internal-vault-render-wiki-go-l892]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createmodulehealtharticle--internal-vault-render-wiki-go-l692]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createsymboltaxonomyarticle--internal-vault-render-wiki-go-l405]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/makewikifrontmatter--internal-vault-render-wiki-go-l43]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/rendersourcebulletlist--internal-vault-render-wiki-go-l978]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/rendersourceindex--internal-vault-render-wiki-go-l178]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]] via `contains` (syntactic)
