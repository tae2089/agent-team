## Harness: Manual Dogfood

**Goal:** Exercise a minimal Agent Team Codex producer-reviewer workflow.

**Trigger:** Use `.agents/skills/manual-dogfood-orchestrator/SKILL.md` for manual dogfood harness runs, reruns, retries, updates, audits, and review follow-ups. Simple questions may be answered directly.

**Model:** Follows the `revfactory/harness` orchestrator/specialist structure adapted to Codex-native skills, agents, artifacts, and Agent Team runtime.

**Orchestrator:** `.agents/skills/manual-dogfood-orchestrator/SKILL.md`
**Agents:** `.codex/agents/`
**Artifacts:** `_workspace/{run_id}/`

**Runtime State:**

- Load `agent-team-shared` first for global runtime rules.
- Use recipe skills for workflow shape: `recipe-agent-team-run-lifecycle` for full runs, `recipe-agent-team-worker-checkpoint` for worker checkpoints, and `recipe-agent-team-operational-audit` for audit/status/cleanup.
- Use service skills for navigation: `agent-team-run`, `agent-team-task`, `agent-team-inbox`, `agent-team-sync`, and `agent-team-ops`.
- Use exact command helper skills for command syntax and flags, for example `agent-team-task-complete`, `agent-team-sync-check`, `agent-team-message-send`, or `agent-team-event-log`.
- `RUN_ID` and `TASK_ID` are orchestrator-owned internal context, not required user input.
- If the user provides an advanced/debug `RUN_ID` or `TASK_ID`, inspect and resume that run/task before creating new state.
- If no runtime context is available, the orchestrator loads `agent-team-run-create` and `agent-team-task-create`, then creates one generated-ID run and generated-ID task records through those helper contracts.
- Do not use runtime state during harness setup, editing, audit-only work, simple one-shot answers, or explicitly local-only runs.
- Orchestrator owns run creation, task creation, evidence aggregation, inbox/sync checks, and artifact integration.
- Workers update only their assigned task.
- Completed tasks require evidence and an artifact path.
- Blocked tasks require a concrete blocked reason.
- `_workspace/` is for artifacts and reports only.

**Change History:**
| Date | Change | Target | Reason |
| --- | --- | --- | --- |
| 2026-05-10 | Initial tiny harness | all | Repository-local skill verification |
