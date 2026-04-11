---
blast_radius: 1
centrality: 0.0723
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 213
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 7
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/cli/generate_output.go"
stage: "raw"
start_line: 207
symbol_kind: "function"
symbol_name: "completedSuffix"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: completedSuffix"
type: "source"
---

# Codebase Symbol: completedSuffix

Source file: [[kodebase-go/raw/codebase/files/internal/cli/generate_output.go|internal/cli/generate_output.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0723
- LOC: 7
- Dead export: false
- Smells: None

## Signature
```text
func completedSuffix(event kgenerate.Event) string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/handlestagecompleted--internal-cli-generate-output-go-l162]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/generate_output.go|internal/cli/generate_output.go]] via `contains` (syntactic)
