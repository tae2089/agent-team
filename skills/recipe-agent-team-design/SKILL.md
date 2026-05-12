---
name: recipe-agent-team-design
description: "Discovery-driven design spec producer. Socratic interview picks one subdomain; produces design brief + subdomain artifacts. Text/markdown only. Use for UI screens, character/level design, logos, icon systems, brand identity. Skip for backend modules (architecture-design), code, post-run learning, or rendered images."
metadata:
  version: 1.2.0
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

# Agent Team Design Recipe

Discovery-driven design spec producer. Routes by `subdomain` (`ui`, `icon-illustration`, `character`, `environment`, `logo-branding`). Runs a Socratic interview, then dispatches to the matching reference. Output: text/markdown spec, not rendered assets.

## When To Use / Not To Use

Use when requester needs structured specs, multiple workers need a single source of truth, or subdomain still ambiguous.

Skip when: backend modules (`architecture-design`), terminology cleanup (`terminology-context`), fuzzy goal (`planning-grill` first), implementation code, or rendered binary assets.

## Supported Subdomains

| Subdomain           | Reference                         | Output root                                           |
| ------------------- | --------------------------------- | ------------------------------------------------------ |
| `ui`                | `references/ui.md`                | `_workspace/{run_id}/design/ui/`                       |
| `icon-illustration` | `references/icon-illustration.md` | `_workspace/{run_id}/design/icon-illustration/`        |
| `character`         | `references/character.md`         | `_workspace/{run_id}/design/character/{character_id}/` |
| `environment`       | `references/environment.md`       | `_workspace/{run_id}/design/environment/{level_id}/`   |
| `logo-branding`     | `references/logo-branding.md`     | `_workspace/{run_id}/design/logo-branding/`            |

To add: create `references/<subdomain>.md` mirroring `references/ui.md` structure (Output Layout, Pipeline Order, Sub-Artifact Templates, Acceptance Criteria, Do Not, Tips, Hand Off). Append a row above.

## Pipeline

1. **Discovery interview** — Uncertainty Scan, then 7 phases in routing-aware gap order → `design-brief.md`.
2. **Subdomain dispatch** — load the matching reference, follow its pipeline.
3. **Acceptance** — both brief and subdomain artifacts must pass.

Phases may run out of numeric order. Coverage is fixed — every phase produces a capture, a recorded refusal, or a recorded upstream skip.

## Artifact Target

When `RUN_ID` exists, write durable outputs to the subdomain's `Output root` listed in the Supported Subdomains table above. `design-brief.md` always lives at `_workspace/{run_id}/design/design-brief.md`.

When no run exists, do not invent a run ID. Keep the same brief and subdomain structure in the final response instead of writing files, and state that durable files were skipped because no run was active.

## Step 1: Discovery Interview

Layered Socratic conversation. Never accept surface answers — probe until concrete.

### Probe Format

Every probe = a single 4-line block. One question per turn. Compound questions forbidden.

```md
현재 이해: <one-sentence summary of what is decided so far>
막힌 결정: <single biggest uncertainty right now>
추천 답안: <concrete starter answer the requester can react to>
질문: <one concrete question>
```

`추천 답안` is required. After each answer, restate `현재 이해` before the next probe. If an answer resolves a later phase, jump forward.

### Discovery Depth

Use the smallest depth that still produces an actionable spec:

- **Fast path**: if upstream inputs already identify a specific user/story, subdomain, priority, and constraints, ask only the highest-gap probes needed to fill missing captures. Cite upstream material for skipped phases.
- **Full path**: use when the request is ambiguous, multiple workers need a durable source of truth, or design decisions have meaningful tradeoffs.

Executed phases use at least the probes listed in the phase table. Prefer two techniques when the phase has real ambiguity; one technique is acceptable for routing or priority phases whose capture is already clear.

### Probe Toolkit

| Technique            | When to apply                     | Pattern                                                |
| -------------------- | --------------------------------- | ------------------------------------------------------ |
| 5 Whys cascade       | Generic surface answer            | "Why does that matter?" × up to 5                      |
| Specificity demand   | Abstract adjectives               | "Describe one user doing one task. No adjectives."     |
| Anti-reference       | Vague style direction             | "Show me a similar problem solved badly. What fails?"  |
| Story extraction     | Generic audience                  | "Walk through one user's day before and after."        |
| Forced choice        | Multiple priorities claimed       | "Keep only one. Which?"                                |
| Assumption probe     | "obviously" / "of course" appears | "List 3 assumptions. Which would break if false?"      |
| Trade-off surface    | Goals contradict                  | "To win X, what is acceptable to lose? Name the cost." |
| Failure rehearsal    | Confidence too high               | "How fails one year from now? Three scenarios."        |
| Constraint inversion | Timid decisions                   | "If unconstrained, what changes? Why blocked now?"     |

