---
blast_radius: 4
centrality: 0.1018
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 92
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: false
is_long_function: false
language: "go"
loc: 21
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/config/config.go"
stage: "raw"
start_line: 72
symbol_kind: "function"
symbol_name: "decodeFile"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: decodeFile"
type: "source"
---

# Codebase Symbol: decodeFile

Source file: [[kodebase-go/raw/codebase/files/internal/config/config.go|internal/config/config.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 4
- External references: 0
- Centrality: 0.1018
- LOC: 21
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func decodeFile(path string, cfg *Config) error {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/load--internal-config-config-go-l59]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/config/config.go|internal/config/config.go]] via `contains` (syntactic)
