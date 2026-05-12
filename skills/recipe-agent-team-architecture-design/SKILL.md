---
name: recipe-agent-team-architecture-design
description: "Recipe: Produce architecture design artifacts (modules, interfaces, task contracts, candidates) before coding. Use for 'design the architecture', 'propose a structure', 'evaluate candidates', 'write a technical design', 'break this into implementation tasks'. Do not use for vague goals, terminology alignment, implementation, code review, or post-run learning; use planning-grill, terminology-context, coding workflow, review workflow, or compound-learning instead."
metadata:
  version: 1.0.0
  openclaw:
    category: "recipe"
    domain: "agent-orchestration"
  requires:
    bins:
      - agent-team
    skills:
      - agent-team-shared
      - agent-team-run
      - agent-team-task
---

# Agent Team Architecture Design

Use this recipe to produce architecture design artifacts for downstream implementation.

Stop when the selected design has affected files or modules, proposed module and interface shape, migration strategy, implementation task contracts, acceptance criteria, risks, and open questions. Implementation belongs to a separate coding workflow that consumes these artifacts.

This recipe has two phases:

1. Candidate discovery: identify and compare one to three viable architecture directions.
2. Selected technical design: expand the chosen candidate into a coding-agent-ready design artifact.

## Boundary

Use `recipe-agent-team-planning-grill` first when the goal, terminology, or acceptance criteria are still vague. Use this recipe only after the architecture question is concrete enough to evaluate candidates or write a selected technical design.

## Reference Loading

Read `references/language.md` when evaluating architecture candidates, writing technical designs, or discussing modules, interfaces, implementations, depth, seams, adapters, leverage, or locality.

Use `references/architecture-interview-router.md` to choose interview files. Load only two to four unique interview files selected by the router. Do not load every interview by default.

## Artifact Contract

Write design outputs under the run artifact root when a `RUN_ID` exists:

```text
_workspace/{run_id}/architecture-candidates.md
_workspace/{run_id}/technical-design.md
_workspace/{run_id}/implementation-tasks.json
_workspace/{run_id}/acceptance.md
```

Use `architecture-candidates.md` when multiple structural options exist. Use `technical-design.md` for the selected design. Use `implementation-tasks.json` only when downstream agents or tooling need structured task contracts. Use `acceptance.md` as a separate file when any of these are true: there are four or more distinct acceptance criteria, criteria span multiple agents or tasks, or the orchestrator needs to verify completion independently of individual task evidence. Otherwise, keep acceptance criteria inside `technical-design.md`.

If no run exists, summarize the design in the final response and defer file creation until there is something durable to preserve.

## Workflow

1. Read the planning artifact first when one exists.
2. Inspect the relevant code, tests, docs, and prior architecture notes before proposing a design.
3. Identify current friction: scattered logic, duplicated rules, weak ownership, hard-to-test behavior, unclear interfaces, or unstable dependencies.
4. Use `references/architecture-interview-router.md` to choose two to four interview lenses when candidates have meaningful tradeoffs.
5. Write `architecture-candidates.md` with one to three candidates unless the user or planning artifact already selected a direction.
6. Pass the selection gate before writing `technical-design.md`.
7. For the selected design, define the new structure, interface contracts, data flow, migration sequence, tests, and rollout risks.
8. Convert the selected design into implementation task contracts for downstream coding agents.
9. Record decisions and unresolved questions in the design artifact.

## Agent-Team Checks

Before writing a design for an agent-team workflow, check:

- existing harness artifacts under `_workspace/{run_id}/`
- prior planning, terminology, acceptance, review, or compound-learning artifacts
- whether the design preserves the state/artifact boundary
- whether task metadata remains compact and machine-readable
- whether workers receive runtime IDs from the orchestrator instead of creating them

When current artifacts contradict the proposed direction, record the contradiction as a design risk or open question. If the contradiction changes the design materially, ask the user before finalizing.

## Phase 1: Candidate Discovery

Use candidate discovery when the best architecture direction is not already clear.

Each candidate should include:

- name
- current friction it addresses
- proposed module/interface shape
- affected files or modules
- expected benefits
- migration cost
- risks
- tests or verification it would enable
- recommendation status: recommended, viable, risky, or rejected

Do not let downstream coding agents treat candidates as implementation instructions. Candidates are options until one passes the selection gate.

Use this structure for `_workspace/{run_id}/architecture-candidates.md`:

