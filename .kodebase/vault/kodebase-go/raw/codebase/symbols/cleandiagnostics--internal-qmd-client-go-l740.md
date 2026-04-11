---
blast_radius: 11
centrality: 0.3196
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 744
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 5
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/qmd/client.go"
stage: "raw"
start_line: 740
symbol_kind: "function"
symbol_name: "cleanDiagnostics"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: cleanDiagnostics"
type: "source"
---

# Codebase Symbol: cleanDiagnostics

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 11
- External references: 0
- Centrality: 0.3196
- LOC: 5
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func cleanDiagnostics(output string) string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/cleanoutput--internal-qmd-client-go-l736]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/run--internal-qmd-client-go-l408]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `contains` (syntactic)
