# Align real CLI adapter arguments with installed CLIs

Goal: Make the Phase 2 CLI argument builders match the real CLI surfaces available on this machine.

Context:
- `claude --help` confirms `--print`, `--output-format stream-json`, `--verbose`, `--system-prompt`, and `--resume`.
- Claude also exposes `--include-partial-messages`, which is needed for true partial-message streaming with stream-json.
- `codex --help` shows non-interactive execution is `codex exec`.
- `codex exec --help` does not expose `--quiet` or `--system-prompt`; the current codex args include unsupported flags.
- `gemini` and `opencode` were not installed during review, so keep those adapters conservative until their flags can be verified.

Do:
1. Update codex argument construction to invoke `codex exec` with valid flags; remove unsupported `--quiet` and `--system-prompt`.
2. Check `codex exec resume --help` and either implement resume argument construction or explicitly document why SupportsResume remains false.
3. If codex exposes session IDs in JSON/event output, parse them into Result.SessionID; otherwise document the limitation in code comments.
4. Update codex tests to assert valid `exec`/resume shape and absence of unsupported flags.
5. Update claude stream-json args to include `--include-partial-messages` if needed for real streaming, and ensure parsing avoids duplicate final text.
6. Verify or clearly mark gemini/opencode flag assumptions; they should still use replay fallback for transcript context.
7. Update the config known-CLI validation message to include every accepted CLI, including `echo`, preferably generated from KnownCLIs.

Do NOT:
- Do not add broad CLI integration tests that require live auth or network access.
- Do not remove unverified adapters; keep them conservative and documented.

Done when:
- `claude --help`, `codex --help`, and `codex exec --help` assumptions are reflected in tests or comments.
- `go test ./internal/agent/... ./internal/config/...` passes.
- `go test ./...` passes.
- `go vet ./...` passes.
- `go build ./...` passes.

Files: internal/agent/claude.go, internal/agent/codex.go, internal/agent/gemini.go, internal/agent/opencode.go, internal/agent/agent_test.go, internal/config/config.go, internal/config/config_test.go
