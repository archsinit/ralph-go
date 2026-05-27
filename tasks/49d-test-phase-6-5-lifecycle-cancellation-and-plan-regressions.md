# TEST: phase 6.5 lifecycle, cancellation, and /plan regressions

Goal: Add deterministic tests that fail on the audited Phase 5.5/6 gaps and pass after tasks 49a-49c.

Context:
- Existing tests pass, but they do not prove TUI quit cancels a blocked engine, turn cancel returns to user input, adapter requests avoid transcript/new-message duplication, or `/plan` writes valid artifacts.
- `internal/plan` currently has no tests.

Do:
1. Add TUI tests for live-handle use, submit delivery, active-turn cancel, idle/second-press quit, wait/shutdown, scroll preservation, mouse scrolling, small terminal layout, and wide-rune/long-token wrapping.
2. Add bridge tests for context-aware submit/cancel delivery, stale cancel scoping, nil/shutdown handle safety, and no silent event drops under normal use.
3. Add engine tests with fake adapters and bounded contexts for:
   - round-robin order across multiple cycles,
   - prefix override applying to exactly the next agent turn,
   - non-duplicated `Transcript` and `NewMessage`, system prompt, and resume ID,
   - cancellation of streaming and blocked adapters returning to user input,
   - `/quit` clean exit,
   - adapter errors and parent context cancellation,
   - `/plan` valid/invalid generation paths.
4. Add plan package tests for strict decode, unknown fields, missing required fields, empty tasks, fence handling if supported, rendering, index write/read round-trip, and safe write failure behavior.
5. Add command-level coverage that `ralph-go --help` lists `plan` exactly once and that factored plan-mode lifecycle wiring handles TUI-first and engine-first exits.
6. Run and document:
   - `gofmt -l .` (must print nothing),
   - `go test -race ./internal/orchestrator/...`,
   - `go test ./internal/session/...`,
   - `go test ./internal/tui/...`,
   - `go test ./internal/orchestrator/...`,
   - `go test ./internal/plan/...`,
   - `go test ./...`,
   - `go vet ./...`,
   - `go build ./...`.

Do NOT:
- Do not use sleeps as the only proof of ordering or cancellation; use channels, contexts, and bounded waits.
- Do not rely on real LLM CLIs, network calls, or a real terminal in automated tests.
- Do not assert brittle ANSI output when a behavior-level assertion is possible.

Done when:
- The new tests fail against the audited implementation and pass after the fixes.
- No tests hang or leak goroutines under `go test -count=1 ./...`.
- All listed commands pass.

Files: internal/tui/tui_test.go, internal/tui/tui_external_test.go, internal/tui/tui_integration_test.go, internal/orchestrator/bridge_test.go, internal/orchestrator/engine_test.go, internal/orchestrator/generate_test.go, internal/plan/plan_test.go, cmd/plan.go, cmd/root.go
