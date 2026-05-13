# Design System Reference

Pattern for designs whose `Output` describes a reusable token catalog. Produces a single `DESIGN.md` file plus optional supporting docs. Downstream Outputs (product UI, brand identity, icon system) cite this catalog instead of inlining values. Cite this file as a `Pattern Hint` in the brief when applicable.

Format follows the `DESIGN.md` open specification from `google-labs-code/design.md` (alpha). YAML frontmatter holds machine-readable tokens; markdown body holds human-readable rationale.

## When To Cite This Pattern

Cite `design-system.md` in a brief's `Pattern Hints` when at least one is true:

- the project lacks a shared token catalog and two or more downstream surfaces (UI, logo, icons) will consume the same values
- a brand refresh changes color/typography/spacing across multiple surfaces
- downstream agents need machine-readable tokens for Tailwind/W3C DTCG export or lint
- the design rationale ("why this palette") must outlive the current surface

Skip this pattern and cite a surface pattern (`ui.md`, `logo-branding.md`, `icon-illustration.md`) directly when only one screen, one logo, or one icon set is in scope and tokens will not be reused.

## Output Layout

```
_workspace/{run_id}/design/{slug}/   # e.g., {slug}=token-catalog or heritage-tokens
├── DESIGN.md           # tokens (YAML frontmatter) + 8-section rationale (markdown)
├── rationale.md        # optional: extended decision notes per token group
└── migration.md        # optional: mapping from prior values to new tokens
```

`DESIGN.md` is the only mandatory file. `rationale.md` and `migration.md` exist only when their content cannot fit inside the corresponding DESIGN.md section.

## Pipeline Order

1. **DESIGN.md frontmatter** — fill token groups in this order: `colors`, `typography`, `rounded`, `spacing`, `elevation`, `components`. Reference earlier tokens from later ones using `{group.name}` syntax.
2. **DESIGN.md body** — write the 8 canonical sections (Overview, Colors, Typography, Layout, Elevation & Depth, Shapes, Components, Do's and Don'ts). Each body section explains the reasoning for the matching frontmatter group.
3. **rationale.md** (optional) — only when a token group needs more than a few paragraphs of justification.
4. **migration.md** (optional) — only when the project already had ad-hoc values that must map to the new tokens. Include a `old_value → new_token` table.

Token version pin: top of `DESIGN.md` frontmatter must include `design_md_version: "alpha"` so downstream agents detect spec drift.

`DESIGN.md` must start with YAML frontmatter. Cite the brief through
`routed_by`, `patterns`, and `brief` fields in frontmatter instead of a
top-line HTML comment. `DESIGN.md` is produced only by the spec recipe with a
`RUN_ID`. The brief is routed by `recipe-agent-team-design-interview`; the
`DESIGN.md` artifact is routed by `recipe-agent-team-design-spec`.

## Sub-Artifact Templates

### DESIGN.md

````markdown
---
design_md_version: "alpha"
routed_by: "recipe-agent-team-design-spec"
patterns: ["design-system.md"]
brief: "_workspace/{run_id}/design/{slug}/design-brief.md"
name: "<system name, e.g., Heritage>"
colors:
  primary: "#1A1C1E"
  secondary: "#FF6B35"
  background: "#FFFFFF"
  surface: "#F5F7FA"
  text:
    primary: "#111827"
    muted: "#6B7280"
  feedback:
    danger: "#DC2626"
    success: "#16A34A"
typography:
  display:
    fontFamily: "Public Sans"
    fontSize: "32px"
    fontWeight: 700
    lineHeight: 1.2
  heading:
    fontFamily: "Public Sans"
    fontSize: "20px"
    fontWeight: 600
    lineHeight: 1.3
  body:
    fontFamily: "Public Sans"
    fontSize: "14px"
    fontWeight: 400
    lineHeight: 1.5
rounded:
  sm: "4px"
  md: "8px"
  lg: "16px"
  full: "9999px"
spacing:
  xs: 4
  sm: 8
  md: 16
  lg: 24
  xl: 32
elevation:
  low: "0 1px 2px rgba(0,0,0,0.06)"
  mid: "0 4px 12px rgba(0,0,0,0.08)"
  high: "0 12px 32px rgba(0,0,0,0.12)"
components:
  button:
    background: "{colors.primary}"
    foreground: "{colors.background}"
    paddingX: "{spacing.md}"
    paddingY: "{spacing.sm}"
    borderRadius: "{rounded.md}"
    typography: "{typography.body}"
  card:
    background: "{colors.surface}"
    padding: "{spacing.lg}"
    borderRadius: "{rounded.lg}"
    elevation: "{elevation.mid}"
---

# <System Name>

## Overview

Brand philosophy, style intent, who the system serves, where it is used. One paragraph.

## Colors

Palette description. Why `primary` is `#1A1C1E` (e.g., "deep ink suggests heritage and trust"). Pairing rules. Contrast notes (WCAG AA/AAA targets). Theme variants if any (light/dark/high-contrast).

