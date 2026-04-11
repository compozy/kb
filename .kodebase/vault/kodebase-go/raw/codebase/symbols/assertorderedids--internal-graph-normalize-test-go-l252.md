---
blast_radius: 1
centrality: 0.0594
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 264
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 13
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/graph/normalize_test.go"
stage: "raw"
start_line: 252
symbol_kind: "function"
symbol_name: "assertOrderedIDs"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: assertOrderedIDs"
type: "source"
---

# Codebase Symbol: assertOrderedIDs

Source file: [[kodebase-go/raw/codebase/files/internal/graph/normalize_test.go|internal/graph/normalize_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0594
- LOC: 13
- Dead export: false
- Smells: None

## Signature
```text
func assertOrderedIDs(t *testing.T, files []models.GraphFile, expected []string) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testnormalizegraphsortscollectionsdeterministically--internal-graph-normalize-test-go-l132]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/graph/normalize_test.go|internal/graph/normalize_test.go]] via `contains` (syntactic)
