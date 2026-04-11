---
blast_radius: 5
centrality: 0.1489
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 395
exported: false
external_reference_count: 2
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 19
outgoing_relation_count: 2
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect.go"
stage: "raw"
start_line: 377
symbol_kind: "function"
symbol_name: "resolveInspectEntity"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: resolveInspectEntity"
type: "source"
---

# Codebase Symbol: resolveInspectEntity

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 5
- External references: 2
- Centrality: 0.1489
- LOC: 19
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func resolveInspectEntity(snapshot vault.VaultSnapshot, query string) (vault.VaultDocument, string, error) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/findinspectfilebysourcepath--internal-cli-inspect-go-l339]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/findsingleinspectsymbolmatch--internal-cli-inspect-go-l350]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/tobacklinksoutput--internal-cli-inspect-backlinks-go-l31]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/todependencyoutput--internal-cli-inspect-deps-go-l31]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]] via `contains` (syntactic)
