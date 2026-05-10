---
name: agent-team-run-list
description: "agent-team: List workflow runs, optionally filtered by status."
metadata:
  version: 1.0.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
  cliHelp: "agent-team run list --help"
---

# run list

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

List runs for operational review or resume decisions.

## Usage

```bash
agent-team run list [--status STATUS] [--limit N] [--after-version VERSION]
```

## Flags

| Flag | JSON key | Required | Default | Description |
|------|----------|----------|---------|-------------|
| `--status`        | `status`        | no       | empty   | Exact run status filter, normally `open` or `closed`; empty means all statuses. |
| `--limit`         | `limit`         | no       | `100`   | Maximum rows returned, capped at 1000.                                          |
| `--after-version` | `after_version` | no       | `0`     | Only include runs changed after this state version.                             |

## Examples

```bash
agent-team run list
agent-team run list --status open
agent-team run list --params '{"status":"closed","limit":50}'
```

## Output

- `data.runs`: array of run records.

## Errors

| Code             | Meaning                               | Action                  |
| ---------------- | ------------------------------------- | ----------------------- |
| `input_conflict` | Same value supplied by flag and JSON. | Provide the value once. |

## See Also

- [agent-team-run-status](../agent-team-run-status/SKILL.md)
- [agent-team-run-close](../agent-team-run-close/SKILL.md)
