---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 7
domain: "kodebase-go"
end_line: 63
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 51
outgoing_relation_count: 5
smells:
  - "dead-export"
  - "feature-envy"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/qmd/client_integration_test.go"
stage: "raw"
start_line: 13
symbol_kind: "function"
symbol_name: "TestQMDClientIndexesAndSearchesTempVault"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestQMDClientIndexesAndSearchesTempVault"
type: "source"
---

# Codebase Symbol: TestQMDClientIndexesAndSearchesTempVault

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client_integration_test.go|internal/qmd/client_integration_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 7
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 51
- Dead export: true
- Smells: `dead-export`, `feature-envy`, `long-function`

## Signature
```text
func TestQMDClientIndexesAndSearchesTempVault(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newclient--internal-qmd-client-go-l154]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withbinarypath--internal-qmd-client-go-l169]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withindexname--internal-qmd-client-go-l176]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sanitizetestidentifier--internal-qmd-client-integration-test-go-l74]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writemarkdownfile--internal-qmd-client-integration-test-go-l65]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/qmd/client_integration_test.go|internal/qmd/client_integration_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client_integration_test.go|internal/qmd/client_integration_test.go]] via `exports` (syntactic)
