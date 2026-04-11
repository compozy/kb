---
blast_radius: 4
centrality: 0.1288
cyclomatic_complexity: 6
domain: "kodebase-go"
end_line: 248
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/vault/pathutils.go"
stage: "raw"
start_line: 246
symbol_kind: "function"
symbol_name: "hasWindowsDrivePrefix"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: hasWindowsDrivePrefix"
type: "source"
---

# Codebase Symbol: hasWindowsDrivePrefix

Source file: [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 6
- Long function: false
- Blast radius: 4
- External references: 0
- Centrality: 0.1288
- LOC: 3
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func hasWindowsDrivePrefix(value string) bool {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/splitcomparablepath--internal-vault-pathutils-go-l254]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `contains` (syntactic)
