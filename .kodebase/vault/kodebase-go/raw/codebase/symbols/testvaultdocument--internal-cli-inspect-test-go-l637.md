---
blast_radius: 9
centrality: 0.2305
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 639
exported: false
external_reference_count: 4
has_smells: true
incoming_relation_count: 10
is_dead_export: false
is_long_function: false
language: "go"
loc: 3
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_test.go"
stage: "raw"
start_line: 637
symbol_kind: "function"
symbol_name: "testVaultDocument"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: testVaultDocument"
type: "source"
---

# Codebase Symbol: testVaultDocument

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 9
- External references: 4
- Centrality: 0.2305
- LOC: 3
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func testVaultDocument(frontmatter map[string]any) vault.VaultDocument {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testinspectfrontmatterhelpers--internal-cli-inspect-helpers-test-go-l51]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtoblastradiusoutputrespectsminimumandtop--internal-cli-inspect-helpers-test-go-l178]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtocouplingoutputfiltersunstableonly--internal-cli-inspect-helpers-test-go-l198]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtosmelloutputfiltersbytype--internal-cli-inspect-helpers-test-go-l152]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtoblastradiusoutputsortsbydescendingblastradius--internal-cli-inspect-test-go-l132]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtocomplexityoutputsortsbydescendingcomplexity--internal-cli-inspect-test-go-l90]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtocouplingoutputsortsbyinstability--internal-cli-inspect-test-go-l166]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtodeadcodeoutputlistsdeadexportsandorphanfiles--internal-cli-inspect-test-go-l55]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtosmelloutputlistssymbolsandfileswithsmells--internal-cli-inspect-test-go-l16]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]] via `contains` (syntactic)
