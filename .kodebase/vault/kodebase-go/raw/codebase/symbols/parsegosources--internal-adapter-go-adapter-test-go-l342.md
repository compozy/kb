---
blast_radius: 12
centrality: 0.2939
cyclomatic_complexity: 6
domain: "kodebase-go"
end_line: 377
exported: false
external_reference_count: 1
has_smells: true
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 36
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter_test.go"
stage: "raw"
start_line: 342
symbol_kind: "function"
symbol_name: "parseGoSources"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: parseGoSources"
type: "source"
---

# Codebase Symbol: parseGoSources

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter_test.go|internal/adapter/go_adapter_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 6
- Long function: false
- Blast radius: 12
- External references: 1
- Centrality: 0.2939
- LOC: 36
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func parseGoSources(t *testing.T, sources map[string]string) []models.ParsedFile {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testgoadapterbuildscrossfilecallrelations--internal-adapter-go-adapter-integration-test-go-l11]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/parsesinglegofile--internal-adapter-go-adapter-test-go-l331]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgoadapterparsessimplegofile--internal-adapter-go-adapter-test-go-l29]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter_test.go|internal/adapter/go_adapter_test.go]] via `contains` (syntactic)
