---
blast_radius: 6
centrality: 0.2251
cyclomatic_complexity: 1
domain: "kodebase-go"
end_line: 38
exported: true
external_reference_count: 2
has_smells: true
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 6
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/config/env.go"
stage: "raw"
start_line: 33
symbol_kind: "function"
symbol_name: "LoadSecretsFromEnv"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: LoadSecretsFromEnv"
type: "source"
---

# Codebase Symbol: LoadSecretsFromEnv

Source file: [[kodebase-go/raw/codebase/files/internal/config/env.go|internal/config/env.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 1
- Long function: false
- Blast radius: 6
- External references: 2
- Centrality: 0.2251
- LOC: 6
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func LoadSecretsFromEnv() Secrets {
```

## Documentation
LoadSecretsFromEnv reads the current process environment into a stable
runtime value object.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/default--internal-config-config-go-l40]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/load--internal-config-config-go-l59]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/config/env.go|internal/config/env.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/config/env.go|internal/config/env.go]] via `exports` (syntactic)
