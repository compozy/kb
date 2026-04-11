---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 218
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
start_line: 213
symbol_kind: "method"
symbol_name: "Status"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "method"
title: "Codebase Symbol: Status"
type: "source"
---

# Codebase Symbol: Status

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
func (client fakeIndexClient) Status(ctx context.Context) (qmd.IndexStatus, error) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/cli/index_test.go|internal/cli/index_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/index_test.go|internal/cli/index_test.go]] via `exports` (syntactic)
