---
blast_radius: 2
centrality: 0.137
cyclomatic_complexity: 5
domain: "kodebase-go"
end_line: 61
exported: true
external_reference_count: 2
has_smells: true
incoming_relation_count: 4
is_dead_export: false
is_long_function: false
language: "go"
loc: 15
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/config/env.go"
stage: "raw"
start_line: 47
symbol_kind: "function"
symbol_name: "LoadDotEnvIfPresent"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: LoadDotEnvIfPresent"
type: "source"
---

# Codebase Symbol: LoadDotEnvIfPresent

Source file: [[kodebase-go/raw/codebase/files/internal/config/env.go|internal/config/env.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 5
- Long function: false
- Blast radius: 2
- External references: 2
- Centrality: 0.137
- LOC: 15
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func LoadDotEnvIfPresent(path string) error {
```

## Documentation
LoadDotEnvIfPresent loads a local dotenv file without overriding env vars
already supplied by the shell or process manager.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testloaddotenvifpresentloadsvalues--internal-config-config-test-go-l157]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testloaddotenvifpresentmissingfileisok--internal-config-config-test-go-l174]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/config/env.go|internal/config/env.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/config/env.go|internal/config/env.go]] via `exports` (syntactic)
