---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 50
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 41
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/adapter/treesitter_test.go"
stage: "raw"
start_line: 10
symbol_kind: "function"
symbol_name: "TestLanguagesInitialize"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestLanguagesInitialize"
type: "source"
---

# Codebase Symbol: TestLanguagesInitialize

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/treesitter_test.go|internal/adapter/treesitter_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 41
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func TestLanguagesInitialize(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/adapter/treesitter_test.go|internal/adapter/treesitter_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/treesitter_test.go|internal/adapter/treesitter_test.go]] via `exports` (syntactic)
