---
afferent_coupling: 0
domain: "kodebase-go"
efferent_coupling: 2
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 1
is_god_file: false
is_orphan_file: true
language: "go"
outgoing_relation_count: 18
smells:
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/config/config_test.go"
stage: "raw"
symbol_count: 8
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/config/config_test.go"
type: "source"
---

# Codebase File: internal/config/config_test.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 0
- Efferent coupling: 2
- Instability: 1
- Entry point: false
- Circular dependency: false
- Smells: `orphan-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/config--internal-config-config-test-go-l1|config (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testdefaultconfighasvaliddefaults--internal-config-config-test-go-l9|TestDefaultConfigHasValidDefaults (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testloadconfigroundtrip--internal-config-config-test-go-l33|TestLoadConfigRoundTrip (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testloademptypathusesdefaults--internal-config-config-test-go-l74|TestLoadEmptyPathUsesDefaults (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testloadrejectsunknownkeys--internal-config-config-test-go-l86|TestLoadRejectsUnknownKeys (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testvalidaterejectsinvalidvalues--internal-config-config-test-go-l106|TestValidateRejectsInvalidValues (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testloaddotenvifpresentloadsvalues--internal-config-config-test-go-l157|TestLoadDotEnvIfPresentLoadsValues (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testloaddotenvifpresentmissingfileisok--internal-config-config-test-go-l174|TestLoadDotEnvIfPresentMissingFileIsOK (function)]] · exported=true

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/config--internal-config-config-test-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testdefaultconfighasvaliddefaults--internal-config-config-test-go-l9]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testloadconfigroundtrip--internal-config-config-test-go-l33]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testloaddotenvifpresentloadsvalues--internal-config-config-test-go-l157]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testloaddotenvifpresentmissingfileisok--internal-config-config-test-go-l174]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testloademptypathusesdefaults--internal-config-config-test-go-l74]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testloadrejectsunknownkeys--internal-config-config-test-go-l86]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testvalidaterejectsinvalidvalues--internal-config-config-test-go-l106]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testdefaultconfighasvaliddefaults--internal-config-config-test-go-l9]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testloadconfigroundtrip--internal-config-config-test-go-l33]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testloaddotenvifpresentloadsvalues--internal-config-config-test-go-l157]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testloaddotenvifpresentmissingfileisok--internal-config-config-test-go-l174]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testloademptypathusesdefaults--internal-config-config-test-go-l74]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testloadrejectsunknownkeys--internal-config-config-test-go-l86]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testvalidaterejectsinvalidvalues--internal-config-config-test-go-l106]]
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `path/filepath`
- `imports` (syntactic) -> `testing`

## Backlinks
None
