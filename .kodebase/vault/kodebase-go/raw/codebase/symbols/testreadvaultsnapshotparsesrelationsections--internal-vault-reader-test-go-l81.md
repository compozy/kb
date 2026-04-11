---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 8
domain: "kodebase-go"
end_line: 135
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 55
outgoing_relation_count: 2
smells:
  - "dead-export"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/vault/reader_test.go"
stage: "raw"
start_line: 81
symbol_kind: "function"
symbol_name: "TestReadVaultSnapshotParsesRelationSections"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestReadVaultSnapshotParsesRelationSections"
type: "source"
---

# Codebase Symbol: TestReadVaultSnapshotParsesRelationSections

Source file: [[kodebase-go/raw/codebase/files/internal/vault/reader_test.go|internal/vault/reader_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 8
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 55
- Dead export: true
- Smells: `dead-export`, `long-function`

## Signature
```text
func TestReadVaultSnapshotParsesRelationSections(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createresolvedvault--internal-vault-reader-test-go-l247]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writemarkdowndocument--internal-vault-reader-test-go-l265]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/reader_test.go|internal/vault/reader_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/reader_test.go|internal/vault/reader_test.go]] via `exports` (syntactic)
