# Game Environment / Level Design Reference

Pattern for designs whose `Output` describes a level, zone, biome, or game environment: concept, layout, set dressing, lighting/mood, and asset manifest. Cite this file as a `Pattern Hint` in the brief when applicable. Text-based — no rendered images.

## Output Layout

```
_workspace/{run_id}/design/{slug}/   # e.g., {slug}=forest-zone-1 or boss-arena
├── concept.md       # narrative role, mood, pacing intent
├── layout.md        # macro map, regions, navigation, encounter beats
├── biome.md         # terrain, vegetation, weather, biome rules
├── lighting.md      # time of day, color palette, atmosphere, key lights
├── set-dressing.md  # props, points of interest, environmental storytelling
└── assets.md        # asset manifest, formats, naming, delivery requirements
```

One brief per level/zone. Pick a stable `{slug}` during Output Capture (e.g., `lv-forest-01`, `boss-arena`).

## Pipeline Order

1. **concept.md** — narrative role of the environment, mood arc, player emotion target, place in overall world.
2. **layout.md** — macro shape, regions, critical path, side paths, encounter beats, sightlines, navigation affordances.
3. **biome.md** — terrain type, vegetation, weather, biome-specific rules (e.g., underwater stamina drain).
4. **lighting.md** — time-of-day plan, palette, key light direction, atmospheric effects (fog, dust, rain).
5. **set-dressing.md** — prop list, points of interest, environmental storytelling beats, hidden details.
6. **assets.md** — concrete deliverable list: meshes, textures, foliage, props, audio cues, VFX hooks.

Each sub-artifact references the previous one (layout cites concept beats; biome cites layout regions; lighting cites biome time-of-day; set-dressing cites layout points-of-interest; assets cites everything).

## Sub-Artifact Templates

### concept.md

```markdown
# Environment Concept — [level_id]

## Identity
- Name: <name>
- Place in world: <region, faction, era>
- Player level / progression slot: <e.g., act 1 mid>

## Narrative Role
- One-line: <"The <adjective> <noun> where <event> happened.">
- Story beat: <what advances here>
- Faction presence: <who lives/contests this place>

## Mood Arc
- Entry: <emotion> — first impression
- Mid: <emotion> — discovery / tension
- Exit: <emotion> — payoff / departure

## Player Fantasy
- Verb: <explore | siege | infiltrate | survive>
- Threat profile: <ambush, environmental, boss>
- Reward profile: <loot, lore, traversal unlock>

## References
- Inspirations: <links / descriptions>
- Avoid: <anti-refs to keep level distinct from neighbors>
```

### layout.md

```markdown
# Layout — [level_id]

## Macro Shape
- Topology: <linear | hub-and-spoke | open | gauntlet | maze>
- Dimensions: <approx footprint, e.g., 600×400m>
- Verticality: <flat | layered (N tiers)>

## Regions

| ID | Name | Function | Connections |
|---|---|---|---|
| R1 | Approach | tension build | → R2 |
| R2 | Threshold | first contact | R1 → R3, R1 → R4 |
| R3 | Heart | climax encounter | R2 → R5 |
| R4 | Side cache | optional reward | R2 |
| R5 | Exit | resolution | R3 → next level |

## Critical Path
- R1 → R2 → R3 → R5

## Side Paths
- R2 → R4 (optional)

## Encounter Beats

| Region | Beat | Type |
|---|---|---|
| R1 | scout patrol | combat-light |
| R2 | ambush from cover | combat-medium |
| R3 | mini-boss | combat-heavy |
| R4 | environmental puzzle | non-combat |
| R5 | exfil with timer | tension |

## Sightlines & Affordances
- Key vista at R2 → R3 transition.
- Hidden alcove in R3 with lore prop.
- Reverse traversal blocked at R3 → R5 (one-way drop).

## Navigation Rules
- No teleporters except R5 exit.
- Backtrack allowed until R3 climax begins.
```

### biome.md

```markdown
# Biome — [level_id]

## Terrain
- Ground: <e.g., wet stone, peat, packed snow>
- Slope rules: <max walkable angle>
- Hazards: <pits, slick surfaces, breakable floors>

## Vegetation
- Canopy: <e.g., dense pine, sparse oak>
- Ground cover: <ferns, moss, snow drifts>
- Density map per region:

| Region | Density | Notes |
|---|---|---|
| R1 | sparse | open sightlines |
| R3 | dense | claustrophobic for boss |

## Weather
- Default state: <e.g., overcast rain>
- Variants: <storm spike, fog rolls>
- Gameplay effect: <visibility −X%, footstep sound +Y%>

## Biome Rules
- Stamina: <e.g., wet ground +10% drain>
- Audio: <muffled by snow>
- VFX: <constant ember motes for volcano biome>
```

### lighting.md

