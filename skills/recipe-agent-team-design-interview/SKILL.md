---
name: recipe-agent-team-design-interview
description: "Deep design interview that produces a design brief. Socratic 4-line probes lock the design's output description, tensions, assumptions, priorities, and failure modes. Output: design-brief.md. Skip for backend architecture, terminology cleanup, fuzzy planning, implementation, review, post-run learning, or rendered assets. Hand off to recipe-agent-team-design-spec for artifact production."
metadata:
  version: 2.0.0
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
      - recipe-agent-team-worker-checkpoint
---

# Agent Team Design Interview Recipe

> **PREREQUISITE:** Load `agent-team-shared` before any CLI command. `agent-team-run` and `agent-team-task` are required inside an agent-team run.

> **Worker integration:** Inside an agent-team task, the worker must follow `recipe-agent-team-worker-checkpoint` for sync check, inbox handling, task start, and completion. This recipe defines brief generation only.

## Output Description

This recipe does not enumerate design types. The brief's `Output` field is plain-text and authoritative. Examples:

- "Solo-fitness app main dashboard screen — onboarded-user view, mobile-first, light/dark themes."
- "Heritage-tone brand identity — logotype, symbol, palette, usage rules for digital + print."
- "Reusable token catalog (DESIGN.md alpha) shared across product UI and marketing surfaces."
- "Playable NPC character — combat role, skills, animation states, and visual style guide."
- "Level 1 zone layout — set dressing, lighting beats, traversal flow."

Long, concrete descriptions beat one-word categories.

`recipe-agent-team-design-spec/references/` ships non-binding patterns (`ui.md`, `logo-branding.md`, `design-system.md`, `character.md`, `environment.md`, `icon-illustration.md`). Cite relevant ones in the brief's `Pattern Hints` field; the spec recipe may compose, adapt, or ignore them.

## Artifact Target

With `RUN_ID`, write the brief to `_workspace/{run_id}/design/{slug}/design-brief.md`. `{slug}` is a short kebab-case identifier captured during Output Capture (e.g., `dashboard`, `heritage-brand`, `token-catalog`, `npc-roster-v1`, `zone-1`).

Without `RUN_ID`, return the same brief structure inline and state durable files were skipped.

## Pipeline

Uncertainty Scan → Discovery Phases → Synthesis Gate → Brief.

## Discovery

Use the smallest depth that still produces an actionable brief.

- **Fast path:** the brief's minimum viable fields exist upstream: Output, Core Story, Priority Ranking, Constraints, and Success Criteria. Score covered phases `0` and probe only remaining applicable gaps.
- **Full path:** request is ambiguous, multiple workers need a durable source of truth, or decisions have meaningful tradeoffs.

Every user-facing probe is a single 4-line block. Korean labels exact; values in the requester's language.

```md
현재 이해: <one-sentence summary of what is decided so far>
막힌 결정: <single biggest uncertainty right now>
추천 답안: <concrete starter answer the requester can react to>
질문: <one concrete question>
```

First probe convention: when no decisions exist yet, write `현재 이해: 인터뷰 시작 (아직 결정 없음)`.

Rules: one question per turn; no compound questions; `추천 답안` is required; restate `현재 이해` after each answer; do not ask what local context already answers; log every executed probe in the brief. See `references/probe-toolkit.md` for technique-to-phase mapping.

## Uncertainty Scan

Inspect task body, `planning-grill.md` or `plan.json`, `terminology.md`, prior design artifacts, and the opening request.

Score phases as `2` no capture, `1` partial signal, `0` fully covered upstream. Emit the scan as a Probe Format block with `추천 답안` set to the ordered phase list. Confirm before starting unless the requester explicitly asked for a quick standalone draft. Record any later phase reordering or reopened gap as an Interview Log re-scan trigger; otherwise record `none`.

The first phase is **Output Capture** unless the upstream task body already contains an unambiguous one-paragraph output description (score `0` and cite the source).

## Discovery Phases

- **Output Capture** — `Output` description, `{slug}`, `Pattern Hints`, brief `Subject` summary. Use literal `none` when no Pattern Hint applies. If output stays vague after two probes, hand off to `recipe-agent-team-planning-grill`.
- **Context & Story** — `Core Story`, audience, problem-feeling, emulate/avoid references. If no concrete user or moment can be named, run planning grill.
- **Specificity Drill** — concrete signals, earn-it examples, fake-it examples, rejected adjectives.
- **Trade-off Surface** — tensions, chosen side, accepted cost, hard constraints.
- **Assumption Probe** — at least three assumptions with evidence, breakage impact, risk flag, upstream input citations.
- **Forced Choice & Priority** — priority ranking, `First 5 Seconds`, success criteria.
- **Failure Rehearsal** — three failure modes with root cause, earliest signal, mitigation seed.

Probes adapt to the captured Output. Pattern Hints name likely artifact patterns only; do not load pattern files during interview, and let the Output description override hint choice.

## Synthesis Gate

After phases, deliver one final Probe Format block:

- `현재 이해`: 3-5 sentences covering output, story, priority, tension, risk
- `막힌 결정`: "brief 잠그기 전 빠진 부분?"
- `추천 답안`: "없음 - 진행"
- `질문`: "이대로 brief 잠그고 design-spec으로 넘어갈까?"

Only proceed after explicit confirmation, unless the requester asked for a quick standalone draft. Record confirmation as either the quote or `quick-draft requested`.

## Brief Contract

Use `references/brief-template.md`. Empty sections are forbidden; record a refusal, risk, or upstream citation instead.

## Acceptance

1. Brief exists at `_workspace/{run_id}/design/{slug}/design-brief.md`, or the final response contains the same structure with no run.
2. Brief's `Output` section contains a one-paragraph description plus `{slug}` and Output Path.
3. Every applicable phase has a capture, refusal, upstream citation, or skip.
4. Interview Log records gap scores for applicable phases plus skip reasons.
5. Synthesis Gate confirmation is recorded as a requester quote or `quick-draft requested`.
6. When running as an agent-team task, call `agent-team task complete` with `--artifact` set to the brief path and `--evidence` citing the brief path plus confirmation state.

A brief with only refusals does not pass. Minimum: explicit Output, constraints, success criteria — captured or upstream-cited.

## Do Not

- Load pattern reference files; cite their paths only.
- Produce artifacts (this recipe ends at the brief).
- Bypass Output Capture without a `0` score and upstream citation.
- Batch multiple user questions per turn.
- Force the Output into a single pattern label when the description spans multiple shapes; record all relevant hints instead.

## Hand Off

| Situation | Hand off |
| --- | --- |
| Brief locked, artifacts needed and `RUN_ID` exists | `recipe-agent-team-design-spec` with brief path |
| Brief locked, artifacts needed but no `RUN_ID` | Ask orchestrator to create a run and materialize the inline brief before spec |
| Output stays vague after two probes | `recipe-agent-team-planning-grill` |
| Term ambiguity blocks Output Capture | `recipe-agent-team-terminology-context` |
| Backend gap surfaces | `recipe-agent-team-architecture-design` |
| Multi-output effort (e.g., shared tokens + consuming surface) | Re-run this recipe per Output; cite earlier briefs as Upstream Inputs |
