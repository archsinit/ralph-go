# Add turn/status indicator

Goal: Show whose turn it is and a waiting indicator while an agent streams.

Do:
1. Add a status line above the input: "waiting for <author>..." during a stream, "your turn" when input expected, idle otherwise.
2. Add a spinner (bubbles/spinner) during agent turns.

Done when:
- go build ./... succeeds.
- Status reflects current phase.
- Spinner animates during streaming only.

Files: internal/tui/tui.go
