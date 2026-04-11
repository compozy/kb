---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 305
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/cli/search_test.go"
stage: "raw"
start_line: 303
symbol_kind: "method"
symbol_name: "Search"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "method"
title: "Codebase Symbol: Search"
type: "source"
---

# Codebase Symbol: Search

Source file: [[kodebase-go/raw/codebase/files/internal/cli/search_test.go|internal/cli/search_test.go]]

## Kind
`method`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 3
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func (client fakeSearchClient) Search(ctx context.Context, options qmd.SearchOptions) ([]qmd.SearchResult, error) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/cli/search_test.go|internal/cli/search_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/search_test.go|internal/cli/search_test.go]] via `exports` (syntactic)
