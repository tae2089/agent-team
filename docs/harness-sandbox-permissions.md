# Harness Sandbox And Permissions

Agent Team harnesses are skill-first. Codex and Gemini agents should load or activate Agent Team runtime skills before running daemonless `agent-team` commands.

The CLI coordinates state only. It does not run workers, spawn LLMs, or bypass the active tool sandbox.

## Codex

Codex harnesses use `.codex/agents/{name}.toml`.

Recommended `sandbox_mode` values:

| Role | Default | Use when |
| --- | --- | --- |
| researcher / analyst | `read-only` | The role only reads repo files and writes no artifacts. |
| planner / docs writer | `workspace-write` | The role writes `_workspace/` artifacts or project docs. |
| coder / reviewer / QA | `workspace-write` | The role edits files, runs tests, or writes review artifacts. |
| operator / release | `danger-full-access` | Only when explicitly required for external paths, release tooling, or trusted local automation. |

Codex headless runs may restrict hidden directories such as `.codex/` or `.agents/` depending on sandbox mode. If a harness-generation dogfood must create those directories in a temporary workspace, run it in a trusted temp directory with an explicit broader sandbox and keep `AGENT_TEAM_STATE_DIR` inside that temp workspace.

## Gemini

Gemini harnesses use `.gemini/agents/{name}.md` and explicit `tools`.

Recommended tools:

| Role | Tools |
| --- | --- |
| direct-only researcher / analyst | `ask_user`, `activate_skill` |
| planner / document producer | `ask_user`, `activate_skill`, plus required file tools |
| coder / docs writer | `ask_user`, `activate_skill`, required file tools, `run_shell_command` only for commands/tests |
| reviewer / QA | `ask_user`, `activate_skill`, file tools, `run_shell_command` when validation commands are needed |
| runtime worker | `ask_user`, `activate_skill`, `run_shell_command`, and required file tools |
| orchestrator / supervisor | `ask_user`, `activate_skill`, `invoke_agent`, plus `run_shell_command` only for durable runtime commands |

Do not use wildcard tool permissions such as `tools: ["*"]`. Ordinary workers should not include `invoke_agent`.

After changing Gemini agents or skills, reload them inside Gemini CLI:

```text
/agents reload
/skills reload
```

## State Directory

Use a state directory that the active sandbox can read and write:

```bash
export AGENT_TEAM_STATE_DIR="$(mktemp -d)"
agent-team init
```

For project-local runs, the default is:

```text
.agent-team/agent-team.db
```

When running in Codex or Gemini, prefer a temp or project-local state directory and avoid pointing `AGENT_TEAM_STATE_DIR` at a home directory path unless the tool sandbox explicitly allows it.
