---
blast_radius: 13
centrality: 0.1611
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 210
exported: true
external_reference_count: 10
has_smells: true
incoming_relation_count: 12
is_dead_export: false
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/vault/pathutils.go"
stage: "raw"
start_line: 208
symbol_kind: "function"
symbol_name: "GetWikiConceptPath"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: GetWikiConceptPath"
type: "source"
---

# Codebase Symbol: GetWikiConceptPath

Source file: [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 13
- External references: 10
- Centrality: 0.1611
- LOC: 3
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func GetWikiConceptPath(articleTitle string) string {
```

## Documentation
GetWikiConceptPath derives the vault document path for a generated wiki concept article.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/createcodebaseoverviewarticle--internal-vault-render-wiki-go-l230]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createcomplexityhotspotsarticle--internal-vault-render-wiki-go-l547]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createdependencyhotspotsarticle--internal-vault-render-wiki-go-l471]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createdirectorymaparticle--internal-vault-render-wiki-go-l349]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createsymboltaxonomyarticle--internal-vault-render-wiki-go-l405]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderconceptindex--internal-vault-render-wiki-go-l134]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderdashboard--internal-vault-render-wiki-go-l75]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/rendersourceindex--internal-vault-render-wiki-go-l178]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderwikiarticle--internal-vault-render-wiki-go-l65]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/buildtopicclaude--internal-vault-writer-go-l359]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `exports` (syntactic)
