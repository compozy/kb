---
blast_radius: 0
centrality: 0.0507
domain: "kodebase-go"
end_line: 20
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
language: "go"
loc: 6
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/config/config.go"
stage: "raw"
start_line: 15
symbol_kind: "struct"
symbol_name: "Config"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "struct"
title: "Codebase Symbol: Config"
type: "source"
---

# Codebase Symbol: Config

Source file: [[kodebase-go/raw/codebase/files/internal/config/config.go|internal/config/config.go]]

## Kind
`struct`

## Static Analysis
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 6
- Dead export: true
- Smells: `dead-export`

## Signature
```text
Config struct {
```

## Documentation
Config contains the complete TOML-backed runtime configuration plus
non-TOML runtime helpers such as secrets.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/config/config.go|internal/config/config.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/config/config.go|internal/config/config.go]] via `exports` (syntactic)
