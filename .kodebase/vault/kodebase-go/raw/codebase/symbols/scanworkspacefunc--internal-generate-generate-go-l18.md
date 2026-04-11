---
blast_radius: 0
centrality: 0.0507
domain: "kodebase-go"
end_line: 18
exported: false
external_reference_count: 0
has_smells: false
incoming_relation_count: 1
is_dead_export: false
language: "go"
loc: 1
outgoing_relation_count: 0
smells:
source_kind: "codebase-symbol"
source_path: "internal/generate/generate.go"
stage: "raw"
start_line: 18
symbol_kind: "type"
symbol_name: "scanWorkspaceFunc"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "type"
title: "Codebase Symbol: scanWorkspaceFunc"
type: "source"
---

# Codebase Symbol: scanWorkspaceFunc

Source file: [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]]

## Kind
`type`

## Static Analysis
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 1
- Dead export: false
- Smells: None

## Signature
```text
scanWorkspaceFunc func(rootPath string, opts ...scanner.Option) (*models.ScannedWorkspace, error)
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]] via `contains` (syntactic)
