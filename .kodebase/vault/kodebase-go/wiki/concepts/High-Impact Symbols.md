---
created: "2026-04-11"
domain: "kodebase-go"
generator: "kodebase"
sources:
  - "[[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/generate.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_backlinks.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_blastradius.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_circulardeps.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_complexity.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_coupling.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_deadcode.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_deps.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_file.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_smells.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_symbol.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/vault/pathutils.go]]"
  - "[[kodebase-go/raw/codebase/symbols/bindinspectsharedflags--internal-cli-inspect-go-l70]]"
  - "[[kodebase-go/raw/codebase/symbols/inspectfrontmatterint--internal-cli-inspect-go-l213]]"
  - "[[kodebase-go/raw/codebase/symbols/inspectfrontmatterstring--internal-cli-inspect-go-l156]]"
  - "[[kodebase-go/raw/codebase/symbols/inspectfrontmatterstringarray--internal-cli-inspect-go-l170]]"
  - "[[kodebase-go/raw/codebase/symbols/namedchildren--internal-adapter-go-adapter-go-l566]]"
  - "[[kodebase-go/raw/codebase/symbols/newgeneratecommand--internal-cli-generate-go-l15]]"
  - "[[kodebase-go/raw/codebase/symbols/newinspectbacklinkscommand--internal-cli-inspect-backlinks-go-l11]]"
  - "[[kodebase-go/raw/codebase/symbols/newinspectblastradiuscommand--internal-cli-inspect-blastradius-go-l20]]"
  - "[[kodebase-go/raw/codebase/symbols/newinspectcirculardepscommand--internal-cli-inspect-circulardeps-go-l12]]"
  - "[[kodebase-go/raw/codebase/symbols/newinspectcomplexitycommand--internal-cli-inspect-complexity-go-l21]]"
  - "[[kodebase-go/raw/codebase/symbols/newinspectcouplingcommand--internal-cli-inspect-coupling-go-l19]]"
  - "[[kodebase-go/raw/codebase/symbols/newinspectdeadcodecommand--internal-cli-inspect-deadcode-go-l19]]"
  - "[[kodebase-go/raw/codebase/symbols/newinspectdepscommand--internal-cli-inspect-deps-go-l11]]"
  - "[[kodebase-go/raw/codebase/symbols/newinspectfilecommand--internal-cli-inspect-file-go-l11]]"
  - "[[kodebase-go/raw/codebase/symbols/newinspectsmellscommand--internal-cli-inspect-smells-go-l19]]"
  - "[[kodebase-go/raw/codebase/symbols/newinspectsymbolcommand--internal-cli-inspect-symbol-go-l21]]"
  - "[[kodebase-go/raw/codebase/symbols/stripmarkdownextension--internal-vault-pathutils-go-l223]]"
  - "[[kodebase-go/raw/codebase/symbols/textof--internal-adapter-go-adapter-go-l588]]"
  - "[[kodebase-go/raw/codebase/symbols/toposixpath--internal-vault-pathutils-go-l15]]"
  - "[[kodebase-go/raw/codebase/symbols/totopicwikilink--internal-vault-pathutils-go-l228]]"
stage: "compiled"
tags:
  - "kodebase-go"
  - "wiki"
  - "codebase"
  - "starter"
title: "High-Impact Symbols"
type: "wiki"
updated: "2026-04-11"
---

# High-Impact Symbols

Blast radius counts both direct and transitive dependents reachable through call and reference chains.

## Top Symbols

