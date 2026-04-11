---
blast_radius: 5
centrality: 0.145
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 55
exported: true
external_reference_count: 1
has_smells: true
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 16
outgoing_relation_count: 1
smells:
  - "bottleneck"
  - "feature-envy"
source_kind: "codebase-symbol"
source_path: "internal/config/config.go"
stage: "raw"
start_line: 40
symbol_kind: "function"
symbol_name: "Default"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: Default"
type: "source"
---

# Codebase Symbol: Default

Source file: [[kodebase-go/raw/codebase/files/internal/config/config.go|internal/config/config.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 5
- External references: 1
- Centrality: 0.145
- LOC: 16
- Dead export: false
- Smells: `bottleneck`, `feature-envy`

## Signature
```text
func Default() Config {
```

## Documentation
Default returns a sane starting configuration.

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/loadsecretsfromenv--internal-config-env-go-l33]]

## Backlinks
- [[kodebase-go/raw/codebase/symbols/load--internal-config-config-go-l59]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testdefaultconfighasvaliddefaults--internal-config-config-test-go-l9]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/config/config.go|internal/config/config.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/config/config.go|internal/config/config.go]] via `exports` (syntactic)
