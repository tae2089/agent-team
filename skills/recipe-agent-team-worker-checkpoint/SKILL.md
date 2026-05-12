---
name: recipe-agent-team-worker-checkpoint
description: "Recipe: Worker checkpoint flow for assigned agent-team tasks — sync, inbox, start/continue, artifact write, complete/block/fail. Use when worker has RUN_ID, TASK_ID, AGENT from orchestrator. Do not use for planning, architecture design, terminology, compound learning, orchestrator management, or unassigned work."
metadata:
  version: 1.1.0
  openclaw:
    category: "recipe"
    domain: "agent-orchestration"
  requires:
    bins:
      - agent-team
    skills:
      - agent-team-shared # shared constants and error codes
      - agent-team-inbox
      - agent-team-sync
      - agent-team-task
---

# Agent Team Worker Checkpoint

Use this recipe when a worker is about to start, checkpoint, or complete assigned work.

## Boundary

Use this recipe only after the orchestrator has provided `RUN_ID`, `TASK_ID`, `AGENT`, artifact root (normally `_workspace/RUN_ID/`), and the recipient name for orchestrator messages. Workers do not create runs, infer task IDs, or choose their own orchestrator recipient.

`task retry` is allowed by the CLI for `blocked`, `in_progress`, and `failed` tasks, but workers should not retry `failed` work without an explicit orchestrator approval message.

## Steps

1. Run sync check for the task. The result includes unread messages and dependency status.

```bash
agent-team sync check --agent AGENT --run RUN_ID --task TASK_ID
```

2. If `data.sync.blocking` is true, inspect `data.sync.issues` and act:

| Issue label               | Action                                                                                                                                                                                                                                                          |
| ------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `unread_messages`         | Read and ack each unread message from `data.sync.unread_messages`, then re-run sync check. If the sync output lacks the full message body or decision context, run `agent-team inbox list --agent AGENT --run RUN_ID --unread` before acking.                    |
| `incomplete_dependencies` | Send a `question` message to the orchestrator via `agent-team message send`. Do not poll. Wait for an inbox reply before re-running sync check.                                                                                                                 |
| `contract_changed`        | Run `agent-team task show --task TASK_ID` to read the updated task contract. Ack the message, then re-run sync check. If the change alters scope, deliverables, or dependencies significantly, send a `question` message to the orchestrator before proceeding. |
| Other / unknown           | Send a `question` message to the orchestrator. Do not proceed.                                                                                                                                                                                                  |

Ack each message individually after handling it:

```bash
# Repeat for each msg_id in data.sync.unread_messages
agent-team inbox ack --msg MSG_ID --agent AGENT
```

Send a question to the orchestrator when blocked or uncertain:

```bash
agent-team message send --run RUN_ID --from AGENT --to ORCHESTRATOR --kind question --body "REASON" --task TASK_ID
```

`ORCHESTRATOR` must come from assigned worker context or task metadata. If it is missing, block with a concrete reason instead of guessing a recipient.

Orchestrator approval for `--force` arrives as an inbox message. Check `data.sync.unread_messages` for a message with `kind: approval` before using `--force`. Without such a message, call `task block` instead.

3. If sync is clear (`data.sync.blocking: false`), start or continue work.

```bash
agent-team task start --task TASK_ID --agent AGENT
```

If `task start` returns `invalid_task_state`, inspect the current task state with `agent-team task show --task TASK_ID`:

| Current state | Action                                                                                                                                                             |
| ------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `in_progress` | Skip `task start` and continue work.                                                                                                                               |
| `done`        | Task already complete. Do not re-complete or overwrite.                                                                                                            |
| `blocked`     | Do not retry by default. Retry only when the orchestrator explicitly approved another attempt in an unread message or the worker prompt. After one retry, re-run sync check before starting. If sync check is still blocking, call `task block` again; do not loop. |
| `failed`      | Task was previously failed. Send a `question` message to the orchestrator. Do not retry unilaterally even though the CLI supports retry from `failed`. Wait for an `approval` inbox reply before proceeding. |
| `cancelled`   | Task was cancelled. Stop work and send a `question` message to the orchestrator. Do not proceed.                                                                   |

```bash
agent-team task retry --task TASK_ID --reason "REASON"
```

4. Perform the assigned work as defined in the task contract. Write outputs to `_workspace/RUN_ID/`.

5. Before calling complete, run one pre-completion checkpoint.

Pre-completion checkpoint:

1. Run `agent-team sync check --agent AGENT --run RUN_ID --task TASK_ID`.
2. If blocking is caused by unread messages, read the messages, handle and ack them, then run sync check once more. If the sync output lacks the full message body, run `agent-team inbox list --agent AGENT --run RUN_ID --unread`.
3. If blocking is caused by incomplete dependencies or an unknown issue, send a `question` message or call `task block`.
4. Do not start a new outer retry loop. If sync remains blocking after the one follow-up check, call `task block` unless explicit orchestrator approval for `--force` is present in unread messages.

6. Write the artifact to `_workspace/RUN_ID/` before calling complete. Do not call complete without the artifact file in place.

If the task was retried and a previous artifact exists at `_workspace/RUN_ID/TASK_ID.md`, overwrite it only when the new output fully supersedes the prior attempt: same scope, same deliverables, no sections carried forward from the old file. Use a suffix (`_workspace/RUN_ID/TASK_ID-ATTEMPT.md`) and pass that path as `--artifact` if: the prior file was partial, the retry adds to rather than replaces prior content, or any section of the old output is still relied upon.

7. Complete with evidence and an artifact path, block with a concrete reason, or fail if work cannot succeed under any conditions.

```bash
agent-team task complete --task TASK_ID --agent AGENT --evidence "verified" --artifact "_workspace/RUN_ID/TASK_ID.md"
agent-team task block --task TASK_ID --agent AGENT --reason "specific blocker"
agent-team task fail --task TASK_ID --agent AGENT --reason "specific failure reason" [--artifact "_workspace/RUN_ID/TASK_ID_failure.md"]
```

Use `task block` when the blocker is external and may be resolved (missing input, pending dependency, awaiting decision). Use `task fail` when the work cannot succeed regardless of retry (unrecoverable error, invalid task contract, tool unavailability).

## Notes

- Evidence should say what was verified, not just that work is done.
- Artifact paths must be run-scoped: `_workspace/RUN_ID/`.
- `--force` requires inbox proof of orchestrator approval. Check unread messages before using it.

## Completion

This recipe is complete when:

- for `done`: pre-completion sync returned `data.sync.blocking: false`, or `--force` had explicit unread-message approval from the orchestrator
- for `done`: all relevant unread messages have been handled and acked
- for `done`: the artifact file exists under `_workspace/RUN_ID/` before `task complete` is called
- for `blocked`: the task transitioned with a concrete reason
- for `failed`: the task transitioned with a concrete reason, and any available diagnostic artifact path was recorded
- the task transitioned to `done`, `blocked`, or `failed` (note: `blocked` is non-terminal; `failed` is terminal for run closure but the orchestrator may still retry before closing the run; this recipe's responsibility ends at the transition)
