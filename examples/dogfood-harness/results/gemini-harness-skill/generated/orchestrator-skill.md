---
name: manual-dogfood-orchestrator
description: "Orchestrator for the dogfood Agent Team Gemini harness. Use for testing the harness-skill-gemini flow."
---

# Manual Dogfood Orchestrator

This skill coordinates the writer-reviewer dogfood flow using Agent Team runtime skills backed by the daemonless `agent-team` CLI.

## Workflow

1. **Run Creation**: Resolve internal runtime context from active session state, an advanced/debug ID, recent open runs, or previous artifacts; if none matches, activate `recipe-agent-team-run-lifecycle` and `agent-team-run-create`, then create a generated-ID run and capture `.data.run.id`.
2. **Task Creation**:
   - Create a generated-ID task for `writer` and capture `.data.task.id`.
   - Create a generated-ID task for `reviewer` that depends on the writer task.
3. **Notice**: Send a compact notice message to `writer`.
4. **Integration**:
   - Wait for `writer` completion.
   - Wait for `reviewer` completion.
5. **Finalization**: Activate `agent-team-run-summary` and `agent-team-run-close`, then check summary and close the run.

## Data Protocol

- RUN_ID and TASK_ID are orchestrator-owned internal context, not required user input.
- This dogfood fixture used `run_harness_skill_gemini` only for deterministic verification.
- Artifacts: `_workspace/{run_id}/`
- Tasks use generated IDs by default.

## Rules

- Use `agent-team` commands only after activating the relevant Agent Team runtime skill.
- Activate `agent-team-shared`, `agent-team-run`, `agent-team-task`, `agent-team-message`, and `agent-team-sync`.
- Activate exact helper skills for commands.
