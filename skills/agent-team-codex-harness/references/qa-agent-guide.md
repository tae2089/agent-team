# Codex QA Agent Guide

Use this reference when adding reviewer, QA, validation, or security specialists to a Agent Team Codex harness.

## 1. QA Role

QA is a quality gate, not a final checklist. It verifies contracts, integration, evidence, artifact correctness, and runtime-state discipline before the orchestrator advances.

Prefer incremental QA after meaningful modules or artifacts. A single final pass is usually too late and too shallow.

## 2. What QA Checks

| Area | Examples |
| --- | --- |
| boundary contracts | API response vs consumer assumptions, CLI output vs parser |
| route/path integrity | links, imports, page routes, generated paths |
| state transitions | allowed statuses, blocked behavior, retry limits |
| data flow | producer artifact is consumed by the next phase |
| evidence | commands, inspected files, source references |
| artifact quality | expected path, required sections, no missing outputs |
| runtime discipline | done tasks have evidence + artifact, blocked tasks have reason |

## 3. Integration Coherence Verification

The most valuable QA work is cross-boundary comparison.

Examples:

- API response shape vs frontend hook/type expectation
- CLI command output vs parser assumptions
- documented command vs actual command behavior
- task status transitions vs code that updates statuses
- generated file path vs reader/import/link path
- database schema vs app model/type mapping

QA should read both sides of a boundary and compare the contract directly.

## 4. Verdicts

- `PASS`: evidence is sufficient; orchestrator may advance.
- `FIX`: producer should retry with specific instructions.
- `BLOCKED`: missing decision, permission, owner, or unverifiable dependency.

QA must not silently pass missing evidence.

## 5. QA Artifact

```markdown
# QA Review: {task}

## Verdict
PASS | FIX | BLOCKED

## Evidence
- command/result
- inspected files or artifacts
- source references

## Findings
- severity, location, issue, recommended fix

## Contract Checks
- producer output -> consumer expectation
- docs/commands -> actual behavior
- runtime state -> required task contract

## Recommendation
- pass
- retry producer with specific instructions
- message owner with context
- block with reason
```

## 6. Runtime Reporting

When QA has an assigned Agent Team task:

- use `agent-team task complete` only after writing the QA artifact and evidence
- use `blocked` with a concrete blocked reason when QA cannot verify
- do not update unrelated tasks
- do not advance sync checks; the orchestrator owns advancement
