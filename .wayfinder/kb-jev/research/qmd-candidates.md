# qmd as the candidate generator

Ticket: [05 — qmd as the candidate generator](../tickets/05-qmd-candidates.md)

Sources: `internal/qmd/client.go`, `internal/cli/{search,index}.go`, `skills/kb/references/cli-search-index.md`, qmd 2.8.3 (`facd35e`): `qmd --help`, the packaged README, and the packaged `dist/` source (`store.js`, `llm.js`, `mcp/server.js`, `cli/qmd.js`). I did not use ctx7; the installed package source is the ground truth for this version.

## Short answer

Yes, but only in one narrow form. qmd has **no "neighbors of document X" API**, so kb has to turn each document into query text. Only two call shapes are cheap enough for per-document candidate generation:

1. **Structured `vec:` query with rerank off** (`qmd query --no-rerank 'vec: <text>'`, or HTTP `POST /query` with `rerank:false`). This skips query expansion and reranking.
2. **Short `lex:` lines** built from the title, aliases, and key terms, fused by RRF in the same call.

The defaults (`qmd query`, `qmd vsearch`) run a local 1.7B query-expansion LLM, and `query` also runs a 0.6B reranker. Uncached, that costs 5–29 s per call, so they don't fit a per-document loop.

## How kb calls qmd today

- `QMDClient` shells out per call: `search` (lexical), `vsearch` (vector), `query` (hybrid). It adds `--json`, `-n`, `--min-score`, `--full`, and `-c <collection>`. There is no `--no-rerank`, `-C`, or structured `lex:/vec:` support, and no `--intent`.
- kb always uses the **global default index** (`~/.cache/qmd/index.sqlite`). `NewClient()` is never given `WithIndexName`. The collection name is the topic slug.
- Hybrid mode first runs a probe, `qmd vsearch "__kb_vector_probe__"`, to detect sqlite-vec. Because `vsearch` runs query expansion, this probe alone cost **2.46 s uncached / 1.16 s cached** on each hybrid `kb search`.
- Results are parsed from `docid, score, file|filepath|displayPath, title, snippet|body|bestChunk`. Stderr is treated as diagnostics.

## Answers to the ticket's questions

### Top-k neighbors for a whole document

- **No native API.** `findSimilarFiles` in `store.js` does Levenshtein matching on file *paths* (used for "did you mean"), not content. There is no more-like-this and no search by docid. `get "#docid"` only retrieves a document.
- **Passing the document as the query:**
  - **Lexical (`search`, `lex:`):** `buildFTS5Query` **ANDs every term**, so a title-plus-paragraph query returns 0 or 1 hit (the document itself). The word `OR` is not an operator. In the measurements, the paragraph query returned 1, 0, and 1 hits. Lex is only useful with 1–3 term queries: the title, aliases, or key phrases. Each `lex:` line is its own ranked list, fused by RRF. Lex lines must be single-line with balanced quotes.
  - **Vector (`vsearch`, `vec:`):** the query is embedded with embeddinggemma-300M and truncated to `QMD_EMBED_CONTEXT_SIZE`, default 2048 tokens, with a warning on stderr. Search runs over **chunks** (about 900 tokens, 15% overlap), and results are deduplicated per document. Structured lines must be single-line: kb has to collapse newlines.
  - **Command-line length:** the query is an argv string. The OS limit (about 1 MB on macOS) is far above the 2048-token embedding cap. The HTTP body has no limit I found. (UNVERIFIED: whether very long `vec:` text degrades quality before it hits the cap.)
- **Self-match:** the source document is in the results and must be filtered out. Recall of itself is poor. With title plus first paragraph as a `vec:` query:
  - jev (130 docs): self was at rank 0 for 7/14 articles and missing from the top 10 for 2/14.
  - 1k-doc topic: self was missing from the top 10 for 4/20.

  The top hits are often hubs like `log.md`, `Dashboard.md`, `Concept Index.md`, and `CLAUDE.md`. kb has to exclude non-article paths, either post-hoc or with a scoped collection or glob.

### Lexical, vector, or hybrid per call

