---
blast_radius: 12
centrality: 0.2082
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 574
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/qmd/client_test.go"
stage: "raw"
start_line: 572
symbol_kind: "function"
symbol_name: "shellQuote"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: shellQuote"
type: "source"
---

# Codebase Symbol: shellQuote

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 12
- External references: 0
- Centrality: 0.2082
- LOC: 3
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func shellQuote(value string) string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/writefakeqmd--internal-qmd-client-test-go-l490]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]] via `contains` (syntactic)
