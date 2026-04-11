---
created: "2026-04-11"
domain: "kodebase-go"
generator: "kodebase"
sources:
  - "[[kodebase-go/raw/codebase/files/.agents/skills/systematic-debugging/condition-based-waiting-example.ts]]"
  - "[[kodebase-go/raw/codebase/files/cmd/kodebase/main.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/adapter/go_adapter_integration_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/adapter/go_adapter_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/adapter/treesitter.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/adapter/treesitter_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter_integration_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/generate.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/generate_output.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/generate_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/index.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/index_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_backlinks.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_blastradius.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_circulardeps.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_complexity.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_coupling.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_deadcode.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_deps.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_file.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_helpers_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_integration_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_smells.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_symbol.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/root.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/search.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/search_index_integration_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/search_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/cli/version.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/config/config.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/config/config_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/config/env.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/generate/events.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/generate/generate.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/generate/generate_integration_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/generate/generate_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/generate/testdata/fixture-go-repo/internal/greeter/greeter.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/generate/testdata/fixture-go-repo/main.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/graph/normalize.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/graph/normalize_integration_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/graph/normalize_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/logger/logger.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/logger/logger_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/metrics/compute.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/metrics/compute_integration_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/metrics/compute_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/models/models.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/models/models_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/output/formatter.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/output/formatter_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/qmd/client.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/qmd/client_integration_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/qmd/client_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/scanner/scanner.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/scanner/scanner_integration_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/scanner/scanner_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/vault/pathutils.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/vault/pathutils_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/vault/query.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/vault/query_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/vault/reader.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/vault/reader_integration_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/vault/reader_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/vault/render.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/vault/render_base.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/vault/render_integration_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/vault/render_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/vault/render_wiki.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/vault/textutils.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/vault/textutils_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/vault/writer.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/vault/writer_integration_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/vault/writer_test.go]]"
  - "[[kodebase-go/raw/codebase/files/internal/version/version.go]]"
  - "[[kodebase-go/raw/codebase/files/magefile.go]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/.agents/skills/systematic-debugging]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/cmd/kodebase]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/internal/adapter]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/internal/cli]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/internal/config]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/internal/generate]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/internal/generate/testdata/fixture-go-repo]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/internal/generate/testdata/fixture-go-repo/internal/greeter]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/internal/graph]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/internal/logger]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/internal/metrics]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/internal/models]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/internal/output]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/internal/qmd]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/internal/scanner]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/internal/vault]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/internal/version]]"
  - "[[kodebase-go/raw/codebase/indexes/directories/root]]"
stage: "compiled"
tags:
  - "kodebase-go"
  - "wiki"
  - "codebase"
  - "starter"
title: "Module Health"
type: "wiki"
updated: "2026-04-11"
---

# Module Health

Instability highlights modules with high outgoing dependence relative to incoming dependence.

## Files

