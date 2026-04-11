---
afferent_coupling: 0
domain: "kodebase-go"
efferent_coupling: 1
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 1
is_god_file: false
is_orphan_file: true
language: "go"
outgoing_relation_count: 19
smells:
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/cli/generate_test.go"
stage: "raw"
symbol_count: 6
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/cli/generate_test.go"
type: "source"
---

# Codebase File: internal/cli/generate_test.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 0
- Efferent coupling: 1
- Instability: 1
- Entry point: false
- Circular dependency: false
- Smells: `orphan-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/cli--internal-cli-generate-test-go-l1|cli (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testgeneratecommandpassesflagsandprintsjsonsummary--internal-cli-generate-test-go-l15|TestGenerateCommandPassesFlagsAndPrintsJSONSummary (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testgeneratecommandwritestexteventstostderrbydefault--internal-cli-generate-test-go-l80|TestGenerateCommandWritesTextEventsToStderrByDefault (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testgeneratecommandwritesjsoneventswhenrequested--internal-cli-generate-test-go-l111|TestGenerateCommandWritesJSONEventsWhenRequested (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testgeneratecommandrejectsinvalidprogressmode--internal-cli-generate-test-go-l167|TestGenerateCommandRejectsInvalidProgressMode (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testgeneratecommandrejectsinvalidlogformat--internal-cli-generate-test-go-l182|TestGenerateCommandRejectsInvalidLogFormat (function)]] · exported=true

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/cli--internal-cli-generate-test-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgeneratecommandpassesflagsandprintsjsonsummary--internal-cli-generate-test-go-l15]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgeneratecommandrejectsinvalidlogformat--internal-cli-generate-test-go-l182]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgeneratecommandrejectsinvalidprogressmode--internal-cli-generate-test-go-l167]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgeneratecommandwritesjsoneventswhenrequested--internal-cli-generate-test-go-l111]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgeneratecommandwritestexteventstostderrbydefault--internal-cli-generate-test-go-l80]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgeneratecommandpassesflagsandprintsjsonsummary--internal-cli-generate-test-go-l15]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgeneratecommandrejectsinvalidlogformat--internal-cli-generate-test-go-l182]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgeneratecommandrejectsinvalidprogressmode--internal-cli-generate-test-go-l167]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgeneratecommandwritesjsoneventswhenrequested--internal-cli-generate-test-go-l111]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgeneratecommandwritestexteventstostderrbydefault--internal-cli-generate-test-go-l80]]
- `imports` (syntactic) -> `bytes`
- `imports` (syntactic) -> `context`
- `imports` (syntactic) -> `encoding/json`
- `imports` (syntactic) -> `kgenerate (github.com/user/go-devstack/internal/generate)`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`
- `imports` (syntactic) -> `reflect`
- `imports` (syntactic) -> `strings`
- `imports` (syntactic) -> `testing`

## Backlinks
None
