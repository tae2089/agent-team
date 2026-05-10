# Changelog

All notable changes to `agent-team` are tracked here.

The project is pre-v1. Until v1, CLI behavior is intended to be stable within the documented command contract, but the SQLite DB schema may still change without migrations.

## Unreleased

### Added

- Public alpha README flow with copy-paste quickstart.
- Install, release, and checksum documentation.
- CI dogfood and manual Codex/Gemini harness dogfood guides.
- Schema export to skill documentation consistency tests.
- Stress profiles for CI smoke and local heavy runs.
- Structured `error.recovery` hints in failure JSON and schema export.
- Error recovery and Codex/Gemini sandbox permission documentation.

### Changed

- Release workflow now uploads OS/architecture binaries and `SHA256SUMS`.
- Skill docs default values are aligned with the command schema contract.
- Heavy stress is reserved for local validation; CI uses a small smoke profile.
- Harness guidance now documents skill-first runtime execution, recovery handling, and Codex/Gemini permission boundaries.

### Tests

- Added status policy regression coverage.
- Added dogfood scripts for basic and harness-style workflows.
- Added stress test coverage for concurrent task and message mutations.
