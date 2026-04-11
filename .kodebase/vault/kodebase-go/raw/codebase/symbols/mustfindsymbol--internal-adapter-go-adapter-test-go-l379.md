---
blast_radius: 17
centrality: 0.3312
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 390
exported: false
external_reference_count: 9
has_smells: true
incoming_relation_count: 18
is_dead_export: false
is_long_function: false
language: "go"
loc: 12
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter_test.go"
stage: "raw"
start_line: 379
symbol_kind: "function"
symbol_name: "mustFindSymbol"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: mustFindSymbol"
type: "source"
---

# Codebase Symbol: mustFindSymbol

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter_test.go|internal/adapter/go_adapter_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 17
- External references: 9
- Centrality: 0.3312
- LOC: 12
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func mustFindSymbol(t *testing.T, symbols []models.SymbolNode, name string) models.SymbolNode {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testgoadapterbuildscrossfilecallrelations--internal-adapter-go-adapter-integration-test-go-l11]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgoadaptercomputescyclomaticcomplexity--internal-adapter-go-adapter-test-go-l218]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgoadapterextractscallrelationsfordirectidentifiers--internal-adapter-go-adapter-test-go-l196]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgoadapterextractsdoccomments--internal-adapter-go-adapter-test-go-l287]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgoadapterextractsmethodsignaturewithreceiver--internal-adapter-go-adapter-test-go-l136]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgoadapterextractstypedeclarations--internal-adapter-go-adapter-test-go-l108]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgoadaptermoduledocusesonlyleadingcomment--internal-adapter-go-adapter-test-go-l312]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgoadapterparsessimplegofile--internal-adapter-go-adapter-test-go-l29]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgoadaptersetsexportedflags--internal-adapter-go-adapter-test-go-l88]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadapterintegrationparsesmultifileproject--internal-adapter-ts-adapter-integration-test-go-l11]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadapterbuildsimportbindingsfordefaultnamedandnamespaceimports--internal-adapter-ts-adapter-test-go-l177]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadaptercomputescyclomaticcomplexity--internal-adapter-ts-adapter-test-go-l285]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadapterextractsclassandmethodsymbols--internal-adapter-ts-adapter-test-go-l140]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadapterhandlesnamedandstarreexports--internal-adapter-ts-adapter-test-go-l236]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadapterparsesjavascriptrequireimports--internal-adapter-ts-adapter-test-go-l99]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadapterparsessimpletypescriptfile--internal-adapter-ts-adapter-test-go-l29]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadapterparsestsxcomponent--internal-adapter-ts-adapter-test-go-l78]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter_test.go|internal/adapter/go_adapter_test.go]] via `contains` (syntactic)