| Call | What runs | LLM work |
| --- | --- | --- |
| `qmd search` | FTS5 BM25 (AND of terms) | none |
| `qmd vsearch` | **query expansion (1.7B)** + embedding of the original and vec/hyde variants; default `minScore 0.3` | expansion |
| `qmd query "<text>"` | BM25 probe. A "strong signal" skips expansion; otherwise it runs expansion → FTS + vector → RRF → rerank of top `-C` (default 40) chunks → position-aware blend | expansion + rerank |
| `qmd query 'lex:..\nvec:..'` | Structured: no expansion. FTS for lex lines + vector for vec/hyde → RRF → rerank | rerank only |
| `... --no-rerank` | RRF scores only | none (the embedding model loads for vec) |

Models (local GGUF via node-llama-cpp, cached in `$XDG_CACHE_HOME/qmd/models`): embeddinggemma-300M-Q8 (313 MB), qwen3-reranker-0.6B-Q8 (610 MB), qmd-query-expansion-1.7B-q4 (1.2 GB). Expansion and rerank results are cached in the index's `llm_cache` table, so repeating the same query is cheap. Everything runs on-device with no API cost; the cost is latency plus GPU/CPU time.

### JSON output fields

`--json` (a legacy alias of `--format json`, still accepted) returns a flat array:

```json
{"docid":"#1b6d15","score":0.75,"file":"qmd://jev/wiki/concepts/Trava DataJud e OAB.md","line":2,"title":"Trava DataJud e OAB","snippet":"@@ -1,4 @@ ..."}
```

- `--full` replaces `snippet` with `body` (the full text). `--explain` (query) adds FTS, vector, RRF, and rerank traces.
- HTTP `/query` returns `{results:[{docid,file,title,score,context,line,snippet}]}` with a URL-encoded `file`.
- Caveats:
  - With a **named index**, `file` gets an `?index=<name>` suffix (`qmd://st/b.md?index=st`). kb's `normalize` would carry that into `Path`.
  - Titles can be wrong for raw sources: `sei-15707-2026.md` came back as title `"3"`. kb should use its own frontmatter title, not qmd's.
- Scores are not comparable across modes:
  - BM25 is normalized and can be `0` on tiny corpora.
  - Vector score is `1/(1+distance)`.
  - `--no-rerank` returns rank-shaped RRF scores (1.00, 0.50, 0.33, 0.25…).
  - Reranked scores are blended.

  Calibrated banding therefore belongs to the decision model, not to qmd scores.

### Latency

Setup: Apple M4 Max; `HOME`, `XDG_CACHE_HOME`, and `XDG_CONFIG_HOME` isolated under `/tmp/qmdbench`, with the models dir symlinked read-only. Query text is the article title plus its first paragraph (206–485 chars).

**Indexing**

| Corpus | `collection add` (FTS) | `embed` |
| --- | --- | --- |
| jev, 130 docs, 3.6 MB | 0.34 s | 2,256 chunks in **59 s** |
| storytelling, 1,023 docs, 74 MB | 2.5 s | 19,836 chunks in **15 m 13 s** |

**Per call, CLI (a new process each time, jev)**

| Call | Uncached | Cached / repeat |
| --- | --- | --- |
| `search` (lex, paragraph) | 0.12–0.14 s (0–1 hits) | same |
| `vsearch` | 2.4–6.6 s | 1.3–2.0 s |
| `query` (default hybrid) | **25.6–28.7 s** first call; 5.7–11.9 s others | 1.2–1.6 s |
| `query --no-rerank` | 11.1–12.0 s (expansion dominates) | 1.2–1.6 s |
| `query 'vec: …'` (rerank C=40) | 8.6–10.0 s | 1.2 s |
| `query -C 10 'vec: …'` | 7.3–7.6 s | — |
| `query --no-rerank 'vec: …'` | **1.19–1.21 s** (≈1.1 s is embed-model load) | same |
| `query --no-rerank 'lex:…\nlex:…\nvec:…'` | 1.5–2.0 s | same |
| 1k docs: `search` | 0.34–0.54 s | — |
| 1k docs: `query --no-rerank 'vec: …'` | 1.76–1.84 s | — |

**Per call, HTTP server (`qmd mcp --http`, models stay loaded)**

