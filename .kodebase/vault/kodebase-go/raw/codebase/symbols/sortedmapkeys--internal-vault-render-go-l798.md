---
blast_radius: 17
centrality: 0.2684
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 805
exported: false
external_reference_count: 8
has_smells: true
incoming_relation_count: 11
is_dead_export: false
is_long_function: false
language: "go"
loc: 8
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/vault/render.go"
stage: "raw"
start_line: 798
symbol_kind: "function"
symbol_name: "sortedMapKeys"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: sortedMapKeys"
type: "source"
---

# Codebase Symbol: sortedMapKeys

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 17
- External references: 8
- Centrality: 0.2684
- LOC: 8
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func sortedMapKeys[V any](values map[string]V) []string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/renderdocuments--internal-vault-render-go-l20]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderfrontmatter--internal-vault-render-go-l130]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderyamlvalue--internal-vault-render-base-go-l295]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createcodebaseoverviewarticle--internal-vault-render-wiki-go-l230]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createcodesmellsarticle--internal-vault-render-wiki-go-l784]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createdirectorymaparticle--internal-vault-render-wiki-go-l349]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createmodulehealtharticle--internal-vault-render-wiki-go-l692]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createsymboltaxonomyarticle--internal-vault-render-wiki-go-l405]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/rendergroupedlinks--internal-vault-render-wiki-go-l964]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/rendersourceindex--internal-vault-render-wiki-go-l178]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]] via `contains` (syntactic)
