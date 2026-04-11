---
blast_radius: 1
centrality: 0.0569
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 405
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 16
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/generate/generate.go"
stage: "raw"
start_line: 390
symbol_kind: "function"
symbol_name: "workspaceLanguages"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: workspaceLanguages"
type: "source"
---

# Codebase Symbol: workspaceLanguages

Source file: [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0569
- LOC: 16
- Dead export: false
- Smells: None

## Signature
```text
func workspaceLanguages(workspace *models.ScannedWorkspace) []models.SupportedLanguage {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/generatewithobserver--internal-generate-generate-go-l88]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]] via `contains` (syntactic)
