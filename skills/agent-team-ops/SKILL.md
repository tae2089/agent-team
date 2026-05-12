---
name: agent-team-ops
description: "Service-level skill for agent-team operational commands. Use for run hygiene, message inspection, event history, schema export, audits, and tooling. Load agent-team-shared first."
metadata:
  version: 1.0.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
    skills:
      - agent-team-shared
  cliHelp: "agent-team event --help && agent-team schema export"
---

# Agent Team Operational Commands

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Operational commands inspect or finalize state without running agents.

## Helper Commands

| Command                                                 | Description                                                          |
| ------------------------------------------------------- | -------------------------------------------------------------------- |
| [`run list`](../agent-team-run-list/SKILL.md)           | Find open or closed runs.                                            |
| [`run summary`](../agent-team-run-summary/SKILL.md)     | Inspect run readiness, blockers, unread messages, and recent events. |
| [`run close`](../agent-team-run-close/SKILL.md)         | Finalize a run after all tasks are terminal.                         |
| [`message list`](../agent-team-message-list/SKILL.md)   | Inspect run-scoped message history.                                  |
| [`event log`](../agent-team-event-log/SKILL.md)         | Inspect append-only state events.                                    |
| [`task stale`](../agent-team-task-stale/SKILL.md)       | Find old in-progress or blocked tasks.                               |
| [`schema export`](../agent-team-schema-export/SKILL.md) | Export the machine-readable CLI contract.                            |

## Command Notes

- Use `event log --after-version` for incremental polling without a daemon.
- Use `message list --unread` to audit unresolved run communication.
- Use `schema export` as the source of truth for generated docs and tooling.

## Discovering Commands

```bash
agent-team run --help
agent-team message --help
agent-team event --help
agent-team schema export
```
