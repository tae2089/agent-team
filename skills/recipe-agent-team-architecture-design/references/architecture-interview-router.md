# Architecture Interview Router

Use this router when architecture candidates have meaningful tradeoffs. These interviews adapt scenario-based architecture evaluation, ADR-style decision capture, and quality-attribute questioning into lightweight candidate reviews for agent-team workflows.

Load only two to four unique interview files by default. If the same interview appears in multiple default sets, load it once. Do not load every interview unless the user explicitly asks for a broad architecture audit.

## Choose Interviews

| Risk Signal | Interview File |
| --- | --- |
| callers know volatile details, expected changes leak through interfaces | `interviews/information-hiding.md` |
| abstraction feels shallow, pass-through, or low leverage | `interviews/deep-modules.md` |
| behavior is hard to test or variation is hidden in callers | `interviews/seams-and-tests.md` |
| migration could be risky, broad, or hard to stage | `interviews/evolutionary-refactoring.md` |
| domain terms, ownership, or model boundaries are unclear | `interviews/domain-model.md` |
| dependencies point both ways or stable policy depends on volatile detail | `interviews/dependency-direction.md` |
| quality attributes are vague or unmeasured | `interviews/quality-scenarios.md` |
| design touches agent-team run/task/message/artifact boundaries | `interviews/runtime-fit.md` |

## Default Sets

For most codebase architecture work, start with:

- `deep-modules.md`
- `information-hiding.md`
- `evolutionary-refactoring.md`

For agent-team runtime or harness work, start with:

- `runtime-fit.md`
- `information-hiding.md`
- `dependency-direction.md`

For domain-heavy work, start with:

- `domain-model.md`
- `information-hiding.md`
- `quality-scenarios.md`

When combining sets, deduplicate interview files before loading them.

## Output Contract

Record interview findings in `architecture-candidates.md` or `technical-design.md`:

```markdown
## Interview Findings

| Lens | Finding | Design Impact |
| ---- | ------- | ------------- |
```

Keep each finding short. The goal is to explain selection pressure, not to paste the full interview.

## Stop Rule

Stop adding interviews when:

- candidate recommendation is clear
- remaining risks are already captured
- another interview would not change module/interface shape, migration plan, or acceptance criteria
