# Codex Manual Dogfood Result

This directory stores a sanitized example result from running `examples/dogfood-harness/prompts/codex.txt` with `codex exec`.

Captured with:

```text
codex-cli 0.130.0
```

Expected files after a successful dogfood:

- `writer.md`
- `review.md`
- `run-summary.json`
- `event-log.json`
- `final-response.md`

The JSON files intentionally omit timestamps and absolute temp paths while preserving run status, task counts, event types, and state versions.
