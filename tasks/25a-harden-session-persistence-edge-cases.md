# Harden session persistence edge cases

Goal: Close Phase 3 persistence gaps before TUI/orchestrator code depends on session state.

Context:
- Phase 3 implemented the session model, append writer, resume loading, and tests.
- Review found several edge cases that weaken the intended crash-safety contract.

Do:
1. Track the next transcript sequence from loaded data as last valid message Seq + 1, not `len(Messages)`, so resumed appends continue correctly even if persisted sequence values are not zero-based or contiguous.
2. Change `Append` so in-memory `Messages` is updated only after the JSON line has been written and `Sync` has succeeded. A failed append must not leave memory ahead of disk.
3. Make transcript loading line-oriented for JSONL. Decode one line at a time, preserve valid messages in order, ignore only a corrupt/truncated trailing line, and return an error for corrupt non-trailing content.
4. Harden `SetAgentSession` atomic rewrite: handle and propagate temp-file close errors, clean up temp files on every failure path, and fsync the session directory after the rename where supported.
5. Consider fsyncing the session directory after first creating `transcript.jsonl` so a freshly created transcript is durable across crashes.
6. Add focused unit tests for resumed sequence assignment, append write-failure memory consistency if feasible, non-trailing corrupt transcript errors, trailing corrupt transcript tolerance, and agent file rewrite cleanup/durability behavior where portable.

Do NOT:
- Do not wire plan, TUI, or orchestrator behavior.
- Do not change the public on-disk file names (`transcript.jsonl`, `agents.json`).
- Do not add noisy logging from session load paths.

Done when:
- `go test ./internal/session/...` passes.
- `go test ./...` passes.
- `go vet ./...` passes.
- Session tests prove persisted and in-memory state cannot diverge on append errors.

Files: internal/session/session.go, internal/session/session_test.go
