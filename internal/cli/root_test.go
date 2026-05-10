package cli

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func runCLI(t *testing.T, args ...string) map[string]any {
	t.Helper()
	cmd := NewRoot()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v\noutput: %s", err, out.String())
	}
	var got map[string]any
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatalf("invalid json output: %v\n%s", err, out.String())
	}
	if ok, _ := got["ok"].(bool); !ok {
		t.Fatalf("command returned ok=false: %s", out.String())
	}
	return got
}

func runCLIFail(t *testing.T, args ...string) error {
	t.Helper()
	cmd := NewRoot()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(args)
	return cmd.Execute()
}

func TestCoreWorkflow(t *testing.T) {
	t.Setenv("AGENT_TEAM_STATE_DIR", t.TempDir())

	runCLI(t, "init")
	runCLI(t, "run", "create", "--params", `{"id":"run_test","title":"test run"}`)
	runCLI(t, "task", "create",
		"--params", `{"id":"task_test","run_id":"run_test","agent":"worker","title":"do work"}`,
		"--json", `{"body":"body","metadata":{"scope":"unit"}}`,
	)
	runCLI(t, "task", "start", "--params", `{"task_id":"task_test","agent":"worker"}`)
	runCLI(t, "task", "complete",
		"--params", `{"task_id":"task_test","agent":"worker","force":false}`,
		"--json", `{"evidence":"verified","artifact":"_workspace/run_test/worker.md"}`,
	)
	status := runCLI(t, "run", "status", "--params", `{"run_id":"run_test"}`)
	data := status["data"].(map[string]any)
	tasks := data["tasks"].(map[string]any)
	if tasks["done"].(float64) != 1 {
		t.Fatalf("expected one done task, got %#v", tasks)
	}
}

func TestHybridInputFlags(t *testing.T) {
	t.Setenv("AGENT_TEAM_STATE_DIR", t.TempDir())

	runCLI(t, "init")
	runCLI(t, "run", "create", "--id", "run_hybrid", "--title", "Hybrid Run")
	runCLI(t, "task", "create", "--id", "task_hybrid", "--run", "run_hybrid", "--agent", "worker", "--title", "Hybrid Task", "--body", "flag body", "--metadata", `{"source":"cli","quality":"high"}`)
	runCLI(t, "task", "start", "--task", "task_hybrid", "--agent", "worker")
	runCLI(t, "task", "complete", "--task", "task_hybrid", "--agent", "worker", "--evidence", "passed", "--artifact", "_workspace/run_hybrid/worker.md")
	show := runCLI(t, "task", "show", "--task", "task_hybrid")
	taskPayload := show["data"].(map[string]any)["task"].(map[string]any)
	if taskPayload["id"] != "task_hybrid" {
		t.Fatalf("expected task_hybrid, got %#v", taskPayload["id"])
	}
}

func TestFlagJSONConflict(t *testing.T) {
	t.Setenv("AGENT_TEAM_STATE_DIR", t.TempDir())

	runCLI(t, "init")
	err := runCLIFail(t, "run", "create", "--id", "run_flag", "--params", `{"id":"run_json","title":"Run"}`)
	assertErrorCode(t, err, "input_conflict")
}

func TestMetadataMustBeObject(t *testing.T) {
	t.Setenv("AGENT_TEAM_STATE_DIR", t.TempDir())

	runCLI(t, "init")
	runCLI(t, "run", "create", "--id", "run_meta", "--title", "Metadata Run")
	err := runCLIFail(t,
		"task", "create",
		"--params", `{"id":"task_meta","run_id":"run_meta","agent":"worker","title":"Meta"}`,
		"--json", `{"body":"x","metadata":[]}`,
	)
	assertErrorCode(t, err, "invalid_json")

	err = runCLIFail(t,
		"task", "create",
		"--id", "task_meta_flag",
		"--run", "run_meta",
		"--agent", "worker",
		"--title", "Meta Flag",
		"--metadata", `[]`,
	)
	assertErrorCode(t, err, "invalid_json")
}

func TestMessageBlocksCompletionUntilAcked(t *testing.T) {
	t.Setenv("AGENT_TEAM_STATE_DIR", t.TempDir())

	runCLI(t, "init")
	runCLI(t, "run", "create", "--params", `{"id":"run_sync","title":"sync run"}`)
	runCLI(t, "task", "create", "--params", `{"id":"task_sync","run_id":"run_sync","agent":"worker","title":"do work"}`)
	runCLI(t, "task", "start", "--params", `{"task_id":"task_sync","agent":"worker"}`)
	runCLI(t, "message", "send",
		"--params", `{"id":"msg_sync","run_id":"run_sync","task_id":"task_sync","from":"planner","to":"worker","kind":"contract_changed"}`,
		"--json", `{"body":"contract changed"}`,
	)

	err := runCLIFail(t, "task", "complete",
		"--params", `{"task_id":"task_sync","agent":"worker","force":false}`,
		"--json", `{"evidence":"verified","artifact":"_workspace/run_sync/worker.md"}`,
	)
	if err == nil {
		t.Fatal("expected sync conflict")
	}

	runCLI(t, "inbox", "ack", "--params", `{"msg_id":"msg_sync","agent":"worker"}`)
	runCLI(t, "task", "complete",
		"--params", `{"task_id":"task_sync","agent":"worker","force":false}`,
		"--json", `{"evidence":"verified","artifact":"_workspace/run_sync/worker.md"}`,
	)
}

