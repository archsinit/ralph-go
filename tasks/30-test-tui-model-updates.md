# TEST: TUI model updates

Goal: Unit-test model state transitions without a real terminal.

Context:
- Drive Update directly with synthetic msgs; do not start a real Program.

Do:
1. Create internal/tui/tui_test.go.
2. Test: send tea.WindowSizeMsg -> width/height updated, layout recomputed.
3. Test: stream start/token/token/end -> one finalized message with concatenated text.
4. Test: submitting input via key Enter -> produces the expected outbound submit (capture via the model's exposed channel/callback).

Done when:
- go test ./internal/tui/... passes.
- Stream concatenation asserted.
- Submit path asserted.

Files: internal/tui/tui_test.go
