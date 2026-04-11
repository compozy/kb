---
blast_radius: 1
centrality: 0.0579
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 342
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 11
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/writer_test.go"
stage: "raw"
start_line: 332
symbol_kind: "function"
symbol_name: "assertDirExists"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: assertDirExists"
type: "source"
---

# Codebase Symbol: assertDirExists

Source file: [[kodebase-go/raw/codebase/files/internal/vault/writer_test.go|internal/vault/writer_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0579
- LOC: 11
- Dead export: false
- Smells: None

## Signature
```text
func assertDirExists(t *testing.T, directoryPath string) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testwritevaultcreatestopicskeletonandmanagedfiles--internal-vault-writer-test-go-l16]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/writer_test.go|internal/vault/writer_test.go]] via `contains` (syntactic)
