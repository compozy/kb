---
blast_radius: 1
centrality: 0.0569
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 452
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 8
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/generate/generate.go"
stage: "raw"
start_line: 445
symbol_kind: "function"
symbol_name: "adapterNames"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: adapterNames"
type: "source"
---

# Codebase Symbol: adapterNames

Source file: [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0569
- LOC: 8
- Dead export: false
- Smells: None

## Signature
```text
func adapterNames(adapters []models.LanguageAdapter) []string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/generatewithobserver--internal-generate-generate-go-l88]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]] via `contains` (syntactic)
