---
blast_radius: 6
centrality: 0.1046
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 215
exported: true
external_reference_count: 4
has_smells: true
incoming_relation_count: 6
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
start_line: 213
symbol_kind: "function"
symbol_name: "GetWikiIndexPath"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: GetWikiIndexPath"
type: "source"
---

# Codebase Symbol: GetWikiIndexPath

Source file: [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 6
- External references: 4
- Centrality: 0.1046
- LOC: 3
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func GetWikiIndexPath(indexTitle string) string {
```

## Documentation
GetWikiIndexPath derives the vault document path for a generated wiki index page.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/renderconceptindex--internal-vault-render-wiki-go-l134]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderdashboard--internal-vault-render-wiki-go-l75]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/rendersourceindex--internal-vault-render-wiki-go-l178]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/buildtopicclaude--internal-vault-writer-go-l359]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `exports` (syntactic)
