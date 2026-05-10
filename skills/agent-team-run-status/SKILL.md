---
name: agent-team-run-status
description: "agent-team: Show one workflow run and grouped task status counts."
metadata:
  version: 1.0.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
  cliHelp: "agent-team run status --help"
---

# run status

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Show one run and task counts grouped by status.

## Usage

```bash
agent-team run status --run RUN_ID
```

## Flags

| Flag | JSON key | Required | Default | Description |
|------|----------|----------|---------|-------------|
| `--run` | `run_id` | yes      | -       | Run to inspect. |

## Examples

```bash
agent-team run status --run run_docs
agent-team run status --params '{"run_id":"run_docs"}'
```

## Output

- `data.run`: run record.
- `data.tasks`: task counts by status.

## Errors

| Code               | Meaning              | Action            |
| ------------------ | -------------------- | ----------------- |
| `validation_error` | `run_id` is missing. | Provide `--run`.  |
| `not_found`        | Run does not exist.  | Check the run ID. |

## See Also

- [agent-team-run-list](../agent-team-run-list/SKILL.md)
- [agent-team-task-list](../agent-team-task-list/SKILL.md)
