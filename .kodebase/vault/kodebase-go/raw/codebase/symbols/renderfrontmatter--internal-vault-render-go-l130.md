---
blast_radius: 2
centrality: 0.0953
cyclomatic_complexity: 11
domain: "kodebase-go"
end_line: 162
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: true
language: "go"
loc: 33
outgoing_relation_count: 1
smells:
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/vault/render.go"
stage: "raw"
start_line: 130
symbol_kind: "function"
symbol_name: "renderFrontmatter"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: renderFrontmatter"
type: "source"
---

# Codebase Symbol: renderFrontmatter

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 11
- Long function: true
- Blast radius: 2
- External references: 0
- Centrality: 0.0953
- LOC: 33
- Dead export: false
- Smells: `long-function`

## Signature
```text
func renderFrontmatter(frontmatter map[string]interface{}) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sortedmapkeys--internal-vault-render-go-l798]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/rendermarkdowndocument--internal-vault-render-go-l164]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]] via `contains` (syntactic)
