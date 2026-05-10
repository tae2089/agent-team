# Codex Harness Skill Dogfood Result

This directory stores a sanitized example result from running `examples/dogfood-harness/prompts/codex-harness-skill.txt` with `codex exec`.

Captured with:

```text
codex-cli 0.130.0
```

This dogfood verifies two contracts:

- `agent-team-codex-harness` can generate a tiny Codex harness skeleton.
- The generated harness can drive daemonless `agent-team` state from run creation to `run_closed`.

The pinned `run_harness_skill_codex` ID is a deterministic test fixture. Production user-facing harness runs should create or resolve `RUN_ID` and `TASK_ID` internally through the orchestrator and should not ask ordinary users to provide raw IDs.

Generated harness files in the temp workspace:

- `AGENTS.md`
- `.codex/agents/writer.toml`
- `.codex/agents/reviewer.toml`
- `.agents/skills/manual-dogfood-orchestrator/SKILL.md`

Expected result files:

- `generated/AGENTS.md`
- `generated/writer-agent.toml`
- `generated/reviewer-agent.toml`
- `generated/orchestrator-skill.md`
- `writer.md`
- `review.md`
- `run-summary.json`
- `event-log.json`
- `final-response.md`

The JSON files intentionally omit timestamps and absolute temp paths while preserving run status, task counts, event types, and state versions.