| Corpus / call | Mean | Median | First call |
| --- | --- | --- | --- |
| jev, `vec`, rerank off (14 articles) | **0.08 s** | 0.07 s | 0.12 s (1.0 s on a cold server) |
| jev, `lex:title` + `vec`, rerank off | 0.14 s | 0.13 s | — |
| jev, `vec`, rerank C=10 (uncached) | 5.2 s | 4.9 s | — |
| jev, `vec`, rerank C=40 (uncached) | 8.7 s | 7.9 s | — |
| 1k docs, `vec`, rerank off (20 docs) | **0.76 s** | 0.71 s | 1.69 s |
| 1k docs, `vec`, rerank C=10 / C=40 (5 docs) | ≈2.0–2.1 s | — | — |

The 1k-doc rerank numbers are partly rerank-cache-warm and should be treated as a lower bound.

**Extrapolated cost of a full relink** (one call per document):

| Path | 130 docs | 1k docs |
| --- | --- | --- |
| HTTP, rerank off | ≈10 s | ≈13 min |
| CLI, rerank off | ≈2.6 min | ≈30 min |
| Default `query` | ≈15–60 min | hours |

UNVERIFIED: why 1k-doc vector search is about 10× slower than 130 docs. sqlite-vec `vec0` looks like brute-force KNN over about 20k chunks, which scales linearly.

### Stale or missing index

| Situation | Behavior |
| --- | --- |
| Named index that doesn't exist | **Silently created**, returns `[]` with rc 0. kb can't tell "no index" from "no match" without checking `status` first. |
| Unknown collection (`-c nope`) | `Collection not found: nope`, rc 1 |
| Documents indexed but not embedded | `vsearch` and `vec:` return `[]`, rc 0, plus a stderr warning (`N documents (100%) need embeddings`). `query` quietly falls back to lexical results. |
| File modified after indexing, no `update` | Search sees the **old content**: the new term isn't found and the old one is. |
| File deleted, no `update` | Still returned as a hit. Only `--full-path` warns on stderr (`could not resolve the file on disk … Run 'qmd update'`). |
| `qmd` binary missing | kb returns `ErrQMDUnavailable` with the install hint. |
| After `qmd update` | Reports `Indexed: 0 new, 1 updated, 0 unchanged, 1 removed`. New chunks still need `qmd embed`, which is incremental by content hash. |
| Rename or removal | The HTTP `/query` endpoint and CLI results keep the stale `qmd://` path until `update`. |

In short: qmd never fails loudly on staleness. kb has to run `qmd update` (sub-second) and `qmd embed` (incremental) before a candidate pass, or check `status` for pending embeddings.

### MCP / server mode

- Transports:
  - `qmd mcp` runs over stdio.
  - `qmd mcp --http [--port N] [--daemon]` adds `POST /mcp`, `POST /query` (alias `/search`), and `GET /health`.
- `/query` takes `{searches:[{type,query}], collections:[..], limit, minScore, candidateLimit, intent, rerank}`:
  - `searches` is required; the HTTP endpoint has no auto-expansion path.
  - It returns up to 10 sub-queries, and the first sub-query gets 2× weight.
  - Unknown parameters are silently ignored.
  - Only `collections` (plural) scopes the search.
- Models stay loaded. Contexts are disposed after 5 min idle, with about a 1 s penalty to recreate them.
- The server serves one index per process (`qmd --index s1k mcp --http`). The daemon's PID file lives in `$XDG_CACHE_HOME/qmd/mcp.pid`, so it's shared with the user's own daemon if one runs on the default port. kb would need its own port and lifecycle.
- The node SDK (`createStore`, `searchVector`, `searchLex`) exists but is JS-only, so it isn't usable from Go without a sidecar.

## Implications for kb candidate generation

1. **Don't use kb's current `Search` path for candidates.** Hybrid mode adds the expansion probe and full rerank (5–29 s uncached per document), and `lex` with document text returns almost nothing. Add a dedicated candidate call that sends a structured query (`vec: <title + lead, single line>` plus `lex:` lines for the title and each alias) with `--no-rerank`. kb's decision model already judges the pairs, and qmd's reranker would duplicate that.
2. **For batches, run the qmd HTTP server for the duration of the run** (start it on a private port, poll `/health`, `POST /query` per document, stop it at the end). This is about 15× cheaper than one CLI process per document (0.08 s vs 1.2 s at 130 docs). The CLI with `--no-rerank` is fine for single-document runs, such as linking one new ingest.
3. **Treat qmd output as a recall set, not a ranking.**
   - Take the top about 20.
   - Drop the document itself.
   - Drop non-article paths (`wiki/index/*`, `log.md`, `CLAUDE.md`/`AGENTS.md`, `outputs/`) or scope the collection to `wiki/`.
   - Map `file` → vault path by stripping `qmd://<collection>/` and `?index=`, then URL-decoding.
   - Ignore qmd's title and score for banding.
