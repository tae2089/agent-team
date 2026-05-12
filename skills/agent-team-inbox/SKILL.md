---
name: agent-team-inbox
description: "Service-level skill for agent-team message and inbox commands. Use for asynchronous daemonless communication, inbox checkpoint checks, acknowledgements, and compact contract-change notifications. Load agent-team-shared first."
metadata:
  version: 1.0.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
    skills:
      - agent-team-shared
  cliHelp: "agent-team message --help && agent-team inbox --help"
---

# Agent Team Message and Inbox Commands

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Messages are compact state-changing notices. Inboxes are recipient-specific message views.

```bash
agent-team message <command> [flags]
agent-team inbox <command> [flags]
```

## Helper Commands

| Command                                               | Description                                         |
| ----------------------------------------------------- | --------------------------------------------------- |
| [`message send`](../agent-team-message-send/SKILL.md) | Send a compact run-scoped message to another agent. |
| [`message list`](../agent-team-message-list/SKILL.md) | Audit messages in a run with optional filters.      |
| [`inbox list`](../agent-team-inbox-list/SKILL.md)     | Read messages for one recipient.                    |
| [`inbox ack`](../agent-team-inbox-ack/SKILL.md)       | Acknowledge a message after handling it.            |

## Recommended Message Kinds

- `progress`
- `dependency_ready`
- `contract_changed`
- `conflict_detected`
- `question`
- `result_note`
- `approval`

## Command Notes

- Keep message bodies compact; place large outputs in artifacts and send the path.
- Workers should check unread inbox messages at task start, before major checkpoints, and before completion.
- Unread relevant messages are treated as blocking sync mismatches during task completion.

## Discovering Commands

```bash
agent-team message --help
agent-team inbox --help
agent-team schema export
```
