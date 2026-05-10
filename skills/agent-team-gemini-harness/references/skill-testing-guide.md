# Gemini Skill Testing and Iteration Guide

Use this reference to validate generated Gemini harnesses and improve them after feedback.

## 1. Test Model

Skill quality is validated through four layers:

| Layer | Purpose |
| --- | --- |
| static checks | verify files, metadata, schemas, placeholders, tools, and references |
| trigger checks | verify the skill is selected for the right user phrasing |
| dry-run checks | verify orchestrator routing, data flow, and failure branches |
| execution checks | compare with-skill behavior against a baseline when feasible |

Documentation-only harness edits can stop at static + dry-run checks. Generated harnesses intended for repeated use should include trigger and execution checks.

## 2. Static Checks

- agent Markdown files exist under `.gemini/agents/`
- required frontmatter fields are present: `name`, `description`, `kind: local`, `model`, `tools`
- no template placeholders remain
- model ids follow `references/schemas/models.md` or project Gemini configuration
- tool lists are explicit and do not use wildcards
- ordinary workers do not include `invoke_agent`
- runtime workers include `run_shell_command` when they call `agent-team`
- skills have valid frontmatter with trigger-rich descriptions
- worker-facing skill descriptions include an orchestrator routing boundary
- referenced files exist
- orchestrator has route, execution mode, data flow, retry, blocked, follow-up, reload reminder, and completion rules
- runtime instructions use `agent-team`
- `_workspace/` is reserved for artifacts and reports
- `GEMINI.md` pointer and change history match actual files

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

## 4. Trigger Validation

For each skill, write:

- 8-10 should-trigger prompts with varied phrasing
- 8-10 should-not-trigger near-miss prompts
- 2-3 follow-up prompts

Near-miss prompts should be close enough to expose boundary issues.

Check for:

- missing common phrases
- overly broad descriptions
- collisions with neighboring skills
- follow-up phrases missing from descriptions
- stale references to non-Gemini tools
- worker skills competing with the orchestrator

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
| tools | no wildcard tools, workers lack `invoke_agent`, runtime workers can run shell commands |

For subjective outputs, use rubric criteria:

- audience fit
- domain specificity
- completeness
- clarity
- actionability
- absence of unsupported claims

## 7. Orchestrator Dry Run

For every route:

1. Identify selected pattern and execution mode: direct, invoked, or hybrid.
2. Name active specialists.
3. Name expected runtime tasks for orchestrated state-backed execution, including how the orchestrator resolves internal runtime context or creates generated-ID run/task records without asking the user for raw IDs.
4. Name expected artifacts.
5. Name required evidence.
6. Name the next consumer of each artifact.
7. Name retry and blocked behavior.
8. Confirm no phase depends on missing data.
9. Confirm file edits require `/agents reload` and `/skills reload`.

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
6. Record material harness changes in `GEMINI.md`.
