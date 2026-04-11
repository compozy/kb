---
blast_radius: 33
centrality: 0.6059
cyclomatic_complexity: 6
domain: "kodebase-go"
end_line: 34
exported: true
external_reference_count: 2
has_smells: true
incoming_relation_count: 10
is_dead_export: false
is_long_function: false
language: "go"
loc: 20
outgoing_relation_count: 0
smells:
  - "bottleneck"
  - "high-blast-radius"
source_kind: "codebase-symbol"
source_path: "internal/vault/pathutils.go"
stage: "raw"
start_line: 15
symbol_kind: "function"
symbol_name: "ToPosixPath"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: ToPosixPath"
type: "source"
---

# Codebase Symbol: ToPosixPath

Source file: [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 6
- Long function: false
- Blast radius: 33
- External references: 2
- Centrality: 0.6059
- LOC: 20
- Dead export: false
- Smells: `bottleneck`, `high-blast-radius`

## Signature
```text
func ToPosixPath(value string) string {
```

## Documentation
ToPosixPath normalizes path separators to forward slashes and trims trailing slashes.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/cleancomparablepath--internal-vault-pathutils-go-l237]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createfileid--internal-vault-pathutils-go-l77]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/createsymbolid--internal-vault-pathutils-go-l169]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/derivetopicslug--internal-vault-pathutils-go-l140]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/getrawsymboldocumentpath--internal-vault-pathutils-go-l186]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/normalizedocumentpathsegment--internal-vault-pathutils-go-l250]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/readvaultsnapshot--internal-vault-reader-go-l62]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/cleantopicrelativepath--internal-vault-writer-go-l270]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] via `exports` (syntactic)
