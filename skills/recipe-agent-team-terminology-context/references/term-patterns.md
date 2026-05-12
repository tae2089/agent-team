# Terminology Patterns

Use this reference when deciding whether terms should be merged, split, renamed, or left as aliases.

## Canonical Term Criteria

Prefer a canonical term when it is:

- already used in a public interface, CLI, API, schema, or durable artifact
- clear to both users and implementers
- specific enough to avoid overloaded interpretation
- stable enough for future tasks to reference
- aligned with existing project vocabulary

Do not choose a new canonical term just because it sounds cleaner. Rename only when the current term is misleading or actively causing errors.

## Alias Pattern

Treat terms as aliases when:

- they refer to the same concept in the same context
- switching between them would not change implementation behavior
- one term is user-facing and another is code-facing

Record the canonical term and keep the alias mapping.

## Overloaded Term Pattern

Split a term when:

- different users or files use it for different concepts
- implementation would differ depending on the meaning
- the term crosses workflow phases with different responsibilities

Record each scoped meaning separately. Add "Avoid Confusing With" notes.

## Vague Noun Pattern

Challenge terms such as manager, handler, processor, context, service, utility, common, core, data, thing, item, object, helper, and flow.

Keep a vague noun only when the repository has an established precise meaning for it. Otherwise ask what behavior or responsibility it owns.

## Public Interface Bias

Prefer existing public names when they are not misleading:

- CLI command and flag names
- API routes and request fields
- schema and table names
- exported type or package names
- durable artifact names

If the public name is misleading, record the mismatch and route renaming through planning or architecture design.

## Multilingual And Mixed-Language Pattern

Preserve user-language terms and code-language terms separately when the project mixes languages, for example Korean user vocabulary with English code names.

Choose the canonical term by audience:

- Use the public code/API/doc term when downstream implementation or review needs exact file, symbol, command, or schema matching.
- Use the user/domain term when it captures business meaning better than the code term.
- Keep both when one is user-facing and the other is code-facing.

Record the mapping instead of translating away useful evidence.

Example:

```markdown
| User Term | Canonical Term | Notes |
| --------- | -------------- | ----- |
| 산출물 | artifact | Use `artifact` in task metadata; keep "산출물" in user-facing rationale. |
```

Avoid romanizing or translating terms unless the repo already does so. If a translation changes meaning, ask the user.

## Mixed Granularity Pattern

Do not merge terms that operate at different levels:

- workflow phase vs task
- artifact vs state record
- user goal vs implementation detail
- module vs function
- domain concept vs storage schema

When levels differ, map them in `context-map.md` instead of forcing one canonical term.

## Question Pattern

Ask a terminology question only when local context cannot resolve the boundary:

```text
Term: ...
Observed uses: ...
Recommended canonical term: ...
Risk if wrong: ...
Question: ...
```

## Output Pattern

Each terminology decision should record:

- canonical term
- aliases
- meaning
- source evidence
- avoid-confusing-with term
- downstream artifact or task impact

Use this compact decision format in `terminology.md`:

```markdown
## Decisions

| ID | Decision | Reason | Impact |
| -- | -------- | ------ | ------ |
| term-001 | Use "compound learning" as the canonical term; keep "compound" as an alias. | The full term is clearer in task contracts and avoids confusing it with integration. | Update planning and review task contracts to reference `compound learning`. |
```

Use this row shape for canonical terms:

```markdown
| Term | Meaning | Use When | Avoid Confusing With | Source |
| ---- | ------- | -------- | -------------------- | ------ |
| compound learning | Capturing reusable lessons from completed workflow evidence. | Referring to post-run reusable knowledge capture. | compound/integration phase | User decision term-001 |
```

For unresolved terminology, use:

```markdown
## Open Questions

- ID: term-q001
  Term: "compound"
  Ambiguity: Could mean integration phase or learning capture.
  Needed From: user
  Blocks: review task vocabulary
```

## Input To Output Example

Input:

```text
User says: "After review, we need compound so the next agent remembers the fix."
Existing artifacts mention: "compound-learning.md"
Code/docs mention: "run summary", "task evidence", "review artifact"
Planned consumers: architecture-design task and review task
```

Output in `terminology.md`:

```markdown
# Terminology

## Canonical Terms

| Term | Meaning | Use When | Avoid Confusing With | Source |
| ---- | ------- | -------- | -------------------- | ------ |
| compound learning | Capturing reusable lessons from completed workflow evidence. | Referring to post-run reusable knowledge capture. | compound/integration phase | User phrase + artifact `compound-learning.md` |

## User Terms

| User Term | Canonical Term | Notes |
| --------- | -------------- | ----- |
| compound | compound learning | Alias only when discussing post-run learning. |

## Code Terms

| Code Term | Canonical Term | Location | Notes |
| --------- | -------------- | -------- | ----- |
| compound-learning.md | compound learning | `_workspace/{run_id}/compound-learning.md` | Artifact name supports canonical term. |

## Decisions

| ID | Decision | Reason | Impact |
| -- | -------- | ------ | ------ |
| term-001 | Use "compound learning" as canonical; keep "compound" as an alias. | Avoids confusing learning capture with integration. | Downstream tasks should use `compound learning`. |

## Open Questions

None.
```

Because the planned consumers are two named tasks, also write `context-map.md`:

```markdown
# Context Map

## Agent Task Vocabulary

| Agent Or Task | Required Terms | Forbidden Or Ambiguous Terms |
| ------------- | -------------- | ---------------------------- |
| architecture-design task | compound learning | compound |
| review task | compound learning | compound |

## Artifact References

| Artifact | Terms It Defines | Terms It Consumes |
| -------- | ---------------- | ----------------- |
| `_workspace/{run_id}/terminology.md` | compound learning |  |
| `_workspace/{run_id}/compound-learning.md` |  | compound learning |
```
