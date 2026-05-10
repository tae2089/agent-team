# Gemini Skill Writing Guide

Use this reference when creating or improving `.gemini/skills/{name}/SKILL.md`.

## 1. Purpose

A Gemini skill contains reusable process knowledge: procedures, checklists, output formats, domain rules, reference-loading instructions, and validation criteria. One-off task details belong in the orchestrator prompt or runtime task record, not permanently in the skill.

Skills should be Gemini-native:

- use YAML frontmatter with `name` and `description`
- keep the main body concise
- use progressive disclosure through `references/`
- use scripts for deterministic repeated work
- avoid unsupported runtime assumptions
- route durable coordination through Agent Team runtime skills backed by daemonless `agent-team` state

## 2. Frontmatter

```yaml
---
name: skill-name
description: "What the skill does, when to use it, follow-up triggers, and boundaries."
---
```

The description must include:

- concrete tasks handled
- likely user phrases
- follow-up phrases: update, refine, rerun, audit, review, partial rerun
- boundaries where another skill or direct answer is better
- for worker-facing skills, an orchestrator routing boundary

Strong:

```yaml
description: "Audits API documentation for stale commands, missing setup steps, broken links, and mismatched examples. Use for docs audit, docs update, command verification, partial rerun of a docs review, or improving a previous documentation artifact. Do not use for implementing the API unless documentation changes are requested."
```

Worker-facing strong:

```yaml
description: "Worker procedure for reviewing API documentation when routed by the docs orchestrator. Use for docs QA, command verification, partial review rerun, and review-driven fixes. For full docs workflows, route through the orchestrator skill."
```

## 3. Body Structure

Recommended structure:

```markdown
# Skill Name

## Purpose
## Inputs
## Workflow
## Output Format
## Validation
## References
```

Use short sections with concrete actions. Do not repeat project-level rules that belong in `GEMINI.md`.

## 4. Writing Principles

| Principle | Rule |
| --- | --- |
| Why-first | Explain the reason behind important rules so the model can generalize. |
| Lean body | Keep `SKILL.md` focused; move long examples/tables to `references/`. |
| Generalize | Turn feedback into reusable criteria, not one-off prompt patches. |
| Testable outputs | Define files, sections, evidence, and verdicts clearly. |
| Gemini-native | Use Gemini files, skills, agents, explicit tools, and available invocation rules; do not reference unsupported tools. |
| Route boundaries | Worker-facing skills should identify their orchestrator and avoid competing with the orchestrator for full workflow requests. |

## 5. Progressive Disclosure

| Level | Load Time | Content |
| --- | --- | --- |
| metadata | always visible | name + trigger-rich description |
| `SKILL.md` | when triggered | main workflow and output contract |
| `references/` | loaded only when needed | detailed examples, schemas, domain variants |
| `scripts/` | executed when useful | deterministic helpers that do not need to be pasted into context |

Split references by domain/provider when only one variant is needed.

If a reference exceeds roughly 300 lines, add a short table of contents.

## 6. Output Contracts

For worker-facing skills, define:

- artifact path convention
- report sections
- evidence required
- commands or files inspected
- pass/fix/blocked vocabulary when relevant
- runtime completion behavior when assigned a Agent Team task

Example:

```markdown
## Output Format

Write `_workspace/{plan}/{task_id}_review.md`:

- Verdict: PASS | FIX | BLOCKED
- Evidence: commands run, files inspected, source links
- Findings: severity, location, issue, recommendation
- Follow-up: retry instructions or block reason
```

## 7. Bundled Scripts

Create `scripts/` only for repeated deterministic work:

- schema validation
- link checking
- fixture generation
- report formatting
- repository-specific scans

Do not bundle scripts for vague reasoning work. The model should reason in the skill; scripts should automate repeatable checks.

## 8. Do Not Include

- one-off user task details
- secrets, credentials, or private tokens
- obsolete non-Gemini tool calls
- `_workspace` as a task board or sync check store
- unsupported assumptions about automatic agent communication
- long copied documentation that should be linked or summarized

## 9. Validation

Before finishing:

- frontmatter is valid
- trigger wording includes likely follow-ups and boundaries
- output contract is testable
- referenced files exist
- no unsupported runtime model is described
- runtime state is only mentioned for assigned durable workflows
