---
name: agent-team-event-log
description: "agent-team: Inspect append-only state events with run, entity, type, and version filters."
metadata:
  version: 1.0.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
  cliHelp: "agent-team event log --help"
---

# event log

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Inspect append-only state events for audit, debugging, or incremental polling.

## Usage

```bash
agent-team event log [--run RUN_ID] [--entity-type TYPE] [--entity ID] [--type EVENT_TYPE] [--after-version VERSION] [--limit N]
```

## Flags

| Flag              | JSON key        | Required | Default | Description                                                                                    |
| ----------------- | --------------- | -------- | ------- | ---------------------------------------------------------------------------------------------- |
| `--run`           | `run_id`        | no       | empty   | Restrict events to one workflow; empty means all runs.                                         |
| `--entity-type`   | `entity_type`   | no       | empty   | Restrict by entity type, usually `run`, `task`, or `message`; empty means all entity types.    |
| `--entity`        | `entity_id`     | no       | empty   | Restrict to one entity ID; empty means all entities.                                           |
| `--type`          | `event_type`    | no       | empty   | Restrict by event type such as `task_reassigned` or `run_closed`; empty means all event types. |
| `--after-version` | `after_version` | no       | `0`     | Return events with `state_version` greater than this value.                                    |
| `--limit`         | `limit`         | no       | `100`   | Maximum events returned; hard cap is `1000`.                                                   |

## Examples

```bash
agent-team event log --run run_docs --limit 50
agent-team event log --entity-type task --entity task_docs
agent-team event log --params '{"run_id":"run_docs","after_version":100,"limit":50}'
```

## Output

- `data.events`: events ordered by `state_version` ascending.

## Errors

| Code               | Meaning                                      | Action                      |
| ------------------ | -------------------------------------------- | --------------------------- |
| `validation_error` | `limit` is less than 1 or greater than 1000. | Use a limit from 1 to 1000. |
| `input_conflict`   | Same value supplied by flag and JSON.        | Provide the value once.     |

## See Also

- [agent-team-run-status](../agent-team-run-status/SKILL.md)
- [agent-team-message-list](../agent-team-message-list/SKILL.md)
