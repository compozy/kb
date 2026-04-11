---
blast_radius: 1
centrality: 0.0524
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 598
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: true
language: "go"
loc: 75
outgoing_relation_count: 5
smells:
  - "feature-envy"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/vault/render.go"
stage: "raw"
start_line: 524
symbol_kind: "function"
symbol_name: "renderRawDirectoryIndex"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: renderRawDirectoryIndex"
type: "source"
---

# Codebase Symbol: renderRawDirectoryIndex

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: true
- Blast radius: 1
- External references: 0
- Centrality: 0.0524
- LOC: 75
- Dead export: false
- Smells: `feature-envy`, `long-function`

## Signature
```text
func renderRawDirectoryIndex(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getrawdirectoryindexpath--internal-vault-pathutils-go-l193]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getrawfiledocumentpath--internal-vault-pathutils-go-l181]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getrawsymboldocumentpath--internal-vault-pathutils-go-l186]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sortsymbolsbylocation--internal-vault-render-go-l783]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tosourcewikilink--internal-vault-render-go-l173]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/renderdocuments--internal-vault-render-go-l20]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]] via `contains` (syntactic)
