# Codex Orchestrator Skill Template

Use this reference when creating `.agents/skills/{orchestrator}/SKILL.md`.

## Required Shape

```markdown
---
name: {domain}-orchestrator
description: "{Domain} Codex harness orchestrator. Routes {domain} build, update, review, partial rerun, resume, audit, and refinement through Codex-native skills, specialist agents when delegation is allowed, artifacts, and agent-team runtime state."
---

# {Domain} Codex Orchestrator

## Route

| Request                                     | Action                                                          |
| ------------------------------------------- | --------------------------------------------------------------- |
| simple answer or file lookup                | answer directly                                                 |
| harness setup/edit/audit                    | inspect local harness files only                                |
| single focused task                         | use direct mode or one worker                                   |
| independent parallel work requested/allowed | use delegated mode                                              |
| orchestrated harness run                    | load runtime skills, resolve/create internal runtime context, and execute the plan |

## Specialist Roster

| Agent     | Role   | Model   | Reasoning | Sandbox   | Output     |
| --------- | ------ | ------- | --------- | --------- | ---------- |
| `{agent}` | {role} | {model} | {effort}  | {sandbox} | {artifact} |

## Execution Modes

- `direct`: the orchestrator performs tightly coupled work using relevant skills/references.
- `delegated`: the orchestrator delegates bounded independent work to Codex agents only when the environment and request allow it.
- `hybrid`: each phase declares whether it is direct or delegated.

## Runtime Contract

- Runtime coordination is skill-first: load runtime recipe/service/helper skills, then use daemonless `agent-team` commands only through those helper contracts.
- Load `agent-team-shared` first for global runtime behavior.
- Use recipes for workflow shape: `recipe-agent-team-terminology-context` for shared vocabulary artifacts, `recipe-agent-team-planning-grill` for pre-execution plan hardening, `recipe-agent-team-architecture-design` for design artifacts before coding, `recipe-agent-team-compound-learning` for reusable learning capture, `recipe-agent-team-run-lifecycle` for full runs, `recipe-agent-team-worker-checkpoint` for worker checkpoints, and `recipe-agent-team-operational-audit` for audit/status/cleanup.
- For exact command behavior, load the relevant helper skill, such as `agent-team-run-create`, `agent-team-task-create`, `agent-team-task-complete`, `agent-team-sync-check`, `agent-team-message-send`, or `agent-team-event-log`.
- Use service skills only for navigation: `agent-team-run`, `agent-team-task`, `agent-team-inbox`, `agent-team-sync`, and `agent-team-ops`.
- Setup, edit, and audit requests do not probe runtime state.
- `RUN_ID` and `TASK_ID` are orchestrator-owned internal context, not required user input.
- Resolve context in this order: active in-session context, advanced/debug user-provided IDs, recent open run plus previous artifacts, user choice among ambiguous recent runs, then a new generated-ID run.
- If no runtime context is available for an orchestrated run, load `agent-team-run-create` and `agent-team-task-create`, then create one run and task records without explicit IDs and capture the returned JSON IDs.
- Workers receive orchestrator-supplied `RUN_ID`, `TASK_ID`, `AGENT`, `ARTIFACT_ROOT`, and optional `AGENT_TEAM_STATE_DIR` only for assigned durable tasks.
- Workers update only their assigned task.
- The orchestrator creates tasks, verifies evidence, checks inbox/sync status, and integrates artifacts.
- `_workspace/{run_id}/` stores artifacts and reports only.

## Workflow

1. Classify the request.
2. Select pattern and execution mode.
3. Load only relevant references.
4. Load `recipe-agent-team-terminology-context` when user language, code terms, or agent task vocabulary need alignment.
5. Load `recipe-agent-team-planning-grill` when the plan, terminology, task boundaries, or acceptance criteria need hardening.
6. Load `recipe-agent-team-architecture-design` when the workflow needs design artifacts before coding.
7. Load `recipe-agent-team-run-lifecycle`, then resolve or create the orchestrator-owned runtime context for orchestrated harness execution.
8. Create active runtime tasks when needed using exact command helper skills, omit explicit task IDs by default, and capture returned task IDs.
9. Execute direct work or delegate bounded specialist tasks when allowed.
10. Collect returned summaries, task records, inbox messages, evidence, and artifacts.
11. Load `recipe-agent-team-compound-learning` when the workflow produced reusable learning for future runs.
12. Retry, block, request clarification, or complete.
13. Report final artifacts and unresolved risks.
```

## Data Flow

Every phase must declare:

- inputs: user prompt, files, previous artifacts, runtime task records
- active mode: direct, delegated, or hybrid
- active specialist(s)
- output artifact path
- evidence required
- next consumer
- failure path

## Completion Rules

Complete only when all required active work has concrete evidence and artifact paths. Do not silently advance past missing evidence.

## Follow-Up Support

Generated orchestrators must support:

- update / modify / refine
- partial rerun
- audit / status / sync
- review / QA
- resume of a durable execution
- natural-language follow-up without asking the user for raw `RUN_ID` or `TASK_ID`

## Test Scenarios

Include at least:

- normal flow: classify -> execute -> verify -> report
- failure flow: worker or phase fails -> retry with changed scope -> block with reason if unresolved
- follow-up flow: previous artifact exists -> partial rerun or refinement
