---
blast_radius: 20
centrality: 0.1487
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 211
exported: false
external_reference_count: 6
has_smells: true
incoming_relation_count: 7
is_dead_export: false
is_long_function: false
language: "go"
loc: 19
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect.go"
stage: "raw"
start_line: 193
symbol_kind: "function"
symbol_name: "inspectFrontmatterBool"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: inspectFrontmatterBool"
type: "source"
---

# Codebase Symbol: inspectFrontmatterBool

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 20
- External references: 6
- Centrality: 0.1487
- LOC: 19
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func inspectFrontmatterBool(document vault.VaultDocument, key string) bool {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/circulardependencyrowsfromflags--internal-cli-inspect-circulardeps-go-l56]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/tocouplingoutput--internal-cli-inspect-coupling-go-l37]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/todeadcodeoutput--internal-cli-inspect-deadcode-go-l32]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/tofilelookupoutput--internal-cli-inspect-file-go-l31]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testinspectfrontmatterhelpers--internal-cli-inspect-helpers-test-go-l51]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/tosymboldetailoutput--internal-cli-inspect-symbol-go-l96]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]] via `contains` (syntactic)
