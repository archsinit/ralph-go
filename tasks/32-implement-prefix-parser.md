# Implement prefix parser

Goal: Parse an optional leading 'name:' on a user message to target one agent next turn.

Context:
- Used to override round-robin for the immediately following turn.

Do:
1. In internal/orchestrator/prefix.go add func ParsePrefix(input string, agents []string) (target string, body string).
2. If input starts with "<agentname>:" where agentname is a configured agent, return that target and the trimmed remainder.
3. Otherwise target empty, body unchanged.
4. Case-insensitive match on agent name; require the colon.

Done when:
- go build ./... succeeds.
- "codex: do x" -> target codex, body "do x".
- "hello" -> target empty.
- Unknown prefix like "foo: x" -> target empty, body unchanged.

Files: internal/orchestrator/prefix.go
