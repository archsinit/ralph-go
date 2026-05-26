# Implement gemini and opencode adapters

Goal: Best-effort adapters for gemini and opencode CLIs following the same pattern.

Context:
- These are lower priority; correctness of pattern matters more than exhaustive flag coverage.

Do:
1. Create internal/agent/gemini.go and internal/agent/opencode.go.
2. Each implements Adapter with a pure buildArgs and runStreaming wiring.
3. Set SupportsResume conservatively (false unless verified).
4. Name() returns "gemini" / "opencode".

Done when:
- go build ./... succeeds.
- Both adapters compile and expose pure buildArgs.

Files: internal/agent/gemini.go, internal/agent/opencode.go