| Symbol | Blast Radius | Direct Dependents | File |
| ------ | ------------ | ----------------- | ---- |
| [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstring--internal-cli-inspect-go-l156|inspectFrontmatterString]] | 42 | 16 | [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]] |
| [[kodebase-go/raw/codebase/symbols/toposixpath--internal-vault-pathutils-go-l15|ToPosixPath]] | 33 | 8 | [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] |
| [[kodebase-go/raw/codebase/symbols/inspectfrontmatterstringarray--internal-cli-inspect-go-l170|inspectFrontmatterStringArray]] | 31 | 10 | [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]] |
| [[kodebase-go/raw/codebase/symbols/stripmarkdownextension--internal-vault-pathutils-go-l223|StripMarkdownExtension]] | 29 | 1 | [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] |
| [[kodebase-go/raw/codebase/symbols/textof--internal-adapter-go-adapter-go-l588|textOf]] | 28 | 16 | [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] |
| [[kodebase-go/raw/codebase/symbols/totopicwikilink--internal-vault-pathutils-go-l228|ToTopicWikiLink]] | 28 | 10 | [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] |
| [[kodebase-go/raw/codebase/symbols/inspectfrontmatterint--internal-cli-inspect-go-l213|inspectFrontmatterInt]] | 27 | 9 | [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]] |
| [[kodebase-go/raw/codebase/symbols/namedchildren--internal-adapter-go-adapter-go-l566|namedChildren]] | 25 | 10 | [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] |
| [[kodebase-go/raw/codebase/symbols/bindinspectsharedflags--internal-cli-inspect-go-l70|bindInspectSharedFlags]] | 25 | 1 | [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]] |
| [[kodebase-go/raw/codebase/symbols/newinspectbacklinkscommand--internal-cli-inspect-backlinks-go-l11|newInspectBacklinksCommand]] | 25 | 1 | [[kodebase-go/raw/codebase/files/internal/cli/inspect_backlinks.go|internal/cli/inspect_backlinks.go]] |
| [[kodebase-go/raw/codebase/symbols/newinspectblastradiuscommand--internal-cli-inspect-blastradius-go-l20|newInspectBlastRadiusCommand]] | 25 | 1 | [[kodebase-go/raw/codebase/files/internal/cli/inspect_blastradius.go|internal/cli/inspect_blastradius.go]] |
| [[kodebase-go/raw/codebase/symbols/newinspectcirculardepscommand--internal-cli-inspect-circulardeps-go-l12|newInspectCircularDepsCommand]] | 25 | 1 | [[kodebase-go/raw/codebase/files/internal/cli/inspect_circulardeps.go|internal/cli/inspect_circulardeps.go]] |
| [[kodebase-go/raw/codebase/symbols/newinspectcomplexitycommand--internal-cli-inspect-complexity-go-l21|newInspectComplexityCommand]] | 25 | 1 | [[kodebase-go/raw/codebase/files/internal/cli/inspect_complexity.go|internal/cli/inspect_complexity.go]] |
| [[kodebase-go/raw/codebase/symbols/newinspectcouplingcommand--internal-cli-inspect-coupling-go-l19|newInspectCouplingCommand]] | 25 | 1 | [[kodebase-go/raw/codebase/files/internal/cli/inspect_coupling.go|internal/cli/inspect_coupling.go]] |
| [[kodebase-go/raw/codebase/symbols/newinspectdeadcodecommand--internal-cli-inspect-deadcode-go-l19|newInspectDeadCodeCommand]] | 25 | 1 | [[kodebase-go/raw/codebase/files/internal/cli/inspect_deadcode.go|internal/cli/inspect_deadcode.go]] |
| [[kodebase-go/raw/codebase/symbols/newinspectdepscommand--internal-cli-inspect-deps-go-l11|newInspectDepsCommand]] | 25 | 1 | [[kodebase-go/raw/codebase/files/internal/cli/inspect_deps.go|internal/cli/inspect_deps.go]] |
| [[kodebase-go/raw/codebase/symbols/newinspectfilecommand--internal-cli-inspect-file-go-l11|newInspectFileCommand]] | 25 | 1 | [[kodebase-go/raw/codebase/files/internal/cli/inspect_file.go|internal/cli/inspect_file.go]] |
| [[kodebase-go/raw/codebase/symbols/newinspectsmellscommand--internal-cli-inspect-smells-go-l19|newInspectSmellsCommand]] | 25 | 1 | [[kodebase-go/raw/codebase/files/internal/cli/inspect_smells.go|internal/cli/inspect_smells.go]] |
| [[kodebase-go/raw/codebase/symbols/newinspectsymbolcommand--internal-cli-inspect-symbol-go-l21|newInspectSymbolCommand]] | 25 | 1 | [[kodebase-go/raw/codebase/files/internal/cli/inspect_symbol.go|internal/cli/inspect_symbol.go]] |
| [[kodebase-go/raw/codebase/symbols/newgeneratecommand--internal-cli-generate-go-l15|newGenerateCommand]] | 24 | 1 | [[kodebase-go/raw/codebase/files/internal/cli/generate.go|internal/cli/generate.go]] |
