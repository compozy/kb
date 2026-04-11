---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 15
domain: "kodebase-go"
end_line: 111
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 59
outgoing_relation_count: 14
smells:
  - "dead-export"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/vault/writer.go"
stage: "raw"
start_line: 53
symbol_kind: "function"
symbol_name: "WriteVault"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: WriteVault"
type: "source"
---

# Codebase Symbol: WriteVault

Source file: [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 15
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 59
- Dead export: true
- Smells: `dead-export`, `long-function`

## Signature
```text
func WriteVault(ctx context.Context, options WriteVaultOptions) (WriteVaultResult, error) {
```

## Documentation
WriteVault persists the rendered markdown and base files for a topic.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/appendlog--internal-vault-writer-go-l529]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/buildtopicclaude--internal-vault-writer-go-l359]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/buildwriterequests--internal-vault-writer-go-l166]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/countwrittendocuments--internal-vault-writer-go-l581]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/ensureagentssymlink--internal-vault-writer-go-l442]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/ensuredirectories--internal-vault-writer-go-l285]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/ensuretopicgitkeeps--internal-vault-writer-go-l454]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/ensuretopicskeleton--internal-vault-writer-go-l126]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newwriteprogressreporter--internal-vault-writer-go-l339]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/removemanagedwikiconcepts--internal-vault-writer-go-l489]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resetmanagedsubtrees--internal-vault-writer-go-l149]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/validatetopic--internal-vault-writer-go-l113]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writefilesinbatches--internal-vault-writer-go-l306]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writetextfile--internal-vault-writer-go-l598]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]] via `exports` (syntactic)
