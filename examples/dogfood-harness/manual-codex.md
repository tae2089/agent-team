# Manual Codex Harness Dogfood

Use this when validating that `agent-team-codex-harness` works with the runtime contract in an interactive Codex environment.

## Setup

```bash
agent-team init
agent-team run create --id run_manual_codex --title "manual codex harness dogfood"
mkdir -p _workspace/run_manual_codex
```

## Assistant Prompt

Ask Codex to load `agent-team-codex-harness` and create a tiny docs harness with:

- an orchestrator
- a writer
- a reviewer
- artifacts under `_workspace/run_manual_codex/`
- runtime state tracked with `agent-team`

Use a small target task, for example:

```text
Use agent-team-codex-harness. Create and run a tiny docs harness dogfood.
The orchestrator should create agent-team tasks for a writer and reviewer,
the writer should produce _workspace/run_manual_codex/writer.md,
the reviewer should produce _workspace/run_manual_codex/review.md,
and the run should be closed only after run summary says it is ready.
```

For headless execution, use the reusable prompt:

```bash
codex exec --skip-git-repo-check -C "$TMP_PROJECT" --sandbox workspace-write -o "$OUT" - < examples/dogfood-harness/prompts/codex.txt
```

## Expected Runtime Shape

The orchestrator should use:

```bash
agent-team task create --run run_manual_codex --agent writer --title "write dogfood artifact"
agent-team task create --run run_manual_codex --agent reviewer --title "review dogfood artifact" --depends-on TASK_ID
agent-team message send --run run_manual_codex --from orchestrator --to writer --kind notice --body "..."
```

Workers should use:

```bash
agent-team inbox list --agent writer --run run_manual_codex --unread
agent-team inbox ack --msg MSG_ID --agent writer
agent-team sync check --agent writer --run run_manual_codex --task TASK_ID
agent-team task start --task TASK_ID --agent writer
agent-team task complete --task TASK_ID --agent writer --evidence "..." --artifact "_workspace/run_manual_codex/writer.md"
```

The orchestrator should finish with:

```bash
agent-team run summary --run run_manual_codex
agent-team event log --run run_manual_codex --limit 100
agent-team run close --run run_manual_codex --reason "manual codex harness dogfood complete"
```

## Acceptance Criteria

- `run summary` reports `close_ready:true` before close.
- `event log` includes task creation, message acknowledgement, task completion, and run close events.
- Writer and reviewer artifacts exist under `_workspace/run_manual_codex/`.
- No daemon or direct peer-to-peer agent channel is required.

Example outputs, when available, are stored under `examples/dogfood-harness/results/codex/`.
