# Stream tokens into in-progress message

Goal: Append streamed tokens live to the current agent message bubble.

Context:
- Orchestrator will translate agent.Token channel into these msgs.

Do:
1. Define tea.Msg types: msgStreamStart{author}, msgStreamToken{text}, msgStreamEnd.
2. On start, create an in-progress message; on token, append and refresh viewport; on end, finalize.
3. Provide a method to feed these from outside (channel -> tea via program.Send).

Done when:
- go build ./... succeeds.
- Tokens visibly accumulate in the live message during a stream.
- Finalized message persists after end.

Files: internal/tui/tui.go
