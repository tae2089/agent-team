# Status Policy

`agent-team` uses a small status model so daemonless workers can make progress without a central scheduler. Statuses are part of the CLI contract and should stay stable unless the command schema version changes.

## Run Statuses

- `open`: work can still be created, started, messaged, retried, or reassigned.
- `closed`: the orchestrator has reviewed the terminal task outcomes and finalized the run.
- `cancelled`: the run itself was aborted and should not be treated as successfully finalized.

`run close` is allowed only for `open` runs where every task is terminal. It may close a run that contains `failed` or `cancelled` tasks; `closed` means "finalized with known outcomes", not "all work succeeded".

`run cancel` is for aborting an active run. A `closed` run cannot be cancelled.

## Task Statuses

Non-terminal statuses:

- `pending`: assigned but not currently being worked.
- `in_progress`: claimed by the assigned agent.
- `blocked`: the assigned agent cannot continue without an external change.

Terminal statuses:

- `done`: completed with evidence and an artifact path.
- `cancelled`: no longer needed, usually by orchestrator decision.
- `failed`: attempted but could not be completed by the assigned agent.

## Transition Policy

- `task start` is allowed from `pending`.
- `task complete` is allowed from `pending`, `in_progress`, and `done`; unread relevant messages or incomplete dependencies block completion unless `--force` is used.
- `task block` is allowed from non-terminal statuses.
- `task fail` is allowed from non-terminal statuses and requires the assigned agent.
- `task cancel` is allowed from non-terminal statuses.
- `task retry` is allowed from `blocked`, `in_progress`, and `failed`; it resets the task to `pending`.
- `task reassign` is allowed from `pending` and `blocked`; blocked tasks reset to `pending` when reassigned.

Every run/task/message mutation emits an event. Retry history, cancellation, failure, run finalization, and message acknowledgements are audited through `event log`.

When a transition is rejected, inspect `error.code` and `error.recovery`. See [errors.md](errors.md) for code-specific recovery actions.
