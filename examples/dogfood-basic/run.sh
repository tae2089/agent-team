#!/bin/sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
STATE_DIR=$(mktemp -d)
WORK_DIR=$(mktemp -d)

export AGENT_TEAM_STATE_DIR="$STATE_DIR"
: "${GOCACHE:=$WORK_DIR/gocache}"
export GOCACHE

cd "$ROOT_DIR"

go run ./cmd/agent-team init
go run ./cmd/agent-team run create --id run_dogfood --title "dogfood basic"
go run ./cmd/agent-team task create \
  --id task_dogfood \
  --run run_dogfood \
  --agent worker \
  --title "write dogfood artifact" \
  --body "Write a small artifact proving the workflow ran."

mkdir -p "$WORK_DIR/_workspace/run_dogfood"
printf '%s\n' "dogfood artifact" > "$WORK_DIR/_workspace/run_dogfood/task_dogfood.md"

go run ./cmd/agent-team inbox list --agent worker --run run_dogfood --unread
go run ./cmd/agent-team sync check --agent worker --run run_dogfood --task task_dogfood
go run ./cmd/agent-team task start --task task_dogfood --agent worker
go run ./cmd/agent-team task complete \
  --task task_dogfood \
  --agent worker \
  --evidence "Artifact written in temporary workspace." \
  --artifact "$WORK_DIR/_workspace/run_dogfood/task_dogfood.md"
go run ./cmd/agent-team run summary --run run_dogfood
go run ./cmd/agent-team run close --run run_dogfood --reason "dogfood workflow complete"

printf 'state_dir=%s\n' "$STATE_DIR"
printf 'work_dir=%s\n' "$WORK_DIR"
