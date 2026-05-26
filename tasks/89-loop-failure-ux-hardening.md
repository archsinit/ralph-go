# Loop failure UX hardening

Goal: Ensure missing CLI, auth fail, and timeout produce clear TUI + ntfy + log output.

Do:
1. Force each failure: missing executor binary, bad auth, task timeout.
2. Confirm a clear TUI status, an ntfy alert, and a master/task log entry for each, with a clean halt.

Done when:
- go build ./... succeeds.
- Each failure mode is clearly reported across TUI, ntfy, and logs.
- Loop halts cleanly.

Files: internal/loop/engine.go, internal/tui/loop.go
