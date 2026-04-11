---
blast_radius: 0
centrality: 0.0507
domain: "kodebase-go"
end_line: 25
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
source_path: "internal/vault/query.go"
stage: "raw"
start_line: 21
symbol_kind: "struct"
symbol_name: "VaultQueryOptions"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "struct"
title: "Codebase Symbol: VaultQueryOptions"
type: "source"
---

# Codebase Symbol: VaultQueryOptions

Source file: [[kodebase-go/raw/codebase/files/internal/vault/query.go|internal/vault/query.go]]

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
VaultQueryOptions struct {
```

## Documentation
VaultQueryOptions control vault and topic discovery from CLI-style flags.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/query.go|internal/vault/query.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/query.go|internal/vault/query.go]] via `exports` (syntactic)
