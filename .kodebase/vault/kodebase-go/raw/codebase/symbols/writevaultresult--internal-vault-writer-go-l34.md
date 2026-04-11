---
blast_radius: 0
centrality: 0.0507
domain: "kodebase-go"
end_line: 38
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
source_path: "internal/vault/writer.go"
stage: "raw"
start_line: 34
symbol_kind: "struct"
symbol_name: "WriteVaultResult"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "struct"
title: "Codebase Symbol: WriteVaultResult"
type: "source"
---

# Codebase Symbol: WriteVaultResult

Source file: [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]]

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
WriteVaultResult struct {
```

## Documentation
WriteVaultResult reports how many managed markdown documents were written.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]] via `exports` (syntactic)
