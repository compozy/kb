---
blast_radius: 1
centrality: 0.0524
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 400
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: true
language: "go"
loc: 89
outgoing_relation_count: 6
smells:
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/vault/render.go"
stage: "raw"
start_line: 312
symbol_kind: "function"
symbol_name: "renderRawFileDocument"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: renderRawFileDocument"
type: "source"
---

# Codebase Symbol: renderRawFileDocument

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: true
- Blast radius: 1
- External references: 0
- Centrality: 0.0524
- LOC: 89
- Dead export: false
- Smells: `long-function`

## Signature
```text
func renderRawFileDocument(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getrawfiledocumentpath--internal-vault-pathutils-go-l181]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getrawsymboldocumentpath--internal-vault-pathutils-go-l186]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderbacklinklist--internal-vault-render-go-l262]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderrelationlist--internal-vault-render-go-l221]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/rendersmelllist--internal-vault-render-go-l300]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tosourcewikilink--internal-vault-render-go-l173]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/renderdocuments--internal-vault-render-go-l20]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]] via `contains` (syntactic)
