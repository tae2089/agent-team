---
name: agent-team-task-show
description: "agent-team: Show one task and its dependency IDs."
metadata:
  version: 1.1.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
  cliHelp: "agent-team task show --help"
---

# task show

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Inspect one task before starting, retrying, reassigning, or integrating work.

## Usage

```bash
agent-team task show --task TASK_ID
```

## Flags

| Flag | JSON key | Required | Default | Description |
|------|----------|----------|---------|-------------|
| `--task` | `task_id` | yes | - | Task to inspect. |

## Examples

```bash
agent-team task show --task task_docs
agent-team task show --params '{"task_id":"task_docs"}'
```

## Output

- `data.task`: task record.
- `data.depends_on`: dependency task IDs.

## Errors

| Code | Meaning | Action |
|------|---------|--------|
| `not_found` | Task does not exist. | Check the task ID. |
| `validation_error` | `task_id` is missing. | Provide `--task`. |

## See Also

- [agent-team-task-list](../agent-team-task-list/SKILL.md)
- [agent-team-sync-check](../agent-team-sync-check/SKILL.md)

