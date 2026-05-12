# Evolutionary Refactoring Interview

Use when migration could be broad, risky, hard to stage, or likely to disrupt existing behavior.

Core lens: architecture should evolve through small verified steps when possible.

## Interview

1. Can this design be introduced in slices?
2. What is the first reversible step?
3. Which behavior must remain stable during migration?
4. Is there a strangler, compatibility wrapper, or branch-by-abstraction path?
5. What tests prove old and new paths agree?
6. What is the rollback or stop point if the design proves wrong?
7. Which part should not be refactored until later?

## Failure Signals

- The design requires a large rewrite before any value appears.
- Migration and behavior changes are mixed together.
- No intermediate state can be tested.
- Rollback would require manual reconstruction.

## Output

Record findings in the `## Interview Findings` table of `architecture-candidates.md` or `technical-design.md`:

| Lens | Finding | Design Impact |
| ---- | ------- | ------------- |
| Evolutionary Refactoring | migration slices / first step / compatibility path | verification at each step / rollback or stop point |
