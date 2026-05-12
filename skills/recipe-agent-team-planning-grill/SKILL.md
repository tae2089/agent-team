---
name: recipe-agent-team-planning-grill
description: "Recipe: Stress-test a plan against repo code, docs, domain terms, and agent-team runtime constraints before execution. Use for 'grill the plan', 'challenge approach', 'refine planning', 'sharpen terminology', 'create acceptance criteria', or turning fuzzy ideas into durable tasks. Skip when plan is concrete with clear acceptance criteria and task contracts; go to architecture-design, coding, review, or compound-learning."
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

# Agent Team Planning Grill

Use this recipe before creating or executing a meaningful agent-team workflow when the plan is ambiguous, high-impact, domain-heavy, or likely to benefit from sharper terminology and acceptance criteria.

The goal is not to debate indefinitely. The goal is to convert vague intent into a concrete, code-aware plan with durable artifacts and task contracts.

## Boundary

For a lightweight planning conversation, you may use this recipe without creating runtime state. Create state only when the user wants durable tracking, execution, follow-up, or artifacts tied to a run.

## Reference Loading

Read `references/grill-patterns.md` whenever this recipe is used. If the reference cannot be loaded, still follow the core grill rules below.

Core grill rules:

- Check goal, scope, vocabulary, constraints, dependencies, acceptance, failure modes, and handoff.
- Ask at most one blocking question at a time.
- Block only when the answer changes scope, owner, task ordering, artifact contract, acceptance criteria, or safety boundary.
- Record non-blocking uncertainty as an open question or risk and proceed with the recommended assumption.
- Route terminology ambiguity to `recipe-agent-team-terminology-context` and structural design ambiguity to `recipe-agent-team-architecture-design`.
- Stop grilling when the next recipe has enough artifact context to proceed.

## Artifact Contract

Write planning outputs under the run artifact root when a `RUN_ID` exists:

```text
_workspace/{run_id}/planning-grill.md
_workspace/{run_id}/plan.json
_workspace/{run_id}/acceptance.md
```

Use `planning-grill.md` for the human-readable decision log and rationale. Use `plan.json` for machine-readable task contracts whenever downstream agents or tooling will consume the plan. Use `acceptance.md` when success criteria are substantial enough to deserve a separate checklist.

If no run exists, summarize decisions in the final response and defer file creation until there is something durable to preserve.

## Workflow

1. Identify the plan, scope, affected system, and intended outcome.
2. Inspect existing code and docs before asking questions that the repository can answer.
3. Build a short decision tree: core domain terms, hard constraints, dependencies, task boundaries, failure modes, acceptance criteria, and irreversible choices.
4. Keep at most one active user-facing question at a time. If the unresolved decision does not block task boundaries or acceptance criteria, record it as a risk or open question and continue.
5. When a term, constraint, or decision is resolved, record it immediately in the planning artifact if a run artifact root exists.
6. Route terminology ambiguity to `recipe-agent-team-terminology-context` and structural design ambiguity to `recipe-agent-team-architecture-design` when those questions outgrow the planning artifact.
7. Convert stable decisions into task contracts only after the relevant ambiguity is resolved.
8. Escalate to ADR-style documentation only for decisions that are costly to reverse, surprising without context, and based on a real tradeoff.

## Planning Checks

Before asking the user, search for:

- root or context docs: `CONTEXT.md`, `CONTEXT-MAP.md`, `AGENTS.md`, `GEMINI.md`, `docs/`, `adr/`, `architecture/`
- domain vocabulary in code: package names, type names, command names, config keys, API routes, schema names
- existing workflow contracts: skills, agent definitions, run/task/message conventions, acceptance criteria, test fixtures

When code or docs contradict the user's plan, surface the contradiction directly and ask which source should change.

Example (Probe Format):

```md
현재 이해: plan puts large planning content into task metadata
막힌 결정: task metadata size boundary — bloated metadata hurts downstream agents
추천 답안: keep task metadata compact; write plan body to _workspace/{run_id}/planning-grill.md and task contracts to _workspace/{run_id}/plan.json (if wrong: workers parse oversized metadata + state churn)
질문: 이 경계 유지할까?
```

## Question Discipline

Ask at most one blocking question at a time unless the user explicitly asks for a full checklist.

Deliver every blocking question as a single 4-line Probe Format block:

