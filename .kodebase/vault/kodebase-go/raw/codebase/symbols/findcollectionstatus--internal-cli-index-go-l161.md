---
blast_radius: 1
centrality: 0.0651
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 172
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 12
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/cli/index.go"
stage: "raw"
start_line: 161
symbol_kind: "function"
symbol_name: "findCollectionStatus"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: findCollectionStatus"
type: "source"
---

# Codebase Symbol: findCollectionStatus

Source file: [[kodebase-go/raw/codebase/files/internal/cli/index.go|internal/cli/index.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0651
- LOC: 12
- Dead export: false
- Smells: None

## Signature
```text
func findCollectionStatus(collections []qmd.CollectionInfo, collectionName string) *qmd.CollectionInfo {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/runindexcommand--internal-cli-index-go-l63]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/index.go|internal/cli/index.go]] via `contains` (syntactic)
