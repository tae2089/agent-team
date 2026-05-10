Done.

Generated harness files:
- `AGENTS.md`
- `.codex/agents/writer.toml`
- `.codex/agents/reviewer.toml`
- `.agents/skills/manual-dogfood-orchestrator/SKILL.md`

Runtime result:
- Run id: `run_harness_skill_codex`
- Final run status: `closed`
- Pre-close summary reported `close_ready:true`
- Event log inspected through `run_closed`
- Reviewer task dependency confirmed: `reviewer_task` depends on `writer_task`

Artifacts:
- Writer: `_workspace/run_harness_skill_codex/writer.md`
- Reviewer: `_workspace/run_harness_skill_codex/review.md`

The generated harness used `agent-team` state correctly: local state was initialized, the run and tasks were recorded, the writer inbox notice was listed and acked, sync checks ran, tasks were started and completed with evidence plus artifact paths, and the run was closed only after readiness was confirmed. No git commit was created.
