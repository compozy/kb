---
blast_radius: 11
centrality: 0.1852
cyclomatic_complexity: 6
domain: "kodebase-go"
end_line: 535
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 12
is_dead_export: false
is_long_function: false
language: "go"
loc: 46
outgoing_relation_count: 1
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/qmd/client_test.go"
stage: "raw"
start_line: 490
symbol_kind: "function"
symbol_name: "writeFakeQMD"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: writeFakeQMD"
type: "source"
---

# Codebase Symbol: writeFakeQMD

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 6
- Long function: false
- Blast radius: 11
- External references: 0
- Centrality: 0.1852
- LOC: 46
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func writeFakeQMD(t *testing.T, options fakeQMDOptions) string {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/shellquote--internal-qmd-client-test-go-l572]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testindexaddconstructsexpectedarguments--internal-qmd-client-test-go-l27]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testindexrejectsinvalidinputs--internal-qmd-client-test-go-l150]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testindexupdateconstructsexpectedarguments--internal-qmd-client-test-go-l64]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testindexwithcontextandembedrunsexpectedcommands--internal-qmd-client-test-go-l97]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testsearchallomitslimitandacceptsmodealias--internal-qmd-client-test-go-l216]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testsearchcontextcancellationstopsrunningcommand--internal-qmd-client-test-go-l343]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testsearchfailureincludesstderrdiagnostics--internal-qmd-client-test-go-l365]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testsearchfullusesbodywhenpresent--internal-qmd-client-test-go-l318]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testsearchhybridmodeusesquerycommand--internal-qmd-client-test-go-l183]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testsearchparsesjsonandnormalizesresults--internal-qmd-client-test-go-l283]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testsearchpasseslimitminscoreandfullflags--internal-qmd-client-test-go-l246]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]] via `contains` (syntactic)
