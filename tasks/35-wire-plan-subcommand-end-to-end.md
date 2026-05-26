# Wire plan subcommand end-to-end

Goal: Connect config, session, adapters, engine, and TUI behind ralph-go plan.

Context:
- This is the integration point for Phases 1–5.

Do:
1. In cmd/plan.go: load+validate config, build adapters from config.Agents via agent.New, Open session dir (config.Paths.SessionDir, support --session to resume), construct the TUI model and an orchestrator Engine bridged to it, run.
2. Implement the bridge that maps engine stream calls to tui stream msgs and engine input requests to tui input submissions.
3. Handle clean shutdown persisting session.

Done when:
- go build ./... succeeds.
- ./ralph-go plan starts the chatroom using echo or real agents per config.
- Conversation persists to the session dir.

Files: cmd/plan.go, internal/orchestrator/engine.go, internal/tui/tui.go
