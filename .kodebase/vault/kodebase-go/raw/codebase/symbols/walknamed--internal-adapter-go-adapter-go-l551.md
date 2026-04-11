---
blast_radius: 20
centrality: 0.3565
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 564
exported: false
external_reference_count: 2
has_smells: true
incoming_relation_count: 7
is_dead_export: false
is_long_function: false
language: "go"
loc: 14
outgoing_relation_count: 2
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter.go"
stage: "raw"
start_line: 551
symbol_kind: "function"
symbol_name: "walkNamed"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: walkNamed"
type: "source"
---

# Codebase Symbol: walkNamed

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 20
- External references: 2
- Centrality: 0.3565
- LOC: 14
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func walkNamed(node *tree_sitter.Node, visit func(*tree_sitter.Node) bool) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/namedchildren--internal-adapter-go-adapter-go-l566]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/walknamed--internal-adapter-go-adapter-go-l551]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/collectnodesbykind--internal-adapter-go-adapter-go-l577]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/computecyclomaticcomplexity--internal-adapter-go-adapter-go-l465]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/extractcalltargetnames--internal-adapter-go-adapter-go-l429]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/walknamed--internal-adapter-go-adapter-go-l551]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/collecttscalltargets--internal-adapter-ts-adapter-go-l1018]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/computetscyclomaticcomplexity--internal-adapter-ts-adapter-go-l1069]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] via `contains` (syntactic)
