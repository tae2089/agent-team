---
name: agent-team-codex-harness
description: "Design, create, evolve, or audit Codex-native harnesses for Agent Team projects. Builds specialized .codex/agents, .agents/skills, an orchestrator skill, AGENTS.md pointers, validation, and an evolution loop. Use for 'set up codex harness', 'create codex automation', 'build specialist agents', harness updates, harness audits, sync, status, and follow-up refinement. Generated harnesses default orchestrated runtime execution to Agent Team runtime skills backed by daemonless `agent-team` state."
metadata:
  version: 1.0.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
    skills:
      - agent-team-shared
---

# Agent Team Codex Harness

Build and evolve a domain-specific Codex harness in the style of `revfactory/harness`: a durable specialist system made of Codex agent definitions, procedural skills, an orchestrator skill, a project pointer, validation, and continuous evolution.

Agent Team preserves the upstream harness product model while using Codex-native primitives:

- Codex uses `.codex/agents/{name}.toml`, `.agents/skills/{name}/SKILL.md`, and `AGENTS.md`.
- Codex skills trigger through name/description metadata and load their `SKILL.md` body only when needed.
- Codex agent definitions are reusable role files. The orchestrator routes work to the right role and may delegate to subagents only when the active Codex environment and user request allow it.
- Codex has no peer-to-peer coordination bus. Coordination is explicit: orchestrator prompts, artifact files, final summaries, and Agent Team runtime skills backed by `agent-team` state.
- Generated harnesses default orchestrated runtime execution to the daemonless `agent-team` run/task/message/inbox/sync model; harness setup, editing, and audit are local filesystem tasks.

## Core Principles

1. Create specialist agent definitions and skills as first-class files. The harness value comes from separating who does the work from how the work is done.
2. Use a specialist roster plus an orchestrator as the default design for multi-specialist work. Codex-native coordination is orchestrator-led routing, not simulated peer-to-peer coordination.
3. Register only a concise harness pointer in `AGENTS.md`: trigger rule, orchestrator path, artifact root, runtime state pointer, and change history.
4. Treat the harness as an evolving system. After execution or feedback, update the relevant agent, skill, orchestrator, reference, and `AGENTS.md` history.
5. Keep construction separate from runtime. Do not call `agent-team` while building, editing, or auditing a harness; use Agent Team runtime skills for orchestrated execution after the harness is run.
6. Preserve upstream architecture patterns as orchestration patterns: pipeline, fan-out/fan-in, expert pool, producer-reviewer, supervisor, and hierarchical.
7. Runtime workers complete through `agent-team task complete` with concrete evidence and an artifact path. Blocked workers record a concrete blocked reason.

## Reference Loading

Read only the references needed for the current request:

- `references/agent-design-patterns.md`: Codex-native execution modes, pattern selection, separation criteria, and orchestration mapping.
- `references/team-examples.md`: concrete specialist roster shapes and output patterns.
- `references/orchestrator-template.md`: orchestrator templates for direct, delegated, and hybrid modes.
- `references/agent-team-runtime-protocol.md`: runtime state commands, completion, blocked behavior, retry, and sync check rules.
- `references/qa-agent-guide.md`: QA/reviewer design and integration-coherence verification.
- `references/skill-writing-guide.md`: skill frontmatter, trigger descriptions, progressive disclosure, examples, bundled resources.
- `references/skill-testing-guide.md`: static checks, trigger tests, with-skill vs baseline, dry-run, assertion evaluation, iteration loop.
- `references/schemas/models.md` and `references/schemas/agent-worker.template.toml`: required before writing `.codex/agents/*.toml`.

## Workflow

### Phase 0: Harness Audit

When this skill triggers, inspect the existing local harness first:

1. Read `PROJECT/.codex/agents/`, `PROJECT/.agents/skills/`, and `PROJECT/AGENTS.md`.
2. Inspect existing `_workspace/` artifact/report directories when relevant to follow-up work.
3. Classify the request:
   - **New build**: no relevant harness exists, or the user asks to create a new domain team.
   - **Extension**: an existing harness exists and the user asks to add/change agents, skills, workflow stages, QA, or triggers.
   - **Operations/maintenance**: audit, sync, status, cleanup, drift detection, or quality improvement of an existing harness.
   - **Runtime execution**: the user asks to run an orchestrated harness workflow, asks for follow-up work on a previous harness result, provides an advanced/debug `RUN_ID`/`TASK_ID`, or asks for tracked/stateful work. Runtime execution is skill-first: load the orchestrator/runtime/helper skills, then use the daemonless `agent-team` CLI only through those skill contracts.
   - **Direct answer**: the request is small enough to answer without a harness update.