4. **Always make the index fresh first:** `kb index` (update + incremental embed) before a candidate pass. If `status` reports pending embeddings, degrade to lex-only rather than silently getting `[]` from `vec`. Initial embedding costs about 1 min per 130 docs (M4 Max), about 15 min for a 1k-doc/74 MB topic. That belongs in the backfill budget (the map's "Not yet specified").
5. **Use the global index with care.** kb currently writes into the user's global qmd index (2.9 GB here). Candidate generation should keep using the topic collection and must not create named indexes, which are created silently and change the result paths.
6. **Mention-level candidates** (an alias or article title appearing in another document's body) are not a qmd job. BM25 AND semantics can confirm that an exact term occurs (`lex: "<alias>"` returns the documents containing it), but kb gets that more cheaply and exactly by scanning markdown in Go.

### Fallback when qmd is absent (pure Go)

The ticket's "everything is optional" rule means kb must still produce candidates without qmd:

- **Alias dictionary scan (mention candidates):** build a map from article titles, filenames, and frontmatter `aliases` to the article. Scan each document body with Aho-Corasick, or a normalized-token trie with accent and case folding (the corpus is PT/EN). Skip existing `[[links]]`, code fences, and frontmatter. This is exact, cheap (ms per document), and is the main source of Link candidates regardless of qmd.
- **In-process BM25 / TF-IDF (pair candidates):** tokenize wiki articles and sources, then score document-to-document similarity with BM25 over the top-N TF-IDF terms of the query document, with OR semantics. Unlike qmd lex, this avoids the AND problem. At 1k documents it fits in memory and needs no dependency. A small hand-written index is enough; if an index is persisted, it goes under `.kb/` as a derived, rebuildable artifact (standing decision).
- **Optional semantic layer:** only through qmd (local embeddings) or a remote embedding API behind the same "optional" switch. The pure-Go path has no embeddings, which is acceptable because the decision model filters candidates anyway.
- Proposed order: alias scan (always) → BM25 pairs (always) → qmd `vec:` neighbors (when qmd is present and the index is fresh), with candidates merged and deduplicated before the decision model sees them.

## Commands used (reproducible)

```bash
W=/tmp/qmdbench; mkdir -p $W/home $W/cache/qmd $W/config
ln -s ~/.cache/qmd/models $W/cache/qmd/models        # reuse model files only; index is isolated
export HOME=$W/home XDG_CACHE_HOME=$W/cache XDG_CONFIG_HOME=$W/config NO_COLOR=1
qmd collection add /Users/pedronauck/Dev/courses/pedronauck/research/jev --name jev   # absolute path: ~ expands to the fake HOME
qmd embed
qmd search  --json -n 10 -c jev "<title. first paragraph>"
qmd vsearch --json -n 10 -c jev "<...>"
qmd query   --json -n 10 -c jev "<...>"
qmd query   --json -n 10 --no-rerank -c jev "vec: <...>"
sqlite3 $W/cache/qmd/index.sqlite "delete from llm_cache"   # force uncached timings
qmd mcp --http --port 18282 &                                 # then POST /query {searches,collections,limit,rerank}
qmd --index s1k collection add .../research/storytelling --name story && qmd --index s1k embed   # 1k-doc run
```

Scripts: `/tmp/qmdbench/bench.sh`, `/tmp/qmdbench/bench_http.py`, `/tmp/qmdbench/bench1k.py` (temporary).

## UNVERIFIED

- Latency on machines without Metal GPU. qmd says `--no-rerank` is "much faster on CPU", so expansion and rerank would be far slower there.
- Whether `vec:` quality drops for long inputs before the 2048-token cap. I only tested title plus first paragraph.
- Whether vector search scales linearly beyond 1k docs; the brute-force sqlite-vec guess is not profiled.
- Neighbor quality compared with a Go BM25 baseline. This needs a labeled pair set: 14 jev articles are too few.
- One `qmd --help` ran against the real HOME before isolation. It only prints help and doesn't index; the real index mtime (16:39) predates the benchmark runs.
