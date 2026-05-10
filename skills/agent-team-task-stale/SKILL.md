---
name: agent-team-task-stale
description: "agent-team: Detect old in-progress or blocked tasks for daemonless operational recovery."
metadata:
  version: 1.0.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
  cliHelp: "agent-team task stale --help"
---

# task stale

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Find stale `in_progress` or `blocked` tasks without requiring a daemon.

## Usage

```bash
agent-team task stale --run RUN_ID --older-than DURATION [--limit N] [--after-version VERSION]
```

## Flags

| Flag              | JSON key        | Required | Default | Description                                          |
| ----------------- | --------------- | -------- | ------- | ---------------------------------------------------- |
| `--run`           | `run_id`        | yes      | -       | Run to inspect.                                      |
| `--older-than`    | `older_than`    | yes      | -       | Go duration such as `30m`, `2h`, or `24h`.           |
| `--limit`         | `limit`         | no       | `100`   | Maximum rows returned, capped at 1000.               |
| `--after-version` | `after_version` | no       | `0`     | Only include tasks changed after this state version. |

## Examples

```bash
agent-team task stale --run run_docs --older-than 2h
agent-team task stale --params '{"run_id":"run_docs","older_than":"30m","limit":50}'
```

## Errors

| Code               | Meaning                             | Action                                      |
| ------------------ | ----------------------------------- | ------------------------------------------- |
| `validation_error` | Run or duration is missing/invalid. | Provide `--run` and a valid `--older-than`. |

## See Also

- [agent-team-run-summary](../agent-team-run-summary/SKILL.md)
- [agent-team-task-reassign](../agent-team-task-reassign/SKILL.md)
