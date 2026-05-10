---
name: agent-team-schema-export
description: "agent-team: Export the machine-readable CLI command and error contract."
metadata:
  version: 1.2.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
  cliHelp: "agent-team schema export --help"
---

# schema export

> **PREREQUISITE:** Read `../agent-team-shared/SKILL.md`.

Export the command contract for tooling, audits, and doc synchronization.

## Usage

```bash
agent-team schema export
```

## Flags

| Flag | JSON key | Required | Default | Description |
|------|----------|----------|---------|-------------|
| none | none | no | none | This command accepts no command-specific inputs. |

## Examples

```bash
agent-team schema export
```

## Output

- `data.schema.command`: CLI name.
- `data.schema.version`: contract version.
- `data.schema.commands`: command list with flags, params, required params, and output keys.
- `data.schema.errors`: stable error codes.
- `data.schema.error_recovery`: recovery summaries, actions, docs, and related skills for each error code.

## Errors

| Code | Meaning | Action |
|------|---------|--------|
| `validation_error` | Positional arguments were provided. | Run without positional arguments. |

## See Also

- [agent-team-shared](../agent-team-shared/SKILL.md)
- [agent-team-ops](../agent-team-ops/SKILL.md)