```markdown
# Lighting — [level_id]

## Time of Day
- Locked: <dawn | day | dusk | night | dynamic>
- Sun angle: <e.g., 25° from horizon, NW>

## Palette

| Token | Hex | Role |
|---|---|---|
| light.key | #FFD08A | Sun warmth |
| light.fill | #6B8FB3 | Sky ambient |
| light.bounce | #3C5470 | Ground bounce |
| light.fog | #88A2B0 | Atmospheric tint |

## Key Lights per Region

| Region | Light source | Intent |
|---|---|---|
| R1 | low sun through canopy | tension |
| R3 | overhead skylight through ruin | reveal boss silhouette |
| R5 | gold-hour open sky | relief |

## Atmosphere
- Fog density: <e.g., 0.02 exponential>
- Particulates: <dust motes, embers, rain streaks>
- Shadow softness: <hard | soft>

## Audio-Light Coupling
- Thunder ↔ flash sync window: <ms>
- Lantern flicker on combat enter.
```

### set-dressing.md

```markdown
# Set Dressing — [level_id]

## Points of Interest

| ID | Region | Type | Story Beat |
|---|---|---|---|
| POI-1 | R2 | wrecked cart | hint of prior battle |
| POI-2 | R3 | shrine | lore prop, optional read |
| POI-3 | R4 | hidden cache | reward |
| POI-4 | R5 | banner | faction presence on exit |

## Prop Categories

| Category | Examples | Density |
|---|---|---|
| structural | walls, beams, pillars | as layout dictates |
| narrative | banners, books, corpses | sparse, deliberate |
| ambient | barrels, crates, tools | medium |
| natural | rocks, logs, stumps | per biome |

## Environmental Storytelling
- POI-1: cart wheel broken outward → ambush from R2 cover (foreshadow).
- POI-2: shrine offerings dated → faction recently passed through.
- POI-4: torn banner → defeat narrative on level exit.

## Forbidden
- No interactable props that block critical path.
- No prop occluding key sightlines defined in layout.md.
```

### assets.md

```markdown
# Assets — [level_id]

## Geometry
| Asset | Format | Notes |
|---|---|---|
| terrain mesh | .fbx | LOD0–2, heightmap source preserved |
| structure kit | .fbx parts | modular walls/roofs/floors |
| props | .fbx | shared prop library where possible |

## Textures
| Set | Resolution | Channels |
|---|---|---|
| terrain | 4096² tile | baseColor, normal, ORM |
| structure | 2048² | PBR set |
| props | 1024² | PBR set |

## Foliage
| Asset | Type | LOD |
|---|---|---|
| canopy tree | mesh + billboard | 4 LODs |
| ground fern | mesh + impostor | 2 LODs |

## VFX Hooks
- POI-1: dust kick on approach.
- R3 entrance: skylight god rays.
- Weather: rain streaks, puddle ripples.

## Audio Cues
| Cue | Trigger |
|---|---|
| ambient-forest | level enter |
| storm-rumble | weather variant |
| combat-stinger | R3 boss enter |

## Delivery
- File naming: `{level_id}_{category}_{name}.{ext}` (e.g., `lv-forest-01_terrain_main.fbx`).
- Source under `{Output Path}/source/` (if any).
- Binary assets do not commit to repo; link externally.
```

## Acceptance Criteria

An environment design set passes when:

1. `concept.md` defines narrative role, mood arc, player fantasy.
2. `layout.md` provides macro shape, regions with IDs, critical path, side paths, encounter beats, sightlines.
3. `biome.md` defines terrain, vegetation, weather, biome rules.
4. `lighting.md` defines time-of-day, palette tokens, key lights per region, atmosphere.
5. `set-dressing.md` lists POIs with story beats, prop categories, environmental storytelling rules.
6. `assets.md` enumerates every deliverable with format and naming.
7. Cross-references resolve: layout regions cite encounter beats; biome cites regions; lighting cites regions; set-dressing cites POIs in regions; assets cites everything upstream.

## Do Not

- Do not render images. Output is text spec.
- Do not skip layout.md — encounter beats anchor everything downstream.
- Do not place props that block the critical path.
- Do not use lighting palettes that contradict biome mood.
- Do not store binary assets inside the spec directory.

## Tips

- One level per subdirectory; copy a baseline level to bootstrap variants.
- Stable IDs (R1, POI-1, level_id) appear in task titles and asset names.
- Pair every POI with one story beat — empty POIs become visual noise.
- Promote shared prop categories to a project-wide prop library doc when patterns emerge.
- Verify reverse-traversal rules in layout.md early; backtracking decisions cascade through encounter design.

## Hand Off

| Situation | Hand off to |
|---|---|
| Spec ready, environment art/build starts | Create build tasks citing `level_id`; worker uses `recipe-agent-team-worker-checkpoint` |
| Characters needed in this environment | Run separate Output passes citing `character.md` as a Pattern Hint |
| Encounter mechanics undefined | Pause and run `recipe-agent-team-planning-grill` for combat rules |
| Naming conflicts with other levels | Run `recipe-agent-team-terminology-context` first |
