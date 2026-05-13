#!/bin/sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
TMP_DIR=$(mktemp -d)

: "${GOCACHE:=$TMP_DIR/gocache}"
export GOCACHE

cd "$ROOT_DIR"

bash scripts/check-skill-crossrefs.sh

go test ./internal/cli -run 'Test(CommandSpecMatchesCobraTree|HelperSkillDocsMatchSpec|ServiceSkillDocsExist|SkillDocsAvoidInvalidBooleanShellExamples|ErrorRecoveryDocsMatchSpec|HarnessSandboxPermissionDocsAreWired|HarnessRuntimeContextDocsAreWired|AgentTeamHarnessSkillsRequireRuntimeStateProtocol|HarnessSkillDogfoodExamplesRecordClosedState)'

go run ./cmd/agent-team schema export > "$TMP_DIR/schema.json"
jq empty "$TMP_DIR/schema.json"

jq empty \
  examples/dogfood-harness/results/codex/run-summary.json \
  examples/dogfood-harness/results/codex/event-log.json \
  examples/dogfood-harness/results/gemini/run-summary.json \
  examples/dogfood-harness/results/gemini/event-log.json \
  examples/dogfood-harness/results/codex-harness-skill/run-summary.json \
  examples/dogfood-harness/results/codex-harness-skill/event-log.json \
  examples/dogfood-harness/results/gemini-harness-skill/run-summary.json \
  examples/dogfood-harness/results/gemini-harness-skill/event-log.json

if grep -R -n -E 'Generated harnesses default orchestrated runtime execution to direct `agent-team`|Runtime execution uses `agent-team` directly|route durable coordination through direct `agent-team` runtime execution' skills; then
	echo "stale direct-CLI harness wording found" >&2
	exit 1
fi

if grep -R -n -E 'If an `RUN_ID` or `TASK_ID` is provided, resume that run/task|They resume a provided `RUN_ID`/`TASK_ID`' skills; then
	echo "stale user-facing runtime ID wording found" >&2
	exit 1
fi

printf 'drift checks passed\n'
