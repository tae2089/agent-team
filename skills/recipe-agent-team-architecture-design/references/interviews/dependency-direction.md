# Dependency Direction Interview

Use when modules call each other both ways, stable policy depends on volatile detail, or a candidate changes ownership boundaries.

Core lens: dependency direction should protect stable policy from volatile implementation details.

## Interview

1. Which module represents stable policy?
2. Which module represents volatile detail?
3. Which direction do dependencies point today?
4. Does the candidate make details depend inward on policy?
5. Are abstractions owned by the stable side or the volatile side?
6. What cycle or hidden dependency remains?
7. What compile-time or runtime boundary should enforce the direction?

## Failure Signals

- Stable rules import volatile adapters directly.
- Both sides must know each other's internals.
- A new interface is owned by the wrong side.
- A dependency cycle is hidden behind callbacks or globals.

## Output

Record findings in the `## Interview Findings` table of `architecture-candidates.md` or `technical-design.md`:

| Lens | Finding | Design Impact |
| ---- | ------- | ------------- |
| Dependency Direction | stable policy module / volatile detail module / desired dependency direction | remaining cycle risk / boundary enforcement |
