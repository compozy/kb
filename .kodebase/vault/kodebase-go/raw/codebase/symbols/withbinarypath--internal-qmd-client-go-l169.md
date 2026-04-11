---
blast_radius: 13
centrality: 0.2154
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 173
exported: true
external_reference_count: 13
has_smells: true
incoming_relation_count: 15
is_dead_export: false
is_long_function: false
language: "go"
loc: 5
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/qmd/client.go"
stage: "raw"
start_line: 169
symbol_kind: "function"
symbol_name: "WithBinaryPath"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: WithBinaryPath"
type: "source"
---

# Codebase Symbol: WithBinaryPath

Source file: [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 13
- External references: 13
- Centrality: 0.2154
- LOC: 5
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func WithBinaryPath(path string) ClientOption {
```

## Documentation
WithBinaryPath overrides the executable used for QMD invocations.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testqmdclientindexesandsearchestempvault--internal-qmd-client-integration-test-go-l13]] via `calls` (syntactic)
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
- [[kodebase-go/raw/codebase/symbols/testsearchreturnserrqmdunavailableformissingbinary--internal-qmd-client-test-go-l14]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] via `exports` (syntactic)
