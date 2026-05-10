---
name: agent-team-task-reassign
description: "agent-team: Reassign a pending or blocked task to another agent."
metadata:
  version: 1.1.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
  cliHelp: "agent-team task reassign --help"
---

# task reassign

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Move a `pending` or `blocked` task to another agent.

## Usage

```bash
agent-team task reassign --task TASK_ID --agent NEW_AGENT --reason TEXT
```

## Flags

| Flag | JSON key | Required | Default | Description |
|------|----------|----------|---------|-------------|
| `--task` | `task_id` | yes | - | Task to move. |
| `--agent` | `agent` | yes | - | New assignee. |
| `--reason` | `reason` | yes | - | Handoff rationale stored in the event payload. |

## Examples

```bash
agent-team task reassign --task task_docs --agent backup-writer --reason "primary worker unavailable"
agent-team task reassign --params '{"task_id":"task_docs","agent":"backup-writer","reason":"primary worker unavailable"}'
```

## Behavior

- Allowed source status: `pending`, `blocked`.
- If the task is `blocked`, reassignment resets it to `pending` and clears `blocked_reason`.

## Errors

| Code | Meaning | Action |
|------|---------|--------|
| `invalid_task_state` | Task is not `pending` or `blocked`. | Retry only eligible tasks or inspect status. |
| `validation_error` | Required field is missing. | Provide task, new agent, and reason. |

## See Also

- [agent-team-task-list](../agent-team-task-list/SKILL.md)
- [agent-team-event-log](../agent-team-event-log/SKILL.md)

