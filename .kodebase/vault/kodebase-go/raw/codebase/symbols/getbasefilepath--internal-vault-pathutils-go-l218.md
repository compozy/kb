---
blast_radius: 1
centrality: 0.0615
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 220
exported: true
external_reference_count: 1
has_smells: false
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/pathutils.go"
stage: "raw"
start_line: 218
symbol_kind: "function"
symbol_name: "GetBaseFilePath"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: GetBaseFilePath"
type: "source"
---

# Codebase Symbol: GetBaseFilePath

Source file: [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 1
- External references: 1
- Centrality: 0.0615
- LOC: 3
- Dead export: false
- Smells: None

## Signature
```text
func GetBaseFilePath(baseName string) string {
```

## Documentation
GetBaseFilePath derives the vault document path for an Obsidian Base definition.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/renderbasefiles--internal-vault-render-base-go-l11]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `exports` (syntactic)
