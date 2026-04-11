---
blast_radius: 0
centrality: 0.0507
domain: "kodebase-go"
end_line: 43
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 1
is_dead_export: false
language: "go"
loc: 1
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/qmd/client.go"
stage: "raw"
start_line: 43
symbol_kind: "type"
symbol_name: "commandContextFunc"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "type"
title: "Codebase Symbol: commandContextFunc"
type: "source"
---

# Codebase Symbol: commandContextFunc

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]]

## Kind
`type`

## Static Analysis
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 1
- Dead export: false
- Smells: None

## Signature
```text
commandContextFunc func(context.Context, string, ...string) *exec.Cmd
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `contains` (syntactic)
