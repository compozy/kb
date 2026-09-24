# Research: `xiaomi/mimo-v2.6-flash` as kb's generation model

Ticket: [02-cheap-llm](../tickets/02-cheap-llm.md) · Researched 2026-09-22 · Author: claude

**Short answer:** Yes, it fits, with conditions. The slug exists and costs about as little as any model does. It supports `response_format` / `structured_outputs`, and reasoning can be turned off. But it went on sale on OpenRouter **the day before this research (2026-09-21)**. Nobody outside Xiaomi has tested it yet, and it has **no PT evaluation**. No live calls were made, because `OPENROUTER_API_KEY` was not set in the environment. Latency and PT/EN output quality are **UNVERIFIED**.

## Existence / slug

- `xiaomi/mimo-v2.6-flash` exists. The pinned canonical slug is `xiaomi/mimo-v2.6-flash-20260921`. OpenRouter lists it as created 2026-09-21 (unix 1790021264). [OR models API]
- Weights: `XiaomiMiMo/MiMo-V2.6-Flash-RL` on Hugging Face, MIT licence. It is a MoE model with 309B total and 15B active parameters, and takes text, image, audio and video input. [OR models API; HF card]
- Other MiMo slugs on OpenRouter: `xiaomi/mimo-v2.6-pro`, `xiaomi/mimo-v2.6-pro-ultraspeed`, `xiaomi/mimo-v2.5-pro`, `xiaomi/mimo-v2.5`. `mimo-v2.5` costs the same as Flash and has been available longer, so it is the natural fallback.

## Pricing (OpenRouter, per 1M tokens)

| Slug | Input | Output | Cache read |
|---|---|---|---|
| mimo-v2.6-flash | $0.14 | $0.28 | $0.0028 |
| mimo-v2.5 | $0.14 | $0.28 | $0.0028 |
| mimo-v2.6-pro / v2.5-pro | $0.435 | $0.87 | $0.0036 |
| mimo-v2.6-pro-ultraspeed | $4.35 | $8.70 | $0.036 |

- Rough per-call cost for kb (my arithmetic, not measured): an alias or summary call with about 400 input and 150 output tokens costs about **$0.0001**. 10k calls cost about $1. If reasoning is left on and adds 500–2000 hidden tokens per call, cost goes up 3–10x and latency goes up a lot.
- OpenRouter reports `supports_implicit_caching: false` for both providers, so do not plan on cache discounts for short prompts.

## Context

- Context window is 1,048,576 tokens. Max completion is 131,072 tokens on Xiaomi's endpoint and 943,718 on DeepInfra's. Neither limit matters for kb's short prompts. [OR endpoints API]
- Two providers serve it, both at fp8: **Xiaomi** (1-day uptime 99.99%) and **DeepInfra** (1-day uptime 99.09%). [OR endpoints API]

## Reasoning toggle

- OpenRouter lists the `reasoning` and `include_reasoning` parameters, with `reasoning.mandatory: false`. That means reasoning can be disabled. [OR models API]
- **Default:**
  - On Xiaomi's own API, deep thinking is **"Enabled by default"** for mimo-v2.6-flash. You turn it off with `"thinking": {"type": "disabled"}`. [Xiaomi deep-thinking docs]
  - The OpenRouter models API does not publish a `default_enabled` field for this slug, so the default through OpenRouter is **UNVERIFIED**. It could differ by provider: vLLM-style hosts such as DeepInfra usually have thinking off unless `enable_thinking` is sent (**UNVERIFIED** for DeepInfra).
