# Planning Grill Patterns

Use these patterns when a plan needs more structure than one or two clarifying questions.

## Decision Dimensions

Check the plan across these dimensions:

- **Goal**: the user-visible outcome and why it matters.
- **Scope**: what is included, what is explicitly out of scope, and which files or systems are likely affected.
- **Vocabulary**: terms that must be canonical before agents start work.
- **Constraints**: compatibility, runtime state, artifact layout, permissions, safety, data, and release constraints.
- **Dependencies**: upstream decisions, task ordering, external services, and prerequisite artifacts.
- **Acceptance**: observable conditions that prove the work is complete.
- **Failure Modes**: likely ways the plan can fail, drift, or be misinterpreted.
- **Handoff**: which downstream recipe or worker consumes the plan next.

## Question Pattern

Ask one question at a time using the 4-line Probe Format from the parent SKILL.md:

```md
현재 이해: <one-sentence summary of the plan and what is already decided>
막힌 결정: <unresolved decision + why it matters, one line>
추천 답안: <recommended answer (if wrong: <consequence>)>
질문: <one concrete question, no compound clauses>
```

`추천 답안` is required and must inline the consequence of being wrong. Keep the question narrow enough that the user can answer with one decision. If the answer requires multiple decisions, split into separate probes — one per turn.

Block only when the answer changes scope, owner, task ordering, artifact contract, acceptance criteria, or safety boundary. Otherwise record the uncertainty as an open question or risk and continue with the recommended assumption.

## Handoff Routing Pattern

Route out of planning grill when a question becomes a different kind of work:

- Use `recipe-agent-team-terminology-context` when the blocker is naming, aliases, overloaded terms, or user/code vocabulary mismatch.
- Use `recipe-agent-team-architecture-design` when the blocker is backend module shape, interface placement, migration sequence, or implementation task structure.
- Use `recipe-agent-team-design` when the blocker is a UI/icon/character/environment/logo design brief instead of backend structure.
- Use coding workflow only after task contracts and acceptance criteria are concrete enough for evidence-backed completion.
- Use review workflow when the plan is stable and the concern is correctness, regression risk, or missing verification.
- Use `recipe-agent-team-compound-learning` after execution when the workflow produced reusable guidance.

## Acceptance Pattern

Good acceptance criteria are:

- testable by a worker or reviewer
- tied to behavior, artifact content, or command output
- scoped to one task or one workflow phase
- explicit about what does not count as complete

Avoid criteria that only say the result should be "clean", "robust", "done", or "better".

## Artifact Update Pattern

Record each resolved decision in the planning artifact before moving on. Prefer short bullets:

```markdown
- Decision: Keep planning prose in `_workspace/{run_id}/planning-grill.md`.
  Reason: Task metadata should stay compact and machine-readable.
  Impacts: Downstream tasks reference `planning_grill_ref` instead of copying the prose.
```

## Stop Conditions

Stop grilling when:

- the remaining unknowns are implementation details for a worker
- further questions would not change task boundaries or acceptance criteria
- unresolved risks are explicit and acceptable
- the next recipe has enough artifact context to proceed
