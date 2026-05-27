# Fix engine turn semantics, cancellation, and adapter request contract

Goal: Make the Phase 5 turn engine satisfy the round-robin, prefix-routing, cancellation, transcript, and adapter-input contracts that later plan/loop features will rely on.

Context:
- The engine only watches the cancel channel on user turns; cancellation during an agent turn is effectively unwired.
- `runAgentTurn` can wait indefinitely while draining tokens after context cancellation if an adapter misbehaves, and it does not reliably clear the UI stream before returning a cancellation result.
- Cancellation appends a note after the in-progress stream remains active, which can leave the UI in an inconsistent streaming state.
- `agent.Request.NewMessage` is never populated. For resume-capable CLIs this sends an empty prompt on resumed turns; for replay fallback it can render a transcript ending in a blank new-message section.
- The current streamed-response display path can duplicate agent messages depending on the TUI ownership contract.
- The tests do not assert the requested ordering, prefix override, cancellation behavior, request contents, or absence of deadlocks.

Do:
1. Define the engine-to-adapter request contract and implement it consistently:
   - `NewMessage` should contain the latest message being presented to the agent for this turn.
   - `Transcript` should contain the prior conversation context without duplicating `NewMessage`, unless the contract is intentionally documented otherwise.
   - Empty-conversation agent turns must have a documented prompt behavior.
   - Resume and non-resume/replay adapters must receive useful input in both first and later turns.
2. Monitor cancel signals during every agent turn. Cancel the per-turn context promptly without data races, and distinguish turn-cancel from full app quit according to the TUI/bridge contract.
3. On turn cancellation, stop or clear the UI stream exactly once, append exactly one `(turn cancelled)` transcript entry, display it once, and route control back to a user turn.
4. Avoid indefinite blocking on cancellation or adapter contract violations. The engine should drain what is safe to drain, but context cancellation must be able to unblock the turn.
5. Align agent response ownership with the TUI contract from task 37a so streamed agent output is persisted once and displayed once.
6. Make round-robin state explicit and deterministic after normal agent turns, user turns, prefix overrides, cancellations, and errors.
7. Treat unknown agent slots, missing adapters, empty turn orders, bridge errors, adapter errors, and context cancellation with clear errors and no goroutine leaks.
8. Replace the current loose engine tests with deterministic tests that wait for `Engine.Run` to return under a bounded context and assert all important side effects.

Do NOT:
- Do not implement `/plan`, `/quit`, or plan-file generation.
- Do not change the adapter public interface unless the plan is updated to account for the ripple.
- Do not hide adapter errors by converting them into assistant messages unless this is explicitly documented and tested.

Done when:
- Echo/fake adapters receive non-empty `NewMessage` when prior conversation exists, including resume-style later turns.
- Round-robin order and prefix override are asserted over multiple turns.
- Canceling a blocked/streaming agent turn cancels its context, records one cancellation note, returns to user input, and does not deadlock.
- Agent responses and cancellation notes are appended to session and displayed exactly once.
- `go test ./internal/orchestrator/...` passes without goroutine leaks or sleeps-as-assertions.
- `go test ./...`, `go vet ./...`, `go build ./...`, and `gofmt -l .` pass with no output from gofmt.

Files: internal/orchestrator/engine.go, internal/orchestrator/engine_test.go, internal/orchestrator/prefix.go, internal/orchestrator/prefix_test.go, internal/orchestrator/bridge.go, internal/agent/agent.go
