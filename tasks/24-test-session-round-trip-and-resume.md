# TEST: session round-trip and resume

Goal: Unit-test persistence, reload equality, and resume id recovery.

Do:
1. Create internal/session/session_test.go using t.TempDir().
2. Test: Open empty, Append three messages, reopen -> messages equal in order and content.
3. Test: SetAgentSession then reopen -> ResumeID returns stored value.
4. Test: append a deliberately corrupt trailing line to transcript.jsonl, reopen -> prior messages intact, no panic.

Done when:
- go test ./internal/session/... passes.
- Round-trip equality asserted.
- Corrupt-line tolerance asserted.

Files: internal/session/session_test.go
