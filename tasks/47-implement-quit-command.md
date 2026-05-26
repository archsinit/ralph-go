# Implement /quit command

Goal: Add /quit to save the session and exit the chatroom cleanly.

Do:
1. Detect '/quit' (and ctrl-c app-quit path) -> persist session, stop the tea Program, exit 0.
2. Ensure any in-progress turn is cancelled first.

Done when:
- go build ./... succeeds.
- /quit exits cleanly with session saved.
- No goroutine left running.

Files: internal/tui/tui.go, internal/orchestrator/engine.go
