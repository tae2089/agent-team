# Game / Character Reference

Use when discovery selects domain `game/character`. Produces character design specifications: concept, visual identity, mechanics, animation requirements, and asset manifest. Text-based — no rendered images.

## Output Layout

```
_workspace/{run_id}/design/game/character/{character_id}/
├── concept.md       # lore, personality, role, archetype
├── visual.md        # silhouette, color palette, costume, props, references
├── mechanics.md     # stats, skills, abilities, balance notes
├── animation.md     # key poses, required animations, transitions
└── assets.md        # asset manifest, formats, naming, delivery requirements
```

One subdirectory per character. Use stable `character_id` (e.g., `c-knight-01`).

## Pipeline Order

1. **concept.md** — who the character is, role in the game, narrative hook, personality, archetype, target audience appeal.
2. **visual.md** — silhouette intent, palette, costume/equipment, signature props, references and contrasts to existing roster.
3. **mechanics.md** — stats vector, skills, abilities, cooldowns, balance rationale, counter relationships.
4. **animation.md** — required animation list, key poses, transitions, idle/combat/emote behavior.
5. **assets.md** — concrete deliverable list for art/audio pipeline: model, textures, rig, animations, VFX hooks, audio cues, formats, naming.

Each sub-artifact references the previous one (visual cites concept role; mechanics cites concept archetype; animation cites mechanics abilities; assets cites everything).

## Sub-Artifact Templates

### concept.md

```markdown
# Character Concept — [character_id]

## Identity
- Name: <name>
- Role: <tank | dps | support | utility>
- Archetype: <e.g., honorable knight, rogue trickster>
- Faction / Origin: <group>
- Age / Stature: <range>

## Personality
- Core traits: <3–5 adjectives>
- Voice / mannerism: <1–2 sentences>
- Motivation: <what drives them>
- Flaw / weakness: <story hook>

## Narrative Hook
- Backstory (3–5 sentences).
- Relation to existing roster: <ally, rival, neutral>.
- Player fantasy: <what playing this character feels like>.

## Target Appeal
- Player segment: <e.g., players who like methodical combat>
- Comparable references: <external games/characters>
- Anti-references: <what to avoid>
```

### visual.md

```markdown
# Visual Design — [character_id]

## Silhouette
- Read at distance: <shape language — e.g., wide shoulders + narrow waist + tall plume>
- Iconic features: <3 distinctive elements>

## Palette

| Token | Hex | Role |
|---|---|---|
| primary | #2C3E50 | Main armor |
| accent | #C0392B | Cape lining, energy effects |
| neutral | #BDC3C7 | Trim, leather |
| metal | #95A5A6 | Weapon, plate |

## Costume / Equipment
- Head: <description>
- Body: <description>
- Hands: <description>
- Legs / Feet: <description>
- Weapon: <description>
- Signature prop: <description>

## Style Constraints
- Stylization: <realistic | stylized | toon | pixel>
- Material rules: <e.g., no chrome, fabric matte>
- Wear/age treatment: <pristine | weathered>

## References
- Inspirations: <links or descriptions>
- Avoid: <contrasts to maintain distinctness from existing roster>
```

### mechanics.md

```markdown
# Mechanics — [character_id]

## Stats (normalized 1–10)

| Stat | Value | Note |
|---|---|---|
| Health | 8 | High durability |
| Damage | 5 | Moderate output |
| Mobility | 3 | Slow |
| Range | 2 | Melee focus |
| Utility | 6 | Crowd control |

## Skills

### Skill 1 — Shield Bash
- Type: active, melee
- Cooldown: 6s
- Effect: stun 1.5s, damage 30
- Counter: dodge, ranged kite

### Skill 2 — Phalanx Wall
- Type: active, defensive aura
- Cooldown: 18s
- Effect: party damage reduction 25% for 4s
- Counter: armor-pierce abilities

### Skill 3 — Last Stand
- Type: ultimate
- Cooldown: 90s
- Effect: refuses death once per fight, heals 40%
- Counter: burst before threshold

## Balance Rationale
- Power budget: high survivability, low mobility, low burst.
- Counter relationships: outscaled by ranged sustained DPS; counters melee assassins.
- Skill ceiling: <low | medium | high> — what mastery looks like.
```

