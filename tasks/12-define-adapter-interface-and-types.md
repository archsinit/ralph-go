# Define Adapter interface and types

Goal: Define the agent adapter contract used by both plan and loop modes.

Context:
- Streaming model: tokens on one channel, a single terminal Result on another, both from Invoke.

Do:
1. In internal/agent/agent.go define Capabilities struct { SupportsResume bool }.
2. Define Token struct { Text string } (a streamed chunk).
3. Define Result struct { SessionID string; Err error } returned when a turn completes.
4. Define interface Adapter { Name() string; Capabilities() Capabilities; Invoke(ctx context.Context, req Request) (<-chan Token, <-chan Result) }.
5. Request struct: SystemPrompt string, Transcript []Message, NewMessage string, ResumeSessionID string.
6. Define Message struct { Author, Role, Text string }.
7. Document streaming contract: Token channel closes when done, then Result delivers one value.

Do NOT:
- Do not implement CLI adapters in this task.

Done when:
- go build ./... succeeds.
- Interface and types documented.
- No concrete impl yet beyond declarations.

Files: internal/agent/agent.go
