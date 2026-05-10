---
name: agent-team-task-start
description: "agent-team: Mark an assigned task in progress."
metadata:
  version: 1.0.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
  cliHelp: "agent-team task start --help"
---

# task start

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Mark a task as `in_progress` and record the current `state_version`.

## Usage

```bash
agent-team task start --task TASK_ID --agent AGENT
```

## Flags

| Flag      | JSON key  | Required | Default | Description                                    |
| --------- | --------- | -------- | ------- | ---------------------------------------------- |
| `--task`  | `task_id` | yes      | -       | Task to claim as `in_progress`.                |
| `--agent` | `agent`   | yes      | -       | Caller identity; must match the task assignee. |

## Examples

```bash
agent-team task start --task task_docs --agent writer
agent-team task start --params '{"task_id":"task_docs","agent":"writer"}'
```

## Errors

| Code             | Meaning                         | Action                              |
| ---------------- | ------------------------------- | ----------------------------------- |
| `agent_mismatch` | Agent is not the task assignee. | Reassign or use the assigned agent. |
| `not_found`      | Task does not exist.            | Check the task ID.                  |

## See Also

- [agent-team-sync-check](../agent-team-sync-check/SKILL.md)
- [agent-team-task-complete](../agent-team-task-complete/SKILL.md)
