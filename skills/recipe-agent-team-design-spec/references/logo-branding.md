# Marketing / Logo-Branding Reference

Use when discovery selects subdomain `logo-branding`. Produces brand identity specifications: concept, wordmark, symbol, usage rules, and color/typography palette. Text-based — no rendered logo files.

## Output Layout

```
_workspace/{run_id}/design/logo-branding/
├── concept.md      # brand essence, positioning, audience
├── wordmark.md     # text-only logotype specification
├── symbol.md       # icon/symbol/monogram specification
├── usage.md        # clear space, sizing, do/don't, placement rules
└── palette.md      # color tokens, typography ramp, tone of voice
```

## Pipeline Order

1. **concept.md** — brand essence, mission, voice, audience, positioning vs competitors.
2. **wordmark.md** — text-form logo spec: typeface, kerning intent, weight, casing.
3. **symbol.md** — non-text mark spec: shape, geometry, construction grid, monogram if any.
4. **usage.md** — clear space, minimum size, allowed/forbidden treatments, placement on backgrounds.
5. **palette.md** — color tokens (primary/secondary/neutral/accent), typography pairing, motion if any.

## Sub-Artifact Templates

### concept.md

```markdown
# Brand Concept

## Essence
- One-sentence promise: <"We help <audience> <verb> <outcome>.">
- Three adjectives: <e.g., bold, precise, warm>
- Anti-adjectives: <words the brand is not>

## Audience
- Primary: <segment>
- Secondary: <segment>
- Decision driver: <what makes them choose>

## Positioning
- Category: <where the brand sits>
- Differentiator: <what only this brand does>
- Competitor map: <2–4 competitors and relative position>

## Voice / Tone
- Tone words: <e.g., direct, optimistic>
- Avoid: <e.g., jargon, corporate filler>
- Example phrase: <one sentence that sounds like the brand>
```

### wordmark.md

```markdown
# Wordmark

## Letterforms
- Base typeface: <e.g., Inter Display Bold>
- Casing: <UPPERCASE | Title | lowercase>
- Tracking adjustment: <+/- units>
- Custom adjustments: <e.g., shortened crossbar on A>

## Geometry
- Cap height baseline: <unit>
- Stem weight: <unit or %>
- Optical corrections: <e.g., O slightly tall>

## Variants
- Primary horizontal
- Stacked
- Single-color (black)
- Reversed (white on dark)

## Forbidden
- Stretching, skewing, drop shadow, gradient, outline.
```

### symbol.md

```markdown
# Symbol / Mark

## Concept
- Idea: <e.g., abstract arrow + leaf for "growth + nature">
- Construction: <grid system, ratios>

## Geometry
- Base unit: <unit>
- Ratios: <e.g., 1:1.618>
- Stroke width: <unit>

## Variants
- Full color
- Single color (black, white)
- Outline-only (for embroidery, small sizes)

## Monogram (if applicable)
- Letterforms: <which letters>
- Lockup: <how they combine>

## Forbidden
- Replacing geometry, recoloring outside palette, decorating.
```

### usage.md

```markdown
# Usage Rules

## Clear Space
- Minimum padding: equal to symbol cap height × 0.5 on all sides.

## Minimum Sizes
| Use | Min size |
|---|---|
| Digital | 24px tall (symbol) |
| Print | 12mm tall (symbol) |
| Wordmark | 80px / 30mm wide minimum |

## Allowed Backgrounds
- Solid palette colors only.
- Photographic backgrounds: must use reversed variant on darkened overlay (>=40% black).

## Lockups
- Logo + tagline: tagline sits below at <ratio>.
- Co-branding: 1:1 weight, separator bar between marks.

## Do / Don't
- DO use approved variants only.
- DO respect clear space.
- DON'T outline, gradient, drop-shadow, rotate, recolor, place on cluttered background.
```

### palette.md

```markdown
# Color & Typography Palette

## Color

| Token | Hex | Role |
|---|---|---|
| brand.primary | #0F172A | Primary surface |
| brand.accent | #F97316 | Highlights, CTA |
| brand.neutral.0 | #FFFFFF | Background light |
| brand.neutral.9 | #0B0F14 | Text on light |
| brand.support | #22C55E | Success/positive |

## Typography

| Role | Family | Weight | Size baseline |
|---|---|---|---|
| display | Inter Display | 700 | 48px |
| heading | Inter | 600 | 24px |
| body | Inter | 400 | 16px |
| micro | Inter | 500 | 12px |

## Motion (if applicable)
- Logo reveal: 320ms ease-out
- Hover scale: 1.02
- Disallowed: bounce, rotate, flash

## Tone of Voice Pairing
- Concept tone words ↔ typography weight chart.
- Example: "bold" → display 700; "warm" → body 400 with generous line-height 1.6.
```

## Acceptance Criteria

A logo/branding artifact set passes when:

1. `concept.md` defines essence, audience, positioning, voice with anti-adjectives.
2. `wordmark.md` specifies typeface, casing, geometry, variants, and forbidden treatments.
3. `symbol.md` specifies idea, construction grid, geometry, variants, forbidden treatments.
4. `usage.md` provides clear space, minimum size table, allowed backgrounds, lockup rules, do/don't list.
5. `palette.md` defines tokens for color, typography, and motion paired with tone of voice.
6. Cross-references resolve: usage cites wordmark/symbol variants; palette tokens are referenced by all other files.

## Do Not

- Do not output rendered logo files. Spec is text + structured descriptions.
- Do not allow gradients, drop-shadows, or outlines unless explicitly in variants.
- Do not skip `usage.md` — without rules, downstream applications drift.
- Do not store binary art assets here.

## Tips

- Pair each tone word in `concept.md` with concrete typographic or color decisions.
- Build the symbol on a grid so reproductions stay consistent.
- Always define a single-color and reversed variant; real-world placements demand them.
- Promote palette tokens to a shared `tokens.md` style if logo doubles as product brand.
- When an upstream `design-system` run produced `_workspace/{run_id}/design/design-system/DESIGN.md`, cite tokens via `{group.name}` inside `palette.md` instead of redefining hex values; add only logo-specific tokens (e.g., reversed variants) as deltas.

## Hand Off

| Situation | Hand off to |
|---|---|
| Spec ready, art production starts | Create design/build tasks citing brand tokens; worker uses `recipe-agent-team-worker-checkpoint` |
| Brand voice conflicts with product copy | Run `recipe-agent-team-terminology-context` |
| Logo doubles as product brand and no shared token catalog exists | Pause; run `recipe-agent-team-design-interview` with `subdomain: design-system`, then resume `logo-branding` via `recipe-agent-team-design-spec` citing the new `DESIGN.md` |
| Logo will appear inside product UI | Start a UI subdomain pass using `references/ui.md` for screen integration |