func TestStrictUnknownFields(t *testing.T) {
	t.Setenv("AGENT_TEAM_STATE_DIR", t.TempDir())
	runCLI(t, "init")
	err := runCLIFail(t, "run", "create", "--params", `{"title":"x","extra":true}`)
	if err == nil {
		t.Fatal("expected unknown field failure")
	}
}

func TestRunListAndClose(t *testing.T) {
	t.Setenv("AGENT_TEAM_STATE_DIR", t.TempDir())
	runCLI(t, "init")
	runCLI(t, "run", "create", "--id", "run_ops", "--title", "ops run")
	runCLI(t, "task", "create", "--id", "task_ops_a", "--run", "run_ops", "--agent", "alice", "--title", "first")
	runCLI(t, "task", "create", "--id", "task_ops_b", "--run", "run_ops", "--agent", "alice", "--title", "second")
	runCLI(t, "run", "list")
	list := runCLI(t, "run", "list")
	runs := list["data"].(map[string]any)["runs"].([]any)
	if len(runs) != 1 {
		t.Fatalf("expected one run, got %d", len(runs))
	}
	err := runCLIFail(t, "run", "close", "--run", "run_ops")
	assertErrorCode(t, err, "run_not_ready")
	runCLI(t, "task", "start", "--task", "task_ops_a", "--agent", "alice")
	runCLI(t, "task", "complete", "--task", "task_ops_a", "--agent", "alice", "--force", "--evidence", "ok", "--artifact", "_workspace/run_ops/task_ops_a.md")
	runCLI(t, "task", "start", "--task", "task_ops_b", "--agent", "alice")
	runCLI(t, "task", "complete", "--task", "task_ops_b", "--agent", "alice", "--force", "--evidence", "ok", "--artifact", "_workspace/run_ops/task_ops_b.md")
	closeOut := runCLI(t, "run", "close", "--run", "run_ops", "--reason", "all done")
	runData := closeOut["data"].(map[string]any)["run"].(map[string]any)
	if runData["status"] != "closed" {
		t.Fatalf("expected closed run, got %#v", runData["status"])
	}
	closeAgain := runCLI(t, "run", "close", "--run", "run_ops")
	warnings := closeAgain["warnings"].([]any)
	if len(warnings) == 0 {
		t.Fatalf("expected warning when closing already closed run")
	}
}

func TestTaskReassignAndRetry(t *testing.T) {
	t.Setenv("AGENT_TEAM_STATE_DIR", t.TempDir())
	runCLI(t, "init")
	runCLI(t, "run", "create", "--id", "run_task_ops", "--title", "task ops")
	runCLI(t, "task", "create", "--id", "task_ops_reassign", "--run", "run_task_ops", "--agent", "alice", "--title", "rewrite")
	reassign := runCLI(t, "task", "reassign", "--task", "task_ops_reassign", "--agent", "bob", "--reason", "handoff")
	taskPayload := reassign["data"].(map[string]any)["task"].(map[string]any)
	if taskPayload["agent"] != "bob" {
		t.Fatalf("expected reassigned agent bob, got %#v", taskPayload["agent"])
	}
	runCLI(t, "task", "block", "--task", "task_ops_reassign", "--agent", "bob", "--reason", "dependency missing")
	retry := runCLI(t, "task", "retry", "--task", "task_ops_reassign", "--reason", "retry after fix")
	taskPayload = retry["data"].(map[string]any)["task"].(map[string]any)
	if taskPayload["status"] != "pending" {
		t.Fatalf("expected retry to reset to pending, got %#v", taskPayload["status"])
	}
}

