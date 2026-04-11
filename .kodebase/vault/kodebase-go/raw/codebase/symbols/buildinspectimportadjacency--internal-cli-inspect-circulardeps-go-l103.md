---
blast_radius: 6
centrality: 0.0786
cyclomatic_complexity: 10
domain: "kodebase-go"
end_line: 143
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 41
outgoing_relation_count: 1
smells:
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_circulardeps.go"
stage: "raw"
start_line: 103
symbol_kind: "function"
symbol_name: "buildInspectImportAdjacency"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: buildInspectImportAdjacency"
type: "source"
---

# Codebase Symbol: buildInspectImportAdjacency

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_circulardeps.go|internal/cli/inspect_circulardeps.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 10
- Long function: false
- Blast radius: 6
- External references: 0
- Centrality: 0.0786
- LOC: 41
- Dead export: false
- Smells: `feature-envy`

## Signature
```text
func buildInspectImportAdjacency(snapshot vault.VaultSnapshot) map[string][]string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstring--internal-cli-inspect-go-l156]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/circulardependencyrowsfromfallback--internal-cli-inspect-circulardeps-go-l74]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_circulardeps.go|internal/cli/inspect_circulardeps.go]] via `contains` (syntactic)
