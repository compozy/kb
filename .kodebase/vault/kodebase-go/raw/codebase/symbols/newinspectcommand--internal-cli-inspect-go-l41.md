---
blast_radius: 24
centrality: 0.2208
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 68
exported: false
external_reference_count: 1
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 28
outgoing_relation_count: 11
smells:
  - "bottleneck"
  - "high-blast-radius"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect.go"
stage: "raw"
start_line: 41
symbol_kind: "function"
symbol_name: "newInspectCommand"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: newInspectCommand"
type: "source"
---

# Codebase Symbol: newInspectCommand

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 24
- External references: 1
- Centrality: 0.2208
- LOC: 28
- Dead export: false
- Smells: `bottleneck`, `high-blast-radius`

## Signature
```text
func newInspectCommand() *cobra.Command {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/bindinspectsharedflags--internal-cli-inspect-go-l70]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newinspectbacklinkscommand--internal-cli-inspect-backlinks-go-l11]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newinspectblastradiuscommand--internal-cli-inspect-blastradius-go-l20]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newinspectcirculardepscommand--internal-cli-inspect-circulardeps-go-l12]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newinspectcomplexitycommand--internal-cli-inspect-complexity-go-l21]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newinspectcouplingcommand--internal-cli-inspect-coupling-go-l19]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newinspectdeadcodecommand--internal-cli-inspect-deadcode-go-l19]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newinspectdepscommand--internal-cli-inspect-deps-go-l11]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newinspectfilecommand--internal-cli-inspect-file-go-l11]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newinspectsmellscommand--internal-cli-inspect-smells-go-l19]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newinspectsymbolcommand--internal-cli-inspect-symbol-go-l21]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/newrootcommand--internal-cli-root-go-l14]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]] via `contains` (syntactic)
