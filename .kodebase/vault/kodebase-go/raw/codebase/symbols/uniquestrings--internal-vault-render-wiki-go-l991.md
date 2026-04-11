---
blast_radius: 13
centrality: 0.1528
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 1003
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 10
is_dead_export: false
is_long_function: false
language: "go"
loc: 13
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/vault/render_wiki.go"
stage: "raw"
start_line: 991
symbol_kind: "function"
symbol_name: "uniqueStrings"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: uniqueStrings"
type: "source"
---

# Codebase Symbol: uniqueStrings

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_wiki.go|internal/vault/render_wiki.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 13
- External references: 0
- Centrality: 0.1528
- LOC: 13
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func uniqueStrings(values []string) []string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/createcirculardependenciesarticle--internal-vault-render-wiki-go-l865]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createcodebaseoverviewarticle--internal-vault-render-wiki-go-l230]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createcodesmellsarticle--internal-vault-render-wiki-go-l784]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createcomplexityhotspotsarticle--internal-vault-render-wiki-go-l547]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createdeadcodereportarticle--internal-vault-render-wiki-go-l622]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createhighimpactsymbolsarticle--internal-vault-render-wiki-go-l892]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createmodulehealtharticle--internal-vault-render-wiki-go-l692]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createsymboltaxonomyarticle--internal-vault-render-wiki-go-l405]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/rendersourcebulletlist--internal-vault-render-wiki-go-l978]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_wiki.go|internal/vault/render_wiki.go]] via `contains` (syntactic)
