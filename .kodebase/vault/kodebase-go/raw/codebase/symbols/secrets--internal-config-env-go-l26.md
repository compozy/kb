---
blast_radius: 0
centrality: 0.0507
domain: "kodebase-go"
end_line: 29
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
source_path: "internal/config/env.go"
stage: "raw"
start_line: 26
symbol_kind: "struct"
symbol_name: "Secrets"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "struct"
title: "Codebase Symbol: Secrets"
type: "source"
---

# Codebase Symbol: Secrets

Source file: [[kodebase-go/raw/codebase/files/internal/config/env.go|internal/config/env.go]]

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
Secrets struct {
```

## Documentation
Secrets contains the runtime secrets loaded from the environment.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/config/env.go|internal/config/env.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/config/env.go|internal/config/env.go]] via `exports` (syntactic)
