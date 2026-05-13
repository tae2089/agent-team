# Design Brief Template

Use this exact heading and table structure for `_workspace/{run_id}/design/{slug}/design-brief.md`, or inline in the final response when no `RUN_ID` exists. Empty sections are forbidden; record a refusal, risk, or upstream citation instead.

```markdown
# Design Brief

## Output
- Description: <one-paragraph plain-text description of what is being designed>
- Slug: <kebab-case identifier, e.g., dashboard, heritage-brand, token-catalog>
- Output Path: `_workspace/{run_id}/design/<slug>/` (or "inline final response" when no RUN_ID)
- Pattern Hints: <comma-separated basenames from the spec Pattern Library, e.g., `ui.md`, `design-system.md`; or "none">

## Core Story
- Specific user: <one person>
- Before / After / Pivotal moment: <concrete>

## Subject
- One-line summary / Audience / Problem-feeling

## Specificity
- Concrete signals / Examples that earn the adjective / Anti-examples that fake it
- Rejected adjectives: <adjective> — <reason>

## Tensions
| Axis A | Axis B | Chosen | Accepted cost |
| --- | --- | --- | --- |

## Assumptions
| Statement | Evidence | Breakage impact | Risk? |
| --- | --- | --- | --- |

## Priority Ranking
1. <must-keep>
2. <second>
3. <third>
4. <fourth>
5. <fifth>

## First 5 Seconds
- <single sentence>

## Failure Modes
| Scenario | Root cause | Earliest signal | Mitigation seed |
| --- | --- | --- | --- |

## Constraints
- Platform / Style / Brand-palette / Budget-deadline

## References
- Emulate / Avoid

## Success Criteria
- <derived from priority ranking and first 5 seconds>

## Upstream Inputs
- Planning / Terminology / Prior design: <paths or "n/a">
- Prior pass briefs/artifacts in same RUN_ID: <paths or "n/a">

## Interview Log
- Phase identifiers: use exact Discovery Phase names from SKILL.md (Output Capture, Context & Story, Specificity Drill, Trade-off Surface, Assumption Probe, Forced Choice & Priority, Failure Rehearsal).
- Gap scores: <applicable phase name = 0|1|2; skipped phase = reason>
- Executed order: <phase names in order>
- Re-scan triggers: <none | description>
- Phase captures: <probe used -> insight, or skipped with citation>

## Routed By
- recipe-agent-team-design-interview
- Reason: <Output Capture summary -> chosen Output description>
- Requester confirmation: <quote | quick-draft requested>
- Next step: hand off to `recipe-agent-team-design-spec` with this brief path
```
