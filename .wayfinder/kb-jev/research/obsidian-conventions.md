# Obsidian conventions for aliases and typed relations

Research for ticket 04 (`tickets/04-obsidian-conventions.md`). Question: which frontmatter shapes do Obsidian, Dataview and Breadcrumbs read for **Aliases** and **Relations**, so kb output shows up in their graphs without kb calling any Obsidian API.

Checked 2026-09-22 against: obsidian-help `master` (latest release notes are v1.14.x), `obsidian.d.ts` from obsidian-api `master`, obsidian-dataview `master`, and Breadcrumbs `master` (v4.21.12, commit 90f48e4).

## Findings

### 1. `aliases`

- Put aliases in the `aliases` property as a YAML **list**. The help page says: "Aliases should always be formatted as a list in YAML." [S1]
- v1.9.0 removed the old singular `alias`, `tag` and `cssclass` keys, and those values *must* be lists. v1.9.3 relaxed this so a single string like `aliases: 'X'` is also valid, but a comma-separated string (`'a, b'`) is marked invalid. [S5, S6] **Always write a list.**
- How Obsidian uses an alias:
  - Aliases appear in `[[` link autocomplete. When you pick one, Obsidian writes `[[Artificial Intelligence|AI]]`, meaning the real file name plus the alias as display text. It does not write `[[AI]]`, "to ensure interoperability with other applications using the Wikilink format". [S1]
  - The Backlinks pane shows **unlinked mentions** of aliases, and converting one produces `[[Note|alias]]`. [S1, S3]
  - **UNVERIFIED (not stated in the official docs):** a bare `[[AI]]`, where `AI` is only an alias, does **not** resolve to the note. Obsidian resolves link targets by file name or path, so such a link stays unresolved and clicking it creates `AI.md`. Community reports agree, and the docs' choice of the `[[File|alias]]` format suggests the same. Treat it as true when designing kb.
- Dataview exposes aliases as `file.aliases`. [S9]
- Obsidian's autosuggest for the `aliases` property value is turned off (v1.4.3). This has no effect on kb.

### 2. Wikilinks inside properties (text and list)

- Supported since **v1.4.0** (2023-07-26), when Properties were introduced. Text and list property types "support internal links". [S4]
- **Quoting is required.** The Properties docs say "When using internal links in list properties, surround them with quotes." [S2] The Dataview docs agree. [S8] Unquoted `[[X]]` is valid YAML but means a nested flow sequence, not a string. Verified with kb's own parser: `related: [[X]]` parses to `[[ "X" ]]`, which is a list inside a list.
  ```yaml
  link: "[[Link]]"
  linklist:
    - "[[Link]]"
    - "[[Link2]]"
  ```
- **Backlinks:** since **v1.4.4** (2023-08-23), release notes say "Properties with links will now properly show in backlink entries." [S7] The same release fixed heading links (`[[X#H]]`) in list properties showing as broken.
- **Graph view:** property links appear as edges. The graph-view docs only say "lines represent internal links" [S10] and no release note names the graph. The evidence is indirect:
  - The API field `CachedMetadata.frontmatterLinks` exists `@since 1.4.0`. [S11]
  - A 2026 forum feature request asks for an option to *hide* frontmatter links from the graph, because they currently show there and add noise to the local graph. [S12]
  - Mark "since 1.4.0" as **UNVERIFIED** for the graph specifically. That it works today is well supported.
- The graph view **does not type edges**. A `related` link and a body link look the same, with no labels or per-relation colours. Only plugins (Breadcrumbs, Dataview) or Bases can show which kind of Relation an edge is.
- There is **no toggle** to exclude property links from the graph (see the S12 feature request). Every Relation kb writes adds a graph edge.
- Other behaviour:
  - Obsidian updates property links when the target file is renamed (fixed in v1.5.7).
  - Markdown-style links `[t](x.md)` inside properties have been supported since **v1.11** (updated on rename; file names with spaces fixed in v1.11.3). This matters only for OKF mode, which uses Markdown links.
