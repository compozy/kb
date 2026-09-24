---
id: 4
title: Obsidian conventions for aliases and typed relations
labels: [wayfinder:research]
status: closed
assignee: "claude"
blocked_by: []
---

## Question

Which frontmatter shapes does Obsidian (and Dataview / Breadcrumbs) understand for aliases and typed relations, so kb's output lights up the graph without integration?

Pin down: `aliases` behavior in link resolution; whether wikilinks inside list properties count in backlinks and graph view; Breadcrumbs' expected field names for hierarchy/relations; how Dataview reads link lists. Record findings in `research/obsidian-conventions.md`.

## Resolution

- `aliases` must be a YAML list of plain strings; the singular `alias` key was removed in v1.9.0. Aliases feed autocomplete and unlinked mentions. A bare `[[alias]]` does not resolve (UNVERIFIED in docs), so write `[[stem|alias]]`.
- Quoted wikilinks in text and list properties (`related: ["[[x]]"]`) are real links. They count in backlinks since v1.4.4 and appear as untyped edges in the graph view (frontmatterLinks exists since 1.4.0). There is no toggle to hide them.
- Dataview puts frontmatterLinks into `file.outlinks`/`inlinks` and parses `"[[X]]"` as a Link. Breadcrumbs v4 reads `frontmatterLinks` whose key before `.N` exactly matches a configured edge field (defaults are `up`/`down`/`same`/`next`/`prev`). Custom Relation names and implied rules have to be configured by the user.
- Recommended shape: flat top-level list keys per Relation (`related`, `extends`, `prerequisite`, `contradicts`, `affects`), each entry a quoted wikilink to the file stem or path. No nested maps. Let `frontmatter.Generate` quote the values, which it already does.
- kb gaps: lint does not read frontmatter wikilinks into its link graph (only `sources` is checked, against raw files). It resolves by title, canonical key and prefix, which Obsidian does not. It ignores the `aliases` property.
- Details and sources: [../research/obsidian-conventions.md](../research/obsidian-conventions.md)
