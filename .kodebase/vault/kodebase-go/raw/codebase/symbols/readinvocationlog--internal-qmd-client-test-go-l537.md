---
blast_radius: 6
centrality: 0.1133
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 570
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 7
is_dead_export: false
is_long_function: false
language: "go"
loc: 34
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/qmd/client_test.go"
stage: "raw"
start_line: 537
symbol_kind: "function"
symbol_name: "readInvocationLog"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: readInvocationLog"
type: "source"
---

# Codebase Symbol: readInvocationLog

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 6
- External references: 0
- Centrality: 0.1133
- LOC: 34
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func readInvocationLog(t *testing.T, path string) [][]string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testindexaddconstructsexpectedarguments--internal-qmd-client-test-go-l27]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testindexupdateconstructsexpectedarguments--internal-qmd-client-test-go-l64]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testindexwithcontextandembedrunsexpectedcommands--internal-qmd-client-test-go-l97]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testsearchallomitslimitandacceptsmodealias--internal-qmd-client-test-go-l216]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testsearchhybridmodeusesquerycommand--internal-qmd-client-test-go-l183]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testsearchpasseslimitminscoreandfullflags--internal-qmd-client-test-go-l246]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]] via `contains` (syntactic)
