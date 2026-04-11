---
blast_radius: 0
centrality: 0.0507
domain: "kodebase-go"
end_line: 31
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
language: "go"
loc: 7
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/vault/writer.go"
stage: "raw"
start_line: 25
symbol_kind: "struct"
symbol_name: "WriteVaultOptions"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "struct"
title: "Codebase Symbol: WriteVaultOptions"
type: "source"
---

# Codebase Symbol: WriteVaultOptions

Source file: [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]]

## Kind
`struct`

## Static Analysis
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 7
- Dead export: true
- Smells: `dead-export`

## Signature
```text
WriteVaultOptions struct {
```

## Documentation
WriteVaultOptions bundles the rendered vault payload written to disk.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]] via `exports` (syntactic)
