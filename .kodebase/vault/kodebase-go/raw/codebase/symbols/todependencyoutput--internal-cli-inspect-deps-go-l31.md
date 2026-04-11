---
blast_radius: 2
centrality: 0.137
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 41
exported: false
external_reference_count: 2
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 11
outgoing_relation_count: 2
smells:
  - "bottleneck"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_deps.go"
stage: "raw"
start_line: 31
symbol_kind: "function"
symbol_name: "toDependencyOutput"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: toDependencyOutput"
type: "source"
---

# Codebase Symbol: toDependencyOutput

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_deps.go|internal/cli/inspect_deps.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 2
- External references: 2
- Centrality: 0.137
- LOC: 11
- Dead export: false
- Smells: `bottleneck`, `feature-envy`

## Signature
```text
func toDependencyOutput(snapshot vault.VaultSnapshot, query string) (inspectOutput, error) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createinspectrelationrows--internal-cli-inspect-go-l315]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolveinspectentity--internal-cli-inspect-go-l377]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testtodependencyoutputlistsoutgoingdependenciesforsymbol--internal-cli-inspect-test-go-l414]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtodependencyoutputsupportsexactfilepathlookup--internal-cli-inspect-test-go-l442]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_deps.go|internal/cli/inspect_deps.go]] via `contains` (syntactic)
