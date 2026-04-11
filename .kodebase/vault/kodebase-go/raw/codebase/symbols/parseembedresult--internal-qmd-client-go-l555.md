---
blast_radius: 2
centrality: 0.1047
cyclomatic_complexity: 10
domain: "kodebase-go"
end_line: 596
exported: false
external_reference_count: 1
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 42
outgoing_relation_count: 2
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/qmd/client.go"
stage: "raw"
start_line: 555
symbol_kind: "function"
symbol_name: "parseEmbedResult"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: parseEmbedResult"
type: "source"
---

# Codebase Symbol: parseEmbedResult

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 10
- Long function: false
- Blast radius: 2
- External references: 1
- Centrality: 0.1047
- LOC: 42
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func parseEmbedResult(output string) (EmbedResult, error) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cleanoutput--internal-qmd-client-go-l736]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsehumandurationmilliseconds--internal-qmd-client-go-l701]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/index--internal-qmd-client-go-l232]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testparseembedresultparsessuccessandnowork--internal-qmd-client-test-go-l408]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `contains` (syntactic)
