package cli

import "github.com/tae2089/agent-team/internal/output"

type SchemaSpec struct {
	Command       string                     `json:"command"`
	Version       string                     `json:"version"`
	Commands      []CommandSpec              `json:"commands"`
	Errors        []string                   `json:"errors"`
	ErrorRecovery []output.ErrorRecoverySpec `json:"error_recovery"`
}

type CommandSpec struct {
	Name        string           `json:"name"`
	Subcommands []SubcommandSpec `json:"subcommands"`
}

type SubcommandSpec struct {
	Name            string     `json:"name"`
	Flags           []FlagSpec `json:"flags"`
	Params          []string   `json:"params"`
	RequiredParams  []string   `json:"required_params,omitempty"`
	Output          []string   `json:"output"`
	Supports        []string   `json:"supports,omitempty"`
	OutputWarnings  []string   `json:"output_warnings,omitempty"`
	ValidFromStatus []string   `json:"valid_from_status,omitempty"`
	HelperSkill     string     `json:"helper_skill,omitempty"`
	ServiceSkill    string     `json:"service_skill,omitempty"`
}

type FlagSpec struct {
	Name    string `json:"name"`
	JSONKey string `json:"json_key"`
	Default string `json:"default,omitempty"`
}

func cliSchema() SchemaSpec {
	spec := SchemaSpec{
		Command: "agent-team",
		Version: "1.3.0",
		Commands: []CommandSpec{
			{
				Name: "init",
				Subcommands: []SubcommandSpec{
					sub("", "", "", nil, nil, nil, []string{"state_dir", "db_path"}),
				},
			},
			{
				Name: "run",
				Subcommands: []SubcommandSpec{
					sub("create", "agent-team-run", "agent-team-run-create", []FlagSpec{flag("id", "id", "generated"), flag("title", "title", "")}, []string{"id", "title"}, []string{"title"}, []string{"run"}),
					sub("status", "agent-team-run", "agent-team-run-status", []FlagSpec{flag("run", "run_id", "")}, []string{"run_id"}, []string{"run_id"}, []string{"run", "tasks"}),
					sub("summary", "agent-team-run", "agent-team-run-summary", []FlagSpec{flag("run", "run_id", ""), flag("recent-limit", "recent_limit", "10")}, []string{"run_id", "recent_limit"}, []string{"run_id"}, []string{"summary"}),
					sub("list", "agent-team-run", "agent-team-run-list", []FlagSpec{flag("status", "status", ""), flag("limit", "limit", "100"), flag("after-version", "after_version", "0")}, []string{"status", "limit", "after_version"}, nil, []string{"runs"}),
					sub("close", "agent-team-run", "agent-team-run-close", []FlagSpec{flag("run", "run_id", ""), flag("reason", "reason", "")}, []string{"run_id", "reason"}, []string{"run_id"}, []string{"run"}),
					sub("cancel", "agent-team-run", "agent-team-run-cancel", []FlagSpec{flag("run", "run_id", ""), flag("reason", "reason", "")}, []string{"run_id", "reason"}, []string{"run_id", "reason"}, []string{"run"}),
				},
			},
			{
				Name: "task",
				Subcommands: []SubcommandSpec{
					sub("create", "agent-team-task", "agent-team-task-create", []FlagSpec{flag("id", "id", "generated"), flag("run", "run_id", ""), flag("agent", "agent", ""), flag("title", "title", ""), flag("depends-on", "depends_on", "[]"), flag("body", "body", ""), flag("metadata", "metadata", "{}")}, []string{"id", "run_id", "agent", "title", "depends_on", "body", "metadata"}, []string{"run_id", "agent", "title"}, []string{"task"}),
					sub("list", "agent-team-task", "agent-team-task-list", []FlagSpec{flag("run", "run_id", ""), flag("agent", "agent", ""), flag("status", "status", ""), flag("limit", "limit", "100"), flag("after-version", "after_version", "0")}, []string{"run_id", "agent", "status", "limit", "after_version"}, nil, []string{"tasks"}),
					sub("show", "agent-team-task", "agent-team-task-show", []FlagSpec{flag("task", "task_id", "")}, []string{"task_id"}, []string{"task_id"}, []string{"task", "depends_on"}),
					sub("start", "agent-team-task", "agent-team-task-start", []FlagSpec{flag("task", "task_id", ""), flag("agent", "agent", "")}, []string{"task_id", "agent"}, []string{"task_id", "agent"}, []string{"task"}),
					sub("complete", "agent-team-task", "agent-team-task-complete", []FlagSpec{flag("task", "task_id", ""), flag("agent", "agent", ""), flag("force", "force", "false"), flag("evidence", "evidence", ""), flag("artifact", "artifact", "")}, []string{"task_id", "agent", "force", "evidence", "artifact"}, []string{"task_id", "agent", "evidence", "artifact"}, []string{"task"}),
					sub("block", "agent-team-task", "agent-team-task-block", []FlagSpec{flag("task", "task_id", ""), flag("agent", "agent", ""), flag("reason", "reason", "")}, []string{"task_id", "agent", "reason"}, []string{"task_id", "agent", "reason"}, []string{"task"}),
					sub("reassign", "agent-team-task", "agent-team-task-reassign", []FlagSpec{flag("task", "task_id", ""), flag("agent", "agent", ""), flag("reason", "reason", "")}, []string{"task_id", "agent", "reason"}, []string{"task_id", "agent", "reason"}, []string{"task"}),
					sub("retry", "agent-team-task", "agent-team-task-retry", []FlagSpec{flag("task", "task_id", ""), flag("reason", "reason", "")}, []string{"task_id", "reason"}, []string{"task_id", "reason"}, []string{"task"}),
					sub("cancel", "agent-team-task", "agent-team-task-cancel", []FlagSpec{flag("task", "task_id", ""), flag("reason", "reason", "")}, []string{"task_id", "reason"}, []string{"task_id", "reason"}, []string{"task"}),
					sub("fail", "agent-team-task", "agent-team-task-fail", []FlagSpec{flag("task", "task_id", ""), flag("agent", "agent", ""), flag("reason", "reason", ""), flag("artifact", "artifact", "")}, []string{"task_id", "agent", "reason", "artifact"}, []string{"task_id", "agent", "reason"}, []string{"task"}),
					sub("stale", "agent-team-task", "agent-team-task-stale", []FlagSpec{flag("run", "run_id", ""), flag("older-than", "older_than", ""), flag("limit", "limit", "100"), flag("after-version", "after_version", "0")}, []string{"run_id", "older_than", "limit", "after_version"}, []string{"run_id", "older_than"}, []string{"stale_tasks"}),
				},
			},
			{
				Name: "message",
				Subcommands: []SubcommandSpec{
					sub("send", "agent-team-inbox", "agent-team-message-send", []FlagSpec{flag("id", "id", "generated"), flag("run", "run_id", ""), flag("task", "task_id", ""), flag("from", "from", ""), flag("to", "to", ""), flag("kind", "kind", ""), flag("body", "body", ""), flag("metadata", "metadata", "{}")}, []string{"id", "run_id", "task_id", "from", "to", "kind", "body", "metadata"}, []string{"run_id", "from", "to", "kind", "body"}, []string{"message"}),
					sub("list", "agent-team-inbox", "agent-team-message-list", []FlagSpec{flag("run", "run_id", ""), flag("task", "task_id", ""), flag("from", "from", ""), flag("to", "to", ""), flag("kind", "kind", ""), flag("unread", "unread", "false"), flag("limit", "limit", "100"), flag("after-version", "after_version", "0")}, []string{"run_id", "task_id", "from", "to", "kind", "unread", "limit", "after_version"}, []string{"run_id"}, []string{"messages"}),
				},
			},
			{
				Name: "inbox",
				Subcommands: []SubcommandSpec{
					sub("list", "agent-team-inbox", "agent-team-inbox-list", []FlagSpec{flag("agent", "agent", ""), flag("run", "run_id", ""), flag("unread", "unread", "false")}, []string{"agent", "run_id", "unread"}, []string{"agent"}, []string{"messages"}),
					sub("ack", "agent-team-inbox", "agent-team-inbox-ack", []FlagSpec{flag("msg", "msg_id", ""), flag("agent", "agent", "")}, []string{"msg_id", "agent"}, []string{"msg_id", "agent"}, []string{"message"}),
				},
			},
			{
				Name: "sync",
				Subcommands: []SubcommandSpec{
					sub("check", "agent-team-sync", "agent-team-sync-check", []FlagSpec{flag("agent", "agent", ""), flag("run", "run_id", ""), flag("task", "task_id", "")}, []string{"agent", "run_id", "task_id"}, []string{"agent"}, []string{"sync"}),
				},
			},
			{
				Name: "event",
				Subcommands: []SubcommandSpec{
					sub("log", "agent-team-ops", "agent-team-event-log", []FlagSpec{flag("run", "run_id", ""), flag("entity-type", "entity_type", ""), flag("entity", "entity_id", ""), flag("type", "event_type", ""), flag("after-version", "after_version", "0"), flag("limit", "limit", "100")}, []string{"run_id", "entity_type", "entity_id", "event_type", "after_version", "limit"}, nil, []string{"events"}),
				},
			},
			{
				Name: "schema",
				Subcommands: []SubcommandSpec{
					sub("export", "agent-team-ops", "agent-team-schema-export", nil, nil, nil, []string{"schema"}),
				},
			},
			{
				Name: "version",
				Subcommands: []SubcommandSpec{
					sub("", "", "", nil, nil, nil, []string{"version"}),
				},
			},
		},
		Errors: []string{
			"validation_error",
			"input_conflict",
			"invalid_json",
			"invalid_json_source",
			"not_found",
			"agent_mismatch",
			"sync_conflict",
			"run_not_ready",
			"invalid_run_state",
			"invalid_task_state",
			"internal_error",
		},
	}
	spec.ErrorRecovery = output.ErrorRecoveryCatalog(spec.Errors)
	return spec
}

func sub(name, serviceSkill, helperSkill string, flags []FlagSpec, params, required, output []string) SubcommandSpec {
	return SubcommandSpec{
		Name:           name,
		Flags:          flags,
		Params:         params,
		RequiredParams: required,
		Output:         output,
		ServiceSkill:   serviceSkill,
		HelperSkill:    helperSkill,
	}
}

func flag(name, jsonKey, def string) FlagSpec {
	return FlagSpec{Name: name, JSONKey: jsonKey, Default: def}
}
