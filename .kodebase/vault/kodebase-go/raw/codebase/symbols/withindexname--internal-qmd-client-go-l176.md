---
blast_radius: 2
centrality: 0.068
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 180
exported: true
external_reference_count: 2
has_smells: false
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 5
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/qmd/client.go"
stage: "raw"
start_line: 176
symbol_kind: "function"
symbol_name: "WithIndexName"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: WithIndexName"
type: "source"
---

# Codebase Symbol: WithIndexName

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 2
- External references: 2
- Centrality: 0.068
- LOC: 5
- Dead export: false
- Smells: None

## Signature
```text
func WithIndexName(name string) ClientOption {
```

## Documentation
WithIndexName routes QMD commands to a named SQLite index.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testqmdclientindexesandsearchestempvault--internal-qmd-client-integration-test-go-l13]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testindexaddconstructsexpectedarguments--internal-qmd-client-test-go-l27]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `exports` (syntactic)
