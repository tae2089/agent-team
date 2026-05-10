---
name: agent-team-message-list
description: "agent-team: List messages in a run with optional audit filters."
metadata:
  version: 1.0.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
  cliHelp: "agent-team message list --help"
---

# message list

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Inspect run-scoped message history.

## Usage

```bash
agent-team message list --run RUN_ID [--task TASK_ID] [--from SENDER] [--to RECIPIENT] [--kind KIND] [--unread] [--limit N] [--after-version VERSION]
```

## Flags

| Flag              | JSON key        | Required | Default | Description                                             |
| ----------------- | --------------- | -------- | ------- | ------------------------------------------------------- |
| `--run`           | `run_id`        | yes      | -       | Workflow scope for the query.                           |
| `--task`          | `task_id`       | no       | empty   | Restrict messages to one task; empty means all tasks.   |
| `--from`          | `from`          | no       | empty   | Restrict by sender; empty means all senders.            |
| `--to`            | `to`            | no       | empty   | Restrict by recipient; empty means all recipients.      |
| `--kind`          | `kind`          | no       | empty   | Restrict by message category; empty means all kinds.    |
| `--unread`        | `unread`        | no       | `false` | Return only messages where `acked_at` is empty.         |
| `--limit`         | `limit`         | no       | `100`   | Maximum rows returned, capped at 1000.                  |
| `--after-version` | `after_version` | no       | `0`     | Only include messages changed after this state version. |

## Examples

```bash
agent-team message list --run run_docs
agent-team message list --run run_docs --to writer --unread
agent-team message list --params '{"run_id":"run_docs","kind":"contract_changed","unread":true,"limit":50}'
```

## Output

- `data.messages`: array of message records.

## Errors

| Code               | Meaning                               | Action                  |
| ------------------ | ------------------------------------- | ----------------------- |
| `validation_error` | `run_id` is missing.                  | Provide `--run`.        |
| `input_conflict`   | Same value supplied by flag and JSON. | Provide the value once. |

## See Also

- [agent-team-message-send](../agent-team-message-send/SKILL.md)
- [agent-team-inbox-list](../agent-team-inbox-list/SKILL.md)
