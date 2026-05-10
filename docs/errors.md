# Error Recovery

`agent-team` commands always return JSON. On failure, branch on `error.code`; `error.message` is for humans and may change.

Failure responses include an additive `error.recovery` object with a short summary, actions, suggested commands, docs, and related skills.

## validation_error

Required or invalid input is missing.

Recovery:

- Check the command helper skill or `agent-team COMMAND --help`.
- Provide IDs and routing fields with named flags or `--params`.
- Put richer payloads such as bodies, metadata, evidence, artifacts, and reasons in `--json` when useful.

## input_conflict

The same value was supplied by both a named flag and JSON.

Recovery:

- Provide each field once.
- Prefer `--params` for IDs, filters, booleans, and routing fields.
- Prefer `--json` for body text, metadata, evidence, artifact paths, and reasons.

## invalid_json

JSON is malformed, contains unknown fields, or uses the wrong value type.

Recovery:

- Validate the JSON object before retrying.
- Use keys from `agent-team schema export` or the relevant helper skill.
- Use JSON objects for metadata and structured payloads.

## invalid_json_source

The JSON source argument is invalid.

Recovery:

- Use inline JSON, `@file`, or `-` for stdin.
- Ensure `@file` paths are non-empty and readable.
- Ensure stdin is available when using `-`.

## not_found

The requested run, task, or message does not exist in the active state database.

Recovery:

- List nearby runs, tasks, or messages to find the correct ID.
- Confirm `AGENT_TEAM_STATE_DIR` points at the expected workflow state.
- Run `agent-team init` when this is a new state directory.

Useful commands:

```bash
agent-team run list
agent-team task list --run RUN_ID
agent-team message list --run RUN_ID
```

## agent_mismatch

The command was executed by an agent that does not own the task or message.

Recovery:

- Inspect the task or message to find the assigned owner.
- Retry as the assigned agent when appropriate.
- Ask the orchestrator to reassign ownership when this is a real handoff.

## sync_conflict

Task completion is blocked by unread relevant messages or incomplete dependencies.

Recovery:

- Run `agent-team sync check --agent AGENT --run RUN_ID --task TASK_ID`.
- Read and acknowledge relevant inbox messages.
- Wait for or complete unfinished dependencies.
- Use `--force` only when the orchestrator explicitly approves it.

## run_not_ready

The run still has unfinished tasks.

Recovery:

- Run `agent-team run summary --run RUN_ID`.
- Complete, retry, reassign, cancel, fail, or block unfinished work as policy allows.
- Close the run only after summary reports it is ready to close.

## invalid_run_state

The current run status does not allow the requested operation.

Recovery:

- Inspect the run with `agent-team run status --run RUN_ID`.
- Do not cancel a closed run.
- Create or resume an open run for new work.

## invalid_task_state

The current task status does not allow the requested transition.

Recovery:

- Inspect the task with `agent-team task show --task TASK_ID`.
- Use retry for blocked, in-progress, or failed work.
- Use reassign for pending or blocked work.
- Do not mutate terminal tasks except through run finalization.

## internal_error

The CLI hit an unexpected local failure, commonly state directory, SQLite, filesystem, or install health.

Recovery:

- Confirm the state directory exists and is writable.
- Set `AGENT_TEAM_STATE_DIR` to an isolated writable directory when sandboxed.
- Run `agent-team init` for a new state directory.
- Check [install](install.md) and [harness sandbox permissions](harness-sandbox-permissions.md) when running through Codex or Gemini.
