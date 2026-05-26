# Build loop TUI two-pane model

Goal: Create the loop TUI: a checklist pane and a streaming output pane.

Context:
- Mirror plan TUI patterns; reuse styles.

Do:
1. In internal/tui/loop.go define a loop model: left/top pane lists tasks with live checkbox state, current task highlighted; bottom/right pane streams current executor/reviewer output.
2. Implement Init/Update/View; handle resize.
3. Expose msg types for: task started, output token, checkbox flipped, status changed, stopped.

Done when:
- go build ./... succeeds.
- Both panes render; current task highlighted.
- Resize behaves.

Files: internal/tui/loop.go
