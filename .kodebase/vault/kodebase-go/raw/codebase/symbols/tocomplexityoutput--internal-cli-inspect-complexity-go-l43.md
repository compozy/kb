---
blast_radius: 1
centrality: 0.0723
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 87
exported: false
external_reference_count: 1
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 45
outgoing_relation_count: 5
smells:
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_complexity.go"
stage: "raw"
start_line: 43
symbol_kind: "function"
symbol_name: "toComplexityOutput"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: toComplexityOutput"
type: "source"
---

# Codebase Symbol: toComplexityOutput

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_complexity.go|internal/cli/inspect_complexity.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 1
- External references: 1
- Centrality: 0.0723
- LOC: 45
- Dead export: false
- Smells: `feature-envy`

## Signature
```text
func toComplexityOutput(snapshot vault.VaultSnapshot, top int) inspectOutput {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterint--internal-cli-inspect-go-l213]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstring--internal-cli-inspect-go-l156]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstringarray--internal-cli-inspect-go-l170]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isfunctionlikedocument--internal-cli-inspect-go-l147]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/complexityrowstomaps--internal-cli-inspect-complexity-go-l89]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testtocomplexityoutputsortsbydescendingcomplexity--internal-cli-inspect-test-go-l90]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_complexity.go|internal/cli/inspect_complexity.go]] via `contains` (syntactic)
