# Information Hiding Interview

Use when callers know too much about a volatile decision, configuration, ordering rule, data shape, storage detail, or workflow step.

Core lens: a useful module hides a design decision likely to change.

## Interview

1. What design decision should this module hide?
2. Which callers currently depend on that decision?
3. What future change should be absorbed inside the module?
4. What knowledge must remain visible to callers?
5. What knowledge is leaking across the interface now?
6. Which interface detail is too implementation-specific?
7. What would still force callers to change after this design?

## Failure Signals

- The module only groups code by noun.
- The interface exposes the same details as the implementation.
- Every likely future change still touches callers.
- Callers must preserve ordering, flags, or internal state rules manually.

## Output

Record findings in the `## Interview Findings` table of `architecture-candidates.md` or `technical-design.md`:

| Lens | Finding | Design Impact |
| ---- | ------- | ------------- |
| Information Hiding | hidden decision / leaked knowledge / caller impact | interface adjustment / migration implication |
