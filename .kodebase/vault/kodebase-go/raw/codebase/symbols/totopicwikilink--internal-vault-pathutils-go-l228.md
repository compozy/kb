---
blast_radius: 28
centrality: 0.4837
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 235
exported: true
external_reference_count: 10
has_smells: true
incoming_relation_count: 12
is_dead_export: false
is_long_function: false
language: "go"
loc: 8
outgoing_relation_count: 1
smells:
  - "bottleneck"
  - "high-blast-radius"
source_kind: "codebase-symbol"
source_path: "internal/vault/pathutils.go"
stage: "raw"
start_line: 228
symbol_kind: "function"
symbol_name: "ToTopicWikiLink"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: ToTopicWikiLink"
type: "source"
---

# Codebase Symbol: ToTopicWikiLink

Source file: [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 28
- External references: 10
- Centrality: 0.4837
- LOC: 8
- Dead export: false
- Smells: `bottleneck`, `high-blast-radius`

## Signature
```text
func ToTopicWikiLink(topicSlug, documentPath, label string) string {
```

## Documentation
ToTopicWikiLink formats a topic-scoped Obsidian wiki-link.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/stripmarkdownextension--internal-vault-pathutils-go-l223]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/tosourcewikilink--internal-vault-render-go-l173]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createcodebaseoverviewarticle--internal-vault-render-wiki-go-l230]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createcomplexityhotspotsarticle--internal-vault-render-wiki-go-l547]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createdependencyhotspotsarticle--internal-vault-render-wiki-go-l471]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createdirectorymaparticle--internal-vault-render-wiki-go-l349]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createsymboltaxonomyarticle--internal-vault-render-wiki-go-l405]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderconceptindex--internal-vault-render-wiki-go-l134]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderdashboard--internal-vault-render-wiki-go-l75]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/rendersourceindex--internal-vault-render-wiki-go-l178]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/buildtopicclaude--internal-vault-writer-go-l359]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `exports` (syntactic)
