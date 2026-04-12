---
name: kodebase
description: Generates Obsidian knowledge vaults from source code repositories, inspects code metrics such as complexity, coupling, blast radius, dead code, and circular dependencies, indexes vault content for hybrid retrieval, and searches indexed vaults with lexical or vector queries. Use when analyzing a codebase for code quality, architecture health, symbol relationships, or code smells. Use when the task mentions kodebase, code vault, code knowledge base, code graph analysis, or code metrics inspection. Do not use for general code review, linting, formatting, building Go projects, or writing application code.
---

# Kodebase CLI

## Prerequisites

1. Verify the kodebase binary is available:
   ```
   kodebase version
   ```
2. For search and index commands, verify QMD is installed:
   ```
   qmd --version
   ```
   If missing, install with `npm install -g @tobilu/qmd`.
3. Supported source languages: TypeScript (`.ts`), TSX (`.tsx`), JavaScript (`.js`), JSX (`.jsx`), Go (`.go`).

## Workflow Overview

Kodebase operates as a pipeline. The `generate` command must run before any other command.

**Workflow A -- Code Analysis (no QMD required):**
```
kodebase generate <path> --> kodebase inspect <subcommand>
```

**Workflow B -- Full Pipeline (requires QMD):**
```
kodebase generate <path> --> kodebase index --> kodebase search <query>
```

The vault is stored at `<path>/.kb/vault/<topic-slug>/` by default. Later commands auto-discover this vault by walking up from the current working directory.

## Command Dispatch

Map the user's intent to the correct command:

| Intent | Command |
|--------|---------|
| Analyze a repository for the first time | `kodebase generate <path> --progress never` |
| Find code smells | `kodebase inspect smells --format json` |
| Find dead exports and orphan files | `kodebase inspect dead-code --format json` |
| Rank functions by complexity | `kodebase inspect complexity --format json` |
| Find high-impact symbols (blast radius) | `kodebase inspect blast-radius --min 5 --format json` |
| Find unstable files (coupling) | `kodebase inspect coupling --unstable --format json` |
| Find circular imports | `kodebase inspect circular-deps --format json` |
| Look up a specific symbol | `kodebase inspect symbol <name> --format json` |
| Look up a specific file | `kodebase inspect file <path> --format json` |
| Find what depends on X (incoming refs) | `kodebase inspect backlinks <name-or-path> --format json` |
| Find what X depends on (outgoing deps) | `kodebase inspect deps <name-or-path> --format json` |
| Search the codebase knowledge | `kodebase search "<query>" --format json` |
| Index vault for search | `kodebase index` |

## Step 1: Generate the Vault

Run the generate command to create the knowledge vault from source code.

```
kodebase generate <path> --progress never
```

Always use `--progress never` in agent contexts to prevent TTY progress bars from corrupting stdout.

Parse the JSON output from stdout to extract key values:
- `topicSlug` -- the topic identifier for later commands
- `vaultPath` -- absolute path to the vault root
- `topicPath` -- absolute path to the topic directory
- `filesScanned`, `filesParsed`, `symbolsExtracted` -- summary statistics
- `diagnostics` -- check for warnings or errors

Stderr carries structured stage logs. Do not treat stderr content as failure evidence.

Key flags:
- `--output <dir>` -- override vault root location
- `--topic <slug>` -- override the topic slug
- `--include <pattern>` -- re-include paths that would otherwise be ignored (repeatable)
- `--exclude <pattern>` -- exclude additional paths from scanning (repeatable)
- `--semantic` -- enable semantic analysis when adapters support it

Read `references/cli-generate.md` for the full flag table and output schema.

## Step 2: Inspect the Vault

Run inspect subcommands to analyze code quality and architecture.

**Shared flags for all inspect subcommands:**
- `--format json` -- always use JSON for programmatic parsing
- `--vault <path>` -- explicit vault root (omit to auto-discover from cwd)
- `--topic <slug>` -- explicit topic slug (omit if only one topic exists)

### Tabular Subcommands

These return a list of rows sorted by the primary metric:

1. **smells** -- List symbols and files with detected code smells.
   ```
   kodebase inspect smells --format json
   kodebase inspect smells --type high-complexity --format json
   ```

2. **dead-code** -- List dead exports (symbols with no incoming references) and orphan files (unreachable files).
   ```
   kodebase inspect dead-code --format json
   ```

3. **complexity** -- Rank functions/methods by cyclomatic complexity. Default top 20.
   ```
   kodebase inspect complexity --format json
   kodebase inspect complexity --top 50 --format json
   ```

