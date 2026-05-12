---
name: agent-team-inbox-list
description: "agent-team: Read messages for one recipient agent."
metadata:
  version: 1.0.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
    skills:
      - agent-team-shared
      - agent-team-inbox
  cliHelp: "agent-team inbox list --help"
---

# inbox list

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Read a recipient inbox before starting work, at checkpoints, and before completion.

## Usage

```bash
agent-team inbox list --agent AGENT [--run RUN_ID] [--unread]
```

## Flags

| Flag | JSON key | Required | Default | Description |
|------|----------|----------|---------|-------------|
| `--agent`  | `agent`  | yes      | -       | Recipient inbox to read.                                      |
| `--run`    | `run_id` | no       | empty   | Restrict messages to one workflow; empty means all workflows. |
| `--unread` | `unread` | no       | `false` | Return only messages where `acked_at` is empty.               |

## Examples

```bash
agent-team inbox list --agent writer --run run_docs --unread
agent-team inbox list --params '{"agent":"writer","run_id":"run_docs","unread":true}'
```

## Output

- `data.messages`: array of messages addressed to `agent`.

## Errors

| Code               | Meaning             | Action             |
| ------------------ | ------------------- | ------------------ |
| `validation_error` | `agent` is missing. | Provide `--agent`. |

## See Also

- [agent-team-inbox-ack](../agent-team-inbox-ack/SKILL.md)
- [agent-team-sync-check](../agent-team-sync-check/SKILL.md)
