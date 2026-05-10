---
name: recipe-agent-team-worker-checkpoint
description: "Recipe: Worker checkpoint flow for inbox, sync, start, complete, or block decisions."
metadata:
  version: 1.1.0
  openclaw:
    category: "recipe"
    domain: "agent-orchestration"
  requires:
    bins:
      - agent-team
    skills:
      - agent-team-inbox
      - agent-team-sync
      - agent-team-task
---

# Agent Team Worker Checkpoint

Use this recipe when a worker is about to start, checkpoint, or complete assigned work.

## Prerequisites

Load:

- `agent-team-shared`
- `agent-team-inbox`
- `agent-team-sync`
- `agent-team-task`

## Steps

1. Read unread messages for the current run.

```bash
agent-team inbox list --agent AGENT --run RUN_ID --unread
```

2. Ack messages only after handling them.

```bash
agent-team inbox ack --msg MSG_ID --agent AGENT
```

3. Run sync check for the task.

```bash
agent-team sync check --agent AGENT --run RUN_ID --task TASK_ID
```

4. If sync is clear, start or continue work.

```bash
agent-team task start --task TASK_ID --agent AGENT
```

5. Complete with evidence and an artifact path, or block with a concrete reason.

```bash
agent-team task complete --task TASK_ID --agent AGENT --evidence "verified" --artifact "_workspace/RUN_ID/TASK_ID.md"
agent-team task block --task TASK_ID --agent AGENT --reason "specific blocker"
```

## Notes

- If `data.sync.blocking` is true, do not complete unless the orchestrator explicitly approves `--force`.
- Evidence should say what was verified, not just that work is done.
- Artifact paths should be run-scoped.

