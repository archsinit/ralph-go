# REVIEW: phase 6.5 plan chatroom readiness

Goal: Verify Phase 5.5 and Phase 6 are truly ready before Phase 7 plan integration begins.

Context:
- This is the review gate for the audit follow-up tasks. It should confirm more than compilation: lifecycle, cancellation, request contracts, `/plan` generation, artifact format, and resume behavior must work together.

Do:
1. Run:
   - `gofmt -l .` and confirm it prints nothing,
   - `go test -race ./internal/orchestrator/...`,
   - `go test ./internal/session/...`,
   - `go test ./internal/tui/...`,
   - `go test ./internal/orchestrator/...`,
   - `go test ./internal/plan/...`,
   - `go test ./...`,
   - `go vet ./...`,
   - `go build ./...`.
2. Inspect `internal/tui` and `internal/orchestrator/bridge.go` from another-package perspective and confirm live-handle operation, submit/cancel delivery, idle/second-press quit, shutdown safety, scroll preservation, and small-terminal/wide-rune behavior.
3. Inspect `internal/orchestrator/engine.go` and confirm deterministic turn order, exact prefix override scope, non-duplicated adapter requests, resume IDs, turn-cancel vs app-quit distinction, clear error handling, and no duplicate session/UI messages.
4. Inspect `/plan` generation and `internal/plan` artifacts and confirm the prompt schema, decoder, rendered task files, index parser, and write safety all use one compatible format.
5. Run an echo/fake-agent smoke test that:
   - starts plan mode with a temp session and outDir,
   - submits normal and prefixed messages,
   - cancels a streaming/blocked turn and continues chatting,
   - runs `/plan` with canned valid JSON and reparses the generated files,
   - runs `/plan` with invalid JSON and confirms no artifacts are written,
   - runs `/quit`, reopens the same session, and verifies resume without duplicates.
6. If real CLI smoke testing is available, run a short real design chat, generate a plan, inspect task quality, and document any skipped real-CLI checks with reasons.

Done when:
- Automated checks pass.
- Echo/fake plan mode demonstrates lifecycle, cancellation, `/plan`, `/quit`, artifact generation, and resume behavior.
- Generated artifacts are valid loop inputs.
- No known Phase 5.5 or Phase 6 blocker remains untracked before Phase 7.
