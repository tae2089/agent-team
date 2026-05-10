---
name: agent-team-run-cancel
description: "agent-team: Cancel an open workflow run with a reason."
metadata:
  version: 1.0.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
  cliHelp: "agent-team run cancel --help"
---

# run cancel

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Cancel a run that should no longer proceed.

## Usage

```bash
agent-team run cancel --run RUN_ID --reason TEXT
```

## Flags

| Flag       | JSON key | Required | Default | Description                                         |
| ---------- | -------- | -------- | ------- | --------------------------------------------------- |
| `--run`    | `run_id` | yes      | -       | Run to cancel.                                      |
| `--reason` | `reason` | yes      | -       | Cancellation rationale stored in the event payload. |

## Examples

```bash
agent-team run cancel --run run_docs --reason "superseded by another run"
agent-team run cancel --params '{"run_id":"run_docs","reason":"superseded by another run"}'
```

## Errors

| Code                | Meaning                         | Action                  |
| ------------------- | ------------------------------- | ----------------------- |
| `invalid_run_state` | Closed run cannot be cancelled. | Inspect run status.     |
| `validation_error`  | Required field is missing.      | Provide run and reason. |

## See Also

- [agent-team-run-summary](../agent-team-run-summary/SKILL.md)
- [agent-team-event-log](../agent-team-event-log/SKILL.md)
