---
blast_radius: 10
centrality: 0.159
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 375
exported: false
external_reference_count: 1
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 26
outgoing_relation_count: 1
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect.go"
stage: "raw"
start_line: 350
symbol_kind: "function"
symbol_name: "findSingleInspectSymbolMatch"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: findSingleInspectSymbolMatch"
type: "source"
---

# Codebase Symbol: findSingleInspectSymbolMatch

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 10
- External references: 1
- Centrality: 0.159
- LOC: 26
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func findSingleInspectSymbolMatch(snapshot vault.VaultSnapshot, query string) (vault.VaultDocument, error) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstring--internal-cli-inspect-go-l156]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/resolveinspectentity--internal-cli-inspect-go-l377]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/tosymbollookupoutput--internal-cli-inspect-symbol-go-l41]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]] via `contains` (syntactic)