func TestMessageListByRun(t *testing.T) {
	t.Setenv("AGENT_TEAM_STATE_DIR", t.TempDir())
	runCLI(t, "init")
	runCLI(t, "run", "create", "--id", "run_msg", "--title", "msg run")
	runCLI(t, "message", "send", "--id", "msg_a", "--run", "run_msg", "--from", "planner", "--to", "worker", "--kind", "notice", "--body", "first")
	runCLI(t, "message", "send", "--id", "msg_b", "--run", "run_msg", "--from", "planner", "--to", "worker", "--kind", "update", "--body", "second")
	all := runCLI(t, "message", "list", "--run", "run_msg")
	messages := all["data"].(map[string]any)["messages"].([]any)
	if len(messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(messages))
	}
	unread := runCLI(t, "message", "list", "--run", "run_msg", "--unread")
	unreadMessages := unread["data"].(map[string]any)["messages"].([]any)
	if len(unreadMessages) != 2 {
		t.Fatalf("expected 2 unread messages, got %d", len(unreadMessages))
	}
	kind := runCLI(t, "message", "list", "--run", "run_msg", "--kind", "update")
	kindMessages := kind["data"].(map[string]any)["messages"].([]any)
	if len(kindMessages) != 1 {
		t.Fatalf("expected 1 update message, got %d", len(kindMessages))
	}
}

func TestEventLog(t *testing.T) {
	t.Setenv("AGENT_TEAM_STATE_DIR", t.TempDir())
	runCLI(t, "init")
	runCLI(t, "run", "create", "--id", "run_evt", "--title", "evt run")
	runCLI(t, "task", "create", "--id", "task_evt", "--run", "run_evt", "--agent", "worker", "--title", "task evt")
	runCLI(t, "task", "start", "--task", "task_evt", "--agent", "worker")
	runCLI(t, "task", "complete", "--task", "task_evt", "--agent", "worker", "--force", "--evidence", "ok", "--artifact", "_workspace/run_evt/task_evt.md")
	logOutput := runCLI(t, "event", "log", "--run", "run_evt", "--limit", "20")
	events := logOutput["data"].(map[string]any)["events"].([]any)
	if len(events) < 3 {
		t.Fatalf("expected at least 3 events, got %d", len(events))
	}
}

func TestSchemaExport(t *testing.T) {
	run := runCLI(t, "schema", "export")
	data := run["data"].(map[string]any)
	schema := data["schema"].(map[string]any)
	if schema["command"] != "agent-team" {
		t.Fatalf("unexpected schema command: %#v", schema["command"])
	}
	commands := schema["commands"].([]any)
	if len(commands) == 0 {
		t.Fatalf("expected schema commands")
	}
	recovery := schema["error_recovery"].([]any)
	if len(recovery) == 0 {
		t.Fatalf("expected error recovery specs")
	}
	foundSync := false
	for _, entry := range recovery {
		payload := entry.(map[string]any)
		if payload["code"] != "sync_conflict" {
			continue
		}
		foundSync = true
		guide := payload["recovery"].(map[string]any)
		if len(guide["actions"].([]any)) == 0 || len(guide["skills"].([]any)) == 0 {
			t.Fatalf("sync_conflict recovery is incomplete: %#v", guide)
		}
	}
	if !foundSync {
		t.Fatalf("expected sync_conflict recovery spec")
	}
}

func TestRunSummaryAndTaskStale(t *testing.T) {
	t.Setenv("AGENT_TEAM_STATE_DIR", t.TempDir())
	runCLI(t, "init")
	runCLI(t, "run", "create", "--id", "run_summary", "--title", "summary run")
	runCLI(t, "task", "create", "--id", "task_stale", "--run", "run_summary", "--agent", "worker", "--title", "long task")
	runCLI(t, "task", "start", "--task", "task_stale", "--agent", "worker")
	time.Sleep(2 * time.Millisecond)

	summary := runCLI(t, "run", "summary", "--run", "run_summary")
	summaryData := summary["data"].(map[string]any)["summary"].(map[string]any)
	if summaryData["close_ready"].(bool) {
		t.Fatalf("expected run not to be close ready")
	}
	inProgress := summaryData["in_progress_tasks"].([]any)
	if len(inProgress) != 1 {
		t.Fatalf("expected one in-progress task, got %d", len(inProgress))
	}

	stale := runCLI(t, "task", "stale", "--run", "run_summary", "--older-than", "1ns")
	staleTasks := stale["data"].(map[string]any)["stale_tasks"].([]any)
	if len(staleTasks) != 1 {
		t.Fatalf("expected one stale task, got %d", len(staleTasks))
	}
}

