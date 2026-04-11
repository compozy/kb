---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 24
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 20
outgoing_relation_count: 1
smells:
  - "dead-export"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/models/models_test.go"
stage: "raw"
start_line: 5
symbol_kind: "function"
symbol_name: "TestSupportedLanguages"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestSupportedLanguages"
type: "source"
---

# Codebase Symbol: TestSupportedLanguages

Source file: [[kodebase-go/raw/codebase/files/internal/models/models_test.go|internal/models/models_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 20
- Dead export: true
- Smells: `dead-export`, `feature-envy`

## Signature
```text
func TestSupportedLanguages(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/supportedlanguages--internal-models-models-go-l20]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/models/models_test.go|internal/models/models_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/models/models_test.go|internal/models/models_test.go]] via `exports` (syntactic)
