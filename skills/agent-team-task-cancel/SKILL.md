---
name: agent-team-task-cancel
description: "agent-team: Cancel a non-terminal task with a reason."
metadata:
  version: 1.0.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
    skills:
      - agent-team-shared
      - agent-team-task
  cliHelp: "agent-team task cancel --help"
---

# task cancel

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Cancel a task that should not be completed or retried.

## Usage

```bash
agent-team task cancel --task TASK_ID --reason TEXT
```

## Flags

| Flag | JSON key | Required | Default | Description |
|------|----------|----------|---------|-------------|
| `--task`   | `task_id` | yes      | -       | Task to cancel.                                     |
| `--reason` | `reason`  | yes      | -       | Cancellation rationale stored in the event payload. |

## Examples

```bash
agent-team task cancel --task task_docs --reason "scope removed"
agent-team task cancel --params '{"task_id":"task_docs","reason":"scope removed"}'
```

## Errors

| Code                 | Meaning                    | Action                   |
| -------------------- | -------------------------- | ------------------------ |
| `invalid_task_state` | Task is already terminal.  | Inspect task status.     |
| `validation_error`   | Required field is missing. | Provide task and reason. |

## See Also

- [agent-team-task-stale](../agent-team-task-stale/SKILL.md)
- [agent-team-event-log](../agent-team-event-log/SKILL.md)
