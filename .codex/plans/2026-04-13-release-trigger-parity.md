# Restore `kb` Release Trigger Parity

## Summary

- Remove the `paths-ignore` filter from `on.push` in `.github/workflows/release.yml`.
- Keep the current split-release architecture in `kb`; it is not the cause of the missed production release.
- Match the behavior of the active `Release` workflow in `compozy`, where release-merge commits on `main` still trigger the release workflow.

## Key Changes

- Delete the `paths-ignore` block under `on.push.branches: [main]` in `.github/workflows/release.yml`.
- Preserve existing job-level conditions:
  - `release-pr` still skips `ci(release):` commits.
  - `prepare-release-tag`, `release-split`, and `release` still run only for `ci(release):` merge commits on `main`.
- Leave `.goreleaser.yml` and release helper actions unchanged.

## Test Plan

- Validate the workflow file with `actionlint`.
- Run `make verify` as the repository completion gate.
- After merge of a release PR that only changes `CHANGELOG.md`, confirm a `Release` push run is created for the merge SHA and that tag/release jobs execute.

## Assumptions

- `kb` remains a Go-only release flow without Docker or NPM publishing.
- Release secrets are already available via org-level or environment-level configuration.
- The required parity is behavior after release-PR merge, not a byte-for-byte copy of `compozy` tooling.
