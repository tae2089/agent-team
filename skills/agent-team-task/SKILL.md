---
name: agent-team-task
description: "Service-level skill for agent-team task lifecycle commands. Use when creating, listing, showing, starting, completing, blocking, reassigning, or retrying assigned work. Load agent-team-shared first."
metadata:
  version: 1.1.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
  cliHelp: "agent-team task --help"
---

# Agent Team Task Commands

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Tasks are assigned units of work inside a run.

```bash
agent-team task <command> [flags]
```

## Helper Commands

| Command | Description |
|---------|-------------|
| [`create`](../agent-team-task-create/SKILL.md) | Create a task for an agent. |
| [`list`](../agent-team-task-list/SKILL.md) | List tasks with optional run, agent, or status filters. |
| [`show`](../agent-team-task-show/SKILL.md) | Show one task and its dependencies. |
| [`start`](../agent-team-task-start/SKILL.md) | Mark an assigned task in progress. |
| [`complete`](../agent-team-task-complete/SKILL.md) | Complete a task with evidence and artifact path. |
| [`block`](../agent-team-task-block/SKILL.md) | Mark a task blocked with an actionable reason. |
| [`reassign`](../agent-team-task-reassign/SKILL.md) | Move a pending or blocked task to another agent. |
| [`retry`](../agent-team-task-retry/SKILL.md) | Reset a blocked or in-progress task to pending. |
| [`cancel`](../agent-team-task-cancel/SKILL.md) | Cancel a non-terminal task. |
| [`fail`](../agent-team-task-fail/SKILL.md) | Mark assigned work failed with a reason. |
| [`stale`](../agent-team-task-stale/SKILL.md) | Detect old blocked or in-progress tasks. |

## Command Notes

- Workers should run `agent-team sync check` before completing work.
- `start`, `complete`, and `block` require `--agent` to match the task assignee.
- `complete` requires both `--evidence` and `--artifact`.
- Use `--force` only with explicit orchestrator approval after a reported sync conflict.

## Discovering Commands

```bash
agent-team task --help
agent-team schema export
```
