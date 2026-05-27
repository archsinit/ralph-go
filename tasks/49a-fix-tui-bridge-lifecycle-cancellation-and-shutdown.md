# Fix remaining TUI/bridge lifecycle, cancellation, and shutdown semantics

Goal: Close the remaining Phase 5.5 lifecycle gaps before plan generation is treated as ready.

Context:
- The audit found that `tui.Run` now returns a live handle, but `ctrl+c` always calls the cancel callback and never has a documented idle/second-press quit path.
- `cmd/plan.go` starts the engine while the TUI is alive, but it waits only for the engine; if the TUI exits first, the engine can remain blocked on user input.
- `TUIBridge.RequestCancelSignal` ignores its context, submit/cancel callbacks can silently drop events when buffers fill, and stale cancel signals can affect later turns.
- The TUI still resets `userScrolled` on `StreamStart`, does not forward mouse wheel events to the viewport, uses byte-count wrapping, and can render taller than very small terminals.

Do:
1. Define and implement explicit TUI cancel/quit semantics:
   - first `ctrl+c` during an active agent turn/stream sends a turn-cancel signal without quitting,
   - `ctrl+c` while idle, or a documented second `ctrl+c`, quits the chatroom,
   - `/quit` and app quit paths share the same clean shutdown behavior.
2. Make the bridge event paths reliable under normal use: honor context cancellation, avoid silent submit/cancel drops, and drain or scope stale cancel signals so a previous cancel cannot cancel a later turn.
3. Rework `cmd/plan.go` lifecycle waiting so TUI exit cancels the engine, fatal engine errors quit the TUI and are returned to the CLI, and `/quit` exits successfully.
4. Preserve viewport position when the user has scrolled up and a stream starts or tokens arrive; resume autoscroll only when the viewport is back at the bottom.
5. Forward viewport mouse events where Bubble Tea/Bubbles supports them.
6. Fix small-terminal layout so the viewport/status/input never obviously exceed the reported terminal height, and replace byte-count wrapping/slicing with rune/display-width-aware wrapping that handles wide runes and long tokens safely.
7. Update `cmd/tui-test` to demonstrate live handle startup, submit callback, fake streaming, turn cancel, idle/second-press quit, scrolling, and clean shutdown.

Do NOT:
- Do not change agent adapter behavior in this task except where required by bridge lifecycle wiring.
- Do not write plan artifacts from the TUI package.

Done when:
- External code can start the TUI, use the handle immediately, receive submits/cancels, quit from the UI, and wait for shutdown without hangs.
- TUI quit cancels a blocked engine in plan mode.
- Turn cancel and app quit are distinct and documented in tests.
- Scroll preservation, mouse scroll, small terminal layout, and wide-rune wrapping have behavior-level tests.
- `gofmt -l .`, `go test ./internal/tui/...`, `go test ./internal/orchestrator/...`, `go test ./...`, `go vet ./...`, and `go build ./...` pass.

Files: internal/tui/tui.go, internal/tui/tui_test.go, internal/tui/tui_external_test.go, internal/tui/tui_integration_test.go, internal/orchestrator/bridge.go, internal/orchestrator/bridge_test.go, cmd/plan.go, cmd/tui-test/main.go
