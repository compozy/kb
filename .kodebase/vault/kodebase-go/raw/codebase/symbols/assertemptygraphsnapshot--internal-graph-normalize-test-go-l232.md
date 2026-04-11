---
blast_radius: 1
centrality: 0.0723
cyclomatic_complexity: 6
domain: "kodebase-go"
end_line: 250
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 19
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/graph/normalize_test.go"
stage: "raw"
start_line: 232
symbol_kind: "function"
symbol_name: "assertEmptyGraphSnapshot"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: assertEmptyGraphSnapshot"
type: "source"
---

# Codebase Symbol: assertEmptyGraphSnapshot

Source file: [[kodebase-go/raw/codebase/files/internal/graph/normalize_test.go|internal/graph/normalize_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 6
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0723
- LOC: 19
- Dead export: false
- Smells: None

## Signature
```text
func assertEmptyGraphSnapshot(t *testing.T, snapshot models.GraphSnapshot) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testnormalizegraphreturnsemptysnapshotfornoparsedfiles--internal-graph-normalize-test-go-l10]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/graph/normalize_test.go|internal/graph/normalize_test.go]] via `contains` (syntactic)
