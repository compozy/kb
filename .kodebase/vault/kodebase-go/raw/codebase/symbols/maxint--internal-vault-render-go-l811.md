---
blast_radius: 6
centrality: 0.1298
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 816
exported: false
external_reference_count: 1
has_smells: true
incoming_relation_count: 5
is_dead_export: false
is_long_function: false
language: "go"
loc: 6
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/vault/render.go"
stage: "raw"
start_line: 811
symbol_kind: "function"
symbol_name: "maxInt"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: maxInt"
type: "source"
---

# Codebase Symbol: maxInt

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 6
- External references: 1
- Centrality: 0.1298
- LOC: 6
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func maxInt(left, right int) int {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/createsymbolfrontmatter--internal-vault-render-go-l402]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/defaultsymbolmetrics--internal-vault-render-go-l671]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/renderrawsymboldocument--internal-vault-render-go-l449]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createcomplexityhotspotsarticle--internal-vault-render-wiki-go-l547]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]] via `contains` (syntactic)
