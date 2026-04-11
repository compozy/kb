---
afferent_coupling: 0
domain: "kodebase-go"
efferent_coupling: 6
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 1
is_god_file: false
is_orphan_file: true
language: "go"
outgoing_relation_count: 23
smells:
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/cli/inspect_helpers_test.go"
stage: "raw"
symbol_count: 9
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/cli/inspect_helpers_test.go"
type: "source"
---

# Codebase File: internal/cli/inspect_helpers_test.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 0
- Efferent coupling: 6
- Instability: 1
- Entry point: false
- Circular dependency: false
- Smells: `orphan-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/cli--internal-cli-inspect-helpers-test-go-l1|cli (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testparseinspectoutputformat--internal-cli-inspect-helpers-test-go-l13|TestParseInspectOutputFormat (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testinspectfrontmatterhelpers--internal-cli-inspect-helpers-test-go-l51|TestInspectFrontmatterHelpers (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtosmelloutputfiltersbytype--internal-cli-inspect-helpers-test-go-l152|TestToSmellOutputFiltersByType (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtoblastradiusoutputrespectsminimumandtop--internal-cli-inspect-helpers-test-go-l178|TestToBlastRadiusOutputRespectsMinimumAndTop (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testtocouplingoutputfiltersunstableonly--internal-cli-inspect-helpers-test-go-l198|TestToCouplingOutputFiltersUnstableOnly (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testinspectcommandvalidationerrors--internal-cli-inspect-helpers-test-go-l214|TestInspectCommandValidationErrors (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testinspectparenthelplistsallsubcommands--internal-cli-inspect-helpers-test-go-l260|TestInspectParentHelpListsAllSubcommands (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testversioncommandprintsbuildversion--internal-cli-inspect-helpers-test-go-l281|TestVersionCommandPrintsBuildVersion (function)]] · exported=true

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cli--internal-cli-inspect-helpers-test-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testinspectcommandvalidationerrors--internal-cli-inspect-helpers-test-go-l214]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testinspectfrontmatterhelpers--internal-cli-inspect-helpers-test-go-l51]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testinspectparenthelplistsallsubcommands--internal-cli-inspect-helpers-test-go-l260]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testparseinspectoutputformat--internal-cli-inspect-helpers-test-go-l13]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtoblastradiusoutputrespectsminimumandtop--internal-cli-inspect-helpers-test-go-l178]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtocouplingoutputfiltersunstableonly--internal-cli-inspect-helpers-test-go-l198]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtosmelloutputfiltersbytype--internal-cli-inspect-helpers-test-go-l152]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testversioncommandprintsbuildversion--internal-cli-inspect-helpers-test-go-l281]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testinspectcommandvalidationerrors--internal-cli-inspect-helpers-test-go-l214]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testinspectfrontmatterhelpers--internal-cli-inspect-helpers-test-go-l51]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testinspectparenthelplistsallsubcommands--internal-cli-inspect-helpers-test-go-l260]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testparseinspectoutputformat--internal-cli-inspect-helpers-test-go-l13]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtoblastradiusoutputrespectsminimumandtop--internal-cli-inspect-helpers-test-go-l178]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtocouplingoutputfiltersunstableonly--internal-cli-inspect-helpers-test-go-l198]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testtosmelloutputfiltersbytype--internal-cli-inspect-helpers-test-go-l152]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testversioncommandprintsbuildversion--internal-cli-inspect-helpers-test-go-l281]]
- `imports` (syntactic) -> `bytes`
- `imports` (syntactic) -> `context`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/output`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/vault`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/version`
- `imports` (syntactic) -> `testing`

## Backlinks
None
