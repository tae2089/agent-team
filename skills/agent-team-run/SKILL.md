---
name: agent-team-run
description: "Service-level skill for agent-team run commands. Use when creating, listing, checking, or closing daemonless workflow runs. Load agent-team-shared first."
metadata:
  version: 1.1.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
  cliHelp: "agent-team run --help"
---

# Agent Team Run Commands

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Runs group tasks, inbox messages, event history, and artifacts for one orchestrated workflow.

```bash
agent-team run <command> [flags]
```

## Helper Commands

| Command | Description |
|---------|-------------|
| [`create`](../agent-team-run-create/SKILL.md) | Create a workflow run. |
| [`status`](../agent-team-run-status/SKILL.md) | Show one run and task status counts. |
| [`summary`](../agent-team-run-summary/SKILL.md) | Show operational run summary and close readiness. |
| [`list`](../agent-team-run-list/SKILL.md) | List runs, optionally filtered by status. |
| [`close`](../agent-team-run-close/SKILL.md) | Close a run after all tasks are done. |
| [`cancel`](../agent-team-run-cancel/SKILL.md) | Cancel an open run. |

## Command Notes

- Create one run per user-requested orchestrated workflow.
- Reuse an existing `run_id` only when the user provides it or explicitly asks to resume.
- Use the returned `data.run.id` as `RUN_ID` in task prompts and artifact paths.
- Close a run only after every task is `done`; otherwise `run close` returns `run_not_ready`.

## Discovering Commands

```bash
agent-team run --help
agent-team schema export
```