4. Compare actual files with `AGENTS.md` and the orchestrator skill. Report drift before editing.

Existing extension phase matrix:

| Change Type                   | Phase 1                 | Phase 2        | Phase 3                 | Phase 4                   | Phase 5                             | Phase 6  |
| ----------------------------- | ----------------------- | -------------- | ----------------------- | ------------------------- | ----------------------------------- | -------- |
| Add agent                     | Skip, use Phase 0 facts | placement only | required                | if dedicated skill needed | update orchestrator                 | required |
| Add/update skill              | skip                    | skip           | skip                    | required                  | if wiring changes                   | required |
| Architecture change           | skip                    | required       | affected agents only    | affected skills only      | required                            | required |
| Runtime-state contract change | skip                    | placement only | affected runtime agents | affected skills only      | update orchestrator and `AGENTS.md` | required |

For operations/maintenance requests, jump to Phase 7-5 after Phase 0.

### Phase 1: Domain Analysis

Extract:

1. Domain, project, audience, and business goal.
2. Core work types: generation, validation, editing, implementation, research, analysis, operation, migration, release, etc.
3. Existing harness conflicts or duplication from Phase 0.
4. Codebase/content structure: stack, data model, main modules, docs, workflows, external systems.
5. Required specialist perspectives and quality gates.
6. Whether the work is single-worker, small-team, or a state-backed orchestrated workflow.
7. User constraints: speed vs rigor, safety, permissions, acceptable automation, output format.
8. User proficiency signals from the conversation. Adjust communication style without weakening technical correctness.

If domain, deliverable, acceptance criteria, or safety boundary cannot be inferred from files and user context, ask before designing the team.

### Phase 2: Specialist Architecture Design

#### 2-1. Execution Mode Selection

The upstream harness defaults to a coordinated specialist system. Agent Team Codex keeps that product model using Codex-native coordination.

| Mode                      | When To Use                                                                                                       | Codex/Agent Team Shape                                                                                                         |
| ------------------------- | ----------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| **Direct orchestrator**   | One specialist view is enough, the task is small, or the user did not ask for delegation                          | The orchestrator handles the work using the relevant skill and records artifacts/state when needed                          |
| **Delegated specialists** | The user asks for agents/delegation/parallel work, or the runtime explicitly supports delegation for this harness | The orchestrator delegates bounded, independent work to Codex agents, then integrates returned summaries and artifacts      |
| **Hybrid**                | Different phases need different coordination shapes                                                               | Direct orchestration for tightly coupled phases; delegated specialists for independent research, build lanes, review, or QA |

Decision order:

1. Start with direct orchestration. This is Codex-native and avoids unnecessary delegation.
2. Add delegated specialists only when they materially improve quality, parallelism, isolation, or review, and when the active Codex environment permits delegation.
3. Use hybrid mode when phase boundaries are clear, for example parallel research followed by direct synthesis, or direct build followed by delegated QA.
4. Follow the active Codex tool policy. Only orchestrators, supervisors, or explicit hierarchical leads may delegate. Ordinary workers do not spawn subagents.

#### 2-2. Architecture Pattern Selection

Split work into specialist domains and choose one or more patterns:

- **Pipeline**: sequential dependencies.
- **Fan-out/fan-in**: independent parallel research or generation before synthesis.
- **Expert pool**: route to one specialist based on classification.
- **Producer-reviewer**: generate, review, then revise or pass.
- **Supervisor**: central lead dynamically batches, assigns, verifies, and advances tasks.
- **Hierarchical**: explicit leads split work for subteams or lanes.
- **Handoff**: ownership moves to another workspace or specialist.

#### 2-3. Agent Separation Criteria

Split a specialist only when at least one axis is strong:

- distinct expertise or vocabulary
- distinct permission/sandbox requirements
- reusable process knowledge deserving a skill
- independent parallel work that reduces wall-clock time
- meaningful quality gate or adversarial review
- context isolation improves focus or reduces risk

Do not split just to make the harness look larger. Prefer 2-4 focused specialists unless the workflow clearly needs more.

### Phase 3: Agent Definition Creation

Every generated agent must be written as `PROJECT/.codex/agents/{name}.toml`. Do not rely on ad hoc role prompts alone. The file makes the role reusable, auditable, and discoverable in future sessions.

