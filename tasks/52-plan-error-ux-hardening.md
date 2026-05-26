# Plan error UX hardening

Goal: Make CLI-missing, auth-fail, and timeout errors clear in the TUI.

Do:
1. Detect adapter Invoke errors (missing binary, auth) and show a clear status line, keep app alive.
2. Add a sane per-turn timeout for plan agents (config or constant) with a clear message on exceed.
3. Ensure /plan generation errors are distinguishable from turn errors.

Done when:
- go build ./... succeeds.
- Each failure mode yields a clear, distinct message.
- App stays responsive after a failed turn.

Files: internal/orchestrator/engine.go, internal/tui/tui.go
