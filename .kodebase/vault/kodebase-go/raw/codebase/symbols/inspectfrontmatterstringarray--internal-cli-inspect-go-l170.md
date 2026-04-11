---
blast_radius: 31
centrality: 0.2079
cyclomatic_complexity: 7
domain: "kodebase-go"
end_line: 191
exported: false
external_reference_count: 10
has_smells: true
incoming_relation_count: 11
is_dead_export: false
is_long_function: false
language: "go"
loc: 22
outgoing_relation_count: 0
smells:
  - "bottleneck"
  - "high-blast-radius"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect.go"
stage: "raw"
start_line: 170
symbol_kind: "function"
symbol_name: "inspectFrontmatterStringArray"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: inspectFrontmatterStringArray"
type: "source"
---

# Codebase Symbol: inspectFrontmatterStringArray

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 7
- Long function: false
- Blast radius: 31
- External references: 10
- Centrality: 0.2079
- LOC: 22
- Dead export: false
- Smells: `bottleneck`, `high-blast-radius`

## Signature
```text
func inspectFrontmatterStringArray(document vault.VaultDocument, key string) []string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/toblastradiusoutput--internal-cli-inspect-blastradius-go-l47]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/tocirculardependencyrow--internal-cli-inspect-circulardeps-go-l145]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/tocomplexityoutput--internal-cli-inspect-complexity-go-l43]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/tocouplingoutput--internal-cli-inspect-coupling-go-l37]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/todeadcodeoutput--internal-cli-inspect-deadcode-go-l32]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/tofilelookupoutput--internal-cli-inspect-file-go-l31]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testinspectfrontmatterhelpers--internal-cli-inspect-helpers-test-go-l51]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/tosmelloutput--internal-cli-inspect-smells-go-l37]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/tosymboldetailoutput--internal-cli-inspect-symbol-go-l96]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/tosymbolsummaryoutput--internal-cli-inspect-symbol-go-l55]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]] via `contains` (syntactic)
