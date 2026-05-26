# Show progress, retries, timing, last event

Goal: Surface retry count, elapsed-vs-timeout, and the last ntfy event in a status area.

Do:
1. Add a status line: 'task i/N', current attempt/max, elapsed vs task_timeout, last ntfy event text, and overall state (running/halted/done + reason).
2. Update a timer each second during a task.

Done when:
- go build ./... succeeds.
- Status reflects retries, timing, last event, and state.
- Halt reason shown on stop.

Files: internal/tui/loop.go
