---
blast_radius: 1
centrality: 0.0569
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 434
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
source_path: "internal/generate/generate.go"
stage: "raw"
start_line: 423
symbol_kind: "function"
symbol_name: "filesForAdapter"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: filesForAdapter"
type: "source"
---

# Codebase Symbol: filesForAdapter

Source file: [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0569
- LOC: 12
- Dead export: false
- Smells: None

## Signature
```text
func filesForAdapter(files []models.ScannedSourceFile, languageAdapter models.LanguageAdapter) []models.ScannedSourceFile {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/generatewithobserver--internal-generate-generate-go-l88]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]] via `contains` (syntactic)
