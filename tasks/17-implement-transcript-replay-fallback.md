# Implement transcript replay fallback

Goal: Provide a shared helper that renders transcript context into a prompt when resume is unavailable.

Context:
- Keeps token cost predictable and behavior consistent across CLIs without resume.

Do:
1. In internal/agent/replay.go add func renderReplay(req Request) string.
2. Format: system prompt, then each Message as "Author: Text", then the new message, clearly delimited.
3. Adapters with SupportsResume=false or empty ResumeSessionID use this to build stdin/prompt.

Done when:
- go build ./... succeeds.
- Output deterministic and readable.
- Used by at least one adapter path.

Files: internal/agent/replay.go
