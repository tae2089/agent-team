---
name: agent-team-run-summary
description: "agent-team: Show operational summary for one run, including blockers, unread messages, recent events, and close readiness."
metadata:
  version: 1.0.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
    skills:
      - agent-team-shared
      - agent-team-run
  cliHelp: "agent-team run summary --help"
---

# run summary

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Inspect one run as an operational dashboard.

## Usage

```bash
agent-team run summary --run RUN_ID [--recent-limit N]
```

## Flags

| Flag | JSON key | Required | Default | Description |
|------|----------|----------|---------|-------------|
| `--run`          | `run_id`       | yes      | -       | Run to summarize.                                  |
| `--recent-limit` | `recent_limit` | no       | `10`    | Number of recent events to include, capped at 100. |

## Examples

```bash
agent-team run summary --run run_docs
agent-team run summary --params '{"run_id":"run_docs","recent_limit":20}'
```

## Errors

| Code               | Meaning              | Action            |
| ------------------ | -------------------- | ----------------- |
| `validation_error` | `run_id` is missing. | Provide `--run`.  |
| `not_found`        | Run does not exist.  | Check the run ID. |

## See Also

- [agent-team-run-status](../agent-team-run-status/SKILL.md)
- [agent-team-event-log](../agent-team-event-log/SKILL.md)
