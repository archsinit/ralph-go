# Implement cancel current turn

Goal: Allow ctrl-c to cancel an in-progress agent turn without exiting the app.

Context:
- Engine owns the per-turn context; TUI sends the cancel signal.

Do:
1. Thread a context per agent turn; on a cancel signal from the UI, cancel that turn's context.
2. On cancel: stop streaming, append a short "(turn cancelled)" note to transcript, return control to the next appropriate turn (user).
3. Distinguish app-quit (second ctrl-c or a quit key) from turn-cancel.

Done when:
- go build ./... succeeds.
- Cancel stops the stream and keeps the app running.
- Transcript remains consistent after cancel.

Files: internal/orchestrator/engine.go, internal/tui/tui.go
