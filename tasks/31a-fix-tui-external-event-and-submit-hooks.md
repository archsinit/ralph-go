# Fix TUI orchestrator-facing API and submit contract

Goal: Make the Phase 4 TUI usable from Phase 5 orchestration code without reaching into package-private state or message types.

Context:
- `internal/tui` currently defines stream message types as unexported `msgStream*` structs.
- `Run()` creates a local model/program and returns only an error, so callers cannot retain a handle for `Program.Send` or `Model.StreamTokens`.
- `Model.StreamTokens` only works if `m.program` has already been set, but external packages cannot set that field.
- Pressing Enter appends a local user message and clears input, but does not emit the submitted text to an orchestrator callback/channel as required by task 30.
- Existing TUI tests are in package `tui`, so they can use internals and do not prove the external API is usable.

Do:
1. Define the public TUI contract needed by the orchestrator:
   - how user submits are delivered out of the TUI,
   - how finalized/history messages are delivered into the TUI,
   - how stream start/token/end events are delivered into the TUI,
   - who owns appending the submitted user message to the transcript view.
2. Add a public handle/API for running the Bubble Tea program while allowing external code to send stream and message events. Examples: exported event constructors plus a `Program.Send` handle, or a TUI handle with `StartStream`, `SendToken`, `EndStream`, and `AddMessage` methods.
3. Add a submit callback or channel configured by `NewModel`/options. Pressing Enter must send exactly the submitted text to this hook and clear the input.
4. Prevent duplicate user messages by documenting and testing whether the TUI appends submitted text immediately or waits for the orchestrator to echo/persist it back through the message-ingest API.
5. Support initial transcript/history injection so resumed sessions can render existing `session.Messages` without accessing private fields.
6. Add at least one external-package test (`package tui_test`) or equivalent compile-time usage check proving another internal package can use the public API without unexported types.
7. Update existing TUI tests so submit behavior asserts the outbound hook/channel, not just local mutation.

Do NOT:
- Do not implement the Phase 5 round-robin orchestrator.
- Do not parse agent prefixes or route messages between agents.
- Do not add session persistence side effects inside the TUI.

Done when:
- External code can run the TUI, receive user submissions, render initial/finalized messages, and feed live stream events using exported APIs only.
- The submit ownership contract is documented and covered by tests.
- `go test ./internal/tui/...` passes.
- `go test ./...`, `go vet ./...`, and `go build ./...` pass.

Files: internal/tui/tui.go, internal/tui/tui_test.go, internal/tui/tui_integration_test.go