Use `references/schemas/agent-worker.template.toml` and fill every placeholder.

Required fields:

- `name`
- `description`
- `developer_instructions`
- `model`
- `sandbox_mode`
- `model_reasoning_effort`

Model and reasoning rules:

- Use `references/schemas/models.md` as the single source of truth.
- Do not hard-code stale model ids outside that registry.
- Choose reasoning effort by responsibility: low for simple helpers, medium for research/docs, high for code/review/QA, xhigh for orchestrator/architect roles where supported.

Agent instructions must include:

- core role and non-goals
- input protocol from the orchestrator
- expected output path and format
- artifact/evidence requirements
- runtime reporting behavior when assigned a Agent Team task
- blocked behavior and escalation path
- collaboration boundaries
- previous-artifact behavior for follow-up or partial rerun

QA/reviewer agents:

- QA checks integration contracts, not just existence.
- QA should run incrementally after meaningful modules or artifacts, not only once at the end.
- See `references/qa-agent-guide.md` for cross-boundary verification patterns.

### Phase 4: Skill Creation

Create procedure skills in `PROJECT/.agents/skills/{name}/SKILL.md` only when they add reusable process knowledge. A specialist can have no dedicated skill if its behavior is fully captured in the agent definition.

#### 4-1. Skill Structure

```text
skill-name/
  SKILL.md              # required: YAML frontmatter + Markdown body
  references/           # optional: conditional reference docs
  scripts/              # optional: deterministic/repeated helpers
  assets/               # optional: templates, static assets, examples
```

#### 4-2. Trigger-Rich Description

The description is the primary trigger signal. It must include:

- what the skill does
- concrete user phrases and task types
- follow-up phrases: update, rerun, refine, audit, review, partial rerun
- near-miss boundaries where another skill or direct answer is better
- for worker-facing skills, a routing boundary that says the skill is normally used when routed by the orchestrator

Weak:

```yaml
description: "Processes documents."
```

Strong:

```yaml
description: "Reviews API documentation for accuracy, missing setup steps, broken commands, and stale references. Use when the user asks to update docs, audit docs, validate examples, rerun a docs review, or improve a previous documentation artifact."
```

Worker-facing strong:

```yaml
description: "Worker procedure for reviewing API documentation when routed by the docs orchestrator. Use for docs QA, command verification, partial review rerun, and review-driven fixes. For full docs workflows, route through the orchestrator skill."
```

#### 4-3. Body Writing Rules

| Principle              | Rule                                                                                 |
| ---------------------- | ------------------------------------------------------------------------------------ |
| Explain why            | Prefer reasoned rules over unexplained ALWAYS/NEVER commands.                        |
| Keep it lean           | Keep `SKILL.md` focused. Move long examples/tables to `references/`.                 |
| Generalize             | Convert feedback into general principles, not one-off patches.                       |
| Bundle repeated code   | Put deterministic helper scripts in `scripts/` when repeated across tests.           |
| Use commands carefully | Commands must match the repo and platform. Do not invent unsupported Codex features. |

#### 4-4. Progressive Disclosure

Use the same progressive disclosure model as upstream:

| Level         | Load Time        | Size Target                                 |
| ------------- | ---------------- | ------------------------------------------- |
| Metadata      | always available | short name + trigger-rich description       |
| `SKILL.md`    | when triggered   | concise main workflow                       |
| `references/` | only when needed | detailed examples, schemas, domain variants |

If a reference grows past roughly 300 lines, add a short table of contents. If a domain has variants, split references by variant and load only the relevant file.

#### 4-5. Skill-Agent Wiring

- One agent may use zero, one, or many skills.
- One skill may be shared by multiple agents.
- Skills describe how work is done; agents describe who owns which judgment and output.
- The orchestrator owns when each skill/agent is used.

### Phase 5: Integration and Orchestration

The orchestrator is a skill that coordinates the full harness. It connects specialists, skills, artifacts, runtime state, retries, follow-ups, and final reporting.

For existing harness extension, update the existing orchestrator instead of creating a parallel one unless the user explicitly wants a new orchestrator. Add new agent trigger keywords to the orchestrator description when needed.

#### 5-0. Orchestrator Patterns

**Direct orchestrator pattern:**

```text
[orchestrator]
  - classifies request
  - selects pattern and active specialists
  - performs tightly coupled work directly using relevant skills/references
  - creates or resumes Agent Team runtime state for orchestrated harness execution
  - writes artifacts and final summaries
```

