---
name: recipe-agent-team-design-interview
description: "Deep design interview that produces a design brief. Socratic 4-line probes pick a subdomain, surface tensions, assumptions, priorities, and failure modes. Output: design-brief.md. Skip for backend architecture, terminology cleanup, fuzzy planning, implementation, review, post-run learning, or rendered assets. Hand off to recipe-agent-team-design-spec for subdomain artifacts."
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

# Agent Team Design Interview Recipe

Produces a design brief through deep Socratic interview. Routes to one `subdomain` and locks decisions before any artifact production. Hands off to `recipe-agent-team-design-spec` for subdomain artifact generation.

## Use / Skip

Use for design discovery, decision locking, subdomain routing, or brief generation before subdomain spec work.

Skip backend architecture, terminology cleanup, fuzzy planning, implementation, review, post-run learning, and binary/rendered assets. Skip when a usable brief already exists; go directly to `recipe-agent-team-design-spec`.

## Subdomains

Routing target only. This recipe does not load any subdomain reference. Subdomain identifiers must match `recipe-agent-team-design-spec` exactly.

| Subdomain | Used for |
| --- | --- |
| `design-system` | Reusable token catalog (DESIGN.md). Pick first when two or more surfaces share tokens, when a machine-readable catalog is missing, or when downstream specs should cite tokens instead of inlining values. |
| `ui` | Screens, flows, wireframes, components, product interface. |
| `icon-illustration` | Icon set, illustration system, taxonomy, SVG manifest. |
| `character` | Playable/NPC character, roster, skills, animation, assets. |
| `environment` | Level, zone, biome, layout, set dressing, lighting. |
| `logo-branding` | Logo, wordmark, symbol, brand voice, usage rules. |

Table order is not priority.

## Artifact Target

With `RUN_ID`, write the brief to the chosen subdomain's output root as `design-brief.md`. Output roots are defined by `recipe-agent-team-design-spec`; this recipe only writes the brief file there. Without `RUN_ID`, return the same brief structure inline and state durable files were skipped.

## Pipeline

Uncertainty Scan → Routing Gate → Discovery Phases → Synthesis Gate → Brief.

## Discovery

Use the smallest depth that still produces an actionable brief.

- **Fast path:** all four exist upstream: user/story, subdomain, priority, and constraints. If only some exist, use hybrid probing for missing phases.
- **Full path:** request is ambiguous, multiple workers need a durable source of truth, or decisions have meaningful tradeoffs.

Every user-facing probe is a single 4-line block. Use these Korean labels exactly; write values in the requester's language.

```md
현재 이해: <one-sentence summary of what is decided so far>
막힌 결정: <single biggest uncertainty right now>
추천 답안: <concrete starter answer the requester can react to>
질문: <one concrete question>
```

First probe convention: when no decisions exist yet, write `현재 이해: 인터뷰 시작 (아직 결정 없음)`.

Rules: one question per turn; no compound questions; `추천 답안` is required; restate `현재 이해` after each answer; do not ask what local context already answers; log every executed probe in the brief. Read `references/probe-toolkit.md` for technique-to-phase mapping.

## Uncertainty Scan

Before probing, inspect task body, `planning-grill.md` or `plan.json`, `terminology.md`, prior design artifacts, and the opening request.

Score phases as `2` no capture, `1` partial signal, `0` fully covered upstream. Emit the scan as a Probe Format block with `추천 답안` set to the ordered phase list. Confirm before starting unless the requester explicitly asked for a quick standalone draft.

## Routing Gate

Run `Subdomain Routing` after the scan and before any subdomain-specific probe. It is a gate, not a discovery phase. Capture `route.subdomain`; refuse to advance without an explicit subdomain.

Routing signals:

