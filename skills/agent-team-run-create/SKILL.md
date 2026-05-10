---
name: agent-team-run-create
description: "agent-team: Create a daemonless workflow run."
metadata:
  version: 1.0.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
  cliHelp: "agent-team run create --help"
---

# run create

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Create a run that groups tasks, messages, events, and artifacts for one workflow.

## Usage

```bash
agent-team run create --title TITLE [--id RUN_ID]
```

## Flags

| Flag      | JSON key | Required | Default   | Description                                                   |
| --------- | -------- | -------- | --------- | ------------------------------------------------------------- |
| `--title` | `title`  | yes      | -         | Short workflow label shown in run lists and status output.    |
| `--id`    | `id`     | no       | generated | Stable run ID. Use only when scripts or docs need a known ID. |

## Examples

```bash
agent-team run create --title "docs refactor"
agent-team run create --id run_docs_refactor --title "docs refactor"
agent-team run create --params '{"id":"run_docs_refactor","title":"docs refactor"}'
```

## Errors

| Code               | Meaning                               | Action                                  |
| ------------------ | ------------------------------------- | --------------------------------------- |
| `validation_error` | `title` is missing.                   | Provide `--title` or `{"title":"..."}`. |
| `input_conflict`   | Same value supplied by flag and JSON. | Provide the value once.                 |

## See Also

- [agent-team-shared](../agent-team-shared/SKILL.md)
- [agent-team-run-status](../agent-team-run-status/SKILL.md)