**Delegated specialist pattern:**

```text
[orchestrator]
  - decomposes independent work into bounded tasks
  - delegates through Codex-supported subagent calls only when allowed
  - passes context through prompts and artifact paths
  - collects evidence/artifacts/messages
  - integrates results and verifies evidence before advancing
```

**Single-worker pattern:**

```text
[orchestrator]
  - invokes one focused specialist only when useful, otherwise handles directly
  - collects the returned result and/or artifact
  - records one durable task result when runtime execution is active
```

**Hybrid pattern:**

- Parallel collection through delegated specialists, then direct synthesis.
- Direct build phase, then independent QA.
- Phase-specific routing where each phase declares its execution mode.

Generated orchestrators must explicitly state the execution mode for each phase: direct, delegated, or hybrid.

#### 5-1. Data Transfer Protocol

| Strategy           | Agent Team/Codex Mechanism               | Use When                                           |
| ------------------ | ------------------------------------- | -------------------------------------------------- |
| Artifact-based     | `_workspace/{run_id}/...` files       | large, structured, auditable outputs               |
| Return summary     | worker final response to orchestrator | compact result aggregation                         |
| Runtime task state | `agent-team task`                  | durable status, evidence, artifact source of truth |
| Runtime messages   | `agent-team message`               | compact progress, warnings, discoveries            |
| Runtime sync       | `agent-team sync`                  | inbox, dependency, and state-version checks        |

Recommended combination for durable workflows: runtime task state for official progress/result, `_workspace/` for artifacts, compact summaries for communication.

File naming convention:

- `_workspace/{run_id}/00_input/...`
- `_workspace/{run_id}/{phase}_{agent}_{artifact}.{ext}`
- `_workspace/{run_id}/{task_id}_result.md`
- final user-requested output path when applicable

`_workspace/` is not the task board or sync check store. It stores artifacts, reports, logs, generated outputs, and inputs.

#### 5-2. Error Handling

Include error handling in the orchestrator:

- Retry failed workers at most two times after the initial failure.
- Change the prompt or scope on retry; do not blindly repeat the same failing call.
- If retry budget is exhausted, mark the task blocked with a concrete reason.
- Do not advance past missing evidence.
- Preserve conflicting findings with source attribution instead of deleting inconvenient results.
- Report missing or skipped branches in the final output.

#### 5-3. Roster Size Guidelines

| Work Size           | Recommended Specialists | Work Per Specialist |
| ------------------- | ----------------------- | ------------------- |
| Small, 5-10 tasks   | 2-3                     | 3-5 tasks           |
| Medium, 10-20 tasks | 3-5                     | 4-6 tasks           |
| Large, 20+ tasks    | 5-7 with clear lanes    | 4-5 tasks           |

More agents increase coordination cost. A focused roster of three often outperforms a diffuse roster of five.

#### 5-4. `AGENTS.md` Harness Pointer

After building or materially changing a harness, update `AGENTS.md` with a concise pointer. Do not duplicate the full operating manual.

```markdown
## Harness: {domain_name}

**Goal:** {one-line harness goal}

**Trigger:** Use `.agents/skills/{orchestrator-skill-name}/SKILL.md` for {domain} work. Simple questions may be answered directly.

**Model:** Follows the `revfactory/harness` orchestrator/specialist structure adapted to Codex-native skills, agents, artifacts, and Agent Team runtime.

**Orchestrator:** `.agents/skills/{orchestrator-skill-name}/SKILL.md`
**Agents:** `.codex/agents/`
**Artifacts:** `_workspace/{run_id}/` when runtime execution is active; use a project-local `_workspace/{domain_name}/` root only for harness setup or audit work with no runtime run.

**Runtime State:**

- Load `agent-team-shared` first for global runtime rules.
- Use recipe/persona skills for workflow shape: `recipe-agent-team-terminology-context` for shared vocabulary artifacts, `recipe-agent-team-planning-grill` for pre-execution plan hardening, `recipe-agent-team-architecture-design` for backend design before coding, `persona-agent-team-designer` for visual/UI/icon/character/environment/logo/design-system work, `recipe-agent-team-compound-learning` for reusable learning capture, `recipe-agent-team-run-lifecycle` for full runs, `recipe-agent-team-worker-checkpoint` for worker checkpoints, and `recipe-agent-team-operational-audit` for audit/status/cleanup.
- Use service skills for navigation: `agent-team-run`, `agent-team-task`, `agent-team-inbox`, `agent-team-sync`, and `agent-team-ops`.
- Use exact command helper skills for command syntax and flags, for example `agent-team-task-complete`, `agent-team-sync-check`, `agent-team-message-send`, or `agent-team-event-log`.
- `RUN_ID` and `TASK_ID` are orchestrator-owned internal context, not required user input.
- If the user provides an advanced/debug `RUN_ID` or `TASK_ID`, inspect and resume that run/task before creating new state.
- If no runtime context is available, the orchestrator loads `agent-team-run-create` and `agent-team-task-create`, then creates one generated-ID run and generated-ID task records through those helper contracts.
- Do not use runtime state during harness setup, editing, audit-only work, simple one-shot answers, or explicitly local-only runs.
- Orchestrator owns run creation, task creation, evidence aggregation, inbox/sync checks, and artifact integration.
- Workers update only their assigned task.
- Completed tasks require evidence and an artifact path.
- Blocked tasks require a concrete blocked reason.
- `_workspace/` is for artifacts and reports only.

**Change History:**
| Date | Change | Target | Reason |
| --- | --- | --- | --- |
| {YYYY-MM-DD} | Initial harness | all | - |
```

