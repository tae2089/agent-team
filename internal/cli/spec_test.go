package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommandSpecMatchesCobraTree(t *testing.T) {
	root := NewRoot()
	actual := map[string]map[string]bool{}
	for _, command := range root.Commands() {
		if command.Hidden {
			continue
		}
		children := map[string]bool{}
		for _, child := range command.Commands() {
			if child.Hidden {
				continue
			}
			children[strings.Fields(child.Use)[0]] = true
		}
		actual[strings.Fields(command.Use)[0]] = children
	}

	for _, command := range cliSchema().Commands {
		children, ok := actual[command.Name]
		if !ok {
			t.Fatalf("command %s exists in spec but not cobra tree", command.Name)
		}
		for _, subcommand := range command.Subcommands {
			if subcommand.Name == "" {
				continue
			}
			if !children[subcommand.Name] {
				t.Fatalf("subcommand %s %s exists in spec but not cobra tree", command.Name, subcommand.Name)
			}
		}
	}
}

func TestHelperSkillDocsMatchSpec(t *testing.T) {
	skillsDir := filepath.Join("..", "..", "skills")
	requiredSections := []string{"## Usage", "## Flags", "## Examples", "## Errors", "## See Also"}
	for _, command := range cliSchema().Commands {
		for _, subcommand := range command.Subcommands {
			if subcommand.HelperSkill == "" {
				continue
			}
			path := filepath.Join(skillsDir, subcommand.HelperSkill, "SKILL.md")
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read helper skill %s: %v", subcommand.HelperSkill, err)
			}
			body := string(raw)
			if got := frontmatterValue(body, "name"); got != subcommand.HelperSkill {
				t.Fatalf("%s frontmatter name = %q, want %q", path, got, subcommand.HelperSkill)
			}
			wantCLIHelp := commandHelp(command.Name, subcommand.Name)
			if got := frontmatterValue(body, "cliHelp"); got != wantCLIHelp {
				t.Fatalf("%s cliHelp = %q, want %q", path, got, wantCLIHelp)
			}
			for _, section := range requiredSections {
				if !strings.Contains(body, section) {
					t.Fatalf("%s missing %s", path, section)
				}
			}
			if !strings.Contains(body, "| Flag | JSON key | Required | Default | Description |") {
				t.Fatalf("%s missing standard flags table", path)
			}
			for _, flag := range subcommand.Flags {
				if !strings.Contains(body, "`--"+flag.Name+"`") {
					t.Fatalf("%s missing flag --%s", path, flag.Name)
				}
				if flag.JSONKey != "" && !strings.Contains(body, "`"+flag.JSONKey+"`") {
					t.Fatalf("%s missing json key %s", path, flag.JSONKey)
				}
			}
			assertFlagTableMatchesSpec(t, path, body, subcommand)
		}
	}
}

func TestServiceSkillDocsExist(t *testing.T) {
	skillsDir := filepath.Join("..", "..", "skills")
	seen := map[string]bool{}
	for _, command := range cliSchema().Commands {
		for _, subcommand := range command.Subcommands {
			if subcommand.ServiceSkill == "" || seen[subcommand.ServiceSkill] {
				continue
			}
			seen[subcommand.ServiceSkill] = true
			path := filepath.Join(skillsDir, subcommand.ServiceSkill, "SKILL.md")
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read service skill %s: %v", subcommand.ServiceSkill, err)
			}
			if got := frontmatterValue(string(raw), "name"); got != subcommand.ServiceSkill {
				t.Fatalf("%s frontmatter name = %q, want %q", path, got, subcommand.ServiceSkill)
			}
		}
	}
}

func TestSkillDocsAvoidInvalidBooleanShellExamples(t *testing.T) {
	skillsDir := filepath.Join("..", "..", "skills")
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		path := filepath.Join(skillsDir, entry.Name(), "SKILL.md")
		raw, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		body := string(raw)
		if strings.Contains(body, "--force true") || strings.Contains(body, "--unread true") {
			t.Fatalf("%s contains invalid boolean shell example", path)
		}
	}
}

