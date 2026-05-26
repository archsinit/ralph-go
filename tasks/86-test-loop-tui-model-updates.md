# TEST: loop TUI model updates

Goal: Unit-test loop model state transitions.

Context:
- Drive Update with synthetic msgs; no real Program.

Do:
1. Create internal/tui/loop_test.go.
2. Test: task-started msg highlights the right task; output tokens accumulate; checkbox-flipped msg marks done; status msg updates the status line; resize recomputes layout.

Done when:
- go test ./internal/tui/... passes (loop parts).
- Transitions asserted.

Files: internal/tui/loop_test.go
