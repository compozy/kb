---
afferent_coupling: 1
domain: "kodebase-go"
efferent_coupling: 1
has_circular_dependency: false
has_smells: true
incoming_relation_count: 0
instability: 0.5
is_god_file: true
is_orphan_file: false
language: "go"
outgoing_relation_count: 32
smells:
  - "god-file"
source_kind: "codebase-file"
source_path: "internal/scanner/scanner_test.go"
stage: "raw"
symbol_count: 16
tags:
  - "kodebase-go"
  - "raw"
  - "codebase"
  - "file"
  - "go"
title: "Codebase File: internal/scanner/scanner_test.go"
type: "source"
---

# Codebase File: internal/scanner/scanner_test.go

Generated raw snapshot for this source file. This note is intended to be cited by compiled wiki pages.

## Language
`go`

## Static Analysis
- Afferent coupling: 1
- Efferent coupling: 1
- Instability: 0.5
- Entry point: false
- Circular dependency: false
- Smells: `god-file`

## Module Notes
None

## Symbols
- [[kodebase-go/raw/codebase/symbols/scanner--internal-scanner-scanner-test-go-l1|scanner (package)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceroutessupportedfilesbylanguage--internal-scanner-scanner-test-go-l13|TestScanWorkspaceRoutesSupportedFilesByLanguage (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceexcludesnodemodulesbydefault--internal-scanner-scanner-test-go-l55|TestScanWorkspaceExcludesNodeModulesByDefault (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceexcludesgitdirectorybydefault--internal-scanner-scanner-test-go-l71|TestScanWorkspaceExcludesGitDirectoryByDefault (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testscanworkspacerespectsgitignorepatterns--internal-scanner-scanner-test-go-l87|TestScanWorkspaceRespectsGitIgnorePatterns (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testscanworkspacerespectsnestedgitignorepatterns--internal-scanner-scanner-test-go-l104|TestScanWorkspaceRespectsNestedGitIgnorePatterns (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceincludepatternrestrictsresults--internal-scanner-scanner-test-go-l123|TestScanWorkspaceIncludePatternRestrictsResults (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceexcludepatternremovesmatches--internal-scanner-scanner-test-go-l144|TestScanWorkspaceExcludePatternRemovesMatches (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceignoresunsupportedextensions--internal-scanner-scanner-test-go-l160|TestScanWorkspaceIgnoresUnsupportedExtensions (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testscanworkspaceemptydirectoryreturnsemptyworkspace--internal-scanner-scanner-test-go-l176|TestScanWorkspaceEmptyDirectoryReturnsEmptyWorkspace (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/testscanworkspacegroupsfilesbylanguage--internal-scanner-scanner-test-go-l196|TestScanWorkspaceGroupsFilesByLanguage (function)]] · exported=true
- [[kodebase-go/raw/codebase/symbols/scantestworkspace--internal-scanner-scanner-test-go-l218|scanTestWorkspace (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/groupedcounts--internal-scanner-scanner-test-go-l229|groupedCounts (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/groupedpaths--internal-scanner-scanner-test-go-l237|groupedPaths (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/scannedpaths--internal-scanner-scanner-test-go-l250|scannedPaths (function)]] · exported=false
- [[kodebase-go/raw/codebase/symbols/writetestfile--internal-scanner-scanner-test-go-l259|writeTestFile (function)]] · exported=false

## Outgoing Relations
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/groupedcounts--internal-scanner-scanner-test-go-l229]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/groupedpaths--internal-scanner-scanner-test-go-l237]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scannedpaths--internal-scanner-scanner-test-go-l250]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scanner--internal-scanner-scanner-test-go-l1]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/scantestworkspace--internal-scanner-scanner-test-go-l218]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testscanworkspaceemptydirectoryreturnsemptyworkspace--internal-scanner-scanner-test-go-l176]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testscanworkspaceexcludepatternremovesmatches--internal-scanner-scanner-test-go-l144]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testscanworkspaceexcludesgitdirectorybydefault--internal-scanner-scanner-test-go-l71]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testscanworkspaceexcludesnodemodulesbydefault--internal-scanner-scanner-test-go-l55]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testscanworkspacegroupsfilesbylanguage--internal-scanner-scanner-test-go-l196]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testscanworkspaceignoresunsupportedextensions--internal-scanner-scanner-test-go-l160]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testscanworkspaceincludepatternrestrictsresults--internal-scanner-scanner-test-go-l123]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testscanworkspacerespectsgitignorepatterns--internal-scanner-scanner-test-go-l87]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testscanworkspacerespectsnestedgitignorepatterns--internal-scanner-scanner-test-go-l104]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testscanworkspaceroutessupportedfilesbylanguage--internal-scanner-scanner-test-go-l13]]
- `contains` (syntactic) -> [[kodebase-go/raw/codebase/symbols/writetestfile--internal-scanner-scanner-test-go-l259]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testscanworkspaceemptydirectoryreturnsemptyworkspace--internal-scanner-scanner-test-go-l176]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testscanworkspaceexcludepatternremovesmatches--internal-scanner-scanner-test-go-l144]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testscanworkspaceexcludesgitdirectorybydefault--internal-scanner-scanner-test-go-l71]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testscanworkspaceexcludesnodemodulesbydefault--internal-scanner-scanner-test-go-l55]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testscanworkspacegroupsfilesbylanguage--internal-scanner-scanner-test-go-l196]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testscanworkspaceignoresunsupportedextensions--internal-scanner-scanner-test-go-l160]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testscanworkspaceincludepatternrestrictsresults--internal-scanner-scanner-test-go-l123]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testscanworkspacerespectsgitignorepatterns--internal-scanner-scanner-test-go-l87]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testscanworkspacerespectsnestedgitignorepatterns--internal-scanner-scanner-test-go-l104]]
- `exports` (syntactic) -> [[kodebase-go/raw/codebase/symbols/testscanworkspaceroutessupportedfilesbylanguage--internal-scanner-scanner-test-go-l13]]
- `imports` (syntactic) -> `github.com/user/go-devstack/internal/models`
- `imports` (syntactic) -> `os`
- `imports` (syntactic) -> `path/filepath`
- `imports` (syntactic) -> `reflect`
- `imports` (syntactic) -> `sort`
- `imports` (syntactic) -> `testing`

## Backlinks
None
