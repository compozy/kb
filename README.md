<div align="center">

# kb

### Build and maintain topic-based knowledge bases.

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![CI](https://img.shields.io/github/actions/workflow/status/pedronauck/kodebase-go/ci.yaml?branch=main&label=CI)](https://github.com/pedronauck/kodebase-go/actions)
[![Go](https://img.shields.io/badge/Go-1.24-00ADD8.svg)](https://go.dev/)

[Install](#install) &#8226; [See It Work](#see-it-work) &#8226; [Features](#features) &#8226; [Commands](#commands) &#8226; [Contributing](#contributing)

</div>

---

`kb` is a single-binary Go CLI for building and maintaining topic-based knowledge bases in the [Karpathy KB](https://github.com/karpathy) pattern. It handles the non-LLM workflow: topic scaffolding, multi-source ingestion (URLs, files, YouTube, codebases, bookmarks), structural linting, codebase analysis, QMD indexing/search, and KB-oriented inspection commands. LLM compilation stays in your agent layer.

No SaaS. No cloud. Just markdown.

---

## Install

> [!NOTE]
> `kb` works even better with its companion skill in [`skills/`](skills/). Install it with `npx skills add https://github.com/compozy/kb --skill kb`.

### Homebrew

```bash
brew install compozy/kb/kb
```

### npm

```bash
npm install -g @compozy/kb
```

### Go

```bash
go install github.com/compozy/kb/cmd/kb@latest
```

### Build from source

```bash
git clone https://github.com/compozy/kb.git
cd kb
make build
# binary is at bin/kb
```

**Optional** -- for semantic search capabilities:

```bash
npm install -g @tobilu/qmd
```

> [!NOTE]
> **Requirements:** Go >= 1.24. The `search` and `index` commands require [QMD](https://github.com/tobilu/qmd) to be installed separately. The `ingest url` command requires a [Firecrawl](https://firecrawl.dev) API key. The `ingest youtube --stt` fallback requires an [OpenRouter](https://openrouter.ai) API key. All other commands work standalone.

<details>
<summary><strong>What it touches</strong></summary>

- **Creates files** in `.kb/vault/` inside the target repository (or a custom `--vault` path)
- **Reads** source files in the target repository (never modifies them)
- **Network calls** -- `ingest url` calls the Firecrawl API; `ingest youtube --stt` calls OpenRouter. All other commands are fully local.
- **No telemetry** -- nothing is sent anywhere
- **Uninstall:** Remove the `kb` binary from your `PATH` and delete the `.kb/` directory

</details>

---

## See It Work

Create a topic and ingest content from multiple sources:

```bash
# scaffold a new topic
$ kb topic new rust-lang "Rust Language" programming
{
  "slug": "rust-lang",
  "title": "Rust Language",
  "domain": "programming"
}

# ingest a web article
$ kb ingest url https://doc.rust-lang.org/book/ch04-01-what-is-ownership.html --topic rust-lang

# ingest a local PDF
$ kb ingest file ./rust-reference.pdf --topic rust-lang

# ingest a YouTube video transcript
$ kb ingest youtube https://www.youtube.com/watch?v=... --topic rust-lang

# ingest a codebase snapshot with full analysis
$ kb ingest codebase ./my-rust-project --topic rust-lang

# lint the topic for structural issues
$ kb lint rust-lang
```

Analyze codebase snapshots from the terminal:

```bash
$ kb inspect complexity --top 5

 symbol_name       | cyclomatic_complexity | loc | source_path
 computeMetrics    | 12                    | 89  | src/compute-metrics.ts
 parseTypeScript   | 9                     | 67  | src/adapters/typescript.ts
 normalizeGraph    | 8                     | 45  | src/normalize-graph.ts
 renderDocuments   | 7                     | 112 | src/render-documents.ts
 scanWorkspace     | 6                     | 54  | src/scan-workspace.ts
```

```bash
$ kb inspect dead-code

 kind   | name           | source_path                | reason
 symbol | oldHelper      | src/utils.ts               | dead-export
 file   | unused-cfg.ts  | src/config/unused-cfg.ts   | orphan-file
```

Search your vault with natural language (requires QMD):

```bash
$ kb search "error handling patterns" --limit 3

 title                  | score | path
 Error Handling Guide   | 0.89  | wiki/concepts/Error Handling.md
 src/error-boundary.ts  | 0.74  | raw/codebase/files/src/error-boundary.ts.md
 handleError            | 0.61  | raw/codebase/symbols/handleError--src-utils-l42.md
```

---

## Features

**Topic-based knowledge bases** -- `kb` organizes knowledge into topics, each with its own `raw/`, `wiki/`, `outputs/`, and `bases/` directories. Scaffold a topic with `kb topic new`, then ingest content from any supported source.

**Multi-source ingestion** -- Ingest web articles via Firecrawl, local files (PDF, DOCX, XLSX, PPTX, EPUB, HTML, CSV, JSON, XML, plain text, images with OCR), YouTube transcripts (with optional OpenRouter STT fallback), codebases, and bookmark clusters. Each source type goes through a converter registry that normalizes content to frontmatter-annotated markdown.

**Codebase analysis** -- Point `kb ingest codebase` at a repository and it generates an [Obsidian](https://obsidian.md) vault layer with every symbol, file, and dependency relationship mapped into interconnected markdown notes. It computes cyclomatic complexity, blast radius, coupling, instability, and dead code detection, then compiles wiki articles and interactive [Base](https://obsidian.md/blog/bases/) views.

**Structural linting** -- `kb lint` checks topics for missing frontmatter, broken wikilinks, orphaned files, and other structural health issues. Reports can be saved as markdown in `outputs/reports/`.

**10 inspect subcommands** -- Query your codebase like a database, from the terminal. Rank functions by complexity, find dead exports nobody imports, trace dependency chains, detect circular imports. Three output formats (table, JSON, TSV), zero external dependencies.

**Obsidian-native output** -- Every generated note uses wikilinks, YAML frontmatter, and backlink-aware cross-references. Base views give you filterable, sortable tables inside Obsidian -- Symbol Explorer, Complexity Hotspots, Danger Zone, Module Health, and more.

**Semantic search** -- Index your vault with QMD for hybrid, lexical, or vector search across all documentation. Useful for onboarding, architecture review, or feeding context to LLMs.

**AI-friendly output** -- Every vault includes `CLAUDE.md` and `AGENTS.md` as schema documents. The structured markdown output is designed for direct consumption by LLMs -- frontmatter metadata, consistent structure, and explicit cross-references.

**Single binary** -- No runtime dependencies. One `kb` binary handles everything. Built in Go for fast startup and low memory usage.

---

## Why kb

Most code analysis tools give you a dashboard. `kb` gives you a knowledge base.

**SonarQube and CodeClimate** tell you your code has problems. `kb` tells you which problems, where they connect, and gives you a structured workspace to reason about them. The output is markdown you own, not a SaaS dashboard that disappears when you cancel the subscription.

**Sourcegraph** is excellent for code search across repositories. `kb` is for understanding a single repository deeply -- its architecture, its coupling patterns, its risk surface -- and building a persistent knowledge artifact around that understanding.

**Obsidian and Notion** are great note-taking tools. `kb` automates the scaffolding and ingestion so you start with structure instead of a blank page, then extend the vault with your own notes and analysis.

The key difference: `kb` outputs compound. A SonarQube scan from last month is stale data. A `kb` vault from last month is a knowledge base you've been building on for 30 days.

---

## Commands

### `kb topic`

Scaffold and manage knowledge base topics.

```bash
kb topic new <slug> <title> <domain>   # Create a new topic
kb topic list                           # List all topics in the vault
kb topic info <slug>                    # Show metadata for a topic
```

### `kb ingest`

Ingest source material into an existing topic.

```bash
kb ingest url <url> --topic <slug>                # Scrape a web URL (requires Firecrawl)
kb ingest file <path> --topic <slug>              # Convert and ingest a local file
kb ingest youtube <url> --topic <slug> [--stt]    # Extract a YouTube transcript
kb ingest codebase <path> --topic <slug>          # Analyze a codebase
kb ingest bookmarks <path> --topic <slug>         # Ingest a bookmark-cluster markdown file
```

**Supported file formats** for `ingest file`: PDF, DOCX, XLSX, PPTX, EPUB, HTML, CSV, JSON, XML, plain text (`.txt`, `.md`), and images (PNG, JPG, TIFF, BMP, GIF -- with optional OCR via Tesseract).

| Ingest subcommand | `--topic` | Additional flags |
| --- | --- | --- |
| `url` | required | -- |
| `file` | required | -- |
| `youtube` | required | `--stt` (enable OpenRouter STT fallback) |
| `codebase` | required | `--include`, `--exclude`, `--semantic`, `--progress`, `--log-format` |
| `bookmarks` | required | -- |

### `kb lint`

Check a topic for structural KB issues.

```bash
kb lint [<slug>] [--format table|json|tsv] [--save] [--topic <slug>]
```

`--save` writes a markdown report to `outputs/reports/<date>-lint.md`.

### `kb inspect`

Query codebase vault data using frontmatter and extracted metrics.

```bash
kb inspect <subcommand> [options]
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

### `kb search`

Semantic search across vault documents. Requires [QMD](https://github.com/tobilu/qmd).

```bash
kb search <query> [options]
```

| Flag           | Default | Description                              |
| -------------- | ------- | ---------------------------------------- |
| `--lex`        | false   | Use BM25 keyword search only             |
| `--vec`        | false   | Use vector similarity search only        |
| `--limit`      | 10      | Maximum results to return                |
| `--full`       | false   | Show full document instead of snippet    |
| `--min-score`  | --      | Minimum similarity threshold             |
| `--all`        | false   | Return all matches above threshold       |
| `--collection` | --      | Explicit QMD collection name             |
| `--format`     | table   | Output format: `table`, `json`, or `tsv` |

Default mode is hybrid (BM25 + vector similarity). Use `--lex` or `--vec` to restrict.

### `kb index`

Create or update a QMD collection for semantic indexing. Requires [QMD](https://github.com/tobilu/qmd).

```bash
kb index [options]
```

| Flag            | Default | Description                                |
| --------------- | ------- | ------------------------------------------ |
| `--vault`       | --      | Vault root path                            |
| `--topic`       | --      | Topic slug inside vault                    |
| `--name`        | --      | Override derived QMD collection name       |
| `--embed`       | true    | Run embedding after syncing files          |
| `--context`     | --      | Attach context to improve search relevance |
| `--force-embed` | false   | Force re-embedding of all documents        |

### `kb version`

Print build version metadata.

---

## What Gets Generated

```text
.kb/vault/
  <topic-slug>/
    CLAUDE.md                    # Schema document for LLMs
    AGENTS.md                    # Agent-facing project reference
    log.md                       # Append-only operation log
    raw/
      codebase/                  # Machine-generated codebase snapshot
        files/                   #   One markdown note per source file
        symbols/                 #   One markdown note per extracted symbol
        indexes/
          directories/           #   Directory-level inventories
          languages/             #   Language-level inventories
      <ingested-sources>.md      # Ingested articles, transcripts, documents
    wiki/
      concepts/                  # Synthesized wiki articles (codebase topics)
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
        Source Index.md          # Reverse index
    outputs/                     # Your analysis outputs (preserved)
      reports/                   # Lint reports (when --save is used)
    bases/                       # Obsidian Base views
```

- **`raw/`** -- Machine-generated source snapshots and ingested documents. Codebase snapshots refresh on every run; ingested documents are append-only.
- **`wiki/`** -- Starter articles synthesized from codebase metrics. Managed areas refresh; your additions are preserved.
- **`outputs/`** -- A place for your own briefings, queries, diagrams, and reports. Never touched by `kb`.
- **`bases/`** -- Obsidian Base `.base` files for interactive table/card/list views of metrics.

---

## Supported Languages (Codebase Analysis)

| Language   | Extensions    | Parser        | Relation Confidence |
| ---------- | ------------- | ------------- | ------------------- |
| TypeScript | `.ts`, `.tsx` | `tree-sitter` | syntactic           |
| JavaScript | `.js`, `.jsx` | `tree-sitter` | syntactic           |
| Go         | `.go`         | `tree-sitter` | syntactic           |

Want to add a language? See [CONTRIBUTING.md](CONTRIBUTING.md#adding-a-new-language-adapter).

---

## Supported File Formats (Ingest)

| Format     | Extensions                        | Notes                                  |
| ---------- | --------------------------------- | -------------------------------------- |
| PDF        | `.pdf`                            | Native text extraction via pdfcpu      |
| DOCX       | `.docx`                           | XML-based extraction                   |
| XLSX       | `.xlsx`                           | Sheet-to-markdown table conversion     |
| PPTX       | `.pptx`                           | Slide text extraction                  |
| EPUB       | `.epub`                           | Chapter extraction with HTML-to-MD     |
| HTML       | `.html`, `.htm`                   | HTML-to-markdown conversion            |
| CSV        | `.csv`                            | Table conversion                       |
| JSON       | `.json`                           | Pretty-printed code block              |
| XML        | `.xml`                            | Pretty-printed code block              |
| Plain text | `.txt`, `.md`                     | Pass-through                           |
| Images     | `.png`, `.jpg`, `.tiff`, `.bmp`, `.gif` | Optional OCR via Tesseract       |

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

## Configuration

`kb` is primarily configured through CLI flags. Optional runtime configuration is loaded from a TOML file and environment variables.

| Variable            | Source     | Description                                |
| ------------------- | ---------- | ------------------------------------------ |
| `APP_CONFIG`        | env        | Path to TOML config file                   |
| `FIRECRAWL_API_KEY` | env / TOML | Firecrawl API key for `ingest url`         |
| `FIRECRAWL_API_URL` | env / TOML | Firecrawl API endpoint                     |
| `OPENROUTER_API_KEY`| env / TOML | OpenRouter API key for `ingest youtube --stt` |
| `OPENROUTER_API_URL`| env / TOML | OpenRouter API endpoint                    |

See [`config.example.toml`](config.example.toml) for the full TOML schema.

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
| `make build`           | Build binary to `bin/kb`                 |
| `make verify`          | fmt -> lint -> test -> build -> boundaries|
| `make deps`            | Run `go mod tidy`                        |

See [CONTRIBUTING.md](CONTRIBUTING.md) for code style, testing requirements, and how to add a new language adapter.

---

## Contributing

`kb` is MIT-licensed and built in the open. We welcome contributions of all kinds:

- **Language adapters** -- Add support for Python, Rust, Java, or any language with a tree-sitter grammar
- **File converters** -- Add support for new file formats in the converter registry
- **New code smell detectors** -- The metrics engine is designed to be extended
- **Wiki article templates** -- Better starter articles mean better vaults out of the box
- **Bug reports and feature requests** -- [Open an issue](https://github.com/pedronauck/kodebase-go/issues), we read them all

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

---

## License

MIT -- see [LICENSE](LICENSE).
