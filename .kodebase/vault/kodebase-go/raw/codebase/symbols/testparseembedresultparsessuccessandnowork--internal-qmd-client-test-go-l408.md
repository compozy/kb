---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 8
domain: "kodebase-go"
end_line: 426
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 19
outgoing_relation_count: 1
smells:
  - "dead-export"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/qmd/client_test.go"
stage: "raw"
start_line: 408
symbol_kind: "function"
symbol_name: "TestParseEmbedResultParsesSuccessAndNoWork"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestParseEmbedResultParsesSuccessAndNoWork"
type: "source"
---

# Codebase Symbol: TestParseEmbedResultParsesSuccessAndNoWork

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 8
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 19
- Dead export: true
- Smells: `dead-export`, `feature-envy`

## Signature
```text
func TestParseEmbedResultParsesSuccessAndNoWork(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parseembedresult--internal-qmd-client-go-l555]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]] via `exports` (syntactic)
