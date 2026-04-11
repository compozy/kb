---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 10
domain: "kodebase-go"
end_line: 58
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 47
outgoing_relation_count: 1
smells:
  - "dead-export"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/vault/reader_integration_test.go"
stage: "raw"
start_line: 12
symbol_kind: "function"
symbol_name: "TestReadVaultSnapshotRoundTripsWriterOutput"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestReadVaultSnapshotRoundTripsWriterOutput"
type: "source"
---

# Codebase Symbol: TestReadVaultSnapshotRoundTripsWriterOutput

Source file: [[kodebase-go/raw/codebase/files/internal/vault/reader_integration_test.go|internal/vault/reader_integration_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 10
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 47
- Dead export: true
- Smells: `dead-export`, `feature-envy`

## Signature
```text
func TestReadVaultSnapshotRoundTripsWriterOutput(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testwritevaultinputs--internal-vault-writer-test-go-l260]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/reader_integration_test.go|internal/vault/reader_integration_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/reader_integration_test.go|internal/vault/reader_integration_test.go]] via `exports` (syntactic)
