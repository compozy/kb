---
blast_radius: 1
centrality: 0.0651
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 319
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 10
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/writer_test.go"
stage: "raw"
start_line: 310
symbol_kind: "function"
symbol_name: "filterOutDocument"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: filterOutDocument"
type: "source"
---

# Codebase Symbol: filterOutDocument

Source file: [[kodebase-go/raw/codebase/files/internal/vault/writer_test.go|internal/vault/writer_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0651
- LOC: 10
- Dead export: false
- Smells: None

## Signature
```text
func filterOutDocument(documents []models.RenderedDocument, relativePath string) []models.RenderedDocument {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testwritevaultremovesstalemanagedwikiconceptsonly--internal-vault-writer-test-go-l187]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/writer_test.go|internal/vault/writer_test.go]] via `contains` (syntactic)
