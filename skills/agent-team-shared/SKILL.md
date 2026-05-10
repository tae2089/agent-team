---
name: agent-team-shared
description: "Shared runtime procedure for the agent-team CLI. Use for daemonless state, global flags, JSON input/output, state directory rules, error handling, and shell usage conventions. Load this before command-specific agent-team skills."
metadata:
  version: 1.2.0
  openclaw:
    category: "agent-orchestration"
  requires:
    bins:
      - agent-team
  cliHelp: "agent-team --help"
---

# Agent Team Shared Runtime

Use `agent-team` for daemonless coordination between orchestrators and workers. The CLI stores state only; it does not run agents, watch files, or require a background daemon.

## Installation

`agent-team` must be available on `$PATH`.

```bash
agent-team --help
```

Initialize state once per project or custom state directory:

```bash
agent-team init
```

## State

- Default state: `.agent-team/agent-team.db`
- Override state directory: `AGENT_TEAM_STATE_DIR=/path/to/state`
- Artifact convention: `_workspace/{run_id}/`
- SQLite is the source of truth for runs, tasks, messages, inbox, sync status, and event history.
- Do not use direct peer-to-peer agent communication. Send compact messages with `agent-team message send`; recipients check inboxes at checkpoints.

## CLI Syntax

```bash
agent-team <resource> <command> [flags]
```

Examples:

```bash
agent-team run create --title "docs refactor"
agent-team task create --run run_docs --agent writer --title "draft docs"
agent-team inbox list --agent writer --run run_docs --unread
agent-team task complete --task task_docs --agent writer --evidence "Verified locally." --artifact "_workspace/run_docs/task_docs.md"
```

## Global Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--params` | `{}` | Strict JSON command arguments. Accepts inline JSON, `@file.json`, or `-` for stdin. |
| `--json` | `{}` | Strict JSON payload. Use for richer content such as task bodies, metadata, evidence, artifact paths, message body, and reasons. |

Unknown JSON fields are rejected. If a value is supplied by both a named flag and JSON, the CLI returns `input_conflict`.

## Shell Tips

- Boolean flags are presence-based in shell usage: use `--unread` or `--force`.
- Use JSON for explicit false values: `--params '{"force":false}'`.
- Quote JSON with single quotes in POSIX shells.
- Use `@file.json` for large JSON objects.
- Use `-` to read a JSON object from stdin.
- Do not pass positional arguments; use flags, `--params`, or `--json`.

## Common Aliases

| Flag | JSON key | Meaning |
|------|----------|---------|
| `--run` | `run_id` | Workflow scope. |
| `--task` | `task_id` | Task scope. |
| `--msg` | `msg_id` | Inbox message scope. |
| `--agent` | `agent` | Logical agent name. |
| `--id` | `id` | Optional caller-provided stable identifier. |

## Common Value Semantics

| Value | Meaning |
|-------|---------|
| `metadata` | JSON object for compact machine-readable context. Put large content in artifacts. |
| `body` | Human-readable task or message content. |
| `reason` | Operational explanation persisted in state or event payloads. |
| `evidence` | Completion proof in text. State what was checked, changed, or verified. |
| `artifact` | Path to produced work, normally under `_workspace/{run_id}/`. |

## JSON Output

Success:

```json
{"ok":true,"state_version":12,"data":{},"warnings":[]}
```

Failure:

```json
{"ok":false,"state_version":12,"error":{"code":"sync_conflict","message":"...","details":{},"recovery":{"summary":"...","actions":[],"commands":[],"docs":[],"skills":[]}}}
```

Agents should branch on `ok` and `error.code`. Use `error.recovery` for operator-facing next steps and related helper skills. Human-readable messages are not stable interfaces.

## Error Codes

| Code | Meaning | Typical action |
|------|---------|----------------|
| `validation_error` | Required or invalid input. | Fix flags or JSON. |
| `input_conflict` | Same field supplied by flag and JSON. | Provide the value once. |
| `invalid_json` | JSON payload is malformed or has unknown fields. | Fix JSON shape. |
| `invalid_json_source` | `@file` or stdin JSON source is invalid. | Fix the source argument. |
| `not_found` | Run, task, or message does not exist. | Check IDs. |
| `agent_mismatch` | Agent does not own the task/message. | Reassign or use the correct agent. |
| `sync_conflict` | Unread messages or incomplete dependencies. | Run sync check, resolve, retry. |
| `run_not_ready` | Run has unfinished tasks. | Finish, retry, or reassign tasks. |
| `invalid_run_state` | Run status does not allow operation. | Inspect run status. |
| `invalid_task_state` | Task status does not allow operation. | Inspect task and choose retry/reassign/block. |
| `internal_error` | Local filesystem, SQLite, state, install, or sandbox failure. | Check state dir permissions, run init, or inspect install/sandbox docs. |

Full recovery guidance lives in `docs/errors.md`. `agent-team schema export` exposes the same `error_recovery` catalog for agents and CI drift checks.

## Skill Routing

| Need | Load |
|------|------|
| Run lifecycle | `agent-team-run` |
| Task lifecycle | `agent-team-task` |
| Messages and inbox | `agent-team-inbox` |
| Drift checks before completion | `agent-team-sync` |
| Audit, logs, schema, operational cleanup | `agent-team-ops` |

Use `agent-team schema export` to inspect the current machine-readable command contract.
