---
blast_radius: 2
centrality: 0.0939
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 90
exported: false
external_reference_count: 2
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 44
outgoing_relation_count: 5
smells:
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_blastradius.go"
stage: "raw"
start_line: 47
symbol_kind: "function"
symbol_name: "toBlastRadiusOutput"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: toBlastRadiusOutput"
type: "source"
---

# Codebase Symbol: toBlastRadiusOutput

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_blastradius.go|internal/cli/inspect_blastradius.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 2
- External references: 2
- Centrality: 0.0939
- LOC: 44
- Dead export: false
- Smells: `feature-envy`

## Signature
```text
func toBlastRadiusOutput(snapshot vault.VaultSnapshot, minimum, top int) inspectOutput {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterfloat--internal-cli-inspect-go-l254]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterint--internal-cli-inspect-go-l213]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstring--internal-cli-inspect-go-l156]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstringarray--internal-cli-inspect-go-l170]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/blastradiusrowstomaps--internal-cli-inspect-blastradius-go-l92]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testtoblastradiusoutputrespectsminimumandtop--internal-cli-inspect-helpers-test-go-l178]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtoblastradiusoutputsortsbydescendingblastradius--internal-cli-inspect-test-go-l132]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_blastradius.go|internal/cli/inspect_blastradius.go]] via `contains` (syntactic)
