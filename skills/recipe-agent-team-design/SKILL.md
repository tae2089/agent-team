---
name: recipe-agent-team-design
description: "Recipe: Discovery-driven design recipe organized by vertical (product / game / marketing) and subdomain. Interviews the requester to clarify vertical and subdomain (e.g., product/ui, game/character, game/environment, marketing/logo-branding), then produces a design brief plus subdomain-specific specification artifacts. Framework-agnostic, text/markdown only. Use for 'design the UI', 'design a character', 'design a level', 'create a logo', 'spec an icon set', 'wireframe the screens', 'brand identity', 'visual design', 'environment design', 'plan the interface'. Do not use for backend module structure (use architecture-design), implementation code, code review, post-run learning, or rendered image production."
metadata:
  version: 1.1.0
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

Discovery-driven design spec producer for agent-team workflows. Organizes design domains under verticals (`product`, `game`, `marketing`). Runs a 2-step interview to identify vertical and subdomain, then dispatches to the matching reference to produce text/markdown specs. Output is structured spec, not rendered assets.

## When To Use

- Requester needs structured design specs before implementation or asset production.
- Multiple downstream workers need a single source of truth for the design.
- Vertical/subdomain not yet decided; need an interview to clarify scope first.

## When Not To Use

- Backend module/interface structure → use `recipe-agent-team-architecture-design`.
- Pure terminology cleanup → use `recipe-agent-team-terminology-context`.
- Plan is fuzzy at the goal level → use `recipe-agent-team-planning-grill` first.
- Direct coding/implementation → out of scope.
- Rendered images, Figma files, binary assets → out of scope (spec only).

## Supported Verticals and Subdomains

| Vertical | Subdomain | Domain key | Reference file | Output subdir |
|---|---|---|---|---|
| product | UI / screens | `product/ui` | `references/product/ui.md` | `_workspace/{run_id}/design/product/ui/` |
| product | Icon / illustration system | `product/icon-illustration` | `references/product/icon-illustration.md` | `_workspace/{run_id}/design/product/icon-illustration/` |
| game | Character | `game/character` | `references/game/character.md` | `_workspace/{run_id}/design/game/character/{character_id}/` |
| game | Environment / level | `game/environment` | `references/game/environment.md` | `_workspace/{run_id}/design/game/environment/{level_id}/` |
| marketing | Logo / branding | `marketing/logo-branding` | `references/marketing/logo-branding.md` | `_workspace/{run_id}/design/marketing/logo-branding/` |

To add a new subdomain:
- Existing vertical → create `references/<vertical>/<subdomain>.md` and append a row above.
- New vertical → create `references/<vertical>/` directory, add subdomain reference files, append rows.

Subdomain references must follow the standard structure (Output Layout, Pipeline Order, Sub-Artifact Templates, Acceptance Criteria, Do Not, Tips, Hand Off).

## Pipeline

1. **Discovery interview** — produce `design-brief.md` (always).
2. **Domain dispatch** — pick the matching reference file by `vertical/subdomain` key, load it on demand, follow its pipeline.
3. **Acceptance** — verify both the brief and the subdomain artifacts pass acceptance criteria.

Do not start subdomain-specific production before the brief is written.

## Step 1: Discovery Interview

Ask the requester (or read from task body) these questions. Skip any already answered in the task body or upstream artifacts.

| # | Question | Why |
|---|---|---|
| 1 | Which vertical? (`product` / `game` / `marketing`) | Top-level routing |
| 2 | Which subdomain within that vertical? | Pick reference file |
| 3 | Who is the audience / user / player? | Frames decisions |
| 4 | What problem does it solve, or what feeling should it create? | Drives concept |
| 5 | Hard constraints (platform, style, brand, budget, deadline)? | Bounds the spec |
| 6 | References to emulate? References to avoid? | Calibrates style |
| 7 | Success criteria for this design? | Defines acceptance |
| 8 | Existing artifacts to consume (plan, terminology, prior design)? | Pins inputs |

If question 1 or 2 is ambiguous, present the supported table and require the requester to pick a `vertical/subdomain` pair. Refuse to dispatch without both.

## Step 2: Write `design-brief.md`

Produce `_workspace/{run_id}/design/design-brief.md` with this structure:

```markdown
# Design Brief

## Domain
- Vertical: <product | game | marketing>
- Subdomain: <e.g., ui | character | environment | logo-branding>
- Domain key: <vertical/subdomain>
- Reference file: references/<vertical>/<subdomain>.md
- Output subdir: _workspace/{run_id}/design/<vertical>/<subdomain>/

## Subject
- One-line summary: <what is being designed>
- Audience: <segment>
- Problem / feeling: <core intent>

## Constraints
- Platform / medium: <where it appears>
- Style direction: <e.g., minimal, playful, ornate>
- Brand / palette anchors: <existing tokens or "fresh">
- Budget / deadline: <if relevant>

## References
- Emulate: <list>
- Avoid: <list>

## Success Criteria
- <bullet list — what makes this design 'done'>

## Upstream Inputs
- Plan: <path or "n/a">
- Terminology: <path or "n/a">
- Prior design: <path or "n/a">

## Routed By
- recipe-agent-team-design
- Reason: <one sentence from the interview>
```

## Step 3: Domain Dispatch

Load only the chosen subdomain reference file (do not load other verticals or subdomains). Follow its Pipeline Order, fill the Sub-Artifact Templates, satisfy Acceptance Criteria, observe Do Not, apply Tips, follow Hand Off.

Every sub-artifact must include a top-line comment citing the brief:

```markdown
<!-- routed by recipe-agent-team-design; key: <vertical>/<subdomain>; brief: _workspace/{run_id}/design/design-brief.md -->
```

## Acceptance Criteria (Recipe-Level)

A design run passes when:

1. `design-brief.md` exists with all sections filled and `vertical/subdomain` key set.
2. The chosen subdomain reference's Acceptance Criteria all pass.
3. Subdomain artifacts cite the brief in their top-line comment.
4. No subdomain references other than the chosen one were loaded into context.
5. Task body in `agent-team task complete` cites the brief path as the primary artifact.

## Do Not

- Do not skip the discovery interview; confirm `vertical/subdomain` explicitly.
- Do not load multiple subdomain references at once; one per run.
- Do not render binary assets (PNG, SVG output, Figma). Spec is text.
- Do not bypass `design-brief.md`.
- Do not mix subdomains in one run; create a new run for a second subdomain.

## Tips

- For ambiguous requests, the brief should narrow to one `vertical/subdomain`. If multiple are needed, propose sequential runs.
- Reuse upstream `plan.md`, `terminology.md`, or prior design artifacts when present — cite them in the brief's Upstream Inputs section.
- Stable IDs (screen IDs, character IDs, level IDs, icon names) belong in the relevant subdomain artifact, not the brief.
- When adding a new subdomain, mirror the structure of `references/product/ui.md` to keep the recipe consistent.
- Cross-vertical sharing (e.g., logo tokens used in product UI) → cite the upstream artifact path in the consumer's Upstream Inputs.

## Hand Off

| Situation | Hand off to |
|---|---|
| Spec ready, implementation/production begins | Create downstream tasks citing brief + subdomain artifacts; worker uses `recipe-agent-team-worker-checkpoint` |
| Discovery reveals plan-level gaps | Pause and run `recipe-agent-team-planning-grill` |
| Discovery reveals term ambiguity | Pause and run `recipe-agent-team-terminology-context` |
| Discovery reveals backend structure gaps | Pause and run `recipe-agent-team-architecture-design` |
