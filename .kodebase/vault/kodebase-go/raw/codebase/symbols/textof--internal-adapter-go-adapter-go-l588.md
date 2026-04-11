---
blast_radius: 28
centrality: 0.6087
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 593
exported: false
external_reference_count: 11
has_smells: true
incoming_relation_count: 17
is_dead_export: false
is_long_function: false
language: "go"
loc: 6
outgoing_relation_count: 0
smells:
  - "bottleneck"
  - "high-blast-radius"
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter.go"
stage: "raw"
start_line: 588
symbol_kind: "function"
symbol_name: "textOf"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: textOf"
type: "source"
---

# Codebase Symbol: textOf

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 28
- External references: 11
- Centrality: 0.6087
- LOC: 6
- Dead export: false
- Smells: `bottleneck`, `high-blast-radius`

## Signature
```text
func textOf(node *tree_sitter.Node, source []byte) string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/extractattachedcomment--internal-adapter-go-adapter-go-l395]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/extractimports--internal-adapter-go-adapter-go-l264]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/formatgosignature--internal-adapter-go-adapter-go-l372]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/operatortext--internal-adapter-go-adapter-go-l503]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/resolvegosymbolname--internal-adapter-go-adapter-go-l354]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/extractcommonjsexports--internal-adapter-ts-adapter-go-l583]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/extractrequirebindings--internal-adapter-ts-adapter-go-l644]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/extracttsexports--internal-adapter-ts-adapter-go-l484]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/extracttsimports--internal-adapter-ts-adapter-go-l391]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/formattsreturntype--internal-adapter-ts-adapter-go-l968]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/formattssignature--internal-adapter-ts-adapter-go-l929]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/formattsvariabletypesuffix--internal-adapter-ts-adapter-go-l986]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/isdefaultexport--internal-adapter-ts-adapter-go-l1285]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/matchcommonjsexporttarget--internal-adapter-ts-adapter-go-l1290]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/resolvetssymbolname--internal-adapter-ts-adapter-go-l880]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/resolvetsvariablename--internal-adapter-ts-adapter-go-l916]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] via `contains` (syntactic)
