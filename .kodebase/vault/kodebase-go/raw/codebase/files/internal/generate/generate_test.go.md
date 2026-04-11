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
outgoing_relation_count: 33
smells:
  - "orphan-file"
source_kind: "codebase-file"
source_path: "internal/generate/generate_test.go"
stage: "raw"
symbol_count: 15
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/generate/generate_test.go"
type: "source"
---

# Codebase File: internal/generate/generate_test.go

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
- [[kodebase-go/raw/codebase/symbols/generate--internal-generate-generate-test-go-l1|generate (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/fakeadapter--internal-generate-generate-test-go-l15|fakeAdapter (struct)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/supports--internal-generate-generate-test-go-l23|Supports (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/parsefiles--internal-generate-generate-test-go-l27|ParseFiles (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/parsefileswithprogress--internal-generate-generate-test-go-l39|ParseFilesWithProgress (method)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testrunnergeneratecallspipelinestagesinorder--internal-generate-generate-test-go-l58|TestRunnerGenerateCallsPipelineStagesInOrder (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testselectadaptersforgoonlyworkspace--internal-generate-generate-test-go-l172|TestSelectAdaptersForGoOnlyWorkspace (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testselectadaptersformixedworkspace--internal-generate-generate-test-go-l191|TestSelectAdaptersForMixedWorkspace (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testrunnergeneratesummaryreportscounts--internal-generate-generate-test-go-l210|TestRunnerGenerateSummaryReportsCounts (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testgeneraterequiresrootpath--internal-generate-generate-test-go-l325|TestGenerateRequiresRootPath (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testgeneraterespectscanceledcontext--internal-generate-generate-test-go-l337|TestGenerateRespectsCanceledContext (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testrunnergenerateemitsparseandwriteprogressevents--internal-generate-generate-test-go-l352|TestRunnerGenerateEmitsParseAndWriteProgressEvents (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/filterevents--internal-generate-generate-test-go-l459|filterEvents (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/firstevent--internal-generate-generate-test-go-l469|firstEvent (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testclock--internal-generate-generate-test-go-l478|testClock (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/fakeadapter--internal-generate-generate-test-go-l15]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/filterevents--internal-generate-generate-test-go-l459]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/firstevent--internal-generate-generate-test-go-l469]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/generate--internal-generate-generate-test-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsefiles--internal-generate-generate-test-go-l27]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsefileswithprogress--internal-generate-generate-test-go-l39]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/supports--internal-generate-generate-test-go-l23]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testclock--internal-generate-generate-test-go-l478]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgeneraterequiresrootpath--internal-generate-generate-test-go-l325]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgeneraterespectscanceledcontext--internal-generate-generate-test-go-l337]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrunnergeneratecallspipelinestagesinorder--internal-generate-generate-test-go-l58]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrunnergenerateemitsparseandwriteprogressevents--internal-generate-generate-test-go-l352]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrunnergeneratesummaryreportscounts--internal-generate-generate-test-go-l210]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testselectadaptersforgoonlyworkspace--internal-generate-generate-test-go-l172]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testselectadaptersformixedworkspace--internal-generate-generate-test-go-l191]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsefiles--internal-generate-generate-test-go-l27]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/parsefileswithprogress--internal-generate-generate-test-go-l39]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/supports--internal-generate-generate-test-go-l23]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgeneraterequiresrootpath--internal-generate-generate-test-go-l325]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testgeneraterespectscanceledcontext--internal-generate-generate-test-go-l337]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrunnergeneratecallspipelinestagesinorder--internal-generate-generate-test-go-l58]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrunnergenerateemitsparseandwriteprogressevents--internal-generate-generate-test-go-l352]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testrunnergeneratesummaryreportscounts--internal-generate-generate-test-go-l210]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testselectadaptersforgoonlyworkspace--internal-generate-generate-test-go-l172]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testselectadaptersformixedworkspace--internal-generate-generate-test-go-l191]]
- `imports` (syntactic) -> `context`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/scanner`
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/vault`
- `imports` (syntactic) -> `reflect`
- `imports` (syntactic) -> `strings`
- `imports` (syntactic) -> `testing`
- `imports` (syntactic) -> `time`

## Backlinks
None
