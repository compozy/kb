---
blast_radius: 2
centrality: 0.073
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 63
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 21
outgoing_relation_count: 1
smells:
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/vault/render_wiki.go"
stage: "raw"
start_line: 43
symbol_kind: "function"
symbol_name: "makeWikiFrontmatter"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: makeWikiFrontmatter"
type: "source"
---

# Codebase Symbol: makeWikiFrontmatter

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_wiki.go|internal/vault/render_wiki.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.073
- LOC: 21
- Dead export: false
- Smells: `feature-envy`

## Signature
```text
func makeWikiFrontmatter(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tosourcewikilink--internal-vault-render-go-l173]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/renderwikiarticle--internal-vault-render-wiki-go-l65]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_wiki.go|internal/vault/render_wiki.go]] via `contains` (syntactic)
