---
blast_radius: 4
centrality: 0.1673
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 54
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 11
outgoing_relation_count: 3
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_circulardeps.go"
stage: "raw"
start_line: 44
symbol_kind: "function"
symbol_name: "circularDependencyRows"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: circularDependencyRows"
type: "source"
---

# Codebase Symbol: circularDependencyRows

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_circulardeps.go|internal/cli/inspect_circulardeps.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 4
- External references: 0
- Centrality: 0.1673
- LOC: 11
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func circularDependencyRows(snapshot vault.VaultSnapshot) []map[string]any {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/circulardependencyrowsfromfallback--internal-cli-inspect-circulardeps-go-l74]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/circulardependencyrowsfromflags--internal-cli-inspect-circulardeps-go-l56]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sortcirculardependencyrows--internal-cli-inspect-circulardeps-go-l155]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/tocirculardepsoutput--internal-cli-inspect-circulardeps-go-l27]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_circulardeps.go|internal/cli/inspect_circulardeps.go]] via `contains` (syntactic)