- API shape: each list entry is indexed separately in `frontmatterLinks` with key `field.N` (for example `related.0`). Breadcrumbs relies on this: it runs `key.split(".")[0]`. [S13]

### 3. Dataview

- A quoted `"[[X]]"` string in frontmatter is parsed as a **Link** value, not text. [S8]
- `file.outlinks` includes frontmatter links. In `src/data-import/markdown-file.ts`, the loop over `metadata.frontmatterLinks` pushes them into the page's links next to body links. [S14] So `file.inlinks`, `FROM [[X]]` and `outgoing([[X]])` all count Relations.
- Typed queries work on the property itself. Examples: `WHERE contains(related, [[X]])`, `TABLE extends, prerequisite FROM "topic/wiki"`. Links are compared by resolved path. **UNVERIFIED:** exact behaviour of `contains` with unresolved links.

### 4. Breadcrumbs (v4)

- The `typed_link` builder is the main edge source. For every file it goes through Obsidian's `frontmatterLinks`, takes the field name before the `.`, and keeps the link only if that name **exactly matches** a configured edge-field label (set membership, case-sensitive). It resolves targets with Obsidian's resolver. Unresolved targets become unresolved nodes. [S13] Body inline fields `field:: [[X]]` are also read, without the Dataview plugin.
- **Default edge fields:** `up`, `down`, `same`, `next`, `prev`. Default groups are `ups`, `downs`, `sames`, `nexts`, `prevs`. [S15]
- Custom names (`related`, `extends`, `prerequisite`, `contradicts`, `affects`) work only after the user adds them as edge fields in plugin settings (stored in `.obsidian/plugins/breadcrumbs/data.json`). Breadcrumbs does not discover them.
- **Implied relations** are per-field transitive rules: `chain: [{field}]`, `close_field`, `close_reversed`, `rounds`. The defaults are `up`↔`down` reversed, `same`↔`same`, and `next`↔`prev` reversed. `self_is_sibling: ["same"]`. [S15] v4 supports any chain, for example `[up, up] -> up` or `[extends] -> reversed extended_by`.
- v4 changes: nodes are keyed by full path, not basename (fixes basename collisions), and implied relations can be set per hierarchy. [S16]
- Nested objects break Breadcrumbs. With `relations: {extends: [...]}`, the frontmatter key becomes `relations.extends.0`, Breadcrumbs keeps only `relations`, and no field matches. **UNVERIFIED:** whether Obsidian indexes links nested inside objects at all. The property editor does not support nested objects.

### 5. What kb does today

- **Writing.** `frontmatter.Generate` uses `gopkg.in/yaml.v3` (`internal/frontmatter/frontmatter.go`) and quotes wikilink strings automatically, with single quotes:
  ```yaml
  related:
      - '[[Artificial Intelligence]]'
  up: '[[Parent]]'
  ```
  Single quotes are valid YAML and Obsidian accepts them. Aliases with `:` also get quoted. The schema doc `skills/kb/references/frontmatter-schemas.md` shows double-quoted `"[[...]]"` for `sources` and `informed_by`. Both forms are fine.
- **Parsing.** `normalizeValue` turns all-string lists into `[]string`. An unquoted `[[X]]` becomes a nested list. `GetStringSlice` then returns nil, and lint reports `frontmatter field "sources" must be a list of strings`, but only for fields listed in that path's `schemaSpec.listFields` (`internal/lint/lint.go` `schemaForPath`). For wiki concepts that is `tags` and `sources`.
- **kb lint does *not* scan frontmatter wikilinks into the link graph.**
  - `loadVault` builds `file.links` from `extractDocumentWikilinks(body, values)`, which reads the body only. `values` is used just to pick the extraction mode by `source_kind`.
  - `wikilinkPattern` runs on the body after `markdownBody` strips the leading frontmatter.
  - The only frontmatter links lint checks are `sources`, in `findSourceIssues`: they must resolve to a raw file (`resolveTarget(ref, true)`), and staleness is compared with `scraped`.
  - `informed_by` is checked for shape but not resolved.
  - Result: a Relation pointing to a missing article is **not** reported as a dead link, and Relations do **not** count as incoming links for orphan detection. Obsidian does count them, so an article kb calls an orphan may not be one in Obsidian's graph, and the reverse.
