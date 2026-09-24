---
id: 1
title: Decision model API contract
labels: [wayfinder:research]
status: closed
assignee: "claude"
blocked_by: []
---

## Question

What exactly does the Jev Decisions API accept and return through OpenRouter (`POST /api/alpha/decisions`)?

Pin down: request and response schema for `choice`, `noul` and `score`; how option descriptions are passed; whether a `choice` can be multi-label or only returns a distribution; limits (questions per call, options per choice, 32k-token state, payload size); model slugs (`typesafe/jev-1.13` vs `~typesafe/jev-latest`) and pinning; error codes, known hangs and 503s; pricing; rate limits observed.

Sources: docs.typesafe.ai, OpenRouter API reference, and the local topic `~/Dev/courses/pedronauck/research/jev` (raw + wiki). Record findings in `research/decision-model-api.md`.

## Resolution

- Endpoint is `POST https://openrouter.ai/api/alpha/decisions`, not chat completions. The body is `{model, state, questions}` plus optional `provider`, `session_id`, `trace` and `user`. The response is `{id, model, provider, answers, usage{input_tokens, output_tokens, cost}}`. The shape is confirmed from recorded real responses; no live call was made because no key was available.
- Question types: `noul` returns `noul` (P(yes), no confidence). `choice` returns `choice`, `probabilities` (sums to 1) and `confidence`, over at most 255 options described by `criteria` (key and description both reach the model). `score` returns `score` (expected level index), `legend`, `probabilities` and `confidence`, over 2–10 ordered levels. A Choice is single-label only, so multi-label means one Noul per label.
- Limits: 64k tokens per request, with 32k for state plus the longest question. A 43k-token request was observed to succeed. There is no documented cap on questions per call, and the byte threshold for 413 is UNVERIFIED.
- Slugs: pin `typesafe/jev-1.13`. The response reports the dated permaslug `typesafe/jev-1.13-20260917`. The alias `~typesafe/jev-latest` exists but has not been tested.
- Errors: OpenRouter returns `{error:{code,message}}` with 400/401/402/403/404/413/429/500/502/503/524/529. Calls can hang, so set a timeout of about 15 s. Retry 408, 429 and 5xx with jittered backoff (2 retries). A timeout or invalid answer means undecided, never no.
- Pricing: $0.042 per 1M input tokens and free output; `usage.cost` gives the exact USD cost per call. Latency is p50 about 0.35 s. No 429 was seen at 12–48 concurrent calls.
- Details, examples, sources and kb implications: [decision-model-api.md](../research/decision-model-api.md).

## Addendum

The research-enrich session (4a04893e, 2026-09-22) saw 2 × 429 and 1 invalid receipt in a 16-worker run, and generation calls hanging >10 min behind OpenRouter keep-alive whitespace: retries and a total per-attempt deadline are required, not optional.
