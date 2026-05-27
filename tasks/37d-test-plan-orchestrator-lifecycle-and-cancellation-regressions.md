# TEST: plan orchestrator lifecycle and cancellation regressions

Goal: Add regression coverage that fails on the current Phase 4.5/5 gaps and proves the follow-up fixes are complete.

Context:
- Current automated checks pass, but they do not catch that `cmd/plan.go` starts the engine only after the TUI exits.
- Engine tests mostly smoke-run goroutines and log observations instead of asserting order, prefix routing, cancellation, adapter request contents, and clean shutdown.
- TUI public API tests compile exported types, but do not prove a live handle is usable while the Bubble Tea program is running.
- `gofmt -l .` currently reports unformatted orchestrator files.

Do:
1. Add TUI lifecycle tests using Bubble Tea test options where possible (for example without the renderer/input) to prove external code can start a live TUI, obtain a handle immediately, send finalized messages, stream tokens, send cancel, quit, and wait for completion.
2. Add bridge tests for submit delivery, cancel delivery, context cancellation, startup/shutdown handle safety, and no silent cancel drops under normal use.
3. Replace weak engine smoke tests with deterministic assertions:
   - turn order over several agent/user cycles,
   - prefix override routes only the immediately following agent turn,
   - adapter `Request.NewMessage`, `Transcript`, system prompt, and resume ID are correct,
   - agent cancellation cancels the adapter context and returns to user,
   - agent responses and cancellation notes are persisted/displayed exactly once,
   - unknown agents and empty turn orders return clear errors.
4. Add command-level regression coverage for duplicate command registration. A help-rendering test is sufficient if a full subprocess test is unnecessary.
5. Add a plan-mode wiring smoke test using echo/fake adapters and a temp session dir. Prefer testing factored wiring directly over fragile real-terminal automation.
6. Run and document these commands in the task output:
   - `gofmt -l .` (must print nothing),
   - `go test ./internal/tui/...`,
   - `go test ./internal/orchestrator/...`,
   - `go test ./...`,
   - `go vet ./...`,
   - `go build ./...`.

Do NOT:
- Do not rely on real network-backed CLIs in automated tests.
- Do not use arbitrary sleeps as the only proof of ordering or cancellation; use channels, contexts, and bounded timeouts.
- Do not assert brittle ANSI escape sequences unless no better behavior-level assertion is available.

Done when:
- The tests fail against the pre-fix behavior described above and pass after tasks 37a–37c are implemented.
- No test goroutine leaks or hangs under `go test -count=1 ./...`.
- All listed commands pass, with no files reported by `gofmt -l .`.

Files: internal/tui/tui_external_test.go, internal/tui/tui_test.go, internal/orchestrator/bridge.go, internal/orchestrator/engine_test.go, cmd/plan.go, cmd/root.go
