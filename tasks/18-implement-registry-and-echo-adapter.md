# Implement registry and echo adapter

Goal: Add an adapter registry keyed by CLI name plus a network-free echo adapter for tests.

Context:
- echo enables full plan/loop tests without real CLIs or network.
- Add "echo" to the known-CLI set used by config validation.

Do:
1. In internal/agent/registry.go add func New(cli string, a config.Agent) (Adapter, error) mapping cli -> adapter (claude/codex/gemini/opencode/echo).
2. Add internal/agent/echo.go: echoAdapter streams back the NewMessage split into a few tokens, sets a fake SessionID, no process spawn.
3. Registry returns error for unknown cli.

Done when:
- go build ./... succeeds.
- New("echo", ...) returns the echo adapter.
- Unknown cli returns clear error.

Files: internal/agent/registry.go, internal/agent/echo.go
