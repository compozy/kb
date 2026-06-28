## Memory Ledger — OKF skill sync

- **Goal:** Update the built-in `kb` skill (`skills/kb/`) to document the OKF dual-mode feature implemented under `.compozy/tasks/okf-dual-mode`. Success = skill accurately covers `mode: wiki|okf`, `kb topic new --mode okf`, `kb promote`, `kb okf check`, `[okf].types`, OKF bundle layout, and the producer-field/relative-link contract.
- **Constraints/Assumptions:** Doc-only change to `skills/kb/`. No code changes. Must match the SHIPPED behavior (verified against source, not just the techspec). Run `make verify` only if any Go is touched (it is not), but at least sanity-check.
- **Key decisions:**
  - Skill is at `skills/kb/SKILL.md` (the "built-in kb skill"), not `.agents/skills/...`.
  - Add OKF coverage to SKILL.md (frontmatter desc, intro, dedicated section, quick ref, dispatch, errors, constraints) + new `references/okf-mode.md` + OKF concept schema in `references/frontmatter-schemas.md`.
- **State:** DONE. Docs updated and verified against shipped behavior.
- **Done:** Verified impl against source; updated `skills/kb/SKILL.md` (desc, intro mode model, new OKF Dual-Mode section, quick-ref topic `--mode` + OKF block, dispatch rows, Procedure 6 distill loop, auto-log clarification, error rows, constraints); added `skills/kb/references/okf-mode.md`; updated `references/frontmatter-schemas.md` (topic.yaml `mode`, OKF concept schema, quick-ref rows); added OKF note to `references/architecture.md`. Confirmed `topic.yaml` writes `mode`.
- **Now:** —
- **Next:** Optional — user may want a commit (not requested).
- **Open questions:** none (UNCONFIRMED: whether skill is embedded into binary — appears to be a standalone repo skill dir; no embed.FS found referencing skills/kb).
- **Working set:** `skills/kb/SKILL.md`, `skills/kb/references/okf-mode.md` (new), `skills/kb/references/frontmatter-schemas.md`.

### Shipped OKF facts (ground truth)
- Topic mode: `mode: wiki|okf` in `topic.yaml` (empty→wiki). `kb topic new <slug> <title> <domain> --mode wiki|okf` (default wiki).
- OKF scaffold = FLAT: topic dir + `CLAUDE.md` (OKF template) + `topic.yaml` (`mode: okf`) + `AGENTS.md` symlink + `index.md` (frontmatter `okf_version: "0.1"`, body `# OKF Bundle Index`) + `log.md` (`# Directory Update Log` + Initialization entry under ISO date `## YYYY-MM-DD`). No raw/wiki/outputs/bases, no .gitkeep.
- `kb promote <wiki-doc> --to <okf-topic> --type <T> [--description <text>]` (both `--to` and `--type` REQUIRED). Mechanical, non-LLM, non-destructive. Writes JSON ConceptResult.
  - Errors if `--to` is not a `mode: okf` topic ("target topic must use mode okf").
  - Concept filename = `<slugify(base(source-path))>.md` at bundle root; collisions get `-2`, `-3`.
  - Frontmatter written: `description, tags(optional), timestamp(RFC3339 UTC), title, type` (alphabetical). title←source title or humanized key; description←--description else first body sentence else empty+warning; tags←source tags if present. Wiki stage markers dropped.
  - `[[wikilinks]]`→relative markdown links `[label](key.md[#anchor])`. Unresolved targets still emitted + recorded in `unresolvedLinks`.
  - Auto regenerates `index.md` (type-grouped bullets) and inserts newest-first `log.md` entry. Source untouched.
- `kb okf check <topic> [--strict] [--format table|json|tsv]`. Lenient per OKF §9. Columns: severity, kind, filePath, target, message. Non-zero exit on errors (and on warnings under --strict).
  - ERROR: concept missing/empty `type`, unparseable frontmatter, root non-`index.md` declaring okf_version, root index.md with non-`okf_version` keys, log `## ` heading not YYYY-MM-DD.
  - WARNING (→error w/ --strict): missing producer field title/description/timestamp; `type` outside `[okf].types` (only when vocab non-empty).
  - Excluded files (never need type): `index.md`, `log.md`, `CLAUDE.md`, `AGENTS.md`, symlinks, README/LICENSE/NOTICE/ATTRIBUTION, dotdirs.
- Config: `[okf].types = []` (empty ⇒ no vocab warnings). Also surfaced in `config.example.toml`.
