# Build TUI model skeleton

Goal: Create the Bubble Tea model with a transcript viewport and an input box.

Context:
- bubbles components: viewport, textinput.
- Keep orchestrator wiring out of this task; expose hooks via messages/callbacks.

Do:
1. In internal/tui/tui.go define model with: viewport (bubbles/viewport) for transcript, textinput (bubbles/textinput) for input, width/height, messages []renderedMsg.
2. Implement Init, Update (handle tea.WindowSizeMsg, tea.KeyMsg for submit and quit), View (viewport on top, input below).
3. Add NewModel(...) and a Run(...) entrypoint that starts tea.Program.

Do NOT:
- Do not connect agents here.

Done when:
- go build ./... succeeds.
- Running a temporary harness shows transcript area + input, accepts typing, quits on ctrl-c.
- Resize adjusts layout.

Files: internal/tui/tui.go