### animation.md

```markdown
# Animation — [character_id]

## Required Animations

| Animation | Frames / Duration | Notes |
|---|---|---|
| idle | loop 4s | weight on back foot, helmet tilts |
| walk | loop 1s | grounded, heavy |
| run | loop 0.7s | armor jostle |
| attack-light | 0.4s | overhead chop |
| attack-heavy | 0.9s | windup + slam |
| skill-1 shield-bash | 0.6s | shield raise → forward step |
| skill-2 phalanx-wall | 1.0s | plant feet → ground shockwave |
| ultimate-last-stand | 1.4s | kneel → roar → glow up |
| hit-reaction-light | 0.3s | minimal flinch |
| hit-reaction-heavy | 0.6s | stagger back one step |
| death | 1.2s | collapse forward |
| emote-victory | loop 2s | shield raised |
| emote-taunt | loop 2.5s | beckoning hand |

## Key Poses
- Stance: feet shoulder-width, shield braced.
- Charge: shield forward, body lowered.
- Ultimate apex: shield ground-planted, free hand to sky.

## Transition Rules
- idle ↔ walk: 6 frames blend.
- attack chains: light → light → heavy within 0.8s.
- Cancel windows: skill windups cancel into block on input.
```

### assets.md

```markdown
# Assets — [character_id]

## 3D / Visual
| Asset | Format | Resolution | Notes |
|---|---|---|---|
| body mesh | .fbx | <= 25k tris | LOD0, retopo from sculpt |
| body texture | PBR set | 2048² | baseColor, normal, ORM |
| weapon mesh | .fbx | <= 4k tris | separate skinning |
| rig | .fbx | humanoid | match team standard skeleton |

## Animation
| Asset | Format | Length |
|---|---|---|
| animation set | .fbx clips | per animation.md table |

## VFX / FX Hooks
- shield-bash: impact spark + dust ring
- phalanx-wall: ground hex pattern + dome shimmer
- last-stand: aura ring + screen-edge tint

## Audio
| Cue | Trigger |
|---|---|
| footstep-heavy | per stride |
| shield-bash-impact | skill 1 hit |
| phalanx-cast | skill 2 cast |
| ultimate-roar | skill 3 cast |
| death | death anim start |

## Delivery
- File naming: `{character_id}_{asset}_{lod}.{ext}` (e.g., `c-knight-01_body_lod0.fbx`).
- Source under `_workspace/{run_id}/design/game/character/{character_id}/source/` (if any).
- External binary assets do not commit to repo; link externally.
```

## Acceptance Criteria

A game character design set passes when:

1. `concept.md` defines role, archetype, personality, motivation, and player fantasy.
2. `visual.md` provides silhouette intent, palette tokens, costume parts, and style constraints with refs and anti-refs.
3. `mechanics.md` lists stats, all skills (cooldown/effect/counter), and balance rationale.
4. `animation.md` covers every state and ability with duration/notes.
5. `assets.md` enumerates every deliverable with format and naming so art/audio pipeline can plan.
6. Cross-references resolve: animation cites mechanics, assets cites animation + visual.

## Do Not

- Do not render images here. Output is text spec.
- Do not skip concept/visual to jump to mechanics; balance rationale needs role context.
- Do not invent stats outside the team's normalized scale.
- Do not omit counter relationships in `mechanics.md`.
- Do not store binary art assets inside the spec directory.

## Tips

- One character per subdirectory keeps multi-character runs tidy.
- Use `character_id` consistently — task titles, file names, asset names.
- For roster expansions, copy an existing character directory as scaffolding and diff.
- Keep palette tokens reusable; promote shared palettes to a roster-wide `palette.md` when patterns emerge.

## Hand Off

| Situation | Hand off to |
|---|---|
| Spec ready, art/animation pipeline starts | Create coding/art tasks citing `character_id`; worker uses `recipe-agent-team-worker-checkpoint` |
| Mechanics need re-balancing against existing roster | Pause and run `recipe-agent-team-planning-grill` |
| Naming conflicts with existing characters | Run `recipe-agent-team-terminology-context` first |
