# Gemini Specialist Roster Examples

Use these examples to pick a concrete Agent Team Gemini harness shape. They are orchestrator-routed specialist rosters using Gemini-native skills, agent definitions with explicit tools, artifacts, and optional invocation.

## Example 1: Parallel Research

Pattern: `fan_out_fan_in`

Specialists:

- `primary-source-researcher`
- `ecosystem-researcher`
- `synthesis-writer`

Flow:

1. Orchestrator records the research question and artifact root.
2. If invocation is appropriate and the orchestrator has `invoke_agent`, run the researchers independently; otherwise perform the two lanes directly using the same output contracts.
3. Researchers write artifacts under `_workspace/{plan}/research/`.
4. Orchestrator verifies evidence and passes artifacts to `synthesis-writer`.
5. Synthesis writer produces the final report.

## Example 2: Feature Build With Review

Pattern: `producer_reviewer`

Specialists:

- `feature-builder`
- `qa-reviewer`

Flow:

1. Builder produces code changes and a result artifact.
2. Reviewer verifies tests, behavior, integration contracts, and tool-use assumptions.
3. Reviewer returns `PASS`, `FIX`, or `BLOCKED`.
4. Orchestrator either accepts, retries builder with specific instructions, or blocks with a reason.

## Example 3: Documentation Update

Pattern: `pipeline`

Specialists:

- `source-auditor`
- `docs-writer`
- `docs-reviewer`

Flow:

1. Auditor gathers source evidence.
2. Writer updates docs with cited commands/files.
3. Reviewer checks commands, links, and accuracy.
4. Orchestrator reports changed files, verification, and reload reminder if harness files changed.

## Example 4: Large Migration

Pattern: `supervisor`

Specialists:

- `migration-supervisor`
- `migrator-a`
- `migrator-b`
- `compat-reviewer`

Flow:

1. Supervisor creates scoped, non-overlapping tasks.
2. Migrators complete assigned slices with evidence and artifacts.
3. Compatibility reviewer checks integration across slices.
4. Supervisor advances only after required evidence is present.

## Example 5: Ownership Transfer

Pattern: `message_coordination`

Specialists:

- `diagnostician`
- `domain-owner`
- `reviewer`

Flow:

1. Diagnostician identifies root cause and owner.
2. Orchestrator requests a message-based coordination with reason and artifact context.
3. Domain owner completes the task.
4. Reviewer validates final evidence.

## Example 6: Direct Orchestrator Only

Pattern: `expert_pool` or narrow `pipeline`

Specialists:

- no invoked worker required
- orchestrator loads the relevant skill/reference directly

Flow:

1. Orchestrator classifies the request as narrow.
2. Orchestrator performs the work directly.
3. If durable execution exists, record one task result with evidence and artifact.
4. Report concise outcome.