func TestErrorRecoveryDocsMatchSpec(t *testing.T) {
	root := filepath.Join("..", "..")
	errorsDoc := readString(t, filepath.Join(root, "docs", "errors.md"))
	sharedSkill := readString(t, filepath.Join(root, "skills", "agent-team-shared", "SKILL.md"))
	schema := cliSchema()
	if len(schema.ErrorRecovery) != len(schema.Errors) {
		t.Fatalf("schema has %d error recovery entries, want %d", len(schema.ErrorRecovery), len(schema.Errors))
	}
	recoveryByCode := map[string]bool{}
	for _, spec := range schema.ErrorRecovery {
		if spec.Recovery.Summary == "" {
			t.Fatalf("%s has empty recovery summary", spec.Code)
		}
		if len(spec.Recovery.Actions) == 0 {
			t.Fatalf("%s has no recovery actions", spec.Code)
		}
		if len(spec.Recovery.Docs) == 0 {
			t.Fatalf("%s has no recovery docs", spec.Code)
		}
		if len(spec.Recovery.Skills) == 0 {
			t.Fatalf("%s has no recovery skills", spec.Code)
		}
		recoveryByCode[spec.Code] = true
	}
	for _, code := range schema.Errors {
		if !recoveryByCode[code] {
			t.Fatalf("%s missing schema recovery entry", code)
		}
		if !strings.Contains(errorsDoc, "## "+code) {
			t.Fatalf("docs/errors.md missing heading for %s", code)
		}
		if !strings.Contains(sharedSkill, "`"+code+"`") {
			t.Fatalf("agent-team-shared missing error code %s", code)
		}
	}
	for _, required := range []string{"error.recovery", "docs/errors.md", "agent-team schema export"} {
		if !strings.Contains(sharedSkill, required) {
			t.Fatalf("agent-team-shared missing %q", required)
		}
	}
}

func TestHarnessSandboxPermissionDocsAreWired(t *testing.T) {
	root := filepath.Join("..", "..")
	doc := readString(t, filepath.Join(root, "docs", "harness-sandbox-permissions.md"))
	requiredDoc := []string{
		"Codex",
		"Gemini",
		"sandbox_mode",
		"tools",
		"workspace-write",
		"danger-full-access",
		"invoke_agent",
		"run_shell_command",
		"AGENT_TEAM_STATE_DIR",
	}
	for _, want := range requiredDoc {
		if !strings.Contains(doc, want) {
			t.Fatalf("harness sandbox docs missing %q", want)
		}
	}
	cases := []struct {
		path string
		want []string
	}{
		{
			path: filepath.Join(root, "skills", "agent-team-codex-harness", "references", "agent-design-patterns.md"),
			want: []string{"docs/harness-sandbox-permissions.md", "workspace-write", "danger-full-access", "AGENT_TEAM_STATE_DIR"},
		},
		{
			path: filepath.Join(root, "skills", "agent-team-gemini-harness", "references", "agent-design-patterns.md"),
			want: []string{"docs/harness-sandbox-permissions.md", "invoke_agent", "run_shell_command", "wildcard tool permissions"},
		},
		{
			path: filepath.Join(root, "skills", "agent-team-codex-harness", "references", "schemas", "agent-worker.template.toml"),
			want: []string{"Sandbox permission must match the role", "error.recovery"},
		},
		{
			path: filepath.Join(root, "skills", "agent-team-gemini-harness", "references", "schemas", "agent-worker.template.md"),
			want: []string{"Tool permissions must match the role", "error.recovery"},
		},
	}
	for _, tc := range cases {
		body := readString(t, tc.path)
		for _, want := range tc.want {
			if !strings.Contains(body, want) {
				t.Fatalf("%s missing %q", tc.path, want)
			}
		}
	}
}

