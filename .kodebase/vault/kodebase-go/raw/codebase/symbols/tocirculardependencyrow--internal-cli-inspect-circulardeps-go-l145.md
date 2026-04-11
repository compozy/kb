---
blast_radius: 7
centrality: 0.1203
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 153
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 9
outgoing_relation_count: 4
smells:
  - "bottleneck"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_circulardeps.go"
stage: "raw"
start_line: 145
symbol_kind: "function"
symbol_name: "toCircularDependencyRow"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: toCircularDependencyRow"
type: "source"
---

# Codebase Symbol: toCircularDependencyRow

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_circulardeps.go|internal/cli/inspect_circulardeps.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 7
- External references: 0
- Centrality: 0.1203
- LOC: 9
- Dead export: false
- Smells: `bottleneck`, `feature-envy`

## Signature
```text
func toCircularDependencyRow(document vault.VaultDocument) map[string]any {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterfloat--internal-cli-inspect-go-l254]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterint--internal-cli-inspect-go-l213]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstring--internal-cli-inspect-go-l156]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstringarray--internal-cli-inspect-go-l170]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/circulardependencyrowsfromfallback--internal-cli-inspect-circulardeps-go-l74]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/circulardependencyrowsfromflags--internal-cli-inspect-circulardeps-go-l56]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_circulardeps.go|internal/cli/inspect_circulardeps.go]] via `contains` (syntactic)
