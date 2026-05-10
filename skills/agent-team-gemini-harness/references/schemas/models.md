# Gemini Model Guidance

Use Gemini CLI-supported model ids verified by the local Gemini configuration or current Gemini CLI documentation.

Do not invent model ids. If the project already standardizes model ids in Gemini settings, follow that configuration.

## Recommended Tiers

| Role | Tier |
| --- | --- |
| Researcher / analyst | gemini-3.1-flash-preview or gemini-3.1-pro-preview |
| Writer / planner | gemini-3.1-flash-preview or gemini-3.1-pro-preview |
| Coder / reviewer / QA | gemini-3.1-pro-preview when available |
| Architect / orchestrator | gemini-3.1-pro-preview |

## Tool Policy Reminder

Model choice and tool choice are separate:

- give workers only the tools they need
- ordinary workers do not include `invoke_agent`
- runtime workers that call `agent-team` need `run_shell_command`
- orchestrators that invoke specialists need `invoke_agent`
- never use wildcard tool permissions

## Update Protocol

When a new model is released, update only this file:

1. Modify model IDs in the table above (add a new model row or replace an existing ID).
2. Apply the same change to the `model` field in `references/schemas/agent-worker.template.md`.
3. All agents subsequently created by the harness will use the new ID.
4. **Existing generated agents require manual updates** (edit `.gemini/agents/*.md` directly).
