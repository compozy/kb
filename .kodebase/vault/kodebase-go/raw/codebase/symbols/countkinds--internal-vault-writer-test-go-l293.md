---
blast_radius: 2
centrality: 0.0723
cyclomatic_complexity: 6
domain: "kodebase-go"
end_line: 308
exported: false
external_reference_count: 1
has_smells: false
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 16
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/writer_test.go"
stage: "raw"
start_line: 293
symbol_kind: "function"
symbol_name: "countKinds"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: countKinds"
type: "source"
---

# Codebase Symbol: countKinds

Source file: [[kodebase-go/raw/codebase/files/internal/vault/writer_test.go|internal/vault/writer_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 6
- Long function: false
- Blast radius: 2
- External references: 1
- Centrality: 0.0723
- LOC: 16
- Dead export: false
- Smells: None

## Signature
```text
func countKinds(documents []models.RenderedDocument) vault.WriteVaultResult {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testwritevaultintegrationpersistsfullrenderedvault--internal-vault-writer-integration-test-go-l14]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testwritevaultcreatestopicskeletonandmanagedfiles--internal-vault-writer-test-go-l16]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/writer_test.go|internal/vault/writer_test.go]] via `contains` (syntactic)
