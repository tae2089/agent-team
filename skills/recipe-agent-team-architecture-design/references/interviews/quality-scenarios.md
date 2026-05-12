# Quality Scenarios Interview

Use when quality attributes such as reliability, performance, security, operability, usability, compatibility, or maintainability are vague.

Core lens: quality requirements should be tested through concrete scenarios.

## Interview

1. Which quality attribute matters for this design?
2. What event or change triggers the scenario?
3. Under what conditions does it occur?
4. What system response is required?
5. How will the response be measured or observed?
6. What test, check, or review proves the scenario?
7. Which candidate best satisfies the scenario with the least extra complexity?

## Failure Signals

- The design claims to be robust, scalable, or safe without a scenario.
- Acceptance criteria do not measure the quality claim.
- The candidate optimizes for a quality attribute the user did not need.
- The scenario cannot be verified by tests, logs, artifacts, or review.

## Output

Record findings in the `## Interview Findings` table of `architecture-candidates.md` or `technical-design.md`:

| Lens | Finding | Design Impact |
| ---- | ------- | ------------- |
| Quality Scenarios | quality attribute / concrete scenario / measurable response | verification method / candidate impact |
