---
name: recipe-agent-team-operational-audit
description: "Recipe: Read-only daemonless audit of agent-team runs, messages, events, stale tasks, unresolved inboxes, and CLI schema. Use for operational inspection, recovery assessment, runtime debugging, trace review, or contract export. Do not use to create runs, dispatch workers, mutate state, plan, design, align terms, compound learnings, or implement."
metadata:
  version: 1.1.0
  openclaw:
    category: "recipe"
    domain: "agent-orchestration"
  requires:
    bins:
      - agent-team
    skills:
      - agent-team-shared
      - agent-team-ops
---

# Agent Team Operational Audit

Use this recipe to inspect state, unresolved communication, event history, or the command contract.

## Boundary

Use this recipe for audit and diagnosis only. Required runtime skills are declared in metadata; load `agent-team-shared` first when executing commands. Do not mutate state from this recipe except for a separate user-approved recovery step using the run lifecycle or task command skills.

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
agent-team event log --run RUN_ID
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
- If handing off to run lifecycle for closure, close runs only after all tasks are in a terminal state: `done`, `failed`, or `cancelled`. Tasks still `in_progress`, `blocked`, or `pending` will cause `run_not_ready`.

## Completion

This recipe is complete when:

- active runs, stale tasks, unresolved messages, and event history relevant to the question have been inspected
- findings cite concrete `run_id`, `task_id`, or `msg_id` values, not summaries alone
- no state mutation occurred from this recipe
- recovery handoff (if any) names the lifecycle or task recipe that will mutate state, with explicit user approval recorded