```markdown
# Architecture Candidates

## Goal

## Current Friction

## Candidate A: Name

## Candidate B: Name

## Candidate C: Name

## Interview Findings

## Recommendation

## Selection Gate
```

## Candidate Evaluation

Evaluate each candidate with concise evidence:

- **Locality**: related behavior moves closer together or gains a clearer owner.
- **Interface leverage**: callers get a smaller, more stable contract.
- **Testability**: behavior can be verified without excessive setup or unrelated systems.
- **Migration cost**: change can be staged without hiding risky rewrites.
- **Runtime fit**: the design preserves agent-team boundaries between state, messages, artifacts, and worker execution.

Avoid recommending a candidate just because it looks cleaner. Tie every recommendation to concrete friction in the current codebase or workflow.

Use the language reference consistently: describe design units as modules, caller-facing contracts as interfaces, variation points as seams, and concrete seam fillers as adapters.

## Selection Gate

Write `technical-design.md` only after one condition is true:

- the user chose a candidate
- the planning artifact already selected a direction
- one candidate is clearly recommended and the tradeoff is low-risk enough to proceed without another question

If candidates have meaningful tradeoffs, ask the user to choose before writing the selected design.

Record the selection in the design artifact:

```markdown
## Selected Candidate

- Candidate: ...
- Reason: ...
- Rejected Alternatives: ...
- Selection Source: user | planning artifact | low-risk recommendation
```

## Phase 2: Selected Technical Design

The selected design turns one candidate into a concrete implementation contract for coding agents. It should not keep multiple alternatives alive except as explicitly rejected alternatives or risks.

## Technical Design Format

Use this structure for `_workspace/{run_id}/technical-design.md`:

```markdown
# Technical Design

## Goal

## Inputs

## Design Language

## Selected Candidate

## Interview Findings

## Current Structure

## Selected Design

## Affected Files And Modules

## Interfaces And Data Flow

## Migration Plan

## Implementation Tasks

## Acceptance Criteria

## Risks And Open Questions
```

Keep the design concrete enough that a coding agent can implement from it without rediscovering the architecture.

## Implementation Task Contract

`technical-design.md` always contains an `## Implementation Tasks` section with human-readable task descriptions. `implementation-tasks.json` is the machine-readable companion and is written only when downstream agents or tooling need structured task contracts. When both exist, `implementation-tasks.json` is the authoritative contract for task IDs, scope, dependencies, and verification; the prose in `technical-design.md` provides rationale and context.

Each downstream task should include:

- `id`
- `title`
- `owner_agent`
- `scope`
- `non_scope`
- `required_behavior`
- `verification`
- `artifact`
- `depends_on`
- `acceptance_ref`
- `risk_notes`

Use task metadata only for compact machine-readable fields such as `design_item_id`, `candidate_id`, `acceptance_ref`, `artifact_contract`, or dependency IDs. Put full design prose in artifacts.

## Implementation Tasks JSON Shape

Use this minimum JSON shape when writing `_workspace/{run_id}/implementation-tasks.json`:

```json
{
  "version": 1,
  "source_design": "_workspace/RUN_ID/technical-design.md",
  "acceptance_ref": "_workspace/RUN_ID/acceptance.md",
  "tasks": [
    {
      "id": "design-task-1",
      "title": "Short imperative title",
      "owner_agent": "coder",
      "scope": ["files or modules in scope"],
      "non_scope": ["explicit exclusions"],
      "depends_on": [],
      "required_behavior": ["observable behavior or artifact requirement"],
      "verification": ["test, command, review, or artifact check"],
      "artifact": "_workspace/RUN_ID/design-task-1.md",
      "risk_notes": ["risks the worker must not silently resolve"]
    }
  ]
}
```

Keep task IDs stable within the design artifact. Use concrete file or module paths in `scope` when known.

## Completion

The architecture design is complete when:

- the selected design is traceable to current code or docs
- affected modules and interfaces are explicit
- design language is used consistently enough for downstream agents to preserve the intended shape
- migration can be executed as ordered coding tasks
- acceptance criteria are testable
- risks and unresolved questions are visible
- durable outputs are written under `_workspace/{run_id}/` when a run exists; when no run exists, the final response contains the selected design, affected modules, migration plan, implementation tasks, and open questions in full

After completion, hand off to a separate coding workflow that consumes the design artifacts and task contracts.
