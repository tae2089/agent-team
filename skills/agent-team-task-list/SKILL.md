---
name: agent-team-task-list
description: "agent-team: List tasks with optional run, agent, or status filters."
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
  cliHelp: "agent-team task list --help"
---

# task list

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

List tasks for planning, assignment, worker pickup, or operational review.

## Usage

```bash
agent-team task list [--run RUN_ID] [--agent AGENT] [--status STATUS] [--limit N] [--after-version VERSION]
```

## Flags

| Flag | JSON key | Required | Default | Description |
|------|----------|----------|---------|-------------|
| `--run`           | `run_id`        | no       | empty   | Limit tasks to one run; empty means all runs.        |
| `--agent`         | `agent`         | no       | empty   | Limit tasks to one assignee; empty means all agents. |
| `--status`        | `status`        | no       | empty   | Exact task status filter; empty means all statuses.  |
| `--limit`         | `limit`         | no       | `100`   | Maximum rows returned, capped at 1000.               |
| `--after-version` | `after_version` | no       | `0`     | Only include tasks changed after this state version. |

## Examples

```bash
agent-team task list --run run_docs
agent-team task list --run run_docs --agent writer --status pending
agent-team task list --params '{"run_id":"run_docs","status":"blocked","limit":50}'
```

## Output

- `data.tasks`: array of task records.

## Errors

| Code             | Meaning                               | Action                  |
| ---------------- | ------------------------------------- | ----------------------- |
| `input_conflict` | Same value supplied by flag and JSON. | Provide the value once. |

## See Also

- [agent-team-task-show](../agent-team-task-show/SKILL.md)
- [agent-team-run-status](../agent-team-run-status/SKILL.md)
