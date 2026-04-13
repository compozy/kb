# Ingest Codebase Command Reference

## Usage

```
kb ingest codebase <path> [flags]
```

The `<path>` argument is the root directory of the source repository to analyze (required).

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--topic` | string | `""` | Topic slug for the ingested codebase (derived from directory name if omitted) |
| `--output` | string | `""` | Vault root where the generated topic will be written. Defaults to `<path>/.kb/vault` |
| `--title` | string | `""` | Override the generated topic title |
| `--domain` | string | `""` | Override the generated topic domain |
| `--include` | string[] | `nil` | Re-include a path pattern that would otherwise be ignored; repeatable |
| `--exclude` | string[] | `nil` | Exclude an additional path pattern from scanning; repeatable |
| `--semantic` | bool | `false` | Enable semantic analysis when the underlying adapters support it |
| `--progress` | string | `auto` | Progress rendering mode: `auto`, `always`, or `never` |
| `--log-format` | string | `text` | Stderr event format: `text` or `json` |

## Non-Interactive Usage

When invoking from an agent context, always set `--progress never` to prevent TTY progress bars from corrupting stdout output.

```
kb ingest codebase /path/to/repo --topic my-project --progress never
```

## Pipeline Stages

The codebase ingestion pipeline executes these stages in order:

1. **scan** -- Discover source files by language
2. **select_adapters** -- Choose language parsers (tree-sitter for TS/JS, Go parser)
3. **parse** -- Extract AST nodes, symbols, and relations
4. **normalize** -- Merge per-file graphs into a unified snapshot, resolve imports
5. **metrics** -- Compute complexity, coupling, blast radius, dead code, smells
6. **render** -- Generate markdown documents and Base definitions
7. **write** -- Persist vault files to disk

## Supported Languages

| Language | Extensions | Adapter |
|----------|-----------|---------|
| TypeScript | `.ts` | tree-sitter |
| TSX | `.tsx` | tree-sitter |
| JavaScript | `.js` | tree-sitter |
| JSX | `.jsx` | tree-sitter |
| Go | `.go` | tree-sitter |

## Output Schema (GenerationSummary)

The command writes JSON to stdout. Parse the following fields:

```
{
  "command": string,           // always "generate"
  "rootPath": string,          // absolute path to the analyzed repository
  "vaultPath": string,         // absolute path to the vault root
  "topicPath": string,         // absolute path to the topic directory
  "topicSlug": string,         // topic identifier (use for --topic in later commands)
  "filesScanned": int,         // total files discovered
  "filesParsed": int,          // files successfully parsed
  "filesSkipped": int,         // files skipped (unsupported or excluded)
  "symbolsExtracted": int,     // total symbols extracted
  "relationsEmitted": int,     // total relation edges
  "rawDocumentsWritten": int,  // per-file markdown documents
  "wikiDocumentsWritten": int, // concept wiki articles
  "indexDocumentsWritten": int, // index pages
  "timings": {
    "scanMillis": int,
    "selectAdaptersMillis": int,
    "parseMillis": int,
    "normalizeMillis": int,
    "metricsMillis": int,
    "renderMillis": int,
    "writeMillis": int,
    "totalMillis": int
  },
  "diagnostics": [             // structured warnings/errors
    {
      "code": string,
      "severity": "warning" | "error",
      "stage": "scan" | "parse" | "render" | "write" | "validate",
      "message": string,
      "filePath": string?,
      "language": string?,
      "detail": string?
    }
  ]
}
```

## Vault Structure

After ingestion, the vault directory contains:

```
<vaultPath>/<topicSlug>/
  raw-codebase/     # One markdown file per source file with frontmatter and code
  wiki-concept/     # Compiled concept articles
  wiki-index/       # Index pages for navigation
  *.base            # Obsidian Base view definitions (YAML)
  CLAUDE.md         # Topic marker file
```

## Default Path Derivation

- If `--output` is omitted: vault path defaults to `<rootPath>/.kb/vault`
- If `--topic` is omitted: topic slug is derived from the repository directory name
- Full topic path: `<vaultPath>/<topicSlug>/`
