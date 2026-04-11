---
blast_radius: 5
centrality: 0.0982
cyclomatic_complexity: 7
domain: "kodebase-go"
end_line: 101
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 28
outgoing_relation_count: 3
smells:
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_circulardeps.go"
stage: "raw"
start_line: 74
symbol_kind: "function"
symbol_name: "circularDependencyRowsFromFallback"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: circularDependencyRowsFromFallback"
type: "source"
---

# Codebase Symbol: circularDependencyRowsFromFallback

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_circulardeps.go|internal/cli/inspect_circulardeps.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 7
- Long function: false
- Blast radius: 5
- External references: 0
- Centrality: 0.0982
- LOC: 28
- Dead export: false
- Smells: None

## Signature
```text
func circularDependencyRowsFromFallback(snapshot vault.VaultSnapshot) []map[string]any {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstring--internal-cli-inspect-go-l156]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/buildinspectimportadjacency--internal-cli-inspect-circulardeps-go-l103]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tocirculardependencyrow--internal-cli-inspect-circulardeps-go-l145]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/circulardependencyrows--internal-cli-inspect-circulardeps-go-l44]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_circulardeps.go|internal/cli/inspect_circulardeps.go]] via `contains` (syntactic)
