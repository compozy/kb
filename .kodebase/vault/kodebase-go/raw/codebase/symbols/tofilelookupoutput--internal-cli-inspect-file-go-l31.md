---
blast_radius: 2
centrality: 0.1083
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 53
exported: false
external_reference_count: 2
has_smells: true
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 23
outgoing_relation_count: 9
smells:
  - "bottleneck"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_file.go"
stage: "raw"
start_line: 31
symbol_kind: "function"
symbol_name: "toFileLookupOutput"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: toFileLookupOutput"
type: "source"
---

# Codebase Symbol: toFileLookupOutput

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_file.go|internal/cli/inspect_file.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 2
- External references: 2
- Centrality: 0.1083
- LOC: 23
- Dead export: false
- Smells: `bottleneck`, `feature-envy`

## Signature
```text
func toFileLookupOutput(snapshot vault.VaultSnapshot, sourcePath string) (inspectOutput, error) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createinspectdetailoutput--internal-cli-inspect-go-l300]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createinspectrelationrows--internal-cli-inspect-go-l315]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/findinspectfilebysourcepath--internal-cli-inspect-go-l339]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterbool--internal-cli-inspect-go-l193]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterfloat--internal-cli-inspect-go-l254]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterint--internal-cli-inspect-go-l213]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstring--internal-cli-inspect-go-l156]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstringarray--internal-cli-inspect-go-l170]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectsymbolsforfile--internal-cli-inspect-go-l411]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testtofilelookupoutputincludescontainedsymbolsandmetrics--internal-cli-inspect-test-go-l312]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtofilelookupoutputreturnsdescriptiveerrorforunknownpath--internal-cli-inspect-test-go-l373]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_file.go|internal/cli/inspect_file.go]] via `contains` (syntactic)
