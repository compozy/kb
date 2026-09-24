---
id: 3
title: Prior art for decision-model linking in LLM wikis
labels: [wayfinder:research]
status: closed
assignee: "claude"
blocked_by: []
---

## Question

How do existing tools use Jev (or similar decision models) to link, tag, and dedupe notes, and what can kb copy?

Read: s1m (github.com/mikekelly/s1m), kaku.md autotag (@gemama0), jev-reranker (@hotchpotch), Granite (@Shpigford), Obsidian LLM Hub (@takeshy), and the tweet list gathered in this conversation. Focus on how they build the state, choose candidates, phrase questions, gate on confidence, and write results back. Record findings in `research/prior-art.md`.

## Resolution

Findings are in [research/prior-art.md](../research/prior-art.md). Code was read for s1m, jev-reranker, DocJev, Obsidian LLM Hub and an extra find, jev-the-janitor. Only public pages and tweets were available for kaku.md, Granite and @dekokun, so their details are marked UNVERIFIED.

- Every tool follows the same division of labour: code proposes Candidates (parsed links, a closed tag vocabulary, sibling titles, page boundaries) and Jev only judges them. Nothing invents targets.
- Recall favours one state per document with an indexed Noul per candidate (`links[{index}]`). s1m measured recall 0.83 for this against 0.57–0.72 for a single Choice over all links. Choice fits closed single-answer questions such as an OKF type, with a `none` or `other` option.
- The Confidence bands in use: DocJev decides at 0.5 and flags ±0.1 for review; the janitor uses a 0.55 floor plus the top-1 minus top-2 margin; s1m picked 0.6 from a labelled sweep. Scores are ordinal and not calibrated (a scent of 0.8–0.9 reached a gold page 22% of the time). kb should calibrate on a labelled set and store the question hash and model version with each decision.
- Write-back to copy is the janitor's: one namespaced frontmatter block, a byte-preserving atomic splice, quarantine by move plus a manifest, a review-pile file, and a `locked: true` field that stops re-judging. Exact duplicates are found by local hash; semantic duplicates need pairwise Nouls on pairs code has already shortlisted.
- Patterns to avoid: failing hard when the gateway errors (LLM Hub), gating on the argmax alone, compound questions, and asking the model for date arithmetic.
- OpenRouter exposes Jev at `/api/alpha/decisions` with model `~typesafe/jev-latest`. This was seen in LLM Hub code only, so it is UNVERIFIED.
