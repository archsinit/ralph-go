# Implement load and resume

Goal: Rebuild a Session from disk so plan mode can resume mid-conversation.

Context:
- Resume IDs feed agent.Request.ResumeSessionID.

Do:
1. Ensure Open reads transcript.jsonl into Messages in order and agents.json into AgentSessions.
2. Handle missing/partial last line gracefully (ignore a trailing corrupt line, log nothing fatal).
3. Expose (s *Session) ResumeID(agent string) string.

Done when:
- go build ./... succeeds.
- Open on a populated dir restores messages and agent ids.
- A truncated final line does not crash load.

Files: internal/session/session.go