func TestHarnessRuntimeContextDocsAreWired(t *testing.T) {
	root := filepath.Join("..", "..")
	doc := readString(t, filepath.Join(root, "docs", "harness-runtime-context.md"))
	requiredDoc := []string{
		"should not need to provide `RUN_ID` or `TASK_ID`",
		"orchestrator owns them",
		"captures them from JSON command output",
		"capture `.data.run.id`",
		"capture `.data.task.id`",
		"not ask ordinary users to paste raw runtime IDs",
		"Workers do not create, infer, or ask for runtime IDs",
	}
	for _, want := range requiredDoc {
		if !strings.Contains(doc, want) {
			t.Fatalf("harness runtime context docs missing %q", want)
		}
	}

	cases := []struct {
		path string
		want []string
	}{
		{
			path: filepath.Join(root, "README.md"),
			want: []string{"docs/harness-runtime-context.md", "Generated Codex/Gemini harness orchestrators create or resolve those IDs"},
		},
		{
			path: filepath.Join(root, "skills", "agent-team-codex-harness", "references", "agent-team-runtime-protocol.md"),
			want: []string{"orchestrator-owned internal context", "capture `.data.run.id`", "capture `.data.task.id`"},
		},
		{
			path: filepath.Join(root, "skills", "agent-team-gemini-harness", "references", "agent-team-runtime-protocol.md"),
			want: []string{"orchestrator-owned internal context", "capture `.data.run.id`", "capture `.data.task.id`"},
		},
		{
			path: filepath.Join(root, "skills", "agent-team-codex-harness", "references", "schemas", "agent-worker.template.toml"),
			want: []string{"orchestrator-supplied internal context", "Do not create them, infer them, or ask the user for them"},
		},
		{
			path: filepath.Join(root, "skills", "agent-team-gemini-harness", "references", "schemas", "agent-worker.template.md"),
			want: []string{"orchestrator-supplied internal context", "Do not create them, infer them, or ask the user for them"},
		},
	}
	for _, tc := range cases {
		body := readString(t, tc.path)
		for _, want := range tc.want {
			if !strings.Contains(body, want) {
				t.Fatalf("%s missing %q", tc.path, want)
			}
		}
	}
}

func TestAgentTeamHarnessSkillsRequireRuntimeStateProtocol(t *testing.T) {
	skillsDir := filepath.Join("..", "..", "skills")
	cases := []struct {
		name              string
		pointerFile       string
		agentDir          string
		skillDir          string
		activationVerb    string
		forbiddenPatterns []string
	}{
		{
			name:           "agent-team-codex-harness",
			pointerFile:    "AGENTS.md",
			agentDir:       ".codex/agents",
			skillDir:       ".agents/skills",
			activationVerb: "load",
			forbiddenPatterns: []string{
				"Generated harnesses default orchestrated runtime execution to harmony",
				"harmony state",
				"Generated harnesses default orchestrated runtime execution to direct `agent-team`",
				"Runtime execution uses `agent-team` directly",
				"route durable coordination through direct `agent-team` runtime execution",
				"If an `RUN_ID` or `TASK_ID` is provided, resume that run/task.",
			},
		},
		{
			name:           "agent-team-gemini-harness",
			pointerFile:    "GEMINI.md",
			agentDir:       ".gemini/agents",
			skillDir:       ".gemini/skills",
			activationVerb: "activate",
			forbiddenPatterns: []string{
				"Generated harnesses default orchestrated runtime execution to harmony",
				"harmony state",
				"Generated harnesses default orchestrated runtime execution to direct `agent-team`",
				"Runtime execution uses `agent-team` directly",
				"route durable coordination through direct `agent-team` runtime execution",
				"If an `RUN_ID` or `TASK_ID` is provided, resume that run/task.",
			},
		},
	}
	requiredRuntime := []string{
		"skill-first",
		"daemonless `agent-team`",
		"recipe/service/helper skills",
		"Agent Team runtime skills",
		"agent-team-run-create",
		"agent-team-task-create",
		"agent-team message",
		"agent-team inbox",
		"agent-team sync",
		"agent-team task complete",
		"agent-team run close",
		"_workspace/{run_id}/",
		"AGENT_TEAM_STATE_DIR",
		"RUN_ID` and `TASK_ID` are orchestrator-owned internal context",
		"not required user input",
		"capture `.data.run.id`",
		"capture `.data.task.id`",
		"Do not tell workers to \"just run the CLI\"",
		"Do not call `agent-team` while building, editing, or auditing a harness",
		"use Agent Team runtime skills for orchestrated execution after the harness is run",
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			skillPath := filepath.Join(skillsDir, tc.name, "SKILL.md")
			raw, err := os.ReadFile(skillPath)
			if err != nil {
				t.Fatalf("read %s: %v", skillPath, err)
			}
			body := string(raw)
			for _, want := range []string{tc.pointerFile, tc.agentDir, tc.skillDir, "agent-team-runtime-protocol.md"} {
				if !strings.Contains(body, want) {
					t.Fatalf("%s missing harness construction contract %q", skillPath, want)
				}
			}
			for _, forbidden := range tc.forbiddenPatterns {
				if strings.Contains(body, forbidden) {
					t.Fatalf("%s contains forbidden stale runtime text %q", skillPath, forbidden)
				}
			}

			protocolPath := filepath.Join(skillsDir, tc.name, "references", "agent-team-runtime-protocol.md")
			protocolRaw, err := os.ReadFile(protocolPath)
			if err != nil {
				t.Fatalf("read %s: %v", protocolPath, err)
			}
			protocol := string(protocolRaw)
			for _, want := range requiredRuntime {
				if !strings.Contains(body+"\n"+protocol, want) {
					t.Fatalf("%s runtime protocol missing %q", tc.name, want)
				}
			}
			if !strings.Contains(protocol, tc.activationVerb+" `agent-team-task-complete`") && !strings.Contains(protocol, strings.Title(tc.activationVerb)+" `agent-team-task-complete`") {
				t.Fatalf("%s runtime protocol does not use expected activation verb %q for command helper skills", tc.name, tc.activationVerb)
			}
		})
	}
}

