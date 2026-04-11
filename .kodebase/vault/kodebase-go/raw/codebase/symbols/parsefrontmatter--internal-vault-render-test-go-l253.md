---
blast_radius: 3
centrality: 0.0903
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 273
exported: false
external_reference_count: 2
has_smells: false
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 21
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/render_test.go"
stage: "raw"
start_line: 253
symbol_kind: "function"
symbol_name: "parseFrontmatter"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: parseFrontmatter"
type: "source"
---

# Codebase Symbol: parseFrontmatter

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_test.go|internal/vault/render_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 3
- External references: 2
- Centrality: 0.0903
- LOC: 21
- Dead export: false
- Smells: None

## Signature
```text
func parseFrontmatter(t *testing.T, body string) (map[string]interface{}, string) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsintegrationproducesfulldocumentset--internal-vault-render-integration-test-go-l14]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testrenderdocumentsbodieshavevalidfrontmatterandkinds--internal-vault-render-test-go-l196]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testwritevaultcreatestopicskeletonandmanagedfiles--internal-vault-writer-test-go-l16]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_test.go|internal/vault/render_test.go]] via `contains` (syntactic)
