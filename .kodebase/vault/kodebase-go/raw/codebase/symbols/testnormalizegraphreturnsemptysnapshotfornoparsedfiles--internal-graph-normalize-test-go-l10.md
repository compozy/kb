---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 20
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 11
outgoing_relation_count: 2
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/graph/normalize_test.go"
stage: "raw"
start_line: 10
symbol_kind: "function"
symbol_name: "TestNormalizeGraphReturnsEmptySnapshotForNoParsedFiles"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestNormalizeGraphReturnsEmptySnapshotForNoParsedFiles"
type: "source"
---

# Codebase Symbol: TestNormalizeGraphReturnsEmptySnapshotForNoParsedFiles

Source file: [[kodebase-go/raw/codebase/files/internal/graph/normalize_test.go|internal/graph/normalize_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 11
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func TestNormalizeGraphReturnsEmptySnapshotForNoParsedFiles(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizegraph--internal-graph-normalize-go-l17]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/assertemptygraphsnapshot--internal-graph-normalize-test-go-l232]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/graph/normalize_test.go|internal/graph/normalize_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/graph/normalize_test.go|internal/graph/normalize_test.go]] via `exports` (syntactic)