| File | Ca | Ce | Instability |
| ---- | -- | -- | ----------- |
| [[kodebase-go/raw/codebase/files/internal/cli/inspect_helpers_test.go|internal/cli/inspect_helpers_test.go]] | 0 | 6 | 1 |
| [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter.go|internal/adapter/ts_adapter.go]] | 0 | 2 | 1 |
| [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter_integration_test.go|internal/adapter/ts_adapter_integration_test.go]] | 0 | 2 | 1 |
| [[kodebase-go/raw/codebase/files/internal/config/config_test.go|internal/config/config_test.go]] | 0 | 2 | 1 |
| [[kodebase-go/raw/codebase/files/internal/scanner/scanner_integration_test.go|internal/scanner/scanner_integration_test.go]] | 0 | 2 | 1 |
| [[kodebase-go/raw/codebase/files/internal/vault/writer.go|internal/vault/writer.go]] | 0 | 2 | 1 |
| [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter_integration_test.go|internal/adapter/go_adapter_integration_test.go]] | 0 | 1 | 1 |
| [[kodebase-go/raw/codebase/files/internal/adapter/treesitter_test.go|internal/adapter/treesitter_test.go]] | 0 | 1 | 1 |
| [[kodebase-go/raw/codebase/files/internal/cli/generate_test.go|internal/cli/generate_test.go]] | 0 | 1 | 1 |
| [[kodebase-go/raw/codebase/files/internal/cli/index_test.go|internal/cli/index_test.go]] | 0 | 1 | 1 |
| [[kodebase-go/raw/codebase/files/internal/cli/search_index_integration_test.go|internal/cli/search_index_integration_test.go]] | 0 | 1 | 1 |
| [[kodebase-go/raw/codebase/files/internal/cli/search_test.go|internal/cli/search_test.go]] | 0 | 1 | 1 |
| [[kodebase-go/raw/codebase/files/internal/generate/generate_integration_test.go|internal/generate/generate_integration_test.go]] | 0 | 1 | 1 |
| [[kodebase-go/raw/codebase/files/internal/generate/generate_test.go|internal/generate/generate_test.go]] | 0 | 1 | 1 |
| [[kodebase-go/raw/codebase/files/internal/graph/normalize_integration_test.go|internal/graph/normalize_integration_test.go]] | 0 | 1 | 1 |
| [[kodebase-go/raw/codebase/files/internal/graph/normalize_test.go|internal/graph/normalize_test.go]] | 0 | 1 | 1 |
| [[kodebase-go/raw/codebase/files/internal/logger/logger_test.go|internal/logger/logger_test.go]] | 0 | 1 | 1 |
| [[kodebase-go/raw/codebase/files/internal/metrics/compute_integration_test.go|internal/metrics/compute_integration_test.go]] | 0 | 1 | 1 |
| [[kodebase-go/raw/codebase/files/internal/metrics/compute_test.go|internal/metrics/compute_test.go]] | 0 | 1 | 1 |
| [[kodebase-go/raw/codebase/files/internal/models/models_test.go|internal/models/models_test.go]] | 0 | 1 | 1 |
| [[kodebase-go/raw/codebase/files/internal/qmd/client_integration_test.go|internal/qmd/client_integration_test.go]] | 0 | 1 | 1 |
| [[kodebase-go/raw/codebase/files/internal/qmd/client_test.go|internal/qmd/client_test.go]] | 0 | 1 | 1 |
| [[kodebase-go/raw/codebase/files/internal/vault/reader.go|internal/vault/reader.go]] | 0 | 1 | 1 |
| [[kodebase-go/raw/codebase/files/internal/vault/reader_integration_test.go|internal/vault/reader_integration_test.go]] | 0 | 1 | 1 |
| [[kodebase-go/raw/codebase/files/internal/vault/render_integration_test.go|internal/vault/render_integration_test.go]] | 0 | 1 | 1 |
| [[kodebase-go/raw/codebase/files/internal/vault/writer_integration_test.go|internal/vault/writer_integration_test.go]] | 0 | 1 | 1 |
| [[kodebase-go/raw/codebase/files/internal/cli/inspect_test.go|internal/cli/inspect_test.go]] | 1 | 11 | 0.9167 |
| [[kodebase-go/raw/codebase/files/internal/vault/render_base.go|internal/vault/render_base.go]] | 1 | 2 | 0.6667 |
| [[kodebase-go/raw/codebase/files/internal/vault/render_wiki.go|internal/vault/render_wiki.go]] | 1 | 2 | 0.6667 |
| [[kodebase-go/raw/codebase/files/internal/vault/render.go|internal/vault/render.go]] | 2 | 2 | 0.5 |
| [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter.go|internal/adapter/go_adapter.go]] | 1 | 1 | 0.5 |
| [[kodebase-go/raw/codebase/files/internal/adapter/ts_adapter_test.go|internal/adapter/ts_adapter_test.go]] | 1 | 1 | 0.5 |
| [[kodebase-go/raw/codebase/files/internal/cli/index.go|internal/cli/index.go]] | 1 | 1 | 0.5 |
| [[kodebase-go/raw/codebase/files/internal/config/config.go|internal/config/config.go]] | 1 | 1 | 0.5 |
| [[kodebase-go/raw/codebase/files/internal/scanner/scanner_test.go|internal/scanner/scanner_test.go]] | 1 | 1 | 0.5 |
| [[kodebase-go/raw/codebase/files/internal/cli/inspect.go|internal/cli/inspect.go]] | 12 | 10 | 0.4545 |
| [[kodebase-go/raw/codebase/files/internal/cli/root.go|internal/cli/root.go]] | 6 | 5 | 0.4545 |
| [[kodebase-go/raw/codebase/files/internal/cli/inspect_backlinks.go|internal/cli/inspect_backlinks.go]] | 2 | 1 | 0.3333 |
| [[kodebase-go/raw/codebase/files/internal/cli/inspect_circulardeps.go|internal/cli/inspect_circulardeps.go]] | 2 | 1 | 0.3333 |
| [[kodebase-go/raw/codebase/files/internal/cli/inspect_complexity.go|internal/cli/inspect_complexity.go]] | 2 | 1 | 0.3333 |
| [[kodebase-go/raw/codebase/files/internal/cli/inspect_deadcode.go|internal/cli/inspect_deadcode.go]] | 2 | 1 | 0.3333 |
| [[kodebase-go/raw/codebase/files/internal/cli/inspect_deps.go|internal/cli/inspect_deps.go]] | 2 | 1 | 0.3333 |
| [[kodebase-go/raw/codebase/files/internal/cli/inspect_file.go|internal/cli/inspect_file.go]] | 2 | 1 | 0.3333 |
| [[kodebase-go/raw/codebase/files/internal/cli/inspect_symbol.go|internal/cli/inspect_symbol.go]] | 2 | 1 | 0.3333 |
| [[kodebase-go/raw/codebase/files/internal/vault/writer_test.go|internal/vault/writer_test.go]] | 2 | 1 | 0.3333 |
| [[kodebase-go/raw/codebase/files/internal/cli/inspect_blastradius.go|internal/cli/inspect_blastradius.go]] | 3 | 1 | 0.25 |
| [[kodebase-go/raw/codebase/files/internal/cli/inspect_coupling.go|internal/cli/inspect_coupling.go]] | 3 | 1 | 0.25 |
| [[kodebase-go/raw/codebase/files/internal/cli/inspect_smells.go|internal/cli/inspect_smells.go]] | 3 | 1 | 0.25 |
| [[kodebase-go/raw/codebase/files/.agents/skills/systematic-debugging/condition-based-waiting-example.ts|.agents/skills/systematic-debugging/condition-based-waiting-example.ts]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/cmd/kodebase/main.go|cmd/kodebase/main.go]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/adapter/go_adapter_test.go|internal/adapter/go_adapter_test.go]] | 3 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/adapter/treesitter.go|internal/adapter/treesitter.go]] | 3 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/cli/generate.go|internal/cli/generate.go]] | 1 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/cli/generate_output.go|internal/cli/generate_output.go]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/cli/inspect_integration_test.go|internal/cli/inspect_integration_test.go]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/cli/search.go|internal/cli/search.go]] | 2 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/cli/version.go|internal/cli/version.go]] | 1 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/config/env.go|internal/config/env.go]] | 2 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/generate/events.go|internal/generate/events.go]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/generate/generate.go|internal/generate/generate.go]] | 2 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/generate/testdata/fixture-go-repo/internal/greeter/greeter.go|internal/generate/testdata/fixture-go-repo/internal/greeter/greeter.go]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/generate/testdata/fixture-go-repo/main.go|internal/generate/testdata/fixture-go-repo/main.go]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/graph/normalize.go|internal/graph/normalize.go]] | 2 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/logger/logger.go|internal/logger/logger.go]] | 1 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/metrics/compute.go|internal/metrics/compute.go]] | 2 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/models/models.go|internal/models/models.go]] | 1 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/output/formatter.go|internal/output/formatter.go]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/output/formatter_test.go|internal/output/formatter_test.go]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/qmd/client.go|internal/qmd/client.go]] | 2 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/scanner/scanner.go|internal/scanner/scanner.go]] | 2 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/vault/pathutils.go|internal/vault/pathutils.go]] | 5 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/vault/pathutils_test.go|internal/vault/pathutils_test.go]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/vault/query.go|internal/vault/query.go]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/vault/query_test.go|internal/vault/query_test.go]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/vault/reader_test.go|internal/vault/reader_test.go]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/vault/render_test.go|internal/vault/render_test.go]] | 2 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/vault/textutils.go|internal/vault/textutils.go]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/vault/textutils_test.go|internal/vault/textutils_test.go]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/internal/version/version.go|internal/version/version.go]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/files/magefile.go|magefile.go]] | 0 | 0 | 0 |

