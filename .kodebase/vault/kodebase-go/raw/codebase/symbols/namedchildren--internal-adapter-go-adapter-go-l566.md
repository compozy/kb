---
blast_radius: 25
centrality: 0.511
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 575
exported: false
external_reference_count: 8
has_smells: true
incoming_relation_count: 11
is_dead_export: false
is_long_function: false
language: "go"
loc: 10
outgoing_relation_count: 0
smells:
  - "bottleneck"
  - "high-blast-radius"
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter.go"
stage: "raw"
start_line: 566
symbol_kind: "function"
symbol_name: "namedChildren"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: namedChildren"
type: "source"
---

# Codebase Symbol: namedChildren

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 25
- External references: 8
- Centrality: 0.511
- LOC: 10
- Dead export: false
- Smells: `bottleneck`, `high-blast-radius`

## Signature
```text
func namedChildren(node *tree_sitter.Node) []tree_sitter.Node {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parsegofile--internal-adapter-go-adapter-go-l161]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/walknamed--internal-adapter-go-adapter-go-l551]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/extractclassmethodsymbols--internal-adapter-ts-adapter-go-l782]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/extractrequirebindings--internal-adapter-ts-adapter-go-l644]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/extracttsexports--internal-adapter-ts-adapter-go-l484]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/extracttsimports--internal-adapter-ts-adapter-go-l391]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/findchildbykind--internal-adapter-ts-adapter-go-l1333]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/formattsreturntype--internal-adapter-ts-adapter-go-l968]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/formattsvariabletypesuffix--internal-adapter-ts-adapter-go-l986]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/parsetsfile--internal-adapter-ts-adapter-go-l299]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] via `contains` (syntactic)
