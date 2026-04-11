---
blast_radius: 27
centrality: 0.1985
cyclomatic_complexity: 17
domain: "kodebase-go"
end_line: 252
exported: false
external_reference_count: 8
has_smells: true
incoming_relation_count: 10
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
start_line: 213
symbol_kind: "function"
symbol_name: "inspectFrontmatterInt"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: inspectFrontmatterInt"
type: "source"
---

# Codebase Symbol: inspectFrontmatterInt

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 17
- Long function: true
- Blast radius: 27
- External references: 8
- Centrality: 0.1985
- LOC: 40
- Dead export: false
- Smells: `bottleneck`, `high-blast-radius`, `long-function`

## Signature
```text
func inspectFrontmatterInt(document vault.VaultDocument, key string) int {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/inspectsymbolsforfile--internal-cli-inspect-go-l411]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/toblastradiusoutput--internal-cli-inspect-blastradius-go-l47]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/tocirculardependencyrow--internal-cli-inspect-circulardeps-go-l145]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/tocomplexityoutput--internal-cli-inspect-complexity-go-l43]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/tocouplingoutput--internal-cli-inspect-coupling-go-l37]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/tofilelookupoutput--internal-cli-inspect-file-go-l31]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testinspectfrontmatterhelpers--internal-cli-inspect-helpers-test-go-l51]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/tosymboldetailoutput--internal-cli-inspect-symbol-go-l96]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/tosymbolsummaryoutput--internal-cli-inspect-symbol-go-l55]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]] via `contains` (syntactic)
