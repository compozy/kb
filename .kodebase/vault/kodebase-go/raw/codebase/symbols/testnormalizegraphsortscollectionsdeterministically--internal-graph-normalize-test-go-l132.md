---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 171
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 40
outgoing_relation_count: 5
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/graph/normalize_test.go"
stage: "raw"
start_line: 132
symbol_kind: "function"
symbol_name: "TestNormalizeGraphSortsCollectionsDeterministically"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestNormalizeGraphSortsCollectionsDeterministically"
type: "source"
---

# Codebase Symbol: TestNormalizeGraphSortsCollectionsDeterministically

Source file: [[kodebase-go/raw/codebase/files/internal/graph/normalize_test.go|internal/graph/normalize_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 40
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func TestNormalizeGraphSortsCollectionsDeterministically(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizegraph--internal-graph-normalize-go-l17]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/assertorderedexternalids--internal-graph-normalize-test-go-l280]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/assertorderedids--internal-graph-normalize-test-go-l252]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/assertorderedsymbolids--internal-graph-normalize-test-go-l266]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsedfilefixture--internal-graph-normalize-test-go-l294]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/graph/normalize_test.go|internal/graph/normalize_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/graph/normalize_test.go|internal/graph/normalize_test.go]] via `exports` (syntactic)
