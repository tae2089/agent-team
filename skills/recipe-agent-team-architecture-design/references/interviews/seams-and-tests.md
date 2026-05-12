# Seams And Tests Interview

Use when behavior is hard to test, variation is embedded in callers, or a candidate proposes a seam, adapter, mock, fake, plugin, or provider abstraction.

Core lens: a seam is useful when behavior can vary without editing callers.

## Interview

1. Where can behavior vary without editing callers?
2. What concrete adapters exist now or are expected soon?
3. Are tests crossing the same interface as production callers?
4. Does the seam make behavior easier to verify?
5. Is the seam hiding real variation or only supporting a hypothetical future?
6. Can tests verify behavior without reaching behind the interface?
7. What production path proves the seam is correctly placed?

## Failure Signals

- There is only one adapter and no near-term variation.
- Tests need privileged access to internals.
- The seam exists only to make mocking easier.
- Callers still branch on concrete implementation details.

## Output

Record findings in the `## Interview Findings` table of `architecture-candidates.md` or `technical-design.md`:

| Lens | Finding | Design Impact |
| ---- | ------- | ------------- |
| Seams And Tests | seam location / concrete adapters / hypothetical seam risk | test strategy impact / interface adjustment |
