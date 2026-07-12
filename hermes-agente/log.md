# Hermes Agent — Log

Chronological, append-only record of every knowledge-base operation on this topic. Never rewrite history.

**Entry format:**

```markdown
## [2026-07-12] <op> | <short description>

(optional) one short paragraph of context, findings, or decisions
```

**Ops:** `bootstrap` · `ingest` · `compile` · `query` · `promote` · `split` · `lint`

**Quick queries:**

```bash
grep "^## \[" log.md | tail -10                # last 10 events
grep "^## \[.*compile" log.md | wc -l          # total compiles
grep "^## \[YYYY-MM" log.md                    # events in a month
```

---

## [2026-07-12] bootstrap | topic scaffolded

Topic `hermes-agente` created via `new-topic.sh`. Domain: `ai-agents`. Ready for ingest.

## [2026-07-12] ingest | 7-passos-obrigat-rios-para-usar-o-hermes-e-criar-agentes-que-trabalham-por-voc-guia-completo.md (youtube-transcript)

Ingested `7 passos obrigatórios para usar o Hermes e criar agentes que trabalham por você (guia completo)` into `hermes-agente/raw/youtube/7-passos-obrigat-rios-para-usar-o-hermes-e-criar-agentes-que-trabalham-por-voc-guia-completo.md`.
