# Harden remaining session sequence and durability edge cases

Goal: Close the remaining Phase 3.5 persistence gaps found during the Phase 4 review pass.

Context:
- The current session package passes the existing tests and implements most Phase 3.5 hardening.
- A few details still do not exactly match the requested crash-safety/resume contract.

Findings to address:
- `loadMessages` sets `nextSeq` to max loaded `Seq` + 1. Task 25a requested last valid message `Seq` + 1. These differ if a transcript contains non-monotonic persisted sequence values.
- `SetAgentSession` mutates `Session.AgentSessions` before the atomic rewrite succeeds. If temp creation/write/sync/close/rename fails, memory can get ahead of disk.
- `Append` ignores deferred `Close` errors after a successful JSON encode and `Sync`; memory is committed before a possible close failure is observed.
- `bufio.Scanner` keeps the default token limit, so a large but valid JSONL message can fail to load with `ErrTooLong`.
- The existing append memory-consistency test covers a successful append only; it does not simulate a failed write/sync/close path.

Do:
1. Decide and implement the sequence contract from task 25a: resumed appends should use the last valid transcript message's `Seq + 1` unless the plan is explicitly amended to require max `Seq + 1` instead.
2. Add a regression test with non-monotonic sequence values to prove the chosen next-sequence behavior.
3. In `SetAgentSession`, prepare a copy of the map and commit it to `Session.AgentSessions` only after successful temp write, temp sync, temp close, rename, and best-effort directory fsync; or otherwise roll back the in-memory mutation on every error path.
4. In `Append`, handle and propagate close errors where possible. Do not append to memory or increment `nextSeq` until all required persistence operations for that append have succeeded.
5. Replace `bufio.Scanner` or increase/configure its buffer so large valid transcript messages load successfully, or document and test a deliberate maximum with a clear error.
6. Add failure-path tests for append memory consistency and agent-session memory consistency using portable techniques where feasible.
7. Preserve the existing JSONL corruption contract: valid messages are preserved in order, only corrupt/truncated trailing content is tolerated, and corrupt non-trailing content returns an error.

Do NOT:
- Do not rename `transcript.jsonl` or `agents.json`.
- Do not add noisy logging from session load paths.
- Do not wire sessions into the TUI or orchestrator here.

Done when:
- The sequence behavior exactly matches the documented contract and has regression coverage.
- Failed appends and failed agent-session rewrites cannot silently leave in-memory state ahead of disk.
- Large valid transcript entries are handled according to the tested contract.
- `go test ./internal/session/...` passes.
- `go test ./...`, `go vet ./...`, and `go build ./...` pass.

Files: internal/session/session.go, internal/session/session_test.go
