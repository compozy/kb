# Local markdown tracker

Issue tracker for wayfinder maps in this repo. Everything is plain markdown, versioned with git.

## Layout

```
.wayfinder/<map-slug>/
├── map.md            # the map issue (label wayfinder:map)
├── tickets/NN-slug.md  # child issues of the map
└── research/<name>.md  # assets produced while resolving tickets
```

## Issue frontmatter

```yaml
id: 3              # identity; unique inside the map
title: <name>      # always refer to the issue by this
labels: [wayfinder:research]
status: open       # open | closed
assignee: ""       # the claim; empty = unclaimed
blocked_by: [1, 2] # ids of tickets that must close first
```

## Wayfinding operations

- **Create map**: write `map.md` with `labels: [wayfinder:map]`.
- **Create ticket**: write `tickets/NN-slug.md`; the next free `NN` is its id. Wire `blocked_by` in a second pass.
- **Claim**: set `assignee` before any work.
- **Resolve**: append a `## Resolution` section (the resolution comment), set `status: closed`, then add a line to the map's Decisions so far.
- **Frontier**: open tickets whose `assignee` is empty and whose every `blocked_by` id is closed:
  `grep -l 'status: open' .wayfinder/<map>/tickets/*.md`, then check `assignee` and `blocked_by`.
- **Blocking** has no native rendering here; `blocked_by` in frontmatter is the convention.
