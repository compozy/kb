---
blast_radius: 1
centrality: 0.0524
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 41
exported: false
external_reference_count: 1
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 23
outgoing_relation_count: 10
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/render_wiki.go"
stage: "raw"
start_line: 19
symbol_kind: "function"
symbol_name: "buildStarterWikiArticles"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: buildStarterWikiArticles"
type: "source"
---

# Codebase Symbol: buildStarterWikiArticles

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_wiki.go|internal/vault/render_wiki.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 1
- External references: 1
- Centrality: 0.0524
- LOC: 23
- Dead export: false
- Smells: None

## Signature
```text
func buildStarterWikiArticles(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createcirculardependenciesarticle--internal-vault-render-wiki-go-l865]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createcodebaseoverviewarticle--internal-vault-render-wiki-go-l230]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createcodesmellsarticle--internal-vault-render-wiki-go-l784]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createcomplexityhotspotsarticle--internal-vault-render-wiki-go-l547]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createdeadcodereportarticle--internal-vault-render-wiki-go-l622]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createdependencyhotspotsarticle--internal-vault-render-wiki-go-l471]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createdirectorymaparticle--internal-vault-render-wiki-go-l349]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createhighimpactsymbolsarticle--internal-vault-render-wiki-go-l892]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createmodulehealtharticle--internal-vault-render-wiki-go-l692]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createsymboltaxonomyarticle--internal-vault-render-wiki-go-l405]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/renderdocuments--internal-vault-render-go-l20]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_wiki.go|internal/vault/render_wiki.go]] via `contains` (syntactic)