```md
현재 이해: <one-sentence summary of the plan and what is already decided>
막힌 결정: <unresolved decision + why it matters, one line>
추천 답안: <recommended answer (if wrong: <consequence>)>
질문: <one concrete question, no compound clauses>
```

Rules:

- One question per turn. Compound questions are forbidden.
- `추천 답안` is required and must inline the consequence of being wrong so the requester sees the cost without a follow-up.
- After each answer, restate `현재 이해` in a one-line acknowledgment before the next probe.
- Do not ask the user to answer what can be discovered from local files, command output, or existing artifacts.
- Block only when the answer would change scope, owner, task ordering, artifact contract, acceptance criteria, or safety boundary. Otherwise record the uncertainty under `Open Questions` or `Risks` and proceed with the recommended assumption.

## Planning Artifact Format

Use this structure for `_workspace/{run_id}/planning-grill.md`:

```markdown
# Planning Grill

## Goal

## Current Facts

## Resolved Terms

## Decisions

## Open Questions

## Task Contracts

## Acceptance Criteria

## Risks
```

Keep entries short and traceable. Prefer dated bullets when the session evolves over multiple turns.

## Task Contract Shape

`planning-grill.md` may summarize task contracts for human review, but `plan.json` is the source of truth for downstream agents when it exists.

If `planning-grill.md` and `plan.json` disagree, update `plan.json` or mark it stale before handoff. Downstream agents should read task contract fields from `plan.json` and use `planning-grill.md` only for rationale, decisions, and open questions.

Each task contract should include:

- `id`
- `title`
- `owner_agent`
- `scope`
- `non_scope`
- `depends_on`
- `required_output`
- `required_evidence`
- `artifact`
- `acceptance_ref`
- `preserve_decisions`
- `risk_notes`

Use task metadata only for compact machine-readable fields such as `plan_item_id`, `depends_on_decision`, `artifact_contract`, or `acceptance_ref`. Put full prose in artifacts.

## Plan JSON Shape

Use this minimum JSON shape when writing `_workspace/{run_id}/plan.json`:

```json
{
  "version": 1,
  "source": "_workspace/RUN_ID/planning-grill.md",
  "acceptance_ref": "_workspace/RUN_ID/acceptance.md",
  "terms_ref": "_workspace/RUN_ID/terminology.md",
  "tasks": [
    {
      "id": "plan-task-1",
      "title": "Short imperative title",
      "owner_agent": "architect",
      "scope": ["files, modules, or workflow areas in scope"],
      "non_scope": ["explicit exclusions"],
      "depends_on": [],
      "required_output": "Concrete output expected from this task",
      "required_evidence": ["what must be verified or shown"],
      "artifact": "_workspace/RUN_ID/plan-task-1.md",
      "acceptance_ref": "_workspace/RUN_ID/acceptance.md",
      "preserve_decisions": ["decision IDs or short decision summaries"],
      "risk_notes": ["risks the worker must not silently resolve"]
    }
  ]
}
```

Task `id` values must be stable within the plan so `depends_on` can reference them. Omit optional refs only when the corresponding artifact does not exist.

## Handoff Metadata

Planning grill handoff metadata should use explicit artifact references:

- terminology work: `terminology_ref` or `context_map_ref`
- architecture design: `planning_grill_ref`, `acceptance_ref`, and unresolved structural questions
- coding tasks: `plan_item_id`, `scope`, `non_scope`, `required_evidence`, `artifact`, and `acceptance_ref`
- review: acceptance criteria, risk list, and decisions the reviewer must preserve
- compound learning: decisions or surprises worth capturing after execution

## ADR Boundary

Suggest an ADR only when all are true:

- changing the decision later would be meaningfully expensive
- future maintainers would not understand the decision from code alone
- the decision chose between real alternatives

If any condition is missing, keep the decision in the planning artifact instead.

## Completion

The planning grill is complete when:

- domain terms used by the plan are precise
- major dependencies and task boundaries are explicit
- code and docs have been checked for contradictions
- acceptance criteria are testable
- unresolved questions are either answered, marked as risks, or turned into blocked tasks
- any durable outputs are written under `_workspace/{run_id}/`

Do not start execution just because the plan sounds plausible. Start execution only after the task contracts and acceptance criteria are concrete enough for a worker to complete with evidence.
