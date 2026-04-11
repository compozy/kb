---
blast_radius: 1
centrality: 0.0615
cyclomatic_complexity: 17
domain: "kodebase-go"
end_line: 262
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: true
language: "go"
loc: 102
outgoing_relation_count: 8
smells:
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter.go"
stage: "raw"
start_line: 161
symbol_kind: "function"
symbol_name: "parseGoFile"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: parseGoFile"
type: "source"
---

# Codebase Symbol: parseGoFile

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 17
- Long function: true
- Blast radius: 1
- External references: 0
- Centrality: 0.0615
- LOC: 102
- Dead export: false
- Smells: `long-function`

## Signature
```text
func parseGoFile(parser *tree_sitter.Parser, file models.ScannedSourceFile) (parsedGoFile, error) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createfileid--internal-adapter-go-adapter-go-l669]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/creategoparsediagnostic--internal-adapter-go-adapter-go-l520]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/creategosymbol--internal-adapter-go-adapter-go-l318]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractcalltargetnames--internal-adapter-go-adapter-go-l429]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractimports--internal-adapter-go-adapter-go-l264]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/extractleadingcomment--internal-adapter-go-adapter-go-l595]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getgosymbolkind--internal-adapter-go-adapter-go-l302]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/namedchildren--internal-adapter-go-adapter-go-l566]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parsefileswithprogress--internal-adapter-go-adapter-go-l62]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] via `contains` (syntactic)
