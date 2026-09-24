---
id: 2
title: xiaomi/mimo-v2.6-flash as the generation model
labels: [wayfinder:research]
status: closed
assignee: "claude"
blocked_by: []
---

## Question

Is `xiaomi/mimo-v2.6-flash` on OpenRouter fit to be kb's generation model for aliases, `summary`/`description`, and contradiction rewrites?

Pin down: that the slug exists, pricing, context window, reasoning toggle and its default, structured/JSON output support, latency, and any reported quality in PT and EN for short extraction-style generation. Record findings in `research/cheap-llm.md`.

## Resolution

Fit, with conditions. The slug `xiaomi/mimo-v2.6-flash` exists (canonical `-20260921`, listed 2026-09-21). It costs $0.14 per 1M input and $0.28 per 1M output tokens, has a 1M-token context, and is served by Xiaomi and DeepInfra (both fp8).
- `response_format` and `structured_outputs` are supported. Reasoning is not mandatory: disable it with `reasoning: {enabled: false}` (`exclude` only hides it).
- Upstream Xiaomi turns thinking ON by default; OpenRouter's default is UNVERIFIED. kb must always send the disable flag and assert `reasoning_tokens == 0`.
- Latency and PT quality are UNVERIFIED. The model is 1 day old, benchmarks are vendor-only and agentic/code-focused, the model card lists only EN/ZH, and no key was available for live calls.
- Cost is about $0.0001 per alias or summary call with reasoning off. `xiaomi/mimo-v2.5` is a same-price fallback.
- Follow-up: run a small PT/EN alias and summary eval before locking in the default, and keep the slug configurable.

Details: [../research/cheap-llm.md](../research/cheap-llm.md)
