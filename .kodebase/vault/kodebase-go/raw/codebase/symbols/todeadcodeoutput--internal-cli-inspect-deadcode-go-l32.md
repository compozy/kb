---
blast_radius: 1
centrality: 0.0723
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 79
exported: false
external_reference_count: 1
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 48
outgoing_relation_count: 4
smells:
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_deadcode.go"
stage: "raw"
start_line: 32
symbol_kind: "function"
symbol_name: "toDeadCodeOutput"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: toDeadCodeOutput"
type: "source"
---

# Codebase Symbol: toDeadCodeOutput

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_deadcode.go|internal/cli/inspect_deadcode.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 1
- External references: 1
- Centrality: 0.0723
- LOC: 48
- Dead export: false
- Smells: `feature-envy`

## Signature
```text
func toDeadCodeOutput(snapshot vault.VaultSnapshot) inspectOutput {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterbool--internal-cli-inspect-go-l193]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstring--internal-cli-inspect-go-l156]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstringarray--internal-cli-inspect-go-l170]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/deadcoderowstomaps--internal-cli-inspect-deadcode-go-l81]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testtodeadcodeoutputlistsdeadexportsandorphanfiles--internal-cli-inspect-test-go-l55]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_deadcode.go|internal/cli/inspect_deadcode.go]] via `contains` (syntactic)
