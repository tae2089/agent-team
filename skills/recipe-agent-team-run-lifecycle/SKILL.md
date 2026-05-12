---
name: recipe-agent-team-run-lifecycle
description: "Recipe: Orchestrate full daemonless agent-team run lifecycle — create, execute, monitor, recover, close. Use for run/task creation, worker dispatch, monitoring, inbox handling, recovery, or closure. Do not use for worker checkpoints, planning grill, architecture design, terminology, compound learning, audit, or implementation outside a run."
metadata:
  version: 1.3.0
  openclaw:
    category: "recipe"
    domain: "agent-orchestration"
  requires:
    bins:
      - agent-team
    skills:
      - agent-team-shared
      - agent-team-run
      - agent-team-task
      - agent-team-inbox
---

# Agent Team Run Lifecycle

Use this recipe when an orchestrator needs to run a complete daemonless workflow.

## Boundary

Use this recipe only from the orchestrator role. Required runtime skills are declared in metadata; load `agent-team-shared` first when executing commands. Workers do not follow this recipe directly.

Workers load [`recipe-agent-team-worker-checkpoint`](../recipe-agent-team-worker-checkpoint/SKILL.md) independently after the orchestrator gives them `RUN_ID`, `TASK_ID`, `AGENT`, artifact root, and orchestrator recipient.

## Steps

### 1. Initialize state.

```bash
agent-team init
```

Idempotent: safe to re-run if state already exists.

### 2. Create one run.

```bash
agent-team run create --title "workflow title"
```

Use the returned `data.run.id` as `RUN_ID` in all subsequent commands.

### 3. Create all assigned tasks.

```bash
agent-team task create --run RUN_ID --agent AGENT --title "task title" --body "task contract"
```

When a task depends on another task's output, declare the dependency at creation time:

```bash
agent-team task create --run RUN_ID --agent AGENT --title "task title" --body "task contract" --depends-on UPSTREAM_TASK_ID
```

Workers will hit `incomplete_dependencies` during sync check if their upstream tasks are not yet `done`. Create all tasks before dispatching workers so dependency references are valid from the start.

Worker artifacts are written to `_workspace/RUN_ID/`. Reference the concrete run-scoped path in task contracts when workers need to share outputs.

### 4. Dispatch workers.

Workers independently follow [`recipe-agent-team-worker-checkpoint`](../recipe-agent-team-worker-checkpoint/SKILL.md).

The dispatch mechanism (spawning subagents, prompt format, harness invocation) is determined by the active environment. If using a harness, load the appropriate one before dispatching:

- [`agent-team-codex-harness`](../agent-team-codex-harness/SKILL.md) - for Codex-based workers
- [`agent-team-gemini-harness`](../agent-team-gemini-harness/SKILL.md) - for Gemini-based workers

Each worker prompt must include:

- `RUN_ID` and `TASK_ID`
- the task contract (body)
- artifact output path: `_workspace/RUN_ID/`
- instruction to follow `recipe-agent-team-worker-checkpoint`

### 5. Monitor until all tasks are terminal.

Poll with `run summary` to inspect blockers, unread messages, and close readiness. Poll after each worker dispatch and at each monitoring pass; do not poll continuously. Wait until a worker subagent turn completes or on a fixed interval, such as every 2-5 minutes in long-running workflows.

```bash
agent-team run summary --run RUN_ID
```

Check the orchestrator's own inbox at each monitoring pass. Workers send `question` messages when blocked on dependencies or contract ambiguity; these must be answered before the task can unblock.

```bash
agent-team inbox list --agent ORCHESTRATOR_AGENT --run RUN_ID --unread
```

Ack each message after handling it:

```bash
agent-team inbox ack --msg MSG_ID --agent ORCHESTRATOR_AGENT
```

When a worker sends a `question` message requesting `--force` or override approval, respond with an `approval` message before the worker can proceed:

```bash
agent-team message send --run RUN_ID --from ORCHESTRATOR_AGENT --to WORKER_AGENT --kind approval --body "approved: proceed with --force"
```

Detect tasks that have stopped progressing:

```bash
agent-team task stale --run RUN_ID --older-than 2h
```

Act on non-terminal task states:

| Task state            | Action                                                                    |
| --------------------- | ------------------------------------------------------------------------- |
| `in_progress`         | No action. Worker is active.                                              |
| `in_progress` (stale) | Worker likely dead. Retry to reset to `pending`, then reassign. Cancel if unrecoverable. |
| `blocked`             | Read reason with `task show`. Send `approval` inbox message, retry, or reassign. |
| `failed`              | Retry to reset to `pending`, then reassign. Cancel if unrecoverable.      |
| `pending`             | No worker claimed it. Reassign or cancel if no worker will claim it.      |

Recovery commands:

```bash
agent-team task show --task TASK_ID
agent-team task retry --task TASK_ID --reason "REASON"
agent-team task reassign --task TASK_ID --agent NEW_AGENT --reason "REASON"
agent-team task cancel --task TASK_ID --reason "REASON"
```

Repeat this step until all tasks are in a terminal state (`done`, `failed`, or `cancelled`), then proceed to Step 6.

### 6. Close the run.

`run close` requires every task to be in a terminal state: `done`, `failed`, or `cancelled`. Tasks still `in_progress`, `blocked`, or `pending` will cause `run_not_ready`.

```bash
agent-team run close --run RUN_ID --reason "all tasks complete"
```

If `run close` returns `run_not_ready`, repeat step 5. At least one task is not yet terminal.

## Abort Path

Cancel at any point if the workflow should not proceed.

`run cancel` does not auto-cancel open tasks. Cancel each non-terminal task first.

1. Identify non-terminal tasks:

```bash
agent-team task list --run RUN_ID
```

Returns all tasks. Cancel only non-terminal ones: `in_progress`, `blocked`, `pending`. Skip `done`, `failed`, `cancelled`.

2. Cancel each non-terminal task:

```bash
agent-team task cancel --task TASK_ID --reason "run aborted"
```

3. Cancel the run:

```bash
agent-team run cancel --run RUN_ID --reason "reason for abort"
```

## Notes

- Do not close a run with unfinished tasks.
- Create all tasks before dispatching workers.
- Use `_workspace/RUN_ID/` for artifacts in worker prompts and task contracts.
- Use messages for contract changes instead of direct peer-to-peer updates.

## Completion

This recipe is complete when:

- every task is in a terminal state: `done`, `failed`, or `cancelled`
- `run close` or `run cancel` returned successfully
- the run is no longer in `open` state
