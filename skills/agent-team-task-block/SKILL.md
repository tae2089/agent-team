---
name: agent-team-task-block
description: "agent-team: Mark an assigned task blocked with an actionable reason."
metadata:
  version: 1.0.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
  cliHelp: "agent-team task block --help"
---

# task block

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Mark a task `blocked` when the assigned agent cannot proceed.

## Usage

```bash
agent-team task block --task TASK_ID --agent AGENT --reason TEXT
```

## Flags

| Flag       | JSON key  | Required | Default | Description                                                                  |
| ---------- | --------- | -------- | ------- | ---------------------------------------------------------------------------- |
| `--task`   | `task_id` | yes      | -       | Task to mark `blocked`.                                                      |
| `--agent`  | `agent`   | yes      | -       | Reporting agent; must match the task assignee.                               |
| `--reason` | `reason`  | yes      | -       | Actionable blocker: missing input, decision, dependency, or error condition. |

## Examples

```bash
agent-team task block --task task_docs --agent writer --reason "Missing API schema."
agent-team task block --params '{"task_id":"task_docs","agent":"writer","reason":"Missing API schema."}'
```

## Errors

| Code               | Meaning                         | Action                              |
| ------------------ | ------------------------------- | ----------------------------------- |
| `agent_mismatch`   | Agent is not the task assignee. | Reassign or use the assigned agent. |
| `validation_error` | Required field is missing.      | Provide task, agent, and reason.    |

## See Also

- [agent-team-task-reassign](../agent-team-task-reassign/SKILL.md)
- [agent-team-task-retry](../agent-team-task-retry/SKILL.md)
