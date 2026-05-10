package output

import (
	"encoding/json"
	"fmt"
	"io"
)

type Response struct {
	OK           bool     `json:"ok"`
	StateVersion int64    `json:"state_version"`
	Data         any      `json:"data,omitempty"`
	Warnings     []string `json:"warnings"`
}

type ErrorResponse struct {
	OK           bool      `json:"ok"`
	StateVersion int64     `json:"state_version"`
	Error        ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code     string         `json:"code"`
	Message  string         `json:"message"`
	Details  any            `json:"details,omitempty"`
	Recovery *RecoveryGuide `json:"recovery,omitempty"`
}

type RecoveryGuide struct {
	Summary  string   `json:"summary"`
	Actions  []string `json:"actions"`
	Commands []string `json:"commands,omitempty"`
	Docs     []string `json:"docs,omitempty"`
	Skills   []string `json:"skills,omitempty"`
}

type ErrorRecoverySpec struct {
	Code     string        `json:"code"`
	Recovery RecoveryGuide `json:"recovery"`
}

type CodedError interface {
	error
	Code() string
	Details() any
}

type AppError struct {
	ErrCode string
	Msg     string
	Extra   any
}

func (e AppError) Error() string {
	return e.Msg
}

func (e AppError) Code() string {
	return e.ErrCode
}

func (e AppError) Details() any {
	return e.Extra
}

func NewError(code, message string, details any) AppError {
	return AppError{ErrCode: code, Msg: message, Extra: details}
}

func Write(w io.Writer, stateVersion int64, data any, warnings []string) error {
	if warnings == nil {
		warnings = []string{}
	}
	return encode(w, Response{
		OK:           true,
		StateVersion: stateVersion,
		Data:         data,
		Warnings:     warnings,
	})
}

func WriteError(w io.Writer, stateVersion int64, err error) {
	code := "internal_error"
	details := any(nil)
	if coded, ok := err.(CodedError); ok {
		code = coded.Code()
		details = coded.Details()
	}
	_ = encode(w, ErrorResponse{
		OK:           false,
		StateVersion: stateVersion,
		Error: ErrorBody{
			Code:     code,
			Message:  err.Error(),
			Details:  details,
			Recovery: RecoveryFor(code),
		},
	})
}

func RecoveryFor(code string) *RecoveryGuide {
	recovery, ok := recoveryCatalog[code]
	if !ok {
		recovery = recoveryCatalog["internal_error"]
	}
	return &RecoveryGuide{
		Summary:  recovery.Summary,
		Actions:  append([]string(nil), recovery.Actions...),
		Commands: append([]string(nil), recovery.Commands...),
		Docs:     append([]string(nil), recovery.Docs...),
		Skills:   append([]string(nil), recovery.Skills...),
	}
}

func ErrorRecoveryCatalog(codes []string) []ErrorRecoverySpec {
	specs := make([]ErrorRecoverySpec, 0, len(codes))
	for _, code := range codes {
		recovery := RecoveryFor(code)
		specs = append(specs, ErrorRecoverySpec{Code: code, Recovery: *recovery})
	}
	return specs
}

