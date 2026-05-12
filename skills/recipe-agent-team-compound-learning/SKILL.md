---
name: recipe-agent-team-compound-learning
description: "Recipe: Capture reusable learnings from completed agent-team runs, tasks, reviews, or fixes so future agents skip re-discovery. Use for 'compound learning', 'document what worked', 'capture a reusable pattern', 'summarize lessons', 'update solution docs', 'close the loop'. Do not use for in-progress runs, uncertain outcomes, trivial fixes, implementation, planning, architecture design, or code review."
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
      - agent-team-ops
---

# Agent Team Compound Learning

Use this recipe at the end of a non-trivial run, task, review, bug fix, architecture decision, or workflow iteration to capture reusable learning for future agents.

The goal is not to archive everything. The goal is to turn fresh context into a concise, discoverable learning that makes the next similar unit of work easier.

## Boundary

Use this recipe after coding, review, compound/integration, debugging, or a successful planning/design iteration that produced a reusable convention or decision.

## Artifact Contract

Write run-scoped capture under the run artifact root when a `RUN_ID` exists:

```text
_workspace/{run_id}/compound-learning.md
```

When the learning is broadly reusable beyond the current run, propose or write a durable solution doc:

```text
docs/solutions/{category}/{slug}.md
```

Use `_workspace/{run_id}/compound-learning.md` for session-local synthesis. Use `docs/solutions/` only for verified, reusable knowledge that future workflows should search before planning or coding.

If no run exists, summarize the learning in the final response and defer file creation until there is something durable to preserve.

Use these solution doc categories unless the repository already has a different taxonomy:

- `agent-workflow`: agent coordination, run/task/message/sync, harness behavior
- `architecture`: module shape, boundaries, design patterns, migration strategy
- `testing`: test strategy, fixtures, flaky tests, verification patterns
- `tooling`: CLI, build, install, release, scripts, developer tools
- `docs`: documentation structure, examples, stale references
- `data`: schemas, migrations, serialization, persistence
- `integration`: external services, APIs, compatibility, cross-module behavior
- `security`: credentials, permissions, sandboxing, access control

Prefer an existing category over inventing a near-duplicate. Add a new category only when none of these describes the learning.

## Workflow

1. Gather evidence from run summary, task records, event log, artifacts, reviews, and the current conversation.
2. Decide whether the learning is worth capturing. Skip trivial typos, obvious one-line fixes, purely mechanical formatting, or unverified solutions.
3. Classify the learning as bug track or knowledge track.
4. Search existing solution docs and relevant artifacts for overlap.
5. If overlap is high, update or recommend updating the existing doc instead of creating a duplicate.
6. Write the run-scoped compound learning artifact.
7. When appropriate, write or propose a `docs/solutions/{category}/{slug}.md` doc.
8. Surface any refresh or follow-up recommendation with a narrow scope.

## Evidence Sources

Prefer concrete evidence over memory:

- `agent-team run summary --run RUN_ID`
- `agent-team event log --run RUN_ID`
- task evidence and artifact paths
- review artifacts
- planning/design artifacts
- test output or verification notes
- relevant changed files

Do not present a solution as verified unless there is evidence that it worked.

## Capture Criteria

Capture when at least one is true:

- the issue took meaningful investigation
- the root cause was non-obvious
- a reusable workflow, convention, or architectural pattern emerged
- a previous assumption was corrected
- future agents are likely to make the same mistake
- review found a pattern broader than a single line
- the work changed how future planning, design, coding, or review should proceed

Skip when the work is too local to reuse or the outcome is still uncertain.

## Tracks

Use bug track for failures, regressions, broken tests, runtime issues, integration problems, data issues, security issues, performance issues, and debugging workflows.

Use knowledge track for architecture patterns, design decisions, conventions, workflow practices, tool behavior, testing strategy, and agent coordination lessons.

## Compound Learning Format

Use this structure for `_workspace/{run_id}/compound-learning.md`:

```markdown
# Compound Learning

## Summary

## Track

## Evidence

## What Happened

## What Worked

## What Did Not Work

## Reusable Guidance

## Future Search Terms

## Solution Doc Decision

## Follow-Ups
```

Keep the run-scoped artifact concise. Link to source artifacts instead of copying long excerpts.

`Future Search Terms` should contain terms a future agent is likely to search when facing the same issue. Include concrete command names, file paths, error codes, API names, module names, domain terms, and symptom phrases. Avoid generic terms such as "bug", "fix", "error", "issue", "broken", or "works" unless paired with a concrete identifier.

`Solution Doc Decision` must be one of:

- `promote`: create a new `docs/solutions/{category}/{slug}.md`
- `update_existing`: update or recommend updating an existing solution doc
- `skip`: keep the learning only in the run-scoped artifact

Include one sentence explaining the decision.

## Solution Doc Format

Use this structure for `docs/solutions/{category}/{slug}.md` when the learning should outlive the run:

```markdown
---
title: ""
category: ""
problem_type: ""
source_run: ""
tags: []
created: "YYYY-MM-DD"
---

# Title

## Context

## Symptoms Or Trigger

For bug-track docs, describe the observed symptom. For knowledge-track docs, describe the situation that makes the guidance relevant.

## Root Cause Or Principle

## Working Solution

## What Did Not Work

## When To Apply

## Prevention Or Review Checklist

## References
```

## Overlap Check

Before creating a new solution doc, search `docs/solutions/` if it exists.

Compare potential overlap by:

- problem statement or situation
- root cause or principle
- solution approach
- referenced files, modules, or workflow steps
- prevention or review guidance

If most dimensions overlap, update the existing doc or recommend a targeted refresh. If only some dimensions overlap, create the new doc and mention the related doc. If no solution store exists, create it only when the learning is clearly reusable.

## Discoverability

Knowledge compounds only if future agents can find it.

If `docs/solutions/` is created or materially used and root instructions such as `AGENTS.md` or `GEMINI.md` do not mention it, propose the smallest discoverability addition that matches the existing file style.

Apply the discoverability edit only when one of these is true:

- the user explicitly asked to update project instructions or make the learning discoverable
- the current task already includes editing `AGENTS.md`, `GEMINI.md`, README, or harness instructions
- the solution doc will be hard to find without a pointer and the repository already uses root instruction pointers for similar resources

## Downstream Contract

When a later task should consume the learning, include compact metadata such as:

- `compound_learning_ref`
- `solution_doc_ref`
- `learning_track`
- `refresh_scope`

Do not duplicate the full learning in task metadata. Put full prose in artifacts.

## Completion

Compound learning is complete when:

- the learning is classified as bug track or knowledge track
- claims are backed by run, task, review, artifact, or test evidence
- reusable guidance is explicit
- overlap with existing docs has been checked when relevant
- durable outputs are written under `_workspace/{run_id}/` when a run exists
- any `docs/solutions/` promotion or refresh recommendation is clear

After completion, feed the learning back into terminology context, planning grill, architecture design, coding, or review in the next workflow iteration.
