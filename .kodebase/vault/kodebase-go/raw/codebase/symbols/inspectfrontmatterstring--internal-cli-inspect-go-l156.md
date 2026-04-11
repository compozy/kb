---
blast_radius: 42
centrality: 0.623
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 168
exported: false
external_reference_count: 12
has_smells: true
incoming_relation_count: 17
is_dead_export: false
is_long_function: false
language: "go"
loc: 13
outgoing_relation_count: 0
smells:
  - "bottleneck"
  - "high-blast-radius"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect.go"
stage: "raw"
start_line: 156
symbol_kind: "function"
symbol_name: "inspectFrontmatterString"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: inspectFrontmatterString"
type: "source"
---

# Codebase Symbol: inspectFrontmatterString

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 42
- External references: 12
- Centrality: 0.623
- LOC: 13
- Dead export: false
- Smells: `bottleneck`, `high-blast-radius`

## Signature
```text
func inspectFrontmatterString(document vault.VaultDocument, key string) string {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/findinspectfilebysourcepath--internal-cli-inspect-go-l339]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/findsingleinspectsymbolmatch--internal-cli-inspect-go-l350]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/inspectsymbolsforfile--internal-cli-inspect-go-l411]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/isfunctionlikedocument--internal-cli-inspect-go-l147]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/toblastradiusoutput--internal-cli-inspect-blastradius-go-l47]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/buildinspectimportadjacency--internal-cli-inspect-circulardeps-go-l103]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/circulardependencyrowsfromfallback--internal-cli-inspect-circulardeps-go-l74]] via `calls` (syntactic)
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
