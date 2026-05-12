---
name: agent-team-message-send
description: "agent-team: Send a compact run-scoped message to another agent."
metadata:
  version: 1.0.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
  cliHelp: "agent-team message send --help"
---

# message send

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Send compact state-changing information, warnings, questions, or contract updates.

## Usage

```bash
agent-team message send --run RUN_ID --from SENDER --to RECIPIENT --kind KIND --body TEXT [--task TASK_ID] [--id MSG_ID] [--metadata JSON_OBJECT]
```

## Flags

| Flag | JSON key | Required | Default | Description |
|------|----------|----------|---------|-------------|
| `--run`      | `run_id`   | yes      | -         | Workflow scope.                               |
| `--from`     | `from`     | yes      | -         | Sender agent or orchestrator role.            |
| `--to`       | `to`       | yes      | -         | Recipient agent; drives `inbox list --agent`. |
| `--kind`     | `kind`     | yes      | -         | Machine-readable message category.            |
| `--body`     | `body`     | yes      | -         | Compact human-readable content.               |
| `--task`     | `task_id`  | no       | empty     | Optional task scope.                          |
| `--id`       | `id`       | no       | generated | Stable message ID.                            |
| `--metadata` | `metadata` | no       | `{}`      | JSON object for structured context.           |

## Examples

```bash
agent-team message send --run run_docs --task task_docs --from planner --to writer --kind contract_changed --body "Use assignee instead of owner."
agent-team message send --params '{"run_id":"run_docs","from":"planner","to":"writer","kind":"question"}' --json '{"body":"Do you need API schema?","metadata":{"severity":"normal"}}'
```

## Recommended Kinds

- `progress`
- `dependency_ready`
- `contract_changed`
- `conflict_detected`
- `question`
- `result_note`
- `approval`

## Errors

| Code               | Meaning                          | Action                                         |
| ------------------ | -------------------------------- | ---------------------------------------------- |
| `validation_error` | Required field is missing.       | Provide run, from, to, kind, and body.         |
| `invalid_json`     | `metadata` is not a JSON object. | Use an object such as `{"severity":"normal"}`. |

## See Also

- [agent-team-inbox-list](../agent-team-inbox-list/SKILL.md)
- [agent-team-message-list](../agent-team-message-list/SKILL.md)
