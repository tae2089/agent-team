# agent-team

`agent-team` is a daemonless CLI for coordinating agent workflows through durable local state.

It stores runs, tasks, messages, inbox acknowledgements, sync checks, and event history in SQLite. It does **not** run workers, spawn LLMs, watch files, or require a background daemon.

## Install

Install with Go:

```bash
go install github.com/tae2089/agent-team/cmd/agent-team@latest
```

Install with npm:

```bash
npm install -g @tae2089/agent-team
```

From this repository:

```bash
go install ./cmd/agent-team
```

Verify the binary:

```bash
agent-team version
agent-team --help
```

Release binaries and checksum verification are documented in [docs/install.md](docs/install.md).

## Quickstart

Copy and paste this from a repository checkout. If `agent-team` is already installed, replace `go run ./cmd/agent-team` with `agent-team`.

```bash
export AGENT_TEAM_STATE_DIR="$(mktemp -d)"
go run ./cmd/agent-team init
go run ./cmd/agent-team run create --id run_docs --title "docs workflow"
go run ./cmd/agent-team task create --id task_docs --run run_docs --agent writer --title "draft docs" --body "Write the first draft."
go run ./cmd/agent-team task start --task task_docs --agent writer
mkdir -p _workspace/run_docs
printf '%s\n' "Draft complete." > _workspace/run_docs/task_docs.md
go run ./cmd/agent-team task complete --task task_docs --agent writer --evidence "Draft written and checked." --artifact "_workspace/run_docs/task_docs.md"
go run ./cmd/agent-team run summary --run run_docs
go run ./cmd/agent-team run close --run run_docs --reason "all tasks terminal"
```

## State And Artifacts

Default state path:

```text
.agent-team/agent-team.db
```

Override it per workflow:

```bash
export AGENT_TEAM_STATE_DIR=/tmp/agent-team-state
```

Artifacts should live outside the DB:

```text
_workspace/{run_id}/
```

Harness users normally do not provide `RUN_ID` or `TASK_ID`. Generated Codex/Gemini harness orchestrators create or resolve those IDs as internal runtime context and pass them to workers. See [docs/harness-runtime-context.md](docs/harness-runtime-context.md).

### Pre-release DB Stability

`agent-team` is pre-v1. The SQLite schema under `.agent-team/agent-team.db` is not yet a stable compatibility contract.

- Until v1, breaking DB schema changes may require deleting the state directory and running `agent-team init` again.
- Do not store irreplaceable workflow artifacts only in the DB. Keep durable outputs under `_workspace/{run_id}/`.
- Additive compatibility helpers may exist, but this project intentionally does not maintain a versioned migration framework before release.

## Input Model

Every command rejects positional arguments. Use named flags, `--params`, and `--json`.

```bash
agent-team task complete \
  --params '{"task_id":"task_docs","agent":"writer","force":false}' \
  --json '{"evidence":"Draft written.","artifact":"_workspace/run_docs/task_docs.md"}'
```

- `--params`: IDs, filters, statuses, booleans, and routing values.
- `--json`: richer payloads such as task bodies, metadata, evidence, artifact paths, message bodies, and reasons.
- If a value appears in both a named flag and JSON, the command returns `input_conflict`.
- Boolean shell flags are presence-based: use `--force` or `--unread`.
- Failure responses include `error.recovery` hints. See [docs/errors.md](docs/errors.md).

## Operational Commands

```bash
agent-team run list --status open
agent-team run summary --run run_docs
agent-team task stale --run run_docs --older-than 2h
agent-team task retry --task task_docs --reason "dependency fixed"
agent-team task reassign --task task_docs --agent backup-writer --reason "handoff"
agent-team message list --run run_docs --unread
agent-team event log --run run_docs --limit 100
agent-team schema export
```

`schema export` is the machine-readable command contract used by tests and skill documentation consistency checks. Run/task status rules are documented in [docs/status-policy.md](docs/status-policy.md).

## Harness Integration

The repository includes gwscli-style skills under `skills/`.

- Service skills: `agent-team-run`, `agent-team-task`, `agent-team-inbox`, `agent-team-sync`, `agent-team-ops`
- Command helper skills: one per CLI command
- Recipes: terminology context, planning grill, architecture design, compound learning, run lifecycle, worker checkpoint, and operational audit
- Persona skills: `persona-agent-team-planner` for terminology/planning/backend architecture routing; `persona-agent-team-designer` for design interview/spec routing
- Harness skills: `agent-team-codex-harness` and `agent-team-gemini-harness`
- Backlog: add a cross-reference linter for skill names, handoffs, and artifact paths; current checks cover command contracts, markdown links, and formatting.

Generated harnesses should load `agent-team-shared`, choose a recipe/service skill, then load exact command helper skills for syntax and flags.

The CI dogfood example exercises the runtime contract without requiring an external LLM process:

```bash
sh examples/dogfood-harness/run.sh
```

Manual Codex and Gemini harness dogfood guides are in [examples/dogfood-harness/README.md](examples/dogfood-harness/README.md).
Sandbox and permission guidance is in [docs/harness-sandbox-permissions.md](docs/harness-sandbox-permissions.md).

## Development

```bash
GOCACHE=/tmp/gocache GOMODCACHE=/tmp/gomodcache go test ./...
sh scripts/ci-drift-check.sh
sh examples/dogfood-basic/run.sh
sh examples/dogfood-harness/run.sh
sh examples/stress/ci-smoke.sh
```

Run a heavier local stress profile. The heavy profile also exercises retry, reassign, cancel, and fail transitions:

```bash
sh examples/stress/local-heavy.sh
AGENT_TEAM_STRESS_N=100 sh examples/stress/concurrency.sh
```

Release tags use the `v*` shape. The release workflow builds OS/architecture binaries, injects the tag version into `agent-team version`, and uploads `SHA256SUMS`. See [docs/release.md](docs/release.md) and [CHANGELOG.md](CHANGELOG.md).