## Directories

| Directory | Ca | Ce | Instability |
| --------- | -- | -- | ----------- |
| [[kodebase-go/raw/codebase/indexes/directories/root|.]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/indexes/directories/.agents/skills/systematic-debugging|.agents/skills/systematic-debugging]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/indexes/directories/cmd/kodebase|cmd/kodebase]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/indexes/directories/internal/adapter|internal/adapter]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/indexes/directories/internal/cli|internal/cli]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/indexes/directories/internal/config|internal/config]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/indexes/directories/internal/generate|internal/generate]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/indexes/directories/internal/generate/testdata/fixture-go-repo|internal/generate/testdata/fixture-go-repo]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/indexes/directories/internal/generate/testdata/fixture-go-repo/internal/greeter|internal/generate/testdata/fixture-go-repo/internal/greeter]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/indexes/directories/internal/graph|internal/graph]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/indexes/directories/internal/logger|internal/logger]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/indexes/directories/internal/metrics|internal/metrics]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/indexes/directories/internal/models|internal/models]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/indexes/directories/internal/output|internal/output]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/indexes/directories/internal/qmd|internal/qmd]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/indexes/directories/internal/scanner|internal/scanner]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/indexes/directories/internal/vault|internal/vault]] | 0 | 0 | 0 |
| [[kodebase-go/raw/codebase/indexes/directories/internal/version|internal/version]] | 0 | 0 | 0 |
