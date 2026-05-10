# Dogfood Harness

**Goal:** Verify the `agent-team-gemini-harness` skill by running a writer-reviewer flow.

**Trigger:** Use `.gemini/skills/manual-dogfood-orchestrator/SKILL.md` for dogfood work.

**Model:** Coordinated specialist system using Gemini agents and Agent Team runtime.

**Orchestrator:** `.gemini/skills/manual-dogfood-orchestrator/SKILL.md`
**Agents:** `.gemini/agents/`
**Artifacts:** `_workspace/run_harness_skill_gemini/`

**Runtime State:**
- Activate `agent-team-shared` first.
- Use `agent-team-run`, `agent-team-task`, `agent-team-inbox`, `agent-team-sync`, and `agent-team-ops` for navigation before command helper skills.
- Activate exact command helper skills such as `agent-team-run-create`, `agent-team-task-create`, `agent-team-task-complete`, and `agent-team-run-close` before using their command forms.
- `RUN_ID` and `TASK_ID` are orchestrator-owned internal context, not required user input.
- Orchestrator owns run/task creation and run closure.
- Workers (writer, reviewer) update their assigned tasks.

**Reload:** Run `/agents reload` and `/skills reload` inside Gemini CLI.

**Change History:**
| Date | Change | Target | Reason |
| --- | --- | --- | --- |
| 2026-05-10 | Initial dogfood harness | all | Skill verification |
