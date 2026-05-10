# Codex Skill Testing and Iteration Guide

Use this reference to validate generated Codex harnesses and improve them after feedback.

## 1. Test Model

Skill quality is validated through four layers:

| Layer | Purpose |
| --- | --- |
| static checks | verify files, metadata, schemas, placeholders, and references |
| trigger checks | verify the skill is selected for the right user phrasing |
| dry-run checks | verify orchestrator routing, data flow, and failure branches |
| execution checks | compare with-skill behavior against a baseline when feasible |

Documentation-only harness edits can stop at static + dry-run checks. Generated harnesses intended for repeated use should include trigger and execution checks.

## 2. Static Checks

- agent TOML files exist under `.codex/agents/`
- required TOML fields are present: `name`, `description`, `developer_instructions`, `model`, `sandbox_mode`, `model_reasoning_effort`
- no template placeholders remain
- model ids match `references/schemas/models.md`
- sandbox modes match responsibility
- ordinary workers do not spawn subagents
- skills have valid frontmatter with trigger-rich descriptions
- referenced files exist
- orchestrator has route, execution mode, data flow, retry, blocked, follow-up, and completion rules
- runtime instructions use `agent-team`
- `_workspace/` is reserved for artifacts and reports
- `AGENTS.md` pointer and change history match actual files

## 3. Test Prompt Design

Prompts should sound like real user requests. Avoid artificial one-word prompts.

Cover:

- one normal core case
- one edge case
- one follow-up or partial-rerun case
- one near-miss where another skill or direct answer should win

Good prompt traits:

- concrete file names or expected output
- natural wording, including casual phrasing
- enough ambiguity to test routing
- realistic constraints, such as "only update the API docs section"

Weak prompts:

- "test the skill"
- "do docs"
- obviously unrelated prompts used as near-misses

## 4. Trigger Validation

For each skill, write:

- 8-10 should-trigger prompts with varied phrasing
- 8-10 should-not-trigger near-miss prompts
- 2-3 follow-up prompts

Near-miss prompts should be close enough to expose boundary issues. For example, a docs-review skill should reject a request to implement a feature unless the user specifically asks for documentation review.

Check for:

- missing common phrases
- overly broad descriptions
- collisions with neighboring skills
- follow-up phrases missing from descriptions
- stale references to non-Codex tools

## 5. With-Skill vs Baseline Evaluation

When feasible, compare:

- **with-skill**: run with the generated skill instructions and references
- **baseline**: same prompt without using the skill

Evaluate whether the skill improves:

- correctness
- completeness
- artifact structure
- evidence quality
- handling of edge cases
- follow-up behavior
- runtime-state discipline

Do not overfit fixes to one prompt. Generalize the improvement into the skill's rules, examples, output contracts, or references.

## 6. Assertion-Based Evaluation

For objective outputs, define assertions before running:

| Output Type | Example Assertions |
| --- | --- |
| files | expected files exist, no unexpected files, correct paths |
| docs | command examples present, links resolve, source references cited |
| code | tests pass, imports compile, public API unchanged where required |
| data | schema matches, required keys present, counts/ranges correct |
| state | done tasks have evidence + artifact, blocked tasks have reason |

For subjective outputs, use rubric criteria:

- audience fit
- domain specificity
- completeness
- clarity
- actionability
- absence of unsupported claims

## 7. Orchestrator Dry Run

For every route:

1. Identify selected pattern and execution mode: direct, delegated, or hybrid.
2. Name active specialists.
3. Name expected runtime tasks for orchestrated state-backed execution, including how the orchestrator resolves internal runtime context or creates generated-ID run/task records without asking the user for raw IDs.
4. Name expected artifacts.
5. Name required evidence.
6. Name the next consumer of each artifact.
7. Name retry and blocked behavior.
8. Confirm no phase depends on missing data.

## 8. Runtime Scenarios

Normal:

- resolve internal context or create generated-ID execution
- create active tasks when durable state is active
- workers produce evidence and artifacts
- orchestrator verifies and advances
- final report cites artifacts

Failure:

- worker or phase fails
- orchestrator retries with changed scope or prompt
- retry budget is exhausted
- task blocks with concrete reason
- workflow does not advance silently

Handoff:

- worker identifies ownership boundary
- message coordination includes reason and artifact context
- target accepts or rejects
- orchestrator records next action

Follow-up:

- previous artifact exists
- user requests partial rerun or refinement
- orchestrator preserves unaffected work
- changed scope is documented

## 9. Iteration Loop

1. Run static and dry-run checks.
2. Run trigger checks for ambiguous skills.
3. Run with-skill vs baseline when the skill is high-value or risky.
4. Convert failures into general rules or examples.
5. Re-run the smallest test set that covers the change.
6. Record material harness changes in `AGENTS.md`.
