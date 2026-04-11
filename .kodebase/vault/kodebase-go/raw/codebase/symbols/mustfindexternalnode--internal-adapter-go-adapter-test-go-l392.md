---
blast_radius: 1
centrality: 0.0651
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 403
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 12
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter_test.go"
stage: "raw"
start_line: 392
symbol_kind: "function"
symbol_name: "mustFindExternalNode"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: mustFindExternalNode"
type: "source"
---

# Codebase Symbol: mustFindExternalNode

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter_test.go|internal/adapter/go_adapter_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0651
- LOC: 12
- Dead export: false
- Smells: None

## Signature
```text
func mustFindExternalNode(t *testing.T, nodes []models.ExternalNode, source string) models.ExternalNode {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testgoadapterextractsimportrelations--internal-adapter-go-adapter-test-go-l157]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter_test.go|internal/adapter/go_adapter_test.go]] via `contains` (syntactic)
