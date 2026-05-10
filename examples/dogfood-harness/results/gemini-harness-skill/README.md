# Gemini Harness Skill Dogfood Result

This directory stores a sanitized example result from running `examples/dogfood-harness/prompts/gemini-harness-skill.txt` with `gemini -p`.

Captured with:

```text
gemini 0.41.2
```

This dogfood verifies two contracts:

- `agent-team-gemini-harness` can generate a tiny Gemini harness skeleton.
- The generated harness can drive daemonless `agent-team` state from run creation to `run_closed`.

The pinned `run_harness_skill_gemini` ID is a deterministic test fixture. Production user-facing harness runs should create or resolve `RUN_ID` and `TASK_ID` internally through the orchestrator and should not ask ordinary users to provide raw IDs.

Generated harness files in the temp workspace:

- `GEMINI.md`
- `.gemini/agents/writer.md`
- `.gemini/agents/reviewer.md`
- `.gemini/skills/manual-dogfood-orchestrator/SKILL.md`

Expected result files:

- `generated/GEMINI.md`
- `generated/writer-agent.md`
- `generated/reviewer-agent.md`
- `generated/orchestrator-skill.md`
- `writer.md`
- `review.md`
- `run-summary.json`
- `event-log.json`
- `final-response.md`

The JSON files intentionally omit timestamps and absolute temp paths while preserving run status, task counts, event types, and state versions.
