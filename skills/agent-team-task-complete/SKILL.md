---
name: agent-team-task-complete
description: "agent-team: Complete an assigned task with evidence and artifact path."
metadata:
  version: 1.0.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
    skills:
      - agent-team-shared
      - agent-team-task
  cliHelp: "agent-team task complete --help"
---

# task complete

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md` and run `agent-team sync check` first.

Mark an assigned task as `done` with concrete evidence and a run-scoped artifact path.

## Usage

```bash
agent-team task complete --task TASK_ID --agent AGENT --evidence TEXT --artifact PATH [--force]
```

## Flags

| Flag | JSON key | Required | Default | Description |
|------|----------|----------|---------|-------------|
| `--task`     | `task_id`  | yes      | -       | Task to mark `done`.                                                |
| `--agent`    | `agent`    | yes      | -       | Completing agent; must match the task assignee.                     |
| `--evidence` | `evidence` | yes      | -       | Verification summary. Mention tests, inspected files, or decisions. |
| `--artifact` | `artifact` | yes      | -       | Result file or directory path, normally under `_workspace/{run_id}/`. |
| `--force`    | `force`    | no       | `false` | Bypass sync conflict after explicit orchestrator approval.          |

## Examples

```bash
agent-team task complete --task task_docs --agent writer --evidence "Links verified." --artifact "_workspace/run_docs/task_docs.md"
agent-team task complete --task task_docs --agent writer --evidence "Accepted conflict." --artifact "_workspace/run_docs/task_docs.md" --force
agent-team task complete --params '{"task_id":"task_docs","agent":"writer","force":true}' --json '{"evidence":"Links verified.","artifact":"_workspace/run_docs/task_docs.md"}'
```

## Errors

| Code               | Meaning                                     | Action                                       |
| ------------------ | ------------------------------------------- | -------------------------------------------- |
| `agent_mismatch`   | Agent is not the task assignee.             | Reassign or use the assigned agent.          |
| `sync_conflict`    | Unread messages or incomplete dependencies. | Run sync check, resolve issues, then retry.  |
| `validation_error` | Required field is missing.                  | Provide task, agent, evidence, and artifact. |

## See Also

- [agent-team-sync-check](../agent-team-sync-check/SKILL.md)
- [agent-team-task-block](../agent-team-task-block/SKILL.md)
