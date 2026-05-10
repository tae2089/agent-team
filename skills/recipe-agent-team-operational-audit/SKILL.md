---
name: recipe-agent-team-operational-audit
description: "Recipe: Audit agent-team runs, messages, events, and schema without a daemon."
metadata:
  version: 1.0.0
  openclaw:
    category: "recipe"
    domain: "agent-orchestration"
  requires:
    bins:
      - agent-team
    skills:
      - agent-team-ops
---

# Agent Team Operational Audit

Use this recipe to inspect state, unresolved communication, event history, or the command contract.

## Prerequisites

Load:

- `agent-team-shared`
- `agent-team-ops`

## Steps

1. Find active runs.

```bash
agent-team run list --status open
```

2. Inspect run status and summary.

```bash
agent-team run status --run RUN_ID
agent-team run summary --run RUN_ID
```

3. Inspect stale or unresolved work.

```bash
agent-team task stale --run RUN_ID --older-than 2h
```

4. Inspect unresolved run messages.

```bash
agent-team message list --run RUN_ID --unread
```

5. Inspect event history.

```bash
agent-team event log --run RUN_ID --limit 100
```

6. Poll incrementally when needed.

```bash
agent-team event log --run RUN_ID --after-version STATE_VERSION --limit 100
```

7. Export the CLI contract for tooling or docs.

```bash
agent-team schema export
```

## Notes

- Use `event log` for auditability, not for replacing task or inbox commands.
- Use `schema export` as the command contract source of truth.
- Close runs only after `run status` shows all tasks are `done`.
