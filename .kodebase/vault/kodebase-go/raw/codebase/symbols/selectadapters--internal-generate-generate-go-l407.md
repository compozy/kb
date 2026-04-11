---
blast_radius: 3
centrality: 0.1432
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 421
exported: false
external_reference_count: 2
has_smells: true
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 15
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/generate/generate.go"
stage: "raw"
start_line: 407
symbol_kind: "function"
symbol_name: "selectAdapters"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: selectAdapters"
type: "source"
---

# Codebase Symbol: selectAdapters

Source file: [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 3
- External references: 2
- Centrality: 0.1432
- LOC: 15
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func selectAdapters(languages []models.SupportedLanguage, adapters []models.LanguageAdapter) []models.LanguageAdapter {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/generatewithobserver--internal-generate-generate-go-l88]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testselectadaptersforgoonlyworkspace--internal-generate-generate-test-go-l172]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testselectadaptersformixedworkspace--internal-generate-generate-test-go-l191]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]] via `contains` (syntactic)
