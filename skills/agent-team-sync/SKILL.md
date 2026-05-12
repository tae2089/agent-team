---
name: agent-team-sync
description: "Service-level skill for agent-team sync checks. Use before workers complete tasks, when an orchestrator suspects drift, or when inbox, dependency, or state-version mismatches must be detected without a daemon. Load agent-team-shared first."
metadata:
  version: 1.0.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
    skills:
      - agent-team-shared
  cliHelp: "agent-team sync --help"
---

# Agent Team Sync Checks

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Sync checks detect unresolved inbox messages and incomplete dependencies without a daemon.

```bash
agent-team sync <command> [flags]
```

## Helper Commands

| Command                                      | Description                                                   |
| -------------------------------------------- | ------------------------------------------------------------- |
| [`check`](../agent-team-sync-check/SKILL.md) | Check whether an agent can safely proceed or complete a task. |

## Completion Rule

Workers must run sync check before `task complete`. If `data.sync.blocking` is true, resolve the issue first:

- read and ack relevant inbox messages
- wait for dependencies to complete
- ask the orchestrator for a decision
- block the task with a concrete reason

`--force` on `task complete` is reserved for explicit orchestrator approval.

## Discovering Commands

```bash
agent-team sync --help
agent-team schema export
```
