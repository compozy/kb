---
blast_radius: 1
centrality: 0.0939
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 22
exported: true
external_reference_count: 1
has_smells: false
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/models/models.go"
stage: "raw"
start_line: 20
symbol_kind: "function"
symbol_name: "SupportedLanguages"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: SupportedLanguages"
type: "source"
---

# Codebase Symbol: SupportedLanguages

Source file: [[kodebase-go/raw/codebase/files/internal/models/models.go|internal/models/models.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 1
- External references: 1
- Centrality: 0.0939
- LOC: 3
- Dead export: false
- Smells: None

## Signature
```text
func SupportedLanguages() []SupportedLanguage {
```

## Documentation
SupportedLanguages returns every supported language constant in stable order.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testsupportedlanguages--internal-models-models-test-go-l5]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/models/models.go|internal/models/models.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/models/models.go|internal/models/models.go]] via `exports` (syntactic)