- `design-system`: two or more token consumers, missing catalog, shared colors/type/spacing, or brand refresh across surfaces.
- `ui`: screens, flows, wireframes, components, or product interface.
- `icon-illustration`: icon set, illustration system, taxonomy, or SVG manifest.
- `character`: playable/NPC character, roster, skills, animation, or assets.
- `environment`: level, zone, biome, layout, set dressing, or lighting.
- `logo-branding`: logo, wordmark, symbol, brand voice, or usage rules.

Tiebreaker: when shared-token signals match alongside a surface signal, recommend `design-system` first and run the consumer subdomain as a separate later interview. When two surface signals match, ask the requester to pick one.

Refuse escape: when no signal matches after two routing probe attempts, hand off to `recipe-agent-team-planning-grill` rather than looping.

Run discovery phases by descending gap. `gap_size = 0` phases skip probing but must cite the upstream path in the brief.

## Discovery Phases

- **Context & Story** captures `core_story.{user,before,after,pivotal_moment}`. For `design-system`, adapt to primary token consumer, consuming surfaces, before/after token workflow, and pivotal handoff moment. If no concrete user or moment can be named, stop and run planning grill.
- **Specificity Drill** captures concrete signals, earn-it examples, fake-it examples, and rejected adjectives.
- **Trade-off Surface** captures tensions, chosen side, and accepted cost.
- **Assumption Probe** captures at least three assumptions with evidence, breakage impact, and risk flag.
- **Forced Choice & Priority** captures priority ranking and `first_5_seconds`.
- **Failure Rehearsal** captures three failure modes with root cause, earliest signal, and mitigation seed. For `design-system`, include token drift, export/lint failure, or downstream surface misuse.

## Synthesis Gate

After phases, deliver one final Probe Format block:

- `현재 이해`: 3-5 sentences covering story, subdomain, priority, tension, risk
- `막힌 결정`: "brief 잠그기 전 빠진 부분?"
- `추천 답안`: "없음 - 진행"
- `질문`: "이대로 brief 잠그고 design-spec으로 넘어갈까?"

Only proceed after explicit confirmation, unless the requester asked for a quick standalone draft. Record confirmation as either the quote or `quick-draft requested`.

## Brief Contract

Use `references/brief-template.md` for `design-brief.md` or the inline no-run equivalent. The template is authoritative for headings and tables. The brief records Subdomain, Core Story, Subject, Specificity, Tensions, Assumptions, Priority, First 5 Seconds, Failure Modes, Constraints, References, Success Criteria, Upstream Inputs, Interview Log, and Routed By.

Valid subdomains are `design-system`, `ui`, `icon-illustration`, `character`, `environment`, and `logo-branding`. Empty sections are forbidden; record a refusal, risk, or upstream citation instead.

## Acceptance

1. Brief exists at the subdomain output root, or the final response contains the same structure with no run.
2. Brief sets a valid subdomain, target output root, and links to the matching `recipe-agent-team-design-spec` reference filename.
3. Every applicable phase has a capture, refusal, upstream citation, or skip.
4. Interview Log records gap scores for applicable phases plus skip reasons.
5. Synthesis Gate confirmation is recorded.
6. In an agent-team task, completion evidence cites the brief path.

A brief with only refusals does not pass. Minimum useful capture is explicit subdomain, constraints, and success criteria, either captured or upstream-cited.

## Do Not

- Load any subdomain reference from `recipe-agent-team-design-spec`.
- Produce subdomain artifacts (this recipe ends at the brief).
- Skip discovery entirely.
- Batch multiple user questions per turn.
- Pick more than one subdomain in a single interview.
- Bypass the brief; the spec recipe consumes it.

## Hand Off

| Situation | Hand off |
| --- | --- |
| Brief locked, subdomain artifacts needed | `recipe-agent-team-design-spec` with brief path |
| Routing refuses after two probe attempts | `recipe-agent-team-planning-grill` |
| Term ambiguity blocks routing | `recipe-agent-team-terminology-context` |
| Backend gap surfaces | `recipe-agent-team-architecture-design` |
| Multi-surface token sharing detected | This recipe again with `design-system` first, then a consumer subdomain pass |
