---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 3
domain: "kodebase-go"
end_line: 129
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 14
outgoing_relation_count: 2
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/vault/render_test.go"
stage: "raw"
start_line: 116
symbol_kind: "function"
symbol_name: "TestRenderDocumentsCodebaseOverviewContainsSummary"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestRenderDocumentsCodebaseOverviewContainsSummary"
type: "source"
---

# Codebase Symbol: TestRenderDocumentsCodebaseOverviewContainsSummary

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_test.go|internal/vault/render_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 3
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 14
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func TestRenderDocumentsCodebaseOverviewContainsSummary(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/finddocument--internal-vault-render-test-go-l240]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderfixturedocuments--internal-vault-render-test-go-l232]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/render_test.go|internal/vault/render_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_test.go|internal/vault/render_test.go]] via `exports` (syntactic)
