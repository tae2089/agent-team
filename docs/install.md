# Install

## Go Install

```bash
go install github.com/tae2089/agent-team/cmd/agent-team@latest
agent-team version
agent-team --help
```

From a local checkout:

```bash
go install ./cmd/agent-team
```

## Release Binary

Download the binary for your OS and architecture from the GitHub release page.

Artifact names use this shape:

```text
agent-team_vX.Y.Z_{goos}_{goarch}
agent-team_vX.Y.Z_windows_{goarch}.exe
SHA256SUMS
```

Verify checksums on Linux:

```bash
sha256sum -c SHA256SUMS
```

Verify checksums on macOS:

```bash
shasum -a 256 -c SHA256SUMS
```

Install the binary somewhere on `PATH`:

```bash
chmod +x agent-team_vX.Y.Z_darwin_arm64
mkdir -p "$HOME/bin"
mv agent-team_vX.Y.Z_darwin_arm64 "$HOME/bin/agent-team"
agent-team version
```

## Source Build

```bash
go build -o agent-team ./cmd/agent-team
./agent-team version
```

## First Check

```bash
agent-team init
agent-team run create --title "install check"
agent-team run list --status open
```

By default, state is created under `.agent-team/agent-team.db`. Set `AGENT_TEAM_STATE_DIR` to isolate workflows.

If a command fails, inspect `error.code` and `error.recovery` in the JSON response. Recovery guidance is documented in [errors.md](errors.md). Codex/Gemini sandbox and permission guidance is documented in [harness-sandbox-permissions.md](harness-sandbox-permissions.md).
