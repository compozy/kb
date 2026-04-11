---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 19
domain: "kodebase-go"
end_line: 150
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 100
outgoing_relation_count: 6
smells:
  - "dead-export"
  - "feature-envy"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_helpers_test.go"
stage: "raw"
start_line: 51
symbol_kind: "function"
symbol_name: "TestInspectFrontmatterHelpers"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestInspectFrontmatterHelpers"
type: "source"
---

# Codebase Symbol: TestInspectFrontmatterHelpers

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_helpers_test.go|internal/cli/inspect_helpers_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 19
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 100
- Dead export: true
- Smells: `dead-export`, `feature-envy`, `long-function`

## Signature
```text
func TestInspectFrontmatterHelpers(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterbool--internal-cli-inspect-go-l193]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterfloat--internal-cli-inspect-go-l254]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterint--internal-cli-inspect-go-l213]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstring--internal-cli-inspect-go-l156]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstringarray--internal-cli-inspect-go-l170]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testvaultdocument--internal-cli-inspect-test-go-l637]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_helpers_test.go|internal/cli/inspect_helpers_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_helpers_test.go|internal/cli/inspect_helpers_test.go]] via `exports` (syntactic)
