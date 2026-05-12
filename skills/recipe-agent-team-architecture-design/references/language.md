# Architecture Design Language

Use this reference when evaluating architecture candidates or writing technical design artifacts. The purpose is consistent design language, not vocabulary for its own sake.

## Terms

**Module**

A design unit with a caller-facing contract and internal code. It can be a function, class, package, command, workflow slice, or larger subsystem.

Use this term when discussing what owns behavior. Avoid switching between unit, component, service, or layer unless the repository already uses those names for concrete things.

**Interface**

Everything a caller must understand to use a module correctly. This includes signatures, inputs, outputs, ordering rules, invariants, errors, configuration, and performance expectations.

Use this term for the full caller contract, not only language-level interfaces or public methods.

**Implementation**

The code behind a module's interface. Use this term for internal behavior that callers should not need to know.

**Depth**

How much useful behavior sits behind a module's interface. A deeper module gives callers more capability while asking them to learn less.

Do not measure depth by implementation line count. Measure it by how much caller complexity disappears.

**Seam**

A place where behavior can vary without editing the caller. The seam is where an interface is placed.

Use this term when discussing where variation or substitution should happen.

**Adapter**

A concrete implementation that fills a seam. Use this term for the role a concrete piece plays at a seam, not for its size or internal complexity.

**Leverage**

The caller benefit of depth: one interface or implementation decision pays off across many call sites, tests, or workflow steps.

**Locality**

The maintainer benefit of depth: knowledge, change, bugs, and verification concentrate in one owner instead of spreading through callers.

## Principles

- Depth belongs to the module interface, not to the amount of internal code.
- A module can have private internal seams for its own implementation and tests while exposing one smaller external interface.
- If deleting a module removes almost no caller complexity, the module may be too shallow.
- If deleting a module pushes complexity into many callers, the module was earning its place.
- Callers and tests should usually cross the same interface. If tests need to reach behind the interface, the module shape may be wrong.
- One concrete adapter often means the seam is only hypothetical. Multiple real adapters are stronger evidence that the seam earns its cost.

## Candidate Review Prompts

Use these prompts in architecture candidate reviews:

- What module owns the behavior after this design?
- What does a caller need to know at the interface?
- What complexity moves behind the interface?
- Where is the seam, and what actually varies across it?
- Which adapters exist now or are expected soon?
- What caller complexity disappears if this module is deep enough?
- What maintenance knowledge becomes local?
- What tests become simpler because they can use the same interface as callers?

## Rejected Framings

- Do not treat depth as a ratio of implementation lines to interface lines.
- Do not use interface to mean only a language keyword or method list.
- Avoid boundary when the design question is really about a seam or caller interface.
