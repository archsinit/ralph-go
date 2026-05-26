# Implement codex adapter

Goal: Adapter that invokes the codex CLI non-interactively with session resume when possible.

Context:
- Verify codex exec/non-interactive flags at build; keep buildArgs pure.

Do:
1. In internal/agent/codex.go define codexAdapter implementing Adapter.
2. Mirror the claude adapter structure: pure buildArgs(req) []string, Capabilities, Invoke via runStreaming.
3. Set SupportsResume per codex CLI reality; if unknown, false and rely on replay.
4. Parse session id if the CLI exposes one.

Done when:
- go build ./... succeeds.
- Name() returns "codex".
- buildArgs is pure and testable.

Files: internal/agent/codex.go