func TestCancelFailAndTerminalClose(t *testing.T) {
	t.Setenv("AGENT_TEAM_STATE_DIR", t.TempDir())
	runCLI(t, "init")
	runCLI(t, "run", "create", "--id", "run_terminal", "--title", "terminal run")
	runCLI(t, "task", "create", "--id", "task_cancel", "--run", "run_terminal", "--agent", "worker", "--title", "cancel me")
	runCLI(t, "task", "create", "--id", "task_fail", "--run", "run_terminal", "--agent", "worker", "--title", "fail me")
	cancelled := runCLI(t, "task", "cancel", "--task", "task_cancel", "--reason", "not needed")
	cancelledTask := cancelled["data"].(map[string]any)["task"].(map[string]any)
	if cancelledTask["status"] != "cancelled" {
		t.Fatalf("expected cancelled task, got %#v", cancelledTask["status"])
	}
	failed := runCLI(t, "task", "fail", "--task", "task_fail", "--agent", "worker", "--reason", "failed verification", "--artifact", "_workspace/run_terminal/task_fail.md")
	failedTask := failed["data"].(map[string]any)["task"].(map[string]any)
	if failedTask["status"] != "failed" {
		t.Fatalf("expected failed task, got %#v", failedTask["status"])
	}
	closed := runCLI(t, "run", "close", "--run", "run_terminal", "--reason", "all tasks terminal")
	runPayload := closed["data"].(map[string]any)["run"].(map[string]any)
	if runPayload["status"] != "closed" {
		t.Fatalf("expected closed run, got %#v", runPayload["status"])
	}
}

func TestStatusPolicyTransitions(t *testing.T) {
	t.Setenv("AGENT_TEAM_STATE_DIR", t.TempDir())
	runCLI(t, "init")

	runCLI(t, "run", "create", "--id", "run_policy", "--title", "policy run")
	runCLI(t, "task", "create", "--id", "task_done", "--run", "run_policy", "--agent", "worker", "--title", "done task")
	runCLI(t, "task", "create", "--id", "task_failed", "--run", "run_policy", "--agent", "worker", "--title", "failed task")
	runCLI(t, "task", "create", "--id", "task_cancelled", "--run", "run_policy", "--agent", "worker", "--title", "cancelled task")

	runCLI(t, "task", "start", "--task", "task_done", "--agent", "worker")
	runCLI(t, "task", "complete", "--task", "task_done", "--agent", "worker", "--force", "--evidence", "ok", "--artifact", "_workspace/run_policy/task_done.md")
	runCLI(t, "task", "fail", "--task", "task_failed", "--agent", "worker", "--reason", "verification failed")
	runCLI(t, "task", "cancel", "--task", "task_cancelled", "--reason", "no longer needed")

	err := runCLIFail(t, "task", "retry", "--task", "task_done", "--reason", "should not retry done")
	assertErrorCode(t, err, "invalid_task_state")

	retried := runCLI(t, "task", "retry", "--task", "task_failed", "--reason", "retry failed work")
	retriedTask := retried["data"].(map[string]any)["task"].(map[string]any)
	if retriedTask["status"] != "pending" {
		t.Fatalf("expected failed retry to reset pending, got %#v", retriedTask["status"])
	}
	runCLI(t, "task", "cancel", "--task", "task_failed", "--reason", "retired after retry")

	err = runCLIFail(t, "task", "reassign", "--task", "task_done", "--agent", "reviewer", "--reason", "should not reassign done")
	assertErrorCode(t, err, "invalid_task_state")

	closed := runCLI(t, "run", "close", "--run", "run_policy", "--reason", "all outcomes terminal")
	runPayload := closed["data"].(map[string]any)["run"].(map[string]any)
	if runPayload["status"] != "closed" {
		t.Fatalf("expected closed run, got %#v", runPayload["status"])
	}

	err = runCLIFail(t, "run", "cancel", "--run", "run_policy", "--reason", "should not cancel closed")
	assertErrorCode(t, err, "invalid_run_state")
}

func TestListPagination(t *testing.T) {
	t.Setenv("AGENT_TEAM_STATE_DIR", t.TempDir())
	runCLI(t, "init")
	runCLI(t, "run", "create", "--id", "run_page", "--title", "page run")
	runCLI(t, "message", "send", "--id", "msg_page_a", "--run", "run_page", "--from", "planner", "--to", "worker", "--kind", "notice", "--body", "first")
	runCLI(t, "message", "send", "--id", "msg_page_b", "--run", "run_page", "--from", "planner", "--to", "worker", "--kind", "notice", "--body", "second")
	runCLI(t, "message", "send", "--id", "msg_page_c", "--run", "run_page", "--from", "planner", "--to", "worker", "--kind", "notice", "--body", "third")
	page := runCLI(t, "message", "list", "--run", "run_page", "--limit", "2")
	messages := page["data"].(map[string]any)["messages"].([]any)
	if len(messages) != 2 {
		t.Fatalf("expected page size 2, got %d", len(messages))
	}
}

func assertErrorCode(t *testing.T, err error, code string) {
	t.Helper()
	coder, ok := any(err).(interface {
		Code() string
	})
	if !ok {
		t.Fatalf("expected coded error, got %T: %v", err, err)
	}
	if coder.Code() != code {
		t.Fatalf("expected error code %s, got %s: %v", code, coder.Code(), err)
	}
}
