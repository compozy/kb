---
blast_radius: 3
centrality: 0.0867
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 354
exported: false
external_reference_count: 1
has_smells: false
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 11
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/writer_test.go"
stage: "raw"
start_line: 344
symbol_kind: "function"
symbol_name: "assertFileExists"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: assertFileExists"
type: "source"
---

# Codebase Symbol: assertFileExists

Source file: [[kodebase-go/raw/codebase/files/internal/vault/writer_test.go|internal/vault/writer_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 3
- External references: 1
- Centrality: 0.0867
- LOC: 11
- Dead export: false
- Smells: None

## Signature
```text
func assertFileExists(t *testing.T, filePath string) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testwritevaultintegrationpersistsfullrenderedvault--internal-vault-writer-integration-test-go-l14]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testwritevaultcreatestopicskeletonandmanagedfiles--internal-vault-writer-test-go-l16]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testwritevaultremovesstalemanagedwikiconceptsonly--internal-vault-writer-test-go-l187]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/writer_test.go|internal/vault/writer_test.go]] via `contains` (syntactic)
