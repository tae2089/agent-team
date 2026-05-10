---
name: manual-dogfood-orchestrator
description: "Orchestrates a tiny manual Agent Team Codex dogfood workflow with writer and reviewer specialists. Use for manual dogfood harness runs, agent-team runtime verification, rerun, retry, update, modify, refine, audit, review, partial rerun, previous-result follow-up, and state inspection. Simple local-only questions should be answered directly instead of creating runtime state."
---

# Manual Dogfood Orchestrator

Use this skill to run a minimal producer-reviewer harness with Agent Team runtime skills, state, and artifact files.

## Execution Mode

Hybrid:

- Direct orchestration for run creation, task creation, notices, state checks, summary, close, and event-log inspection.
- Delegated specialist responsibilities for the writer and reviewer task bodies. In Codex sessions where subagent delegation is unavailable or unnecessary, the orchestrator may execute the worker procedures directly while preserving the runtime task contract.

## Runtime Loading

Before runtime commands, load:

- `agent-team-shared`
- `agent-team-run`, `agent-team-task`, `agent-team-inbox`, `agent-team-sync`, and `agent-team-ops`
- `recipe-agent-team-run-lifecycle`, `recipe-agent-team-worker-checkpoint`, and `recipe-agent-team-operational-audit`
- Exact helpers for commands being used, especially `agent-team-run-create`, `agent-team-task-create`, `agent-team-message-send`, `agent-team-inbox-list`, `agent-team-inbox-ack`, `agent-team-sync-check`, `agent-team-task-start`, `agent-team-task-complete`, `agent-team-run-summary`, `agent-team-run-close`, and `agent-team-event-log`

## Normal Flow

1. If needed, load `agent-team-shared` and initialize local state.
2. Resolve internal runtime context from active session state, an advanced/debug ID, recent open runs, or previous artifacts; if none matches, load `recipe-agent-team-run-lifecycle` and `agent-team-run-create`, then create a generated-ID run and capture `.data.run.id`.
3. Load `agent-team-task-create`, then create a generated-ID writer task assigned to `writer` and capture `.data.task.id`.
4. Load `agent-team-task-create`, then create a generated-ID reviewer task assigned to `reviewer` with `--depends-on` pointing at the writer task.
5. Load `agent-team-message-send`, then send a compact notice message from `orchestrator` to `writer`.
6. Worker checkpoint for writer:
   - load `agent-team-inbox-list`, then run `agent-team inbox list --agent writer --run RUN_ID --unread`
   - ack unread messages
   - load `agent-team-sync-check`, then run `agent-team sync check --agent writer --run RUN_ID --task WRITER_TASK_ID`
   - load `agent-team-task-start`, then run `agent-team task start --task WRITER_TASK_ID --agent writer`
   - write `_workspace/RUN_ID/writer.md`
   - sync check again
   - complete with concrete evidence and artifact path
7. Worker checkpoint for reviewer:
   - load `agent-team-sync-check`, then run `agent-team sync check --agent reviewer --run RUN_ID --task REVIEWER_TASK_ID`
   - load `agent-team-task-start`, then run `agent-team task start --task REVIEWER_TASK_ID --agent reviewer`
   - inspect `_workspace/RUN_ID/writer.md`
   - write `_workspace/RUN_ID/review.md`
   - sync check again
   - complete with concrete evidence and artifact path
8. Load `agent-team-run-summary`, then run `agent-team run summary --run RUN_ID` and verify `close_ready:true`.
9. Load `agent-team-run-close`, then close the run with `agent-team run close`.
10. Load `agent-team-event-log`, then inspect `agent-team event log --run RUN_ID`.

## Error Handling

- Do not advance past a missing writer artifact or missing completion evidence.
- Retry a failed worker at most two times after changing the scope or prompt.
- If retry budget is exhausted, block the task with a concrete reason.
- Preserve conflicting findings in the reviewer artifact instead of deleting them.
- Do not close the run unless summary reports close readiness.

## Artifacts

- Writer: `_workspace/{run_id}/writer.md`
- Reviewer: `_workspace/{run_id}/review.md`

## Follow-Up Behavior

`RUN_ID` and `TASK_ID` are internal runtime context owned by the orchestrator. If the user provides an advanced/debug ID, resume and inspect state before creating new work. Otherwise resolve by active session, recent open runs, and previous artifacts; ask the user to choose by title/status/artifact summary only when several candidates match. For partial reruns, preserve existing artifacts, create replacement task records, and write updated artifacts under the same run root only when the user asks to update that run.

## Test Scenarios

- Normal flow: create a run, create writer and reviewer tasks, send notice, complete both tasks with evidence and artifact paths, confirm `close_ready:true`, close, and inspect the event log.
- Failure flow: writer artifact is missing; reviewer does not start, writer is retried or blocked, and the run remains open until terminal task state makes closure valid.
