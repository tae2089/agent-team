# Dogfood Harness

This directory has two dogfood layers.

## CI Dogfood

`run.sh` is deterministic and safe for CI. It exercises the real `agent-team` CLI contract without spawning Codex, Gemini, or any external LLM process.

The workflow shape mirrors a harness runtime:

1. an orchestrator creates a run and role-specific tasks
2. a planner writes an implementation plan artifact
3. a worker checks sync, acknowledges a contract-change message, and writes a result artifact
4. a reviewer depends on the worker task and writes a review artifact
5. the orchestrator inspects summary and event history, then closes the run

Run from the repository root:

```bash
sh examples/dogfood-harness/run.sh
```

The script uses a temporary `AGENT_TEAM_STATE_DIR` and writes artifacts under a temporary directory.

## Manual Harness Dogfood

Manual dogfood checks that the Codex/Gemini harness skills use the runtime contract correctly in an actual assistant workflow:

- [manual-codex.md](manual-codex.md)
- [manual-gemini.md](manual-gemini.md)

These guides intentionally stay out of CI because they require an interactive assistant CLI and may involve LLM execution.

Reusable headless prompts live under `prompts/`. Sanitized example outputs live under `results/` when a manual run completes successfully:

```text
examples/dogfood-harness/prompts/codex.txt
examples/dogfood-harness/prompts/gemini.txt
examples/dogfood-harness/prompts/codex-harness-skill.txt
examples/dogfood-harness/prompts/gemini-harness-skill.txt
examples/dogfood-harness/results/codex/
examples/dogfood-harness/results/gemini/
examples/dogfood-harness/results/codex-harness-skill/
examples/dogfood-harness/results/gemini-harness-skill/
```

The committed result examples are intentionally sanitized: timestamps and absolute temporary paths are removed, while event ordering, task counts, run status, and artifact paths are preserved.

The `*-harness-skill.txt` prompts verify the harness skills themselves: they require the assistant to read the repository-local harness skill, generate a tiny skill-first harness skeleton, and then prove that the generated harness can drive `agent-team` runtime state to a closed run through runtime skill contracts.