func TestHarnessSkillDogfoodExamplesRecordClosedState(t *testing.T) {
	resultsDir := filepath.Join("..", "..", "examples", "dogfood-harness", "results")
	cases := []struct {
		dir             string
		prompt          string
		runID           string
		generatedFiles  []string
		requiredContent []string
	}{
		{
			dir:    "codex-harness-skill",
			prompt: "codex-harness-skill.txt",
			runID:  "run_harness_skill_codex",
			generatedFiles: []string{
				"generated/AGENTS.md",
				"generated/writer-agent.toml",
				"generated/reviewer-agent.toml",
				"generated/orchestrator-skill.md",
			},
			requiredContent: []string{
				"agent-team-codex-harness",
				"AGENTS.md",
				".codex/agents/writer.toml",
				".agents/skills/manual-dogfood-orchestrator/SKILL.md",
			},
		},
		{
			dir:    "gemini-harness-skill",
			prompt: "gemini-harness-skill.txt",
			runID:  "run_harness_skill_gemini",
			generatedFiles: []string{
				"generated/GEMINI.md",
				"generated/writer-agent.md",
				"generated/reviewer-agent.md",
				"generated/orchestrator-skill.md",
			},
			requiredContent: []string{
				"agent-team-gemini-harness",
				"GEMINI.md",
				".gemini/agents/writer.md",
				".gemini/skills/manual-dogfood-orchestrator/SKILL.md",
			},
		},
	}
	requiredEvents := []string{
		"run_created",
		"task_created",
		"message_sent",
		"message_acked",
		"task_started",
		"task_completed",
		"run_closed",
	}
	for _, tc := range cases {
		t.Run(tc.dir, func(t *testing.T) {
			base := filepath.Join(resultsDir, tc.dir)
			for _, rel := range append([]string{"README.md", "writer.md", "review.md", "run-summary.json", "event-log.json", "final-response.md"}, tc.generatedFiles...) {
				if _, err := os.Stat(filepath.Join(base, rel)); err != nil {
					t.Fatalf("missing %s: %v", filepath.Join(base, rel), err)
				}
			}

			readme := readString(t, filepath.Join(base, "README.md"))
			for _, want := range tc.requiredContent {
				if !strings.Contains(readme, want) {
					t.Fatalf("%s README missing %q", tc.dir, want)
				}
			}

			summary := readString(t, filepath.Join(base, "run-summary.json"))
			for _, want := range []string{
				`"prompt": "` + tc.prompt + `"`,
				`"run_id": "` + tc.runID + `"`,
				`"harness_files_generated": true`,
				`"pre_close_close_ready": true`,
				`"final_run_status": "closed"`,
				`"state_version": 10`,
				`"done": 2`,
			} {
				if !strings.Contains(summary, want) {
					t.Fatalf("%s run-summary missing %q", tc.dir, want)
				}
			}

			eventLog := readString(t, filepath.Join(base, "event-log.json"))
			for _, eventType := range requiredEvents {
				if !strings.Contains(eventLog, `"event_type": "`+eventType+`"`) {
					t.Fatalf("%s event-log missing event_type %q", tc.dir, eventType)
				}
			}
			if !strings.Contains(eventLog, `"state_version": 10`) {
				t.Fatalf("%s event-log missing final state_version 10", tc.dir)
			}
		})
	}
}

func readString(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(raw)
}

type skillFlagRow struct {
	Flag     string
	JSONKey  string
	Required string
	Default  string
}

