---
blast_radius: 4
centrality: 0.1069
cyclomatic_complexity: 7
domain: "kodebase-go"
end_line: 283
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 14
outgoing_relation_count: 1
smells:
  - "bottleneck"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/vault/writer.go"
stage: "raw"
start_line: 270
symbol_kind: "function"
symbol_name: "cleanTopicRelativePath"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: cleanTopicRelativePath"
type: "source"
---

# Codebase Symbol: cleanTopicRelativePath

Source file: [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 7
- Long function: false
- Blast radius: 4
- External references: 0
- Centrality: 0.1069
- LOC: 14
- Dead export: false
- Smells: `bottleneck`, `feature-envy`

## Signature
```text
func cleanTopicRelativePath(relativePath string) (string, error) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/toposixpath--internal-vault-pathutils-go-l15]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/validatebasefile--internal-vault-writer-go-l242]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/validaterendereddocument--internal-vault-writer-go-l200]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]] via `contains` (syntactic)
