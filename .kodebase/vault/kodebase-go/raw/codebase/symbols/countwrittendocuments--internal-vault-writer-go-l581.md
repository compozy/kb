---
blast_radius: 1
centrality: 0.0538
cyclomatic_complexity: 6
domain: "kodebase-go"
end_line: 596
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 16
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/writer.go"
stage: "raw"
start_line: 581
symbol_kind: "function"
symbol_name: "countWrittenDocuments"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: countWrittenDocuments"
type: "source"
---

# Codebase Symbol: countWrittenDocuments

Source file: [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 6
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0538
- LOC: 16
- Dead export: false
- Smells: None

## Signature
```text
func countWrittenDocuments(documents []models.RenderedDocument) WriteVaultResult {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/writevault--internal-vault-writer-go-l53]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]] via `contains` (syntactic)
