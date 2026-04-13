# kb Release NPM Design

## Goal

Add npm publication for the `kb` CLI under `@compozy/kb` using the existing GoReleaser-based release pipeline, and update installation documentation to show Homebrew, npm, and Go install paths clearly.

## Decisions

- Keep the current split GoReleaser release flow in `kb`.
- Let GoReleaser publish the npm package through `npms`; do not add a separate `npm publish` script step.
- Provide npm authentication to the release job with `NPM_TOKEN`/`NODE_AUTH_TOKEN`.
- Keep the package lightweight: the npm package is only a delivery wrapper for the released CLI binaries.

## Files

- `.goreleaser.yml`: add `npms` metadata and include npm install instructions in the release body.
- `.github/actions/setup-node/action.yml`: reusable Node/npm setup with registry scope configuration.
- `.github/workflows/release.yml`: configure Node/npm auth for the production release job.
- `README.md`: document installation via Homebrew, npm, Go, and source build.

## Verification

- Validate workflow syntax with `actionlint`.
- Run `make verify`.
