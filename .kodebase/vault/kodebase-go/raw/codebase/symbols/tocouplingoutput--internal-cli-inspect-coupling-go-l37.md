---
blast_radius: 2
centrality: 0.0939
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 73
exported: false
external_reference_count: 2
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 37
outgoing_relation_count: 6
smells:
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_coupling.go"
stage: "raw"
start_line: 37
symbol_kind: "function"
symbol_name: "toCouplingOutput"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: toCouplingOutput"
type: "source"
---

# Codebase Symbol: toCouplingOutput

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_coupling.go|internal/cli/inspect_coupling.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 2
- External references: 2
- Centrality: 0.0939
- LOC: 37
- Dead export: false
- Smells: `feature-envy`

## Signature
```text
func toCouplingOutput(snapshot vault.VaultSnapshot, unstableOnly bool) inspectOutput {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterbool--internal-cli-inspect-go-l193]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterfloat--internal-cli-inspect-go-l254]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterint--internal-cli-inspect-go-l213]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstring--internal-cli-inspect-go-l156]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstringarray--internal-cli-inspect-go-l170]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/couplingrowstomaps--internal-cli-inspect-coupling-go-l75]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testtocouplingoutputfiltersunstableonly--internal-cli-inspect-helpers-test-go-l198]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtocouplingoutputsortsbyinstability--internal-cli-inspect-test-go-l166]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_coupling.go|internal/cli/inspect_coupling.go]] via `contains` (syntactic)
