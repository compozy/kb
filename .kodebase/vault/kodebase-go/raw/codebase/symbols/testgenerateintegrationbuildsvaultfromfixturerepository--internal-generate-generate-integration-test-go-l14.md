---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 11
domain: "kodebase-go"
end_line: 66
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 53
outgoing_relation_count: 1
smells:
  - "dead-export"
  - "feature-envy"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/generate/generate_integration_test.go"
stage: "raw"
start_line: 14
symbol_kind: "function"
symbol_name: "TestGenerateIntegrationBuildsVaultFromFixtureRepository"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestGenerateIntegrationBuildsVaultFromFixtureRepository"
type: "source"
---

# Codebase Symbol: TestGenerateIntegrationBuildsVaultFromFixtureRepository

Source file: [[kodebase-go/raw/codebase/files/internal/generate/generate_integration_test.go|internal/generate/generate_integration_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 11
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 53
- Dead export: true
- Smells: `dead-export`, `feature-envy`, `long-function`

## Signature
```text
func TestGenerateIntegrationBuildsVaultFromFixtureRepository(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newrunner--internal-generate-generate-go-l61]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/generate/generate_integration_test.go|internal/generate/generate_integration_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/generate/generate_integration_test.go|internal/generate/generate_integration_test.go]] via `exports` (syntactic)
