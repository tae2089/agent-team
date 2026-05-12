# Domain Model Interview

Use when names, ownership, bounded context, or business meaning are unclear.

Core lens: architecture should preserve coherent domain language inside each context.

## Interview

1. What domain concept does this module represent?
2. Are user terms, code terms, and artifact terms aligned?
3. Is this design crossing a bounded context?
4. Which rules belong inside the domain model rather than callers?
5. Which terms are overloaded or misleading?
6. Does the interface expose domain behavior or technical plumbing?
7. Should terminology context be updated before coding?

## Failure Signals

- Generic names hide business meaning.
- One module mixes concepts from separate contexts.
- Callers enforce domain rules in scattered conditionals.
- Downstream tasks use different terms for the same concept.

## Output

Record findings in the `## Interview Findings` table of `architecture-candidates.md` or `technical-design.md`:

| Lens | Finding | Design Impact |
| ---- | ------- | ------------- |
| Domain Model | canonical terms / bounded context or ownership boundary / domain rules moved behind the interface | terminology artifact updates / naming adjustment |
