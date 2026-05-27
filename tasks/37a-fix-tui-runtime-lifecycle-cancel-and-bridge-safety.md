# Fix TUI runtime lifecycle, cancel signaling, and bridge safety

Goal: Make the plan TUI controllable by the orchestrator while it is running, and make cancel/bridge behavior match the Phase 5 contract.

Context:
- `tui.Run` currently blocks in `Program.Run()` and returns the `Handle` only after the Bubble Tea program exits. `cmd/plan.go` therefore cannot start the engine while the TUI is alive, and `cmd/tui-test` cannot feed fake streams until after the UI quits.
- Pressing `ctrl+c` in the TUI always returns `tea.Quit`; task 34 requires canceling the current agent turn without exiting the app.
- `TUIBridge.SendCancelSignal` is not wired to any TUI key path, uses an unbuffered channel with a non-blocking send, and can silently drop cancellation.
- Bridge methods call `b.handle.*` without nil/race protection; startup ordering changes could otherwise panic.
- Stream/message ownership is currently ambiguous: `StreamEnd` finalizes the in-progress message in the TUI, while the engine also calls `AddMessage` after persisting the same agent response.
- Some Phase 4.5 rendering promises are still fragile: stream start forces autoscroll by clearing `userScrolled`, mouse scroll is not forwarded, very small terminal layout can exceed the terminal, and wrapping still uses byte counts.

Do:
1. Introduce a TUI lifecycle API that returns a usable handle before the program exits. One acceptable shape is `Start(opts ...Option) (*Handle, error)` or `NewProgram/RunAsync`, plus `Handle.Wait()`/`Err()`/`Quit()`; keep a blocking compatibility wrapper only if it cannot be confused with the non-blocking API.
2. Add an explicit cancel callback/channel option for the TUI. During an active agent stream/turn, the first `ctrl+c` should emit cancel to the orchestrator instead of quitting. Define and test the behavior for a second `ctrl+c` and for `ctrl+c` while idle.
3. Wire `TUIBridge` to the TUI submit and cancel hooks. Make cancellation buffered or otherwise reliable under normal use, honor context cancellation in bridge wait paths, and avoid silent drops unless documented and tested.
4. Make all bridge stream/message methods safe during startup/shutdown: protect handle access or return clear errors instead of nil panics.
5. Decide and document TUI stream ownership. Either `StreamEnd` finalizes the streamed message and the engine must not call `AddMessage` for that response, or `StreamEnd` only clears the in-progress stream and persisted messages are echoed back through `AddMessage`. Implement one contract and add duplicate-prevention tests.
6. Preserve a user-scrolled viewport when a stream starts or tokens arrive; autoscroll should resume only after the user returns to the bottom. Forward viewport mouse events where supported.
7. Tighten small-terminal layout and wrapping enough that the viewport/status/input fit the terminal and rendered lines do not obviously exceed display width for long tokens or wide runes.
8. Update `cmd/tui-test` so it demonstrates the non-blocking handle, submit callback, fake streaming, cancellation, scrolling, and shutdown.

Do NOT:
- Do not implement the Phase 5.6 plan-file generator or Phase 6 `/plan` command.
- Do not add session persistence side effects inside the TUI.
- Do not make the TUI know about concrete agent adapters.

Done when:
- External code can start the TUI, immediately retain a live handle, receive submits/cancels, send messages, send stream start/token/end events, wait for shutdown, and quit cleanly.
- A TUI `ctrl+c` can cancel a running turn without exiting the app.
- Agent responses are displayed exactly once.
- Autoscroll, resize, wrapping, and scrolling regressions are covered by tests.
- `go test ./internal/tui/...` passes.
- `go test ./internal/orchestrator/...`, `go test ./...`, `go vet ./...`, `go build ./...`, and `gofmt -l .` pass with no output from gofmt.

Files: internal/tui/tui.go, internal/tui/tui_test.go, internal/tui/tui_integration_test.go, internal/tui/tui_external_test.go, internal/orchestrator/bridge.go, cmd/tui-test/main.go