- **Link resolution is looser than Obsidian's.** `resolveTarget` tries, in order: vault or topic path, file stem, **`title`**, a **canonical token key** (lowercased, stopwords removed), then a **canonical prefix match** (for example the part of the title before `:`). Obsidian resolves only by path or file name. So kb lint accepts `[[Some Title]]` links that Obsidian shows as unresolved (the same gap as bare aliases in §1).
- **The `aliases` property is ignored.** `linkAliasesForFile` builds its alias index from paths, stem, `title` and the part of `title` before the `:`. It never reads frontmatter `aliases`.

## Recommended YAML shape

For wiki articles, in wiki mode:

```yaml
---
title: Artificial Intelligence
aliases:              # always a list, even with one entry; plain strings, no [[ ]]
  - AI
  - IA
  - Inteligência Artificial
related:              # symmetric
  - "[[machine-learning]]"
extends:              # this article builds on the target
  - "[[neural-networks]]"
prerequisite:         # read the target first
  - "[[linear-algebra]]"
contradicts:
  - "[[symbolic-ai-supremacy|Symbolic AI supremacy]]"
sources:
  - "[[raw/articles/some-source]]"
---
```

Sources use the `affects` Relation:

```yaml
affects:
  - "[[artificial-intelligence]]"
```

Rules:

1. **Top-level flat keys, one per Relation type, each a list of quoted wikilinks.** Do not use nested maps (`relations: {...}`). They break Breadcrumbs' `field.N` parsing and the Obsidian property editor.
2. **Lowercase snake_case field names** that match what the user configures in Breadcrumbs, since matching is exact. Do not reuse Breadcrumbs' default `up`/`down`/`same`/`next`/`prev` unless kb means that hierarchy.
3. **The link target is the file stem or vault-relative path, not the title.** Obsidian resolves only paths and file names. Use `|Display` for readable text. Use a vault-relative path when the stem is not unique across the vault; kb vaults hold several topics, so stems can collide. **UNVERIFIED assumption:** the kb vault root is the Obsidian vault root.
4. **Insert body Links with aliases as `[[target|alias]]`, never `[[alias]]`.** This matches Obsidian's own behaviour and keeps links resolved in Obsidian.
5. **Let `frontmatter.Generate` do the quoting.** Never write YAML by hand; CLAUDE.md already requires this.
6. **Breadcrumbs implied rules, if documented for users:**
   - `related`: chain `[related]` → close `related`, reversed (symmetric).
   - `extends`: → `extended_by`, reversed.
   - `prerequisite`: → `enables`, reversed.
   - `contradicts`: symmetric.
   - `affects`: → `affected_by`, reversed.

## Implications for kb

- Relations written as quoted wikilink lists appear in Obsidian **backlinks** (since v1.4.4), the **graph view** (untyped edges), **Dataview** `file.inlinks`/`outlinks` and typed property queries, and **Breadcrumbs** once the user registers the field names. No Obsidian API integration is needed, which fits the map's "Obsidian only renders it" rule.
- **Lint and Obsidian disagree** on two points, and the spec must decide how to handle them:
  - Obsidian counts frontmatter links as links; kb lint does not. The spec should decide whether lint adds Relation fields to the link graph (dead-link and orphan checks) and adds them to the wiki `listFields` so shape errors are caught.
  - kb resolves by title, canonical key and prefix; Obsidian does not. A link kb accepts can be unresolved in Obsidian. The spec could add a lint warning such as "resolves in kb but not in Obsidian".