## Typography

Font choice rationale. Pairing rules (display vs body). Type scale steps and the ratio. Web font loading strategy (woff2, fallback stack).

## Layout

Spacing base unit (4px or 8px). Scale progression rule (linear, geometric). Grid system if any. Breakpoint targets.

## Elevation & Depth

Layering model: how many z-levels, what each represents (low = inline, mid = overlay, high = modal). Shadow construction (offset/blur/color rationale).

## Shapes

Border-radius scale rationale. When to use `sm` vs `lg` vs `full`. Geometry intent (e.g., rounded for friendly, sharp for technical).

## Components

For each component in the frontmatter `components` group, explain:

- intended usage (where this component appears)
- token bindings (which tokens it composes)
- variants and states (default, hover, disabled, loading)
- accessibility requirements (focus ring, contrast, keyboard)

## Do's and Don'ts

- Do: cite tokens by name in downstream specs (`{colors.primary}`), never inline hex.
- Do: extend the system by adding new tokens at the same scale tier as existing ones.
- Don't: introduce one-off hex values inside component specs.
- Don't: skip the `design_md_version` pin — alpha spec may change.
- Don't: store rendered assets here; this file is text-only.
````

### rationale.md (optional)

```markdown
# Token Rationale

## Color decisions

- Why `#1A1C1E` over pure black: pure black causes harsh contrast on OLED; `#1A1C1E` softens without losing legibility.
- Rejected alternatives: `#000000` (too harsh), `#222222` (too warm).

## Typography decisions

- Public Sans over Inter: open-license, similar metrics, US-government provenance reinforces trust positioning.
```

### migration.md (optional)

```markdown
# Migration Map

| Old value | Location(s) | New token | Notes |
|---|---|---|---|
| `#1A1B1E` | header.css, modal.css | `{colors.primary}` | Old value was a typo; standardize on `#1A1C1E`. |
| `Inter` 14px regular | body text across all surfaces | `{typography.body}` | No metric change; font swap only. |
```

## Acceptance Criteria

A design-system artifact passes when:

1. `DESIGN.md` exists with `design_md_version` pinned and all eight body sections present (no empty section; record "deferred" with a reason if intentionally skipped).
2. Frontmatter includes `routed_by`, `patterns`, `brief`, at least `colors`, `typography`, `spacing`, and one `components` entry. `rounded` and `elevation` may be empty groups with a recorded reason.
3. Every value in `components` references existing tokens via `{group.name}` syntax — no inline hex/px inside `components`.
4. Token names are scale-tier consistent (e.g., spacing uses `xs/sm/md/lg/xl` or `1/2/3/4/6/8`, not a mix).
5. Each body section includes at least one rationale sentence with "because", "so that", or an explicit trade-off/avoidance; value-only restatement is rejected.
6. If `migration.md` exists, every `old_value` row maps to a token defined in `DESIGN.md` frontmatter.

## Do Not

- Do not embed pixel-perfect mockups or rendered assets — this pattern produces tokens and rationale, not visuals.
- Do not bypass the `{group.name}` reference syntax inside `components`.
- Do not invent component variants here without a downstream consumer surface; add them in the consuming surface's reference and link back.
- Do not duplicate a token under two names — pick one canonical name and alias if necessary.
- Do not promote the `design_md_version` past `alpha` without verifying the upstream spec.

## Tips

- Start with `colors` and `typography`; downstream patterns (`ui.md`, `logo-branding.md`, `icon-illustration.md`) depend on them most heavily.
- When a downstream Output needs a value the system does not yet have, add a token to `DESIGN.md` rather than inlining; rerun this Output to update.
- For multi-theme systems (light/dark), store variants as nested keys (`colors.background.light`, `colors.background.dark`) and document switching rules in the Colors body section.
- Reference `{group.name}` even inside `rationale.md` examples so search-and-replace stays safe.
- Run a manual lint pass: every `{group.name}` reference must resolve to a defined token.

## Hand Off

| Situation | Hand off to |
|---|---|
| Tokens ready, screens needed | Run `recipe-agent-team-design-interview` for a new Output describing the screen set (Pattern Hint `ui.md`), then `recipe-agent-team-design-spec`; cite `DESIGN.md` in Upstream Inputs |
| Tokens ready, logo needed | Run `recipe-agent-team-design-interview` for a brand-identity Output (Pattern Hint `logo-branding.md`), then `recipe-agent-team-design-spec`; cite `DESIGN.md` in Upstream Inputs |
| Tokens ready, icon system needed | Run `recipe-agent-team-design-interview` for an icon-system Output (Pattern Hint `icon-illustration.md`), then `recipe-agent-team-design-spec`; cite `DESIGN.md` in Upstream Inputs |
| Token consumer needs Tailwind/W3C DTCG export | Note in handoff that `DESIGN.md` is exportable via the `design.md` CLI's `export` command (alpha) |
| Downstream surface needs a value not in tokens | Pause that surface; rerun this Output to add the missing token |
