---
blast_radius: 2
centrality: 0.0939
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 864
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 3
is_dead_export: false
is_long_function: false
language: "go"
loc: 21
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/cli/inspect_test.go"
stage: "raw"
start_line: 844
symbol_kind: "function"
symbol_name: "testFileDocumentForCycle"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: testFileDocumentForCycle"
type: "source"
---

# Codebase Symbol: testFileDocumentForCycle

Source file: [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 2
- External references: 0
- Centrality: 0.0939
- LOC: 21
- Dead export: false
- Smells: None

## Signature
```text
func testFileDocumentForCycle(relativePath, sourcePath string, frontmatter map[string]any, targets ...string) vault.VaultDocument {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testtocirculardepsoutputfallsbacktosccdetectionforlegacyvaults--internal-cli-inspect-test-go-l508]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testtocirculardepsoutputlistsfileswithcirculardependencyflag--internal-cli-inspect-test-go-l467]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]] via `contains` (syntactic)
