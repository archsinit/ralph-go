# Wire loop subcommand end-to-end

Goal: Connect everything behind ralph-go loop with its TUI bridge.

Context:
- Integration point for Phases 8–13.

Do:
1. In cmd/loop.go: load+validate config, resolve project dir (cwd) and plan path (default ./plan.md, --plan flag), build executor+reviewer adapters, ntfy client, loggers, construct the loop Engine and the loop TUI, run.
2. Pass OutDir/LogDir/SessionDir from config.
3. Clean shutdown persists logs and leaves consistent state.

Done when:
- go build ./... succeeds.
- ./ralph-go loop runs against a plan.md using echo agents.
- All wiring present; no nil deref.

Files: cmd/loop.go, internal/loop/engine.go
