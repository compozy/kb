---
blast_radius: 2
centrality: 0.0795
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 330
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 10
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/writer_test.go"
stage: "raw"
start_line: 321
symbol_kind: "function"
symbol_name: "readFile"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: readFile"
type: "source"
---

# Codebase Symbol: readFile

Source file: [[kodebase-go/raw/codebase/files/internal/vault/writer_test.go|internal/vault/writer_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.0795
- LOC: 10
- Dead export: false
- Smells: None

## Signature
```text
func readFile(t *testing.T, filePath string) string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testwritevaultcreatesclaudemanifestandappendonlylog--internal-vault-writer-test-go-l102]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testwritevaultcreatestopicskeletonandmanagedfiles--internal-vault-writer-test-go-l16]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/writer_test.go|internal/vault/writer_test.go]] via `contains` (syntactic)
