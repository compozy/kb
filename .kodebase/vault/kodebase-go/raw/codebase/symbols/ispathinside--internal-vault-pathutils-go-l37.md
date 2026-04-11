---
blast_radius: 2
centrality: 0.0965
cyclomatic_complexity: 10
domain: "kodebase-go"
end_line: 74
exported: true
external_reference_count: 1
has_smells: false
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 38
outgoing_relation_count: 2
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/pathutils.go"
stage: "raw"
start_line: 37
symbol_kind: "function"
symbol_name: "IsPathInside"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: IsPathInside"
type: "source"
---

# Codebase Symbol: IsPathInside

Source file: [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 10
- Long function: false
- Blast radius: 2
- External references: 1
- Centrality: 0.0965
- LOC: 38
- Dead export: false
- Smells: None

## Signature
```text
func IsPathInside(parentPath, targetPath string) bool {
```

## Documentation
IsPathInside reports whether targetPath is the same as or nested under parentPath.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cleancomparablepath--internal-vault-pathutils-go-l237]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/splitcomparablepath--internal-vault-pathutils-go-l254]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/validatetopic--internal-vault-writer-go-l113]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `exports` (syntactic)
