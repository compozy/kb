---
afferent_coupling: 2
domain: "kodebase-go"
efferent_coupling: 0
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 0
is_god_file: true
is_orphan_file: false
language: "go"
outgoing_relation_count: 43
smells:
  - "god-file"
source_kind: "codebase-file"
source_path: "internal/scanner/scanner.go"
stage: "raw"
symbol_count: 25
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/scanner/scanner.go"
type: "source"
---

# Codebase File: internal/scanner/scanner.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 2
- Efferent coupling: 0
- Instability: 0
- Entry point: false
- Circular dependency: false
- Smells: `god-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/scanner--internal-scanner-scanner-go-l1|scanner (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/scanoptions--internal-scanner-scanner-go-l36|ScanOptions (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/option--internal-scanner-scanner-go-l43|Option (type)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/scanner--internal-scanner-scanner-go-l46|Scanner (struct)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/ignorerule--internal-scanner-scanner-go-l50|ignoreRule (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/newscanner--internal-scanner-scanner-go-l57|NewScanner (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/withoutputpath--internal-scanner-scanner-go-l69|WithOutputPath (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/withincludepatterns--internal-scanner-scanner-go-l76|WithIncludePatterns (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/withexcludepatterns--internal-scanner-scanner-go-l83|WithExcludePatterns (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/scanworkspace--internal-scanner-scanner-go-l90|ScanWorkspace (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/scanworkspace--internal-scanner-scanner-go-l95|ScanWorkspace (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/supportedlanguage--internal-scanner-scanner-go-l193|supportedLanguage (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/resolveoutputpath--internal-scanner-scanner-go-l212|resolveOutputPath (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/sameorinsidepath--internal-scanner-scanner-go-l225|sameOrInsidePath (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/shouldskipharddirectory--internal-scanner-scanner-go-l238|shouldSkipHardDirectory (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/collectignorerules--internal-scanner-scanner-go-l243|collectIgnoreRules (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/collectgitignorepaths--internal-scanner-scanner-go-l264|collectGitIgnorePaths (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/readgitignorerules--internal-scanner-scanner-go-l298|readGitIgnoreRules (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/builduserrules--internal-scanner-scanner-go-l308|buildUserRules (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/normalizepattern--internal-scanner-scanner-go-l335|normalizePattern (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/buildrules--internal-scanner-scanner-go-l344|buildRules (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/shouldkeeppattern--internal-scanner-scanner-go-l374|shouldKeepPattern (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/isnegatedpattern--internal-scanner-scanner-go-l383|isNegatedPattern (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/isignored--internal-scanner-scanner-go-l388|isIgnored (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/scopepath--internal-scanner-scanner-go-l409|scopePath (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/buildrules--internal-scanner-scanner-go-l344]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/builduserrules--internal-scanner-scanner-go-l308]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/collectgitignorepaths--internal-scanner-scanner-go-l264]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/collectignorerules--internal-scanner-scanner-go-l243]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/ignorerule--internal-scanner-scanner-go-l50]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isignored--internal-scanner-scanner-go-l388]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/isnegatedpattern--internal-scanner-scanner-go-l383]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newscanner--internal-scanner-scanner-go-l57]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/normalizepattern--internal-scanner-scanner-go-l335]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/option--internal-scanner-scanner-go-l43]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/readgitignorerules--internal-scanner-scanner-go-l298]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/resolveoutputpath--internal-scanner-scanner-go-l212]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/sameorinsidepath--internal-scanner-scanner-go-l225]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scanner--internal-scanner-scanner-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scanner--internal-scanner-scanner-go-l46]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scanoptions--internal-scanner-scanner-go-l36]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scanworkspace--internal-scanner-scanner-go-l90]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scanworkspace--internal-scanner-scanner-go-l95]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scopepath--internal-scanner-scanner-go-l409]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/shouldkeeppattern--internal-scanner-scanner-go-l374]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/shouldskipharddirectory--internal-scanner-scanner-go-l238]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/supportedlanguage--internal-scanner-scanner-go-l193]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withexcludepatterns--internal-scanner-scanner-go-l83]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withincludepatterns--internal-scanner-scanner-go-l76]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withoutputpath--internal-scanner-scanner-go-l69]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/newscanner--internal-scanner-scanner-go-l57]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/option--internal-scanner-scanner-go-l43]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scanner--internal-scanner-scanner-go-l46]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scanoptions--internal-scanner-scanner-go-l36]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scanworkspace--internal-scanner-scanner-go-l90]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scanworkspace--internal-scanner-scanner-go-l95]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withexcludepatterns--internal-scanner-scanner-go-l83]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withincludepatterns--internal-scanner-scanner-go-l76]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/withoutputpath--internal-scanner-scanner-go-l69]]
- `imports` (syntactic) -> `errors`
- `imports` (syntactic) -> `fmt`
- `imports` (syntactic) -> `ignore (github.com/sabhiram/go-gitignore)`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`
- `imports` (syntactic) -> `io/fs`
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `path/filepath`
- `imports` (syntactic) -> `sort`
- `imports` (syntactic) -> `strings`

## Backlinks
None
