---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 201
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 29
outgoing_relation_count: 2
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/graph/normalize_test.go"
stage: "raw"
start_line: 173
symbol_kind: "function"
symbol_name: "TestNormalizeGraphOrdersDiagnosticsByStageFilePathAndMessage"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestNormalizeGraphOrdersDiagnosticsByStageFilePathAndMessage"
type: "source"
---

# Codebase Symbol: TestNormalizeGraphOrdersDiagnosticsByStageFilePathAndMessage

Source file: [[kodebase-go/raw/codebase/files/internal/graph/normalize_test.go|internal/graph/normalize_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 29
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func TestNormalizeGraphOrdersDiagnosticsByStageFilePathAndMessage(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizegraph--internal-graph-normalize-go-l17]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsedfilefixture--internal-graph-normalize-test-go-l294]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/graph/normalize_test.go|internal/graph/normalize_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/graph/normalize_test.go|internal/graph/normalize_test.go]] via `exports` (syntactic)
