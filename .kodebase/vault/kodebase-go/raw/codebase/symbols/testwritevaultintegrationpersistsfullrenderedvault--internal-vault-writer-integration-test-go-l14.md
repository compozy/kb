---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 7
domain: "kodebase-go"
end_line: 50
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 37
outgoing_relation_count: 3
smells:
  - "dead-export"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/vault/writer_integration_test.go"
stage: "raw"
start_line: 14
symbol_kind: "function"
symbol_name: "TestWriteVaultIntegrationPersistsFullRenderedVault"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestWriteVaultIntegrationPersistsFullRenderedVault"
type: "source"
---

# Codebase Symbol: TestWriteVaultIntegrationPersistsFullRenderedVault

Source file: [[kodebase-go/raw/codebase/files/internal/vault/writer_integration_test.go|internal/vault/writer_integration_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 7
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 37
- Dead export: true
- Smells: `dead-export`, `feature-envy`

## Signature
```text
func TestWriteVaultIntegrationPersistsFullRenderedVault(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/assertfileexists--internal-vault-writer-test-go-l344]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/countkinds--internal-vault-writer-test-go-l293]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testwritevaultinputs--internal-vault-writer-test-go-l260]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/writer_integration_test.go|internal/vault/writer_integration_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/writer_integration_test.go|internal/vault/writer_integration_test.go]] via `exports` (syntactic)
