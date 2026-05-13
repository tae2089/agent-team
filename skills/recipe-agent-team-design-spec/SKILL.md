---
name: recipe-agent-team-design-spec
description: "Subdomain spec producer for locked design briefs. Consumes design-brief.md and writes subdomain artifacts via exactly one reference (ui/icon-illustration/character/environment/logo-branding/design-system). Skip when no brief exists; run recipe-agent-team-design-interview first. Skip for backend architecture, code, review, or post-run learning."
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

# Agent Team Design Spec Recipe

Consumes a locked `design-brief.md` and produces subdomain artifacts. Loads exactly one reference per pass. Does not interview.

## Use / Skip

Use when a brief from `recipe-agent-team-design-interview` exists and subdomain artifacts must be generated.

Skip when no brief exists (run the interview recipe first), or when the work is backend architecture, code, review, post-run learning, or rendered binary assets.

## Prerequisites

- `design-brief.md` exists at a subdomain output root, or the equivalent brief is provided inline.
- Brief sets a valid `subdomain` from the Subdomains table.
- Brief includes Acceptance source data (constraints, success criteria, priorities, tensions, failure modes).

If a prerequisite fails, stop and hand off to `recipe-agent-team-design-interview`.

## Subdomains

| Subdomain | Reference | Output root |
| --- | --- | --- |
| `design-system` | `references/design-system.md` | `_workspace/{run_id}/design/design-system/` |
| `ui` | `references/ui.md` | `_workspace/{run_id}/design/ui/` |
| `icon-illustration` | `references/icon-illustration.md` | `_workspace/{run_id}/design/icon-illustration/` |
| `character` | `references/character.md` | `_workspace/{run_id}/design/character/{character_id}/` |
| `environment` | `references/environment.md` | `_workspace/{run_id}/design/environment/{level_id}/` |
| `logo-branding` | `references/logo-branding.md` | `_workspace/{run_id}/design/logo-branding/` |

Brief lives in the same output root as `design-brief.md`.

## Pipeline

Validate brief → load one reference → produce artifacts → acceptance.

## Validate Brief

Before loading any reference, confirm:

- `subdomain` field is one of the six valid values.
- Output root in the brief matches the Subdomains table for that subdomain.
- Synthesis Gate confirmation is recorded.

Mismatch or missing fields → stop and hand off to interview recipe with the specific gap.

## Dispatch

Load only the reference for the brief's `subdomain`. Follow its Pipeline Order, Sub-Artifact Templates, Acceptance Criteria, Do Not, Tips, and Hand Off.

Every sub-artifact cites the brief according to its reference. Default citation:

```markdown
<!-- routed by recipe-agent-team-design-spec; subdomain: <subdomain>; brief: <brief path> -->
```

Exception: `design-system` `DESIGN.md` cites the brief through YAML frontmatter (`routed_by`, `subdomain`, `brief` fields) instead of an HTML comment — see `references/design-system.md`.

## Acceptance

1. Subdomain artifacts exist under the brief's output root.
2. The chosen reference's acceptance criteria all pass.
3. Artifacts cite the brief using the citation rule in the chosen reference.
4. No reference other than the brief's `subdomain` reference was loaded.
5. In an agent-team task, completion evidence cites the brief path plus the produced artifact paths.

## Multi-Subdomain Runs

A single `RUN_ID` may host multiple subdomain passes (e.g., `design-system` then `ui`). Each pass:

- has its own brief and its own output root from the Subdomains table
- loads only its subdomain reference
- may cite earlier pass artifacts via Upstream Inputs (e.g., a `ui` pass cites the prior `design-system` `DESIGN.md`)

Passes are sequential. Do not parallelize references in the same context window.

## Do Not

- Run without a brief.
- Load more than one subdomain reference per pass.
- Re-interview the requester; that responsibility belongs to the interview recipe.
- Render PNG, SVG output, Figma, or binary assets.
- Promote `design_md_version` past `alpha` without verifying the upstream `google-labs-code/design.md` spec.

## Hand Off

| Situation | Hand off |
| --- | --- |
| Artifacts ready, implementation needed | Downstream tasks cite brief + artifacts; workers use `recipe-agent-team-worker-checkpoint` |
| Brief is incomplete or ambiguous | `recipe-agent-team-design-interview` |
| Brief reveals plan gaps mid-spec | `recipe-agent-team-planning-grill` |
| Brief reveals term ambiguity | `recipe-agent-team-terminology-context` |
| Brief reveals backend gaps | `recipe-agent-team-architecture-design` |
| Missing shared token catalog | Hand off to interview with `design-system` recommended; resume here for the consumer subdomain after `DESIGN.md` exists |