var recoveryCatalog = map[string]RecoveryGuide{
	"validation_error": {
		Summary: "Fix missing or invalid command input.",
		Actions: []string{
			"Check the command helper skill or --help output for required fields.",
			"Provide IDs and routing fields through flags or --params.",
			"Provide rich text, metadata, evidence, artifact paths, or reasons through --json when useful.",
		},
		Commands: []string{"agent-team COMMAND --help", "agent-team schema export"},
		Docs:     []string{"docs/errors.md#validation_error"},
		Skills:   []string{"agent-team-shared"},
	},
	"input_conflict": {
		Summary: "Remove duplicate values supplied by both flags and JSON.",
		Actions: []string{
			"Pick either the named flag or the JSON key for the conflicting field.",
			"Keep IDs and filters in --params when using JSON-heavy commands.",
		},
		Commands: []string{"agent-team schema export"},
		Docs:     []string{"docs/errors.md#input_conflict"},
		Skills:   []string{"agent-team-shared"},
	},
	"invalid_json": {
		Summary: "Fix malformed JSON, unknown fields, or wrong value types.",
		Actions: []string{
			"Validate the JSON object before retrying.",
			"Use keys from the command helper skill or schema export only.",
			"Use an object for metadata and other structured payloads.",
		},
		Commands: []string{"agent-team schema export"},
		Docs:     []string{"docs/errors.md#invalid_json"},
		Skills:   []string{"agent-team-shared"},
	},
	"invalid_json_source": {
		Summary: "Fix the JSON source reference.",
		Actions: []string{
			"Use inline JSON, @file, or - for stdin.",
			"Ensure @file paths are non-empty and readable.",
			"Ensure stdin is available when using -.",
		},
		Commands: []string{"agent-team COMMAND --json @payload.json", "agent-team COMMAND --json -"},
		Docs:     []string{"docs/errors.md#invalid_json_source"},
		Skills:   []string{"agent-team-shared"},
	},
	"not_found": {
		Summary: "Verify the run, task, or message ID and state directory.",
		Actions: []string{
			"List nearby runs, tasks, or messages to find the correct ID.",
			"Confirm AGENT_TEAM_STATE_DIR points at the expected workflow state.",
			"Initialize state when this is a new workspace.",
		},
		Commands: []string{"agent-team run list", "agent-team task list --run RUN_ID", "agent-team message list --run RUN_ID"},
		Docs:     []string{"docs/errors.md#not_found"},
		Skills:   []string{"agent-team-shared", "agent-team-run-list", "agent-team-task-list", "agent-team-message-list"},
	},
	"agent_mismatch": {
		Summary: "Use the assigned agent or ask the orchestrator to reassign ownership.",
		Actions: []string{
			"Inspect the task or message to find the assigned owner.",
			"Retry as the assigned agent when appropriate.",
			"Use task reassign only from the orchestrator or owner handoff path.",
		},
		Commands: []string{"agent-team task show --task TASK_ID", "agent-team task reassign --task TASK_ID --agent AGENT --reason REASON"},
		Docs:     []string{"docs/errors.md#agent_mismatch"},
		Skills:   []string{"agent-team-task-show", "agent-team-task-reassign"},
	},
	"sync_conflict": {
		Summary: "Resolve inbox or dependency drift before completing the task.",
		Actions: []string{
			"Run sync check for the task and assigned agent.",
			"Ack unread relevant messages after reading them.",
			"Wait for or complete unfinished dependencies.",
			"Use --force only when the orchestrator explicitly approves it.",
		},
		Commands: []string{"agent-team sync check --agent AGENT --run RUN_ID --task TASK_ID", "agent-team inbox list --agent AGENT --run RUN_ID --unread", "agent-team inbox ack --msg MSG_ID --agent AGENT"},
		Docs:     []string{"docs/errors.md#sync_conflict"},
		Skills:   []string{"agent-team-sync-check", "agent-team-inbox-list", "agent-team-inbox-ack", "agent-team-task-complete"},
	},
	"run_not_ready": {
		Summary: "Finish or terminally resolve all tasks before closing the run.",
		Actions: []string{
			"Inspect run summary for pending, in-progress, or blocked tasks.",
			"Complete, retry, reassign, cancel, fail, or block work as policy allows.",
			"Close the run only after summary reports close readiness.",
		},
		Commands: []string{"agent-team run summary --run RUN_ID", "agent-team task list --run RUN_ID"},
		Docs:     []string{"docs/errors.md#run_not_ready"},
		Skills:   []string{"agent-team-run-summary", "agent-team-task-list", "agent-team-run-close"},
	},
	"invalid_run_state": {
		Summary: "Choose an operation allowed by the current run status.",
		Actions: []string{
			"Inspect the run status.",
			"Do not cancel a closed run.",
			"Create or resume an open run for new work.",
		},
		Commands: []string{"agent-team run status --run RUN_ID", "agent-team run list --status open"},
		Docs:     []string{"docs/errors.md#invalid_run_state"},
		Skills:   []string{"agent-team-run-status", "agent-team-run-list"},
	},
	"invalid_task_state": {
		Summary: "Choose a task transition allowed by the current task status.",
		Actions: []string{
			"Inspect the task status and dependencies.",
			"Use retry for blocked, in-progress, or failed work.",
			"Use reassign for pending or blocked work.",
			"Do not mutate terminal tasks except through run finalization.",
		},
		Commands: []string{"agent-team task show --task TASK_ID", "agent-team task retry --task TASK_ID --reason REASON", "agent-team task reassign --task TASK_ID --agent AGENT --reason REASON"},
		Docs:     []string{"docs/errors.md#invalid_task_state"},
		Skills:   []string{"agent-team-task-show", "agent-team-task-retry", "agent-team-task-reassign"},
	},
	"internal_error": {
		Summary: "Check local state, filesystem permissions, and installation health.",
		Actions: []string{
			"Confirm the state directory exists and is writable.",
			"Set AGENT_TEAM_STATE_DIR to an isolated writable directory when sandboxed.",
			"Run agent-team init for a new state directory.",
			"Check install and sandbox permission docs when running through Codex or Gemini.",
		},
		Commands: []string{"agent-team init", "agent-team version"},
		Docs:     []string{"docs/errors.md#internal_error", "docs/install.md", "docs/harness-sandbox-permissions.md"},
		Skills:   []string{"agent-team-shared"},
	},
}

func encode(w io.Writer, value any) error {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(value); err != nil {
		return fmt.Errorf("encode json response: %w", err)
	}
	return nil
}
