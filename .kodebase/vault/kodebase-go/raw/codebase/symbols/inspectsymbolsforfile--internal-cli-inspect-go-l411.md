---
blast_radius: 3
centrality: 0.061
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 452
exported: false
external_reference_count: 1
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 42
outgoing_relation_count: 2
smells:
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect.go"
stage: "raw"
start_line: 411
symbol_kind: "function"
symbol_name: "inspectSymbolsForFile"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: inspectSymbolsForFile"
type: "source"
---

# Codebase Symbol: inspectSymbolsForFile

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 3
- External references: 1
- Centrality: 0.061
- LOC: 42
- Dead export: false
- Smells: None

## Signature
```text
func inspectSymbolsForFile(snapshot vault.VaultSnapshot, sourcePath string) []string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterint--internal-cli-inspect-go-l213]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstring--internal-cli-inspect-go-l156]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/tofilelookupoutput--internal-cli-inspect-file-go-l31]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]] via `contains` (syntactic)
