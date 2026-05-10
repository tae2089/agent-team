# Dogfood Basic

This example runs a complete daemonless workflow:

1. initialize state
2. create a run
3. create a worker task
4. simulate worker start, sync, completion
5. inspect summary
6. close the run

Run from the repository root:

```bash
sh examples/dogfood-basic/run.sh
```

The script uses a temporary `AGENT_TEAM_STATE_DIR` and writes artifacts under a temporary directory.

