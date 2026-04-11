---
blast_radius: 2
centrality: 0.0965
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 487
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 15
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/writer.go"
stage: "raw"
start_line: 473
symbol_kind: "function"
symbol_name: "ensureGitkeep"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: ensureGitkeep"
type: "source"
---

# Codebase Symbol: ensureGitkeep

Source file: [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.0965
- LOC: 15
- Dead export: false
- Smells: None

## Signature
```text
func ensureGitkeep(directoryPath string) error {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writetextfile--internal-vault-writer-go-l598]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/ensuretopicgitkeeps--internal-vault-writer-go-l454]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]] via `contains` (syntactic)
