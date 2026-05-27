# REVIEW: phase 5.5 plan orchestrator readiness

Goal: Verify the TUI/session/orchestrator stack is truly ready for the plan-format and checklist-generation phases.

Context:
- This review is a gate after the Phase 5 follow-up fixes. It should confirm the app is not merely compiling, but that the interactive plan loop can actually run while the TUI is active.

Do:
1. Run:
   - `gofmt -l .` and confirm it prints nothing,
   - `go test ./internal/session/...`,
   - `go test ./internal/tui/...`,
   - `go test ./internal/orchestrator/...`,
   - `go test ./...`,
   - `go vet ./...`,
   - `go build ./...`.
2. Inspect `internal/tui` and `internal/orchestrator/bridge.go` from another-package perspective and confirm external orchestration code can start a live TUI, retain a handle, send messages/streams, receive submits, receive cancels, wait for shutdown, and avoid nil/race panics.
3. Inspect `internal/orchestrator/engine.go` and confirm:
   - turn order is deterministic,
   - prefix overrides affect exactly the next agent turn,
   - agent requests contain useful `NewMessage`/`Transcript`/resume information,
   - cancellation during agent streaming cancels the turn without exiting the app,
   - session transcript and TUI display receive each finalized message exactly once.
4. Inspect `cmd/plan.go` and `cmd/root.go` and confirm:
   - `plan` is registered once,
   - TUI and engine lifecycles run concurrently,
   - TUI quit cancels the engine,
   - engine errors shut down the UI and return a useful CLI error,
   - resume shows existing transcript without duplication.
5. Confirm the Phase 5 turn-order model no longer blocks Phase 5.6 from defining the task-spec `internal/plan.Plan` type.
6. Run a manual smoke test with a temp echo-agent config:
   - start `ralph-go plan`,
   - submit several messages,
   - use an agent prefix override,
   - cancel one streaming/blocked turn,
   - quit cleanly,
   - reopen the same `--session` directory and verify transcript resume.
7. Document any real-CLI smoke checks skipped and why.

Done when:
- All automated checks pass.
- Manual echo-agent plan mode runs without startup deadlock, duplicate messages, cancel failure, or dirty shutdown.
- No known Phase 4.5/5 blocker remains untracked before Phase 5.6 and Phase 6 begin.
