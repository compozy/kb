---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 7
domain: "kodebase-go"
end_line: 219
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: false
language: "go"
loc: 24
outgoing_relation_count: 2
smells:
  - "dead-export"
source_kind: "codebase-symbol"
source_path: "internal/vault/render_test.go"
stage: "raw"
start_line: 196
symbol_kind: "function"
symbol_name: "TestRenderDocumentsBodiesHaveValidFrontmatterAndKinds"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: TestRenderDocumentsBodiesHaveValidFrontmatterAndKinds"
type: "source"
---

# Codebase Symbol: TestRenderDocumentsBodiesHaveValidFrontmatterAndKinds

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render_test.go|internal/vault/render_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 7
- Long function: false
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 24
- Dead export: true
- Smells: `dead-export`

## Signature
```text
func TestRenderDocumentsBodiesHaveValidFrontmatterAndKinds(t *testing.T) {
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsefrontmatter--internal-vault-render-test-go-l253]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderfixturedocuments--internal-vault-render-test-go-l232]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/render_test.go|internal/vault/render_test.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render_test.go|internal/vault/render_test.go]] via `exports` (syntactic)