Do not put full agent lists, full skill manuals, or detailed runtime procedures in `AGENTS.md`. Those belong in `.codex/agents/`, `.agents/skills/`, and references.

#### 5-5. Follow-Up Support

The orchestrator must support follow-up work, not just initial creation.

1. Include follow-up trigger phrases in the orchestrator description:
   - rerun, retry, update, modify, refine, audit, review
   - only rerun `{part}`
   - based on previous result
   - improve the previous output
2. Add a context check phase:
   - active in-session runtime context exists: resume that run/task
   - advanced/debug `RUN_ID` or `TASK_ID` present: inspect and resume/check state before creating new work
   - `_workspace/` missing: initial run
   - `_workspace/` exists and user asks partial change: partial rerun
   - `_workspace/` exists and user provides a new input: new run, archive or namespace previous artifacts
   - multiple plausible open/recent runs exist: ask the user to choose by title, status, and artifact summary, not by raw ID
   - orchestrated harness run with no state id: load the runtime recipe/helper skills and create a new Agent Team runtime execution through them
3. Add previous-artifact behavior to agent definitions:
   - read existing artifacts when provided
   - preserve unaffected sections
   - revise only the requested scope
   - cite what changed

### Phase 6: Validation and Testing

Validate the generated harness. See `references/skill-testing-guide.md` for detailed methodology.

#### 6-1. Structure Validation

- All agent TOML files exist in `.codex/agents/`.
- All required TOML fields are present.
- No placeholders remain.
- Model ids and reasoning values match `references/schemas/models.md`.
- Skill frontmatter has `name` and trigger-rich `description`.
- References named by skills exist.
- Orchestrator exists and has route, data flow, error handling, follow-up support, and test scenarios.
- Runtime instructions use `agent-team` and reserve `_workspace/` for artifacts.
- `AGENTS.md` pointer and change history match actual files.

#### 6-2. Execution-Mode Validation

- **Direct**: the orchestrator has enough context, references, and validation to complete without unnecessary delegation.
- **Delegated**: every specialist has a clear task, input, output, artifact path, ownership boundary, and aggregation point.
- **Hybrid**: each phase declares its mode and phase boundaries do not break data transfer.
- **Durable runtime**: active tasks, dependencies, sync checks, and blocked behavior are defined.

#### 6-3. Skill Execution Tests

For generated skills, when feasible:

1. Write 2-3 realistic test prompts per skill.
2. Compare with-skill vs baseline behavior when useful.
3. Evaluate output quality qualitatively and, where possible, with assertions.
4. Generalize improvements into the skill rather than patching one test case.
5. Bundle repeated helper code into `scripts/`.

#### 6-4. Trigger Validation

For each skill:

- Write 8-10 should-trigger prompts with varied phrasing.
- Write 8-10 should-not-trigger near-miss prompts.
- Near-miss prompts should be realistically ambiguous, not obviously unrelated.
- Check collisions with existing skills and update descriptions when boundaries are weak.

#### 6-5. Orchestrator Dry Run

- Review phase order and routing logic.
- Verify every worker input is produced by a prior phase or user input.
- Verify every required output is consumed, stored, or reported.
- Verify retry, blocked, and partial-rerun branches are executable.
- Add `## Test Scenarios` to the orchestrator skill: at least one normal flow and one failure flow.

