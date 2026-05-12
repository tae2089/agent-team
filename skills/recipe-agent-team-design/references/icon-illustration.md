# Product / Icon-Illustration Reference

Use when discovery selects subdomain `icon-illustration`. Produces icon set or illustration system specifications: grid, style rules, taxonomy, naming, accessibility, and asset manifest. Text-based — no rendered icons.

## Output Layout

```
_workspace/{run_id}/design/icon-illustration/
├── style.md         # visual rules, geometry, stroke, corner radius
├── taxonomy.md      # category tree + canonical names + synonyms
├── usage.md         # sizing, color application, pairing with text, a11y
└── manifest.md      # icon/illustration list with status, references, deliverables
```

## Pipeline Order

1. **style.md** — visual rules: grid, stroke weight, corner radius, terminals, fill rules, perspective.
2. **taxonomy.md** — category structure, canonical names, synonyms, deprecation list.
3. **usage.md** — sizing scale, color tokens, pairing with text, accessibility (focus, contrast, alt text), state variants (filled/outlined/duotone).
4. **manifest.md** — concrete list of icons or illustrations to produce, status, references, deliverable formats.

## Sub-Artifact Templates

### style.md

```markdown
# Style Rules

## Grid
- Base grid: 24×24
- Keyline shapes: square 18×18, circle 18ø, vertical 16×20, horizontal 20×16
- Padding from edge: 2px

## Stroke
- Default weight: 2px
- Endcap: round
- Joins: round
- Open shapes allowed: yes (for outlined set)

## Corner Radius
- Default: 2px
- Allowed: 0, 1, 2, 4 (no free values)

## Terminals & Joins
- Maintain optical balance over geometric precision.
- All terminals snap to half-pixel grid.

## Fill Rules
- Outlined: 2px stroke, no fill.
- Filled: solid fill only, no inner stroke.
- Duotone (optional): primary 100%, secondary 40% tint.

## Perspective (for illustrations)
- Default: flat / orthographic
- If isometric: 30° axes
- No 3-point perspective
```

### taxonomy.md

```markdown
# Taxonomy

## Categories
- action/ — verbs the user performs (e.g., add, edit, share)
- object/ — nouns the system represents (e.g., file, user, device)
- status/ — system states (e.g., success, warning, loading)
- navigation/ — wayfinding (e.g., arrow, chevron, menu)
- communication/ — messaging (e.g., chat, mail, notification)

## Canonical Names
- Singular, lowercase, hyphen-separated: `arrow-right`, `user-circle`, `file-text`.
- No plurals: `file` not `files`.
- No platform prefixes: `share` not `ios-share`.

## Synonyms Map
| Canonical | Synonyms (map to canonical) |
|---|---|
| trash | bin, delete, remove |
| user | profile, account, person |
| settings | gear, cog, preferences |

## Deprecation
| Old name | Replacement | Removal date |
|---|---|---|
| user-old | user-circle | <yyyy-mm-dd> |
```

### usage.md

```markdown
# Usage

## Sizing Scale

| Size | Pixel | Use |
|---|---|---|
| xs | 12 | Inline with body text (rare) |
| sm | 16 | Inline, dense UI |
| md | 20 | Default UI buttons |
| lg | 24 | Standalone affordances |
| xl | 32 | Empty states, hero |

## Color Application
- Inherit current text color by default (`currentColor`).
- Status icons use semantic tokens (`color.fg.success`, `color.fg.danger`, etc.).
- Duotone secondary: 40% tint of primary.

## Pairing with Text
- Icon size = body line-height × 1.0.
- Gap between icon and text: spacing token `space.2` (8px).
- Alignment: icon optical center matches text x-height.

## Accessibility
- Decorative icon: `aria-hidden="true"`.
- Meaningful icon: provide `aria-label` or visible text alternative.
- Minimum contrast: WCAG AA against background.
- Focus state for interactive icons: 2px outline using `color.accent`.

## State Variants
- default
- hover
- active
- disabled (50% opacity)
- selected (filled variant)
```

### manifest.md

```markdown
# Icon / Illustration Manifest

## Status Legend
- planned, in-progress, in-review, ready, deprecated

## Icons

| Name | Category | Status | Notes |
|---|---|---|---|
| arrow-right | navigation | ready | — |
| arrow-left | navigation | ready | mirror of arrow-right |
| user-circle | object | ready | replaces deprecated user-old |
| settings | action | in-review | gear motif, match radius rules |
| trash | action | planned | reference: Material trash-can |

## Illustrations

| Name | Status | Use | Notes |
|---|---|---|---|
| empty-inbox | ready | Inbox empty state | flat, palette: brand neutral |
| onboarding-welcome | in-progress | Onboarding hero | 16:9, illustration grid |
| error-404 | planned | 404 page | playful, contained palette |

## Deliverables
- Format: SVG (single-file, optimized).
- Naming: `{category}_{name}.svg` (e.g., `navigation_arrow-right.svg`).
- Source files (if any) under `_workspace/{run_id}/design/icon-illustration/source/`.
- Optimization: stroke not expanded, viewBox preserved, no inline styles.
```

## Acceptance Criteria

An icon/illustration system passes when:

1. `style.md` defines grid, stroke, corner radius, terminals, fills, perspective.
2. `taxonomy.md` provides categories, canonical naming rules, synonym map, deprecation table.
3. `usage.md` defines sizing scale, color application, text pairing, accessibility, state variants.
4. `manifest.md` lists every planned icon/illustration with category, status, and notes.
5. Cross-references resolve: manifest items follow taxonomy names; usage cites style decisions.

## Do Not

- Do not render images here. Output is text spec.
- Do not allow ad-hoc icons outside the taxonomy.
- Do not skip accessibility rules in `usage.md`.
- Do not store binary asset files inside the spec directory.

## Tips

- Lock geometry rules early; downstream icon production breaks when style drifts mid-set.
- Use synonyms map to keep search/index aligned without renaming.
- For large sets, batch by category in `manifest.md` to track progress.
- Promote icon palette tokens up to product `palette.md` if shared.

## Hand Off

| Situation | Hand off to |
|---|---|
| Spec ready, asset production starts | Create design/build tasks citing icon names; worker uses `recipe-agent-team-worker-checkpoint` |
| Icon set will live inside product UI | Start a UI subdomain pass using `references/ui.md` for screen integration |
| Naming overlaps with existing component glossary | Run `recipe-agent-team-terminology-context` first |
