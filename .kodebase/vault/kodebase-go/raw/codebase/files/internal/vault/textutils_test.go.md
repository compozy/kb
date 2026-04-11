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
outgoing_relation_count: 13
smells:
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/vault/textutils_test.go"
stage: "raw"
symbol_count: 6
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/vault/textutils_test.go"
type: "source"
---

# Codebase File: internal/vault/textutils_test.go

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
- [[kodebase-go/raw/codebase/symbols/vault-test--internal-vault-textutils-test-go-l1|vault_test (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testnormalizecomment--internal-vault-textutils-test-go-l9|TestNormalizeComment (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testextractleadingcommentfromgosource--internal-vault-textutils-test-go-l24|TestExtractLeadingCommentFromGoSource (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testextractleadingcommentfromtssource--internal-vault-textutils-test-go-l40|TestExtractLeadingCommentFromTSSource (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testextractleadingcommentignoresnonleadingcomments--internal-vault-textutils-test-go-l56|TestExtractLeadingCommentIgnoresNonLeadingComments (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/teststripquotes--internal-vault-textutils-test-go-l70|TestStripQuotes (function)]] · exported=true

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testextractleadingcommentfromgosource--internal-vault-textutils-test-go-l24]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testextractleadingcommentfromtssource--internal-vault-textutils-test-go-l40]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testextractleadingcommentignoresnonleadingcomments--internal-vault-textutils-test-go-l56]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnormalizecomment--internal-vault-textutils-test-go-l9]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/teststripquotes--internal-vault-textutils-test-go-l70]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/vault-test--internal-vault-textutils-test-go-l1]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testextractleadingcommentfromgosource--internal-vault-textutils-test-go-l24]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testextractleadingcommentfromtssource--internal-vault-textutils-test-go-l40]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testextractleadingcommentignoresnonleadingcomments--internal-vault-textutils-test-go-l56]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testnormalizecomment--internal-vault-textutils-test-go-l9]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/teststripquotes--internal-vault-textutils-test-go-l70]]
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/vault`
- `imports` (syntactic) -> `testing`

## Backlinks
None
