---
blast_radius: 23
centrality: 0.1332
cyclomatic_complexity: 17
domain: "kodebase-go"
end_line: 293
exported: false
external_reference_count: 6
has_smells: true
incoming_relation_count: 7
is_dead_export: false
is_long_function: true
language: "go"
loc: 40
outgoing_relation_count: 0
smells:
  - "bottleneck"
  - "high-blast-radius"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect.go"
stage: "raw"
start_line: 254
symbol_kind: "function"
symbol_name: "inspectFrontmatterFloat"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: inspectFrontmatterFloat"
type: "source"
---

# Codebase Symbol: inspectFrontmatterFloat

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 17
- Long function: true
- Blast radius: 23
- External references: 6
- Centrality: 0.1332
- LOC: 40
- Dead export: false
- Smells: `bottleneck`, `high-blast-radius`, `long-function`

## Signature
```text
func inspectFrontmatterFloat(document vault.VaultDocument, key string) float64 {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/toblastradiusoutput--internal-cli-inspect-blastradius-go-l47]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/tocirculardependencyrow--internal-cli-inspect-circulardeps-go-l145]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/tocouplingoutput--internal-cli-inspect-coupling-go-l37]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/tofilelookupoutput--internal-cli-inspect-file-go-l31]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testinspectfrontmatterhelpers--internal-cli-inspect-helpers-test-go-l51]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/tosymboldetailoutput--internal-cli-inspect-symbol-go-l96]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]] via `contains` (syntactic)
