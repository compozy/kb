---
blast_radius: 7
centrality: 0.2449
cyclomatic_complexity: 2
domain: "kodebase-go"
end_line: 212
exported: false
external_reference_count: 0
has_smells: true
incoming_relation_count: 8
is_dead_export: false
is_long_function: false
language: "go"
loc: 7
outgoing_relation_count: 0
smells:
  - "bottleneck"
source_kind: "codebase-symbol"
source_path: "internal/vault/query_test.go"
stage: "raw"
start_line: 206
symbol_kind: "function"
symbol_name: "mkdirAll"
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "symbol"
  - "go"
  - "function"
title: "Codebase Symbol: mkdirAll"
type: "source"
---

# Codebase Symbol: mkdirAll

Source file: [[kodebase-go/raw/codebase/files/internal/vault/query_test.go|internal/vault/query_test.go]]

## Kind
`function`

## Static Analysis
- Cyclomatic complexity: 2
- Long function: false
- Blast radius: 7
- External references: 0
- Centrality: 0.2449
- LOC: 7
- Dead export: false
- Smells: `bottleneck`

## Signature
```text
func mkdirAll(t *testing.T, path string) {
```

## Documentation
None

## Outgoing Relations
None

## Backlinks
- [[kodebase-go/raw/codebase/symbols/testlistavailabletopicsreturnssortedtopics--internal-vault-query-test-go-l178]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryautoresolvessingletopic--internal-vault-query-test-go-l72]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryerrorswhenexplicittopicismissing--internal-vault-query-test-go-l144]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryerrorswhenmultipletopicsexist--internal-vault-query-test-go-l95]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryfindsvaultbywalkingup--internal-vault-query-test-go-l12]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryprefersexplicitvault--internal-vault-query-test-go-l41]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryusesexplicittopic--internal-vault-query-test-go-l119]] via `calls` (syntactic)
- [[kodebase-go/raw/codebase/files/internal/vault/query_test.go|internal/vault/query_test.go]] via `contains` (syntactic)
