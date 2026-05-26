# Implement append and flush writer

Goal: Persist each message immediately and update agent session ids durably.

Context:
- Per-turn flush is the crash-safety guarantee.

Do:
1. Add (s *Session) Append(m Message) error: assign Seq, set TS if zero, append in memory, write the JSON line to transcript.jsonl, flush/Sync.
2. Add (s *Session) SetAgentSession(agent, id string) error: update map, rewrite agents.json atomically (temp + rename).
3. Open existing transcript.jsonl to set next Seq correctly.

Done when:
- go build ./... succeeds.
- After Append, the line exists on disk immediately.
- agents.json rewrite is atomic.

Files: internal/session/session.go
