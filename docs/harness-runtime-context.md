# Harness Runtime Context

Harness users should not need to provide `RUN_ID` or `TASK_ID`.

Those IDs are internal Agent Team runtime handles. The generated orchestrator owns them, captures them from JSON command output, and passes them to workers as task context. Workers use the IDs only after the orchestrator assigns a durable runtime task.

## User-Facing Flow

Users can ask for work in natural language:

```text
Run the docs harness on this repo and produce a review.
```

The orchestrator then:

1. Loads or activates the generated orchestrator skill.
2. Loads or activates Agent Team runtime recipe/service/helper skills.
3. Resolves active context or decides that a new runtime run is needed.
4. Runs `agent-team run create` without `--id`; the orchestrator must capture `.data.run.id`.
5. Runs `agent-team task create` without `--id`; the orchestrator must capture `.data.task.id`.
6. Passes `RUN_ID`, `TASK_ID`, `AGENT`, and `ARTIFACT_ROOT` to assigned workers.
7. Reports artifacts and unresolved risks to the user.

The final response may include a compact tracking line such as `run_id=...` for debugging, but it is not a required input for normal use.

## Follow-Up Resolution

When a user asks for a follow-up without IDs, the orchestrator resolves context in this order:

1. Active in-session runtime context from the current conversation.
2. Advanced/debug `RUN_ID` or `TASK_ID` if the user explicitly provides one.
3. Recent open runs and previous `_workspace/` artifacts.
4. A user choice between ambiguous candidates, shown by title, status, and artifact summary.
5. A new generated-ID run when no existing context matches.

Do not ask ordinary users to paste raw runtime IDs just to continue a harness workflow.

## Worker Boundary

Workers do not create, infer, or ask for runtime IDs. If a worker lacks `RUN_ID`, `TASK_ID`, or `AGENT`, it should report the missing context to the orchestrator instead of creating state on its own.
