# Gemini Specialist Design Patterns

Use this reference to choose the specialist roster, execution mode, and tool permissions for a Agent Team Gemini harness. It mirrors the `revfactory/harness` pattern catalog while using Gemini-native primitives.

## Gemini Execution Model

Gemini-native harnesses are built from:

- `.gemini/agents/{name}.md`: reusable role definitions with explicit frontmatter tools
- `.gemini/skills/{name}/SKILL.md`: reusable process knowledge
- `GEMINI.md`: project-level pointers and compact operating rules
- `_workspace/{run_id}/`: artifacts, reports, inputs, generated outputs
- Agent Team runtime skills backed by daemonless `agent-team` state for orchestrated harness execution

The orchestrator skill is the router and integrator. It may invoke Gemini agents only when its frontmatter includes `invoke_agent`; otherwise it performs direct orchestration and uses the specialist roster as role guidance.

## Execution Modes

| Mode | Use When | Shape |
| --- | --- | --- |
| `direct` | One specialist view is enough, the task is tightly coupled, or invocation is unnecessary | Orchestrator works directly using relevant skills/references and writes artifacts |
| `invoked` | Independent tasks benefit from isolation, review, or parallel work, and the orchestrator has `invoke_agent` | Orchestrator gives each specialist a bounded prompt, artifact path, and output contract, then integrates results |
| `hybrid` | Some phases are tightly coupled and others are independent | Declare the mode per phase and pass data through artifacts/summaries |

Decision order:

1. Start with `direct`. This is the safest Gemini-native default.
2. Use `invoked` only when parallelism, isolation, specialist judgment, or independent review materially improves quality.
3. Use `hybrid` when a workflow has clear phase boundaries.
4. Ordinary workers never include `invoke_agent`. Only orchestrators, supervisors, or explicit hierarchical leads may use it.

## Core Patterns

| Pattern | Use When | Gemini-Native Runtime Shape |
| --- | --- | --- |
| `pipeline` | Steps are sequential | Orchestrator completes or invokes one step, verifies artifact/evidence, then advances |
| `fan_out_fan_in` | Specialists can work independently before synthesis | Orchestrator runs independent lanes directly or via invoked workers, then synthesizes |
| `expert_pool` | One specialist should be selected by classification | Orchestrator records routing rationale and uses only the selected specialist |
| `producer_reviewer` | Output quality improves through review | Producer artifact feeds reviewer; reviewer returns PASS/FIX/BLOCKED |
| `supervisor` | Many tasks need batching or dynamic assignment | Supervisor/orchestrator owns task creation, evidence verification, and sync check updates |
| `hierarchical` | A lead must split work across lanes | Only explicit leads may include `invoke_agent`; write ownership boundaries before invoking work |
| `message_coordination` | Ownership or contract information must move between specialists | Use `agent-team message send` with a concrete reason and artifact context |

## Separation Criteria

Split a specialist only when it has a real reason to exist:

- distinct expertise, vocabulary, or judgment
- distinct tool/permission requirements
- reusable process knowledge deserving a skill
- independent parallel work
- quality gate or adversarial review
- context isolation that reduces risk

Do not split merely to make the harness look larger. Prefer a small roster with clear boundaries.

## Gemini Tool Guide

| Agent Type | Tools |
| --- | --- |
| direct-only researcher / analyst | `ask_user`, `activate_skill` |
| planner / document producer | `ask_user`, `activate_skill`; add local file tools required by the project |
| coder / docs writer | `ask_user`, `activate_skill`, project file tools; add `run_shell_command` only when commands/tests are needed |
| reviewer / QA | `ask_user`, `activate_skill`, project file tools, `run_shell_command` when validation commands are needed |
| runtime worker | `ask_user`, `activate_skill`, `run_shell_command` |
| orchestrator / supervisor / hierarchical lead | `ask_user`, `activate_skill`, `invoke_agent`; add `run_shell_command` only for durable runtime commands |

Never use wildcard tool permissions such as `tools: ["*"]`.

Workers should not include `invoke_agent`. Add `run_shell_command` only when the role must run tests or execute Agent Team runtime commands. Add file tools only when the role writes artifacts or edits project files. Human-facing guidance lives in `docs/harness-sandbox-permissions.md`.

Runtime state rules live in `GEMINI.md`; orchestrators activate `agent-team-shared`, choose the relevant recipe/service skill, and activate exact command helper skills for orchestrated harness execution. `RUN_ID` and `TASK_ID` are orchestrator-owned internal context, not required user input. Orchestrators resolve active context, accept explicit IDs only as an advanced/debug escape hatch, or create a generated-ID run and tasks when no context exists. Workers update only the state task assigned by the orchestrator.

## Roster Size

| Work Size | Recommended Roster |
| --- | --- |
| Narrow task | 1 specialist or direct orchestrator |
| Small workflow | 2-4 specialists |
| Large workflow | More than 4 only with independent lanes and a real aggregation step |

Every extra specialist increases routing, tool-permission, artifact, and validation cost.
