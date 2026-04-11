---
blast_radius: 0
centrality: 0.0507
domain: "kodebase-go"
end_line: 26
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
language: "go"
loc: 4
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/config/config.go"
stage: "raw"
start_line: 23
symbol_kind: "struct"
symbol_name: "AppConfig"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "struct"
title: "Codebase Symbol: AppConfig"
type: "source"
---

# Codebase Symbol: AppConfig

Source file: [[kodebase-go/raw/codebase/files/internal/config/config.go|internal/config/config.go]]

## Kind
`struct`

## Static Analysis
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 4
- Dead export: true
- Smells: `dead-export`

## Signature
```text
AppConfig struct {
```

## Documentation
AppConfig contains the application identity and environment.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/config/config.go|internal/config/config.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/config/config.go|internal/config/config.go]] via `exports` (syntactic)
