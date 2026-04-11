---
created: "2026-04-11"
domain: "kodebase-go"
generator: "kodebase"
sources:
  - "[[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_helpers_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/generate/generate.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/metrics/compute.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/qmd/client.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/vault/reader.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/vault/render_base.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/vault/render_wiki.go]]"
  - "[[kodebase-go/raw/codebase/symbols/computeapproxcentrality--internal-metrics-compute-go-l261]]"
  - "[[kodebase-go/raw/codebase/symbols/computemetrics--internal-metrics-compute-go-l23]]"
  - "[[kodebase-go/raw/codebase/symbols/createcodesmellsarticle--internal-vault-render-wiki-go-l784]]"
  - "[[kodebase-go/raw/codebase/symbols/extractcommonjsexports--internal-adapter-ts-adapter-go-l583]]"
  - "[[kodebase-go/raw/codebase/symbols/extractrequirebindings--internal-adapter-ts-adapter-go-l644]]"
  - "[[kodebase-go/raw/codebase/symbols/extracttsexports--internal-adapter-ts-adapter-go-l484]]"
  - "[[kodebase-go/raw/codebase/symbols/extracttsimports--internal-adapter-ts-adapter-go-l391]]"
  - "[[kodebase-go/raw/codebase/symbols/frontmatterint--internal-vault-reader-go-l376]]"
  - "[[kodebase-go/raw/codebase/symbols/generatewithobserver--internal-generate-generate-go-l88]]"
  - "[[kodebase-go/raw/codebase/symbols/index--internal-qmd-client-go-l232]]"
  - "[[kodebase-go/raw/codebase/symbols/inspectfrontmatterfloat--internal-cli-inspect-go-l254]]"
  - "[[kodebase-go/raw/codebase/symbols/inspectfrontmatterint--internal-cli-inspect-go-l213]]"
  - "[[kodebase-go/raw/codebase/symbols/matchcommonjsexporttarget--internal-adapter-ts-adapter-go-l1290]]"
  - "[[kodebase-go/raw/codebase/symbols/parsefileswithprogress--internal-adapter-go-adapter-go-l62]]"
  - "[[kodebase-go/raw/codebase/symbols/parsefileswithprogress--internal-adapter-ts-adapter-go-l93]]"
  - "[[kodebase-go/raw/codebase/symbols/parsegofile--internal-adapter-go-adapter-go-l161]]"
  - "[[kodebase-go/raw/codebase/symbols/parseindexstatus--internal-qmd-client-go-l598]]"
  - "[[kodebase-go/raw/codebase/symbols/parsetsfile--internal-adapter-ts-adapter-go-l299]]"
  - "[[kodebase-go/raw/codebase/symbols/renderyamlvalue--internal-vault-render-base-go-l295]]"
  - "[[kodebase-go/raw/codebase/symbols/testinspectfrontmatterhelpers--internal-cli-inspect-helpers-test-go-l51]]"
stage: "compiled"
tags:
  - "kodebase-go"
  - "wiki"
  - "codebase"
  - "starter"
title: "Complexity Hotspots"
type: "wiki"
updated: "2026-04-11"
---

# Complexity Hotspots

These functions have the highest measured cyclomatic complexity in the current codebase snapshot.

## Top Functions

| Symbol | Complexity | LOC | File |
| ------ | ---------- | --- | ---- |
| [[kodebase-go/raw/codebase/symbols/computemetrics--internal-metrics-compute-go-l23|ComputeMetrics]] | 52 | 237 | [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]] |
| [[kodebase-go/raw/codebase/symbols/parsefileswithprogress--internal-adapter-ts-adapter-go-l93|ParseFilesWithProgress]] | 37 | 205 | [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] |
| [[kodebase-go/raw/codebase/symbols/parseindexstatus--internal-qmd-client-go-l598|parseIndexStatus]] | 30 | 89 | [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] |
| [[kodebase-go/raw/codebase/symbols/extractrequirebindings--internal-adapter-ts-adapter-go-l644|extractRequireBindings]] | 25 | 121 | [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] |
| [[kodebase-go/raw/codebase/symbols/computeapproxcentrality--internal-metrics-compute-go-l261|computeApproxCentrality]] | 24 | 106 | [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]] |
| [[kodebase-go/raw/codebase/symbols/renderyamlvalue--internal-vault-render-base-go-l295|renderYAMLValue]] | 23 | 68 | [[kodebase-go/raw/codebase/files/internal/vault/render_base.go|internal/vault/render_base.go]] |
| [[kodebase-go/raw/codebase/symbols/extractcommonjsexports--internal-adapter-ts-adapter-go-l583|extractCommonJSExports]] | 22 | 60 | [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] |
| [[kodebase-go/raw/codebase/symbols/parsefileswithprogress--internal-adapter-go-adapter-go-l62|ParseFilesWithProgress]] | 20 | 98 | [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] |
| [[kodebase-go/raw/codebase/symbols/matchcommonjsexporttarget--internal-adapter-ts-adapter-go-l1290|matchCommonJSExportTarget]] | 20 | 42 | [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] |
| [[kodebase-go/raw/codebase/symbols/generatewithobserver--internal-generate-generate-go-l88|GenerateWithObserver]] | 19 | 217 | [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]] |
| [[kodebase-go/raw/codebase/symbols/testinspectfrontmatterhelpers--internal-cli-inspect-helpers-test-go-l51|TestInspectFrontmatterHelpers]] | 19 | 100 | [[kodebase-go/raw/codebase/files/internal/cli/inspect_helpers_test.go|internal/cli/inspect_helpers_test.go]] |
| [[kodebase-go/raw/codebase/symbols/extracttsexports--internal-adapter-ts-adapter-go-l484|extractTSExports]] | 18 | 98 | [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] |
| [[kodebase-go/raw/codebase/symbols/createcodesmellsarticle--internal-vault-render-wiki-go-l784|createCodeSmellsArticle]] | 18 | 80 | [[kodebase-go/raw/codebase/files/internal/vault/render_wiki.go|internal/vault/render_wiki.go]] |
| [[kodebase-go/raw/codebase/symbols/parsegofile--internal-adapter-go-adapter-go-l161|parseGoFile]] | 17 | 102 | [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] |
| [[kodebase-go/raw/codebase/symbols/extracttsimports--internal-adapter-ts-adapter-go-l391|extractTSImports]] | 17 | 92 | [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] |
| [[kodebase-go/raw/codebase/symbols/frontmatterint--internal-vault-reader-go-l376|frontmatterInt]] | 17 | 41 | [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]] |
| [[kodebase-go/raw/codebase/symbols/inspectfrontmatterfloat--internal-cli-inspect-go-l254|inspectFrontmatterFloat]] | 17 | 40 | [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]] |
| [[kodebase-go/raw/codebase/symbols/inspectfrontmatterint--internal-cli-inspect-go-l213|inspectFrontmatterInt]] | 17 | 40 | [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]] |
| [[kodebase-go/raw/codebase/symbols/parsetsfile--internal-adapter-ts-adapter-go-l299|parseTSFile]] | 16 | 78 | [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] |
| [[kodebase-go/raw/codebase/symbols/index--internal-qmd-client-go-l232|Index]] | 15 | 74 | [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] |

## Cross-References

Compare these hotspots against [[kodebase-go/wiki/concepts/High-Impact Symbols|High-Impact Symbols]] to distinguish locally complex functions from high-blast-radius functions.
