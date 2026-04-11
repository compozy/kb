---
blast_radius: 0
centrality: 0.0507
domain: "kodebase-go"
end_line: 53
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
language: "go"
loc: 1
outgoing_relation_count: 0
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/models/models.go"
stage: "raw"
start_line: 53
symbol_kind: "type"
symbol_name: "DiagnosticSeverity"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "type"
title: "Codebase Symbol: DiagnosticSeverity"
type: "source"
---

# Codebase Symbol: DiagnosticSeverity

Source file: [[kodebase-go/raw/codebase/files/internal/models/models.go|internal/models/models.go]]

## Kind
`type`

## Static Analysis
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 1
- Dead export: true
- Smells: `dead-export`

## Signature
```text
DiagnosticSeverity string
```

## Documentation
DiagnosticSeverity indicates the severity of a structured diagnostic.

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/models/models.go|internal/models/models.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/models/models.go|internal/models/models.go]] via `exports` (syntactic)
