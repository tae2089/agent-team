# Agent Team Runtime Protocol

Use this reference when writing runtime contracts or worker reporting instructions.

## Boundary

Harness setup, editing, and audit are local filesystem tasks. Orchestrated runtime execution is skill-first: load the generated orchestrator skill, load the relevant Agent Team runtime recipe/service/helper skills, then execute daemonless `agent-team` commands only through those helper contracts. The CLI coordinates state only; it does not run workers.

Artifacts and reports live under `_workspace/{run_id}/`. SQLite state lives in `.agent-team/agent-team.db` unless `AGENT_TEAM_STATE_DIR` is set. `RUN_ID` and `TASK_ID` are orchestrator-owned internal context, not required user input.

## Skill Loading

Generated harnesses must tell orchestrators and workers to load skills before any runtime command:

- `agent-team-shared` for global runtime rules, JSON conventions, state directory behavior, and errors.
- Service skills for navigation: `agent-team-run`, `agent-team-task`, `agent-team-inbox`, `agent-team-sync`, and `agent-team-ops`.
- Recipe skills for workflow shape: `recipe-agent-team-terminology-context`, `recipe-agent-team-planning-grill`, `recipe-agent-team-architecture-design`, `recipe-agent-team-compound-learning`, `recipe-agent-team-run-lifecycle`, `recipe-agent-team-worker-checkpoint`, and `recipe-agent-team-operational-audit`.
- Exact command helper skills for command syntax, flags, examples, and errors.

Load helper skills for exact command behavior:

| Runtime action | Helper skill |
|----------------|--------------|
| Create/resume/finalize run | `agent-team-run-create`, `agent-team-run-status`, `agent-team-run-summary`, `agent-team-run-close` |
| Create/list/show tasks | `agent-team-task-create`, `agent-team-task-list`, `agent-team-task-show` |
| Start assigned work | `agent-team-task-start` |
| Send progress or coordination | `agent-team-message-send` |
| Read and ack inbox | `agent-team-inbox-list`, `agent-team-inbox-ack` |
| Check drift before completion | `agent-team-sync-check` |
| Complete assigned work | `agent-team-task-complete` |
| Block assigned work | `agent-team-task-block` |
| Retry, reassign, cancel, or fail work | `agent-team-task-retry`, `agent-team-task-reassign`, `agent-team-task-cancel`, `agent-team-task-fail` |
| Detect stale work | `agent-team-task-stale` |
| Audit state | `agent-team-run-list`, `agent-team-message-list`, `agent-team-event-log`, `agent-team-schema-export` |

## Orchestrator Context Resolution

For a user-facing harness run, the user should not need to provide raw runtime IDs.

1. Reuse active in-session runtime context when the current request is a follow-up.
2. Accept an advanced/debug `RUN_ID` or `TASK_ID` only as an escape hatch.
3. Inspect recent open runs and previous `_workspace/` artifacts when the user refers to earlier work without IDs.
4. If multiple recent runs match, ask the user to choose by title, status, and artifact summary, not by raw ID.
5. If no context matches, load `agent-team-run-create`, create a run without `--id`, capture `.data.run.id`, then create tasks without `--id` and capture `.data.task.id`.

## Worker Context

```text
RUN_ID=...
TASK_ID=...
AGENT=...
ARTIFACT_ROOT=_workspace/{run_id}/
AGENT_TEAM_STATE_DIR=...        # optional
```

## Skill-Mediated Command Shapes

These command examples are not standalone operating instructions. They are the shell forms that the named helper skills authorize after the helper skill has been loaded. Workers receive the IDs below from the orchestrator; they do not infer, create, or ask the user for them.

Load `agent-team-task-start` before start:

```bash
agent-team task start --task TASK_ID --agent AGENT
```

Load `agent-team-message-send` before progress or coordination messages:

```bash
agent-team message send \
  --run RUN_ID --task TASK_ID --from AGENT --to supervisor --kind progress \
  --body "Short update."
```

Load `agent-team-inbox-list` and `agent-team-inbox-ack` before inbox checkpoints:

```bash
agent-team inbox list --agent AGENT --run RUN_ID --unread
agent-team inbox ack --msg MSG_ID --agent AGENT
```

Load `agent-team-sync-check` before sync check:

```bash
agent-team sync check --agent AGENT --run RUN_ID --task TASK_ID
```

Load `agent-team-task-complete` before completion:

```bash
agent-team task complete \
  --task TASK_ID --agent AGENT \
  --evidence "VERIFIABLE_EVIDENCE" \
  --artifact "_workspace/RUN_ID/TASK_ID_result.md"
```

Load `agent-team-task-block` before blocking:

```bash
agent-team task block --task TASK_ID --agent AGENT --reason "REASON"
```

Load `agent-team-run-summary` and `agent-team-run-close` before finalization:

```bash
agent-team run summary --run RUN_ID
agent-team run close --run RUN_ID --reason "All tasks reached terminal status and artifacts were integrated."
```

Use compact JSON only when it makes command construction safer:

```bash
agent-team task complete \
  --params '{"task_id":"TASK_ID","agent":"AGENT","force":false}' \
  --json '{"evidence":"VERIFIABLE_EVIDENCE","artifact":"_workspace/RUN_ID/TASK_ID_result.md"}'
```

## Rules

- Done tasks require evidence and an artifact path.
- Blocked tasks require a concrete reason.
- Workers update only their assigned task.
- Workers check inbox and sync before completion.
- Workers do not use `--force` unless the orchestrator explicitly approves it.
- On command failure, branch on `error.code`; use `error.recovery` for next steps and load the helper skills listed there when relevant.
- The orchestrator creates runs/tasks, verifies evidence, checks inbox/sync status, and integrates artifacts.
- The orchestrator owns runtime ID generation and captures generated IDs from JSON command output before assigning worker context.
- The orchestrator closes a run only after `agent-team run summary` reports it is ready to close.
- Operational flows should load `agent-team-ops` or `recipe-agent-team-operational-audit`.
- Do not tell workers to "just run the CLI"; tell them which runtime skill to load and which task context to use.
- Do not treat `_workspace/` as a task board; it is for artifacts and reports only.
