#!/bin/sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
STATE_DIR=$(mktemp -d)
WORK_DIR=$(mktemp -d)
BIN="$WORK_DIR/agent-team"
PROFILE=${AGENT_TEAM_STRESS_PROFILE:-smoke}

case "$PROFILE" in
	smoke)
		DEFAULT_N=12
		DEFAULT_RUNS=1
		;;
	heavy)
		DEFAULT_N=50
		DEFAULT_RUNS=3
		;;
	*)
		printf 'unknown AGENT_TEAM_STRESS_PROFILE: %s\n' "$PROFILE" >&2
		exit 1
		;;
esac

N=${AGENT_TEAM_STRESS_N:-$DEFAULT_N}
RUNS=${AGENT_TEAM_STRESS_RUNS:-$DEFAULT_RUNS}

export AGENT_TEAM_STATE_DIR="$STATE_DIR"
: "${GOCACHE:=$WORK_DIR/gocache}"
export GOCACHE

cd "$ROOT_DIR"
go build -o "$BIN" ./cmd/agent-team

agent_team() {
	"$BIN" "$@"
}

wait_all() {
	for pid in "$@"; do
		wait "$pid"
	done
}

count_status() {
	run_id=$1
	status=$2
	agent_team task list --run "$run_id" --status "$status" --limit 1000 | grep -o "\"status\":\"$status\"" | wc -l | tr -d ' '
}

assert_event() {
	run_id=$1
	event_type=$2
	events=$(agent_team event log --run "$run_id" --type "$event_type" --limit 1000)
	if ! printf '%s\n' "$events" | grep -q "\"event_type\":\"$event_type\""; then
		printf 'expected %s event for %s\n%s\n' "$event_type" "$run_id" "$events" >&2
		exit 1
	fi
}

run_heavy_mixed_ops() {
	run_index=$1
	run_id=$2

	block_retry="task_stress_${run_index}_retry_blocked"
	fail_retry="task_stress_${run_index}_retry_failed"
	reassign_task="task_stress_${run_index}_reassign"
	cancel_task="task_stress_${run_index}_cancel"
	fail_terminal="task_stress_${run_index}_fail_terminal"

	agent_team task create --id "$block_retry" --run "$run_id" --agent ops_a --title "retry blocked task" >/dev/null
	agent_team task start --task "$block_retry" --agent ops_a >/dev/null
	agent_team task block --task "$block_retry" --agent ops_a --reason "dependency not ready" >/dev/null
	agent_team task retry --task "$block_retry" --reason "dependency fixed" >/dev/null
	agent_team task start --task "$block_retry" --agent ops_a >/dev/null
	block_artifact="$WORK_DIR/_workspace/$run_id/$block_retry.md"
	printf '%s\n' "blocked retry artifact" > "$block_artifact"
	agent_team task complete --task "$block_retry" --agent ops_a --evidence "blocked task retried and completed" --artifact "$block_artifact" >/dev/null

	agent_team task create --id "$fail_retry" --run "$run_id" --agent ops_b --title "retry failed task" >/dev/null
	agent_team task start --task "$fail_retry" --agent ops_b >/dev/null
	agent_team task fail --task "$fail_retry" --agent ops_b --reason "first attempt failed" --artifact "$WORK_DIR/_workspace/$run_id/$fail_retry.failed.md" >/dev/null
	agent_team task retry --task "$fail_retry" --reason "retry after failed attempt" >/dev/null
	agent_team task start --task "$fail_retry" --agent ops_b >/dev/null
	fail_retry_artifact="$WORK_DIR/_workspace/$run_id/$fail_retry.md"
	printf '%s\n' "failed retry artifact" > "$fail_retry_artifact"
	agent_team task complete --task "$fail_retry" --agent ops_b --evidence "failed task retried and completed" --artifact "$fail_retry_artifact" >/dev/null

	agent_team task create --id "$reassign_task" --run "$run_id" --agent ops_c --title "reassign task" >/dev/null
	agent_team task reassign --task "$reassign_task" --agent ops_d --reason "handoff to available worker" >/dev/null
	agent_team task start --task "$reassign_task" --agent ops_d >/dev/null
	reassign_artifact="$WORK_DIR/_workspace/$run_id/$reassign_task.md"
	printf '%s\n' "reassign artifact" > "$reassign_artifact"
	agent_team task complete --task "$reassign_task" --agent ops_d --evidence "reassigned task completed" --artifact "$reassign_artifact" >/dev/null

	agent_team task create --id "$cancel_task" --run "$run_id" --agent ops_e --title "cancel task" >/dev/null
	agent_team task cancel --task "$cancel_task" --reason "scope removed" >/dev/null

	agent_team task create --id "$fail_terminal" --run "$run_id" --agent ops_f --title "terminal fail task" >/dev/null
	agent_team task start --task "$fail_terminal" --agent ops_f >/dev/null
	agent_team task fail --task "$fail_terminal" --agent ops_f --reason "kept failed for final outcome" --artifact "$WORK_DIR/_workspace/$run_id/$fail_terminal.md" >/dev/null
}

