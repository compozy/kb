---
blast_radius: 0
centrality: 0.0507
domain: "kodebase-go"
end_line: 27
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 1
is_dead_export: false
language: "go"
loc: 1
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/generate/generate.go"
stage: "raw"
start_line: 27
symbol_kind: "type"
symbol_name: "writeVaultFunc"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "type"
title: "Codebase Symbol: writeVaultFunc"
type: "source"
---

# Codebase Symbol: writeVaultFunc

Source file: [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]]

## Kind
`type`

## Static Analysis
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 1
- Dead export: false
- Smells: None

## Signature
```text
writeVaultFunc func(ctx context.Context, options vault.WriteVaultOptions) (vault.WriteVaultResult, error)
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]] via `contains` (syntactic)
