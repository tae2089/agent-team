#!/bin/sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
STATE_DIR=$(mktemp -d)
WORK_DIR=$(mktemp -d)
BIN="$WORK_DIR/agent-team"
export AGENT_TEAM_STATE_DIR="$STATE_DIR"
: "${GOCACHE:=$WORK_DIR/gocache}"
export GOCACHE

cd "$ROOT_DIR"
go build -o "$BIN" ./cmd/agent-team

agent_team() {
	"$BIN" "$@"
}

assert_contains() {
	value=$1
	needle=$2
	if ! printf '%s\n' "$value" | grep -q "$needle"; then
		printf 'expected output to contain %s\noutput: %s\n' "$needle" "$value" >&2
		exit 1
	fi
}

agent_team init
RUN_JSON=$(agent_team run create --title "harness dogfood")
RUN_ID=$(printf '%s\n' "$RUN_JSON" | jq -r '.data.run.id')

PLAN_ARTIFACT="$WORK_DIR/_workspace/$RUN_ID/plan.md"
WORK_ARTIFACT="$WORK_DIR/_workspace/$RUN_ID/worker.md"
REVIEW_ARTIFACT="$WORK_DIR/_workspace/$RUN_ID/review.md"

mkdir -p "$WORK_DIR/_workspace/$RUN_ID"

PLAN_TASK_JSON=$(agent_team task create \
	--run "$RUN_ID" \
	--agent planner \
	--title "write implementation plan" \
	--body "Create a short implementation plan artifact.")
PLAN_TASK=$(printf '%s\n' "$PLAN_TASK_JSON" | jq -r '.data.task.id')

WORK_TASK_JSON=$(agent_team task create \
	--run "$RUN_ID" \
	--agent worker \
	--title "implement planned change" \
	--depends-on "$PLAN_TASK" \
	--body "Read planner guidance, check sync, and produce a worker artifact.")
WORK_TASK=$(printf '%s\n' "$WORK_TASK_JSON" | jq -r '.data.task.id')

REVIEW_TASK_JSON=$(agent_team task create \
	--run "$RUN_ID" \
	--agent reviewer \
	--title "review worker artifact" \
	--depends-on "$WORK_TASK" \
	--body "Review the worker artifact and close the loop.")
REVIEW_TASK=$(printf '%s\n' "$REVIEW_TASK_JSON" | jq -r '.data.task.id')

printf '%s\n' "Plan: use agent-team runtime commands for every role transition." > "$PLAN_ARTIFACT"
agent_team task start --task "$PLAN_TASK" --agent planner
agent_team task complete --task "$PLAN_TASK" --agent planner --evidence "Plan artifact written." --artifact "$PLAN_ARTIFACT"

agent_team message send \
	--id msg_harness_contract \
	--run "$RUN_ID" \
	--task "$WORK_TASK" \
	--from orchestrator \
	--to worker \
	--kind contract_changed \
	--body "Use the plan artifact before completing the worker task."

sync_output=$(agent_team sync check --agent worker --run "$RUN_ID" --task "$WORK_TASK")
assert_contains "$sync_output" '"blocking":true'

agent_team inbox ack --msg msg_harness_contract --agent worker
agent_team task start --task "$WORK_TASK" --agent worker
printf '%s\n' "Worker: implemented according to $PLAN_ARTIFACT." > "$WORK_ARTIFACT"
agent_team task complete --task "$WORK_TASK" --agent worker --evidence "Worker artifact written after inbox acknowledgement." --artifact "$WORK_ARTIFACT"

agent_team task start --task "$REVIEW_TASK" --agent reviewer
printf '%s\n' "Review: worker artifact is present and the runtime contract held." > "$REVIEW_ARTIFACT"
agent_team task complete --task "$REVIEW_TASK" --agent reviewer --evidence "Review artifact written." --artifact "$REVIEW_ARTIFACT"

summary=$(agent_team run summary --run "$RUN_ID")
assert_contains "$summary" '"close_ready":true'

events=$(agent_team event log --run "$RUN_ID" --limit 50)
assert_contains "$events" '"event_type":"message_acked"'
assert_contains "$events" '"event_type":"task_completed"'

agent_team schema export >/dev/null
agent_team run close --run "$RUN_ID" --reason "harness dogfood workflow complete"

printf 'state_dir=%s\n' "$STATE_DIR"
printf 'work_dir=%s\n' "$WORK_DIR"
printf 'run_id=%s\n' "$RUN_ID"
