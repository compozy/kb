---
blast_radius: 29
centrality: 0.462
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 225
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: true
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 0
smells:
  - "bottleneck"
  - "dead-export"
  - "high-blast-radius"
source_kind: "codebase-symbol"
source_path: "internal/vault/pathutils.go"
stage: "raw"
start_line: 223
symbol_kind: "function"
symbol_name: "StripMarkdownExtension"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: StripMarkdownExtension"
type: "source"
---

# Codebase Symbol: StripMarkdownExtension

Source file: [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 29
- External references: 0
- Centrality: 0.462
- LOC: 3
- Dead export: true
- Smells: `bottleneck`, `dead-export`, `high-blast-radius`

## Signature
```text
func StripMarkdownExtension(documentPath string) string {
```

## Documentation
StripMarkdownExtension removes a trailing .md extension from a document path.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/totopicwikilink--internal-vault-pathutils-go-l228]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `exports` (syntactic)