4. **blast-radius** -- Rank symbols by transitive dependent count.
   ```
   kodebase inspect blast-radius --format json
   kodebase inspect blast-radius --min 10 --top 20 --format json
   ```

5. **coupling** -- Rank files by instability (Ce / (Ca + Ce)).
   ```
   kodebase inspect coupling --format json
   kodebase inspect coupling --unstable --format json
   ```

6. **circular-deps** -- List files participating in circular import chains. Returns a message row if no cycles exist.
   ```
   kodebase inspect circular-deps --format json
   ```

### Detail Lookup Subcommands

These return field-value pairs for a single matched entity:

7. **symbol \<name\>** -- Case-insensitive substring match. Returns detail fields for a single match, or a summary table for multiple matches.
   ```
   kodebase inspect symbol parseConfig --format json
   ```

8. **file \<path\>** -- Exact source path lookup. Use the source-relative path as stored in vault frontmatter.
   ```
   kodebase inspect file src/config.ts --format json
   ```

### Relation Subcommands

These return relation edges (`target_path`, `type`, `confidence`):

9. **backlinks \<name-or-path\>** -- Incoming references. Accepts a symbol name or file path.
   ```
   kodebase inspect backlinks parseConfig --format json
   ```

10. **deps \<name-or-path\>** -- Outgoing dependencies. Accepts a symbol name or file path.
    ```
    kodebase inspect deps src/config.ts --format json
    ```

Read `references/cli-inspect.md` for all column schemas and flag details.

## Step 3: Index the Vault (Optional)

Index the vault content into QMD for search. This step requires QMD on PATH.

```
kodebase index
```

The command is idempotent: it checks whether the collection already exists and chooses `add` (create) or `update` (refresh) automatically.

Key flags:
- `--embed` (default true) -- run embedding after syncing files
- `--force-embed` -- force re-embedding all documents
- `--context <text>` -- attach human context to improve search relevance
- `--name <name>` -- override the derived collection name

Parse the JSON output for:
- `updateResult.indexed` -- number of documents indexed
- `status.totalDocuments` -- total documents in the collection
- `status.hasVectorIndex` -- whether vector search is available
- `embedResult` -- embedding summary (null if `--embed=false`)

Read `references/cli-search-index.md` for the full output schema.

## Step 4: Search the Vault (Optional)

Search indexed vault content with QMD. Requires a prior `kodebase index` run.

```
kodebase search "<query>" --format json
```

**Search modes:**
- Hybrid (default) -- combines lexical and vector search
- Lexical (`--lex`) -- BM25 keyword search only
- Vector (`--vec`) -- embedding-based semantic search

The `--lex` and `--vec` flags are mutually exclusive. Omit both for hybrid mode.

Key flags:
- `--limit N` (default 10) -- maximum results
- `--min-score N` -- minimum relevance threshold
- `--full` -- return full document content instead of snippets
- `--all` -- return all matches above the minimum score

Output columns: `path`, `score`, `preview`.

Read `references/cli-search-index.md` for full details and example invocations.

## Output Format Selection

All `inspect` and `search` commands support `--format`:
- **json** -- always use for programmatic parsing
- **table** -- human-readable aligned columns (default)
- **tsv** -- tab-separated for piping to Unix tools

The `generate` and `index` commands always output JSON to stdout.

Read `references/output-formats.md` for format examples and empty result handling.

## Error Handling

Common errors and recovery:

| Error | Recovery |
|-------|----------|
| `unable to find a vault from <path>` | Run `kodebase generate <path>` first |
| `QMD is not available to kodebase` | Run `npm install -g @tobilu/qmd` |
| `no topics were found` | Run `kodebase generate` to populate the vault |
| `multiple topics were found` | Re-run with `--topic <slug>` |
| `no symbols matched "<query>"` | Use `inspect smells` or `inspect complexity` to discover valid names |
| `no file matched "<path>"` | Use exact source-relative path from vault frontmatter |

Read `references/error-handling.md` for the full error catalog with causes and recovery steps.

## Constraints

### MUST DO
- Run `kodebase generate` before any inspect, search, or index command
- Use `--format json` when parsing output programmatically
- Use `--progress never` when running `generate` in a non-interactive context
- Parse stdout only for command output; treat stderr as diagnostics
- Use the `topicSlug` from generate output for subsequent `--topic` flags

### MUST NOT DO
- Pass both `--lex` and `--vec` to `search`
- Pass `--force-embed` with `--embed=false` to `index`
- Treat stderr content as failure evidence for `generate`
- Assume vault location without running `generate` or checking for `.kb/vault/`
- Use relative paths like `./src/config.ts` for `inspect file` -- use `src/config.ts` instead
