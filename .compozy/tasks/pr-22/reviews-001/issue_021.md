---
provider: coderabbit
pr: "22"
round: 1
round_created_at: 2026-06-27T02:47:23.093614Z
status: resolved
file: internal/youtube/channel_test.go
line: 39
severity: nitpick
author: coderabbitai[bot]
provider_ref: review:4583686378,nitpick_hash:93b3f348125c
review_hash: 93b3f348125c
source_review_id: "4583686378"
source_review_submitted_at: "2026-06-27T01:56:04Z"
---

# Issue 021: Assert the invalid-URL contract, not just any error.
## Review Comment

The `wantErr` branch only checks `err != nil`, so this still passes if `NormalizeChannelURL` starts returning the wrong error kind or message. Please assert `*Error`/`ErrorKindInvalidURL` and the expected diagnostic for the rejected input. As per path instructions, `MUST have specific error assertions (ErrorContains, ErrorAs)`.

<!-- cr-comment:v1:5f27f16c7cd403e330e41aa3 -->

_Source: Path instructions_

---

## Triage

- Decision: `VALID`
- Notes: `TestNormalizeChannelURL` currently treats all rejected channel inputs as equivalent by checking only `err != nil`. That would allow a regression where rejected channel URLs return a generic error, a different `ErrorKind`, or an unhelpful diagnostic while the test still passes. The owning invariant is the YouTube channel URL normalization contract: invalid or single-video inputs must fail with the package `*Error`, `ErrorKindInvalidURL`, and the expected diagnostic for the rejected input. The canonical suite is `internal/youtube/channel_test.go` because it already owns the valid/invalid `NormalizeChannelURL` table.
- Fix approach: Extend the invalid cases in `TestNormalizeChannelURL` with expected error message substrings, assert `errors.As` into `*Error`, assert `ErrorKindInvalidURL`, and assert the diagnostic contains the case-specific text.
- Resolution: Updated `TestNormalizeChannelURL` so invalid cases assert `errors.As` into `*Error`, `ErrorKindInvalidURL`, and case-specific diagnostic substrings.
- Verification: `rtk go test ./internal/youtube` exited 0 with `Go test: 37 passed in 1 packages`; `rtk make verify` exited 0 with `0 issues`, `DONE 829 tests, 1 skipped`, and `OK: all package boundaries respected`.
