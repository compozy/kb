---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 181
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 8
outgoing_relation_count: 1
smells:
  - "dead-export"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/config/config_test.go"
stage: "raw"
start_line: 174
symbol_kind: "function"
symbol_name: "TestLoadDotEnvIfPresentMissingFileIsOK"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestLoadDotEnvIfPresentMissingFileIsOK"
type: "source"
---

# Codebase Symbol: TestLoadDotEnvIfPresentMissingFileIsOK

Source file: [[kodebase-go/raw/codebase/files/internal/config/config_test.go|internal/config/config_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 8
- Dead export: true
- Smells: `dead-export`, `feature-envy`

## Signature
```text
func TestLoadDotEnvIfPresentMissingFileIsOK(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/loaddotenvifpresent--internal-config-env-go-l47]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/config/config_test.go|internal/config/config_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/config/config_test.go|internal/config/config_test.go]] via `exports` (syntactic)
