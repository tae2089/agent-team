---
name: agent-team-task-create
description: "agent-team: Create an assigned task inside a run."
metadata:
  version: 1.0.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
  cliHelp: "agent-team task create --help"
---

# task create

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Create a task assigned to one agent, optionally with dependencies.

## Usage

```bash
agent-team task create --run RUN_ID --agent AGENT --title TITLE [--id TASK_ID] [--depends-on TASK_ID] [--body TEXT] [--metadata JSON_OBJECT]
```

## Flags

| Flag           | JSON key     | Required | Default   | Description                                                   |
| -------------- | ------------ | -------- | --------- | ------------------------------------------------------------- |
| `--run`        | `run_id`     | yes      | -         | Owning workflow run.                                          |
| `--agent`      | `agent`      | yes      | -         | Assignee expected to start, block, or complete the task.      |
| `--title`      | `title`      | yes      | -         | Compact task label. Put detailed instructions in `--body`.    |
| `--id`         | `id`         | no       | generated | Stable task ID.                                               |
| `--depends-on` | `depends_on` | no       | `[]`      | Task ID that must be `done`; repeat flag or use a JSON array. |
| `--body`       | `body`       | no       | empty     | Worker contract: scope, expected output, and constraints.     |
| `--metadata`   | `metadata`   | no       | `{}`      | JSON object for structured context.                           |

## Examples

```bash
agent-team task create --run run_docs --agent writer --title "draft API docs" --body "Use the research artifact."
agent-team task create --params '{"run_id":"run_docs","agent":"writer","title":"draft API docs","depends_on":["task_research"]}' --json '{"body":"Use the research artifact.","metadata":{"scope":"api"}}'
```

## Errors

| Code               | Meaning                                   | Action                                   |
| ------------------ | ----------------------------------------- | ---------------------------------------- |
| `validation_error` | `run_id`, `agent`, or `title` is missing. | Provide required fields.                 |
| `invalid_json`     | `metadata` is not a JSON object.          | Use an object such as `{"scope":"api"}`. |

## See Also

- [agent-team-task-start](../agent-team-task-start/SKILL.md)
- [agent-team-task-list](../agent-team-task-list/SKILL.md)
