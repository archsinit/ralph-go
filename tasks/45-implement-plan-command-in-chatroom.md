# Implement /plan command in chatroom

Goal: Add the /plan command that asks the designated agent for a JSON task list.

Context:
- Reuse the agent adapter; this is a special directed turn, not a normal round-robin turn.

Do:
1. In the orchestrator/TUI input handling, detect a line equal to '/plan'.
2. On /plan: send the configured Plan.PlanAgent a message instructing JSON-only output matching the TaskSpec schema, including the full transcript as context (resume or replay).
3. Embed the schema and rules in the agent's prompt (JSON only, no fences, required fields).
4. Collect the agent's full response text (not streamed into the normal transcript bubble necessarily; may show a 'generating plan...' status).

Done when:
- go build ./... succeeds.
- Typing /plan triggers a single plan-generation turn.
- The raw JSON response is captured for parsing.

Files: internal/orchestrator/engine.go, internal/tui/tui.go
