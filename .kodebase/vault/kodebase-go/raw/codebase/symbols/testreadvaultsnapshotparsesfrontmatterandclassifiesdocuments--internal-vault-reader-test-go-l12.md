---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 9
domain: "kodebase-go"
end_line: 79
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 68
outgoing_relation_count: 2
smells:
  - "dead-export"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/vault/reader_test.go"
stage: "raw"
start_line: 12
symbol_kind: "function"
symbol_name: "TestReadVaultSnapshotParsesFrontmatterAndClassifiesDocuments"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestReadVaultSnapshotParsesFrontmatterAndClassifiesDocuments"
type: "source"
---

# Codebase Symbol: TestReadVaultSnapshotParsesFrontmatterAndClassifiesDocuments

Source file: [[kodebase-go/raw/codebase/files/internal/vault/reader_test.go|internal/vault/reader_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 9
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 68
- Dead export: true
- Smells: `dead-export`, `long-function`

## Signature
```text
func TestReadVaultSnapshotParsesFrontmatterAndClassifiesDocuments(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createresolvedvault--internal-vault-reader-test-go-l247]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writemarkdowndocument--internal-vault-reader-test-go-l265]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/reader_test.go|internal/vault/reader_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/reader_test.go|internal/vault/reader_test.go]] via `exports` (syntactic)
