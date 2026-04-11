---
blast_radius: 1
centrality: 0.0524
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 171
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 8
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/render.go"
stage: "raw"
start_line: 164
symbol_kind: "function"
symbol_name: "renderMarkdownDocument"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: renderMarkdownDocument"
type: "source"
---

# Codebase Symbol: renderMarkdownDocument

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0524
- LOC: 8
- Dead export: false
- Smells: None

## Signature
```text
func renderMarkdownDocument(document models.RenderedDocument) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderfrontmatter--internal-vault-render-go-l130]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/renderdocuments--internal-vault-render-go-l20]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]] via `contains` (syntactic)
