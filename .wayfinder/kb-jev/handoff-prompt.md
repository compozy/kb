# kb-jev implementation handoff prompt

Full task contract for the /goal driving this implementation. Re-read this file after every compaction and before each ship step.

Goal: Implement the full kb-jev spec (~/Dev/compozy/kb/.wayfinder/kb-jev/spec.md, revision 3) in the kb repo, get a final review from Codex GPT 6 Astra xhigh via herdr, open a detailed PR, wait for CI, squash-merge, then wait for the release PR and squash-merge it so a new kb version ships.

Context
- Repo: ~/Dev/compozy/kb (Go CLI). Read AGENTS.md first; `make verify` (fmt -> lint -> test -> build -> boundaries) is the non-negotiable gate and `make lint` must report zero findings. Also run `make test-integration`.
- Spec: .wayfinder/kb-jev/spec.md is the source of truth, plus map.md, tickets/ and research/ in the same folder. Revision 3 is the latest: frontmatter keys have NO `ai_` prefix (triage, triage_reason, genre, depth, relevance, concepts, summary, criterion, entities, questions, quality, related/extends/prerequisite/example_of/contradicts, affects, supersedes, aliases, locked); body hash / contract / bank versions live in <topic>/.decisions/state.jsonl; collisions are handled by the owned-key list plus value-hash check (§6). Do not reintroduce `ai_` anywhere.
- Jev/typesafe skills: ~/Dev/courses/pedronauck/.agents/skills/jev-engineering and ~/Dev/courses/pedronauck/.agents/skills/typesafe-ai. Load them before writing question banks or the decisions client.
- Real vault for validation: ~/Dev/courses/pedronauck/research (kb.toml lives there). OPENROUTER_API_KEY is in the kb .env / environment. Keep real-run spend small (use --budget) and never write to that vault outside a dedicated test topic or a copy; prefer running real flows on a copy of one small topic in a temp dir.

Scope
- Implement everything in the spec, following its implementation order (§ "order of implementation"): internal/decisions + config, internal/questions banks, corpus + byte-preserving writer + state.jsonl, gates on ingest (with quarantine + inbound-reference ledger), classify + criterion + topic vocabulary, link, find, review/calibrate, lint additions, and the skill/docs updates the spec lists (skills/kb/references/*, README, AGENTS.md CLI tables).
- Tests exactly as the spec's test strategy describes (unit table-driven + integration with fake decisions/generation/Firecrawl servers). Follow the repo test conventions; fix production code when a test exposes a bug.
- Final validation must be a real run of the new commands against a small real topic copy with the real Jev/OpenRouter API, not only tests. Report cost.

Git hygiene
- The working tree already has unrelated changes I made (.agents/skills/* deletions/modifications from `npx skills add`, .peer-reviews/ deletions, README.md edit). Do NOT stage, revert, or touch them. Work on a new branch (e.g. feat/kb-jev) and only stage files you changed. Commit the .wayfinder/ spec folder together with the implementation. Never use git restore/checkout/reset/clean/stash.
- Conventional commits (feat:, fix:, docs:, ...).

Execution
- Use subagents on Opus (model opus, not Sonnet/Haiku) for parallel exploration and for independent implementation slices with non-overlapping file claims; you own integration and verification. Keep a ledger per your CLAUDE.md rules because this is long.

Final review with Codex via herdr (use the herdr and herdr-orchestration skills)
- When implementation is complete and `make verify` + integration tests pass, split a new pane next to yours with --no-focus and start Codex there with herdr agent start (kind codex), cwd ~/Dev/compozy/kb, args: --yolo -m gpt-6-astra -c model_reasoning_effort="xhigh".
- Prompt it for a thorough final review of the full branch diff against main and against the spec (correctness, spec conformance, data-loss risks in the frontmatter writer/quarantine, concurrency, error handling, tests, docs), reporting only actionable findings with file:line.
- Wait on it with herdr agent wait, read its report, fix every valid finding (explain any you reject), re-run make verify, and close the Codex pane when done.

Ship
1. Push the branch and open a PR to main with gh. The description must be detailed and extensive: motivation, summary of the spec, every new command and flag, frontmatter keys and state files, architecture/package map, gates and thresholds, cost model, migration for existing vaults, test coverage, real-run results with cost, Codex review findings and how each was handled, and known limitations.
2. Wait for all CI checks to pass (gh pr checks --watch). If something fails, fix, push, and wait again.
3. Squash-merge the PR (gh pr merge --squash --delete-branch).
4. The push to main triggers the release workflow, which creates/updates a "ci(release): Release vX.Y.Z" PR (there is already an open release PR #28 for v0.0.11; it should be updated). Wait for it to be updated and for its checks to pass, review that the version bump and changelog are correct, then squash-merge it too.
5. If anything is missing or broken after merge (CI on main, release workflow, release PR checks), fix it and push directly to main, then repeat until the release is published. Confirm the GitHub release/tag exists at the end.

Finish with a short report: PR URL, release PR URL, released version/tag, what was implemented, real-run results and cost, Codex findings handled, anything left undone.
