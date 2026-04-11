---
blast_radius: 1
centrality: 0.0939
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 41
exported: false
external_reference_count: 1
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 11
outgoing_relation_count: 2
smells:
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_backlinks.go"
stage: "raw"
start_line: 31
symbol_kind: "function"
symbol_name: "toBacklinksOutput"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: toBacklinksOutput"
type: "source"
---

# Codebase Symbol: toBacklinksOutput

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_backlinks.go|internal/cli/inspect_backlinks.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 1
- External references: 1
- Centrality: 0.0939
- LOC: 11
- Dead export: false
- Smells: `feature-envy`

## Signature
```text
func toBacklinksOutput(snapshot vault.VaultSnapshot, query string) (inspectOutput, error) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createinspectrelationrows--internal-cli-inspect-go-l315]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolveinspectentity--internal-cli-inspect-go-l377]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testtobacklinksoutputlistsincomingreferencesforsymbol--internal-cli-inspect-test-go-l385]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_backlinks.go|internal/cli/inspect_backlinks.go]] via `contains` (syntactic)
