# Deep Modules Interview

Use when a candidate introduces a module, service, helper, package, or layer and you need to test whether it earns its abstraction cost.

Core lens: a deep module gives callers substantial capability through a small interface.

## Interview

1. What caller complexity disappears behind this interface?
2. Is the interface smaller than the behavior it unlocks?
3. Is this module doing real work or mostly forwarding calls?
4. What decisions become local to the module?
5. What must a caller still know after the change?
6. If this module were deleted, how much complexity would move into callers?
7. Could the same depth be achieved by improving an existing module?

## Failure Signals

- The candidate adds a pass-through layer.
- The caller must still understand the same internals.
- The module is named after a vague role rather than owned behavior.
- The abstraction exists mainly to make files look organized.

## Output

Record findings in the `## Interview Findings` table of `architecture-candidates.md` or `technical-design.md`:

| Lens | Finding | Design Impact |
| ---- | ------- | ------------- |
| Deep Modules | caller complexity removed / interface size and shape / module depth judgment / existing module that could absorb the behavior | design adjustment |
