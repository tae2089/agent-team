---
name: recipe-agent-team-design-spec
description: "Generic spec producer for locked design briefs. Consumes design-brief.md, reads the Output description, consults pattern hints from references/, and produces artifacts under the brief's Output Path. Skip when no brief exists; run recipe-agent-team-design-interview first. Skip for backend architecture, code, review, or post-run learning."
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
---

# Agent Team Design Spec Recipe

> **PREREQUISITE:** Load `agent-team-shared` before any CLI command. `agent-team-run` and `agent-team-task` are required inside an agent-team run.

> **Worker integration:** Inside an agent-team task, the worker must follow `recipe-agent-team-worker-checkpoint` for sync check, inbox handling, task start, and completion. This recipe defines artifact generation only.

## Prerequisites

- `RUN_ID` exists. No inline fallback. If `RUN_ID` is missing, stop and ask the orchestrator to create one via `recipe-agent-team-run-lifecycle`; do not hand off to interview for this case.
- If only an inline brief exists after a run is created, materialize it unchanged to its Output Path as `design-brief.md`, updating only concrete path fields, then resume.
- `design-brief.md` exists at `_workspace/{run_id}/design/{slug}/design-brief.md`.
- Brief's `Output` section sets a one-paragraph description, a `Slug`, and an `Output Path`.
- Brief includes Acceptance source data (constraints, success criteria, priorities, tensions, failure modes).

If `RUN_ID` is missing, hand off to `recipe-agent-team-run-lifecycle` or the orchestrator. If brief content, Output, Slug, Output Path, source data, or confirmation is missing, hand off to `recipe-agent-team-design-interview` with the specific gap.

## Pattern Library

`references/` holds non-binding artifact patterns. `Pattern Hints` use basenames only (for example, `ui.md`) and resolve relative to `references/`. The brief's `Output` description is authoritative; load patterns only when `Pattern Hints` cites them, and only as guidance for artifact shape, file naming, and acceptance heuristics.

| Pattern Hint | Typical use |
| --- | --- |
| `ui.md` | Screens, flows, wireframes, components, surface tokens |
| `logo-branding.md` | Brand identity, wordmark, symbol, usage rules, palette |
| `design-system.md` | Reusable token catalog (DESIGN.md alpha spec) shared across surfaces |
| `character.md` | Playable/NPC character roster, skills, animation, assets |
| `environment.md` | Level/zone/biome layout, set dressing, lighting |
| `icon-illustration.md` | Icon set, illustration system, taxonomy, SVG manifest |

When the Output spans multiple patterns (e.g., shared tokens + a UI screen), default to separate Outputs and sequential passes. Compose patterns inside one Output Path only when the requester or task explicitly asks for one combined Output; record pattern precedence in the artifact citation (e.g., `patterns: ui.md > design-system.md`) and resolve conflicts by Output description first, then first-listed pattern.

## Pipeline

Validate brief → consult pattern hints → produce artifacts → acceptance.

## Validate Brief

Before producing any artifact, confirm:

- Brief's `Output Path` equals `_workspace/{run_id}/design/{slug}/` (the brief's own directory).
- Brief's `Output` description is a non-empty paragraph (not a single label).
- Brief's `Slug` is kebab-case and matches the directory name.
- Synthesis Gate confirmation is recorded as a requester quote or `quick-draft requested`.

Any mismatch → stop and hand off to the interview recipe with the specific gap.

## Dispatch

1. Read brief's `Output` description and `Pattern Hints`.
2. Load only patterns cited in `Pattern Hints`, resolving basenames under `references/`. If hints are `none`, infer artifact shape from the Output description.
3. Produce artifacts under the brief's Output Path. File names and section structure come from cited patterns when applicable, or from the Output description when not.
4. Every artifact cites the brief at the top:

```markdown
<!-- routed by recipe-agent-team-design-spec; brief: <brief path>; patterns: <patterns used or "none"> -->
```

Exception: a `DESIGN.md` (token catalog) cites the brief through YAML frontmatter (`routed_by`, `brief`, `patterns`) instead — see `references/design-system.md`.

## Acceptance

1. Every artifact exists under the brief's Output Path.
2. Each artifact addresses at least one Priority Ranking entry and at least one Success Criterion; orphan artifacts are rejected.
3. Cited patterns are respected for file naming and section structure; deviations are documented in the artifact's citation.
4. Every Pattern Hint listed in the brief is a valid Pattern Library basename and is either used or explicitly rejected with a reason.
5. When running as an agent-team task, the worker calls `agent-team task complete` with `--artifact` set to the Output Path directory (e.g., `_workspace/{run_id}/design/dashboard/`), and `--evidence` listing every produced filename plus the brief path. A single-file Output (e.g., DESIGN.md) may point `--artifact` at the file directly.

## Multi-Output Runs

A single `RUN_ID` may host multiple design Outputs (e.g., a token catalog followed by a UI screen that cites it). Each pass:

- has its own brief and Output Path
- consults only the patterns its brief hints at
- may cite earlier pass artifacts via Upstream Inputs

Passes are sequential. Do not parallelize patterns in the same context window.

## Do Not

- Re-interview the requester; that responsibility belongs to the interview recipe.
- Render PNG, SVG, Figma, or binary assets.
- Promote `design_md_version` past `alpha` without verifying the upstream `google-labs-code/design.md` spec.
- Force the Output into a pattern shape the brief did not endorse.

## Hand Off

| Situation | Hand off |
| --- | --- |
| Artifacts ready, implementation needed | Downstream tasks cite brief + artifacts; workers use `recipe-agent-team-worker-checkpoint` |
| Brief incomplete or ambiguous | `recipe-agent-team-design-interview` |
| Plan gaps surface mid-spec | `recipe-agent-team-planning-grill` |
| Term ambiguity surfaces | `recipe-agent-team-terminology-context` |
| Backend gaps surface | `recipe-agent-team-architecture-design` |
| Output requires a shared token catalog the brief did not produce | Interview for an Output of `reusable token catalog`; resume here for the consumer Output after `DESIGN.md` exists |
