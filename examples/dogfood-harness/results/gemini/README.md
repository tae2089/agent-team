# Gemini Manual Dogfood Result

This directory stores a sanitized example result from running `examples/dogfood-harness/prompts/gemini.txt` with `gemini -p`.

Captured with:

```text
gemini 0.41.2
```

Expected files after a successful dogfood:

- `writer.md`
- `review.md`
- `run-summary.json`
- `event-log.json`
- `final-response.md`

The JSON files intentionally omit timestamps and absolute temp paths while preserving run status, task counts, event types, and state versions.
