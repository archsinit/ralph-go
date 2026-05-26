# Fix adapter stream contract and replay fallback

Goal: Close the Phase 2 adapter lifecycle and transcript replay gaps before session/orchestrator work depends on them.

Context:
- The documented Adapter contract says the Token channel closes first, then exactly one Result is delivered.
- Current implementations send Result while the Token channel is still open because of defer ordering and explicit result sends.
- CLI adapter cancellation paths can return without sending an explicit error Result.
- Adapters without resume support, and resume-capable adapters called without a ResumeSessionID, must not drop transcript context.

Do:
1. Update runStreaming, CLI adapter wrappers, and echo adapter so every Invoke path closes the Token channel before delivering exactly one Result, then closes the Result channel.
2. Ensure cancellation returns a terminal Result with ctx.Err() or a wrapped cancellation error; do not allow callers to observe a closed, empty result channel or nil error on cancellation.
3. Surface scanner errors from stdout as Result.Err, while preserving stderr-tail errors for failed processes.
4. Use renderReplay consistently whenever an adapter lacks resume support or a request has transcript context without a ResumeSessionID.
5. Ensure claude does not drop Request.Transcript when ResumeSessionID is empty.
6. Keep echo deterministic and make concatenated tokens reconstruct NewMessage exactly.
7. Add focused tests for channel ordering, exactly-one terminal Result, cancellation errors, scanner/process errors, and replay use in no-resume paths.

Do NOT:
- Do not implement session persistence.
- Do not wire plan or loop orchestration.
- Do not require real network-backed CLI calls in unit tests.

Done when:
- `go test ./internal/agent/...` passes.
- `go test ./...` passes.
- `go vet ./...` passes.
- `go build ./...` passes.
- Tests prove the Token channel closes before Result delivery on success and cancellation.

Files: internal/agent/exec.go, internal/agent/claude.go, internal/agent/codex.go, internal/agent/gemini.go, internal/agent/opencode.go, internal/agent/echo.go, internal/agent/agent_test.go
