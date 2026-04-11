---
blast_radius: 11
centrality: 0.2453
cyclomatic_complexity: 6
domain: "kodebase-go"
end_line: 406
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 17
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/adapter/ts_adapter_test.go"
stage: "raw"
start_line: 390
symbol_kind: "function"
symbol_name: "languageForPath"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: languageForPath"
type: "source"
---

# Codebase Symbol: languageForPath

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter_test.go|internal/adapter/ts_adapter_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 6
- Long function: false
- Blast radius: 11
- External references: 0
- Centrality: 0.2453
- LOC: 17
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func languageForPath(t *testing.T, relativePath string) models.SupportedLanguage {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/parsetslikesources--internal-adapter-ts-adapter-test-go-l340]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter_test.go|internal/adapter/ts_adapter_test.go]] via `contains` (syntactic)
