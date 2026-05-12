---
name: agent-team-task-retry
description: "agent-team: Reset blocked, in-progress, or failed work to pending."
metadata:
  version: 1.0.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
  cliHelp: "agent-team task retry --help"
---

# task retry

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Move a `blocked`, `in_progress`, or `failed` task back to `pending` for another attempt.

## Usage

```bash
agent-team task retry --task TASK_ID --reason TEXT
```

## Flags

| Flag | JSON key | Required | Default | Description |
|------|----------|----------|---------|-------------|
| `--task`   | `task_id` | yes      | -       | Task to reset to `pending`.                  |
| `--reason` | `reason`  | yes      | -       | Retry rationale stored in the event payload. |

## Examples

```bash
agent-team task retry --task task_docs --reason "API schema is now available"
agent-team task retry --params '{"task_id":"task_docs","reason":"API schema is now available"}'
```

## Behavior

- Allowed source status: `blocked`, `in_progress`, `failed`.
- Clears `evidence`, `artifact`, `blocked_reason`, and `started_version`.

## Errors

| Code                 | Meaning                                            | Action                                         |
| -------------------- | -------------------------------------------------- | ---------------------------------------------- |
| `invalid_task_state` | Task is not `blocked`, `in_progress`, or `failed`. | Inspect task and choose the correct operation. |
| `validation_error`   | Required field is missing.                         | Provide task and reason.                       |

## See Also

- [agent-team-task-show](../agent-team-task-show/SKILL.md)
- [agent-team-event-log](../agent-team-event-log/SKILL.md)
