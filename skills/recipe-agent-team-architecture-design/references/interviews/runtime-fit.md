# Agent Team Runtime Fit Interview

Use when a candidate touches agent-team run, task, message, inbox, sync, event, artifact, harness, or skill boundaries.

Core lens: agent-team coordinates durable state and artifacts; it does not run workers or turn the database into a document store.

## Interview

1. Which concern belongs in SQLite state, and which belongs in `_workspace/{run_id}/` artifacts?
2. Does the design preserve compact task metadata?
3. Are large prose, plans, designs, reviews, and learnings stored as artifacts?
4. Does the design keep workers from creating or inferring runtime IDs?
5. Are messages used for coordination instead of peer-to-peer side channels?
6. Does sync checking still catch unread messages or dependency drift?
7. Does run close readiness remain clear and based on terminal task states?

## Failure Signals

- Planning or design prose is pushed into task metadata.
- Workers create their own runs or task IDs.
- `_workspace/` is treated as a task board instead of artifact storage.
- Runtime state and harness construction concerns are mixed.
- The design weakens evidence or artifact requirements.

## Output

Record findings in the `## Interview Findings` table of `architecture-candidates.md` or `technical-design.md`:

| Lens | Finding | Design Impact |
| ---- | ------- | ------------- |
| Runtime Fit | state vs artifact boundary / task metadata shape / worker/orchestrator responsibility split | sync or message impact / run close impact |
