# REVIEW: phase 3.5 session durability readiness

Goal: Verify session persistence is ready for Phase 4+ consumers after the Phase 3.5 hardening work.

Context:
- Phase 3.5 should address sequence continuity, append failure consistency, JSONL corruption handling, and atomic agent-session rewrites.

Do:
1. Run `go test ./internal/session/...`, `go test ./...`, `go vet ./...`, and `go build ./...`.
2. Manually inspect `internal/session/session.go` to confirm `Append` commits to memory only after durable disk write succeeds.
3. Manually inspect transcript loading to confirm only a corrupt final JSONL line is ignored; corrupt middle lines must produce an error.
4. Manually inspect `SetAgentSession` to confirm temp file cleanup, close error handling, atomic rename, and best-effort directory fsync are implemented.
5. Simulate resume from a transcript whose final valid message has a non-`len(Messages)-1` sequence and confirm the next append uses last Seq + 1.

Done when:
- Tests, vet, and build pass.
- Reviewer is satisfied that a crash or failed disk write cannot silently make session memory, transcript, or agent-session metadata inconsistent.
- Any remaining durability limitations are documented in this task or in a new follow-up before Phase 4 starts.
