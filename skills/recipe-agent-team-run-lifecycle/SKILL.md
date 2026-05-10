---
name: recipe-agent-team-run-lifecycle
description: "Recipe: Create, execute, and close a daemonless agent-team workflow run."
metadata:
  version: 1.1.0
  openclaw:
    category: "recipe"
    domain: "agent-orchestration"
  requires:
    bins:
      - agent-team
    skills:
      - agent-team-run
      - agent-team-task
      - agent-team-sync
---

# Agent Team Run Lifecycle

Use this recipe when an orchestrator needs to run a complete daemonless workflow.

## Prerequisites

Load:

- `agent-team-shared`
- `agent-team-run`
- `agent-team-task`
- `agent-team-sync`

## Steps

1. Initialize state.

```bash
agent-team init
```

2. Create one run.

```bash
agent-team run create --title "workflow title"
```

3. Create assigned tasks.

```bash
agent-team task create --run RUN_ID --agent AGENT --title "task title" --body "task contract"
```

4. Workers start tasks after checking inbox and sync.

```bash
agent-team inbox list --agent AGENT --run RUN_ID --unread
agent-team sync check --agent AGENT --run RUN_ID --task TASK_ID
agent-team task start --task TASK_ID --agent AGENT
```

5. Workers complete or block tasks.

```bash
agent-team task complete --task TASK_ID --agent AGENT --evidence "verified" --artifact "_workspace/RUN_ID/TASK_ID.md"
agent-team task block --task TASK_ID --agent AGENT --reason "specific blocker"
```

6. Close the run only after all tasks are done.

```bash
agent-team run status --run RUN_ID
agent-team run close --run RUN_ID --reason "all tasks complete"
```

## Notes

- Do not close a run with unfinished tasks.
- Use `_workspace/{run_id}/` for artifacts.
- Use messages for contract changes instead of direct peer-to-peer updates.

