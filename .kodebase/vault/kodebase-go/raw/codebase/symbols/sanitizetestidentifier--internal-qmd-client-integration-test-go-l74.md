---
blast_radius: 1
centrality: 0.0594
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 77
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 4
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/qmd/client_integration_test.go"
stage: "raw"
start_line: 74
symbol_kind: "function"
symbol_name: "sanitizeTestIdentifier"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: sanitizeTestIdentifier"
type: "source"
---

# Codebase Symbol: sanitizeTestIdentifier

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client_integration_test.go|internal/qmd/client_integration_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0594
- LOC: 4
- Dead export: false
- Smells: None

## Signature
```text
func sanitizeTestIdentifier(value string) string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testqmdclientindexesandsearchestempvault--internal-qmd-client-integration-test-go-l13]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client_integration_test.go|internal/qmd/client_integration_test.go]] via `contains` (syntactic)
