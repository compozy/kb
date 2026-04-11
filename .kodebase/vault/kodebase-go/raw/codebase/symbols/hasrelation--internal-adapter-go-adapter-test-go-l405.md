---
blast_radius: 11
centrality: 0.1946
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 413
exported: false
external_reference_count: 8
has_smells: true
incoming_relation_count: 12
is_dead_export: false
is_long_function: false
language: "go"
loc: 9
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter_test.go"
stage: "raw"
start_line: 405
symbol_kind: "function"
symbol_name: "hasRelation"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: hasRelation"
type: "source"
---

# Codebase Symbol: hasRelation

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter_test.go|internal/adapter/go_adapter_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 11
- External references: 8
- Centrality: 0.1946
- LOC: 9
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func hasRelation(relations []models.RelationEdge, fromID string, toID string, relationType models.RelationType) bool {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testgoadapterbuildscrossfilecallrelations--internal-adapter-go-adapter-integration-test-go-l11]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgoadapterextractscallrelationsfordirectidentifiers--internal-adapter-go-adapter-test-go-l196]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgoadapterextractsimportrelations--internal-adapter-go-adapter-test-go-l157]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgoadapterparsessimplegofile--internal-adapter-go-adapter-test-go-l29]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadapterintegrationparsesmultifileproject--internal-adapter-ts-adapter-integration-test-go-l11]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadapterbuildsimportbindingsfordefaultnamedandnamespaceimports--internal-adapter-ts-adapter-test-go-l177]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadapterextractsclassandmethodsymbols--internal-adapter-ts-adapter-test-go-l140]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadapterhandlesnamedandstarreexports--internal-adapter-ts-adapter-test-go-l236]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadapterparsesjavascriptrequireimports--internal-adapter-ts-adapter-test-go-l99]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadapterparsessimpletypescriptfile--internal-adapter-ts-adapter-test-go-l29]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtsadapterparsestsxcomponent--internal-adapter-ts-adapter-test-go-l78]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter_test.go|internal/adapter/go_adapter_test.go]] via `contains` (syntactic)
