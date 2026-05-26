# Implement round-robin turn engine

Goal: Drive turns over the configured order, blocking on user slots and invoking agents otherwise.

Context:
- Keep transport to the TUI behind the bridge interface so tests can supply a fake bridge.

Do:
1. In internal/orchestrator/engine.go define Engine holding config, session, adapters map, and a UI bridge interface.
2. UI bridge interface: methods to stream start/token/end and to request user input (returns a channel or blocks).
3. Loop over Plan.TurnOrder cyclically: if slot == "user", request input; else invoke that agent.
4. On agent turn: build agent.Request (system prompt from config, transcript from session, resume id), call Invoke, forward tokens to UI, append final message + SetAgentSession from Result.
5. On a user message with a prefix target, make the next turn that target instead of the cyclic next.

Done when:
- go build ./... succeeds.
- With echo adapters and scripted input, order is honored.
- Prefix override redirects the next turn.

Files: internal/orchestrator/engine.go