- **Alias resolution is a design choice.** If kb adds frontmatter `aliases` to `aliasIndex`, a bare `[[alias]]` would pass kb lint but stay unresolved in Obsidian. The safer contract: aliases feed Candidate generation (mention matching), and kb always writes `[[stem|alias]]`.
- **Breadcrumbs needs setup.** It does nothing with `related`/`extends`/... until they are configured. kb could document a `data.json` snippet or template, but should not write `.obsidian/` (out of scope under the map's no-Obsidian-API rule). Dataview and backlinks need no setup.
- **Every Relation adds graph noise**, with no way for users to hide property links from the graph. This argues for writing only Relations in the apply band and keeping review-band Candidates out of frontmatter (in the `.kb/` index or quarantine).
- **Keep `aliases` a plain string list**, with no wikilinks, no comma-joined strings, and never the singular `alias` key (dropped in v1.9.0).
- **OKF mode is separate.** OKF concept bodies use relative Markdown links. Markdown links in properties work only since Obsidian v1.11, so OKF Relations, if any, need their own decision. Do not reuse the wiki-mode shape.

## Sources

- [S1] Obsidian Help, Aliases: https://help.obsidian.md/aliases (repo `en/Linking notes and files/Aliases.md`)
- [S2] Obsidian Help, Properties (List / Text format): https://help.obsidian.md/properties
- [S3] Obsidian Help, Backlinks: https://help.obsidian.md/plugins/backlinks
- [S4] Obsidian v1.4.0 release notes (2023-07-26), "About properties > Internal links", deprecation of `alias`: `obsidian-help/Release notes/v1.4.0.md`
- [S5] Obsidian v1.9.0 release notes: removal of `tag`/`alias`/`cssclass`, values must be lists
- [S6] Obsidian v1.9.3 release notes: single strings accepted for `aliases`/`tags`/`cssclasses`
- [S7] Obsidian v1.4.4 release notes (2023-08-23): "Backlinks: Properties with links will now properly show in backlink entries"; also v1.5.7 (rename updates property links), v1.11 (Markdown links in properties)
- [S8] Dataview docs, Types of metadata (quoted links in YAML): https://blacksmithgu.github.io/obsidian-dataview/annotation/types-of-metadata/ (via ctx7 `/blacksmithgu/obsidian-dataview`)
- [S9] Dataview docs, Metadata on pages (`file.inlinks`, `file.outlinks`, `file.aliases`): https://blacksmithgu.github.io/obsidian-dataview/annotation/metadata-pages/
- [S10] Obsidian Help, Graph view: https://help.obsidian.md/plugins/graph
- [S11] `obsidian.d.ts`: `CachedMetadata.frontmatterLinks?: FrontmatterLinkCache[]` `@since 1.4.0`; `MetadataCache.resolvedLinks`: https://github.com/obsidianmd/obsidian-api
- [S12] Obsidian Forum, "An option to disable frontmatter links to show up in Graph View": https://forum.obsidian.md/t/an-option-to-disable-frontmatter-links-to-show-up-in-graph-view/114940
- [S13] Breadcrumbs `src/graph/builders/explicit/typed_link.ts`: https://github.com/SkepticMystic/breadcrumbs
- [S14] Dataview `src/data-import/markdown-file.ts` (frontmatterLinks → page links): https://github.com/blacksmithgu/obsidian-dataview
- [S15] Breadcrumbs `src/const/settings.ts` (`DEFAULT_SETTINGS.edge_fields`, `implied_relations.transitive`)
- [S16] Breadcrumbs `V4.md` (typed frontmatter links, full-path nodes, per-hierarchy implied relations)
- kb code: `/Users/pedronauck/Dev/compozy/kb/internal/frontmatter/frontmatter.go`, `/Users/pedronauck/Dev/compozy/kb/internal/lint/lint.go` (`wikilinkPattern`, `extractDocumentWikilinks`, `findSourceIssues`, `schemaForPath`, `resolveTarget`, `linkAliasesForFile`), `/Users/pedronauck/Dev/compozy/kb/skills/kb/references/frontmatter-schemas.md`
