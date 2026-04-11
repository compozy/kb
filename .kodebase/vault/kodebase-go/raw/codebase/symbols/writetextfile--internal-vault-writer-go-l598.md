---
blast_radius: 5
centrality: 0.2274
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 604
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 5
is_dead_export: false
is_long_function: false
language: "go"
loc: 7
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/vault/writer.go"
stage: "raw"
start_line: 598
symbol_kind: "function"
symbol_name: "writeTextFile"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: writeTextFile"
type: "source"
---

# Codebase Symbol: writeTextFile

Source file: [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 5
- External references: 0
- Centrality: 0.2274
- LOC: 7
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func writeTextFile(targetPath, body string) error {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/appendlog--internal-vault-writer-go-l529]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/ensuregitkeep--internal-vault-writer-go-l473]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/writefilesinbatches--internal-vault-writer-go-l306]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/writevault--internal-vault-writer-go-l53]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]] via `contains` (syntactic)
