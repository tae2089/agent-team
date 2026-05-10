---
name: agent-team-task-fail
description: "agent-team: Mark assigned work failed with a reason and optional artifact."
metadata:
  version: 1.0.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
  cliHelp: "agent-team task fail --help"
---

# task fail

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Mark assigned work failed when it reached a terminal unsuccessful outcome.

## Usage

```bash
agent-team task fail --task TASK_ID --agent AGENT --reason TEXT [--artifact PATH]
```

## Flags

| Flag | JSON key | Required | Default | Description |
|------|----------|----------|---------|-------------|
| `--task`     | `task_id`  | yes      | -       | Task to mark failed.                           |
| `--agent`    | `agent`    | yes      | -       | Reporting agent; must match task assignee.     |
| `--reason`   | `reason`   | yes      | -       | Failure rationale stored in the event payload. |
| `--artifact` | `artifact` | no       | empty   | Optional failure report path.                  |

## Examples

```bash
agent-team task fail --task task_docs --agent writer --reason "tests failed" --artifact "_workspace/run_docs/task_docs_failure.md"
agent-team task fail --params '{"task_id":"task_docs","agent":"writer","reason":"tests failed"}'
```

## Errors

| Code                 | Meaning                         | Action                              |
| -------------------- | ------------------------------- | ----------------------------------- |
| `agent_mismatch`     | Agent is not the task assignee. | Reassign or use the assigned agent. |
| `invalid_task_state` | Task is already terminal.       | Inspect task status.                |

## See Also

- [agent-team-task-retry](../agent-team-task-retry/SKILL.md)
- [agent-team-event-log](../agent-team-event-log/SKILL.md)
