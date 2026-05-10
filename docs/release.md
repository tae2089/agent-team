# Release Policy

`agent-team` is currently pre-v1. Releases are intended for alpha users who can tolerate state directory resets when the DB schema changes.

## Version Sources

- `internal/version/version.go` contains the development default, for example `1.2.0-dev`.
- Release builds inject the tag version with Go linker flags:

```bash
go build -trimpath \
  -ldflags "-s -w -X github.com/tae2089/agent-team/internal/version.Version=${VERSION}" \
  -o agent-team ./cmd/agent-team
```

`agent-team version` must report the injected release version.

## Tag Shape

Stable release tags:

```text
vMAJOR.MINOR.PATCH
```

Pre-release tags:

```text
vMAJOR.MINOR.PATCH-alpha.N
vMAJOR.MINOR.PATCH-beta.N
vMAJOR.MINOR.PATCH-rc.N
```

## Artifacts

The release workflow builds:

```text
agent-team_vX.Y.Z_linux_amd64
agent-team_vX.Y.Z_linux_arm64
agent-team_vX.Y.Z_darwin_amd64
agent-team_vX.Y.Z_darwin_arm64
agent-team_vX.Y.Z_windows_amd64.exe
agent-team_vX.Y.Z_windows_arm64.exe
SHA256SUMS
```

## Checklist

Before tagging:

```bash
go test ./...
sh examples/dogfood-basic/run.sh
sh examples/dogfood-harness/run.sh
sh examples/stress/ci-smoke.sh
AGENT_TEAM_STRESS_PROFILE=heavy AGENT_TEAM_STRESS_N=20 AGENT_TEAM_STRESS_RUNS=2 sh examples/stress/concurrency.sh
```

Then:

1. Update `CHANGELOG.md`.
2. Confirm `internal/version/version.go` has the next `-dev` default if needed.
3. Create and push a `v*` tag.
4. Confirm GitHub Release artifacts and `SHA256SUMS` are present.
5. Download one binary, verify checksum, and confirm `agent-team version`.

## DB Stability

Until v1, the SQLite schema is not a compatibility guarantee. Breaking DB changes may require deleting `.agent-team/` or a custom `AGENT_TEAM_STATE_DIR` and running `agent-team init` again.