- **How to disable on OpenRouter:** send `"reasoning": {"enabled": false}`, or `{"effort": "none"}`. `{"exclude": true}` only hides the reasoning from the response; the model still reasons and you are still billed for it. [OR reasoning-tokens docs]
- Always send the disable flag explicitly. Do not rely on the default, since the providers may differ.
- While thinking is on, Xiaomi forces `temperature=1.0` and `top_p=0.95`. A custom temperature only works with thinking off. [Xiaomi docs]
- Tool calls across multiple turns with thinking on require sending `reasoning_content` back, or the API returns a 400. kb's use is single-turn with no tools, so this does not apply. [Xiaomi docs]
- Known pitfall: a Hermes Agent bug where the thinking-disable flag was never forwarded and tokens were wasted. This is a reason to check `usage.reasoning_tokens == 0` in tests. [hermes-agent#27325]

## Structured output

- Both providers list `response_format` and `structured_outputs` in `supported_parameters`, so JSON-schema mode should work through OpenRouter. [OR endpoints API]
- Vercel AI Gateway also lists JSON mode. [Vercel]
- Tool choice differs by provider. Xiaomi's endpoint does **not** support `tool_choice: function` or `none`; DeepInfra supports all of them. Prefer `response_format` over forced tool calls.
- Real JSON conformance: **UNVERIFIED** (no live call was made).

## Latency

- **UNVERIFIED.** OpenRouter's `latency_last_30m` and `throughput_last_30m` are still `null` because the model is too new.
- Tabbit reports about 140–160 tok/s on self-hosted vLLM/SGLang, not through OpenRouter.
- The ~134 tok/s Artificial Analysis figure you will see quoted is for **Pro**, not Flash. [orcarouter]
- With reasoning off and ~150 output tokens, sub-2-second calls are plausible, but that has not been measured.

## Quality notes

- All benchmarks are **vendor-reported** and focus on agentic tasks and code: DeepSWE 67.9, OSWorld-Verified 80.8, Terminal Bench 2.1 87.6, CyberGym 95.1. [HF card; kingy.ai]
- None have been reproduced independently. BenchLM does not rank it yet. [orcarouter; benchlm]
- **No multilingual or PT evaluation found.** The HF card lists only English and Chinese.
- No published evaluation covers short extraction-style generation (aliases, 2-sentence summaries, rewrites). PT quality is **UNVERIFIED**.
- The older MiMo-V2-Flash did well on SWE-bench Multilingual, but that measures programming languages, not natural languages.

## Risks

- **Brand new** (1 day old). Pricing, providers and defaults may change, and no one has independent data on it yet.
- **Reasoning default is unclear**, and upstream (Xiaomi) the default is ON. If kb forgets to disable it, costs go up silently and the model ignores the temperature kb sets.
- **Only two providers, both fp8.** DeepInfra's uptime is noticeably lower (99.1%).
- **PT quality unknown.** The model is aimed at EN/ZH, and PT aliases and summaries could be weak or drift into English.
- **Documentation inconsistencies** in the release: drafter layer counts differ, and one mention gives 159B parameters instead of 309B. [orcarouter] This does not directly affect API use.
- A reported tendency to fall into repetitive retry loops in agentic flows [tabbit]. This is less relevant to single-shot generation.
- **Privacy:** vault text goes to Xiaomi and DeepInfra via OpenRouter. This connects to the map's open privacy item.

## Implications for kb

- Keep the model slug in config (provider-neutral key, e.g. `[generation].model`). Default it to `xiaomi/mimo-v2.6-flash`, and document `xiaomi/mimo-v2.5` as a same-price fallback.
- Always send `reasoning: {enabled: false}` for alias, summary and rewrite calls. In a smoke or integration test, assert that `usage.reasoning_tokens` is 0 (or absent).
- Use `response_format` with a JSON schema (strict), and validate the output in Go anyway. On a schema failure, retry once and then skip; never write malformed frontmatter.
- Optionally pin the provider with `provider.order: ["xiaomi"]` and `allow_fallbacks: true` for uptime, and pin the canonical dated slug for reproducibility.
- Budget: about $0.0001 per call with reasoning off. A cost ceiling per run (an open map item) can be derived from the `usage` / `cost` fields in OpenRouter responses.
- Before committing, run a small PT/EN eval set: 10–20 concept titles for aliases (acronym, synonyms, PT↔EN) and 10 PT paragraphs for summaries. Check language fidelity (a PT summary stays in PT). Being unverified on PT is the main gap.
- Everything must stay optional: with no key, a provider error, or schema failure, kb skips generation and behaves as it does today.

## Live calls

Not run. `OPENROUTER_API_KEY` was not set in the research environment. To reproduce, use the two tests from the ticket brief (alias JSON for "Confidence-gated routing"; a 2-sentence PT summary) with `reasoning.enabled=false`. Record `usage`, the OpenRouter `cost`, and wall-clock latency.

## Sources

- OpenRouter models API: https://openrouter.ai/api/v1/models (filtered `mimo`)
- OpenRouter endpoints API: https://openrouter.ai/api/v1/models/xiaomi/mimo-v2.6-flash/endpoints
- OpenRouter model page: https://openrouter.ai/xiaomi/mimo-v2.6-flash
- OpenRouter reasoning docs: https://openrouter.ai/docs/guides/best-practices/reasoning-tokens
- Xiaomi deep-thinking docs: https://mimo.mi.com/docs/en-US/quick-start/usage-guide/text-generation/deep-thinking
- Xiaomi V2.6 release notes: https://mimo.mi.com/docs/en-US/updates/model
- HF model card: https://huggingface.co/XiaomiMiMo/MiMo-V2.6-Flash-RL
- Orca Router release analysis: https://www.orcarouter.ai/blog/xiaomi-mimo-v2-6-flash-release
- Kingy AI benchmarks: https://kingy.ai/blog/mimo-v2-6-pro-benchmarks-specs-comparison/
- Tabbit Pro vs Flash: https://go.tabbit.ai/blog/mimo-v2-6-pro-vs-mimo-v2-6-flash
- BenchLM: https://benchlm.ai/models/mimo-v2-6-flash
- Vercel AI Gateway: https://vercel.com/ai-gateway/models/mimo-v2.6-flash
- Hermes Agent thinking bug: https://github.com/NousResearch/hermes-agent/issues/27325