run_once() {
	run_index=$1
	run_id="run_stress_$run_index"

	agent_team run create --id "$run_id" --title "concurrency stress $run_index" >/dev/null

	pids=""
	i=1
	while [ "$i" -le "$N" ]; do
		agent_team task create --id "task_stress_${run_index}_$i" --run "$run_id" --agent "worker_$i" --title "stress task $i" >/dev/null &
		pids="$pids $!"
		i=$((i + 1))
	done
	wait_all $pids

	pids=""
	i=1
	while [ "$i" -le "$N" ]; do
		agent_team message send --id "msg_stress_${run_index}_$i" --run "$run_id" --task "task_stress_${run_index}_$i" --from orchestrator --to "worker_$i" --kind notice --body "stress message $i" >/dev/null &
		pids="$pids $!"
		i=$((i + 1))
	done
	wait_all $pids

	pids=""
	i=1
	while [ "$i" -le "$N" ]; do
		(
			agent_team inbox ack --msg "msg_stress_${run_index}_$i" --agent "worker_$i" >/dev/null
			agent_team task start --task "task_stress_${run_index}_$i" --agent "worker_$i" >/dev/null
			artifact="$WORK_DIR/_workspace/$run_id/task_stress_$i.md"
			mkdir -p "$(dirname "$artifact")"
			printf '%s\n' "stress artifact $run_index/$i" > "$artifact"
			agent_team task complete --task "task_stress_${run_index}_$i" --agent "worker_$i" --evidence "stress task $i completed" --artifact "$artifact" >/dev/null
		) &
		pids="$pids $!"
		i=$((i + 1))
	done
	wait_all $pids

	if [ "$PROFILE" = "heavy" ]; then
		run_heavy_mixed_ops "$run_index" "$run_id"
	fi

	expected_done=$N
	if [ "$PROFILE" = "heavy" ]; then
		expected_done=$((N + 3))
	fi
	done_count=$(count_status "$run_id" done)
	if [ "$done_count" -ne "$expected_done" ]; then
		printf 'expected %s done tasks for %s, got %s\n' "$expected_done" "$run_id" "$done_count" >&2
		agent_team run summary --run "$run_id" >&2
		exit 1
	fi
	if [ "$PROFILE" = "heavy" ]; then
		cancelled_count=$(count_status "$run_id" cancelled)
		failed_count=$(count_status "$run_id" failed)
		if [ "$cancelled_count" -lt 1 ] || [ "$failed_count" -lt 1 ]; then
			printf 'expected at least one cancelled and failed task for %s, got cancelled=%s failed=%s\n' "$run_id" "$cancelled_count" "$failed_count" >&2
			agent_team run summary --run "$run_id" >&2
			exit 1
		fi
		assert_event "$run_id" task_retried
		assert_event "$run_id" task_reassigned
		assert_event "$run_id" task_cancelled
		assert_event "$run_id" task_failed
	fi

	summary=$(agent_team run summary --run "$run_id")
	if ! printf '%s\n' "$summary" | grep -q '"close_ready":true'; then
		printf 'expected %s to be close ready\n%s\n' "$run_id" "$summary" >&2
		exit 1
	fi

	agent_team event log --run "$run_id" --limit 1000 >/dev/null
	agent_team run close --run "$run_id" --reason "stress workflow complete" >/dev/null
	assert_event "$run_id" run_closed
}

agent_team init >/dev/null

run_index=1
while [ "$run_index" -le "$RUNS" ]; do
	run_once "$run_index"
	run_index=$((run_index + 1))
done

printf 'profile=%s\n' "$PROFILE"
printf 'runs=%s\n' "$RUNS"
printf 'tasks_per_run=%s\n' "$N"
printf 'state_dir=%s\n' "$STATE_DIR"
printf 'work_dir=%s\n' "$WORK_DIR"