Log each executed probe in the brief.

### Uncertainty Scan (before phases)

1. Inspect upstream sources (task body, `planning-grill.md` or `plan.json`, `terminology.md`, prior design artifacts, requester's opening message).
2. Score each phase: **2** = no capture available, **1** = partial signal, **0** = upstream covers every capture.
3. Emit a scan summary as a Probe Format block (`추천 답안` = ordered phase list). Confirm before starting.
4. Run Phase 2 before subdomain-specific probing, then run remaining phases in descending-gap order. Probes within a phase ordered by the same rule.
5. Re-scan once mid-interview if an answer changes later-phase heat.

Skip rule: `gap_size = 0` phases skip probing but must cite the upstream path in the brief's matching section. Empty sections forbidden.

Routing rule: Phase 2 (Subdomain Routing) is the only ordering override. It runs before any subdomain-specific probe in later phases, even when another phase has a higher gap score.

### Phases

Each phase: deliver probes using the Probe Format. Bail-out triggers terminate the recipe and route elsewhere.

| #   | Phase                    | Goal                                          | Probe sequence                                                                                    | Probes                                        | Capture                                                                        | Bail-out                                                      |
| --- | ------------------------ | --------------------------------------------- | ------------------------------------------------------------------------------------------------- | --------------------------------------------- | ------------------------------------------------------------------------------ | ------------------------------------------------------------- |
| 1   | Context & Story          | Anchor decisions in a concrete human scenario | (a) name one specific user (b) day **before** (c) day **after** (d) pivotal moment                | specificity demand on each, 5 Whys on pivotal | `core_story.{user, before, after, pivotal_moment}`                             | No specific user/moment → stop, run `planning-grill`          |
| 2   | Subdomain Routing        | Pick `subdomain` key                          | One block. `추천 답안` infers from pivotal moment; re-ask with table if rejected                  | forced choice; assumption probe if ambiguous  | `domain.subdomain`                                                             | Refuse to advance without explicit subdomain                  |
| 3   | Specificity Drill        | Convert each adjective into measurable facts  | Per adjective: (a) 2 earn-it examples (b) 2 fake-it examples (c) concrete element separating them | specificity demand, anti-reference            | `specificity[] = {adjective, concrete_signals[], examples[], anti_examples[]}` | Adjective resists after 3 probes → mark `rejected_adjective`  |
| 4   | Trade-off Surface        | Force explicit cost acknowledgment            | Per tension: `추천 답안` picks winner + accepted cost; `질문` confirms cost in user-impact terms  | trade-off surface, forced choice              | `tensions[] = {axis_a, axis_b, chosen, accepted_cost}`                         | Refuses any cost → re-ask with constraint inversion           |
| 5   | Assumption Probe         | Surface hidden invariants                     | ≥ 3 assumptions. Per assumption: (a) name it (b) evidence + breakage impact                       | assumption probe, 5 Whys                      | `assumptions[] = {statement, evidence, breakage_impact}`                       | No evidence + high impact → flag `risk`, recommend validation |
| 6   | Forced Choice & Priority | Collapse multi-priority into a ranking        | (a) "Keep only one. Rank the next four." (b) "One feeling within 5 seconds of first contact?"     | forced choice; trade-off surface if tied      | `priority_ranking[]`, `first_5_seconds`                                        | "All equal" → re-ask with budget metaphor                     |
| 7   | Failure Rehearsal        | Pre-mortem                                    | Per scenario (target 3): (a) failure description (b) earliest signal                              | failure rehearsal, 5 Whys                     | `failure_modes[] = {scenario, root_cause, earliest_signal, mitigation_seed}`   | No imagined failure → constraint inversion                    |

### Synthesis

After all phases, deliver a final Probe Format block:

- `현재 이해`: 3–5 sentence restatement (Core Story + Subdomain + top priority + biggest tension + biggest risk)
- `막힌 결정`: "brief 잠그기 전 빠진 부분?"
- `추천 답안`: "없음 — 진행"
- `질문`: "이대로 brief 잠그고 Step 2로 넘어갈까?"

Capture the confirmation quote in `Routed By → Requester confirmation`. Only proceed after explicit confirmation, unless the requester explicitly asked for a quick standalone draft instead of an interactive interview; in that case mark `Requester confirmation: quick-draft requested`.

## Step 2: Write `design-brief.md`

When `RUN_ID` exists, produce `_workspace/{run_id}/design/design-brief.md`. Without `RUN_ID`, include the same brief structure in the final response. Every section maps to a phase capture. Empty sections forbidden — record refusal or upstream cite instead.

```markdown
# Design Brief

## Domain

- Subdomain: <ui | icon-illustration | character | environment | logo-branding>
- Reference file: references/<subdomain>.md
- Output root: <path from Artifact Target table, or "inline final response" when no RUN_ID>

## Core Story

- Specific user: <one person>
- Before / After / Pivotal moment: <concrete>

## Subject

- One-line summary / Audience / Problem-feeling

## Specificity

### Adjective: <name>

- Concrete signals / Examples (earn) / Anti-examples (fake)

### Rejected adjectives

- <adjective> — reason

## Tensions

| Axis A | Axis B | Chosen | Accepted cost |

## Assumptions

| Statement | Evidence | Breakage impact | Risk? |

## Priority Ranking

1. <must-keep> ... 5. <fifth>

## First 5 Seconds

- <single sentence>

## Failure Modes

| Scenario | Root cause | Earliest signal | Mitigation seed |

## Constraints

- Platform / Style / Brand-palette / Budget-deadline

## References

- Emulate / Avoid (from Specificity)

## Success Criteria

- <derived from Priority Ranking + First 5 Seconds>

## Upstream Inputs

- Planning (`planning-grill.md`/`plan.json`) / Terminology / Prior design (paths or "n/a")

## Interview Log

### Uncertainty Scan

- Initial order: <e.g., "P2 → P1 → P3 → P6 → P4 → P5 (skipped, cites planning-grill.md) → P7">
- Gap scores: <P1=2, P2=2, P3=1, ...>
- Re-scan triggers: <none | description>

### Phase Captures

- Phase N: <probes used → insight ; or "skipped: cites <path>">

## Routed By

- recipe-agent-team-design
- Reason: <Core Story → Domain pick>
- Requester confirmation: <quoted "yes proceed">
```

## Step 3: Subdomain Dispatch

Load only the chosen subdomain reference (no other subdomains). Follow its Pipeline Order, Sub-Artifact Templates, Acceptance Criteria, Do Not, Tips, Hand Off.

Every sub-artifact starts with:

```markdown
<!-- routed by recipe-agent-team-design; subdomain: <subdomain>; brief: _workspace/{run_id}/design/design-brief.md -->
```

When no `RUN_ID` exists, keep the same top-line comment but use `brief: inline final response`.

## Acceptance Criteria

A design run passes when:

1. `design-brief.md` exists with all sections filled and `subdomain` key set, or the final response contains the same structure when no `RUN_ID` exists.
2. Every phase produced a capture, a recorded refusal, or an explicit upstream citation in the brief.
3. `Interview Log → Uncertainty Scan` records gap scores for all 7 phases, executed order, and any re-scan triggers.
4. `Interview Log → Phase Captures` cites every executed probe, any upstream citation, and any fast-path skip.
5. The chosen subdomain reference's Acceptance Criteria all pass.
6. Subdomain artifacts cite the brief in their top-line comment.
7. No other subdomain references were loaded into context.
8. If running inside an agent-team task, `agent-team task complete` body cites the brief path as primary artifact.

## Do Not

- Skip the discovery interview or batch multiple questions per turn.
- Load multiple subdomain references at once.
- Render binary assets (PNG, SVG output, Figma).
- Bypass `design-brief.md` or mix subdomains in one run.

## Tips

- Narrow to one `subdomain` per run; propose sequential runs if multiple needed.
- Reuse upstream `planning-grill.md` / `plan.json` / `terminology.md` / prior design — cite in Upstream Inputs.
- Stable IDs (screen, character, level, icon) belong in subdomain artifacts, not brief.
- Cross-subdomain sharing (e.g., logo tokens in product UI) → cite the upstream path in the consumer's Upstream Inputs.

## Hand Off

| Situation                        | Hand off to                                                                                          |
| -------------------------------- | ---------------------------------------------------------------------------------------------------- |
| Spec ready, production starts    | Downstream tasks cite brief + subdomain artifacts; worker uses `recipe-agent-team-worker-checkpoint` |
| Discovery reveals plan gaps      | `recipe-agent-team-planning-grill`                                                                   |
| Discovery reveals term ambiguity | `recipe-agent-team-terminology-context`                                                              |
| Discovery reveals backend gaps   | `recipe-agent-team-architecture-design`                                                              |
