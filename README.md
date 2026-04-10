<div align="center">

# Kodebase

### Turn any codebase into a knowledge base.

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![CI](https://img.shields.io/github/actions/workflow/status/pedronauck/kodebase-go/ci.yaml?branch=main&label=CI)](https://github.com/pedronauck/kodebase-go/actions)
[![Go](https://img.shields.io/badge/Go-1.24-00ADD8.svg)](https://go.dev/)

[Install](#install) &#8226; [See It Work](#see-it-work) &#8226; [Features](#features) &#8226; [Commands](#commands) &#8226; [Contributing](#contributing)

</div>

---

Kodebase parses your TypeScript, JavaScript, or Go repository and generates a structured [Obsidian](https://obsidian.md) vault -- complete with every symbol, file, and dependency relationship mapped into interconnected markdown notes. It computes cyclomatic complexity, blast radius, coupling, instability, and dead code detection, then compiles 10 starter wiki articles and 11 interactive [Base](https://obsidian.md/blog/bases/) views so you can explore your architecture the same way you explore your notes.

No SaaS. No cloud. Just markdown.

---

## Install

```bash
# build from source
git clone https://github.com/pedronauck/kodebase-go.git
cd kodebase-go
make build
# binary is at bin/kodebase
```

**Optional** -- for semantic search capabilities:

```bash
npm install -g @tobilu/qmd
```

> [!NOTE]
> **Requirements:** Go >= 1.24. The `search` and `index` commands require [QMD](https://github.com/tobilu/qmd) to be installed separately. All other commands work standalone.

<details>
<summary><strong>What it touches</strong></summary>

- **Creates files** in `.kodebase/vault/` inside the target repository (or a custom `--output` path)
- **Reads** source files in the target repository (never modifies them)
- **No network calls** -- all analysis is local. The optional QMD integration runs a local index at `~/.qmd/`
- **No telemetry** -- nothing is sent anywhere
- **Uninstall:** Remove the `kodebase` binary from your `PATH` and delete the `.kodebase/` directory

</details>

---

## See It Work

Generate a knowledge vault from any repository:

```bash
$ kodebase generate ./my-project
```

```json
{
  "filesScanned": 147,
  "filesParsed": 142,
  "symbolsExtracted": 891,
  "relationsEmitted": 2347,
  "rawDocumentsWritten": 1033,
  "wikiDocumentsWritten": 10,
  "indexDocumentsWritten": 3
}
```

Query metrics from the terminal -- no external services required:

```bash
$ kodebase inspect complexity --top 5

 symbol_name       | cyclomatic_complexity | loc | source_path
 computeMetrics    | 12                    | 89  | src/compute-metrics.ts
 parseTypeScript   | 9                     | 67  | src/adapters/typescript.ts
 normalizeGraph    | 8                     | 45  | src/normalize-graph.ts
 renderDocuments   | 7                     | 112 | src/render-documents.ts
 scanWorkspace     | 6                     | 54  | src/scan-workspace.ts
```

```bash
$ kodebase inspect dead-code

 kind   | name           | source_path                | reason
 symbol | oldHelper      | src/utils.ts               | dead-export
 file   | unused-cfg.ts  | src/config/unused-cfg.ts   | orphan-file
```

Search your vault with natural language (requires QMD):

```bash
$ kodebase search "error handling patterns" --limit 3

 title                  | score | path
 Error Handling Guide   | 0.89  | wiki/concepts/Error Handling.md
 src/error-boundary.ts  | 0.74  | raw/codebase/files/src/error-boundary.ts.md
 handleError            | 0.61  | raw/codebase/symbols/handleError--src-utils-l42.md
```

---

## Features

**Knowledge base generation** -- Point Kodebase at a repository and it generates a [Karpathy-style](https://github.com/karpathy) Obsidian vault with raw source snapshots in `raw/codebase/`, 10 synthesized wiki articles in `wiki/`, and 11 Base views in `bases/`. Every file, symbol, and dependency becomes a linked markdown note you can browse, search, and extend.

**Metrics without a dashboard** -- Kodebase computes cyclomatic complexity, blast radius, coupling, instability, and betweenness centrality at the symbol and file level. Seven code smell detectors flag specific problems, not just a letter grade.

**10 inspect subcommands** -- Query your codebase like a database, from the terminal. Rank functions by complexity, find dead exports nobody imports, trace dependency chains, detect circular imports. Three output formats (table, JSON, TSV), zero external dependencies.

**Obsidian-native output** -- Every generated note uses wikilinks, YAML frontmatter, and backlink-aware cross-references. The 11 Base views give you filterable, sortable tables inside Obsidian -- Symbol Explorer, Complexity Hotspots, Danger Zone, Module Health, and more.

**Semantic search** -- Index your vault with QMD for hybrid, lexical, or vector search across all documentation. Useful for onboarding, architecture review, or feeding context to LLMs.

**Safe reruns** -- The `raw/codebase/` layer is machine-managed and refreshes on every run. Your manual notes, `outputs/` directory, and custom wiki pages are preserved across reruns.

**AI-friendly output** -- Every vault includes `CLAUDE.md` and `AGENTS.md` as schema documents. The structured markdown output is designed for direct consumption by LLMs -- frontmatter metadata, consistent structure, and explicit cross-references.

**Single binary** -- No runtime dependencies. One `kodebase` binary handles everything. Built in Go for fast startup and low memory usage.

---

## Why Kodebase

Most code analysis tools give you a dashboard. Kodebase gives you a knowledge base.

**SonarQube and CodeClimate** tell you your code has problems. Kodebase tells you which problems, where they connect, and gives you a structured workspace to reason about them. The output is markdown you own, not a SaaS dashboard that disappears when you cancel the subscription.

**Sourcegraph** is excellent for code search across repositories. Kodebase is for understanding a single repository deeply -- its architecture, its coupling patterns, its risk surface -- and building a persistent knowledge artifact around that understanding.

**Static analyzers** produce reports. Kodebase produces a living document system where analysis results are woven into wiki articles, linked to raw source snapshots, and searchable with semantic queries. The output gets better as you add your own notes.

The key difference: Kodebase outputs compound. A SonarQube scan from last month is stale data. A Kodebase vault from last month is a knowledge base you've been building on for 30 days.

---

## Commands

### `generate`

Parse a repository and create a knowledge vault.

```bash
kodebase generate <root> [options]
```

| Argument       | Type       | Required | Description                                              |
| -------------- | ---------- | -------- | -------------------------------------------------------- |
| `root`         | positional | yes      | Path to the repository root to scan                      |
| `--output`     | string     | no       | Vault root where the topic folder will be created        |
| `--topic`      | string     | no       | Override the generated topic slug                        |
| `--title`      | string     | no       | Override the generated topic title                       |
| `--domain`     | string     | no       | Override the topic domain used in frontmatter            |
| `--include`    | string[]   | no       | Re-include path patterns that would otherwise be ignored |
| `--exclude`    | string[]   | no       | Exclude additional path patterns from scanning           |
| `--semantic`   | boolean    | no       | Enable semantic relation extraction                      |
| `--progress`   | boolean    | no       | Show progress bar during generation                      |
| `--log-format` | string     | no       | Log output format (`text` or `json`)                     |

Output: JSON summary on stdout with vault path, scan counts, generated document counts, and diagnostics.

### `inspect`

Query vault data using frontmatter and extracted metrics. No external dependencies required.

```bash
kodebase inspect <subcommand> [options]
```

| Category | Subcommand      | Description                                          |
| -------- | --------------- | ---------------------------------------------------- |
| Metrics  | `smells`        | List symbols and files with detected code smells     |
| Metrics  | `dead-code`     | List dead exports and orphan files                   |
| Metrics  | `complexity`    | Rank functions by cyclomatic complexity              |
| Metrics  | `blast-radius`  | Rank symbols by blast radius (transitive dependents) |
| Metrics  | `coupling`      | Rank files by instability (efferent/afferent)        |
| Graph    | `backlinks`     | Show what references a given symbol                  |
| Graph    | `deps`          | Show outgoing relations for a file                   |
| Graph    | `circular-deps` | List files participating in circular dependencies    |
| Lookup   | `symbol`        | Fuzzy-match and detail view of a symbol              |
| Lookup   | `file`          | Exact lookup of a file by source path                |

Shared flags: `--format table|json|tsv` (default: `table`), `--vault <path>`, `--topic <slug>`.

### `search`

Semantic search across vault documents. Requires [QMD](https://github.com/tobilu/qmd).

```bash
kodebase search <query> [options]
```

| Argument       | Type       | Default | Description                              |
| -------------- | ---------- | ------- | ---------------------------------------- |
| `query`        | positional | --      | Search query string                      |
| `--lex`        | boolean    | false   | Use BM25 keyword search only             |
| `--vec`        | boolean    | false   | Use vector similarity search only        |
| `--limit`      | number     | 10      | Maximum results to return                |
| `--full`       | boolean    | false   | Show full document instead of snippet    |
| `--min-score`  | number     | --      | Minimum similarity threshold             |
| `--all`        | boolean    | false   | Return all matches above threshold       |
| `--collection` | string     | --      | Explicit QMD collection name             |
| `--format`     | enum       | table   | Output format: `table`, `json`, or `tsv` |

Default mode is hybrid (BM25 + vector similarity). Use `--lex` or `--vec` to restrict.

### `index`

Create or update a QMD collection for semantic indexing. Requires [QMD](https://github.com/tobilu/qmd). Also available as `kodebase index-vault`.

```bash
kodebase index [options]
```

| Argument        | Type    | Default | Description                                |
| --------------- | ------- | ------- | ------------------------------------------ |
| `--vault`       | string  | --      | Vault root path                            |
| `--topic`       | string  | --      | Topic slug inside vault                    |
| `--name`        | string  | --      | Override derived QMD collection name       |
| `--embed`       | boolean | true    | Run embedding after syncing files          |
| `--context`     | string  | --      | Attach context to improve search relevance |
| `--force-embed` | boolean | false   | Force re-embedding of all documents        |

---

## What Gets Generated

```text
.kodebase/vault/
  my-repo/
    CLAUDE.md                    # Schema document for LLMs
    AGENTS.md -> CLAUDE.md       # Symlink for Codex parity
    log.md                       # Append-only operation log
    raw/
      codebase/
        files/                   # One markdown note per source file
        symbols/                 # One markdown note per extracted symbol
        indexes/
          directories/           # Directory-level inventories
          languages/             # Language-level inventories
    wiki/
      concepts/                  # 10 generated wiki articles
        Codebase Overview.md
        Directory Map.md
        Symbol Taxonomy.md
        Dependency Hotspots.md
        Complexity Hotspots.md
        Module Health.md
        Dead Code Report.md
        Code Smells.md
        Circular Dependencies.md
        High-Impact Symbols.md
      index/
        Dashboard.md             # Landing page
        Concept Index.md         # Article listing
        Source Index.md          # Reverse index (raw -> articles)
    outputs/                     # Your analysis outputs (preserved)
    bases/                       # 11 Obsidian Base views
```

- **`raw/`** -- Machine-generated source snapshots. Refreshed on every run.
- **`wiki/`** -- Starter articles synthesized from metrics. Managed areas refresh; your additions are preserved.
- **`outputs/`** -- A place for your own briefings, queries, diagrams, and reports. Never touched by Kodebase.
- **`bases/`** -- Obsidian Base `.base` files for interactive table/card/list views of metrics.

---

## Supported Languages

| Language   | Extensions    | Parser        | Relation Confidence |
| ---------- | ------------- | ------------- | ------------------- |
| TypeScript | `.ts`, `.tsx` | `tree-sitter` | syntactic           |
| JavaScript | `.js`, `.jsx` | `tree-sitter` | syntactic           |
| Go         | `.go`         | `tree-sitter` | syntactic           |

Want to add a language? See [CONTRIBUTING.md](CONTRIBUTING.md#adding-a-new-language-adapter).

---

## Tracked Relations

| Relation     | Description                                  |
| ------------ | -------------------------------------------- |
| `imports`    | File imports another file or external module |
| `exports`    | File exports a symbol                        |
| `calls`      | Function or method calls another symbol      |
| `references` | Code references a symbol                     |
| `declares`   | File declares a symbol                       |
| `contains`   | File structurally contains a symbol          |

---

## Detected Code Smells

| Smell               | Scope  | Condition                                              |
| ------------------- | ------ | ------------------------------------------------------ |
| `dead-export`       | symbol | Exported but never referenced from outside its file    |
| `long-function`     | symbol | Function with > 50 LOC or cyclomatic complexity > 10   |
| `high-blast-radius` | symbol | More than 20 transitive dependents                     |
| `bottleneck`        | symbol | Betweenness centrality > 0.1                           |
| `feature-envy`      | symbol | More references to another file's symbols than its own |
| `god-file`          | file   | More than 15 symbols or efferent coupling > 10         |
| `orphan-file`       | file   | Zero afferent coupling and not an entry point          |

---

## Excluded Paths

The scanner:

- Always skips `.git`, `.hg`, `.svn`, symlinks, and the configured vault root itself
- Applies default convenience ignores for `vendor/`, `.turbo/`, `.next/`, `node_modules/`, `dist/`, `build/`, and `coverage/`
- Respects `.gitignore` files found from the scan root downward, including nested ones
- Applies `--exclude` patterns after repository ignore rules
- Applies `--include` patterns last as explicit re-includes

`--include` and `--exclude` use `.gitignore`-style patterns evaluated against paths relative to the scan root.

---

## Development

**Prerequisites:** [Go](https://go.dev) >= 1.24

```bash
git clone https://github.com/pedronauck/kodebase-go.git
cd kodebase-go
make verify    # format + lint + test + build + boundaries
```

| Command                | Description                              |
| ---------------------- | ---------------------------------------- |
| `make fmt`             | Format all Go files with gofmt           |
| `make lint`            | Run golangci-lint with zero tolerance    |
| `make test`            | Unit tests with race detector            |
| `make test-integration`| Unit + integration tests                 |
| `make build`           | Build binary to `bin/kodebase`           |
| `make verify`          | fmt -> lint -> test -> build -> boundaries|
| `make deps`            | Run `go mod tidy`                        |

See [CONTRIBUTING.md](CONTRIBUTING.md) for code style, testing requirements, and how to add a new language adapter.

---

## Contributing

Kodebase is MIT-licensed and built in the open. We welcome contributions of all kinds:

- **Language adapters** -- Add support for Python, Rust, Java, or any language with a tree-sitter grammar
- **New code smell detectors** -- The metrics engine is designed to be extended
- **Wiki article templates** -- Better starter articles mean better vaults out of the box
- **Bug reports and feature requests** -- [Open an issue](https://github.com/pedronauck/kodebase-go/issues), we read them all

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

---

## License

MIT -- see [LICENSE](LICENSE).
