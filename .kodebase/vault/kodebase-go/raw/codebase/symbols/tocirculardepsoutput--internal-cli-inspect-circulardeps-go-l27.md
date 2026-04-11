---
blast_radius: 3
centrality: 0.137
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 42
exported: false
external_reference_count: 3
has_smells: true
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 16
outgoing_relation_count: 1
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_circulardeps.go"
stage: "raw"
start_line: 27
symbol_kind: "function"
symbol_name: "toCircularDepsOutput"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: toCircularDepsOutput"
type: "source"
---

# Codebase Symbol: toCircularDepsOutput

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_circulardeps.go|internal/cli/inspect_circulardeps.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 3
- External references: 3
- Centrality: 0.137
- LOC: 16
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func toCircularDepsOutput(snapshot vault.VaultSnapshot) inspectOutput {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/circulardependencyrows--internal-cli-inspect-circulardeps-go-l44]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testtocirculardepsoutputfallsbacktosccdetectionforlegacyvaults--internal-cli-inspect-test-go-l508]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtocirculardepsoutputlistsfileswithcirculardependencyflag--internal-cli-inspect-test-go-l467]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtocirculardepsoutputshowsmessagewhennocycles--internal-cli-inspect-test-go-l538]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_circulardeps.go|internal/cli/inspect_circulardeps.go]] via `contains` (syntactic)
