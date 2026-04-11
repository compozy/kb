---
blast_radius: 5
centrality: 0.0982
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 72
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 17
outgoing_relation_count: 2
smells:
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_circulardeps.go"
stage: "raw"
start_line: 56
symbol_kind: "function"
symbol_name: "circularDependencyRowsFromFlags"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: circularDependencyRowsFromFlags"
type: "source"
---

# Codebase Symbol: circularDependencyRowsFromFlags

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_circulardeps.go|internal/cli/inspect_circulardeps.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 5
- External references: 0
- Centrality: 0.0982
- LOC: 17
- Dead export: false
- Smells: None

## Signature
```text
func circularDependencyRowsFromFlags(snapshot vault.VaultSnapshot) ([]map[string]any, bool) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterbool--internal-cli-inspect-go-l193]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tocirculardependencyrow--internal-cli-inspect-circulardeps-go-l145]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/circulardependencyrows--internal-cli-inspect-circulardeps-go-l44]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_circulardeps.go|internal/cli/inspect_circulardeps.go]] via `contains` (syntactic)
