---
blast_radius: 0
centrality: 0.0507
domain: "kodebase-go"
end_line: 21
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
language: "go"
loc: 5
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/vault/reader.go"
stage: "raw"
start_line: 17
symbol_kind: "struct"
symbol_name: "VaultRelation"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "struct"
title: "Codebase Symbol: VaultRelation"
type: "source"
---

# Codebase Symbol: VaultRelation

Source file: [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]]

## Kind
`struct`

## Static Analysis
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 5
- Dead export: true
- Smells: `dead-export`

## Signature
```text
VaultRelation struct {
```

## Documentation
VaultRelation describes a parsed relation link in a vault markdown document.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]] via `exports` (syntactic)
