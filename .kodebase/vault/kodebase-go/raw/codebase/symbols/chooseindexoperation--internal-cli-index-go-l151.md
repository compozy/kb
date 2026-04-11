---
blast_radius: 1
centrality: 0.0651
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 159
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 9
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/cli/index.go"
stage: "raw"
start_line: 151
symbol_kind: "function"
symbol_name: "chooseIndexOperation"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: chooseIndexOperation"
type: "source"
---

# Codebase Symbol: chooseIndexOperation

Source file: [[kodebase-go/raw/codebase/files/internal/cli/index.go|internal/cli/index.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0651
- LOC: 9
- Dead export: false
- Smells: None

## Signature
```text
func chooseIndexOperation(status qmd.IndexStatus, collectionName string) qmd.IndexOperation {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/runindexcommand--internal-cli-index-go-l63]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/index.go|internal/cli/index.go]] via `contains` (syntactic)
