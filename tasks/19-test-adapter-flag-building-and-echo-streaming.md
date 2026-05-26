# TEST: adapter flag building and echo streaming

Goal: Unit-test pure buildArgs functions and end-to-end echo streaming.

Do:
1. Create internal/agent/agent_test.go.
2. Test buildArgs for claude and codex: given a Request, assert expected flags/positional args present (resume flag present only when ResumeSessionID set).
3. Test renderReplay output contains system prompt, prior messages, new message in order.
4. Test echo adapter: Invoke yields tokens reconstructing NewMessage and a terminal Result with no error and a SessionID.
5. Use context.Background and a short timeout.

Done when:
- go test ./internal/agent/... passes.
- Resume flag presence asserted both ways.
- Echo round-trip verified.

Files: internal/agent/agent_test.go
