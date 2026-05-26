# Define transcript and session model

Goal: Model the chatroom transcript and per-agent session metadata for persistence.

Context:
- Crash-safety achieved by appending to transcript.jsonl per turn (next task).

Do:
1. In internal/session/session.go define Message { Seq int; Author, Role, Text string; TS time.Time }.
2. Define Session { Dir string; Messages []Message; AgentSessions map[string]string } (agent name -> CLI session id).
3. Define constructor Open(dir string) (*Session, error) creating dir if absent and loading existing state if present.
4. Decide on-disk files: transcript.jsonl (one Message per line) and agents.json (the map).

Do NOT:
- Do not write append/flush logic yet beyond what Open needs to load.

Done when:
- go build ./... succeeds.
- Open on empty dir returns an empty usable Session.
- Types documented.

Files: internal/session/session.go
