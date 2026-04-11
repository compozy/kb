---
blast_radius: 2
centrality: 0.0631
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 284
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 21
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/reader.go"
stage: "raw"
start_line: 264
symbol_kind: "function"
symbol_name: "parseRelations"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: parseRelations"
type: "source"
---

# Codebase Symbol: parseRelations

Source file: [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.0631
- LOC: 21
- Dead export: false
- Smells: None

## Signature
```text
func parseRelations(section string, pattern *regexp.Regexp) []VaultRelation {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parsevaultdocument--internal-vault-reader-go-l223]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]] via `contains` (syntactic)
