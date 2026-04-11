---
afferent_coupling: 0
domain: "kodebase-go"
efferent_coupling: 0
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 0
is_god_file: false
is_orphan_file: true
language: "go"
outgoing_relation_count: 24
smells:
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/vault/query_test.go"
stage: "raw"
symbol_count: 11
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/vault/query_test.go"
type: "source"
---

# Codebase File: internal/vault/query_test.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 0
- Efferent coupling: 0
- Instability: 0
- Entry point: false
- Circular dependency: false
- Smells: `orphan-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/vault-test--internal-vault-query-test-go-l1|vault_test (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryfindsvaultbywalkingup--internal-vault-query-test-go-l12|TestResolveVaultQueryFindsVaultByWalkingUp (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryprefersexplicitvault--internal-vault-query-test-go-l41|TestResolveVaultQueryPrefersExplicitVault (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryautoresolvessingletopic--internal-vault-query-test-go-l72|TestResolveVaultQueryAutoResolvesSingleTopic (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryerrorswhenmultipletopicsexist--internal-vault-query-test-go-l95|TestResolveVaultQueryErrorsWhenMultipleTopicsExist (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryusesexplicittopic--internal-vault-query-test-go-l119|TestResolveVaultQueryUsesExplicitTopic (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryerrorswhenexplicittopicismissing--internal-vault-query-test-go-l144|TestResolveVaultQueryErrorsWhenExplicitTopicIsMissing (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryerrorsclearlywhennovaultisfound--internal-vault-query-test-go-l164|TestResolveVaultQueryErrorsClearlyWhenNoVaultIsFound (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testlistavailabletopicsreturnssortedtopics--internal-vault-query-test-go-l178|TestListAvailableTopicsReturnsSortedTopics (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/mkdirall--internal-vault-query-test-go-l206|mkdirAll (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/writetopicmarker--internal-vault-query-test-go-l214|writeTopicMarker (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/mkdirall--internal-vault-query-test-go-l206]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testlistavailabletopicsreturnssortedtopics--internal-vault-query-test-go-l178]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryautoresolvessingletopic--internal-vault-query-test-go-l72]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryerrorsclearlywhennovaultisfound--internal-vault-query-test-go-l164]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryerrorswhenexplicittopicismissing--internal-vault-query-test-go-l144]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryerrorswhenmultipletopicsexist--internal-vault-query-test-go-l95]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryfindsvaultbywalkingup--internal-vault-query-test-go-l12]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryprefersexplicitvault--internal-vault-query-test-go-l41]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryusesexplicittopic--internal-vault-query-test-go-l119]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/vault-test--internal-vault-query-test-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writetopicmarker--internal-vault-query-test-go-l214]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testlistavailabletopicsreturnssortedtopics--internal-vault-query-test-go-l178]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryautoresolvessingletopic--internal-vault-query-test-go-l72]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryerrorsclearlywhennovaultisfound--internal-vault-query-test-go-l164]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryerrorswhenexplicittopicismissing--internal-vault-query-test-go-l144]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryerrorswhenmultipletopicsexist--internal-vault-query-test-go-l95]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryfindsvaultbywalkingup--internal-vault-query-test-go-l12]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryprefersexplicitvault--internal-vault-query-test-go-l41]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testresolvevaultqueryusesexplicittopic--internal-vault-query-test-go-l119]]
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/vault`
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `path/filepath`
- `imports` (syntactic) -> `strings`
- `imports` (syntactic) -> `testing`

## Backlinks
None
