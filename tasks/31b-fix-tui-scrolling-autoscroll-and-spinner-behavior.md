# Fix TUI viewport, spinner, rendering, and harness behavior

Goal: Make the Phase 4 TUI behavior match the interaction and rendering promises before Phase 5 wires real orchestration into it.

Context:
- `Model.Update` marks `userScrolled`, but it does not forward scroll keys/mouse events to `viewport.Update`, so the transcript viewport does not actually scroll.
- Autoscroll state is inferred from key names before viewport movement, rather than from the viewport's actual bottom position after movement.
- Spinner animation is not started with the normal `spinner.Tick` command flow, so the waiting indicator may not animate during streams.
- Author styling only distinguishes hard-coded names (`user`, `claude`, `codex`); other configured agents fall back to the same style.
- Wrapping uses a simple byte/word count and does not split very long words; rendered lines can exceed the viewport width.
- The temporary `cmd/tui-test` harness currently starts a blank TUI and does not demonstrate fake streams, status changes, or scrolling.

Do:
1. Forward viewport-relevant key and mouse messages to `viewport.Update`.
2. Track `userScrolled` from actual viewport position: when the user scrolls away from the bottom, new tokens/messages should not force-scroll; when the user returns to bottom, autoscroll should resume.
3. Start spinner animation with `spinner.Tick` when a stream starts, continue ticking only while streaming, and stop after stream end.
4. Keep status behavior deterministic: waiting status while streaming, `your turn` when input is expected, and a documented idle state if distinct from input.
5. Clamp layout dimensions for small terminal sizes so viewport height never becomes invalid and viewport + separator + status + input fit without clipping.
6. Improve author styling so every configured/seen author gets a stable distinguishable style, not a shared unknown-agent color.
7. Improve wrapping to respect terminal display width and handle long unbroken tokens reasonably.
8. Expand the temporary harness to feed fake streams (or echo-adapter-style streams), so manual review can verify typing, live streaming, scrolling, resize, status, and spinner behavior.
9. Add unit tests for scrolling up, staying scrolled up while tokens arrive, resuming autoscroll at bottom, spinner tick command behavior, stable author style assignment, and long-message wrapping.

Do NOT:
- Do not implement the orchestrator turn engine.
- Do not connect real agent CLIs here.
- Do not make brittle tests that depend on exact ANSI escape sequences unless unavoidable.

Done when:
- Manual scrolling actually moves the transcript viewport.
- Autoscroll, resize, wrapping, author styling, status, and spinner behavior match the Phase 4 task descriptions.
- The harness can visibly demonstrate live token accumulation and status/spinner changes.
- `go test ./internal/tui/...` passes.
- `go test ./...`, `go vet ./...`, and `go build ./...` pass.

Files: internal/tui/tui.go, internal/tui/tui_test.go, internal/tui/tui_integration_test.go, cmd/tui-test/main.go
