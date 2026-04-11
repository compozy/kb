---
blast_radius: 1
centrality: 0.0939
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 87
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
source_path: "internal/cli/generate_output.go"
stage: "raw"
start_line: 80
symbol_kind: "function"
symbol_name: "isTerminalWriter"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: isTerminalWriter"
type: "source"
---

# Codebase Symbol: isTerminalWriter

Source file: [[kodebase-go/raw/codebase/files/internal/cli/generate_output.go|internal/cli/generate_output.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0939
- LOC: 8
- Dead export: false
- Smells: None

## Signature
```text
func isTerminalWriter(writer io.Writer) bool {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/newgenerateobserver--internal-cli-generate-output-go-l57]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/generate_output.go|internal/cli/generate_output.go]] via `contains` (syntactic)