func assertFlagTableMatchesSpec(t *testing.T, path, body string, subcommand SubcommandSpec) {
	t.Helper()
	rows := parseFlagRows(t, path, body)
	if len(subcommand.Flags) == 0 {
		if len(rows) != 1 || rows[0].Flag != "none" || rows[0].JSONKey != "none" || rows[0].Required != "no" || rows[0].Default != "none" {
			t.Fatalf("%s expected a single none flag row, got %#v", path, rows)
		}
		return
	}
	if len(rows) != len(subcommand.Flags) {
		t.Fatalf("%s has %d flag rows, want %d", path, len(rows), len(subcommand.Flags))
	}
	byName := map[string]skillFlagRow{}
	for _, row := range rows {
		if row.Flag == "none" {
			t.Fatalf("%s has none row for command with flags", path)
		}
		if byName[row.Flag].Flag != "" {
			t.Fatalf("%s has duplicate flag row %s", path, row.Flag)
		}
		byName[row.Flag] = row
	}
	required := map[string]bool{}
	for _, key := range subcommand.RequiredParams {
		required[key] = true
	}
	for _, flag := range subcommand.Flags {
		flagName := "--" + flag.Name
		row, ok := byName[flagName]
		if !ok {
			t.Fatalf("%s missing flag table row for %s", path, flagName)
		}
		if row.JSONKey != flag.JSONKey {
			t.Fatalf("%s %s json key = %q, want %q", path, flagName, row.JSONKey, flag.JSONKey)
		}
		wantRequired := "no"
		if required[flag.JSONKey] {
			wantRequired = "yes"
		}
		if row.Required != wantRequired {
			t.Fatalf("%s %s required = %q, want %q", path, flagName, row.Required, wantRequired)
		}
		wantDefault := expectedSkillDefault(flag, required[flag.JSONKey])
		if row.Default != wantDefault {
			t.Fatalf("%s %s default = %q, want %q", path, flagName, row.Default, wantDefault)
		}
	}
}

func parseFlagRows(t *testing.T, path, body string) []skillFlagRow {
	t.Helper()
	section := markdownSection(body, "## Flags")
	if section == "" {
		t.Fatalf("%s missing flags section", path)
	}
	rows := []skillFlagRow{}
	for _, line := range strings.Split(section, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "|") || !strings.HasSuffix(line, "|") {
			continue
		}
		if strings.Contains(line, "|------") || strings.Contains(line, "| Flag | JSON key |") {
			continue
		}
		cells := splitMarkdownRow(line)
		if len(cells) < 5 {
			t.Fatalf("%s malformed flags row: %s", path, line)
		}
		rows = append(rows, skillFlagRow{
			Flag:     normalizeSkillCell(cells[0]),
			JSONKey:  normalizeSkillCell(cells[1]),
			Required: normalizeSkillCell(cells[2]),
			Default:  normalizeSkillCell(cells[3]),
		})
	}
	if len(rows) == 0 {
		t.Fatalf("%s has no flag rows", path)
	}
	return rows
}

func splitMarkdownRow(line string) []string {
	trimmed := strings.Trim(line, "|")
	raw := strings.Split(trimmed, "|")
	cells := make([]string, 0, len(raw))
	for _, cell := range raw {
		cells = append(cells, strings.TrimSpace(cell))
	}
	return cells
}

func normalizeSkillCell(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, "`")
	return value
}

func markdownSection(body, heading string) string {
	start := strings.Index(body, heading)
	if start < 0 {
		return ""
	}
	rest := body[start+len(heading):]
	next := strings.Index(rest, "\n## ")
	if next >= 0 {
		rest = rest[:next]
	}
	return rest
}

func expectedSkillDefault(flag FlagSpec, required bool) string {
	if flag.Default != "" {
		return flag.Default
	}
	if required {
		return "-"
	}
	return "empty"
}

func commandHelp(commandName, subcommandName string) string {
	parts := []string{"agent-team", commandName}
	if subcommandName != "" {
		parts = append(parts, subcommandName)
	}
	parts = append(parts, "--help")
	return strings.Join(parts, " ")
}

func frontmatterValue(body, key string) string {
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, key+":") {
			value := strings.TrimSpace(strings.TrimPrefix(trimmed, key+":"))
			return strings.Trim(value, `"`)
		}
	}
	return ""
}
