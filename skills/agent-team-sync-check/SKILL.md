---
name: agent-team-sync-check
description: "agent-team: Check unresolved inbox messages and incomplete dependencies before proceeding."
metadata:
  version: 1.1.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
  cliHelp: "agent-team sync check --help"
---

# sync check

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Check whether an agent has unresolved messages or incomplete task dependencies.

## Usage

```bash
agent-team sync check --agent AGENT [--run RUN_ID] [--task TASK_ID]
```

## Flags

| Flag | JSON key | Required | Default | Description |
|------|----------|----------|---------|-------------|
| `--agent` | `agent` | yes | - | Agent whose inbox and task context are checked. |
| `--run` | `run_id` | no | empty | Restrict check to one workflow; empty means all workflows. |
| `--task` | `task_id` | no | empty | Also check incomplete dependencies for this task; empty skips task-specific dependency checks. |

## Examples

```bash
agent-team sync check --agent writer --run run_docs --task task_docs
agent-team sync check --params '{"agent":"writer","run_id":"run_docs","task_id":"task_docs"}'
```

## Output

- `data.sync.blocking`: true when completion should stop.
- `data.sync.unread_messages`: relevant unread messages.
- `data.sync.issues`: concise issue labels.
- top-level `warnings`: includes a blocking sync warning when applicable.

## Errors

| Code | Meaning | Action |
|------|---------|--------|
| `validation_error` | `agent` is missing. | Provide `--agent`. |

## See Also

- [agent-team-inbox-list](../agent-team-inbox-list/SKILL.md)
- [agent-team-task-complete](../agent-team-task-complete/SKILL.md)
