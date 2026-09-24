---
id: 5
title: qmd as the candidate generator
labels: [wayfinder:research]
status: closed
assignee: "claude"
blocked_by: []
---

## Question

Can qmd give kb the candidates it needs, per document, cheaply?

Pin down, from `internal/qmd` and qmd's own docs/CLI: top-k neighbors for a whole document (not a query string), lexical vs vector vs hybrid per call, JSON output fields (path, score, snippet), latency for a 1k-document collection, and behavior when the index is stale or missing. Record findings in `research/qmd-candidates.md`.

## Resolution

qmd can supply candidates, but only in one narrow shape. It has no "neighbors of doc X" API, so kb has to turn each document into query text.
- Lexical search ANDs every term, so paragraph-length queries return 0–1 hits. Lex only works for short title or alias lines.
- The default `query` and `vsearch` run a local 1.7B query-expansion LLM, and `query` also runs a 0.6B reranker. Uncached, that is 5–29 s per call on an M4 Max, so neither fits a per-document loop.
- Cheap path: structured `vec:` (+ `lex:` title/alias lines) with `--no-rerank`. Per call that is 1.2 s over the CLI (model load) and 0.08 s (130 docs) / 0.76 s (1k docs) through a run-scoped `qmd mcp --http` server.
- Results include the document itself and hub files (`log.md`, indexes). kb must filter them, map `qmd://` paths back to vault paths (URL-decode, strip `?index=`), and ignore qmd titles and scores for banding.
- Staleness never fails loudly: deleted files are still returned, edits are invisible until `update`, missing embeddings give `[]` with a warning, and a missing named index is silently created. kb must run `update` and `embed` first. Embedding took 59 s for 130 docs and 15 min for 1k docs/74 MB.
- Fallback without qmd: a pure-Go alias dictionary scan for mentions plus in-memory OR-semantics BM25/TF-IDF for pairs. qmd `vec:` neighbors are an optional recall layer merged on top.

Details, commands, and numbers: [research/qmd-candidates.md](../research/qmd-candidates.md)
