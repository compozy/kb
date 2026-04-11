---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 11
domain: "kodebase-go"
end_line: 88
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 75
outgoing_relation_count: 3
smells:
  - "dead-export"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/graph/normalize_integration_test.go"
stage: "raw"
start_line: 14
symbol_kind: "function"
symbol_name: "TestNormalizeGraphMergesOverlappingImportsAcrossParsedFiles"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestNormalizeGraphMergesOverlappingImportsAcrossParsedFiles"
type: "source"
---

# Codebase Symbol: TestNormalizeGraphMergesOverlappingImportsAcrossParsedFiles

Source file: [[kodebase-go/raw/codebase/files/internal/graph/normalize_integration_test.go|internal/graph/normalize_integration_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 11
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 75
- Dead export: true
- Smells: `dead-export`, `long-function`

## Signature
```text
func TestNormalizeGraphMergesOverlappingImportsAcrossParsedFiles(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizegraph--internal-graph-normalize-go-l17]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/hasrelation--internal-graph-normalize-integration-test-go-l99]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writefixturefile--internal-graph-normalize-integration-test-go-l90]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/graph/normalize_integration_test.go|internal/graph/normalize_integration_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/graph/normalize_integration_test.go|internal/graph/normalize_integration_test.go]] via `exports` (syntactic)
