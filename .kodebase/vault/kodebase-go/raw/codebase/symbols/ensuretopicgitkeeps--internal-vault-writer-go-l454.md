---
blast_radius: 1
centrality: 0.0538
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 471
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 18
outgoing_relation_count: 1
smells:
source_kind: "codebase-symbol"
source_path: "internal/vault/writer.go"
stage: "raw"
start_line: 454
symbol_kind: "function"
symbol_name: "ensureTopicGitkeeps"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: ensureTopicGitkeeps"
type: "source"
---

# Codebase Symbol: ensureTopicGitkeeps

Source file: [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 1
- External references: 0
- Centrality: 0.0538
- LOC: 18
- Dead export: false
- Smells: None

## Signature
```text
func ensureTopicGitkeeps(topicPath string) error {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/ensuregitkeep--internal-vault-writer-go-l473]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/writevault--internal-vault-writer-go-l53]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]] via `contains` (syntactic)
