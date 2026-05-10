# Codex Specialist Design Patterns

Use this reference to choose the specialist roster and orchestration shape for a Agent Team Codex harness. It mirrors the `revfactory/harness` pattern catalog while using Codex-native primitives.

## Codex Execution Model

Codex-native harnesses are built from:

- `.codex/agents/{name}.toml`: reusable role definitions
- `.agents/skills/{name}/SKILL.md`: reusable process knowledge
- `AGENTS.md`: project-level pointers and compact operating rules
- `_workspace/{plan}/`: artifacts, reports, inputs, generated outputs
- Agent Team runtime skills backed by daemonless `agent-team` state for orchestrated harness execution

Codex does not require or expose a peer-to-peer coordination bus. The orchestrator skill is the router and integrator. It may delegate to Codex agents only when the active Codex environment and user request allow delegation; otherwise it performs direct orchestration and uses the specialist roster as role guidance.

## Execution Modes

| Mode        | Use When                                                                                                        | Shape                                                                                                            |
| ----------- | --------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- |
| `direct`    | One specialist view is enough, the task is tightly coupled, or delegation is not explicitly available           | Orchestrator works directly using relevant skills/references and writes artifacts                                |
| `delegated` | Independent tasks can run in parallel or benefit from context isolation/review, and Codex delegation is allowed | Orchestrator gives each specialist a bounded prompt, artifact path, and output contract, then integrates results |
| `hybrid`    | Some phases are tightly coupled and others are independent                                                      | Declare the mode per phase and pass data through artifacts/summaries                                             |

Decision order:

1. Start with `direct`. This is the safest Codex-native default.
2. Use `delegated` only when parallelism, isolation, specialist judgment, or independent review materially improves quality.
3. Use `hybrid` when a workflow has clear phase boundaries.
4. Ordinary workers never spawn subagents. Only orchestrators, supervisors, or explicit hierarchical leads may delegate when the active environment permits it.

## Core Patterns

| Pattern             | Use When                                            | Codex-Native Runtime Shape                                                                |
| ------------------- | --------------------------------------------------- | ----------------------------------------------------------------------------------------- |
| `pipeline`          | Steps are sequential                                | Orchestrator completes or delegates one step, verifies artifact/evidence, then advances   |
| `fan_out_fan_in`    | Specialists can work independently before synthesis | Orchestrator runs independent lanes directly or via delegated workers, then synthesizes   |
| `expert_pool`       | One specialist should be selected by classification | Orchestrator records routing rationale and uses only the selected specialist              |
| `producer_reviewer` | Output quality improves through review              | Producer artifact feeds reviewer; reviewer returns PASS/FIX/BLOCKED                       |
| `supervisor`        | Many tasks need batching or dynamic assignment      | Supervisor/orchestrator owns task creation, evidence verification, and sync check updates |
| `hierarchical`      | A lead must split work across lanes                 | Only explicit leads may delegate; write ownership boundaries before spawning work         |
| `message_coordination` | Ownership or contract information must move between specialists | Use `agent-team message send` with a concrete reason and artifact context                           |

## Separation Criteria

Split a specialist only when it has a real reason to exist:

- distinct expertise, vocabulary, or judgment
- distinct permission/sandbox requirements
- reusable process knowledge deserving a skill
- independent parallel work
- quality gate or adversarial review
- context isolation that reduces risk

Do not split merely to make the harness look larger. Prefer a small roster with clear boundaries.

## Permission Guide

| Agent Type                  | Suggested `sandbox_mode`                           |
| --------------------------- | -------------------------------------------------- |
| researcher / analyst        | `read-only`                                        |
| planner / document producer | `workspace-write`                                  |
| coder / docs writer         | `workspace-write`                                  |
| reviewer / QA               | `workspace-write`                                  |
| operator / deployer         | `danger-full-access` only when explicitly required |

Use `workspace-write` for most generated workers because they write `_workspace/` artifacts. Use `danger-full-access` only when the user explicitly approves release, external path, or trusted local automation. Codex headless runs may block hidden harness directories such as `.codex/` and `.agents/` in restricted modes; for harness-generation dogfood, use a trusted temp workspace and keep `AGENT_TEAM_STATE_DIR` inside that workspace. Human-facing guidance lives in `docs/harness-sandbox-permissions.md`.

Runtime state rules live in `AGENTS.md`; orchestrators load `agent-team-shared`, choose the relevant recipe/service skill, and load exact command helper skills for orchestrated harness execution. `RUN_ID` and `TASK_ID` are orchestrator-owned internal context, not required user input. Orchestrators resolve active context, accept explicit IDs only as an advanced/debug escape hatch, or create a generated-ID run and tasks when no context exists. Workers update only the state task assigned by the orchestrator.

## Roster Size

| Work Size      | Recommended Roster                                                  |
| -------------- | ------------------------------------------------------------------- |
| Narrow task    | 1 specialist or direct orchestrator                                 |
| Small workflow | 2-4 specialists                                                     |
| Large workflow | More than 4 only with independent lanes and a real aggregation step |

Every extra specialist increases routing, artifact, and validation cost.
