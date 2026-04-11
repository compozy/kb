---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 4
domain: "kodebase-go"
end_line: 536
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 29
outgoing_relation_count: 2
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_test.go"
stage: "raw"
start_line: 508
symbol_kind: "function"
symbol_name: "TestToCircularDepsOutputFallsBackToSCCDetectionForLegacyVaults"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestToCircularDepsOutputFallsBackToSCCDetectionForLegacyVaults"
type: "source"
---

# Codebase Symbol: TestToCircularDepsOutputFallsBackToSCCDetectionForLegacyVaults

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 4
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 29
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func TestToCircularDepsOutputFallsBackToSCCDetectionForLegacyVaults(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/tocirculardepsoutput--internal-cli-inspect-circulardeps-go-l27]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testfiledocumentforcycle--internal-cli-inspect-test-go-l844]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]] via `exports` (syntactic)
