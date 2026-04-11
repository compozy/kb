---
blast_radius: 3
centrality: 0.0918
cyclomatic_complexity: 7
domain: "kodebase-go"
end_line: 279
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 26
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/pathutils.go"
stage: "raw"
start_line: 254
symbol_kind: "function"
symbol_name: "splitComparablePath"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: splitComparablePath"
type: "source"
---

# Codebase Symbol: splitComparablePath

Source file: [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 7
- Long function: false
- Blast radius: 3
- External references: 0
- Centrality: 0.0918
- LOC: 26
- Dead export: false
- Smells: None

## Signature
```text
func splitComparablePath(value string) (drive string, parts []string, absolute bool) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/haswindowsdriveprefix--internal-vault-pathutils-go-l246]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/ispathinside--internal-vault-pathutils-go-l37]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `contains` (syntactic)
