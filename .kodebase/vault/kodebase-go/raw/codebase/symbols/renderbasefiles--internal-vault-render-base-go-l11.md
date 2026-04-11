---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 195
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 185
outgoing_relation_count: 4
smells:
  - "dead-export"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/vault/render_base.go"
stage: "raw"
start_line: 11
symbol_kind: "function"
symbol_name: "RenderBaseFiles"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: RenderBaseFiles"
type: "source"
---

# Codebase Symbol: RenderBaseFiles

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_base.go|internal/vault/render_base.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 185
- Dead export: true
- Smells: `dead-export`, `long-function`

## Signature
```text
func RenderBaseFiles(metrics models.MetricsResult) []models.BaseFile {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/getbasefilepath--internal-vault-pathutils-go-l218]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/andfilter--internal-vault-render-base-go-l213]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createbasefile--internal-vault-render-base-go-l202]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/exprfilter--internal-vault-render-base-go-l209]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/render_base.go|internal/vault/render_base.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_base.go|internal/vault/render_base.go]] via `exports` (syntactic)
