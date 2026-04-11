---
blast_radius: 1
centrality: 0.0538
cyclomatic_complexity: 9
domain: "kodebase-go"
end_line: 518
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 30
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/writer.go"
stage: "raw"
start_line: 489
symbol_kind: "function"
symbol_name: "removeManagedWikiConcepts"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: removeManagedWikiConcepts"
type: "source"
---

# Codebase Symbol: removeManagedWikiConcepts

Source file: [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 9
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0538
- LOC: 30
- Dead export: false
- Smells: None

## Signature
```text
func removeManagedWikiConcepts(topicPath string) error {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/hasmanagedgenerator--internal-vault-writer-go-l520]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/writevault--internal-vault-writer-go-l53]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]] via `contains` (syntactic)
