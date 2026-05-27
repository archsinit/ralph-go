# Fix plan command wiring and turn-order model ownership

Goal: Make `ralph-go plan` actually run an interactive plan chatroom end-to-end and remove structural blockers before the plan-format package work starts.

Context:
- `cmd/plan.go` calls `tui.Run(...)` before `engine.Run(...)`. Because `tui.Run` blocks until the TUI exits, the engine is not running while the user is in the chatroom.
- `cmd/root.go` and `cmd/plan.go` both call `rootCmd.AddCommand(planCmd)`, so `ralph-go --help` shows the `plan` command twice.
- Phase 5 introduced `internal/plan.Plan` as a turn-order wrapper. Phase 5.6 tasks are intended to define `internal/plan.Plan` for generated task specs, so the current type/package ownership will collide with the next phase.
- Plan-mode startup/shutdown currently lacks clear error propagation between the TUI, bridge, session, and engine.

Do:
1. Register `planCmd` exactly once. Verify `go run . --help` lists one `plan` command.
2. Wire `cmd/plan.go` to the non-blocking TUI lifecycle from task 37a:
   - load and validate config,
   - open or resume the session directory,
   - render existing session messages as initial TUI history,
   - build adapters,
   - create the bridge,
   - start the TUI and engine concurrently,
   - propagate TUI quit to engine cancellation,
   - propagate fatal engine errors to TUI shutdown,
   - return a useful error to the CLI.
3. Ensure submitted user messages are persisted before being echoed to the TUI, and ensure resumed transcript messages are not duplicated on startup.
4. Resolve the turn-order model ownership before Phase 5.6:
   - Prefer moving/renaming the Phase 5 turn-order wrapper into `internal/orchestrator` (for example `TurnPlan`, `TurnOrder`, or accepting `[]string` directly), or otherwise rename it so `internal/plan.Plan` remains available for the task-spec plan package.
   - Update imports and tests accordingly.
5. Preserve `--session` resume support and `config.Paths.SessionDir` behavior. If multiple sessions under a base directory are desired later, leave that to a future task rather than changing semantics silently.
6. Add a lightweight plan-mode smoke path using echo/fake adapters where possible without requiring a real terminal. If a full TUI subprocess test is too brittle, factor wiring enough to unit-test the dependencies.

Do NOT:
- Do not implement Phase 5.6 task file generation.
- Do not implement Phase 6 `/plan` or `/quit` chat commands.
- Do not require real Claude/Codex/Gemini/OpenCode CLIs for automated tests.

Done when:
- `ralph-go --help` contains exactly one `plan` command.
- `ralph-go plan --config <echo config> --session <dir>` starts a live chatroom whose submits reach the engine while the TUI remains open.
- Quitting the TUI cancels the engine and exits cleanly.
- Existing session transcript messages render on resume without duplicate entries.
- The future `internal/plan.Plan` task-spec type can be added without conflicting with Phase 5 turn-order code.
- `go test ./...`, `go vet ./...`, `go build ./...`, and `gofmt -l .` pass with no output from gofmt.

Files: cmd/root.go, cmd/plan.go, internal/orchestrator/engine.go, internal/orchestrator/bridge.go, internal/plan/plan.go, internal/plan/doc.go
