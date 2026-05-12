---
name: agent-team-inbox-ack
description: "agent-team: Acknowledge a message after handling it."
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
  cliHelp: "agent-team inbox ack --help"
---

# inbox ack

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Acknowledge a message after incorporating, answering, or explicitly deferring it.

## Usage

```bash
agent-team inbox ack --msg MSG_ID --agent AGENT
```

## Flags

| Flag | JSON key | Required | Default | Description |
|------|----------|----------|---------|-------------|
| `--msg`   | `msg_id` | yes      | -       | Message to acknowledge.                                       |
| `--agent` | `agent`  | yes      | -       | Recipient acknowledging the message; must match message `to`. |

## Examples

```bash
agent-team inbox ack --msg msg_contract --agent writer
agent-team inbox ack --params '{"msg_id":"msg_contract","agent":"writer"}'
```

## Errors

| Code               | Meaning                             | Action                                          |
| ------------------ | ----------------------------------- | ----------------------------------------------- |
| `agent_mismatch`   | Agent is not the message recipient. | Use the recipient agent or inspect the message. |
| `not_found`        | Message does not exist.             | Check the message ID.                           |
| `validation_error` | Required field is missing.          | Provide message ID and agent.                   |

## See Also

- [agent-team-inbox-list](../agent-team-inbox-list/SKILL.md)
- [agent-team-sync-check](../agent-team-sync-check/SKILL.md)
