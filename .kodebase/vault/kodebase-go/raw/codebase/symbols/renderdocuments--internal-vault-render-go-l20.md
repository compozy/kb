---
blast_radius: 0
centrality: 0.0507
cyclomatic_complexity: 7
domain: "kodebase-go"
end_line: 128
exported: true
external_reference_count: 0
has_smells: true
incoming_relation_count: 2
is_dead_export: true
is_long_function: true
language: "go"
loc: 109
outgoing_relation_count: 26
smells:
  - "dead-export"
  - "long-function"
source_kind: "codebase-symbol"
source_path: "internal/vault/render.go"
stage: "raw"
start_line: 20
symbol_kind: "function"
symbol_name: "RenderDocuments"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: RenderDocuments"
type: "source"
---

# Codebase Symbol: RenderDocuments

Source file: [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 7
- Long function: true
- Blast radius: 0
- External references: 0
- Centrality: 0.0507
- LOC: 109
- Dead export: true
- Smells: `dead-export`, `long-function`

## Signature
```text
func RenderDocuments(
```

## Documentation
None

## Outgoing Relations
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createdocumentlookup--internal-vault-render-go-l177]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/createexternalnodelookup--internal-vault-render-go-l191]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/defaultdirectorymetrics--internal-vault-render-go-l678]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/defaultfilemetrics--internal-vault-render-go-l667]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/defaultsymbolmetrics--internal-vault-render-go-l671]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/groupfilesbydirectory--internal-vault-render-go-l743]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/groupfilesbylanguage--internal-vault-render-go-l759]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/grouprelationsbysource--internal-vault-render-go-l719]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/grouprelationsbytarget--internal-vault-render-go-l727]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/groupsymbolsbydirectory--internal-vault-render-go-l751]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/groupsymbolsbyfile--internal-vault-render-go-l735]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/groupsymbolsbykind--internal-vault-render-go-l775]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/groupsymbolsbylanguage--internal-vault-render-go-l767]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/rendermarkdowndocument--internal-vault-render-go-l164]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderrawdirectoryindex--internal-vault-render-go-l524]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderrawfiledocument--internal-vault-render-go-l312]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderrawlanguageindex--internal-vault-render-go-l600]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderrawsymboldocument--internal-vault-render-go-l449]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sortedmapkeys--internal-vault-render-go-l798]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sortsymbolsbylocation--internal-vault-render-go-l783]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/toprelationhotspotfiles--internal-vault-render-go-l682]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/buildstarterwikiarticles--internal-vault-render-wiki-go-l19]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderconceptindex--internal-vault-render-wiki-go-l134]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderdashboard--internal-vault-render-wiki-go-l75]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/rendersourceindex--internal-vault-render-wiki-go-l178]]
- `calls` (syntactic) -> [[kodebase-go/raw/codebase/symbols/renderwikiarticle--internal-vault-render-wiki-go-l65]]

## Backlinks
- [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]] via `contains` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]] via `exports` (syntactic)
