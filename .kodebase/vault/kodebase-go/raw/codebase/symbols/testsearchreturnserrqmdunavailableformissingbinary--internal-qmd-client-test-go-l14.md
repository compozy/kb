---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 25
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 12
outgoing_relation_count: 2
smells:
  - "dead-export"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/qmd/client_test.go"
stage: "raw"
start_line: 14
symbol_kind: "function"
symbol_name: "TestSearchReturnsErrQMDUnavailableForMissingBinary"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestSearchReturnsErrQMDUnavailableForMissingBinary"
type: "source"
---

# Codebase Symbol: TestSearchReturnsErrQMDUnavailableForMissingBinary

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 12
- Dead export: true
- Smells: `dead-export`, `feature-envy`

## Signature
```text
func TestSearchReturnsErrQMDUnavailableForMissingBinary(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newclient--internal-qmd-client-go-l154]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withbinarypath--internal-qmd-client-go-l169]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]] via `exports` (syntactic)
