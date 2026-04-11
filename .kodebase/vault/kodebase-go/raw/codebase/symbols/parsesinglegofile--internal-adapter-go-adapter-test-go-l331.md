---
blast_radius: 9
centrality: 0.2521
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 340
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 10
is_dead_export: false
is_long_function: false
language: "go"
loc: 10
outgoing_relation_count: 1
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/adapter/go_adapter_test.go"
stage: "raw"
start_line: 331
symbol_kind: "function"
symbol_name: "parseSingleGoFile"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: parseSingleGoFile"
type: "source"
---

# Codebase Symbol: parseSingleGoFile

Source file: [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter_test.go|internal/adapter/go_adapter_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 9
- External references: 0
- Centrality: 0.2521
- LOC: 10
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func parseSingleGoFile(t *testing.T, source string) models.ParsedFile {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsegosources--internal-adapter-go-adapter-test-go-l342]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testgoadaptercomputescyclomaticcomplexity--internal-adapter-go-adapter-test-go-l218]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgoadapterextractscallrelationsfordirectidentifiers--internal-adapter-go-adapter-test-go-l196]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgoadapterextractsdoccomments--internal-adapter-go-adapter-test-go-l287]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgoadapterextractsimportrelations--internal-adapter-go-adapter-test-go-l157]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgoadapterextractsmethodsignaturewithreceiver--internal-adapter-go-adapter-test-go-l136]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgoadapterextractstypedeclarations--internal-adapter-go-adapter-test-go-l108]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgoadaptermoduledocusesonlyleadingcomment--internal-adapter-go-adapter-test-go-l312]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgoadapterproducesdiagnosticsforparseerrors--internal-adapter-go-adapter-test-go-l260]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testgoadaptersetsexportedflags--internal-adapter-go-adapter-test-go-l88]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter_test.go|internal/adapter/go_adapter_test.go]] via `contains` (syntactic)