Run project tests when harness edits affect executable code. For documentation-only harness edits, static validation is enough unless the user asks for full verification.

### Phase 7: Harness Evolution

Harnesses are living systems.

#### 7-1. Collect Feedback

After a harness execution, leave room for feedback:

- output quality
- specialist roster
- workflow order
- missing reviewer or owner
- trigger misses or false positives
- artifact structure
- runtime-state friction

Do not force feedback, but capture it when present.

#### 7-2. Feedback Routing

| Feedback Type          | Update Target                               | Example                                     |
| ---------------------- | ------------------------------------------- | ------------------------------------------- |
| output quality         | relevant skill                              | "analysis is shallow" -> add depth criteria |
| role mismatch          | agent TOML                                  | add/remove/split/merge specialist           |
| workflow order         | orchestrator skill                          | move QA earlier                             |
| specialist composition | orchestrator + agents                       | merge overlapping specialists               |
| trigger miss           | skill description                           | add follow-up phrase                        |
| runtime confusion      | `AGENTS.md` + orchestrator + state protocol | clarify task completion or sync handling          |

#### 7-3. Change History

Record every material change in `AGENTS.md`:

```markdown
**Change History:**
| Date | Change | Target | Reason |
| --- | --- | --- | --- |
| 2026-05-09 | Initial harness | all | - |
| 2026-05-09 | Add QA reviewer | `.codex/agents/qa-reviewer.toml` | artifact verification was weak |
```

The history prevents regression and explains why the harness evolved.

#### 7-4. Evolution Triggers

Propose evolution when:

- the same feedback repeats
- an agent repeatedly fails in the same way
- users bypass the orchestrator manually
- runtime state is updated inconsistently
- artifact paths or result contracts confuse workers
- a workflow has grown beyond the original specialist design

#### 7-5. Operations and Maintenance Workflow

Use this branch for audit, sync, status, cleanup, or drift repair:

1. Audit `.codex/agents/`, `.agents/skills/`, orchestrator skill, `AGENTS.md`, and relevant artifacts.
2. Compare actual agent files with the orchestrator's team table.
3. Compare actual skill directories with orchestrator and agent references.
4. Report drift before editing.
5. Apply one coherent change at a time.
6. Update `AGENTS.md` change history.
7. Re-run affected validation checks.

## Output Checklist

Confirm before completion:

- [ ] `PROJECT/.codex/agents/` contains all required agent TOML files.
- [ ] `PROJECT/.agents/skills/` contains generated skills with valid `SKILL.md` frontmatter.
- [ ] One orchestrator skill exists for multi-agent harnesses.
- [ ] Execution mode is explicit: direct, delegated, or hybrid.
- [ ] Every generated agent has model, reasoning, sandbox, input/output, evidence, artifact, and blocked behavior.
- [ ] Ordinary workers do not spawn subagents.
- [ ] No unsupported non-Codex tools are referenced as executable commands.
- [ ] Existing agents/skills do not conflict with new ones.
- [ ] Skill descriptions are trigger-rich and include follow-up keywords.
- [ ] Worker-facing skill descriptions include an orchestrator routing boundary.
- [ ] `SKILL.md` bodies are lean; large details are in `references/`.
- [ ] Trigger tests and near-miss tests are written or intentionally deferred with a reason.
- [ ] Orchestrator dry-run has no missing data links.
- [ ] `AGENTS.md` contains the harness pointer and change history.
- [ ] Runtime state contract points to `agent-team-shared`, the relevant recipe/service skill, and exact command helper skills.
- [ ] Worker runtime instructions do not rely only on service-level skills for command syntax.
- [ ] Operational flows reference `agent-team-ops` or `recipe-agent-team-operational-audit`.
- [ ] Runtime examples use presence-based boolean flags in shell commands, not separate boolean words after the flag.
- [ ] `_workspace/` is used only for artifacts and reports.

## References

- Harness patterns: `references/agent-design-patterns.md`
- Specialist examples: `references/team-examples.md`
- Orchestrator template: `references/orchestrator-template.md`
- Skill writing guide: `references/skill-writing-guide.md`
- Skill testing guide: `references/skill-testing-guide.md`
- QA agent guide: `references/qa-agent-guide.md`
- Agent Team runtime protocol: `references/agent-team-runtime-protocol.md`
- Codex agent schemas: `references/schemas/`
