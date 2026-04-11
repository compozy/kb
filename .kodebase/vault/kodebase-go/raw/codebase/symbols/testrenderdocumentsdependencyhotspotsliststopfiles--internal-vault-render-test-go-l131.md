---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 140
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 10
outgoing_relation_count: 2
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/vault/render_test.go"
stage: "raw"
start_line: 131
symbol_kind: "function"
symbol_name: "TestRenderDocumentsDependencyHotspotsListsTopFiles"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestRenderDocumentsDependencyHotspotsListsTopFiles"
type: "source"
---

# Codebase Symbol: TestRenderDocumentsDependencyHotspotsListsTopFiles

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_test.go|internal/vault/render_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 10
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func TestRenderDocumentsDependencyHotspotsListsTopFiles(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/finddocument--internal-vault-render-test-go-l240]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderfixturedocuments--internal-vault-render-test-go-l232]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/render_test.go|internal/vault/render_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_test.go|internal/vault/render_test.go]] via `exports` (syntactic)
