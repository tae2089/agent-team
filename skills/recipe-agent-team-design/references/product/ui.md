# Product / UI Reference

Use when discovery selects domain `product/ui`. Produces interface specifications for screens, flows, components, and design tokens. Framework-agnostic — no JSX/Vue/Svelte syntax.

## Output Layout

```
_workspace/{run_id}/design/product/ui/
├── flows.md         # information architecture + user flows
├── wireframes.md    # screen layouts (ASCII/markdown)
├── components.md    # reusable component specs
└── tokens.md        # design tokens (color/spacing/typography/theme)
```

## Pipeline Order

Run sub-steps in order. Skip a step only when its artifact already exists and inputs unchanged.

1. **flows.md** — information architecture, navigation graph, primary user journeys, state transitions, error/empty states. Defines *what screens exist and how users move between them*.
2. **wireframes.md** — per-screen layout, region partitions, content hierarchy, responsive breakpoints. Defines *what each screen contains*.
3. **components.md** — reusable components extracted from wireframes, props, variants, states, accessibility notes. Defines *what building blocks exist*.
4. **tokens.md** — color palette, spacing scale, typography ramp, theme variants (light/dark), motion. Defines *visual language constants*.

Later sub-artifacts reference earlier ones (wireframes cite flow IDs, components cite wireframe regions, tokens are referenced by components).

## Sub-Artifact Templates

### flows.md

```markdown
# Information Architecture & User Flows

## Site Map
- [F1] Home
  - [F1.1] Dashboard
  - [F1.2] Settings
    - [F1.2.1] Profile
    - [F1.2.2] Notifications
- [F2] Auth
  - [F2.1] Login
  - [F2.2] Signup

## Primary Flows

### Flow: First-time onboarding
1. [F2.2] Signup → submit credentials
2. [F1] Welcome screen → confirm role
3. [F1.1] Dashboard (empty state)

### State Transitions
| From | Event | To |
|---|---|---|
| [F2.1] Login (idle) | submit valid | [F1.1] Dashboard |
| [F2.1] Login (idle) | submit invalid | [F2.1] Login (error) |

### Error / Empty States
- [F1.1] Dashboard: empty, loading, error, success.
- [F2.1] Login: idle, submitting, error.
```

### wireframes.md

```markdown
# Wireframes

## [F1.1] Dashboard

Breakpoints: mobile (<640px), tablet (640–1024px), desktop (>1024px).

### Desktop layout

\`\`\`
+--------------------------------------------+
| TopBar: logo | search | user-menu          |
+----------+---------------------------------+
| SideNav  | Main                            |
| - Home   | [Card: KPI A] [Card: KPI B]     |
| - Reports| [Table: recent activity]        |
| - Settings| [Empty state if no data]       |
+----------+---------------------------------+
\`\`\`

### Mobile layout

\`\`\`
+----------------------+
| TopBar (collapsed)   |
+----------------------+
| Main                 |
| [Card: KPI A]        |
| [Card: KPI B]        |
| [Table: activity]    |
+----------------------+
| BottomNav            |
+----------------------+
\`\`\`

Regions: TopBar, SideNav, BottomNav, Main, Card, Table, EmptyState.
```

### components.md

```markdown
# Components

## Card

Used in: Dashboard (KPI tiles), Reports list.

| Prop | Type | Required | Default | Description |
|---|---|---|---|---|
| title | string | yes | — | Card heading |
| value | string \| number | yes | — | Primary metric |
| trend | "up" \| "down" \| "flat" | no | "flat" | Arrow indicator |
| onClick | () => void | no | — | Click handler |

States: default, hover, loading, error.
Variants: compact, expanded.
A11y: role="article", keyboard focusable when clickable, aria-label from title.

## Table

Used in: Dashboard (recent activity), Reports.

| Prop | Type | Required | Default |
|---|---|---|---|
| rows | Row[] | yes | — |
| columns | Column[] | yes | — |
| emptyState | ReactNode | no | <Empty/> |

States: idle, loading, error, empty.
Sort/filter/pagination requirements documented per usage.
```

### tokens.md

```markdown
# Design Tokens

## Color

| Token | Light | Dark | Usage |
|---|---|---|---|
| color.bg.surface | #FFFFFF | #0B0F14 | Card/panel background |
| color.bg.canvas | #F5F7FA | #050810 | Page background |
| color.fg.primary | #111827 | #F8FAFC | Body text |
| color.fg.muted | #6B7280 | #9CA3AF | Secondary text |
| color.accent | #2563EB | #60A5FA | Primary action |
| color.danger | #DC2626 | #F87171 | Error/destructive |

## Spacing (4px base)

| Token | Value |
|---|---|
| space.1 | 4px |
| space.2 | 8px |
| space.3 | 12px |
| space.4 | 16px |
| space.6 | 24px |
| space.8 | 32px |

## Typography

| Token | Family | Size | Weight | Line height |
|---|---|---|---|---|
| text.display | Inter | 32px | 700 | 1.2 |
| text.heading | Inter | 20px | 600 | 1.3 |
| text.body | Inter | 14px | 400 | 1.5 |
| text.caption | Inter | 12px | 400 | 1.4 |

## Motion

| Token | Duration | Easing |
|---|---|---|
| motion.fast | 120ms | ease-out |
| motion.standard | 200ms | ease-in-out |
| motion.slow | 320ms | ease-in-out |

## Theme Variants

- light (default)
- dark
- high-contrast (optional)
```

## Acceptance Criteria

A UI design artifact set passes when:

1. `flows.md` lists every screen with a stable ID, primary flows, and key states (loading/empty/error/success).
2. `wireframes.md` covers every screen ID from flows.md with at least one breakpoint layout and labeled regions.
3. `components.md` extracts every region/element reused across two or more wireframes, with props/states/variants/a11y notes.
4. `tokens.md` defines color, spacing, typography, and motion tokens that the wireframes/components reference.
5. Cross-references resolve: wireframe regions cite token names; component specs cite wireframe regions; flows cite wireframe IDs.
6. No framework-specific syntax (no JSX/Vue templates/Svelte). Pure markdown + ASCII.

## Do Not

- Do not write implementation code (JSX, CSS, etc.). This reference ends at the spec.
- Do not embed pixel-perfect designs — wireframes are structural, not visual mockups.
- Do not bypass earlier sub-artifacts; later ones depend on them.
- Do not duplicate component specs inside wireframes; reference by name.
- Do not store assets (images, fonts) inside the output directory; link externally if needed.

## Tips

- Use stable screen IDs (`[F1.1]` style) so downstream coding workers can grep.
- Mirror flow IDs in task titles for traceability (e.g., "Implement [F1.1] Dashboard").
- For complex states, add a small state diagram in `flows.md` rather than burying transitions in prose.
- Component props should be minimal — defer optionality decisions to consuming screens.
- Reuse token names rather than literal hex/px values in `components.md`.

## Hand Off

| Situation | Hand off to |
|---|---|
| UI spec ready, implementation needed | Create coding tasks citing flow/component/token IDs; worker uses `recipe-agent-team-worker-checkpoint` |
| UI spec depends on missing backend contract | Pause and run `recipe-agent-team-architecture-design` |
| Terms inside spec ambiguous | Stop and run `recipe-agent-team-terminology-context` first |
