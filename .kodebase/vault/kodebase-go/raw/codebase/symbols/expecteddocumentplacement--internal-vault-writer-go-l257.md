---
blast_radius: 3
centrality: 0.0788
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 268
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 12
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/writer.go"
stage: "raw"
start_line: 257
symbol_kind: "function"
symbol_name: "expectedDocumentPlacement"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: expectedDocumentPlacement"
type: "source"
---

# Codebase Symbol: expectedDocumentPlacement

Source file: [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 3
- External references: 0
- Centrality: 0.0788
- LOC: 12
- Dead export: false
- Smells: None

## Signature
```text
func expectedDocumentPlacement(kind models.DocumentKind) (models.ManagedArea, string, error) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/validaterendereddocument--internal-vault-writer-go-l200]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]] via `contains` (syntactic)
