---
blast_radius: 1
centrality: 0.0594
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 72
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 8
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/qmd/client_integration_test.go"
stage: "raw"
start_line: 65
symbol_kind: "function"
symbol_name: "writeMarkdownFile"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: writeMarkdownFile"
type: "source"
---

# Codebase Symbol: writeMarkdownFile

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client_integration_test.go|internal/qmd/client_integration_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0594
- LOC: 8
- Dead export: false
- Smells: None

## Signature
```text
func writeMarkdownFile(t *testing.T, rootPath, name, body string) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testqmdclientindexesandsearchestempvault--internal-qmd-client-integration-test-go-l13]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client_integration_test.go|internal/qmd/client_integration_test.go]] via `contains` (syntactic)
