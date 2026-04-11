---
blast_radius: 2
centrality: 0.0939
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 85
exported: false
external_reference_count: 2
has_smells: false
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 49
outgoing_relation_count: 4
smells:
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_smells.go"
stage: "raw"
start_line: 37
symbol_kind: "function"
symbol_name: "toSmellOutput"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: toSmellOutput"
type: "source"
---

# Codebase Symbol: toSmellOutput

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_smells.go|internal/cli/inspect_smells.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 2
- External references: 2
- Centrality: 0.0939
- LOC: 49
- Dead export: false
- Smells: None

## Signature
```text
func toSmellOutput(snapshot vault.VaultSnapshot, smellType string) inspectOutput {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstring--internal-cli-inspect-go-l156]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstringarray--internal-cli-inspect-go-l170]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/includesmellrow--internal-cli-inspect-smells-go-l87]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/smellrowstomaps--internal-cli-inspect-smells-go-l104]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testtosmelloutputfiltersbytype--internal-cli-inspect-helpers-test-go-l152]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtosmelloutputlistssymbolsandfileswithsmells--internal-cli-inspect-test-go-l16]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_smells.go|internal/cli/inspect_smells.go]] via `contains` (syntactic)
