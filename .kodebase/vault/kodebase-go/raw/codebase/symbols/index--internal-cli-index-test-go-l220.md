---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 225
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 6
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/cli/index_test.go"
stage: "raw"
start_line: 220
symbol_kind: "method"
symbol_name: "Index"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "method"
title: "Codebase Symbol: Index"
type: "source"
---

# Codebase Symbol: Index

Source file: [[kodebase-go/raw/codebase/files/internal/cli/index_test.go|internal/cli/index_test.go]]

## Kind
`method`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 6
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func (client fakeIndexClient) Index(ctx context.Context, options qmd.IndexOptions) (qmd.IndexResult, error) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/cli/index_test.go|internal/cli/index_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/index_test.go|internal/cli/index_test.go]] via `exports` (syntactic)
