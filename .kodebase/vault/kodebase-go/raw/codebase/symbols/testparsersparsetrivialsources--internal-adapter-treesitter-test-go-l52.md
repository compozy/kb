---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 118
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 67
outgoing_relation_count: 0
smells:
  - "dead-export"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/adapter/treesitter_test.go"
stage: "raw"
start_line: 52
symbol_kind: "function"
symbol_name: "TestParsersParseTrivialSources"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestParsersParseTrivialSources"
type: "source"
---

# Codebase Symbol: TestParsersParseTrivialSources

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/treesitter_test.go|internal/adapter/treesitter_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 67
- Dead export: true
- Smells: `dead-export`, `long-function`

## Signature
```text
func TestParsersParseTrivialSources(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/adapter/treesitter_test.go|internal/adapter/treesitter_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/treesitter_test.go|internal/adapter/treesitter_test.go]] via `exports` (syntactic)
