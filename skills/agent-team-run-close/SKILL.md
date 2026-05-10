---
name: agent-team-run-close
description: "agent-team: Close a workflow run after all tasks are done."
metadata:
  version: 1.0.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
  cliHelp: "agent-team run close --help"
---

# run close

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Finalize a run after all tasks are in terminal status.

## Usage

```bash
agent-team run close --run RUN_ID [--reason TEXT]
```

## Flags

| Flag       | JSON key | Required | Default | Description                                            |
| ---------- | -------- | -------- | ------- | ------------------------------------------------------ |
| `--run`    | `run_id` | yes      | -       | Run to close.                                          |
| `--reason` | `reason` | no       | empty   | Closure note stored in the `run_closed` event payload. |

## Examples

```bash
agent-team run close --run run_docs --reason "all tasks complete"
agent-team run close --params '{"run_id":"run_docs","reason":"all tasks complete"}'
```

## Errors

| Code                | Meaning                      | Action                                                                    |
| ------------------- | ---------------------------- | ------------------------------------------------------------------------- |
| `run_not_ready`     | Some tasks are not terminal. | Inspect task counts, then finish, retry, cancel, fail, or reassign tasks. |
| `invalid_run_state` | Run is not open.             | Inspect run status.                                                       |
| `not_found`         | Run does not exist.          | Check the run ID.                                                         |

## See Also

- [agent-team-run-status](../agent-team-run-status/SKILL.md)
- [agent-team-event-log](../agent-team-event-log/SKILL.md)
