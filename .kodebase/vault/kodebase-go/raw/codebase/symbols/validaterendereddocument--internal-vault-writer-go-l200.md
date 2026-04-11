---
blast_radius: 2
centrality: 0.066
cyclomatic_complexity: 8
domain: "kodebase-go"
end_line: 240
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 41
outgoing_relation_count: 2
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/writer.go"
stage: "raw"
start_line: 200
symbol_kind: "function"
symbol_name: "validateRenderedDocument"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: validateRenderedDocument"
type: "source"
---

# Codebase Symbol: validateRenderedDocument

Source file: [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 8
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.066
- LOC: 41
- Dead export: false
- Smells: None

## Signature
```text
func validateRenderedDocument(document models.RenderedDocument) (string, error) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cleantopicrelativepath--internal-vault-writer-go-l270]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/expecteddocumentplacement--internal-vault-writer-go-l257]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/buildwriterequests--internal-vault-writer-go-l166]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]] via `contains` (syntactic)
