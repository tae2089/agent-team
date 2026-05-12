---
name: recipe-agent-team-terminology-context
description: "Recipe: Align user, codebase, docs, and agent vocabulary into durable terminology context artifacts. Use for shared terms, glossary creation, context mapping, vocabulary cleanup, canonical naming, or ambiguity resolution before planning/design/coding/review. Skip when terms are already aligned, naming is unambiguous, or terminology artifact already covers it."
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

# Agent Team Terminology Context

Use this recipe to align the user's language, repository vocabulary, documentation terms, and downstream agent task language.

The output is a durable terminology context that later planning, architecture design, coding, review, and compound workflows can reference. It is especially useful when the same idea has multiple names, one term has multiple meanings, or user-facing language differs from code-facing language.

## Boundary

Use this recipe before `recipe-agent-team-planning-grill` when terminology is already unclear. Use it during planning or design when ambiguity appears mid-workflow.

## Reference Loading

Read `references/term-patterns.md` whenever this recipe is used. If the reference cannot be loaded, still preserve user terms, code terms, canonical terms, aliases, decisions, and open questions separately.

## Artifact Contract

Write terminology outputs under the run artifact root when a `RUN_ID` exists:

```text
_workspace/{run_id}/terminology.md
_workspace/{run_id}/context-map.md
```

Use `terminology.md` for canonical terms, aliases, definitions, and open questions. Use `context-map.md` when the terminology must coordinate across named tasks, named agent roles, artifact paths, code names, docs, or review checks. For a small single-task clarification with no cross-artifact or cross-role mapping, write only `terminology.md`.

Create `context-map.md` when any of these are true:

- two or more named tasks, agent roles, or artifact paths consume the terms
- one task consumes terms from two or more sources, such as user language plus code names, docs, or prior artifacts
- a term must be preserved during implementation or checked during review
- user-facing and code-facing names intentionally differ
- overloaded terms span workflow phases, artifacts, or modules

If no run exists, summarize terminology decisions in the final response and defer file creation until there is something durable to preserve.

## Workflow

1. Collect user terms from the current request and recent workflow artifacts.
2. Inspect relevant code, tests, docs, schemas, CLI commands, skills, and previous artifacts for existing vocabulary.
3. Identify term collisions, aliases, vague nouns, overloaded terms, and terms that differ between user language and code language.
4. Propose canonical terms with short definitions and source evidence.
5. Resolve terms from local context when possible and record open questions only for boundaries that cannot be inferred.
6. Record accepted terms, rejected aliases, and unresolved ambiguities.
7. Update downstream task contracts or planning/design artifacts to reference the canonical terms when appropriate.

## Terminology Checks

Before finalizing terminology, search for:

- root or context docs: `CONTEXT.md`, `CONTEXT-MAP.md`, `AGENTS.md`, `GEMINI.md`, `docs/`, `adr/`, `architecture/`
- public names: commands, flags, API routes, package names, type names, database tables, config keys
- workflow language: skills, agent definitions, `_workspace/{run_id}/` artifacts, acceptance criteria, task bodies
- tests and fixtures that reveal behavioral meaning

Prefer terms already established in public interfaces unless they are misleading. If code-facing names and user-facing names differ, preserve both and map them explicitly.

## Terminology Format

Use this structure for `_workspace/{run_id}/terminology.md`:

```markdown
# Terminology

## Canonical Terms

| Term | Meaning | Use When | Avoid Confusing With | Source |
| ---- | ------- | -------- | -------------------- | ------ |

## User Terms

| User Term | Canonical Term | Notes |
| --------- | -------------- | ----- |

## Code Terms

| Code Term | Canonical Term | Location | Notes |
| --------- | -------------- | -------- | ----- |

## Decisions

## Open Questions
```

Keep definitions short enough that downstream agents will actually use them. Put detailed rationale in decisions only when it affects implementation or review.

## Context Map Format

Use this structure for `_workspace/{run_id}/context-map.md` when the workflow spans multiple artifacts or agents:

```markdown
# Context Map

## Workflow Language

| Workflow Term | Meaning | Used By | Notes |
| ------------- | ------- | ------- | ----- |

## Artifact References

| Artifact | Terms It Defines | Terms It Consumes |
| -------- | ---------------- | ----------------- |

## Agent Task Vocabulary

| Agent Or Task | Required Terms | Forbidden Or Ambiguous Terms |
| ------------- | -------------- | ---------------------------- |

## Code And Docs Mapping

| Canonical Term | Code Name | Doc/User Name | Location |
| -------------- | --------- | ------------- | -------- |

## Terms To Preserve During Implementation

| Term | Must Preserve Because | Review Check |
| ---- | --------------------- | ------------ |
```

## Question Discipline

Ask the user only when the repository cannot answer the term boundary.

Good question shape:

```text
You are using "compound" to mean the integration phase after review. I found no existing repo term for that phase. Should "compound" be the canonical term, or should downstream tasks call it "integration"?
```

Avoid asking broad naming questions when a narrower term boundary is enough.

Block on a terminology question only when the answer would change task scope, module/interface naming, public API names, artifact names, or review criteria. Otherwise record the ambiguity as an open question and continue with the recommended canonical term.

## Next Recipe

After terminology context is complete, hand off explicitly:

- `recipe-agent-team-planning-grill`: when goals, scope, acceptance, or task contracts still need pressure testing.
- `recipe-agent-team-architecture-design`: when terminology is stable and the next decision is module/interface shape or implementation task structure.
- coding workflow: when terminology and task contracts are stable enough for implementation.
- review workflow: when terminology should be checked against an artifact, code change, or design.
- `recipe-agent-team-compound-learning`: when terminology decisions produced reusable guidance after a completed workflow.

Pass `terminology_ref`, optional `context_map_ref`, and any critical `canonical_terms` in downstream metadata.

## Downstream Contract

When later tasks depend on terminology, include compact metadata such as:

- `terminology_ref`
- `canonical_terms`
- `context_map_ref`
- `term_decision_id`

Use `terminology_ref` for the artifact path. Use `canonical_terms` only as an array of one to five term keys that the task must preserve, for example `["compound learning", "architecture candidate"]`. Use `context_map_ref` only when a context map exists. Use `term_decision_id` for one specific decision that affects the task.

Do not duplicate the full glossary inside task metadata. Put full prose in artifacts.

## Completion

The terminology context is complete when:

- key user terms map to canonical workflow or code terms
- overloaded terms are split or explicitly scoped
- aliases are documented
- unresolved terms are listed as open questions
- `context-map.md` is written when the context-map creation criteria are met, or intentionally skipped for a single-task clarification
- downstream artifacts know which terminology file to reference
- durable outputs are written under `_workspace/{run_id}/` when a run exists

After completion, continue with the selected next recipe using the canonical terms.
