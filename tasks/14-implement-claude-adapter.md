# Implement claude adapter

Goal: Adapter that invokes the claude CLI non-interactively with session resume when possible.

Context:
- Exact claude CLI flags must be verified at build time; isolate flag-building in a pure function buildArgs(req) []string for testing.
- If a flag is unknown, prefer a documented placeholder and a TODO, keeping buildArgs pure.

Do:
1. In internal/agent/claude.go define claudeAdapter implementing Adapter.
2. Capabilities: SupportsResume true.
3. Build args for non-interactive one-shot (print mode) including system prompt and the new message; include resume flag when ResumeSessionID set.
4. If resume unsupported in a given invocation, fall back to sending transcript as context (see replay helper task).
5. Parse the returned session id from output/metadata and set Result.SessionID.
6. Map ctx + Request to runStreaming.

Done when:
- go build ./... succeeds.
- Adapter compiles and Name() returns "claude".
- Flag construction unit-testable (see test task).

Files: internal/agent/claude.go
