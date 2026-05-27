# Fix engine turn cancellation and adapter request contracts

Goal: Make the orchestrator engine satisfy the Phase 5.5 contracts that `/plan`, resume, and later loop mode depend on.

Context:
- The audit found that agent requests set `NewMessage` to the latest transcript message while also including that same message in `Transcript`, duplicating input for replay/resume adapters.
- Agent-turn cancellation currently returns `context.Canceled` to `Run`, which stops plan mode instead of returning control to a user turn.
- Cancellation state is shared across goroutines without a clear synchronization contract, and result draining can block forever if an adapter violates the stream/result contract.
- Empty-conversation agent turns, prefix overrides after cancellation, adapter errors, and `/quit` vs turn-cancel semantics are not specified tightly enough.

Do:
1. Define a single engine-to-adapter request builder and use it for normal agent turns and plan-generation turns:
   - `NewMessage` contains only the message being presented for this turn,
   - `Transcript` contains prior context and does not duplicate `NewMessage`,
   - resume IDs and system prompts are populated consistently,
   - empty-conversation agent turns have an explicit, tested behavior.
2. Introduce an internal turn-cancel result/sentinel distinct from app quit and parent context cancellation.
3. On agent turn cancel: cancel the per-turn context, end/clear the UI stream exactly once, append/display exactly one `(turn cancelled)` note if that is the chosen UX, and route control back to a user turn without exiting plan mode.
4. Ensure result/token draining cannot hang forever on cancellation or adapter contract violations; use bounded/context-aware waits and avoid goroutine leaks.
5. Remove data races by communicating cancellation state through channels, contexts, or synchronized fields.
6. Make round-robin state deterministic after normal turns, prefix overrides, cancellation, adapter errors, and `/plan`/`/quit` commands.
7. Keep adapter errors visible as errors; do not silently convert them into assistant messages unless that policy is documented and tested.

Do NOT:
- Do not change the public `agent.Adapter` interface unless the plan is updated for every impacted adapter.
- Do not implement artifact rendering here; use the plan package from Phase 5.6/6.5 task 49c.

Done when:
- Fake adapters observe non-duplicated `Transcript`/`NewMessage` inputs on first, resumed, prefixed, and plan-generation turns.
- Canceling a blocked or streaming agent turn returns to user input and does not exit the app.
- `/quit` exits cleanly and is distinguishable from turn cancel.
- Unknown agents, missing adapters, empty turn orders, bridge errors, adapter errors, and parent context cancellation return clear errors with no goroutine leaks.
- `go test -race ./internal/orchestrator/...`, `go test ./...`, `go vet ./...`, `go build ./...`, and `gofmt -l .` pass.

Files: internal/orchestrator/engine.go, internal/orchestrator/engine_test.go, internal/orchestrator/bridge.go, internal/orchestrator/prefix.go, internal/agent/agent.go
