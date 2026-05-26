# Render messages with per-author styling

Goal: Display each message with an author label and a distinct color per agent.

Context:
- Wrapping must respect current width from WindowSizeMsg.

Do:
1. Add lipgloss styles; assign a color per author name (stable hash or config order).
2. Render label + wrapped text into the viewport; user, claude, codex visually distinct.
3. Auto-scroll viewport to bottom on new message unless user scrolled up.

Done when:
- go build ./... succeeds.
- Authors visually distinguishable.
- Long messages wrap to width; resize re-wraps.

Files: internal/tui/tui.go
