---
blast_radius: 5
centrality: 0.0707
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 205
exported: true
external_reference_count: 3
has_smells: false
incoming_relation_count: 5
is_dead_export: false
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/pathutils.go"
stage: "raw"
start_line: 203
symbol_kind: "function"
symbol_name: "GetRawLanguageIndexPath"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: GetRawLanguageIndexPath"
type: "source"
---

# Codebase Symbol: GetRawLanguageIndexPath

Source file: [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 5
- External references: 3
- Centrality: 0.0707
- LOC: 3
- Dead export: false
- Smells: None

## Signature
```text
func GetRawLanguageIndexPath(language string) string {
```

## Documentation
GetRawLanguageIndexPath derives the vault document path for a raw language index.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/renderrawlanguageindex--internal-vault-render-go-l600]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createcodebaseoverviewarticle--internal-vault-render-wiki-go-l230]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createsymboltaxonomyarticle--internal-vault-render-wiki-go-l405]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `exports` (syntactic)
